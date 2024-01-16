package controllers

import (
	"go_echo/common"
	"go_echo/models"
	"go_echo/services"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func JoinTrade(c echo.Context) error {
	var validate = validator.New()
	var requestBody []models.PredictionTradeJoined
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	for _, v := range requestBody {
		if err := validate.Struct(v); err != nil {
			validationErrors := common.GetValidationErrors(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validation Error", "details": validationErrors})
		}
	}
	err, response := services.JoinTrade(requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	return c.JSON(http.StatusOK, response)
}
