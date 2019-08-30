package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// 只递增的单个数目
func NewCounter(
	namespace string,
	subsystem string,
	name string,
	help string,
) (counter prometheus.Counter, err error) {

	counter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})

	err = prometheus.Register(counter)
	if nil != err {
		logrus.Errorf("register prometheus counter failed. %s.", err)
		return
	}

	return
}

// 可增减的单个数目
func NewGuage(
	namespace string,
	subsystem string,
	name string,
	help string,
) (gauge prometheus.Gauge, err error) {
	gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})

	err = prometheus.Register(gauge)
	if nil != err {
		logrus.Errorf("register prometheus gauge failed. %s.", err)
		return
	}

	return
}

var SummaryObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

// 单个统计分布
func NewSummary(
	namespace string,
	subsystem string,
	name string,
	help string,
) (summary prometheus.Summary, err error) {
	summary = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       name,
		Help:       help,
		Objectives: SummaryObjectives,
	})

	err = prometheus.Register(summary)
	if nil != err {
		logrus.Errorf("register prometheus summary failed. %s.", err)
		return
	}

	return
}

// 只递增的多个数目
func NewCounterVec(
	namespace string,
	subsystem string,
	name string,
	help string,
	labels []string,
) (counterVec *prometheus.CounterVec, err error) {
	counterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	}, labels)

	err = prometheus.Register(counterVec)
	if nil != err {
		logrus.Errorf("register prometheus counter vec failed. %s.", err)
		return
	}

	return
}

// 可增减的多个数目
func NewGaugeVec(
	namespace string,
	subsystem string,
	name string,
	help string,
	labels []string,
) (gaugeVec *prometheus.GaugeVec, err error) {
	gaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	}, labels)

	err = prometheus.Register(gaugeVec)
	if nil != err {
		logrus.Errorf("register prometheus gauge vec failed. %s.", err)
		return
	}

	return
}

// 多个统计
func NewSummaryVec(
	namespace string,
	subsystem string,
	name string,
	help string,
	labels []string,
) (summaryVec *prometheus.SummaryVec, err error) {
	summaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       name,
		Help:       help,
		Objectives: SummaryObjectives,
	}, labels)

	err = prometheus.Register(summaryVec)
	if nil != err {
		logrus.Errorf("register prometheus summary vec failed. %s.", err)
		return
	}

	return
}

func PrometheusGinMetrics(routes gin.IRoutes, metricsPath string) {
	routes.GET(metricsPath, func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Println(err)
			}
		}()

		promhttp.Handler().ServeHTTP(context.Writer, context.Request)

		// don't print access log
		context.Set(ginbuilder.NO_ACCESS_LOG_PRINT, true)

	})
}
