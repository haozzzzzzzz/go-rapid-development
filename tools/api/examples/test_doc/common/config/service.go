package config

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

type ServiceConfigFormat struct {
	MetricsNamespace string `json:"metrics_namespace" yaml:"metrics_namespace" validate:"required"`
}

var ServiceConfig *ServiceConfigFormat

func CheckServiceConfig() {
	if ServiceConfig != nil {
		return
	}

	CheckConsulConfig()

	ServiceConfig = &ServiceConfigFormat{}
	var err error
	err = GetConsulClient().GetYaml(ConsulConfig.KeyPrefix+"/service.yaml", ServiceConfig)
	if nil != err {
		logrus.Fatalf("read service config file failed. %s", err)
		return
	}

	err = validator.New().Struct(ServiceConfig)
	if nil != err {
		logrus.Errorf("validate service config failed. %s.", err)
		return
	}

}
