package es

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	var err error
	client, err := NewClient([]string{"http://localhost:9200"}, nil)
	if nil != err {
		t.Error(err)
		return
	}

	// get cluster info
	res, err := client.Info()
	if nil != err {
		t.Error(err)
		return
	}

	if res.IsError() {
		fmt.Printf("error: %s", res.String())
	}

	r := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&r)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(r)
}
