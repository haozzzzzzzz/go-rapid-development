package es

import (
	"fmt"
	http2 "net/http"
	"strings"
	"testing"
)

type SampleChecker struct {
}

func (m *SampleChecker) Before(req *http2.Request) {
	fmt.Println("before")
}

func (m *SampleChecker) After(resp *http2.Response, err error) {
	fmt.Println("after")
}

func TestNewClient(t *testing.T) {
	var err error
	client, err := NewClient([]string{"http://localhost:9200"}, func() RoundTripChecker {
		return &SampleChecker{}
	})
	if nil != err {
		t.Error(err)
		return
	}

	//// get cluster info
	//resp, err := client.Info()
	//if nil != err {
	//	t.Error(err)
	//	return
	//}
	//
	//if resp.IsError() {
	//	fmt.Printf("error: %s", resp.String())
	//}
	//
	//r := make(map[string]interface{})
	//err = json.NewDecoder(resp.Body).Decode(&r)
	//if nil != err {
	//	t.Error(err)
	//	return
	//}
	//err = resp.Body.Close()
	//if nil != err {
	//	t.Error(err)
	//	return
	//}

	//resp, err = client.Index(
	//	"customer",
	//	strings.NewReader("{\"name\":\"luohao\"}"),
	//)
	//if nil != err {
	//	t.Error(err)
	//	return
	//}
	//defer func() {
	//	err = resp.Body.Close()
	//	if nil != err {
	//		t.Error(err)
	//		return
	//	}
	//}()

	//bulkIndexBody, err := BulkBodyFromLines([]interface{}{
	//	map[string]interface{}{
	//		"index": map[string]interface{}{},
	//	},
	//	map[string]interface{}{
	//		"name": "luohao1",
	//	},
	//})
	//if nil != err {
	//	t.Error(err)
	//	return
	//}

	bulkIndexBody, err := BulkBodyFromIdDocMap("create", "customer", map[string]interface{}{
		"1": map[string]interface{}{
			"name": "l2",
		},
	})

	fmt.Println(bulkIndexBody)
	resp, err := client.Bulk(strings.NewReader(bulkIndexBody),
		client.Bulk.WithIndex("customer"),
	)
	if nil != err {
		t.Error()
		return
	}
	defer func() {
		err = resp.Body.Close()
		if nil != err {
			t.Error(err)
			return
		}
	}()
	fmt.Println(resp)

	//resp, err = client.Search()
}
