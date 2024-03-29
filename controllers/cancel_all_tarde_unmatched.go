package controllers

import (
	"go_echo/common"
	"go_echo/services"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func CancelAllUnMatchedTrade(c echo.Context) error {
	var requestBody common.CancelAllUnMatchedBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	if err := validate.Struct(requestBody); err != nil {
		validationErrors := common.GetValidationErrors(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Validation Error", "details": validationErrors})
	}
	err, response := services.CancelAllUnMatchedTrade(requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	return c.JSON(http.StatusOK, response)
}
