package config

import (
	"github.com/haozzzzzzzz/go-rapid-development/utils/ujson"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

func LoadFileYamlPanic(skip bool, path string, obj interface{}) {
	if skip {
		return
	}

	err := yaml.ReadYamlFromFile(path, obj)
	if nil != err {
		logrus.Panicf("read yaml from file failed. error: %s.", err)
		return
	}

	err = validator.New().Struct(obj)
	if nil != err {
		logrus.Panicf("validate yaml config failed. path: %s, %s", path, err)
		return
	}

	return
}

func LoadFileJsonPanic(skip bool, path string, obj interface{}) {
	if skip {
		return
	}

	err := ujson.ReadJsonFromFile(path, obj)
	if nil != err {
		logrus.Panicf("read yaml from file failed. error: %s.", err)
		return
	}

	err = validator.New().Struct(obj)
	if nil != err {
		logrus.Panicf("validate json config failed. path: %s, %s", path, err)
		return
	}

	return
}
