package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Create a new Echo instance
	e := echo.New()
	e.Use(requestLogger)
	InitializeRoutes(e)
	// Define routes
	// e.GET("/", helloHandler)

	// Start the server on port 8080
	PORT := os.Getenv("PORT")
	log.Fatal(e.Start(PORT))
}

// Handler for the "/" route
// func helloHandler(c echo.Context) error {
// 	user := User{
// 		Name:  "John Doe",
// 		Email: "john@example.com",
// 	}
// 	return c.JSON(http.StatusOK, user)
// }

func requestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Record the start time of the request
		start := time.Now()

		// Call the next middleware or handler
		err := next(c)

		// Record the end time of the request
		end := time.Now()

		// Calculate the response time
		responseTime := end.Sub(start)

		// Log the request information
		log.Printf("[%s] %s - %d - %v", c.Request().Method, c.Path(), c.Response().Status, responseTime)

		return err
	}
}
