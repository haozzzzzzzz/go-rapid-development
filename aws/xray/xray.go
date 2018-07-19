package xray

import (
	"net/http"

	"context"

	"fmt"

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

// seg 需要在调用后close
func NewBackgroundContext(name string) (ctx context.Context, seg *xray.Segment) {
	ctx, seg = xray.BeginSegment(context.Background(), fmt.Sprintf("background_%s", name))
	return
}
