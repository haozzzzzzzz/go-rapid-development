package ginbuilder

// 返回状态码
type ReturnCode struct {
	Code    uint32
	Message string
}

func (m ReturnCode) String() string {
	return m.Message
}

var (
	// 成功
	CodeSuccess = &ReturnCode{
		Code:    0,
		Message: "success",
	}

	// 服务器错误
	CodeErrorServer = &ReturnCode{
		Code:    1000,
		Message: "server error",
	}

	// 请求失败
	CodeErrorFailed = &ReturnCode{
		Code:    1001,
		Message: "request failed",
	}

	// 校验query参数失败
	CodeErrorQueryParams = &ReturnCode{
		Code:    1002,
		Message: "verify query params failed",
	}

	// 校验path参数失败
	CodeErrorPathParams = &ReturnCode{
		Code:    1003,
		Message: "verify path params failed",
	}

	// 校验post参数失败
	CodeErrorPostParams = &ReturnCode{
		Code:    1004,
		Message: "verify post params failed",
	}
)
