package configs

import (
	"context"
	"fmt"
	"go_echo/common"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// }

var (
	ctx    = context.Background()
	client *redis.Client
)

func init() {
	redisURL := common.REDIS_URL
	if redisURL == "" {
		log.Panic("REDIS_URL environment variable is not set.")
	}

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Panic("Error parsing Redis URL:", err)
	}

	client = redis.NewClient(options)

	// Check if the Redis connection is successful
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Panic("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)
}

// hash function

func GetHashKeyValues(key string) (map[string]string, error) {
	result, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Hmset(key string, fields map[string]interface{}) error {
	err := client.HMSet(ctx, key, fields).Err()
	if err != nil {
		return err
	}
	return nil
}

func HashIncrBy(key string, field string, value float64) (float64, error) {
	incrVal, err := client.HIncrByFloat(ctx, key, field, value).Result()
	if err != nil {
		return 0, err
	}
	return incrVal, nil
}

func HashGetByKeyField(key, field string) {
	client.HGet(ctx, key, field)
}

// list functions

func Rpush(key string, value string) error {
	_, err := client.RPush(ctx, key, value).Result()
	if err != nil {
		return nil
	}
	return err
}

func RemoveListElement(key string, count int64, valueToRemove string) (int64, error) {
	removedCount, err := client.LRem(ctx, key, count, valueToRemove).Result()
	if err != nil {
		return 0, err
	}
	return removedCount, nil
}

// string functions

func IncrementBY(key string, value float64) error {
	err := client.IncrByFloat(ctx, key, value)
	if err != nil {
		return nil
	}
	return nil
}

func SetStringValue(key, value string, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

func GetStringValue(key string) (string, error) {
	data, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

// expiry function

func SetWithExpirationDays(key string, expirationDays int) error {
	expirationDuration := time.Duration(expirationDays) * 24 * time.Hour
	_, err := client.Expire(ctx, key, expirationDuration).Result()
	if err != nil {
		return err
	}

	return nil
}
