package task

import (
	"context"
	"github.com/haozzzzzzzz/go-rapid-development/v2/utils/uctx"
	"time"
)

func NewBackgroundContext(name string) (ctx context.Context, cancel func(err error)) {
	rawCtx, rawCancel := uctx.BackgroundWithCancel()
	ctx = rawCtx
	cancel = func(err error) {
		rawCancel()
	}
	return
}

func NewBackgroundContextWithTimeout(name string, timeout time.Duration) (ctx context.Context, cancel func(err error)) {
	rawCtx, rawCancel := uctx.BackgroundWithTimeout(timeout)
	ctx = rawCtx
	cancel = func(err error) {
		rawCancel()
	}
	return
}
