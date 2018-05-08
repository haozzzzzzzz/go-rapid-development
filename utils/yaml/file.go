package yaml

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func ReadYamlFromFile(filePath string, obj interface{}) (err error) {
	byteObj, err := ioutil.ReadFile(filePath)
	if nil != err {
		logrus.Errorf("read %q failed. \n%s.", filePath, err)
		return
	}

	err = yaml.Unmarshal(byteObj, obj)
	if nil != err {
		logrus.Errorf("unmarshal %q yaml file to %q failed. \n%s.", filePath, reflect.TypeOf(obj), err)
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
