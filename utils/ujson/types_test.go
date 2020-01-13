package ujson

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseJson(t *testing.T) {
	strJson := "[]"
	jsonType, jsonObj, err := ParseJson(strJson)
	if nil != err {
		t.Error(err)
		return
	}
	fmt.Printf("%d, %#v, %s\n", jsonType, jsonObj, reflect.TypeOf(jsonObj))
}
