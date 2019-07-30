package parser

import (
	"errors"
	"go/ast"
	"go/importer"
	"go/token"
	"path/filepath"
	"runtime/debug"

	"golang.org/x/tools/go/packages"

	"go/types"

	"go/parser"

	"os"

	"fmt"

	"strings"

	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/api/request"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
)

func (m *ApiParser) ScanApis(
	parseRequestData bool, // 如果parseRequestData会有点慢
	useMod bool, // parseRequestData=true时，生效
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

	// 服务源文件，只能一个pkg一个pkg地解析
	for _, subApiDir := range subApiDir {
		subApis, errParse := ParsePkgApis(m.GoPaths, apiDir, subApiDir, parseRequestData, useMod)
		err = errParse
		if nil != err {
			logrus.Errorf("parse api file dir %q failed. error: %s.", subApiDir, err)
			return
		}

		apis = append(apis, subApis...)
	}

	return
}

func mergeTypesInfos(info *types.Info, infos ...*types.Info) {
	for _, tempInfo := range infos {
		for scopeKey, scopeVal := range tempInfo.Scopes {
			info.Scopes[scopeKey] = scopeVal
		}

		for defKey, defVal := range tempInfo.Defs {
			info.Defs[defKey] = defVal
		}

		for useKey, useVal := range tempInfo.Uses {
			info.Uses[useKey] = useVal
		}

		for implKey, implVal := range tempInfo.Implicits {
			info.Implicits[implKey] = implVal
		}

		for selKey, selVal := range tempInfo.Selections {
			info.Selections[selKey] = selVal
		}

		// do not need to merge InitOrder
	}
	return
}

