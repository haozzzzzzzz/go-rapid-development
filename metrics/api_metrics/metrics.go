package api_metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/metrics"
	"github.com/haozzzzzzzz/go-rapid-development/web/wgin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/xiam/to"
	"time"
)

type ApiMetrics struct {
	Namespace  string
	InstanceId string

	CallTimesCounterVec      *prometheus.CounterVec // 接口调用次数统计
	ReturnCodeCounterVec     *prometheus.CounterVec // 业务码返回统计
	HttpStatusCodeCounterVec *prometheus.CounterVec // http状态码返回统计
	SpentTimeSummaryVec      *prometheus.SummaryVec // 耗时
}

func NewApiMetrics(namespace string, instanceId string) (apiMetrics *ApiMetrics, err error) {
	apiMetrics = &ApiMetrics{
		Namespace:  namespace,
		InstanceId: instanceId,
	}

	apiMetrics.CallTimesCounterVec, err = metrics.NewCounterVec(
		namespace,
		"api",
		"call_times",
		"调用次数",
		[]string{"instance", "uri", "type"},
	)
	if err != nil {
		logrus.Errorf("new api metrics call times counter vec failed. error: %s", err)
		return
	}

	apiMetrics.ReturnCodeCounterVec, err = metrics.NewCounterVec(
		namespace,
		"api",
		"return_code_counter",
		"返回码统计",
		[]string{"instance", "uri", "code"},
	)
	if err != nil {
		logrus.Errorf("new api metrics return code counter vec failed. error: %s", err)
		return
	}

	apiMetrics.HttpStatusCodeCounterVec, err = metrics.NewCounterVec(
		namespace,
		"api",
		"http_status_code_counter",
		"http返回状态码统计",
		[]string{"instance", "uri", "status_code"},
	)
	if err != nil {
		logrus.Errorf("new api metrics http status code vec failed. error: %s", err)
		return
	}

	apiMetrics.SpentTimeSummaryVec, err = metrics.NewSummaryVec(
		namespace,
		"api",
		"spent_time",
		"耗时",
		[]string{"instance", "uri"},
	)
	if err != nil {
		logrus.Errorf("new api metrics api spent time vec failed. error: %s", err)
		return
	}

	return
}

// NewGinMiddleware gin中间件
func NewGinMiddleware(
	apiMetrics *ApiMetrics,
	successRetCode int, // 成功的业务码
	skipPaths []string,
) gin.HandlerFunc {
	skipPathMap := make(map[string]bool)
	for _, skipPath := range skipPaths {
		skipPathMap[skipPath] = true
	}

	return func(context *gin.Context) {
		urlPath := context.Request.URL.Path
		if skipPathMap[urlPath] {
			context.Next()
			return
		}

		instanceId := apiMetrics.InstanceId

		apiMetrics.CallTimesCounterVec.WithLabelValues(instanceId, urlPath, "total").Inc()
		startTime := time.Now()

		// before
		context.Next()
		// after

		isHandleFailed := false

		endTime := time.Now()
		statusCode := context.Writer.Status()

		apiMetrics.SpentTimeSummaryVec.WithLabelValues(instanceId, urlPath).Observe(float64(endTime.Sub(startTime).Milliseconds()))
		apiMetrics.HttpStatusCodeCounterVec.WithLabelValues(instanceId, urlPath, to.String(statusCode)).Inc()

		if statusCode != 200 {
			isHandleFailed = true
		}

		retCode, ok := wgin.GetResponseRetCode(context)
		if ok {
			apiMetrics.ReturnCodeCounterVec.WithLabelValues(instanceId, urlPath, to.String(retCode)).Inc()
			if retCode != successRetCode {
				isHandleFailed = true
			}
		}

		if isHandleFailed {
			apiMetrics.CallTimesCounterVec.WithLabelValues(instanceId, urlPath, "error").Inc()
		}
	}
}
