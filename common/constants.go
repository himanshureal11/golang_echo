package common

type TradeConstant struct {
	JOINED_PREDICTION_TRADE string
	TRADE_UNMATCHED         string
	TOTAL_TRADE_JOINED      string
	MATCH_TRADE_PREDICTION  string
	TRADE_ON_SALE           string
	TRADE_SALE_USER_META    string
}

var TRADE_CONSTANT = TradeConstant{
	JOINED_PREDICTION_TRADE: "joined-prediction-trade:",
	TRADE_UNMATCHED:         "trade-umatch:",
	TOTAL_TRADE_JOINED:      "total-trade-joined:",
	MATCH_TRADE_PREDICTION:  "match-trade-prediction:",
	TRADE_ON_SALE:           "trade-on-sale:",
	TRADE_SALE_USER_META:    "trade-on-sale-user-meta:",
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
