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
	Namespace string

	// ec2 instance id
	Ec2InstanceId string

)

func init() {
	var err error

	Namespace = config.ServiceConfig.Namespace
	Ec2InstanceId = config.AWSEc2InstanceIdentifyDocument.InstanceID
}
`
