package yaml

import (
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
)

func ReadYamlFromFile(filePath string, obj interface{}) (err error) {
	byteObj, err := ioutil.ReadFile(filePath)
	if nil != err {
		logrus.Errorf("read %q failed. error: %s.", filePath, err)
		return
	}

	err = yaml.Unmarshal(byteObj, obj)
	if nil != err {
		logrus.Errorf("unmarshal %q yaml file to %q failed. %s.", filePath, reflect.TypeOf(obj), err)
		return
	}
	return
}

func ReadValidateYamlFromFile(filePath string, obj interface{}) (err error) {
	err = ReadYamlFromFile(filePath, obj)
	if err != nil {
		logrus.Errorf("read yaml from file %s error: %s", filePath, err)
		return
	}
	err = validator.New().Struct(obj)
	if err != nil {
		logrus.Errorf("validate yaml obj error: %s", err)
		return
	}
	return
}

func WriteYamlToFile(filePath string, obj interface{}, mode os.FileMode) (err error) {
	byteObj, err := yaml.Marshal(obj)
	if nil != err {
		logrus.Errorf("marshal %q failed. %s.", reflect.TypeOf(obj), err)
		return
	}

	err = ioutil.WriteFile(filePath, byteObj, mode)
	if nil != err {
		logrus.Errorf("write %q failed. %s.", filePath, err)
		return
	}

	return
}
