package code

type ApiCode struct {
	Code    uint32
	Message string
}

func (m *ApiCode) Clone() *ApiCode {
	return &ApiCode{
		Code:    m.Code,
		Message: m.Message,
	}
}

func (m ApiCode) String() string {
	return m.Message
}

var (
	// 成功
	CodeSuccess = &ApiCode{
		Code:    0,
		Message: "success",
	}

	// 服务器错误
	CodeErrorServer = &ApiCode{
		Code:    1000,
		Message: "server error",
	}

	// 请求失败
	CodeErrorRequestFailed = &ApiCode{
		Code:    1001,
		Message: "request failed",
	}

	// 校验query参数失败
	CodeErrorQueryParams = &ApiCode{
		Code:    1002,
		Message: "verify query params failed",
	}

	// 校验path参数失败
	CodeErrorPathParams = &ApiCode{
		Code:    1003,
		Message: "verify path params failed",
	}

	// 校验post参数失败
	CodeErrorPostParams = &ApiCode{
		Code:    1004,
		Message: "verify post params failed",
	}

	// 校验header参数失败
	CodeErrorHeaderParams = &ApiCode{
		Code:    1005,
		Message: "verify header params failed",
	}

	// 数据库错误
	CodeErrorDB = &ApiCode{
		Code:    2000,
		Message: "db error",
	}

	// 数据库查询失败
	CodeErrorDBQueryFailed = &ApiCode{
		Code:    2001,
		Message: "db query failed",
	}

	// 数据库插入失败
	CodeErrorDBInsertFailed = &ApiCode{
		Code:    2002,
		Message: "db insert failed",
	}

	// 数据库更新失败
	CodeErrorDBUpdateFailed = &ApiCode{
		Code:    2003,
		Message: "db update failed",
	}

	// 数据库删除记录失败
	CodeErrorDBDeleteFailed = &ApiCode{
		Code:    2004,
		Message: "db delete failed",
	}

	// redis错误
	CodeErrorRedis = &ApiCode{
		Code:    2100,
		Message: "redis error",
	}

	// redis 获取失败
	CodeErrorRedisGetFailed = &ApiCode{
		Code:    2101,
		Message: "redis get record failed",
	}

	// redis 设置失败
	CodeErrorRedisSetFailed = &ApiCode{
		Code:    2102,
		Message: "redis set record failed",
	}

	// redis 删除失败
	CodeErrorRedisDeleteFailed = &ApiCode{
		Code:    2103,
		Message: "redis delete record failed",
	}
)
