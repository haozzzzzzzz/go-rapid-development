package ginbuilder

import (
	"github.com/haozzzzzzzz/go-rapid-development/v2/api/code"
)

// Response 请求返回基类
type Response struct {
	ReturnCode uint32      `json:"ret"`
	Message    string      `json:"msg"`
	Data       interface{} `json:"data"`
}

func (m *Response) GetApiCode() *code.ApiCode {
	return &code.ApiCode{
		Code:    m.ReturnCode,
		Message: m.Message,
	}
}

func NewResponse(code *code.ApiCode, data interface{}) *Response {
	if data == nil {
		data = make(map[string]interface{})
	}

	return &Response{
		ReturnCode: code.Code,
		Message:    code.Message,
		Data:       data,
	}
}
