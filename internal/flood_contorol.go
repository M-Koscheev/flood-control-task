package internal

import (
	"context"
	"errors"
	"time"
)

type Controller struct {
}

const (
	sleepDuration = 4
)

func (ctr Controller) Check(ctx context.Context, userID int64) (bool, error) {
	time.Sleep(sleepDuration * time.Second)

	if errors.Is(context.DeadlineExceeded, ctx.Err()) {
		return false, errors.New("context timeout exceeded")
	}

	if errors.Is(context.Canceled, ctx.Err()) {
		return false, errors.New("context canceled")
	}

	return true, nil
}
