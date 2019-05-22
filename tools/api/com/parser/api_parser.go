package parser

import (
	"fmt"
	"io/ioutil"

	"sort"

	"os"

	"strings"

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

	// binding routers
	// add imports
	importsMap := make(map[string]string) // package -> alias
	_ = importsMap

	goPath := os.Getenv("GOPATH")
	goPaths := strings.Split(goPath, ":")
	for _, subGoPath := range goPaths {
		for _, apiItem := range apis {
			if strings.Contains(apiItem.PackagePath, subGoPath) {
				relPath := strings.Replace(apiItem.PackagePath, subGoPath+"/src/", "", -1)
				if relPath != "" {
					importsMap[relPath] = apiItem.RelativePackage
				}
			}
		}
	}

	strImports := make([]string, 0)
	for relPath, alias := range importsMap {
		var str string
		if alias == "" {
			str = fmt.Sprintf("    %q", relPath)
		} else {
			str = fmt.Sprintf("    %s %q", alias, relPath)
		}
		strImports = append(strImports, str)
	}

	strRouters := make([]string, 0)
	for _, apiItem := range apis {
		strHandleFunc := apiItem.ApiHandlerFunc
		if apiItem.RelativePackage != "" {
			strHandleFunc = fmt.Sprintf("%s.%s", apiItem.RelativePackage, apiItem.ApiHandlerFunc)
		}

		str := fmt.Sprintf("    engine.Handle(\"%s\", \"%s\", %s.GinHandler)", apiItem.HttpMethod, apiItem.RelativePath, strHandleFunc)
		strRouters = append(strRouters, str)

	}

	routersFileName := fmt.Sprintf("%s/routers.go", m.ApiDir)
	newRoutersText := fmt.Sprintf(routersFileText, strings.Join(strImports, "\n"), strings.Join(strRouters, "\n"))
	err = ioutil.WriteFile(routersFileName, []byte(newRoutersText), 0644)
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

var routersFileText = `package api
import (
"github.com/gin-gonic/gin"
%s
)

// 注意：BindRouters函数体内不能自定义添加任何声明，由api compile命令生成api绑定声明
func BindRouters(engine *gin.Engine) (err error) {
%s
	return
}
`
