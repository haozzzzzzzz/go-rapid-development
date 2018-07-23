package project

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
)

func (m *ProjectSource) generateComMetrics(comDir string) (err error) {
	metricsDir := fmt.Sprintf("%s/metrics", comDir)
	err = os.MkdirAll(metricsDir, project.ProjectDirMode)
	if nil != err {
		logrus.Errorf("make metrics dir %q failed. %s.", metricsDir, err)
		return
	}

	metricsFilePath := fmt.Sprintf("%s/metrics.go", metricsDir)
	err = ioutil.WriteFile(metricsFilePath, []byte(metricsFileText), project.ProjectFileMode)
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
	// ec2 instance id
	Ec2InstanceId string

	// 服务次数统计
	SERVICE_TIMES_COUNTER_VEC   *prometheus.CounterVec
	SERVICE_TIMES_COUNTER_APP_PANIC    prometheus.Counter // app service panic
	SERVICE_TIMES_COUNTER_MANAGE_PANIC prometheus.Counter // manage service panic

	// API次数统计
	API_EXEC_TIMES_COUNTER_VEC             *prometheus.CounterVec
	API_EXEC_TIMES_COUNTER_APP_TOTAL       prometheus.Counter // app api总执行数
	API_EXEC_TIMES_COUNTER_APP_ABNORMAL    prometheus.Counter // app api非正常返回结果次数
	API_EXEC_TIMES_COUNTER_APP_PANIC       prometheus.Counter // app api出错
	API_EXEC_TIMES_COUNTER_MANAGE_TOTAL    prometheus.Counter // manage api总执行数
	API_EXEC_TIMES_COUNTER_MANAGE_ABNORMAL prometheus.Counter // manage api非正常返回结果次数
	API_EXEC_TIMES_COUNTER_MANAGE_PANIC    prometheus.Counter // manage api panic

	// API返回码分布
	API_RETURN_CODE_COUNTER_VEC *prometheus.CounterVec

	// 接口调用次数
	API_URI_CALL_TIMES_COUNTER_VEC *prometheus.CounterVec

	// API响应时间分布
	API_SPENT_TIME_SUMMARY prometheus.Summary
	API_SPENT_TIME_SUMMARY_APP    prometheus.Observer // app api响应
	API_SPENT_TIME_SUMMARY_MANAGE prometheus.Observer // manage api响应

)

func init() {
	var err error

	Ec2InstanceId = config.AWSEc2InstanceIdentifyDocument.InstanceID

	// system次数统计
	SERVICE_TIMES_COUNTER_VEC, err = metrics.NewCounterVec(
		MetricsNameSpace,
		"system",
		"service_times",
		"服务次数统计",
		[]string{"instance", "origin", "type"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	// app service panic
	SERVICE_TIMES_COUNTER_APP_PANIC = SERVICE_TIMES_COUNTER_VEC.WithLabelValues(Ec2InstanceId, "app", "panic")
	// manage service panic
	SERVICE_TIMES_COUNTER_MANAGE_PANIC = SERVICE_TIMES_COUNTER_VEC.WithLabelValues(Ec2InstanceId, "manage", "panic")

    // api次数
	API_EXEC_TIMES_COUNTER_VEC, err = metrics.NewCounterVec(
		MetricsNameSpace,
		"api",
		"exec_times",
		"api发生次数",
		[]string{"instance", "origin", "type"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	API_EXEC_TIMES_COUNTER_APP_TOTAL = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues(Ec2InstanceId, "app", "total")
	API_EXEC_TIMES_COUNTER_APP_ABNORMAL = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues(Ec2InstanceId, "app", "abnormal")
	API_EXEC_TIMES_COUNTER_APP_PANIC = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues(Ec2InstanceId, "app", "panic")
	API_EXEC_TIMES_COUNTER_MANAGE_TOTAL = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues(Ec2InstanceId, "manage", "total")
	API_EXEC_TIMES_COUNTER_MANAGE_ABNORMAL = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues(Ec2InstanceId, "manage", "abnormal")
	API_EXEC_TIMES_COUNTER_MANAGE_PANIC = API_EXEC_TIMES_COUNTER_VEC.WithLabelValues(Ec2InstanceId, "manage", "panic")

	// api返回码
	API_RETURN_CODE_COUNTER_VEC, err = metrics.NewCounterVec(
		MetricsNameSpace,
		"api",
		"return_code",
		"api返回码",
		[]string{"instance", "origin", "code"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	// api调用次数
	API_URI_CALL_TIMES_COUNTER_VEC, err = metrics.NewCounterVec(
		MetricsNameSpace,
		"api",
		"uri_call_times",
		"api uri调用次数",
		[]string{"instance", "origin", "uri"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	// api时间
	API_SPENT_TIME_SUMMARY_VEC, err = metrics.NewSummaryVec(
		MetricsNameSpace,
		"api",
		"spent_time",
		"api消耗时间",
		[]string{"instance", "origin"},
	)
	if nil != err {
		logrus.Fatal(err)
		return
	}

	API_SPENT_TIME_SUMMARY_APP = API_SPENT_TIME_SUMMARY_VEC.WithLabelValues(Ec2InstanceId, "app")
	API_SPENT_TIME_SUMMARY_MANAGE = API_SPENT_TIME_SUMMARY_VEC.WithLabelValues(Ec2InstanceId, "manage")

}
`
