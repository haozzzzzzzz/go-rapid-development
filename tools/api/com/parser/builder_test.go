package parser

import (
	"testing"
)

func TestCreateApiSource(t *testing.T) {
	err := CreateApiSource(&ApiItem{
		ApiHandlerFunc: "Func",
		PackageName:    "api",
		HttpMethod:     "GET",
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
