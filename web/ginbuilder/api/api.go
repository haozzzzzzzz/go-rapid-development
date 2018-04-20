package api

type ApiItem struct {
	HttpMethod        string `validate:"required" json:"http_method" yaml:"http_method"`
	RelativePath      string `validate:"required" json:"relative_path" yaml:"relative_path"`
	ApiHandlerFunc    string `validate:"required" json:"api_handler_func" yaml:"api_handler_func"`
	ApiHandlerPackage string `validate:"required" json:"api_handler_package_func" yaml:"api_handler_package_func"`
	SourceFile        string `validate:"required" json:"source_file" yaml:"source_file"`
}
