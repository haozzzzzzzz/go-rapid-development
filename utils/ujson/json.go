package ujson

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"reflect"

	"github.com/sirupsen/logrus"
)

func UnmarshalJsonFromReader(reader io.Reader, v interface{}) error {

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		logrus.Errorf("unmarshal json from reader's content failed. content: %s. error: %s.", string(bytes), err)
		return err
	}

	return nil
}

func ReadJsonFromFile(filePath string, obj interface{}) (err error) {
	byteObj, err := ioutil.ReadFile(filePath)
	if nil != err {
		logrus.Errorf("read %q json file failed. error: %s.", filePath, err)
		return
	}

	err = json.Unmarshal(byteObj, obj)
	if nil != err {
		logrus.Errorf("unmarshal %q json file failed. error: %s.", filePath, err)
		return
	}

	return
}

func WriteJsonToFile(filePath string, obj interface{}, mode os.FileMode) (err error) {
	byteObj, err := json.Marshal(obj)
	if nil != err {
		logrus.Errorf("marshal %q  failed. error: %s.", reflect.TypeOf(obj), err)
		return
	}

	err = ioutil.WriteFile(filePath, byteObj, mode)
	if nil != err {
		logrus.Errorf("write %q failed. error: %s.", filePath, err)
		return
	}

	return
}

func MarshalPretty(obj interface{}) ([]byte, error){
	return json.MarshalIndent(obj, "", "\t")
}
