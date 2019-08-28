package uctx

import (
	"context"
)

func BackgroundWithCancel() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithCancel(context.Background())
}
