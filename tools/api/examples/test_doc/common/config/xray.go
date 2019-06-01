package config

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-xray-sdk-go/xray"
	// Importing the plugins enables collection of AWS resource information at runtime.
	// Every plugin should be imported after "github.com/aws/aws-xray-sdk-go/xray" library.
	_ "github.com/aws/aws-xray-sdk-go/plugins/ec2"
)

type XRayConfigFormat struct {
	DaemonAddress  string `json:"daemon_address" yaml:"daemon_address" validate:"required"`
	LogLevel       string `json:"log_level" yaml:"log_level"`
	ServiceVersion string `json:"service_version" yaml:"service_version"`
}

var XRayConfig *XRayConfigFormat

func CheckXRayConfig() {
	if XRayConfig != nil {
		return
	}

	CheckConsulConfig()
	XRayConfig = &XRayConfigFormat{}
	var err error
	err = GetConsulClient().GetYaml(ConsulConfig.KeyPrefix+"/xray.yaml", XRayConfig)
	if nil != err {
		logrus.Fatalf("read xray config failed. %s", err)
		return
	}

	err = validator.New().Struct(XRayConfig)
	if nil != err {
		logrus.Fatalf("validate xray config failed. %s", err)
		return
	}

	err = xray.Configure(xray.Config{
		DaemonAddr:     XRayConfig.DaemonAddress,
		LogLevel:       XRayConfig.LogLevel,
		ServiceVersion: XRayConfig.ServiceVersion,
	})
	if nil != err {
		logrus.Fatalf("configure xray config failed. %s", err)
		return
	}
}
