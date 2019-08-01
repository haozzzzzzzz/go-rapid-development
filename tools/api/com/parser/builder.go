package parser

import (
	"os"
	"text/template"

	"bytes"

	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func CreateApiSource(apiItem *ApiItem) (err error) {
	templ, err := template.New(apiItem.ApiHandlerFunc).Parse(apiText)
	if nil != err {
		logrus.Errorf("new api source template failed failed. error: %s.", err)
		return
	}

	buf := bytes.NewBuffer(nil)
	err = templ.Execute(buf, apiItem)
	if nil != err {
		logrus.Errorf("execute template failed. error: %s.", err)
		return
	}

	err = ioutil.WriteFile(apiItem.SourceFile, buf.Bytes(), os.ModePerm)
	if nil != err {
		logrus.Errorf("write api file failed. %s", err)
		return
	}

	return
}

var apiText = `package {{.PackageName}}

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/web/ginbuilder"
)

var {{.ApiHandlerFunc}} ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod:   "{{.HttpMethod}}",
	RelativePaths: []string{
		{{range .RelativePaths}}{{.}},
		{{end}}
	},
	Handle: func(ctx *ginbuilder.Context) (err error) {
		ctx.SuccessReturn(map[string]interface{}{
			"info": "hello, world",
		})
		return
	},
}`
