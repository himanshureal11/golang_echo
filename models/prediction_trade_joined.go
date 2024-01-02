package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PredictionTradeJoined struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PredictionID            primitive.ObjectID `bson:"prediction_id" json:"prediction_id"`
	QuestionEnglish         string             `bson:"question_english" json:"question_english"`
	QuestionHindi           string             `bson:"question_hindi" json:"question_hindi"`
	UserID                  primitive.ObjectID `bson:"user_id" json:"user_id"`
	MatchID                 int                `bson:"match_id" json:"match_id"`
	SeriesID                int                `bson:"series_id" json:"series_id"`
	Sport                   int                `bson:"sport" json:"sport"`
	OptionID                int                `bson:"option_id" json:"option_id"`
	SlotFee                 float64            `bson:"slot_fee" json:"slot_fee"`
	TotalMatched            int                `bson:"total_matched" json:"total_matched"`
	TotalSlot               int                `bson:"total_slot" json:"total_slot"`
	TotalCash               float64            `bson:"total_cash" json:"total_cash"`
	TotalWining             float64            `bson:"total_wining" json:"total_wining"`
	WinReward               float64            `bson:"win_reward" json:"win_reward"`
	WinDistribute           bool               `bson:"win_distribute" json:"win_distribute"`
	Status                  int                `bson:"status" json:"status"`
	IsPredCancel            int                `bson:"is_pred_cancel" json:"is_pred_cancel"`
	IsPredWinset            int                `bson:"is_pred_winset" json:"is_pred_winset"`
	IsSlotCancel            int                `bson:"is_slot_cancel" json:"is_slot_cancel"`
	RefundCashAmount        float64            `bson:"refund_cash_amount" json:"refund_cash_amount"`
	RefundWinAmount         float64            `bson:"refund_win_amount" json:"refund_win_amount"`
	CancelledSlotNumber     int                `bson:"cancelled_slot_number" json:"cancelled_slot_number"`
	SaleFee                 float64            `bson:"sale_fee" json:"sale_fee"`
	SlotsOnSale             int                `bson:"slots_on_sale" json:"slots_on_sale"`
	SaleAdminComssionAmt    float64            `bson:"sale_admin_comssion_amount" json:"sale_admin_comssion_amount"`
	SoldSlots               int                `bson:"sold_slots" json:"sold_slots"`
	SoldSlotsRewards        float64            `bson:"sold_slots_rewards" json:"sold_slots_rewards"`
	MatchedTime             int                `bson:"matched_time" json:"matched_time"`
	InWinning               float64            `bson:"in_winning" json:"in_winning"`
	InCash                  float64            `bson:"in_cash" json:"in_cash"`
	AdminCommissionDeducted float64            `bson:"admin_commission_deducted" json:"admin_commission_deducted"`
	PendingAmount           float64            `bson:"pending_amount" json:"pending_amount"`
}

// Add any additional methods or validations as needed
