package services

import (
	"go_echo/common"
	"go_echo/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var wg sync.WaitGroup

func JoinTrade(data []models.PredictionTradeJoined) (error, common.Response) {
	var joinPredictionData []interface{}
	for i := 0; i < 10000; i++ {
		data[0].ID = primitive.NewObjectID()
		data[0].SetTimestamps()
		joinPredictionData = append(joinPredictionData, data[0])
		// wg.Add(1)
	}
	models.InsertMany(joinPredictionData)
	// wg.Wait()
	return nil, common.Response{}
}
