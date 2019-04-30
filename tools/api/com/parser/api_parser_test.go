package parser

import (
	"fmt"
	"testing"
)

func TestApiParser_ParseApiFile(t *testing.T) {
	apis, err := ParseApiFile(
		"/Users/hao/Documents/Projects/XunLei/video_buddy_service/src/github.com/haozzzzzzzz/go-rapid-development/tools/api/com/parser/temp",
		"/Users/hao/Documents/Projects/XunLei/video_buddy_service/src/github.com/haozzzzzzzz/go-rapid-development/tools/api/com/parser/temp/api/api_temp.go",
	)
	if nil != err {
		t.Error(err)
		return
	}

	_ = apis
	for _, item := range apis {
		fmt.Printf("%#v\n", item.PathData)
		fmt.Printf("%#v\n", item.QueryData)
		fmt.Printf("%#v\n", item.PostData)
		fmt.Printf("%#v\n", item.RespData)
	}
}
