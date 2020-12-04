package request

import (
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/xiam/to"
	"net/url"
	"reflect"
)

// 将form字段解析成struct值
func MapForm(i interface{}, form map[string][]string) (err error) {
	return mapForm(i, form)
}

// 将struct值解析成form
// 目前用于解析uriData和queryData
// 只有简单类型和简单类型的数组才会被解析
func ObjToUrlValues(form url.Values, obj interface{}) (err error) {
	mForm := make(map[string]interface{})
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &mForm,
		TagName:  "form",
	})
	if err != nil {
		logrus.Errorf("new mapstructure decoder failed. error: %s", err)
		return
	}

	err = decoder.Decode(obj)
	if err != nil {
		logrus.Errorf("mapstructure decode struct to map failed. error: %s", err)
		return
	}

	for key, value := range mForm {
		values, ok := form[key]
		if !ok {
			values = make([]string, 0)
		}

		valType := reflect.TypeOf(value)
		valVal := reflect.ValueOf(value)
		kind := valType.Kind()
		switch kind {
		case reflect.Slice, reflect.Array:
			itemLen := valVal.Len()
			for i := 0; i < itemLen; i++ {
				itemVal := valVal.Index(i)
				strVal, errV := valueToString(itemVal.Interface())
				err = errV
				if err != nil {
					logrus.Errorf("slice url value to string failed. error: %s", err)
					return
				}
				values = append(values, strVal)
			}

		default:
			strVal, errV := valueToString(value)
			err = errV
			if err != nil {
				logrus.Errorf("url value to string failed. error: %s", err)
				return
			}
			values = append(values, strVal)
		}

		form[key] = values
	}

	return
}

// 将map值转成struct字段值
func StructMapForm(ptr interface{}, form map[string][]string, tag string) (err error) {
	return mapFormByTag(ptr, form, tag)
}

func valueToString(val interface{}) (ret string, err error) {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool, reflect.Float32, reflect.Float64, reflect.Complex128, reflect.Complex64, reflect.String:
		ret = to.String(val)

	default:
		err = uerrors.Newf("unsupported type for url.Values. type: %s", reflect.TypeOf(val))

	}

	return
}
