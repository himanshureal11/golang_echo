package main

import (
	"go_echo/controllers"

	"github.com/labstack/echo/v4"
)

// InitializeRoutes sets up all the routes for the application
func InitializeRoutes(e *echo.Echo) {
	e.POST("/cancel_all_unmatched_tarde", controllers.CancelAllUnMatchedTrade)
	e.POST("/sale_trade", controllers.SaleTrade)
	e.POST("/cancel_trade", controllers.CancelTrade)
	e.POST("/join_trade", controllers.JoinTrade)
}

// Handler for the "/" route
// func helloHandler(c echo.Context) error {
// 	return c.String(200, "Hello, Echo!")
// }

// // Handler for the "/user" route
// func getUserHandler(c echo.Context) error {
// 	return c.JSON(200, map[string]string{"message": "User data"})
// }
