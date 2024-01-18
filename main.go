package main

import (
	"go_echo/common"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	// Create a new Echo instance
	e := echo.New()
	e.Use(requestLogger)
	InitializeRoutes(e)

	// Start the server on port 8080
	PORT := common.PORT
	e.Logger.Fatal(e.Start(PORT))
}

var response common.Response

func requestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Record the start time of the request
		start := time.Now()
		log.Printf("Option[%s] %s - %d", c.Request().Method, c.Path(), c.Response().Status)
		// Call the next middleware or handler
		err := next(c)
		if err != nil {
			response.Status = false
			response.Message = strings.Split(err.Error(), "Error:")[1]
			c.JSON(http.StatusBadRequest, response)
		}
		// Record the end time of the request
		end := time.Now()
		// Calculate the response time
		responseTime := end.Sub(start)
		// Log the request information
		log.Printf("[%s] %s - %d - %v", c.Request().Method, c.Path(), c.Response().Status, responseTime)
		return err
	}
}
