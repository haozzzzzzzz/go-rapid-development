package config

import (
	"service/common/proxy/operations_api"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

type RemoteApiConfigFormat struct {
	OperationsRpcApiUrl string `yaml:"operations_rpc_api_url" validate:"required"`
}

var RemoteApiConfig *RemoteApiConfigFormat

func CheckRemoteApiConfig() {
	var err error
	if RemoteApiConfig != nil {
		return
	}

	RemoteApiConfig = &RemoteApiConfigFormat{}
	err = GetConsulClient().GetYaml(ConsulConfig.KeyPrefix+"/remote_api.yaml", RemoteApiConfig)
	if nil != err {
		logrus.Fatalf("read remote api config file failed. %s", err)
		return
	}

	err = validator.New().Struct(RemoteApiConfig)
	if nil != err {
		logrus.Fatalf("validate remote api config failed. %s", err)
		return
	}

	err = operations_api.SetApiUrlPrefix(RemoteApiConfig.OperationsRpcApiUrl)
	if nil != err {
		logrus.Fatalf("set operations api url failed. %s", err)
		return
	}

}
