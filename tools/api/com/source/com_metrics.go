package source

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateComMetrics(comDir string) (err error) {
	metricsDir := fmt.Sprintf("%s/metrics", comDir)
	err = os.MkdirAll(metricsDir, proj.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make metrics dir %q failed. %s.", metricsDir, err)
		return
	}

	metricsFilePath := fmt.Sprintf("%s/metrics.go", metricsDir)
	err = ioutil.WriteFile(metricsFilePath, []byte(metricsFileText), proj.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write metrics file %q failed. %s.", metricsFilePath, err)
		return
	}

	return
}

var metricsFileText = `package metrics

import (
	"github.com/haozzzzzzzz/go-rapid-development/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	// 服务次数统计
	SERVICE_TIMES_COUNTER_VEC   *prometheus.CounterVec
	SERVICE_TIMES_COUNTER_PANIC prometheus.Counter

	// API次数统计
	API_EXEC_TIMES_COUNTER_VEC      *prometheus.CounterVec
	API_EXEC_TIMES_COUNTER_TOTAL    prometheus.Counter // api总执行次数
	API_EXEC_TIMES_COUNTER_ABNORMAL prometheus.Counter // api非正常返回结果次数
	API_EXEC_TIMES_COUNTER_PANIC    prometheus.Counter // api panic

	// API返回码分布
	API_RETURN_CODE_COUNTER_VEC *prometheus.CounterVec

	// 接口调用次数
	API_URI_CALL_TIMES_COUNTER_VEC *prometheus.CounterVec

	// API响应时间分布
	API_SPENT_TIME_SUMMARY prometheus.Summary
)

func init() {
	var err error

	// system次数统计
	SERVICE_TIMES_COUNTER_VEC, err = metrics.NewCounterVec(
		constant.ServiceName,
		"system",
		"service_times",
		"服务次数统计",
		[]string{"type"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	SERVICE_TIMES_COUNTER_PANIC = SERVICE_TIMES_COUNTER_VEC.WithLabelValues("panic")

	// api
	API_EXEC_TIMES_COUNTER_VEC, err = metrics.NewCounterVec(
		constant.ServiceName,
		"api",
		"exec_times",
		"api发生次数",
		[]string{"type"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	API_EXEC_TIMES_COUNTER_TOTAL = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues("total")
	API_EXEC_TIMES_COUNTER_ABNORMAL = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues("abnormal")
	API_EXEC_TIMES_COUNTER_PANIC = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues("panic")

	// api返回码
	API_RETURN_CODE_COUNTER_VEC, err = metrics.NewCounterVec(
		constant.ServiceName,
		"api",
		"return_code",
		"api返回码",
		[]string{"code"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	// api调用次数
	API_URI_CALL_TIMES_COUNTER_VEC, err = metrics.NewCounterVec(
		constant.ServiceName,
		"api",
		"uri_call_times",
		"uri调用次数",
		[]string{"uri"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	// api响应时间
	API_SPENT_TIME_SUMMARY, err = metrics.NewSummary(
		constant.ServiceName,
		"api",
		"spent_time",
		"api消耗时间分布",
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

}
`
