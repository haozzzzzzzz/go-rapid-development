package parser

import (
	"errors"
	"path/filepath"

	"go/ast"
	"go/token"
	"runtime/debug"

	"go/importer"
	"go/types"

	"go/parser"

	"os"

	"fmt"

	"strings"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
)

func (m *ApiParser) ScanApis(
	parseRequestData bool, // 如果parseRequestData会有点慢
) (apis []*ApiItem, err error) {
	apis = make([]*ApiItem, 0)
	logrus.Info("Scan api files...")
	defer func() {
		if err == nil {
			logrus.Info("Scan api files completed")
		}
	}()

	apiDir, err := filepath.Abs(m.ApiDir)
	if nil != err {
		logrus.Warnf("get absolute file apiDir failed. \n%s.", err)
		return
	}

	// api文件夹中所有的文件
	subApiDir := make([]string, 0)
	subApiDir, err = file.SearchFileNames(apiDir, func(fileInfo os.FileInfo) bool {
		if fileInfo.IsDir() {
			return true
		} else {
			return false
		}

	}, true)
	subApiDir = append(subApiDir, apiDir)

	// 服务源文件
	for _, subApiDir := range subApiDir {
		subApis, errParse := ParseApis(apiDir, subApiDir, parseRequestData)
		err = errParse
		if nil != err {
			logrus.Errorf("parse api file dir %q failed. error: %s.", subApiDir, err)
			return
		}

		apis = append(apis, subApis...)
	}

	return
}

