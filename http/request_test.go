package http

import (
	"context"
	"fmt"
	"testing"
)

func TestRequest_GetText(t *testing.T) {
	req, err := NewRequest("http://127.0.0.1:18101/home_manage/metrics", context.Background(), RequestClient)
	if nil != err {
		t.Error(err)
		return
	}

	txt, err := req.GetText()
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(txt)
}

func TestRequest_GetJSON(t *testing.T) {
	req, err := NewRequest("http://127.0.0.1:18101/home_manage/v3/tab/feed/list", context.Background(), RequestClient )
	if nil != err {
		t.Error(err)
		return
	}

	m := make(map[string]interface{})
	err = req.GetJSON(&m)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(m)
}

func TestRequest_PostJson(t *testing.T) {
	req, err := NewRequest("http://127.0.0.1:18101/home_manage/v3/tab/feed/add", context.Background(), RequestClient )
	if nil != err {
		t.Error(err)
		return
	}

	m := make(map[string]interface{})
	err = req.PostJson(map[string]interface{}{}, &m)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(m)
}