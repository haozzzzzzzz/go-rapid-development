package parser

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegExp(t *testing.T) {
	docReg, err := regexp.Compile(`(?si:@api_doc_start(.*?)@api_doc_end)`)
	if nil != err {
		t.Error(err)
		return
	}

	strs := docReg.FindAllStringSubmatch(`
/**
@api_doc_start
	{
		"hello": "world"
	}
@api_doc_end

@api_doc_start
	{
		"hello": "world2"
	}
@api_doc_end

@api_doc_start@api_doc_end

**/
`, -1)

	strJsons := make([]string, 0)
	for _, str := range strs {
		strJsons = append(strJsons, str[1])
	}

	for _, strJson := range strJsons {
		fmt.Println(strJson)
	}
}
