package xray

import (
	"net/http"

	"context"

	"fmt"

	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gin-gonic/gin"
)

func XRayGinMiddleware(strSegmentNamer string) func(*gin.Context) {
	return func(context *gin.Context) {
		xray.Handler(xray.NewFixedSegmentNamer(strSegmentNamer), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.Request = r
			context.Next()
		})).ServeHTTP(context.Writer, context.Request)
	}
}

// 需要在调用后调用cancel
func NewBackgroundContext(name string) (
	ctx context.Context,
	seg *xray.Segment,
	cancel func(err error),
) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx, seg = xray.BeginSegment(ctx, fmt.Sprintf("background_%s", name))
	cancel = func(err error) {
		seg.Close(err)
		cancelFunc()
	}

	seg.Sampled = false
	return
}

func NewBackgroundContextWithTimeout(name string, timeout time.Duration) (ctx context.Context, seg *xray.Segment, cancel func(err error)) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	ctx, seg = xray.BeginSegment(ctx, fmt.Sprintf("background_%s", name))
	cancel = func(err error) {
		seg.Close(err)
		cancelFunc()
	}
	seg.Sampled = false
	return
}
