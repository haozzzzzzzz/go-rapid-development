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

	"runtime/debug"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/project"
	"github.com/haozzzzzzzz/go-rapid-development/tools/goimports"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
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
	err = m.MapApi()
	if nil != err {
		logrus.Errorf("mapping api failed. %s.", err)
		return
	}

	return
}

func (m *ApiParser) ScanApis() (apis []*ApiItem, err error) {
	logrus.Info("Scanning api files...")
	defer func() {
		if err == nil {
			logrus.Info("Scan api files completed")
		}
	}()

	apiFileDir, err := filepath.Abs(m.ApiDir)
	if nil != err {
		logrus.Warnf("get absolute file apiFileDir failed. \n%s.", err)
		return
	}

	// api文件夹中所有的文件
	fileNames, err := file.SearchFileNames(m.ApiDir, func(fileInfo os.FileInfo) bool {
		fileName := fileInfo.Name()
		if strings.HasPrefix(fileName, "api_") {
			return true
		}
		return false

	}, true)

	// 服务源文件
	sort.Strings(fileNames)
	for _, fileName := range fileNames {
		fileApis, errParse := ParseApiFile(apiFileDir, fileName)
		err = errParse
		if nil != err {
			logrus.Errorf("parse api file %q failed. error: %s.", fileName, err)
			return
		}

		apis = append(apis, fileApis...)
	}

	return
}

func ParseApiFile(
	apiFileDir string,
	fileName string,
) (apis []*ApiItem, err error) {
	defer func() {
		if iRec := recover(); iRec != nil {
			logrus.Errorf("panic %s. file_name: %s", iRec, fileName)
			debug.PrintStack()
		}
	}()

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
		apiItem.SourceFile = strings.Replace(fileName, apiFileDir, "", 1)
		apiItem.ApiHandlerFunc = objName
		apiItem.ApiHandlerPackage = packageName
		relativeDir := filepath.Dir(fileName)
		relativePackageDir := strings.Replace(relativeDir, apiFileDir, "", 1)
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

// 需要解决调用层级问题
func parseStructType(typeSpec *ast.TypeSpec, maxLevel int, curLevel int) (structType *StructType) {
	if curLevel > maxLevel {
		return
	}

	specStructType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return
	}

	structType = NewStructType()
	structType.Name = typeSpec.Name.Name
	//fmt.Printf("parse_struct_types:%s  max_level: %d cur_level: %d\n", structType.Name, maxLevel, curLevel)
	for _, field := range specStructType.Fields.List {
		exprFieldType := convertExpr(field.Type)

		if field.Names == nil { // 复用了其他struct
			fieldParent := exprFieldType.(*ast.Ident).Obj
			parentStruct := parseStructType(fieldParent.Decl.(*ast.TypeSpec), maxLevel, curLevel+1)
			if parentStruct != nil {
				for _, pField := range parentStruct.Fields {
					structType.Fields = append(structType.Fields, pField)
				}
			}

		} else {

			//&ast.Field{Doc:(*ast.CommentGroup)(nil), Names:[]*ast.Ident{(*ast.Ident)(0xc4202b21e0)}, Type:(*ast.Ident)(0xc4202b2200), Tag:(*ast.BasicLit)(0xc4202b2220), Comment:(*ast.CommentGroup)(nil)}
			structField := NewField()
			structField.Name = field.Names[0].Name

			tagValue := strings.Replace(field.Tag.Value, "`", "", -1)
			strPairs := strings.Split(tagValue, " ")
			for _, pair := range strPairs {
				pair = strings.Replace(pair, "\"", "", -1)
				tagPair := strings.Split(pair, ":")
				structField.Tags[tagPair[0]] = tagPair[1]
			}

			iType := parseTypeSpec(exprFieldType, maxLevel, curLevel+1)
			structField.Type = iType.TypeName()
			structField.TypeSpec = iType

			structType.Fields = append(structType.Fields, structField)
		}

	}

	return
}

func convertExpr(expr ast.Expr) (newExpr ast.Expr) {
	switch expr.(type) {
	case *ast.StarExpr:
		newExpr = expr.(*ast.StarExpr).X

	case *ast.SelectorExpr:
		newExpr = expr.(*ast.SelectorExpr).Sel

	default:
		newExpr = expr

	}

	return
}

// iType 一定不为nil，即使不向下递归也要返回Name
func parseTypeSpec(typeExpr ast.Expr, maxLevel int, curLevel int) (iType IType) {
	typeExpr = convertExpr(typeExpr) // 转换引用类型

	//fmt.Printf("parse_type_spec:%s max_level:%d cur_level:%d\n", typeExpr, maxLevel, curLevel)

	switch typeExpr.(type) {
	case *ast.Ident:
		ident := typeExpr.(*ast.Ident)
		if ident.Obj == nil { // 标准类型
			standType := StandardType(ident.Name)
			iType = standType

		} else { // 自定义类型
			structType := NewStructType()
			structType.Name = ident.Name
			if curLevel <= maxLevel {
				parsed := parseStructType(ident.Obj.Decl.(*ast.TypeSpec), maxLevel, curLevel)
				if parsed != nil {
					structType = parsed
				}
			}

			iType = structType
		}

	case *ast.MapType:
		exprMapType := typeExpr.(*ast.MapType)
		keyName := exprMapType.Key.(*ast.Ident).Name
		mapType := &MapType{
			Key: keyName,
		}

		exprValue := convertExpr(exprMapType.Value)
		if curLevel <= maxLevel {
			valueType := parseTypeSpec(exprMapType.Value, maxLevel, curLevel)
			mapType.ValueSpec = valueType
			mapType.Name = fmt.Sprintf("map[%s]%s", mapType.Key, valueType.TypeName())

		} else {
			value := "interface{}"
			switch exprValue.(type) {
			case *ast.Ident:
				value = exprValue.(*ast.Ident).Name
			default:
				// 如果value的type为map，则设为interface
			}

			mapType.Name = fmt.Sprintf("map[%s]%s", mapType.Key, value)

		}

		iType = mapType

	default:
		fmt.Printf("parse_type_sepc %#v\n", typeExpr)

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

	structData = parseStructType(compositeLit.Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec), 3, 1) // api层次最多3层，避免无穷递归
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
