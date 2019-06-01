package config

import (
	"service/common/http"

	"github.com/sirupsen/logrus"
)

func CheckHttpConfig() {
	var err error
	CheckAWSConfig()
	CheckServiceConfig()

	err = http.InitMetrics(ServiceConfig.MetricsNamespace, AWSEc2InstanceIdentifyDocument.InstanceID)
	if nil != err {
		logrus.Fatalf("init http remote request metrics failed. %s", err)
		return
	}
}
