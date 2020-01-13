package db

import (
	"fmt"
	"testing"
)

func TestJsonInterface(t *testing.T) {
	i := JsonInterface{map[string]interface{}{
		"name": "value",
	}}
	fmt.Println(i)

	j := JsonInterface{[]interface{}{"str1"}}
	fmt.Println(j)
}
