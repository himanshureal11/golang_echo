package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var (
	ctx    = context.Background()
	client *redis.Client
)

func init() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is not set.")
	}

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Error parsing Redis URL:", err)
	}

	client = redis.NewClient(options)

	// Check if the Redis connection is successful
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)
}

func SetStringValue(key, value string, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

func GetHashKeyValues(key string) (map[string]string, error) {
	result, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func RemoveListElement(key string, count int64, valueToRemove string) (int64, error) {
	removedCount, err := client.LRem(ctx, key, count, valueToRemove).Result()
	if err != nil {
		return 0, err
	}
	return removedCount, nil
}

func Hmset(key string, fields map[string]interface{}) error {
	err := client.HMSet(ctx, key, fields).Err()
	if err != nil {
		return err
	}
	return nil
}

func IncrementBY(key string, value float64) error {
	err := client.IncrByFloat(ctx, key, value)
	if err != nil {
		return nil
	}
	return nil
}
