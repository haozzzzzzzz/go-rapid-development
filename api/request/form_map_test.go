package request

import (
	"fmt"
	"testing"
)

//func TestFormMapStruct(t *testing.T) {
//	form := make(map[string][]string)
//	err := FormMapStruct(form, &struct {
//		IntValues     []int     `json:"int_values" form:"int_values"`
//		UintValues    []uint    `json:"uint_values" form:"uint_values"`
//		Uint32Values  []uint32  `json:"uint32_values" form:"uint32_values"`
//		BoolValues    []bool    `json:"bool_values" form:"bool_values"`
//		Float32Values []float32 `json:"float32_values" form:"float32_values"`
//		Uint32Value   uint32    `json:"uint_32_value" form:"uint_32_value"`
//	}{
//		IntValues:     []int{1, 2},
//		UintValues:    []uint{1, 2},
//		Uint32Values:  []uint32{1, 2, 3},
//		BoolValues:    []bool{false, true},
//		Float32Values: []float32{1.0, 2.0, 2.1},
//		Uint32Value:   1,
//	})
//	if nil != err {
//		t.Error(err)
//		return
//	}
//
//	fmt.Println(form)
//}

func TestFormMap(t *testing.T) {
	var err error
	form := make(map[string][]string, 0)
	data := &Data{
		Key:   "456",
		Value: "789",
		Values: []string{
			"1234",
			"4567",
		},
		ArrayValues: []string{
			"abcd",
			"efg",
		},
		CustomBasicType: CustomBasicTypeValue,
	}

	//m := map[string]interface{}{
	//	"h": 123,
	//	"i": "hello",
	//	"a": []int{8, 9, 10},
	//}

	err = ObjToUrlValues(form, data)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(form)

}

type CustomBasicType int

const CustomBasicTypeValue CustomBasicType = 1

type Data struct {
	Key             string          `json:"key" form:"key" binding:"required"`
	Value           string          `json:"value" form:"value" binding:"required"`
	Values          []string        `json:"values" form:"values" binding:"required"`
	ArrayValues     []string        `json:"array_values" form:"array_values"`
	CustomBasicType CustomBasicType `json:"custom_basic_type" form:"custom_basic_type"`
}

func TestStructMapForm(t *testing.T) {
	m := map[string][]string{
		"key":   []string{"lh"},
		"value": []string{"123"},
	}
	data := &Data{}

	err := StructMapForm(data, m, "form")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(data)

}
