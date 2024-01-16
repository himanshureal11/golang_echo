package services

import (
	"go_echo/common"
	"go_echo/models"
)

// var wg sync.WaitGroup

func JoinTrade(data []models.PredictionTradeJoined) (error, common.Response) {
	var joinPredictionData []interface{}
	for _, v := range data {
		v.SetTimestamps()
		joinPredictionData = append(joinPredictionData, v)
	}
	err := models.InsertMany(joinPredictionData)
	if err != nil {
		return err, common.Response{Data: []string{}, Message: "Not Able to join", Status: false}
	}
	return nil, common.Response{Data: []string{}, Message: "Joined Successfully", Status: true}
}
