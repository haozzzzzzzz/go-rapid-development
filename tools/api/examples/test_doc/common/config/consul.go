package config

import (
	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/consul"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
)

type ConsulConfigFormat struct {
	ClientConfig *consul.ClientConfigFormat `yaml:"client_config" validate:"required,required"`
	KeyPrefix    string                     `yaml:"key_prefix" validate:"required"`
}

var ConsulConfig *ConsulConfigFormat
var consulClient *consul.Client

func CheckConsulConfig() {
	if ConsulConfig != nil {
		return
	}

	logrus.Info("checking consul config")
	ConsulConfig = &ConsulConfigFormat{}
	var err error
	err = yaml.ReadYamlFromFile("./config/consul.yaml", ConsulConfig)
	if nil != err {
		logrus.Fatalf("read consul config failed. %s", err)
		return
	}

	err = validator.New().Struct(ConsulConfig)
	if nil != err {
		logrus.Fatalf("validate consul config failed. %s", err)
		return
	}

	consulClient, err = consul.NewClient(ConsulConfig.ClientConfig)
	if nil != err {
		logrus.Fatalf("new consult client failed. error: %s.", err)
		return
	}
}

func GetConsulClient() *consul.Client {
	CheckConsulConfig()
	return consulClient
}
