package db

import (
	"database/sql/driver"
	"encoding/json"

	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StringSlice []string

// 转成db存储值
func (m StringSlice) Value() (value driver.Value, err error) {
	byteStringSlice, err := json.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal string slice failed. %s.", err)
		return
	}

	value = string(byteStringSlice)
	return
}

// 从db存储值中转
func (m *StringSlice) Scan(src interface{}) (err error) {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New(fmt.Sprintf("Incompatible type %q for StringSlice", reflect.TypeOf(src)))

	}

	err = json.Unmarshal([]byte(source), m)
	if nil != err {
		logrus.Errorf("unmarshal string slice failed. %s.", err)
		return
	}

	return
}

type Uint32Map map[string]uint32

// 转成db存储值
func (m Uint32Map) Value() (value driver.Value, err error) {
	byteValue, err := json.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal uint32 map failed. %s.", err)
		return
	}

	value = string(byteValue)
	return
}

func (m *Uint32Map) Scan(src interface{}) (err error) {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New(fmt.Sprintf("Incompatible type %q for Uint32Map", reflect.TypeOf(src)))
	}

	err = json.Unmarshal([]byte(source), m)
	if nil != err {
		logrus.Errorf("unmarshal string slice failed. %s.", err)
		return
	}

	return
}

type InterfaceMap map[string]interface{}

func (m InterfaceMap) Value() (value driver.Value, err error) {
	byteValue, err := json.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal interface map failed. %s.", err)
		return
	}

	value = string(byteValue)

	return
}

func (m *InterfaceMap) Scan(src interface{}) (err error) {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New(fmt.Sprintf("Incompatible type %q for InterfaceMap", reflect.TypeOf(src)))
	}

	err = json.Unmarshal([]byte(source), m)
	if nil != err {
		logrus.Errorf("unmarshal string slice failed. %s.", err)
		return
	}

	return
}
