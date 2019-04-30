package parser

type StructDataField struct {
	Name          string            `json:"name" yaml:"name"`
	Type          string            `json:"type" yaml:"type"`
	Tags          map[string]string `json:"tags" yaml:"tags"`
	SubStructData *StructData       `json:"sub_struct_data" yaml:"sub_struct_data"` // 允许多重嵌套
}

func NewStructDataField() *StructDataField {
	return &StructDataField{
		Tags: make(map[string]string),
	}
}

type StructData struct {
	Name   string             `json:"name" yaml:"name"`
	Fields []*StructDataField `json:"fields" yaml:"fields"`
}

func NewStructData() *StructData {
	return &StructData{
		Fields: make([]*StructDataField, 0),
	}
}

type ApiItem struct {
	ApiHandlerFunc    string      `validate:"required" json:"api_handler_func" yaml:"api_handler_func"`
	ApiHandlerPackage string      `validate:"required" json:"api_handler_package_func" yaml:"api_handler_package_func"`
	SourceFile        string      `validate:"required" json:"source_file" yaml:"source_file"`
	HttpMethod        string      `validate:"required" json:"http_method" yaml:"http_method"`
	RelativePath      string      `validate:"required" json:"relative_path" yaml:"relative_path"`
	RelativePackage   string      `json:"relative_package" yaml:"relative_package"`
	PathData          *StructData `json:"path_data" yaml:"path_data"`
	QueryData         *StructData `json:"query_data" yaml:"query_data"`
	PostData          *StructData `json:"post_data" yaml:"post_data"`
	RespData          *StructData `json:"response_data" yaml:"response_data"`
}
