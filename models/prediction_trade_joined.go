package models

import (
	"context"
	"fmt"
	"go_echo/collections"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SaleTrade struct {
	SaleSlots int     `bson:"sale_slot" json:"sale_slot"`
	SoldSlots int     `bson:"sold_slot" json:"sold_slot"`
	SaleFee   float32 `bson:"sale_fee" json:"sale_fee"`
}

type PredictionTradeJoined struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PredictionID            primitive.ObjectID `bson:"prediction_id,omitempty" json:"prediction_id,omitempty"`
	QuestionEnglish         string             `bson:"question_english" json:"question_english" validate:"omitempty"`
	QuestionHindi           string             `bson:"question_hindi" json:"question_hindi"`
	UserID                  primitive.ObjectID `bson:"user_id" json:"user_id,omitempty"`
	MatchID                 int                `bson:"match_id" json:"match_id" validate:"omitempty"`
	SeriesID                int                `bson:"series_id" json:"series_id" validate:"omitempty"`
	Sport                   int                `bson:"sport" json:"sport" validate:"omitempty"`
	OptionID                int                `bson:"option_id" json:"option_id" validate:"omitempty"`
	SlotFee                 float64            `bson:"slot_fee" json:"slot_fee"`
	TotalMatched            int                `bson:"total_matched" json:"total_matched" default:"0"`
	TotalSlot               int                `bson:"total_slot" json:"total_slot" validate:"omitempty"`
	TotalCash               float64            `bson:"total_cash" json:"total_cash"`
	TotalWining             float64            `bson:"total_wining" json:"total_wining"`
	WinReward               float64            `bson:"win_reward" json:"win_reward" default:"0"`
	RefundCashAmount        float64            `bson:"refund_cash_amount" json:"refund_cash_amount"`
	RefundWinAmount         float64            `bson:"refund_win_amount" json:"refund_win_amount"`
	CancelledSlotNumber     int                `bson:"cancelled_slot_number" json:"cancelled_slot_number"`
	IsSlotCancel            int                `bson:"is_slot_cancel" json:"is_slot_cancel" default:"0"`
	WinDistribute           bool               `bson:"win_distribute" json:"win_distribute" default:"false"`
	Status                  int                `bson:"status" json:"status" enum:"1,0" default:"1"`
	IsPredCancel            int                `bson:"is_pred_cancel" json:"is_pred_cancel" default:"0"`
	IsPredWinset            int                `bson:"is_pred_winset" json:"is_pred_winset" default:"0"`
	SlotsOnSale             int                `bson:"slots_on_sale" json:"slots_on_sale"`
	SaleFee                 float64            `bson:"sale_fee" json:"sale_fee"`
	NewSaleFee              float64            `bson:"new_sale_fee" json:"new_sale_fee"`
	SoldSlots               int                `bson:"sold_slots" json:"sold_slots"`
	SoldSlotsRewards        float64            `bson:"sold_slots_rewards" json:"sold_slots_rewards"`
	SaleAdminComssionAmount int                `bson:"sale_admin_comssion_amount" json:"sale_admin_comssion_amount"`
	BuyFromSale             bool               `bson:"buy_from_sale" json:"buy_from_sale" default:"false"`
	MatchedTime             int                `bson:"matched_time" json:"matched_time"`
	InWinning               float64            `bson:"in_winning" json:"in_winning" default:"0"`
	InCash                  float64            `bson:"in_cash" json:"in_cash" default:"0"`
	PendingAmount           float64            `bson:"pending_amount" json:"pending_amount" validate:"omitempty"`
	ExtraRefundCash         float64            `bson:"extra_refund_Cash" json:"extra_refund_Cash" validate:"omitempty"`
	ExtraSlotMatchedAmount  float64            `bson:"extra_slot_matched_amount" json:"extra_slot_matched_amount" validate:"omitempty"`
	SaleTrade               []SaleTrade        `bson:"sale_trade" json:"sale_trade"`
	CreatedAt               time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt               time.Time          `bson:"updated_at" json:"updated_at"`
}

// Add any additional methods or validations as needed

func (p *PredictionTradeJoined) SetTimestamps() {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
}

func InsertData(data PredictionTradeJoined, wg *sync.WaitGroup) error {
	defer wg.Done()
	data.SetTimestamps()
	_, err := collections.TRADE_JOINED_COLLECTION.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Println(">>>>error", err)
		return err
	}
	return nil
}

func InsertMany(data []interface{}) error {
	_, err := collections.TRADE_JOINED_COLLECTION.InsertMany(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}
