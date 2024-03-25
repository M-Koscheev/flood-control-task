package internal

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

const (
	defaultMeetings   int64 = 10
	defaultTimeToLive       = 2 * time.Minute
)

var meetingsExceeded = errors.New("maximum number of meetings was exceeded")

type Controller struct {
	client      *redis.Client
	maxMeetings int64
	timeToLive  time.Duration
}

func NewController(client *redis.Client, amount int64, timeToLive time.Duration) Controller {
	if amount < 0 {
		amount = defaultMeetings
	}
	if timeToLive < 0 {
		timeToLive = defaultTimeToLive
	}
	return Controller{client: client, maxMeetings: amount, timeToLive: timeToLive}
}

func (ctr Controller) Check(ctx context.Context, userID int64) (bool, error) {
	stringId := strconv.FormatInt(userID, 36)
	_, err := ctr.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipeline := ctr.client.Pipeline()
		amount, err := pipeline.Incr(ctx, stringId).Result()
		if err != nil {
			return err
		} else if amount > ctr.maxMeetings {
			return meetingsExceeded
		}

		curTTL, err := pipeline.TTL(ctx, stringId).Result()
		if err != nil {
			return err
		}
		if curTTL < 0 {
			pipeline.Expire(ctx, stringId, ctr.timeToLive)
		}
		return nil
	})
	if err != nil {
		if errors.Is(meetingsExceeded, err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
