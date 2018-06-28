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
