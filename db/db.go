package db

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func Initialize() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := client.Ping(context.Background()).Result(); err == nil {
		return nil, err
	}

	return client, nil
}
