package api

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ApiParser struct {
	ProjectPath string
	ApiDir      string
}

func NewApiParser(projectPath string) *ApiParser {
	return &ApiParser{
		ProjectPath: projectPath,
		ApiDir:      fmt.Sprintf("%s/api", projectPath),
	}
}

func (m *ApiParser) ScanApi() (apis []*ApiItem, err error) {
	path, err := filepath.Abs(m.ApiDir)
	if nil != err {
		logrus.Warnf("get absolute file path failed. \n%s.", err)
		return
	}

	files, err := file.SearchFileNames(path, func(fileInfo os.FileInfo) bool {
		fileName := fileInfo.Name()
		if strings.HasPrefix(fileName, "api_") {
			return true
		}
		return false
	}, true)

	logrus.Info("Scanning api files...")
	defer func() {
		if err == nil {
			logrus.Info("Scan api files completed")
		}
	}()

	sort.Strings(files)
	for _, fileName := range files {
		fileSet := token.NewFileSet()
		astFile, errParse := parser.ParseFile(fileSet, fileName, nil, parser.AllErrors)
		if nil != errParse {
			err = errParse
			logrus.Warnf("parser parse file %q failed. \n%s.", fileName, err)
			return
		}

		packageName := astFile.Name.Name

		for objName, obj := range astFile.Scope.Objects {
			valueSpec, ok := obj.Decl.(*ast.ValueSpec)
			if !ok {
				continue
			}

			selectorExpr, ok := valueSpec.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			xIdent, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				continue
			}

			selIdent := selectorExpr.Sel

			if xIdent.Name != "ginbuilder" && selIdent.Name != "HandleFunc" {
				continue
			}

			apiItem := new(ApiItem)
			apiItem.SourceFile = strings.Replace(fileName, m.ProjectPath, "", 1)
			apiItem.ApiHandlerFunc = objName
			apiItem.ApiHandlerPackage = packageName

			for _, value := range valueSpec.Values {
				compositeLit, ok := value.(*ast.CompositeLit)
				if !ok {
					continue
				}

				for _, elt := range compositeLit.Elts {
					keyValueExpr, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					keyIdent, ok := keyValueExpr.Key.(*ast.Ident)
					if !ok {
						continue
					}

					valueLit, ok := keyValueExpr.Value.(*ast.BasicLit)
					if !ok {
						continue
					}

					value := strings.Replace(valueLit.Value, "\"", "", -1)
					switch keyIdent.Name {
					case "HttpMethod":
						apiItem.HttpMethod = value
					case "RelativePath":
						apiItem.RelativePath = value
					}
				}
			}

			err = validator.New().Struct(apiItem)
			if nil != err {
				logrus.Errorf("%#v\n invalid", apiItem)
				return
			}

			apis = append(apis, apiItem)
		}

	}

	return
}

func (m *ApiParser) MapApi() (err error) {
	apis, err := m.ScanApi()
	if nil != err {
		logrus.Errorf("Scan api failed. \n%s.", err)
		return
	}

	logrus.Info("Mapping apis ...")
	defer func() {
		if err == nil {
			logrus.Info("Map apis completed")
		}
	}()

	newRoutersFileMiddle := bytes.NewBuffer(nil)
	for _, apiItem := range apis {
		str := fmt.Sprintf("    engine.Handle(\"%s\", \"%s\", %s.%s.HandlerFunc)\n", apiItem.HttpMethod, apiItem.RelativePath, apiItem.ApiHandlerPackage, apiItem.ApiHandlerFunc)
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

	err = ioutil.WriteFile(fmt.Sprintf("%s/.proj/apis.yaml", m.ProjectPath), byteYamlApis, 0644)
	if nil != err {
		logrus.Errorf("write apis.yaml failed. \n%s.", err)
		return
	}

	return
}
