package db_sqlx

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/v2/utils/ujson"
	"reflect"

	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// NULL will present as emtpy string
type AdaptNullString string

func (m AdaptNullString) Value() (value driver.Value, err error) {
	if m == "" {
		return nil, nil
	}

	value = string(m)
	return
}

func (m *AdaptNullString) Scan(src interface{}) (err error) {
	ns := sql.NullString{}
	err = ns.Scan(src)
	if nil != err {
		logrus.Errorf("scan sql null string failed. error: %s.", err)
		return
	}

	*m = AdaptNullString(ns.String)
	return
}

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

// float64 map
type Float64Map map[string]float64

func (m Float64Map) Value() (value driver.Value, err error) {
	byteValue, err := json.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal uint32 map failed. %s.", err)
		return
	}

	value = string(byteValue)
	return
}

func (m *Float64Map) Scan(src interface{}) (err error) {
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

// interface map
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

	if source == "" {
		source = "null"
	}

	err = json.Unmarshal([]byte(source), m)
	if nil != err {
		logrus.Errorf("unmarshal string slice failed. %s.", err)
		return
	}

	return
}

// string map
type StringMap map[string]string

func (m StringMap) Value() (value driver.Value, err error) {
	byteValue, err := json.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal interface map failed. %s.", err)
		return
	}

	value = string(byteValue)

	return
}

func (m *StringMap) Scan(src interface{}) (err error) {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New(fmt.Sprintf("Incompatible type %q for StringMap", reflect.TypeOf(src)))
	}

	if source == "" {
		source = "null"
	}

	err = json.Unmarshal([]byte(source), m)
	if nil != err {
		logrus.Errorf("unmarshal string slice failed. %s.", err)
		return
	}

	return
}

type Uint32Slice []uint32

func (m Uint32Slice) Value() (value driver.Value, err error) {
	byteValue, err := json.Marshal(m)
	if nil != err {
		logrus.Errorf("marshal uint32 map failed. %s.", err)
		return
	}

	value = string(byteValue)
	return
}

func (m *Uint32Slice) Scan(src interface{}) (err error) {
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

type SplitStringSlice []string

func (m SplitStringSlice) Value() (value driver.Value, err error) {
	value = strings.Join(m, ",")
	return
}

func (m *SplitStringSlice) Scan(src interface{}) (err error) {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		err = errors.New(fmt.Sprintf("Incompatible type %q for SplitStringSlice", reflect.TypeOf(src)))
		return
	}

	if source != "" {
		*m = strings.Split(source, ",")
	} else {
		*m = SplitStringSlice{}
	}
	return
}

// json interface
type JsonInterface struct {
	I interface{}
}

func (m JsonInterface) Value() (value driver.Value, err error) {
	byteValue, err := json.Marshal(m.I)
	if nil != err {
		logrus.Errorf("marshal interface map failed. %s.", err)
		return
	}

	value = string(byteValue)

	return
}

func (m *JsonInterface) Scan(src interface{}) (err error) {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New(fmt.Sprintf("Incompatible type %q for InterfaceMap", reflect.TypeOf(src)))
	}

	if source == "" {
		source = "null"
	}

	_, m.I, err = ujson.ParseJson(source)
	if nil != err {
		logrus.Errorf("parse json failed. error: %s.", err)
		return
	}

	return
}
