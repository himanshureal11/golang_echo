package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_echo/collections"
	"go_echo/common"
	"go_echo/configs"
	"go_echo/helper"
	"go_echo/models"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaleTrade(data common.RequestSaleBody) (error, common.Response) {
	var response common.Response = common.Response{
		Status:  false,
		Message: "No Prediction Found",
		Data:    []string{},
	}
	preKey := common.GetPredictionKey(common.TRADE_CONSTANT.MATCH_TRADE_PREDICTION, data.MatchID, data.Sport, data.PredictionID)
	preKeyResult, err := configs.GetHashKeyValues(preKey)

	if err != nil {
		return err, response
	}
	if len(preKeyResult) > 0 {
		if _, ok := preKeyResult["_id"]; ok {
			putOnSaleKey := fmt.Sprintf("%s%d:%s:%s:%d:%.1f", common.TRADE_CONSTANT.PREDICTION_TRADE_PREDEFINED_RATE, data.MatchID, data.PredictionID, data.RecordID, data.OptionID, data.SaleFee)
			joinPredKey := common.GetJoinedTradeKey(common.TRADE_CONSTANT.JOINED_PREDICTION_TRADE, data.MatchID, data.Sport, data.PredictionID, data.UserID, data.RecordID)
			joinSaleTradeKey := fmt.Sprintf("%s%s", common.SLOTS_ON_SALE, data.RecordID)
			joinPredKeyResult, err := configs.GetStringValue(joinSaleTradeKey)
			var userTradeSaleData []models.SaleTrade
			var saleTradeByUser = models.SaleTrade{
				SaleSlots: data.SaleSlots,
				SaleFee:   data.SaleFee,
				SoldSlots: 0,
			}
			var found bool
			if len(joinPredKeyResult) > 0 {
				err = json.Unmarshal([]byte(joinPredKeyResult), &userTradeSaleData)
				if err != nil {
					return err, response
				}
				for i, v := range userTradeSaleData {
					if v.SaleFee == data.SaleFee {
						// If SaleFee matches, update SaleSlots
						userTradeSaleData[i].SaleSlots += data.SaleSlots
						found = true
						break
					}
				}
			}
			if !found {
				userTradeSaleData = append(userTradeSaleData, saleTradeByUser)
			}
			err = updateSaleSlots(data, userTradeSaleData, joinPredKey, joinSaleTradeKey, putOnSaleKey)
			if err != nil {
				response.Message = "You Are not Allowed to sale the slots for this trade"
				return err, response
			}
			response.Status = true
			response.Message = "Trade successfully put on sale"
			return nil, response
		} else {
			return nil, response
		}
	} else {
		return nil, response
	}
}

func CancelSale(data common.CancelSaleRequestData) (error, common.Response) {
	var response common.Response = common.Response{
		Status:  false,
		Message: "No Prediction Found",
		Data:    []string{},
	}
	joinPredKey := common.GetJoinedTradeKey(common.TRADE_CONSTANT.JOINED_PREDICTION_TRADE, data.MatchID, data.Sport, data.PredictionID, data.UserID, data.RecordID)
	res, err := configs.HashGetByKeyField(joinPredKey, "slot_fee")
	if err != nil {
		return err, response
	}
	joinSaleTradeKey := fmt.Sprintf("%s%s", common.SLOTS_ON_SALE, data.RecordID)
	joinSaleTradeKeyResult, err := configs.GetStringValue(joinSaleTradeKey)
	if err != nil {
		response.Message = "Invalid key"
		return err, response
	}
	var userTradeSaleData []models.SaleTrade
	if err := json.Unmarshal([]byte(joinSaleTradeKeyResult), &userTradeSaleData); err != nil {
		response.Message = "Invalid Data"
		return err, response
	}
	keyName := helper.KeyName(float64(data.SaleFee))
	key := fmt.Sprintf("%s%d:%s:%d:%s", common.TRADE_ON_SALE, data.MatchID, data.PredictionID, data.OptionID, keyName)
	removeInKeyData := fmt.Sprintf("%s-%s-%.1f-%s", data.UserID, data.RecordID, data.SaleFee, res)
	configs.RemoveListElement(key, int64(data.CancelSlots), removeInKeyData)

	for i, v := range userTradeSaleData {
		if v.SaleFee == data.SaleFee {
			if userTradeSaleData[i].SaleSlots == 0 {
				response.Message = "Your All slots are sold"
				return errors.New("All Slots Are Matched"), response
			} else {
				data.CancelSlots = userTradeSaleData[i].SaleSlots
				userTradeSaleData[i].SaleSlots = 0
				configs.HashIncrBy(joinPredKey, "slots_on_sale", -float64(data.CancelSlots))
			}
			break
		}
	}
	err = saveSaleTradeDataInRedis(joinSaleTradeKey, userTradeSaleData)
	if err != nil {
		response.Message = "You Are Not Able To Cancel The Sale"
		return err, response
	}
	err = updateDbForSale(userTradeSaleData, data.RecordID, -data.CancelSlots)
	if err != nil {
		response.Message = "You Are Not Able To Cancel The Sale"
		return err, response
	}
	response.Message = "Sale Cancel Successfully"
	response.Status = true
	return nil, response
}

