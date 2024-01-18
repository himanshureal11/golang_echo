package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var PORT string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	PORT = os.Getenv("PORT")
}

// var