func ParseApis(
	rootDir string,
	fileDir string,
	parseRequestData bool,
) (apis []*ApiItem, err error) {
	apis = make([]*ApiItem, 0)
	defer func() {
		if iRec := recover(); iRec != nil {
			logrus.Errorf("panic %s. api_dir: %s", iRec, fileDir)
			debug.PrintStack()
		}
	}()

	fileSet := token.NewFileSet()

	// 检索目录下所有的文件
	astFiles := make([]*ast.File, 0)
	pkgs, err := parser.ParseDir(fileSet, fileDir, nil, parser.AllErrors)
	if nil != err {
		logrus.Errorf("parser parse dir failed. error: %s.", err)
		return
	}
	_ = pkgs

	astFileMap := make(map[string]*ast.File)
	for pkgName, pkg := range pkgs {
		_ = pkgName
		for fileName, file := range pkg.Files {
			astFiles = append(astFiles, file)
			astFileMap[fileName] = file
		}
	}

	for _, astFile := range astFiles { // 遍历语法树
		apiPosFile := fileSet.File(astFile.Pos())
		fileName := apiPosFile.Name()
		for objName, obj := range astFile.Scope.Objects { // 遍历顶层所有变量，寻找HandleFunc
			valueSpec, ok := obj.Decl.(*ast.ValueSpec)
			if !ok {
				continue
			}

			_ = objName
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

			apiItem := &ApiItem{
				SourceFile:        fileName,
				ApiHandlerFunc:    objName,
				ApiHandlerPackage: astFile.Name.Name,
				RelativePaths:     make([]string, 0),
			}

			relativeDir := filepath.Dir(fileName)
			relativePackageDir := strings.Replace(relativeDir, rootDir, "", 1)
			if relativePackageDir != "" { // 子目录
				relativePackageDir = strings.Replace(relativePackageDir, "/", "", 1)
				apiItem.RelativePackage = strings.Replace(relativePackageDir, "/", "_", -1)
			}

			if apiItem.RelativePackage != "" {
				apiItem.PackagePath = relativeDir
			}

			for _, value := range valueSpec.Values { // 遍历属性
				compositeLit, ok := value.(*ast.CompositeLit)
				if !ok {
					continue
				}

				// compositeLit.Elts
				for _, elt := range compositeLit.Elts {
					keyValueExpr, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					keyIdent, ok := keyValueExpr.Key.(*ast.Ident)
					if !ok {
						continue
					}

					//fmt.Printf("%s %#v \n", keyIdent.Name, keyValueExpr.Value)
					switch keyIdent.Name {
					case "HttpMethod":
						valueLit, ok := keyValueExpr.Value.(*ast.BasicLit)
						if !ok {
							break
						}

						value := strings.Replace(valueLit.Value, "\"", "", -1)
						switch value {
						case "GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS", "DELETE", "CONNECT", "TRACE":
						default:
							err = errors.New(fmt.Sprintf("unsupported http method : %s", value))
							logrus.Errorf("mapping unsupported api failed. %s.", err)
							return
						}

						apiItem.HttpMethod = value

					case "RelativePath": // 废弃
						valueLit, ok := keyValueExpr.Value.(*ast.BasicLit)
						if !ok {
							break
						}

						value := strings.Replace(valueLit.Value, "\"", "", -1)
						apiItem.RelativePaths = append(apiItem.RelativePaths, value)

					case "RelativePaths":
						compLit, ok := keyValueExpr.Value.(*ast.CompositeLit)
						if !ok {
							break
						}

						for _, elt := range compLit.Elts {
							basicLit, ok := elt.(*ast.BasicLit)
							if !ok {
								continue
							}

							value := strings.Replace(basicLit.Value, "\"", "", -1)
							apiItem.RelativePaths = append(apiItem.RelativePaths, value)
						}

					case "Handle":
						funcLit, ok := keyValueExpr.Value.(*ast.FuncLit)
						if !ok {
							break
						}

						if parseRequestData == false {
							break
						}

						// types
						typesConf := types.Config{
							Importer: importer.For("source", nil),
							//Importer: importer.Default(),
						}

						info := &types.Info{
							Scopes:     make(map[ast.Node]*types.Scope),
							Defs:       make(map[*ast.Ident]types.Object),
							Uses:       make(map[*ast.Ident]types.Object),
							Types:      make(map[ast.Expr]types.TypeAndValue),
							Implicits:  make(map[ast.Node]types.Object),
							Selections: make(map[*ast.SelectorExpr]*types.Selection),
							InitOrder:  make([]*types.Initializer, 0),
						}
						pkg, errCheck := typesConf.Check(fileDir, fileSet, astFiles, info)
						err = errCheck
						if nil != err {
							logrus.Errorf("check types failed. error: %s.", err)
							return
						}
						_ = pkg

						// parse request data
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
										apiItem.PathData = parseApiRequest(ident, info)
									case "queryData":
										apiItem.QueryData = parseApiRequest(ident, info)
									case "postData":
										apiItem.PostData = parseApiRequest(ident, info)
									case "respData":
										apiItem.RespData = parseApiRequest(ident, info)
									}
								}

							case *ast.ReturnStmt:

							}

						}

					}

					//// property
					//valueLit, ok := keyValueExpr.Value.(*ast.BasicLit)
					//if ok {
					//	value := strings.Replace(valueLit.Value, "\"", "", -1)
					//	switch keyIdent.Name {
					//	case "HttpMethod":
					//		switch value {
					//		case "GET":
					//		case "POST":
					//		case "PUT":
					//		case "PATCH":
					//		case "HEAD":
					//		case "OPTIONS":
					//		case "DELETE":
					//		case "CONNECT":
					//		case "TRACE":
					//		default:
					//			err = errors.New(fmt.Sprintf("unsupported http method : %s", value))
					//			logrus.Errorf("mapping unsupported api failed. %s.", err)
					//			return
					//		}
					//		apiItem.HttpMethod = value
					//
					//	case "RelativePath":
					//		apiItem.RelativePaths = append(apiItem.RelativePaths, value)
					//
					//	case "RelativePaths":
					//		fmt.Printf("aha")
					//
					//	}
					//
					//}
					//
					//// handle func
					//funcLit, ok := keyValueExpr.Value.(*ast.FuncLit)
					//if ok {
					//	switch keyIdent.Name {
					//	case "Handle":
					//		if parseRequestData == false {
					//			break
					//		}
					//
					//		// types
					//		typesConf := types.Config{
					//			Importer: importer.For("source", nil),
					//			//Importer: importer.Default(),
					//		}
					//
					//		info := &types.Info{
					//			Scopes:     make(map[ast.Node]*types.Scope),
					//			Defs:       make(map[*ast.Ident]types.Object),
					//			Uses:       make(map[*ast.Ident]types.Object),
					//			Types:      make(map[ast.Expr]types.TypeAndValue),
					//			Implicits:  make(map[ast.Node]types.Object),
					//			Selections: make(map[*ast.SelectorExpr]*types.Selection),
					//			InitOrder:  make([]*types.Initializer, 0),
					//		}
					//		pkg, errCheck := typesConf.Check(fileDir, fileSet, astFiles, info)
					//		err = errCheck
					//		if nil != err {
					//			logrus.Errorf("check types failed. error: %s.", err)
					//			return
					//		}
					//		_ = pkg
					//
					//		// parse request data
					//		funcBody := funcLit.Body
					//		for _, funcStmt := range funcBody.List {
					//			switch funcStmt.(type) {
					//			case *ast.AssignStmt:
					//				assignStmt := funcStmt.(*ast.AssignStmt)
					//				lhs := assignStmt.Lhs
					//				rhs := assignStmt.Rhs
					//
					//				_ = lhs
					//				_ = rhs
					//
					//				for _, expr := range lhs {
					//					ident, ok := expr.(*ast.Ident)
					//					if !ok {
					//						continue
					//					}
					//
					//					switch ident.Name {
					//					case "pathData":
					//						apiItem.PathData = parseApiRequest(ident, info)
					//					case "queryData":
					//						apiItem.QueryData = parseApiRequest(ident, info)
					//					case "postData":
					//						apiItem.PostData = parseApiRequest(ident, info)
					//					case "respData":
					//						apiItem.RespData = parseApiRequest(ident, info)
					//					}
					//				}
					//
					//			case *ast.ReturnStmt:
					//
					//			}
					//
					//		}
					//	}
					//}

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

func parseApiRequest(
	astIdent *ast.Ident,
	info *types.Info,
) (dataType *StructType) {
	identType := info.Defs[astIdent]
	typeVar, ok := identType.(*types.Var)
	if !ok {
		return
	}

	iType := parseType(typeVar.Type())
	if iType != nil {
		dataType, _ = iType.(*StructType)
	}

	return
}

func parseType(t types.Type) (iType IType) {
	iType = NewBasicType("Unsupported")

	switch t.(type) {
	case *types.Basic:
		iType = NewBasicType(t.(*types.Basic).Name())

	case *types.Pointer:
		iType = parseType(t.(*types.Pointer).Elem())

	case *types.Named:
		tNamed := t.(*types.Named)
		iType = parseType(tNamed.Underlying())

		// 如果是structType
		structType, ok := iType.(*StructType)
		if ok {
			structType.Name = tNamed.Obj().Name()
			iType = structType
		}

	case *types.Struct: // 匿名
		structType := NewStructType()

		tStructType := t.(*types.Struct)

		numFields := tStructType.NumFields()
		for i := 0; i < numFields; i++ {
			field := NewField()

			tField := tStructType.Field(i)
			if !tField.Exported() {
				continue
			}

			// tags
			tagValue := strings.Replace(tStructType.Tag(i), "`", "", -1)
			strPairs := strings.Split(tagValue, " ")
			for _, pair := range strPairs {
				if pair == "" {
					continue
				}

				tagPair := strings.Split(pair, ":")
				field.Tags[tagPair[0]] = tagPair[1]
			}

			// definition
			field.Name = tField.Name()
			fieldType := parseType(tField.Type())
			field.TypeName = fieldType.TypeName()
			field.TypeSpec = fieldType

			structType.Fields = append(structType.Fields, field)

		}

		iType = structType

	case *types.Slice:
		arrType := NewArrayType()
		eltType := parseType(t.(*types.Slice).Elem())
		arrType.EltSpec = eltType
		arrType.EltName = eltType.TypeName()
		arrType.Name = fmt.Sprintf("[]%s", eltType.TypeName())

		iType = arrType

	case *types.Map:
		mapType := NewMapType()
		tMap := t.(*types.Map)
		mapType.ValueSpec = parseType(tMap.Elem())
		mapType.KeySpec = parseType(tMap.Key())
		mapType.Name = fmt.Sprintf("map[%s]%s", mapType.KeySpec.TypeName(), mapType.ValueSpec.TypeName())

		iType = mapType

	case *types.Interface:
		iType = NewInterfaceType()

	default:
		fmt.Printf("parse unsupported type %#v\n", t)

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
