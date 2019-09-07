package config

import (
	"github.com/haozzzzzzzz/go-rapid-development/consul"
	"github.com/haozzzzzzzz/go-rapid-development/utils/ujson"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

func LoadFileYaml(path string, obj interface{}) (err error) {
	err = yaml.ReadYamlFromFile(path, obj)
	if nil != err {
		logrus.Errorf("read yaml from file failed. error: %s.", err)
		return
	}

	err = validator.New().Struct(obj)
	if nil != err {
		logrus.Errorf("validate yaml config failed. path: %s, %s", path, err)
		return
	}

	return
}

func LoadFileYamlPanic(path string, obj interface{}) {
	err := LoadFileYaml(path, obj)
	if nil != err {
		logrus.Panicf("load file yaml failed. path: %s, err: %s", path, err)
		return
	}
}

func LoadFileYamlOrDefaultPanic(
	path string,
	obj interface{},
	defVal interface{},
) {
	obj = defVal
	err := LoadFileYaml(path, obj)
	if nil != err {
		err = nil
		return
	}
}

func LoadFileJson(path string, obj interface{}) (err error) {
	err = ujson.ReadJsonFromFile(path, obj)
	if nil != err {
		logrus.Errorf("read yaml from file failed. error: %s.", err)
		return
	}

	err = validator.New().Struct(obj)
	if nil != err {
		logrus.Errorf("validate json config failed. path: %s, %s", path, err)
		return
	}

	return
}

func LoadFileJsonPanic(path string, obj interface{}) {
	err := LoadFileJson(path, obj)
	if nil != err {
		logrus.Panicf("load file json failed. path: %s, err: %s", path, err)
		return
	}
	return
}

type ConsulLoader struct {
	Client    *consul.Client
	address   string
	keyPrefix string
}

func NewConsulLoader(
	address string,
	keyPrefix string,
) (loader *ConsulLoader, err error) {
	loader = &ConsulLoader{
		address:   address,
		keyPrefix: keyPrefix,
	}
	loader.Client, err = consul.NewClient(&consul.ClientConfigFormat{
		Address: address,
	})
	if nil != err {
		logrus.Errorf("new consul loader failed. error: %s.", err)
		return
	}
	return
}

func NewConsulLoaderPanic(
	address string,
	keyPrefix string,
) (loader *ConsulLoader) {
	loader, err := NewConsulLoader(address, keyPrefix)
	if nil != err {
		logrus.Panicf("new consule loader failed. address: %s, key_prefix: %s, err: %s", address, keyPrefix, err)
		return
	}
	return
}

func (m *ConsulLoader) GetYaml(key string, obj interface{}) (err error) {
	key = m.keyPrefix + key
	err = m.Client.GetYaml(key, obj)
	if nil != err {
		logrus.Errorf("get yaml failed. key: %s, err: %s", key, err)
		return
	}

	err = validator.New().Struct(obj)
	if nil != err {
		logrus.Errorf("validate yaml failed. key: %s, err: %s", key, err)
		return
	}
	return
}

func (m *ConsulLoader) GetYamlPanic(key string, obj interface{}) {
	err := m.GetYaml(key, obj)
	if nil != err {
		logrus.Panicf("get yaml failed. key: %s, error: %s", key, err)
		return
	}
}

func (m *ConsulLoader) GetJson(key string, obj interface{}) (err error) {
	key = m.keyPrefix + key
	err = m.Client.GetJson(key, obj)
	if nil != err {
		logrus.Errorf("get json failed. key: %s, err: %s", key, err)
		return
	}

	err = validator.New().Struct(obj)
	if nil != err {
		logrus.Errorf("validate json failed. key: %s, err: %s", key, err)
		return
	}
	return
}

func (m *ConsulLoader) GetJsonPanic(key string, obj interface{}) {
	err := m.GetJson(key, obj)
	if nil != err {
		logrus.Panicf("get json failed. key: %s, err: %s", key, err)
		return
	}
}

func (m *ConsulLoader) GetString(key string) (value string, err error) {
	key = m.keyPrefix + key
	value, err = m.Client.GetString(key)
	if nil != err {
		logrus.Errorf("get string failed. key: %s, error: %s.", key, err)
		return
	}

	return
}

func (m *ConsulLoader) GetStringPanic(key string) (value string) {
	value, err := m.GetString(key)
	if nil != err {
		logrus.Panicf("get string failed. key: %s, error: %s", key, err)
		return
	}
	return
}
