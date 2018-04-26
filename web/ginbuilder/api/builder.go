package api

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func CreateApiSource(apiItem *ApiItem) (err error) {
	newApiText := fmt.Sprintf(apiText,
		apiItem.ApiHandlerPackage,
		apiItem.ApiHandlerFunc,
		apiItem.HttpMethod,
		apiItem.RelativePath)

	err = ioutil.WriteFile(apiItem.SourceFile, []byte(newApiText), os.ModePerm)
	if nil != err {
		logrus.Errorf("write api file failed. %s", err)
		return
	}

	return
}

var apiText = `package %s

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

var %s ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod:   "%s",
	RelativePath: "%s",
	Handle: func(ctx *ginbuilder.Context) (err error) {
		ctx.SuccessReturn(map[string]interface{}{
			"info": "hello, world",
		})
		return
	},
}`
