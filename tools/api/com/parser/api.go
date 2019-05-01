package parser

type Field struct {
	Name     string            `json:"name" yaml:"name"`
	Type     string            `json:"type" yaml:"type"`
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

// 标准类型
type StandardType string

func (m StandardType) TypeName() string {
	return string(m)
}

// struct
type StructType struct {
	Name   string   `json:"name" yaml:"name"`
	Fields []*Field `json:"fields" yaml:"fields"`
}

func (m *StructType) TypeName() string {
	return m.Name
}

func NewStructType() *StructType {
	return &StructType{
		Fields: make([]*Field, 0),
	}
}

// map
type MapType struct {
	Name      string      `json:"name" yaml:"name"`
	Key       string      `json:"key" yaml:"key"`
	ValueSpec interface{} `json:"value_spec" yaml:"value_spec"`
}

func (m *MapType) TypeName() string {
	return m.Name
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
