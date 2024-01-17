package common

import "fmt"

func GetPredictionKey(prefix string, matchId int, sport int8, predictionId string) string {
	return fmt.Sprintf("%s%d:%d:%s", prefix, matchId, sport, predictionId)
}

func GetJoinedTradeKey(prefix string, matchId int, sport int8, predictionId string, userId string, recordId string) string {
	return fmt.Sprintf("%s%d:%d:%s:%s:%s", prefix, matchId, sport, predictionId, userId, recordId)
}