func ParsePkgApis(
	goPaths []string,
	apiRootDir string,
	apiPackageDir string,
	parseRequestData bool,
	useMod bool,
) (apis []*ApiItem, err error) {
	apis = make([]*ApiItem, 0)
	defer func() {
		if iRec := recover(); iRec != nil {
			logrus.Errorf("panic %s. api_dir: %s", iRec, apiPackageDir)
			debug.PrintStack()
		}
	}()

	fileSet := token.NewFileSet()

	// 检索目录下所有的文件
	astFiles := make([]*ast.File, 0)

	// types
	typesInfo := &types.Info{
		Scopes:     make(map[ast.Node]*types.Scope),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		InitOrder:  make([]*types.Initializer, 0),
	}
	if parseRequestData {
		if !useMod { // not use mod
			parseMode := parser.AllErrors | parser.ParseComments
			astPkgMap, errParse := parser.ParseDir(fileSet, apiPackageDir, nil, parseMode)
			err = errParse
			if nil != err {
				logrus.Errorf("parser parse dir failed. error: %s.", err)
				return
			}
			_ = astPkgMap

			astFileMap := make(map[string]*ast.File)
			for pkgName, astPkg := range astPkgMap {
				_ = pkgName
				for fileName, pkgFile := range astPkg.Files {
					astFiles = append(astFiles, pkgFile)
					astFileMap[fileName] = pkgFile
				}
			}

			typesConf := types.Config{
				Importer: importer.Default(),
			}

			//typesConf.Importer = importer.For("source", nil)

			// cur package
			pkg, errCheck := typesConf.Check(apiPackageDir, fileSet, astFiles, typesInfo)
			err = errCheck
			if nil != err {
				logrus.Errorf("check types failed. error: %s.", err)
				return
			}
			_ = pkg

			// imported
			impPaths := make(map[string]bool)
			for fileName, pkgFile := range astFileMap {
				_ = fileName
				for _, fileImport := range pkgFile.Imports {
					impPath := strings.Replace(fileImport.Path.Value, "\"", "", -1)
					for _, goPath := range goPaths {
						absPath := fmt.Sprintf("%s/src/%s", goPath, impPath)
						if file.PathExists(absPath) {
							impPaths[absPath] = true
						}
					}
				}
			}

			for impPath, _ := range impPaths {
				tempFileSet := token.NewFileSet()
				tempAstFiles := make([]*ast.File, 0)

				tempPkgs, errParse := parser.ParseDir(tempFileSet, impPath, nil, parseMode)
				err = errParse
				if nil != err {
					logrus.Errorf("parser parse dir failed. error: %s.", err)
					return
				}

				for pkgName, pkg := range tempPkgs {
					_ = pkgName
					for _, pkgFile := range pkg.Files {
						//pkgFile.Imports ? 不递归下去了，目前最多两层
						tempAstFiles = append(tempAstFiles, pkgFile)
					}
				}

				// type check imported path
				_, err = typesConf.Check(impPath, tempFileSet, tempAstFiles, typesInfo)
				if nil != err {
					logrus.Errorf("check imported package failed. path: %s, error: %s.", impPath, err)
					return
				}

			}

		} else {
			mode := packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles |
				packages.NeedImports |
				packages.NeedDeps |
				packages.NeedExportsFile |
				packages.NeedTypes |
				packages.NeedSyntax |
				packages.NeedTypesInfo |
				packages.NeedTypesSizes
			pkgs, errLoad := packages.Load(&packages.Config{
				Mode: mode,
				Dir:  apiPackageDir,
				Fset: fileSet,
			})
			err = errLoad
			if nil != err {
				logrus.Errorf("go/packages load failed. error: %s.", err)
				return
			}

			logrus.Infof("scan imports")

			// only one package
			impPaths := make(map[string]bool)
			for _, pkg := range pkgs {
				_, ok := impPaths[pkg.PkgPath]
				if ok {
					continue
				}

				impPaths[pkg.PkgPath] = true

				// 需要合并token.FileSet、ast.File和types.Info
				//logrus.Info(pkg.PkgPath)

				pkg.Fset.Iterate(func(i *token.File) bool {
					fmt.Println(i.Name())
					fmt.Println(i.Base())
					fmt.Println(i.Size())
					f := fileSet.AddFile(i.Name(), fileSet.Base(), i.Size())
					_ = f
					return true
				})

				for _, astFile := range pkg.Syntax {
					astFiles = append(astFiles, astFile)
				}

				mergeTypesInfos(typesInfo, pkg.TypesInfo)
				for _, impPkg := range pkg.Imports { // 依赖
					_, ok := impPaths[impPkg.PkgPath]
					if ok {
						continue
					}

					impPaths[impPkg.PkgPath] = true

					impPkg.Fset.Iterate(func(i *token.File) bool {
						fmt.Println(i.Name())
						f := fileSet.AddFile(i.Name(), fileSet.Base(), i.Size())
						_ = f
						return true
					})

					for _, astFile := range impPkg.Syntax {
						astFiles = append(astFiles, astFile)
					}

					//logrus.Info(impPkg.PkgPath)
					mergeTypesInfos(typesInfo, impPkg.TypesInfo)
				}
			}
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

			fileDir := filepath.Dir(fileName)

			// 设置相对api文件夹的relative path
			pkgDir := strings.Replace(fileDir, apiRootDir, "", 1)
			if pkgDir != "" { // 子目录
				pkgDir = strings.Replace(pkgDir, "/", "", 1)
				apiItem.RelativePackage = strings.Replace(pkgDir, "/", "_", -1)
			}

			if apiItem.RelativePackage != "" {
				apiItem.PackagePath = fileDir
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

					switch keyIdent.Name {
					case "HttpMethod":
						valueLit, ok := keyValueExpr.Value.(*ast.BasicLit)
						if !ok {
							break
						}

						value := strings.Replace(valueLit.Value, "\"", "", -1)
						switch value {
						case request.METHOD_GET,
							request.METHOD_POST,
							request.METHOD_PUT,
							request.METHOD_PATCH,
							request.METHOD_HEAD,
							request.METHOD_OPTIONS,
							request.METHOD_DELETE,
							request.METHOD_CONNECT,
							request.METHOD_TRACE:
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

						// if parse request data
						if parseRequestData == false {
							break
						}

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
										apiItem.PathData = parseApiRequest(typesInfo, ident)
									case "queryData":
										apiItem.QueryData = parseApiRequest(typesInfo, ident)
									case "postData":
										apiItem.PostData = parseApiRequest(typesInfo, ident)
									case "respData":
										apiItem.RespData = parseApiRequest(typesInfo, ident)
									}
								}

							case *ast.ReturnStmt:

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

func parseApiRequest(
	info *types.Info,
	astIdent *ast.Ident,
) (dataType *StructType) {
	identType := info.Defs[astIdent]
	typeVar, ok := identType.(*types.Var)
	if !ok {
		return
	}

	iType := parseType(info, typeVar.Type())
	if iType != nil {
		dataType, _ = iType.(*StructType)
	} else {
		logrus.Warnf("parse api request nill: %#v\n", typeVar)
	}

	return
}

func parseType(
	info *types.Info,
	t types.Type,
) (iType IType) {
	iType = NewBasicType("Unsupported")

	switch t.(type) {
	case *types.Basic:
		iType = NewBasicType(t.(*types.Basic).Name())

	case *types.Pointer:
		iType = parseType(info, t.(*types.Pointer).Elem())

	case *types.Named:
		tNamed := t.(*types.Named)
		iType = parseType(info, tNamed.Underlying())

		// 如果是structType
		structType, ok := iType.(*StructType)
		if ok {
			structType.Name = tNamed.Obj().Name()
			iType = structType
		}

	case *types.Struct: // 匿名
		structType := NewStructType()

		tStructType := t.(*types.Struct)

		typeAstExpr := FindStructAstExprFromInfoTypes(info, tStructType)
		if typeAstExpr == nil {
			logrus.Warnf("cannot found expr of type: %s", tStructType.String())
		}

		numFields := tStructType.NumFields()
		for i := 0; i < numFields; i++ {
			field := NewField()

			tField := tStructType.Field(i)
			if !tField.Exported() {
				continue
			}

			if typeAstExpr != nil { // 找到声明

				astStructType, ok := typeAstExpr.(*ast.StructType)
				if !ok {
					logrus.Printf("parse struct type failed. expr: %s, type: %#v\n\n\n\n\n", typeAstExpr, tStructType)
					return
				}

				astField := astStructType.Fields.List[i]

				// 注释
				if astField.Doc != nil && len(astField.Doc.List) > 0 {
					for _, comment := range astField.Doc.List {
						if field.Description != "" {
							field.Description += "; "
						}

						field.Description += RemoveCommentStartEndToken(comment.Text)
					}
				}

				if astField.Comment != nil && len(astField.Comment.List) > 0 {
					for _, comment := range astField.Comment.List {
						if field.Description != "" {
							field.Description += "; "
						}
						field.Description += RemoveCommentStartEndToken(comment.Text)
					}
				}

			}

			// tags
			tagValue := strings.Replace(tStructType.Tag(i), "`", "", -1)
			strPairs := strings.Split(tagValue, " ")
			for _, pair := range strPairs {
				if pair == "" {
					continue
				}

				tagPair := strings.Split(pair, ":")
				field.Tags[tagPair[0]] = strings.Replace(tagPair[1], "\"", "", -1)
			}

			// definition
			field.Name = tField.Name()
			fieldType := parseType(info, tField.Type())
			field.TypeName = fieldType.TypeName()
			field.TypeSpec = fieldType

			structType.Fields = append(structType.Fields, field)

		}

		iType = structType

	case *types.Slice:
		arrType := NewArrayType()
		eltType := parseType(info, t.(*types.Slice).Elem())
		arrType.EltSpec = eltType
		arrType.EltName = eltType.TypeName()
		arrType.Name = fmt.Sprintf("[]%s", eltType.TypeName())

		iType = arrType

	case *types.Map:
		mapType := NewMapType()
		tMap := t.(*types.Map)
		mapType.ValueSpec = parseType(info, tMap.Elem())
		mapType.KeySpec = parseType(info, tMap.Key())
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

//var stop bool
// struct expr匹配类型
func FindStructAstExprFromInfoTypes(info *types.Info, t *types.Struct) (expr ast.Expr) {
	for tExpr, tType := range info.Types {
		tStruct, ok := tType.Type.(*types.Struct)
		if !ok {
			continue
		}

		if t == tStruct { // 同一组astFiles生成的Types，内存中对象匹配成功
			expr = tExpr
			break
		}

		if t.String() == tStruct.String() { // 如果是不同的astFiles生成的Types，可能astFile1中没有这个类型信息，但是另外一组astFiles导入到info里，这是同一个类型，内存对象不一样，但是整体结构是一样的
			expr = tExpr
			break
		}

		if tStruct.NumFields() == t.NumFields() {
			numFields := tStruct.NumFields()
			notMatch := false
			for i := 0; i < numFields; i++ {
				if tStruct.Tag(i) != t.Tag(i) || tStruct.Field(i).Name() != t.Field(i).Name() {
					notMatch = true
					break
				}
			}

			if notMatch == false {
				expr = tExpr
				break
			}
		}

	}

	return
}

func RemoveCommentStartEndToken(text string) (newText string) {
	newText = strings.Replace(text, "//", "", 1)
	newText = strings.Replace(newText, "/*", "", 1)
	newText = strings.Replace(newText, "*/", "", 1)
	newText = strings.TrimSpace(newText)
	return
}
