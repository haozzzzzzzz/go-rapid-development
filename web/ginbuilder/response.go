package ginbuilder

// 请求返回基类
type Response struct {
	ReturnCode uint32      `json:"ret"`
	Message    string      `json:"msg"`
	Data       interface{} `json:"data"`
}

func NewResponse(code *ReturnCode, data interface{}) *Response {
	if data == nil {
		data = make(map[string]interface{})
	}

	return &Response{
		ReturnCode: code.Code,
		Message:    code.Message,
		Data:       data,
	}
}
