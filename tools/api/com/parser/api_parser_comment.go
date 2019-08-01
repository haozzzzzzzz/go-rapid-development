package parser

import (
	"os"

	"go/parser"
	"go/token"

	"regexp"

	"encoding/json"

	"reflect"

	"strings"

	"fmt"

	"go/ast"

	"path/filepath"

	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/sirupsen/logrus"
)

func ParseApisFromComments(dir string) (
	title string,
	description string,
	version string,
	contact string,
	apis []*ApiItem,
	err error,
) {
	title = "api service name"
	description = "api service description"
	version = "1.0"
	contact = "contact"
	apis = make([]*ApiItem, 0)

	pkgDirs := make([]string, 0)
	pkgDirs, err = file.SearchFileNames(dir, func(fileInfo os.FileInfo) bool {
		if fileInfo.IsDir() {
			return true
		} else {
			return false
		}
	}, true)
	pkgDirs = append(pkgDirs, dir)

	for _, pkgDir := range pkgDirs {
		subApis, errParse := ParseApisFromPkgComment(pkgDir)
		err = errParse
		if nil != err {
			logrus.Errorf("parse apis from pkg comment failed. pkgDir: %s, error: %s.", pkgDir, err)
			return
		}

		apis = append(apis, subApis...)
	}

	return
}

func ParseApisFromPkgComment(pkgDir string) (apis []*ApiItem, err error) {
	apis = make([]*ApiItem, 0)
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, pkgDir, nil, parser.ParseComments)
	if nil != err {
		logrus.Errorf("parse pkg dir failed. pkgDir: %s, error: %s.", pkgDir, err)
		return
	}

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Files {
			for _, commentGroup := range astFile.Comments {
				for _, comment := range commentGroup.List {
					tempApis, errParse := ParseApisFromPkgCommentText(
						fileSet,
						astFile,
						comment,
					)
					err = errParse
					if nil != err {
						logrus.Errorf("parse apis from pkg comment text failed. text: %s, error: %s.", comment.Text, err)
						return
					}

					apis = append(apis, tempApis...)
				}
			}
		}
	}
	return
}

func ParseApisFromPkgCommentText(
	fileSet *token.FileSet,
	astFile *ast.File,
	comment *ast.Comment,
) (apis []*ApiItem, err error) {
	apis = make([]*ApiItem, 0)

	docReg, err := regexp.Compile(`(?si:@api_doc_start(.*?)@api_doc_end)`)
	if nil != err {
		logrus.Errorf("reg compile pkg api comment text failed. error: %s.", err)
		return
	}

	arrStrs := docReg.FindAllStringSubmatch(comment.Text, -1)
	strJsons := make([]string, 0)
	for _, strs := range arrStrs {
		strJsons = append(strJsons, strs[1])
	}

	tokenFile := fileSet.File(astFile.Pos())
	fileName := tokenFile.Name()
	fileDir := filepath.Dir(fileName)
	for _, strJson := range strJsons {
		tempApiItem, errParse := parseCommentTextToApi(strJson)
		err = errParse
		if nil != err {
			logrus.Errorf("parse comment to to api failed. %s, error: %s.", strJson, err)
			return
		}

		if tempApiItem == nil {
			continue
		}

		tempApiItem.PackageName = astFile.Name.Name
		tempApiItem.SourceFile = fileName
		tempApiItem.PackageDir = fileDir

		apis = append(apis, tempApiItem)
	}

	return
}

/**
@api_doc_start
{
	"http_method": "GET",
	"relative_paths": ["/hello_world"],
	"query_data": {
		"name": "姓名|string|required"
	},
	"post_data": {
		"location": "地址|string|required"
	},
	"resp_data": {
	    "a": "a|int",
	    "b": "b|int",
	    "c": {
			"d": "d|string"
		},
		"__c": "c|object",
		"f": [
			"string"
		],
		"__f": "f|object|required",
		"g": [
			{
				"h": "h|string|required"
			}
		],
		"__g": "g|array|required"
	}
}
@api_doc_end
*/

type CommentTextApi struct {
	HttpMethod    string                 `json:"http_method"`
	RelativePaths []string               `json:"relative_paths"`
	PathData      map[string]interface{} `json:"path_data"`
	QueryData     map[string]interface{} `json:"query_data"`
	PostData      map[string]interface{} `json:"post_data"`
	RespData      map[string]interface{} `json:"resp_data"`
}

func (m *CommentTextApi) IsEmpty() bool {
	return m.HttpMethod == "" || len(m.RelativePaths) == 0
}

