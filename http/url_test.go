package http

import (
	"testing"

	"fmt"
)

func TestString(t *testing.T) {
	var err error

	tUrl, err := NewUrlByStrUrl("http://123456.com")
	if nil != err {
		t.Error(err)
		return
	}

	tUrl.QueryValues.Set("name", "Ricky Chen")
	fmt.Println(tUrl.String())
}
