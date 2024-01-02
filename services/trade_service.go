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
		{"user_id", userID},
	}
	projection := bson.D{
		{Key: "prediction_id", Value: 1},
		{Key: "match_id", Value: 1},
		{Key: "sport", Value: 1},
		{Key: "_id", Value: 1},
		{"user_id", 1},
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
	var wg sync.WaitGroup
	for _, element := range trades {
		wg.Add(1)
		go cancelAllSlotsForTheSinglePrediction(element, &wg)
	}
	wg.Wait()
	if err := cursor.Err(); err != nil {
		log.Panic(err)
	}
	response := common.Response{
		Status:  true,
		Message: "All unmatched slots are successfully cancelled",
		Data:    []string{},
	}
	return nil, response
}

func cancelAllSlotsForTheSinglePrediction(data UserProjectedData, wg *sync.WaitGroup) {
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
					}
					update := bson.M{
						"$set": fields,
					}
					err := configs.Hmset(recordKey, fields)
					if err != nil {
						log.Panic(err)
					}
					_, err = collections.TRADE_JOINED_COLLECTION.UpdateOne(context.TODO(), bson.M{"_id": data.ID}, update)
					if err != nil {
						log.Panic(err)
					}
					keyForTotalJoined := fmt.Sprintf("%s%d:%s:%s", common.TRADE_CONSTANT.TOTAL_TRADE_JOINED, data.MatchID, data.PredictionID.Hex(), data.UserID.Hex())
					configs.IncrementBY(keyForTotalJoined, -float64(unMatchedSlots))
					updates := bson.D{
						{Key: "$inc", Value: bson.D{
							{Key: "cash_balance", Value: refundCash},
							{Key: "winning_balance", Value: refundWin},
						}},
					}
					collections.USERS.UpdateOne(context.TODO(), bson.M{"_id": data.UserID}, updates)
				}
			}
		}
	}
}