func SaleOnDifferentPrice(data common.SaleOnDifferentPrice) (error, common.Response) {
	var response common.Response = common.Response{
		Status:  false,
		Message: "No Prediction Found",
		Data:    []string{},
	}
	var cancelSaleBody = common.CancelSaleRequestData{
		RecordID:     data.RecordID,
		SaleFee:      data.OldSaleFee,
		PredictionID: data.PredictionID,
		MatchID:      data.MatchID,
		OptionID:     data.OptionID,
		Sport:        data.Sport,
		CancelSlots:  data.Slots,
		UserID:       data.UserID,
	}
	err, _ := CancelSale(cancelSaleBody)
	if err != nil {
		response.Message = "You Are not able to cancel sale trade"
		return nil, response
	}
	var puttingSaleSlot = common.RequestSaleBody{
		RecordID:     data.RecordID,
		SaleFee:      data.SaleFee,
		PredictionID: data.PredictionID,
		MatchID:      data.MatchID,
		OptionID:     data.OptionID,
		Sport:        data.Sport,
		UserID:       data.UserID,
		SaleSlots:    data.Slots,
	}
	fmt.Println(cancelSaleBody, puttingSaleSlot)
	return nil, response
}

// update in db

func updateDbForSale(saleSlotArray []models.SaleTrade, RecordID string, saleSlots int) error {
	RecordId, err := primitive.ObjectIDFromHex(RecordID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": RecordId}
	update := bson.M{
		"$inc": bson.M{"slots_on_sale": saleSlots},
		"$set": bson.M{"sale_trade": saleSlotArray},
	}
	opts := options.Update().SetUpsert(true)
	_, err = collections.TRADE_JOINED_COLLECTION.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// update in redis

func saveSaleTradeDataInRedis(key string, saleTrade []models.SaleTrade) error {
	jsonString, err := json.Marshal(saleTrade)
	if err != nil {
		return err
	}
	expiration := 30 * 24 * time.Hour
	err = configs.SetStringValue(key, string(jsonString), expiration)
	if err != nil {
		return err
	}
	return nil
}

// pushing key in redis list for sell

func pushKeyInTradeOnSale(key string, value string, wg *sync.WaitGroup) {
	defer wg.Done()
	configs.Rpush(key, value)
}

// updating data for sale slots

func updateSaleSlots(data common.RequestSaleBody, saleSlotArray []models.SaleTrade, joinedUserKey string, joinSaleTradeKey string, putOnSaleKey string) error {
	configs.HashIncrBy(joinedUserKey, "slots_on_sale", float64(data.SaleSlots))
	res, err := configs.HashGetByKeyField(joinedUserKey, "slot_fee")
	if err != nil {
		return err
	}
	configs.Hmset(joinedUserKey, map[string]interface{}{"on_sale": "true"})
	configs.Hmset(putOnSaleKey, map[string]interface{}{"is_put_on_sale": "1"})
	keyName := helper.KeyName(float64(data.SaleFee))
	key := fmt.Sprintf("%s%d:%s:%d:%s", common.TRADE_ON_SALE, data.MatchID, data.PredictionID, data.OptionID, keyName)
	pushInKeyData := fmt.Sprintf("%s-%s-%.1f-%s", data.UserID, data.RecordID, data.SaleFee, res)
	var keySync sync.WaitGroup
	for i := 0; i < data.SaleSlots; i++ {
		keySync.Add(1)
		go pushKeyInTradeOnSale(key, pushInKeyData, &keySync)
	}
	keySync.Wait()
	tradeOnSaleKey := fmt.Sprintf("%s%d:%s:%s:%s", common.TRADE_ON_SALE_USER_META, data.MatchID, data.PredictionID, data.RecordID, data.UserID)
	tradeOnSalePushElement := fmt.Sprintf("%s-%s-%.1f-%d", data.UserID, data.RecordID, data.SaleFee, data.OptionID)
	configs.Rpush(tradeOnSaleKey, tradeOnSalePushElement)
	configs.SetWithExpirationDays(tradeOnSaleKey, 40)
	configs.SetWithExpirationDays(key, 40)
	err = updateDbForSale(saleSlotArray, data.RecordID, data.SaleSlots)
	if err != nil {
		return err
	}
	err = saveSaleTradeDataInRedis(joinSaleTradeKey, saleSlotArray)
	if err != nil {
		return err
	}
	return nil
}
