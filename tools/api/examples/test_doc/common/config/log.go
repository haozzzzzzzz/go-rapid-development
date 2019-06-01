package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/common/dependent"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type LogOutputConfigFormat struct {
	Filedir    string `json:"filedir" yaml:"filedir" validate:"required"`
	MaxSize    int    `json:"max_size" yaml:"max_size" validate:"required"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups" validate:"required"`
	MaxAge     int    `json:"max_age" yaml:"max_age" validate:"required"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

type LogConfigFormat struct {
	LogLevel logrus.Level           `json:"log_level" yaml:"log_level" validate:"required"`
	Output   *LogOutputConfigFormat `json:"output" yaml:"output"`
}

var LogConfig *LogConfigFormat

func CheckLogConfig() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	newLogger := logrus.New()
	newLogger.Formatter = &logrus.JSONFormatter{}

	if LogConfig != nil {
		return
	}

	LogConfig = &LogConfigFormat{}
	var err error
	err = yaml.ReadYamlFromFile("./config/log.yaml", LogConfig)
	if nil != err {
		logrus.Fatalf("read log config file failed. %s", err)
		return
	}

	err = validator.New().Struct(LogConfig)
	if nil != err {
		logrus.Fatalf("validate log config failed. %s", err)
		return
	}

	logrus.SetLevel(LogConfig.LogLevel)
	if LogConfig.Output != nil && dependent.ServiceName != "" {
		filename := fmt.Sprintf("%s/%s/info.log", LogConfig.Output.Filedir, dependent.ServiceName)
		logger := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    LogConfig.Output.MaxSize,
			MaxBackups: LogConfig.Output.MaxBackups,
			MaxAge:     LogConfig.Output.MaxAge,
			Compress:   LogConfig.Output.Compress,
		}
		logrus.SetOutput(logger)
		log.SetOutput(logger)

	} else {
		log.SetOutput(newLogger.Writer())

	}

}
