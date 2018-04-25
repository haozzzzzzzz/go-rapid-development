package ginbuilder

// 返回状态码
type ReturnCode struct {
	uint32
	Message string
}

var (
	CodeNormal = ReturnCode{
		uint32:  0,
		Message: "normal",
	}

	CodeVerifyParamsFailed = ReturnCode{
		uint32:  1000,
		Message: "verify params failed",
	}
)
