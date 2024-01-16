package services

import (
	"go_echo/common"
	"go_echo/models"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var wg sync.WaitGroup

func JoinTrade(data []models.PredictionTradeJoined) (error, common.Response) {
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		data[0].ID = primitive.NewObjectID()
		go models.InsertData(data[0], &wg)
	}
	wg.Wait()
	return nil, common.Response{}
}
