package parser

type Field struct {
	Name     string            `json:"name" yaml:"name"`
	TypeName string            `json:"type_name" yaml:"type_name"`
	Tags     map[string]string `json:"tags" yaml:"tags"`
	TypeSpec interface{}       `json:"type_spec" form:"type_spec"`
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
const TypeClassStandardType = "standard"
const TypeClassStructType = "struct"
const TypeClassMapType = "map"
const TypeClassArrayType = "array"
const TypeClassInterfaceType = "interface"

// 标准类型
type StandardType struct {
	TypeClass string `json:"type_class" yaml:"type_class"`
	Name      string `json:"name" yaml:"name"`
}

func (m StandardType) TypeName() string {
	return string(m.Name)
}

func NewStandardType(name string) *StandardType {
	return &StandardType{
		Name:      name,
		TypeClass: TypeClassStandardType,
	}
}

// struct
type StructType struct {
	TypeClass string   `json:"type_class" yaml:"type_class"`
	Name      string   `json:"name" yaml:"name"`
	Fields    []*Field `json:"fields" yaml:"fields"`
}

func (m *StructType) TypeName() string {
	return m.Name
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
	TypeClass string      `json:"type_class" yaml:"type_class"`
	Name      string      `json:"name" yaml:"name"`
	Key       string      `json:"key" yaml:"key"`
	ValueSpec interface{} `json:"value_spec" yaml:"value_spec"`
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
	TypeClass string      `json:"type_class" yaml:"type_class"`
	Name      string      `json:"name" yaml:"name"`
	EltName   string      `json:"elt_name" yaml:"elt_name"`
	EltSpec   interface{} `json:"elt_spec" yaml:"elt_spec"`
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
	ApiHandlerFunc    string      `validate:"required" json:"api_handler_func" yaml:"api_handler_func"`
	ApiHandlerPackage string      `validate:"required" json:"api_handler_package_func" yaml:"api_handler_package_func"`
	SourceFile        string      `validate:"required" json:"source_file" yaml:"source_file"`
	HttpMethod        string      `validate:"required" json:"http_method" yaml:"http_method"`
	RelativePath      string      `validate:"required" json:"relative_path" yaml:"relative_path"`
	RelativePackage   string      `json:"relative_package" yaml:"relative_package"`
	PathData          *StructType `json:"path_data" yaml:"path_data"`
	QueryData         *StructType `json:"query_data" yaml:"query_data"`
	PostData          *StructType `json:"post_data" yaml:"post_data"`
	RespData          *StructType `json:"response_data" yaml:"response_data"`
}
