package services

import (
	"context"
	"fmt"
	"go_echo/collections"
	"go_echo/common"
	"go_echo/configs"
	"go_echo/helper"
	"log"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserProjectedData struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	PredictionID primitive.ObjectID `bson:"prediction_id" json:"prediction_id"`
	MatchID      int                `bson:"match_id" json:"match_id"`
	Sport        int                `bson:"sport" json:"sport"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
}

type CancelUserData struct {
	SlotFee             float64 `bson:"slot_fee" json:"slot_fee"`
	TotalSlot           int     `bson:"total_slot" json:"total_slot"`
	OptionID            int     `bson:"option_id" json:"option_id"`
	CancelledSlotNumber int     `bson:"cancelled_slot_number" json:"cancelled_slot_number"`
}

func CancelAllUnMatchedTrade(data common.CancelAllUnMatchedBody) (error, common.Response) {
	predictionID, err := primitive.ObjectIDFromHex(data.PredictionID)
	userID, err := primitive.ObjectIDFromHex(data.UserId)
	filter := bson.D{
		{Key: "prediction_id", Value: predictionID},
		{Key: "is_pred_cancel", Value: 0},
		{Key: "is_slot_cancel", Value: 0},
		{Key: "user_id", Value: userID},
	}
	projection := bson.D{
		{Key: "prediction_id", Value: 1},
		{Key: "match_id", Value: 1},
		{Key: "sport", Value: 1},
		{Key: "_id", Value: 1},
		{Key: "user_id", Value: 1},
	}

	options := options.Find().SetProjection(projection)

	cursor, err := collections.TRADE_JOINED_COLLECTION.Find(context.TODO(), filter, options)
	if err != nil {
		log.Panic(err)
	}
	defer cursor.Close(context.TODO())

	var trades []UserProjectedData // Assuming `models` is the package where you defined the `PredictionTradeJoined` struct

	if err := cursor.All(context.TODO(), &trades); err != nil {
		log.Panic(err)
	}
	if len(trades) == 0 {
		response := common.Response{
			Status:  false,
			Message: "No Slots For Cancel",
			Data:    []string{},
		}
		return nil, response
	}
	var updateUsersCh = make(chan mongo.WriteModel, len(trades))
	var updateTradeJoinedCh = make(chan mongo.WriteModel, len(trades))
	var userRefund = make(chan common.TradeTransaction, len(trades))
	var wg sync.WaitGroup
	for _, element := range trades {
		wg.Add(1)
		go cancelAllSlotsForTheSinglePrediction(element, &wg, updateUsersCh, updateTradeJoinedCh, userRefund)
	}
	go func() {
		wg.Wait()
		close(updateTradeJoinedCh)
		close(updateUsersCh)
		close(userRefund)
	}()
	if err := cursor.Err(); err != nil {
		log.Panic(err)
	}
	var updateUsers []mongo.WriteModel
	var updateTradeJoined []mongo.WriteModel
	var tradeTransaction common.TradeTransaction
	for u := range updateUsersCh {
		updateUsers = append(updateUsers, u)
	}
	for r := range userRefund {
		tradeTransaction.InCash += r.InCash
		tradeTransaction.PredictionID = r.PredictionID
		tradeTransaction.UserId = r.UserId
		tradeTransaction.InWinning += r.InWinning
		tradeTransaction.SportType = r.SportType
		tradeTransaction.MatchId = r.MatchId
		tradeTransaction.PredictionType = r.PredictionType
		tradeTransaction.Type = "cancel_all_trade"
	}
	helper.CreateTradeTransaction(tradeTransaction.UserId, tradeTransaction)
	for t := range updateTradeJoinedCh {
		updateTradeJoined = append(updateTradeJoined, t)
	}

	if len(updateTradeJoined) > 0 {
		_, err = collections.TRADE_JOINED_COLLECTION.BulkWrite(context.TODO(), updateTradeJoined)
		if err != nil {
			return err, common.Response{}
		}
	}

	if len(updateUsers) > 0 {
		_, err = collections.USERS.BulkWrite(context.TODO(), updateUsers)
		if err != nil {
			return err, common.Response{}
		}
	}
	tradeTransaction.InCash = 0
	tradeTransaction.InWinning = 0

	helper.CreateTradeTransaction(tradeTransaction.UserId, tradeTransaction)
	response := common.Response{
		Status:  true,
		Message: "All unmatched slots are successfully cancelled",
		Data:    []string{},
	}
	return nil, response
}

func cancelAllSlotsForTheSinglePrediction(data UserProjectedData, wg *sync.WaitGroup, updateUsersCh chan<- mongo.WriteModel, updateTradeJoinedCh chan<- mongo.WriteModel, userRefund chan<- common.TradeTransaction) {
	defer wg.Done()
	var recordKey = fmt.Sprintf("%s%d:%d:%s:%s:%s", common.TRADE_CONSTANT.JOINED_PREDICTION_TRADE, data.MatchID, data.Sport, data.PredictionID.Hex(), data.UserID.Hex(), data.ID.Hex())
	recordKeyData, err := configs.GetHashKeyValues(recordKey)
	if err != nil {
		log.Panic(err)
	}
	if len(recordKeyData) > 0 {
		if _, ok := recordKeyData["_id"]; ok {
			totalSlots, _ := strconv.Atoi(recordKeyData["total_slot"])
			isSlotCancel, _ := strconv.Atoi(recordKeyData["is_slot_cancel"])
			totalMatched, _ := strconv.Atoi(recordKeyData["total_matched"])
			totalCashDeduct, _ := strconv.ParseFloat(recordKeyData["total_cash"], 64)
			totalWiningDeduct, _ := strconv.ParseFloat(recordKeyData["total_wining"], 64)
			slotFee, _ := strconv.ParseFloat(recordKeyData["slot_fee"], 64)

			previousTotalCancelled, _ := strconv.Atoi(recordKeyData["cancelled_slot_number"])
			unMatchedSlots := totalSlots - (totalMatched + previousTotalCancelled)
			refundAmount := float64(unMatchedSlots) * slotFee
			cashRatio := totalCashDeduct / float64(totalSlots)
			winningRatio := totalWiningDeduct / float64(totalSlots)

			refundCash := cashRatio * float64(unMatchedSlots)
			refundWin := winningRatio * float64(unMatchedSlots)
			refundCash = refundAmount * (cashRatio / (cashRatio + winningRatio))
			refundWin = refundAmount - refundCash
			if refundAmount > 0 && unMatchedSlots > 0 && isSlotCancel == 0 {
				preRefund, _ := strconv.ParseFloat(recordKeyData["refund_amount"], 64)
				totalRefund := preRefund + refundAmount
				optionId, _ := strconv.Atoi(recordKeyData["option_id"])
				entryFee, _ := strconv.ParseFloat(recordKeyData["slot_fee"], 64)
				keyNames := helper.KeyName(entryFee)
				removeCountKey := fmt.Sprintf("%s%d:%s:%d:%s", common.TRADE_CONSTANT.TRADE_UNMATCHED, data.MatchID, data.PredictionID.Hex(), optionId, keyNames)
				removeCountValue := fmt.Sprintf("%s-%s", data.UserID.Hex(), data.ID.Hex())
				removeCount, err := configs.RemoveListElement(removeCountKey, -int64(unMatchedSlots), removeCountValue)
				if err != nil {
					log.Panic(err)
				}
				if removeCount > 0 {
					fields := map[string]interface{}{
						"is_slot_cancel":        1,
						"refund_amount":         totalRefund,
						"cancelled_slot_number": unMatchedSlots,
						"refund_cash_amount":    refundCash,
						"refund_win_amount":     refundWin,
					}
					// update := bson.M{
					// 	"$set": fields,
					// }
					err := configs.Hmset(recordKey, fields)
					if err != nil {
						log.Panic(err)
					}
					// _, err = collections.TRADE_JOINED_COLLECTION.UpdateOne(context.TODO(), bson.M{"_id": data.ID}, update)
					// if err != nil {
					// 	log.Panic(err)
					// }
					keyForTotalJoined := fmt.Sprintf("%s%d:%s:%s", common.TRADE_CONSTANT.TOTAL_TRADE_JOINED, data.MatchID, data.PredictionID.Hex(), data.UserID.Hex())
					configs.IncrementBY(keyForTotalJoined, -float64(unMatchedSlots))
					updateTradeJoinedModel := mongo.NewUpdateOneModel()
					updateTradeJoinedModel.SetFilter(bson.M{"_id": data.ID})
					updateTradeJoinedModel.SetUpdate(bson.D{
						{Key: "$set", Value: fields},
					})
					updateTradeJoinedModel.SetUpsert(false)
					// into channel
					updateTradeJoinedCh <- updateTradeJoinedModel
					updateUserModel := mongo.NewUpdateOneModel()
					updateUserModel.SetFilter(bson.M{"_id": data.UserID})
					updateUserModel.SetUpdate(bson.D{
						{Key: "$inc", Value: bson.D{
							{Key: "cash_balance", Value: refundCash},
							{Key: "winning_balance", Value: refundWin},
						}},
					})
					userRefund <- common.TradeTransaction{
						UserId:         data.UserID,
						InWinning:      refundWin,
						InCash:         refundCash,
						CashBalance:    0,
						WinningBalance: 0,
						PredictionType: "trading",
						MatchId:        data.MatchID,
						SportType:      data.Sport,
						PredictionID:   data.PredictionID,
					}
					updateUserModel.SetUpsert(false)
					// into channel
					updateUsersCh <- updateUserModel
					// updates := bson.D{
					// 	{Key: "$inc", Value: bson.D{
					// 		{Key: "cash_balance", Value: refundCash},
					// 		{Key: "winning_balance", Value: refundWin},
					// 	}},
					// }
				}
			}
		}
	}
}
