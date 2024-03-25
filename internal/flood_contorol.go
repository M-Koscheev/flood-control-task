package internal

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"task/db"
	"task/lib"
)

type Controller struct{}

var client *redis.Client = nil

func (ctr Controller) Check(ctx context.Context, userID int64) (bool, error) {
	if client == nil {
		newClient, err := db.Initialize()
		if err != nil {
			return false, err
		}
		client = newClient
	}

	var person lib.UserInfo
	stringId := strconv.FormatInt(userID, 36)
	err := client.Get(ctx, stringId).Scan(&person)
	if err != nil {
		return false, err
	}

	person.MeetsAmount += 1
	if err := client.Set(ctx, stringId, person, 0).Err(); err != nil {
		return false, err
	}

	return true, nil
}
