package parser

import (
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/haozzzzzzzz/go-rapid-development/tools/lib/gofmt"

	"strings"

	"os"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/sirupsen/logrus"
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
	logrus.Info("Map apis ...")
	defer func() {
		if err == nil {
			logrus.Info("Map apis completed")
		}
	}()

	importsMap := make(map[string]string) // package_exported -> alias
	strRouters := make([]string, 0)

	for _, apiItem := range apis {
		// imports
		if apiItem.PackageRelAlias != "" { // 本目录下的不用加入
			if apiItem.PackageExportedPath == "" {
				err = uerrors.Newf("alias require exported path. alias: %s", apiItem.PackageRelAlias)
				return
			}

			importsMap[apiItem.PackageExportedPath] = apiItem.PackageRelAlias
		}

		// handle func binding
		strHandleFunc := apiItem.ApiHandlerFunc
		if apiItem.PackageRelAlias != "" {
			strHandleFunc = fmt.Sprintf("%s.%s", apiItem.PackageRelAlias, apiItem.ApiHandlerFunc)
		}

		for _, uri := range apiItem.RelativePaths {
			str := fmt.Sprintf("    engine.Handle(\"%s\", \"%s\", %s.GinHandler)", apiItem.HttpMethod, uri, strHandleFunc)
			strRouters = append(strRouters, str)
		}

	}

	strImports := make([]string, 0)
	for expPath, alias := range importsMap {
		var str string
		if alias == "" {
			str = fmt.Sprintf("    %q", expPath)
		} else {
			str = fmt.Sprintf("    %s %q", alias, expPath)
		}
		strImports = append(strImports, str)
	}

	routersFileName := fmt.Sprintf("%s/routers.go", m.ApiDir)
	newRoutersText := fmt.Sprintf(routersFileText, strings.Join(strImports, "\n"), strings.Join(strRouters, "\n"))
	newRoutersText, err = gofmt.StrGoFmt(newRoutersText)
	if nil != err {
		logrus.Errorf("go fmt source failed. text: %s, error: %s.", newRoutersText, err)
		return
	}

	err = ioutil.WriteFile(routersFileName, []byte(newRoutersText), 0644)
	if nil != err {
		logrus.Errorf("write new routers file failed. \n%s.", err)
		return
	}

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