func parseCommentTextToApi(
	text string,
) (api *ApiItem, err error) {
	comApi := &CommentTextApi{}
	err = json.Unmarshal([]byte(text), comApi)
	if nil != err {
		logrus.Errorf("unmarshal api failed. error: %s.", err)
		return
	}

	if comApi.IsEmpty() {
		logrus.Warnf("found empty api doc comment, require http_method and relative_paths")
		return
	}

	api = &ApiItem{
		HttpMethod:    comApi.HttpMethod,
		RelativePaths: comApi.RelativePaths,
	}

	api.PathData, err = commentApiRequestDataToStructType(comApi.PathData)
	if nil != err {
		logrus.Errorf("comment text api path data to struct type failed. error: %s.", err)
		return
	}

	api.QueryData, err = commentApiRequestDataToStructType(comApi.QueryData)
	if nil != err {
		logrus.Errorf("comment text api query data to struct type failed. error: %s.", err)
		return
	}

	api.PostData, err = commentApiRequestDataToStructType(comApi.PostData)
	if nil != err {
		logrus.Errorf("comment text api post data to struct type failed. error: %s.", err)
		return
	}

	api.RespData, err = commentApiRequestDataToStructType(comApi.RespData)
	if nil != err {
		logrus.Errorf("comment text api resp data to struct type failed. error: %s.", err)
		return
	}

	return
}

func commentApiRequestDataToStructType(
	mapData map[string]interface{},
) (structType *StructType, err error) {
	if mapData == nil {
		return
	}

	structType = NewStructType()
	for key, typeDesc := range mapData {
		if strings.HasPrefix(key, "__") {
			continue
		}

		keyDesc := mapData[fmt.Sprintf("__%s", key)] // 如果是嵌套类型，则会有一个__key描述这个field在当前struct的属性
		strKeyDesc, _ := keyDesc.(string)
		field, errField := commentApiRequestDataFieldDesc(key, typeDesc, strKeyDesc)
		err = errField
		if nil != err {
			logrus.Errorf("get struct field failed. key: %s, type_desc: %#v error: %s.", key, typeDesc, err)
			return
		}

		structType.AddFields(field)
	}
	return
}

func commentApiRequestDataFieldDesc(key string, fieldTypeDesc interface{}, slaveFieldDesc string) (field *Field, err error) {
	if fieldTypeDesc == nil {
		return
	}

	field = NewField()
	field.Name = key

	strFieldTypeDesc, ok := fieldTypeDesc.(string)
	if !ok {
		strFieldTypeDesc = slaveFieldDesc
	}

	if strFieldTypeDesc != "" {
		vals := [3]string{} // description, type, tags
		splitDefs := strings.Split(strFieldTypeDesc, "|")
		for i := 0; i < 3 && i < len(splitDefs); i++ {
			vals[i] = splitDefs[i]
		}

		field.Description = vals[0]
		field.TypeName = vals[1]
		field.TypeSpec = NewBasicType(vals[1])
		field.Tags["json"] = key
		field.Tags["binding"] = vals[2]
	}

	field.TypeSpec, err = commentApiRequestDataIType(fieldTypeDesc)
	if nil != err {
		logrus.Errorf("parse field type spec failed. error: %s.", err)
		return
	}

	return
}

func commentApiRequestDataIType(
	typeDesc interface{},
) (itype IType, err error) {
	reflectType := reflect.TypeOf(typeDesc)
	switch reflectType.Kind() {
	case reflect.String:
		itype = NewBasicType(typeDesc.(string))

	case reflect.Map:
		mTypeDesc, ok := typeDesc.(map[string]interface{})
		if !ok {
			logrus.Warnf("convert def to map type failed. typeDesc: %#v", typeDesc)
			return
		}

		itype, err = commentApiRequestDataToStructType(mTypeDesc)
		if nil != err {
			logrus.Errorf("field map type desc to struct type failed. error: %s.", err)
			return
		}

	case reflect.Slice:
		sliceTypeDesc, ok := typeDesc.([]interface{})
		if !ok {
			logrus.Warnf("convert def to slice type failed. typeDesc: %#v", typeDesc)
			return
		}

		sliceType := NewArrayType()
		if len(sliceTypeDesc) > 0 {
			sliceType.EltSpec, err = commentApiRequestDataIType(sliceTypeDesc[0])
			if nil != err {
				logrus.Errorf("parse slice type elt spec failed. error: %s.", err)
				return
			}
		}

		itype = sliceType

	default:
		err = uerrors.Newf("unsupported type: %#v", typeDesc)

	}
	return
}
