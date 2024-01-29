package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var PORT string
var MONGO_URL string
var REDIS_URL string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	PORT = os.Getenv("PORT")
	MONGO_URL = os.Getenv("MONGODB_URL")
	REDIS_URL = os.Getenv("REDIS_URL")
}

// var
