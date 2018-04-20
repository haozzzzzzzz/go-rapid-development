package api

import (
	"testing"
)

func TestCreateApi(t *testing.T) {
	apiItem := &ApiItem{}
	err := CreateApiSource("/Users/hao/Documents/Projects/Github/go_lambda_learning/src/github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder/api", apiItem)
	if nil != err {
		t.Error(err)
		return
	}
}
