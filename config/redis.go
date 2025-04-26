package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		fmt.Println("Error: REDIS_URI environment variable is not set")
		return
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: redisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Redis Connection Failed: %v\n", err)
	} else {
		fmt.Println("Redis Connected")
	}
}

// CheckRedisConnection verifies if Redis is connected
func CheckRedisConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis connection error: %v", err)
	}
	return nil
}

// ReconnectRedis attempts to reconnect to Redis
func ReconnectRedis() {
	fmt.Println("Attempting to reconnect to Redis...")
	ConnectRedis()
}
