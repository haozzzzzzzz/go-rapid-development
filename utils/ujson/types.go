package ujson

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
)

// copy from encoding/json scanner.go
type JsonType int

const (
	JsonTypeUnknown JsonType = iota
	JsonTypeEmpty
	JsonTypeBoolean
	JsonTypeNumbers
	JsonTypeString
	JsonTypeArray
	JsonTypeObject
	JsonTypeNull
)

func GetStrJsonType(strJson string) (jsonType JsonType) {
	strJson = strings.TrimSpace(strJson)
	if strJson == "" {
		jsonType = JsonTypeEmpty
		return
	}

	firstC := strJson[0]
	switch firstC {
	case '{':
		jsonType = JsonTypeObject
	case '[':
		jsonType = JsonTypeArray
	case '"':
		jsonType = JsonTypeString
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		jsonType = JsonTypeNumbers
	case 't', 'f':
		jsonType = JsonTypeBoolean
	case 'n':
		jsonType = JsonTypeNull
	default:
		jsonType = JsonTypeUnknown
	}

	return
}

func ParseJson(strJson string) (jsonType JsonType, jsonObj interface{}, err error) {
	jsonType = GetStrJsonType(strJson)
	switch jsonType {
	case JsonTypeObject:
		jsonObj = map[string]interface{}{}

	case JsonTypeArray:
		jsonObj = []interface{}{}

	case JsonTypeNumbers:
		jsonObj = float64(0)

	case JsonTypeBoolean:
		jsonObj = false

	case JsonTypeNull, JsonTypeEmpty, JsonTypeUnknown:
		jsonObj = nil

	default:
		err = errors.New("unknown json type")

	}

	if jsonObj == nil {
		return
	}

	err = json.Unmarshal([]byte(strJson), &jsonObj)
	if nil != err {
		logrus.Errorf("unmarshal json obj failed. error: %s.", err)
		return
	}

	return
}
