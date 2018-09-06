package project

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
)

func (m *ProjectSource) generateCommonMetrics() (err error) {
	metricsDir := fmt.Sprintf("%s/common/metrics", m.ProjectDir)
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
	MetricsNamespace string

	// ec2 instance id
	Ec2InstanceId string

	// 服务次数统计
	SERVICE_TIMES_COUNTER_VEC   *prometheus.CounterVec
	SERVICE_TIMES_COUNTER_APP_PANIC    prometheus.Counter // app service panic
	SERVICE_TIMES_COUNTER_MANAGE_PANIC prometheus.Counter // manage service panic
)

func init() {
	var err error

	MetricsNamespace = config.ServiceConfig.MetricsNamespace
	Ec2InstanceId = config.AWSEc2InstanceIdentifyDocument.InstanceID

	// system次数统计
	SERVICE_TIMES_COUNTER_VEC, err = metrics.NewCounterVec(
		MetricsNamespace,
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

}
`
