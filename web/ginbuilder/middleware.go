package ginbuilder

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func LogAccessMiddleware() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.Set(TRACE_REQUEST_KEY, logrus.Fields{})

		startTime := time.Now()
		context.Next()

		if context.GetBool(NO_ACCESS_LOG_PRINT) {
			return
		}

		// print access log
		duration := time.Now().Sub(startTime)

		var logFields logrus.Fields
		value, exist := context.Get(TRACE_REQUEST_KEY)
		if !exist {
			logFields = logrus.Fields{}
		}

		logFields, ok := value.(logrus.Fields)
		if !ok || logFields == nil {
			logFields = logrus.Fields{}
		}

		method := context.Request.Method
		statusCode := context.Writer.Status()
		path := context.Request.URL.Path
		raw := context.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		logFields["status_code"] = statusCode
		logFields["duration"] = duration
		logFields["path"] = path
		logFields["method"] = method
		logFields["client_ip"] = context.ClientIP()

		logrus.WithFields(logFields).Infof("%s | %d | %v | %s", method, statusCode, duration, path)

	}
}
