package services

import (
	"context"
	"go_echo/collections"
	"go_echo/common"
	"go_echo/configs"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			joinPredKeyResult, err := configs.GetHashKeyValues(joinPredKey)
			if err != nil {
				response.Message = "You don't have slots for sale"
				return nil, response
			}
			if len(joinPredKeyResult) > 1 {
				if _, ok := joinPredKeyResult["_id"]; ok {
					totalMatchedSlot, _ := strconv.Atoi(joinPredKeyResult["total_matched"])
					slotsOnSale, _ := strconv.Atoi(joinPredKeyResult["slots_on_sale"])
					saleSlots, _ := strconv.Atoi(data.SaleSlots)
					// question := joinPredKeyResult["question_english"]
					// adminCommission := joinPredKeyResult["admin_commission"]
					// slotFee := joinPredKeyResult["slot_fee"]
					// optionId := joinPredKeyResult["option_id"]
					if (slotsOnSale + saleSlots) > totalMatchedSlot {
						response.Message = "You Slot's are already on sale for this trade"
						return nil, response
					}
					if totalMatchedSlot >= (slotsOnSale + saleSlots) {
						incrementedValue, err := configs.HashIncrBy(joinPredKey, "slots_on_sale", float64(saleSlots))
						if err != nil {
							response.Message = "Invalid Sale Slot Number"
							return nil, response
						}
						if incrementedValue <= float64(totalMatchedSlot) {
							fields := map[string]any{
								"sale_fee": data.SaleFee,
							}
							update := bson.M{
								"$inc": fields,
							}
							recordId, err := primitive.ObjectIDFromHex(data.RecordID)
							if err != nil {
								response.Message = "invalid request"
								return nil, response
							}
							_, err = collections.TRADE_JOINED_COLLECTION.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: recordId}}, update)

						}
					}

				}
			}
		} else {
			return nil, response
		}
	} else {
		return nil, response
	}
	return nil, response
}

// {
// 	question:"",
// 	pred_id: "",
// 	user_id: "",
// 	record_id: "",
// 	slot_fee: "",
// 	sale_fee: "",
// 	slot_on_sale: "",
// 	pending_amount: "",
// 	sold_slots: "",
// 	sold_slot_reward: "",
// 	sale_admin_commission: "",
// 	admin_commission: "",
// 	option_id: "",
// 	option_name: ""
// 	_id: ""
// }
