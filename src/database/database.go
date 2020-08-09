package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"storage-api/src/config"
)

func NewRedisClient(c config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: c.Password,
		DB:       c.DB,
	})

	// Check if connection has been established.
	_, err := client.Ping(context.Background()).Result()
	return client, err
}

