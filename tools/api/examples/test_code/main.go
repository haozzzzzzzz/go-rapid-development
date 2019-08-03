package main

import (
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/tools/code/com/precompiler"
	"github.com/sirupsen/logrus"
)

func main() {
	filename := "/Users/hao/Documents/Projects/Github/go-rapid-development/tools/api/examples/test_code/src/src.go"
	newFileText, err := precompiler.PrecompileText(filename, nil, map[string]interface{}{
		"USE_IMPORT": true,
	})
	if nil != err {
		logrus.Errorf(" failed. error: %s.", err)
		return
	}
	fmt.Println(newFileText)
}
