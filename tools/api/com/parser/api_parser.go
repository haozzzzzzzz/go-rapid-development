package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"go/ast"
	"go/parser"
	"go/token"
	"sort"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ApiParser struct {
	Service    *project.Service
	ServiceDir string
	ApiDir     string
}

func NewApiParser(service *project.Service) *ApiParser {
	return &ApiParser{
		Service:    service,
		ServiceDir: service.Config.ServiceDir,
		ApiDir:     fmt.Sprintf("%s/api", service.Config.ServiceDir),
	}
}

func (m *ApiParser) ParseRouter() (err error) {
	apis, err := m.ScanApis()
	if nil != err {
		logrus.Errorf("Scan api failed. \n%s.", err)
		return
	}

	err = m.MapApi(apis)
	if nil != err {
		logrus.Errorf("mapping api failed. %s.", err)
		return
	}

	return
}

func (m *ApiParser) MapApi(apis []*ApiItem) (err error) {
	// sort api
	apiUriKeys := make([]string, 0)
	mapApi := make(map[string]*ApiItem)
	for _, oneApi := range apis {
		uriKey := fmt.Sprintf("%s_%s", oneApi.RelativePath, oneApi.HttpMethod)
		apiUriKeys = append(apiUriKeys, uriKey)
		mapApi[uriKey] = oneApi
	}

	sort.Strings(apiUriKeys)

	// new order apis
	apis = make([]*ApiItem, 0)
	for _, apiUriKey := range apiUriKeys {
		apis = append(apis, mapApi[apiUriKey])
	}

	logrus.Info("Mapping apis ...")
	defer func() {
		if err == nil {
			logrus.Info("Map apis completed")
		}
	}()

	newRoutersFileMiddle := bytes.NewBuffer(nil)
	for _, apiItem := range apis {
		strHandleFunc := apiItem.ApiHandlerFunc
		if apiItem.RelativePackage != "" {
			strHandleFunc = fmt.Sprintf("%s.%s", apiItem.RelativePackage, apiItem.ApiHandlerFunc)
		}

		str := fmt.Sprintf("    engine.Handle(\"%s\", \"%s\", %s.GinHandler)\n", apiItem.HttpMethod, apiItem.RelativePath, strHandleFunc)
		newRoutersFileMiddle.Write([]byte(str))
	}

	routersFileName := fmt.Sprintf("%s/routers.go", m.ApiDir)
	fileTokenSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fileTokenSet, routersFileName, nil, parser.AllErrors)
	if nil != err {
		logrus.Errorf("parse routers.go failed. \n%s.", err)
		return
	}

	var lBrace, rBrace token.Pos
	for _, scopeObject := range astFile.Scope.Objects {
		if scopeObject.Name != "BindRouters" {
			continue
		}

		bindRoutersFuncDecl, ok := scopeObject.Decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		lBrace = bindRoutersFuncDecl.Body.Lbrace
		rBrace = bindRoutersFuncDecl.Body.Rbrace

	}

	lBracePos := fileTokenSet.Position(lBrace)
	rBracePos := fileTokenSet.Position(rBrace)

	byteRoutersFile, err := ioutil.ReadFile(routersFileName)
	if nil != err {
		logrus.Warnf("read routers file failed. \n%s.", err)
		return
	}

	newRoutersFileLeft := byteRoutersFile[0:lBracePos.Offset]
	newRoutersFileRight := byteRoutersFile[rBracePos.Offset+1:]

	newRoutersFileText := fmt.Sprintf(`%s {
%s    return
}%s`, string(newRoutersFileLeft), newRoutersFileMiddle.String(), string(newRoutersFileRight))

	err = ioutil.WriteFile(routersFileName, []byte(newRoutersFileText), 0644)
	if nil != err {
		logrus.Errorf("write new routers file failed. \n%s.", err)
		return
	}

	logrus.Info("Doing go imports")
	// do goimports
	goimports.DoGoImports([]string{m.ApiDir}, true)
	logrus.Info("Do go imports completed")

	// save api.yaml
	byteYamlApis, err := yaml.Marshal(apis)
	if nil != err {
		logrus.Errorf("yaml marshal apis failed. \n%s.", byteYamlApis)
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.service/apis.yaml", m.ServiceDir), byteYamlApis, 0644)
	if nil != err {
		logrus.Errorf("write apis.yaml failed. \n%s.", err)
		return
	}

	return
}
