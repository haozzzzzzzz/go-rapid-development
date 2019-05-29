package parser

import (
	"testing"
)

func TestCreateApiSource(t *testing.T) {
	err := CreateApiSource(&ApiItem{
		ApiHandlerFunc:    "Func",
		ApiHandlerPackage: "api",
		HttpMethod:        "GET",
		RelativePaths: []string{
			"Hello",
			"world",
		},
	})
	if nil != err {
		t.Error(err)
		return
	}
}
