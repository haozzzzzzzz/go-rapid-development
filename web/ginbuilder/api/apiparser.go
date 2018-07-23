package api

import (
	"bytes"
	"errors"
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
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/haozzzzzzzz/go-rapid-development/utils/printutil"
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

		// package name
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

				// compositeLit.Elts
				//&ast.KeyValueExpr{Key:(*ast.Ident)(0xc42026e5c0), Colon:167, Value:(*ast.BasicLit)(0xc42026e5e0)}
				//&ast.KeyValueExpr{Key:(*ast.Ident)(0xc42026e600), Colon:192, Value:(*ast.BasicLit)(0xc42026e620)}
				//&ast.KeyValueExpr{Key:(*ast.Ident)(0xc42026e660), Colon:231, Value:(*ast.FuncLit)(0xc4202196e0)}
				for _, elt := range compositeLit.Elts {
					keyValueExpr, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					//&ast.Ident{NamePos:157, Name:"HttpMethod", Obj:(*ast.Object)(nil)}
					//&ast.Ident{NamePos:180, Name:"RelativePath", Obj:(*ast.Object)(nil)}
					//&ast.Ident{NamePos:225, Name:"Handle", Obj:(*ast.Object)(nil)}

					keyIdent, ok := keyValueExpr.Key.(*ast.Ident)
					if !ok {
						continue
					}

					// property
					valueLit, ok := keyValueExpr.Value.(*ast.BasicLit)
					if ok {
						value := strings.Replace(valueLit.Value, "\"", "", -1)
						switch keyIdent.Name {
						case "HttpMethod":
							switch value {
							case "GET":
							case "POST":
							case "PUT":
							case "PATCH":
							case "HEAD":
							case "OPTIONS":
							case "DELETE":
							case "CONNECT":
							case "TRACE":
							default:
								err = errors.New(fmt.Sprintf("unsupported http method : %s", value))
								logrus.Errorf("mapping unsupported api failed. %s.", err)
								return
							}
							apiItem.HttpMethod = value

						case "RelativePath":
							apiItem.RelativePath = value

						}
					}

					// handle func
					funcLit, ok := keyValueExpr.Value.(*ast.FuncLit)
					if ok {
						switch keyIdent.Name {
						case "Handle":
							funcBody := funcLit.Body
							for _, funcStmt := range funcBody.List {
								switch funcStmt.(type) {
								case *ast.AssignStmt:
									assignStmt := funcStmt.(*ast.AssignStmt)
									lhs := assignStmt.Lhs
									rhs := assignStmt.Rhs

									_ = lhs
									_ = rhs

									for _, expr := range lhs {
										ident, ok := expr.(*ast.Ident)
										if !ok {
											continue
										}

										switch ident.Name {
										case "pathData":
											apiItem.PathData = m.parseApiStructData(ident)
										case "queryData":
											apiItem.QueryData = m.parseApiStructData(ident)
										case "postData":
											apiItem.PostData = m.parseApiStructData(ident)
										case "respData":
											apiItem.RespData = m.parseApiStructData(ident)
										}
									}

								case *ast.IfStmt:
								case *ast.ExprStmt:
								case *ast.ReturnStmt:

								}

							}
						}
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

func (m *ApiParser) parseApiStructData(ident *ast.Ident) (structData *StructData) {
	structData = &StructData{
		Name: ident.Name,
	}

	decl := ident.Obj.Decl
	declStmt, ok := decl.(*ast.AssignStmt)
	if !ok {
		return
	}

	for _, expr := range declStmt.Rhs {
		// 指针
		var compositeLit *ast.CompositeLit
		pointerExpr, ok := expr.(*ast.UnaryExpr)
		if ok {
			expr = pointerExpr.X
		}

		// 非指针
		compositeLit, ok = expr.(*ast.CompositeLit)
		if !ok {
			continue
		}

		structType, ok := compositeLit.Type.(*ast.StructType)
		if !ok {
			continue
		}

		// &ast.Field{Doc:(*ast.CommentGroup)(nil), Names:[]*ast.Ident{(*ast.Ident)(0xc42027a7e0)}, Type:(*ast.Ident)(0xc42027a800), Tag:(*ast.BasicLit)(0xc42027a820), Comment:(*ast.CommentGroup)(nil)}
		for _, field := range structType.Fields.List {
			fieldDoc := field.Doc
			fieldNames := field.Names
			fieldType := field.Type
			fieldTag := field.Tag
			fieldComment := field.Comment

			_ = fieldDoc
			_ = fieldNames
			_ = fieldType
			_ = fieldTag
			_ = fieldComment

			name := fieldNames[0].Name
			structDataField := &StructDataField{
				Name: name,
			}
			structData.Fields = append(structData.Fields, structDataField)

			switch fieldType.(type) {
			case *ast.Ident:
				ident := fieldType.(*ast.Ident)
				structDataField.Type = ident.Name
				if ident.Obj != nil {
					printutil.PrintObject(ident.Obj.Decl)
				}

			case *ast.StarExpr:
			case *ast.ArrayType:
			case *ast.MapType:
			case *ast.StructType:

			}

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
		str := fmt.Sprintf("    engine.Handle(\"%s\", \"%s\", %s.%s.GinHandler)\n", apiItem.HttpMethod, apiItem.RelativePath, apiItem.ApiHandlerPackage, apiItem.ApiHandlerFunc)
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

	//logrus.Info("Doing go imports")
	//// do goimports
	//goimports.DoGoImports([]string{m.ApiDir}, true)
	//logrus.Info("Do go imports completed")

	// save api.yaml
	byteYamlApis, err := yaml.Marshal(apis)
	if nil != err {
		logrus.Errorf("yaml marshal apis failed. \n%s.", byteYamlApis)
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/.service/apis.yaml", m.ProjectPath), byteYamlApis, 0644)
	if nil != err {
		logrus.Errorf("write apis.yaml failed. \n%s.", err)
		return
	}

	return
}
