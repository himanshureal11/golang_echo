package common

import "fmt"

func GetPredictionKey(prefix string, matchId string, sport string, predictionId string) string {
	return fmt.Sprintf("%d%d:%d:%d", prefix, matchId, sport, predictionId)
}

func GetJoinedTradeKey(prefix string, matchId string, sport string, predictionId string, userId string, recordId string) string {
	return fmt.Sprintf("%d%d:%d:%d:%d:%d", prefix, matchId, sport, predictionId, userId, recordId)
}
