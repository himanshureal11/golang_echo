package common

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CancelAllUnMatchedBody struct {
	PredictionID string `json:"prediction_id" validate:"required"`
	UserId       string `json:"user_id" validate:"required"`
	MatchID      int    `json:"match_id" validate:"required,min=1"`
	Sport        int8   `json:"sport" validate:"required,min=1,max=4"`
}

type RequestSaleBody struct {
	RecordID     string `json:"record_id" validate:"required,hex,len=24"`
	SaleFee      string `json:"sale_fee" validate:"required,numeric,gte=0.5,lte=9.5"`
	PredictionID string `json:"prediction_id" validate:"required,hex,len=24"`
	MatchID      string `json:"match_id" validate:"required"`
	OptionID     string `json:"option_id" validate:"required,oneof=1 2"`
	Sport        string `json:"sport" validate:"required"`
	SaleSlots    string `json:"sale_slots" validate:"required,numeric,gte=1"`
	UserID       string `json:"user_id" validate:"required,hex,len=24"`
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
