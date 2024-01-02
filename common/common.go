package common

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type CancelAllUnMatchedBody struct {
	PredictionID string `json:"prediction_id" validate:"required"`
	UserId       string `json:"user_id" validate:"required"`
	MatchID      int    `json:"match_id" validate:"required,min=1"`
	Sport        int8   `json:"sport" validate:"required,min=1,max=4"`
}

func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	// Check if it's a validation error
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			// Field name is converted to lowercase and used as the key
			field := strings.ToLower(e.Field())
			// Use tag name as error message
			message := e.Tag()
			errors[field] = message
		}
	}

	return errors
}
