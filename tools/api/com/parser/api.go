package parser

import (
	"fmt"
	"strings"
)

type Field struct {
	Name        string            `json:"name" yaml:"name"`
	TypeName    string            `json:"type_name" yaml:"type_name"`
	Tags        map[string]string `json:"tags" yaml:"tags"`
	TypeSpec    IType             `json:"type_spec" form:"type_spec"`
	Description string            `json:"description" form:"description"`
}

func (m *Field) TagJson() (name string) {
	return m.Tags["json"]
}

func (m *Field) Required() (required bool) {
	strBind := m.Tags["binding"]
	if strings.Contains(strBind, "required") {
		required = true
	}
	return
}

func NewField() *Field {
	return &Field{
		Tags: make(map[string]string),
	}
}

// type interface
type IType interface {
	TypeName() string
}

// 类型分类
const TypeClassBasicType = "basic"
const TypeClassStructType = "struct"
const TypeClassMapType = "map"
const TypeClassArrayType = "array"
const TypeClassInterfaceType = "interface"

// 标准类型
type BasicType struct {
	TypeClass string `json:"type_class" yaml:"type_class"`
	Name      string `json:"name" yaml:"name"`
}

func (m BasicType) TypeName() string {
	return string(m.Name)
}

func NewBasicType(name string) *BasicType {
	return &BasicType{
		Name:      name,
		TypeClass: TypeClassBasicType,
	}
}

// struct
type StructType struct {
	TypeClass   string   `json:"type_class" yaml:"type_class"`
	Name        string   `json:"name" yaml:"name"`
	Fields      []*Field `json:"fields" yaml:"fields"`
	Description string   `json:"description" yaml:"description"`
}

func (m *StructType) TypeName() string {
	return m.Name
}

func (m *StructType) AddFields(filed *Field) {
	m.Fields = append(m.Fields, filed)
}

func NewStructType() *StructType {
	return &StructType{
		TypeClass: TypeClassStructType,
		Name:      TypeClassStructType,
		Fields:    make([]*Field, 0),
	}
}

// map
type MapType struct {
	TypeClass string `json:"type_class" yaml:"type_class"`
	Name      string `json:"name" yaml:"name"`
	KeySpec   IType  `json:"key" yaml:"key"`
	ValueSpec IType  `json:"value_spec" yaml:"value_spec"`
}

func (m *MapType) TypeName() string {
	return m.Name
}

func NewMapType() *MapType {
	return &MapType{
		TypeClass: TypeClassMapType,
		Name:      TypeClassMapType,
	}
}

// array
type ArrayType struct {
	TypeClass string `json:"type_class" yaml:"type_class"`
	Name      string `json:"name" yaml:"name"`
	EltName   string `json:"elt_name" yaml:"elt_name"`
	EltSpec   IType  `json:"elt_spec" yaml:"elt_spec"`
}

func (m *ArrayType) TypeName() string {
	return m.Name
}

func NewArrayType() *ArrayType {
	return &ArrayType{
		TypeClass: TypeClassArrayType,
		Name:      TypeClassArrayType,
	}
}

// interface
type InterfaceType struct {
	TypeClass string `json:"type_class" yaml:"type_class"`
}

func NewInterfaceType() *InterfaceType {
	return &InterfaceType{
		TypeClass: TypeClassInterfaceType,
	}
}

func (m *InterfaceType) TypeName() string {
	return m.TypeClass
}

type ApiItem struct {
	ApiHandlerFunc string `validate:"required" json:"api_handler_func" yaml:"api_handler_func"` // handler 函数名
	SourceFile     string `validate:"required" json:"source_file" yaml:"source_file"`           // 源码

	PackageName         string `validate:"required" json:"api_handler_package" yaml:"api_handler_package"` // handler 所在的包名
	PackageExportedPath string `json:"package_exported" yaml:"package_exported"`                           // 被别人引入的path
	PackageDir          string `json:"package_dir" yaml:"package_dir"`                                     // 包所在的路径
	PackageRelAlias     string `json:"package_rel_alias" yaml:"package_rel_alias"`                         // 包昵称

	HttpMethod    string      `validate:"required" json:"http_method" yaml:"http_method"`
	RelativePaths []string    `validate:"required" json:"relative_paths" yaml:"relative_paths"`
	PathData      *StructType `json:"path_data" yaml:"path_data"`
	QueryData     *StructType `json:"query_data" yaml:"query_data"`
	PostData      *StructType `json:"post_data" yaml:"post_data"`
	RespData      *StructType `json:"response_data" yaml:"response_data"`

	Summary     string `json:"summary" yaml:"summary"`
	Description string `json:"description" yaml:"description"`
}

func (m *ApiItem) PackageFuncName() string {
	return fmt.Sprintf("%s.%s", m.PackageRelAlias, m.ApiHandlerFunc)
}

// 成功返回结构
//type Response struct {
//	ReturnCode uint32      `json:"ret"`
//	Message    string      `json:"msg"`
//	Data       interface{} `json:"data"`
//}
func SuccessResponseStructType(
	respData *StructType,
) (successResp *StructType) {
	successResp = &StructType{
		TypeClass:   TypeClassStructType,
		Name:        "SuccessResponse",
		Description: "api success response",
		Fields: []*Field{
			{
				Name:        "ReturnCode",
				TypeName:    "uint32",
				Description: "result code",
				Tags: map[string]string{
					"json": "ret",
				},
				TypeSpec: NewBasicType("uint32"),
			},
			{
				Name:        "Message",
				TypeName:    "string",
				Description: "result message",
				Tags: map[string]string{
					"json": "msg",
				},
				TypeSpec: NewBasicType("string"),
			},
			{
				Name:        "Data",
				TypeName:    respData.TypeName(),
				Description: "result data",
				Tags: map[string]string{
					"json": "data",
				},
				TypeSpec: respData,
			},
		},
	}
	return
}
