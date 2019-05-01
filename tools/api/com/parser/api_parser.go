package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
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
	err = m.MapApi()
	if nil != err {
		logrus.Errorf("mapping api failed. %s.", err)
		return
	}

	return
}

func (m *ApiParser) ScanApis() (apis []*ApiItem, err error) {
	fileDir, err := filepath.Abs(m.ApiDir)
	if nil != err {
		logrus.Warnf("get absolute file fileDir failed. \n%s.", err)
		return
	}

	files, err := file.SearchFileNames(fileDir, func(fileInfo os.FileInfo) bool {
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
		fileApis, errParse := ParseApiFile(fileDir, fileName)
		err = errParse
		if nil != err {
			logrus.Errorf("parse api file %q failed. error: %s.", fileName, err)
			return
		}

		apis = append(apis, fileApis...)
	}

	return
}

func ParseApiFile(fileDir string, fileName string) (apis []*ApiItem, err error) {
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
		apiItem.SourceFile = strings.Replace(fileName, fileDir, "", 1)
		apiItem.ApiHandlerFunc = objName
		apiItem.ApiHandlerPackage = packageName
		relativeDir := filepath.Dir(fileName)
		relativePackageDir := strings.Replace(relativeDir, fileDir, "", 1)
		if relativePackageDir != "" { // 子目录
			relativePackageDir = strings.Replace(relativePackageDir, "/", "", 1)
			apiItem.RelativePackage = strings.Replace(relativePackageDir, "/", "_", -1)
		}

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
										apiItem.PathData = parseApiStructData(ident)
									case "queryData":
										apiItem.QueryData = parseApiStructData(ident)
									case "postData":
										apiItem.PostData = parseApiStructData(ident)
									case "respData":
										apiItem.RespData = parseApiStructData(ident)
									}
								}

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

	return
}

func ParseStructData(typeSpec *ast.TypeSpec) (structData *StructType) {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return
	}

	structData = NewStructData()
	structData.Name = typeSpec.Name.Name
	for _, field := range structType.Fields.List {
		var exprFieldType ast.Expr
		switch field.Type.(type) {
		case *ast.StarExpr:
			exprFieldType = field.Type.(*ast.StarExpr).X
		default:
			exprFieldType = field.Type
		}

		if field.Names == nil { // 复用了其他struct
			fieldParent := exprFieldType.(*ast.Ident).Obj
			parentStruct := ParseStructData(fieldParent.Decl.(*ast.TypeSpec))
			for _, pField := range parentStruct.Fields {
				structData.Fields = append(structData.Fields, pField)
			}

		} else {

			//&ast.Field{Doc:(*ast.CommentGroup)(nil), Names:[]*ast.Ident{(*ast.Ident)(0xc4202b21e0)}, Type:(*ast.Ident)(0xc4202b2200), Tag:(*ast.BasicLit)(0xc4202b2220), Comment:(*ast.CommentGroup)(nil)}
			structField := NewField()
			structField.Name = field.Names[0].Name
			switch exprFieldType.(type) {
			case *ast.Ident:
				fieldType := exprFieldType.(*ast.Ident)
				structField.Type = fieldType.Name

				tagValue := strings.Replace(field.Tag.Value, "`", "", -1)
				strPairs := strings.Split(tagValue, " ")
				for _, pair := range strPairs {
					pair = strings.Replace(pair, "\"", "", -1)
					tagPair := strings.Split(pair, ":")
					structField.Tags[tagPair[0]] = tagPair[1]
				}

				if fieldType.Obj != nil { // struct
					structField.Spec = ParseStructData(fieldType.Obj.Decl.(*ast.TypeSpec))
				}

			case *ast.MapType:
				fieldType := exprFieldType.(*ast.MapType)
				mapKey := fieldType.Key.(*ast.Ident).Name

				fmt.Printf("%#v\n", fieldType.Value.(*ast.InterfaceType))
				_ = mapKey
			}

			structData.Fields = append(structData.Fields, structField)
		}

	}

	return
}

// TODO
func parseApiStructData(ident *ast.Ident) (structData *StructType) {
	structData = &StructType{
		Name: ident.Name,
	}

	decl := ident.Obj.Decl
	declStmt, ok := decl.(*ast.AssignStmt)
	if !ok {
		return
	}

	// 指针
	var compositeLit *ast.CompositeLit
	expr := declStmt.Rhs[0]
	pointerExpr, ok := expr.(*ast.UnaryExpr)
	if ok {
		expr = pointerExpr.X
	}

	// 非指针
	compositeLit, ok = expr.(*ast.CompositeLit)
	if !ok {
		return
	}

	structData = ParseStructData(compositeLit.Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec))
	return
}

func (m *ApiParser) MapApi() (err error) {
	apis, err := m.ScanApis()
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
