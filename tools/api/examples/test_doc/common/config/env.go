package config

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

type EnvConfigFormat struct {
	Debug bool   `json:"debug" yaml:"debug"`
	Stage string `json:"stage" yaml:"stage" validate:"required"`
}

func (m *EnvConfigFormat) WithStagePrefix(strVal string) string {
	return fmt.Sprintf("%s_%s", m.Stage, strVal)
}

var EnvConfig *EnvConfigFormat

func CheckEnvConfig() {
	if EnvConfig != nil {
		return
	}

	CheckConsulConfig()

	EnvConfig = &EnvConfigFormat{}
	var err error
	err = GetConsulClient().GetYaml(ConsulConfig.KeyPrefix+"/env.yaml", EnvConfig)
	if nil != err {
		logrus.Errorf("read env config file failed. %s.", err)
		return
	}

	err = validator.New().Struct(EnvConfig)
	if nil != err {
		logrus.Fatalf("validate env config failed. %s", err)
		return
	}

}
