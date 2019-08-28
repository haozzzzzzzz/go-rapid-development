package uctx

import (
	"context"
	"time"
)

func BackgroundWithCancel() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func BackgroundWithTimeout(timeout time.Duration) (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
