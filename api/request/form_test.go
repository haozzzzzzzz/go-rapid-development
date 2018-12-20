package request

import (
	"fmt"
	"testing"
)

func TestFormMapStruct(t *testing.T) {
	form := make(map[string][]string)
	err := FormMapStruct(form, &struct {
		IntValues     []int     `json:"int_values" form:"int_values"`
		UintValues    []uint    `json:"uint_values" form:"uint_values"`
		Uint32Values  []uint32  `json:"uint32_values" form:"uint32_values"`
		BoolValues    []bool    `json:"bool_values" form:"bool_values"`
		Float32Values []float32 `json:"float32_values" form:"float32_values"`
		Uint32Value   uint32    `json:"uint_32_value" form:"uint_32_value"`
	}{
		IntValues:     []int{1, 2},
		UintValues:    []uint{1, 2},
		Uint32Values:  []uint32{1, 2, 3},
		BoolValues:    []bool{false, true},
		Float32Values: []float32{1.0, 2.0, 2.1},
		Uint32Value:   1,
	})
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(form)
}

func TestFormMap(t *testing.T) {
	var err error
	form := make(map[string][]string, 0)

	//err = FormMap(form, map[string]string{})
	//if nil != err {
	//	t.Error(err)
	//	return
	//}

	//var i = 1
	//err = FormMapStruct(form, &i)
	//if nil != err {
	//	t.Error(err)
	//	return
	//}

	err = FormMap(form, map[string]string{
		"hello": "world",
	})
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(form)

}
