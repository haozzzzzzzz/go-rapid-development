package task

import (
	"context"
	"github.com/haozzzzzzzz/go-rapid-development/aws/xray"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uctx"
	"time"
)

func NewBackgroundContext(name string) (ctx context.Context, cancel func(err error)) {
	if UseXray {
		ctx, _, cancel = xray.NewBackgroundContext(name)
	} else {
		rawCtx, rawCancel := uctx.BackgroundWithCancel()
		ctx = rawCtx
		cancel = func(err error) {
			rawCancel()
		}
	}
	return
}

func NewBackgroundContextWithTimeout(name string, timeout time.Duration) (ctx context.Context, cancel func(err error)) {
	if UseXray {
		ctx, _, cancel = xray.NewBackgroundContextWithTimeout(name, timeout)

	} else {
		rawCtx, rawCancel := uctx.BackgroundWithTimeout(timeout)
		ctx = rawCtx
		cancel = func(err error) {
			rawCancel()
		}
	}
	return
}
