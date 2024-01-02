package common

type TradeConstant struct {
	JOINED_PREDICTION_TRADE string
	TRADE_UNMATCHED         string
	TOTAL_TRADE_JOINED      string
}

var TRADE_CONSTANT = TradeConstant{
	JOINED_PREDICTION_TRADE: "joined-prediction-trade:",
	TRADE_UNMATCHED:         "trade-umatch:",
	TOTAL_TRADE_JOINED:      "total-trade-joined:",
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// var RESPONSE = response{
// 	Status:  false,
// 	Message: "Invalid Request",
// }
