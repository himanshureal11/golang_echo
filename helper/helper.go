package helper

import (
	"context"
	"go_echo/collections"
	"go_echo/common"
	"log"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func KeyName(key float64) string {
	var keyName string

	switch key {
	case 0.5:
		keyName = "pointfive"
	case 1:
		keyName = "one"
	case 1.5:
		keyName = "onefive"
	case 2:
		keyName = "two"
	case 2.5:
		keyName = "twofive"
	case 3:
		keyName = "three"
	case 3.5:
		keyName = "threefive"
	case 4:
		keyName = "four"
	case 4.5:
		keyName = "fourfive"
	case 5:
		keyName = "five"
	case 5.5:
		keyName = "fivefive"
	case 6:
		keyName = "six"
	case 6.5:
		keyName = "sixfive"
	case 7:
		keyName = "seven"
	case 7.5:
		keyName = "sevenfive"
	case 8:
		keyName = "eight"
	case 8.5:
		keyName = "eightfive"
	case 9:
		keyName = "nine"
	case 9.5:
		keyName = "ninefive"
	default:
		keyName = ""
	}

	return keyName
}

func LogErrorWithStackTrace(err error) error {
	log.Println("Error:", err)
	// Capture the stack trace
	stack := make([]byte, 1024)
	runtime.Stack(stack, false)
	log.Println(string(stack))
	return err
}

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	CashBalance    float64            `bson:"cash_balance"`
	WinningBalance float64            `bson:"winning_balance"`
}

// CreateTradeTransaction creates a trade transaction
func CreateTradeTransaction(id primitive.ObjectID, data common.TradeTransaction, when string) {
	var user User

	// Fetch user data
	filter := bson.D{{Key: "_id", Value: id}}
	projection := bson.D{
		{Key: "_id", Value: 1},
		{Key: "cash_balance", Value: 1},
		{Key: "winning_balance", Value: 1},
	}

	err := collections.USERS.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&user)
	if err != nil {
		log.Panic(err)
	}

	// Adjust data based on "when"
	if when == "after" {
		data.CreditedAmount = 0
	}

	// Set data values
	data.CashBalance = user.CashBalance
	data.WinningBalance = user.WinningBalance
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	// Insert trade transaction
	_, err = collections.PREDICTION_TRANSACTION.InsertOne(context.TODO(), data)
	if err != nil {
		log.Panic(err)
	}
}
