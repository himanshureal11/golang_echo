package controllers

import (
	"go_echo/models"
	"go_echo/services"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func JoinTrade(c echo.Context) error {
	// var validate = validator.New()
	var requestBody []models.PredictionTradeJoined
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	// for _, v := range requestBody {
	// 	if err := validate.Struct(v); err != nil {
	// 		validationErrors := common.GetValidationErrors(err)
	// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validation Error", "details": validationErrors})
	// 	}
	// }
	err, response := services.JoinTrade(requestBody)
	if err != nil {
		log.Println("<<<<<err>>>>>", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	return c.JSON(http.StatusOK, response)
}
