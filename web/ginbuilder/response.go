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
