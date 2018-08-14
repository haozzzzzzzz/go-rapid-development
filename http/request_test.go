package http

import (
	"context"
	"testing"
)

func TestRequest_GetText(t *testing.T) {
	req, err := NewRequest("http://127.0.0.1:18100/home/metrics", context.Background(), RequestClient)
	if nil != err {
		t.Error(err)
		return
	}

	req.GetText()
}
