package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CancelAllUnMatchedBody struct {
	PredictionID string `json:"prediction_id" validate:"required"`
	UserId       string `json:"user_id" validate:"required"`
	MatchID      int    `json:"match_id" validate:"required,min=1"`
	Sport        int8   `json:"sport" validate:"required,min=1,max=4"`
}

type RequestSaleBody struct {
	RecordID     string  `json:"record_id" validate:"required,hexadecimal,len=24"`
	SaleFee      float32 `json:"sale_fee" validate:"required,numeric,gte=0.5,lte=9.5"`
	PredictionID string  `json:"prediction_id" validate:"required,hexadecimal,len=24"`
	MatchID      int     `json:"match_id" validate:"required"`
	OptionID     int8    `json:"option_id" validate:"required,oneof=1 2"`
	Sport        int8    `json:"sport" validate:"required"`
	SaleSlots    int     `json:"sale_slots" validate:"required,numeric,gte=1"`
	UserID       string  `json:"user_id" validate:"required,hexadecimal,len=24"`
}

type CancelSaleRequestData struct {
	RecordID     string  `json:"record_id" validate:"required,hexadecimal,len=24"`
	SaleFee      float32 `json:"sale_fee" validate:"required,numeric,gte=0.5,lte=9.5"`
	PredictionID string  `json:"prediction_id" validate:"required,hexadecimal,len=24"`
	MatchID      int     `json:"match_id" validate:"required"`
	OptionID     int8    `json:"option_id" validate:"required,oneof=1 2"`
	Sport        int8    `json:"sport" validate:"required"`
	CancelSlots  int     `json:"cancel_slots" validate:"required,numeric,gte=1"`
	UserID       string  `json:"user_id" validate:"required,hexadecimal,len=24"`
}

type SaleOnDifferentPrice struct {
	RecordID     string  `json:"record_id" validate:"required,hexadecimal,len=24"`
	SaleFee      float32 `json:"sale_fee" validate:"required,numeric,gte=0.5,lte=9.5"`
	OldSaleFee   float32 `json:"old_sale_fee" validate:"required,numeric,gte=0.5,lte=9.5"`
	PredictionID string  `json:"prediction_id" validate:"required,hexadecimal,len=24"`
	MatchID      int     `json:"match_id" validate:"required"`
	OptionID     int8    `json:"option_id" validate:"required,oneof=1 2"`
	Sport        int8    `json:"sport" validate:"required"`
	Slots        int     `json:"slots" validate:"required,numeric,gte=1"`
	UserID       string  `json:"user_id" validate:"required,hexadecimal,len=24"`
}

type TradeTransaction struct {
	MatchId        int                `bson:"match_id" json:"match_id"`
	SportType      int                `bson:"sport_type" json:"sport_type"`
	PredictionType string             `bson:"prediction_type" json:"prediction_type"`
	CashBalance    float64            `bson:"cash_balance" json:"cash_balance"`
	WinningBalance float64            `bson:"winning_balance" json:"winning_balance"`
	DebitedAmount  float64            `bson:"debited_amount" json:"debited_amount"`
	CreditedAmount float64            `bson:"credited_amount" json:"credited_amount"`
	PredictionID   primitive.ObjectID `bson:"prediction_id" json:"prediction_id"`
	UserId         primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	// Check if it's a validation error
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			fmt.Println(">>>>>errors", e)
			// Field name is converted to lowercase and used as the key
			field := strings.ToLower(e.Field())
			// Use tag name as error message
			message := e.Tag()
			errors[field] = message
		}
	}

	return errors
}
