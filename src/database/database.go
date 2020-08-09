package database

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Check if connection has been established.
	_, err := client.Ping(context.Background()).Result()
	return client, err
}
