package ginbuilder

// 请求返回基类
type ResponseBase struct {
	ReturnCode uint32 `json:"ret"`
	Message    string `json:"msg"`
}

func NewResponseBaseByReturnCode(code *ReturnCode) *ResponseBase {
	return &ResponseBase{
		ReturnCode: code.uint32,
		Message:    code.Message,
	}
}

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
)
