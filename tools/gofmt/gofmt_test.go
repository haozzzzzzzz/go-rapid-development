package gofmt

import (
	"testing"
)

func TestGoFmt(t *testing.T) {
	err := GoFmt("/Users/hao/Documents/Projects/Github/go-rapid-development/tools/gofmt/g.go")
	if nil != err {
		t.Error(err)
		return
	}

}
