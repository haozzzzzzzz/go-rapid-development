package parser

import (
	"fmt"
	"io/ioutil"

	"sort"

	"strings"

	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ApiParser struct {
	Service    *project.Service
	ServiceDir string
	ApiDir     string
	GoPaths    []string
}

func NewApiParser(service *project.Service) *ApiParser {
	goPath := os.Getenv("GOPATH")
	return &ApiParser{
		Service:    service,
		ServiceDir: service.Config.ServiceDir,
		ApiDir:     fmt.Sprintf("%s/api", service.Config.ServiceDir),
		GoPaths:    strings.Split(goPath, ":"),
	}
}

func (m *ApiParser) ParseRouter(
	parseRequestData bool,
	importSource bool,
) (err error) {
	apis, err := m.ScanApis(parseRequestData, importSource)
	if nil != err {
		logrus.Errorf("Scan api failed. \n%s.", err)
		return
	}

	err = m.GenerateRoutersSourceFile(apis)
	if nil != err {
		logrus.Errorf("map api failed. %s.", err)
		return
	}

	return
}

func (m *ApiParser) apiUrlKey(uri string, method string) string {
	return fmt.Sprintf("%s_%s", uri, method)
}

func (m *ApiParser) GenerateRoutersSourceFile(apis []*ApiItem) (err error) {
	// sort api
	sortedApiUriKeys := make([]string, 0)
	mapApi := make(map[string]*ApiItem)
	for _, oneApi := range apis {
		if oneApi.RelativePaths == nil {
			continue
		}

		for _, relPath := range oneApi.RelativePaths {
			uriKey := m.apiUrlKey(relPath, oneApi.HttpMethod)
			sortedApiUriKeys = append(sortedApiUriKeys, uriKey)
			mapApi[uriKey] = oneApi
		}

	}

	sort.Strings(sortedApiUriKeys)

	logrus.Info("Map apis ...")
	defer func() {
		if err == nil {
			logrus.Info("Map apis completed")
		}
	}()

	// binding routers
	// add imports
	importsMap := make(map[string]string) // package -> alias
	_ = importsMap

	goPaths := m.GoPaths
	for _, subGoPath := range goPaths {
		for _, apiItem := range mapApi {
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
	for _, uriKey := range sortedApiUriKeys {
		apiItem := mapApi[uriKey]
		strHandleFunc := apiItem.ApiHandlerFunc
		if apiItem.RelativePackage != "" {
			strHandleFunc = fmt.Sprintf("%s.%s", apiItem.RelativePackage, apiItem.ApiHandlerFunc)
		}

		for _, uri := range apiItem.RelativePaths {
			if uriKey != m.apiUrlKey(uri, apiItem.HttpMethod) {
				continue
			}

			str := fmt.Sprintf("    engine.Handle(\"%s\", \"%s\", %s.GinHandler)", apiItem.HttpMethod, uri, strHandleFunc)
			strRouters = append(strRouters, str)
		}

	}

	routersFileName := fmt.Sprintf("%s/routers.go", m.ApiDir)
	newRoutersText := fmt.Sprintf(routersFileText, strings.Join(strImports, "\n"), strings.Join(strRouters, "\n"))
	err = ioutil.WriteFile(routersFileName, []byte(newRoutersText), 0644)
	if nil != err {
		logrus.Errorf("write new routers file failed. \n%s.", err)
		return
	}

	logrus.Info("Do go imports ...")
	// do goimports
	goimports.DoGoImports([]string{m.ApiDir}, true)
	logrus.Info("Do go imports completed")
	return
}

func (m *ApiParser) SaveApisToFile(apis []*ApiItem) (err error) {
	logrus.Info("Save apis ...")
	defer func() {
		if err == nil {
			logrus.Info("Save apis completed")
		}
	}()

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
