package services

import (
	"context"
	"encoding/json"
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

var response common.Response = common.Response{
	Status:  false,
	Message: "No Prediction Found",
	Data:    []string{},
}

func SaleTrade(data common.RequestSaleBody) (error, common.Response) {
	preKey := common.GetPredictionKey(common.TRADE_CONSTANT.MATCH_TRADE_PREDICTION, data.MatchID, data.Sport, data.PredictionID)
	preKeyResult, err := configs.GetHashKeyValues(preKey)
	if err != nil {
		return err, response
	}
	if len(preKeyResult) > 0 {
		if _, ok := preKeyResult["_id"]; ok {
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
			updateSaleSlots(data, userTradeSaleData, joinPredKey, joinSaleTradeKey)
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

func pushKeyInTradeOnSale(key string, value string, wg *sync.WaitGroup) {
	defer wg.Done()
	configs.Rpush(key, value)
}

func updateSaleSlots(data common.RequestSaleBody, saleSlotArray []models.SaleTrade, joinedUserKey string, joinSaleTradeKey string) {
	configs.HashIncrBy(joinedUserKey, "slots_on_sale", float64(data.SaleSlots))
	keyName := helper.KeyName(float64(data.SaleFee))
	key := fmt.Sprintf("%s%d:%s:%d:%s", common.TRADE_ON_SALE, data.MatchID, data.PredictionID, data.OptionID, keyName)
	pushInKeyData := fmt.Sprintf("%s-%s-%.1f", data.UserID, data.RecordID, data.SaleFee)
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
	RecordId, err := primitive.ObjectIDFromHex(data.RecordID)
	if err != nil {
		return
	}
	filter := bson.M{"_id": RecordId}
	update := bson.M{
		"$inc": bson.M{"slots_on_sale": data.SaleSlots},
		"$set": bson.M{"sale_trade": saleSlotArray},
	}
	opts := options.Update().SetUpsert(true)
	_, err = collections.TRADE_JOINED_COLLECTION.UpdateOne(context.TODO(), filter, update, opts)
	jsonString, err := json.Marshal(saleSlotArray)
	expiration := 30 * 24 * time.Hour
	configs.SetStringValue(joinSaleTradeKey, string(jsonString), expiration)
}

func CancelSale(data common.CancelSaleRequestData) (error, common.Response) {
	joinSaleTradeKey := fmt.Sprintf("%s%s", common.SLOTS_ON_SALE, data.RecordID)
	joinPredKeyResult, err := configs.GetStringValue(joinSaleTradeKey)
	if err != nil {
		response.Message = "Invalid key"
		return err, response
	}
	var userTradeSaleData []models.SaleTrade
	if err := json.Unmarshal([]byte(joinPredKeyResult), &userTradeSaleData); err != nil {
		response.Message = "Invalid Data"
		return err, response
	}
	keyName := helper.KeyName(float64(data.SaleFee))
	key := fmt.Sprintf("%s%d:%s:%d:%s", common.TRADE_ON_SALE, data.MatchID, data.PredictionID, data.OptionID, keyName)
	removeInKeyData := fmt.Sprintf("%s-%s-%.1f", data.UserID, data.RecordID, data.SaleFee)
	configs.RemoveListElement(key, int64(data.CancelSlots), removeInKeyData)
	for i, v := range userTradeSaleData {
		if v.SaleFee == data.SaleFee {
			userTradeSaleData[i].SaleSlots = 0
			break
		}
	}
	jsonString, err := json.Marshal(userTradeSaleData)
	expiration := 30 * 24 * time.Hour
	configs.SetStringValue(joinSaleTradeKey, string(jsonString), expiration)
	return nil, response
}
