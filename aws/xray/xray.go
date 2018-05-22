package xray

import (
	"net/http"

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
