package config

import (
	"service/common/session"

	"github.com/sirupsen/logrus"
)

func CheckSessionConfig() {
	var err error
	CheckAWSConfig()
	CheckServiceConfig()

	err = session.InitMetrics(ServiceConfig.MetricsNamespace, AWSEc2InstanceIdentifyDocument.InstanceID)
	if nil != err {
		logrus.Fatalf("init session metrics failed. error: %s.", err)
		return
	}
}
