package code

import (
	"fmt"
)

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
	return fmt.Sprintf("Code: %d; Message: %s", m.Code, m.Message)
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

	// token无效
	CodeErrorToken = &ApiCode{
		Code:    1006,
		Message: "verify token failed",
	}

	// 请求被拒绝
	CodeErrorRequestRejected = &ApiCode{
		Code:    1007,
		Message: "request was rejected",
	}

	// session无效
	CodeErrorSession = &ApiCode{
		Code:    1008,
		Message: "invalid session",
	}

	// 请求参数错误
	CodeErrorRequestParams = &ApiCode{
		Code:    1009,
		Message: "verify request params failed",
	}

	// 请求过于频繁
	CodeErrorRequestFrequently = &ApiCode{
		Code:    1010,
		Message: "request frequently",
	}

	// proxy
	CodeErrorProxyFailed = &ApiCode{
		Code:    1100,
		Message: "proxy failed",
	}

	CodeErrorProxyRequestFailed = &ApiCode{
		Code:    1101,
		Message: "request porxy failed",
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

	// 没有符合条件的记录
	CodeErrorDBQueryNoRecords = &ApiCode{
		Code:    2005,
		Message: "db query no records",
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

	// dynamodb 错误
	CodeErrorDynamodb = &ApiCode{
		Code:    2200,
		Message: "dynamodb error",
	}

	// dynamodb 获取失败
	CodeErrorDynamodbGetFailed = &ApiCode{
		Code:    2201,
		Message: "dynamodb get record failed",
	}

	// dynamodb 设置失败
	CodeErrorDynamodbSetFailed = &ApiCode{
		Code:    2202,
		Message: "dynamodb set record failed",
	}

	// dynamodb 删除失败
	CodeErrorDynamodbDeleteFailed = &ApiCode{
		Code:    2203,
		Message: "dynamodb delete record failed",
	}

	// dynamodb unmarshal失败
	CodeErrorDynamodbAttributeUnmarshalFailed = &ApiCode{
		Code:    2204,
		Message: "dynamodb unmarshal failed",
	}
)
