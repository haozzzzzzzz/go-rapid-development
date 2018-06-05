package api

type StructDataField struct {
	Name          string
	Type          string
	SubStructData *StructData // 允许多重嵌套
}

type StructData struct {
	Name   string
	Fields []*StructDataField
}

type ApiItem struct {
	ApiHandlerFunc    string `validate:"required" json:"api_handler_func" yaml:"api_handler_func"`
	ApiHandlerPackage string `validate:"required" json:"api_handler_package_func" yaml:"api_handler_package_func"`
	SourceFile        string `validate:"required" json:"source_file" yaml:"source_file"`

	HttpMethod   string      `validate:"required" json:"http_method" yaml:"http_method"`
	RelativePath string      `validate:"required" json:"relative_path" yaml:"relative_path"`
	PathData     *StructData `json:"path_data" yaml:"path_data"`
	QueryData    *StructData `json:"query_data" yaml:"query_data"`
	PostData     *StructData `json:"post_data" yaml:"post_data"`
	RespData     *StructData `json:"response_data" yaml:"response_data"`
}
