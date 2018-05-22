package ginbuilder

import (
	"github.com/haozzzzzzzz/go-rapid-development/api/code"
)

// 返回状态码
type ReturnCode = code.ApiCode

var (
	// 成功
	CodeSuccess = code.CodeSuccess

	// 服务器错误
	CodeErrorServer = code.CodeErrorServer

	// 请求失败
	CodeErrorRequestFailed = code.CodeErrorRequestFailed

	// 校验query参数失败
	CodeErrorQueryParams = code.CodeErrorQueryParams

	// 校验path参数失败
	CodeErrorPathParams = code.CodeErrorPathParams

	// 校验post参数失败
	CodeErrorPostParams = code.CodeErrorPostParams
)
