package code

import (
	"fmt"
	"github.com/sirupsen/logrus"
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

// 状态码
var ApiCodeMap = map[uint32]*ApiCode{}

func NewAddApiCodePanic(code uint32, msg string) (apiCode *ApiCode) {
	apiCode = &ApiCode{
		Code:    code,
		Message: msg,
	}

	AddApiCodePanic(apiCode)
	return
}

func AddApiCodePanic(apiCode *ApiCode) {
	_, ok := ApiCodeMap[apiCode.Code]
	if ok {
		logrus.Panicf("api code exists. code: %s", apiCode)
		return
	}

	ApiCodeMap[apiCode.Code] = apiCode
}

var (
	// 成功
	CodeSuccess = NewAddApiCodePanic(0, "success")
	//&ApiCode{
	//	Code:    0,
	//	Message: "success",
	//}

	// 服务器错误
	CodeErrorServer = NewAddApiCodePanic(1000, "server error")
	//&ApiCode{
	//	Code:    1000,
	//	Message: "server error",
	//}

	// 请求失败
	CodeErrorRequestFailed = NewAddApiCodePanic(1001, "request failed")
	//&ApiCode{
	//	Code:    1001,
	//	Message: "request failed",
	//}

	// 校验query参数失败
	CodeErrorQueryParams = NewAddApiCodePanic(1002, "verify query params failed")
	//&ApiCode{
	//	Code:    1002,
	//	Message: "verify query params failed",
	//}

	// 校验path参数失败
	CodeErrorPathParams = NewAddApiCodePanic(1003, "verify uri params failed")
	//&ApiCode{
	//	Code:    1003,
	//	Message: "verify uri params failed",
	//}

	CodeErrorUriParams = CodeErrorPathParams

	// 校验post参数失败
	CodeErrorPostParams = NewAddApiCodePanic(1004, "verify post params failed")
	//&ApiCode{
	//	Code:    1004,
	//	Message: "verify post params failed",
	//}

	// 校验header参数失败
	CodeErrorHeaderParams = NewAddApiCodePanic(1005, "verify header params failed")
	//&ApiCode{
	//	Code:    1005,
	//	Message: "verify header params failed",
	//}

	// token无效
	CodeErrorToken = NewAddApiCodePanic(1006, "verify token failed")
	//&ApiCode{
	//	Code:    1006,
	//	Message: "verify token failed",
	//}

	// 请求被拒绝
	CodeErrorRequestRejected = NewAddApiCodePanic(1007, "request was rejected")
	//&ApiCode{
	//	Code:    1007,
	//	Message: "request was rejected",
	//}

	// session无效
	CodeErrorSession = NewAddApiCodePanic(1008, "invalid session")
	//&ApiCode{
	//	Code:    1008,
	//	Message: "invalid session",
	//}

	// 请求参数错误
	CodeErrorRequestParams = NewAddApiCodePanic(1009, "verify request params failed")
	//&ApiCode{
	//	Code:    1009,
	//	Message: "verify request params failed",
	//}

	// 请求过于频繁
	CodeErrorRequestFrequently = NewAddApiCodePanic(1010, "request frequently")
	//&ApiCode{
	//	Code:    1010,
	//	Message: "request frequently",
	//}

	// 请求条件不符
	CodeErrorConditionCheckFailed = NewAddApiCodePanic(1011, "condition check failed")
	//&ApiCode{
	//	Code:    1011,
	//	Message: "condition check failed",
	//}

	// 请求资源不存在
	CodeErrorResourceNotExists = NewAddApiCodePanic(1012, "resource not exists")
	//&ApiCode{
	//	Code:    1012,
	//	Message: "resource not exists",
	//}

	CodeErrorPostFormParams = NewAddApiCodePanic(1013, "verify post form params failed")
	//&ApiCode{
	//	Code:    1013,
	//	Message: "verify post form params failed",
	//}

	// 请求参数不匹配
	CodeErrorRequestParamNotMatch = NewAddApiCodePanic(1014, "request params not match")
	//&ApiCode{
	//	Code:    1014,
	//	Message: "request params not match",
	//}

	// 需要授权
	CodeErrorAuthRequired = NewAddApiCodePanic(1015, "require auth")
	//&ApiCode{
	//	Code:    1015,
	//	Message: "require auth",
	//}

	// proxy
	CodeErrorProxyFailed = NewAddApiCodePanic(1100, "proxy failed")
	//&ApiCode{
	//	Code:    1100,
	//	Message: "proxy failed",
	//}

	// 请求api失败
	CodeErrorProxyRequestFailed = NewAddApiCodePanic(1101, "request porxy failed")
	//&ApiCode{
	//	Code:    1101,
	//	Message: "request porxy failed",
	//}

	// 数据库错误
	CodeErrorDB = NewAddApiCodePanic(2000, "db error")
	//&ApiCode{
	//	Code:    2000,
	//	Message: "db error",
	//}

	// 数据库查询失败
	CodeErrorDBQueryFailed = NewAddApiCodePanic(2001, "db query failed")
	//&ApiCode{
	//	Code:    2001,
	//	Message: "db query failed",
	//}

	// 数据库插入失败
	CodeErrorDBInsertFailed = NewAddApiCodePanic(2002, "db insert failed")
	//&ApiCode{
	//	Code:    2002,
	//	Message: "db insert failed",
	//}

	// 数据库更新失败
	CodeErrorDBUpdateFailed = NewAddApiCodePanic(2003, "db update failed")
	//&ApiCode{
	//	Code:    2003,
	//	Message: "db update failed",
	//}

	// 数据库删除记录失败
	CodeErrorDBDeleteFailed = NewAddApiCodePanic(2004, "db delete failed")
	//&ApiCode{
	//	Code:    2004,
	//	Message: "db delete failed",
	//}

	// 没有符合条件的记录
	CodeErrorDBQueryNoRecords = NewAddApiCodePanic(2005, "db query no records")
	//&ApiCode{
	//	Code:    2005,
	//	Message: "db query no records",
	//}

	// 记录已存在
	CodeErrorDBRecordExists = NewAddApiCodePanic(2006, "db record exists")
	//&ApiCode{
	//	Code:    2006,
	//	Message: "db record exists",
	//}

	// redis错误
	CodeErrorRedis = NewAddApiCodePanic(2100, "redis error")
	//&ApiCode{
	//	Code:    2100,
	//	Message: "redis error",
	//}

	// redis 获取失败
	CodeErrorRedisGetFailed = NewAddApiCodePanic(2101, "redis get record failed")
	//&ApiCode{
	//	Code:    2101,
	//	Message: "redis get record failed",
	//}

	// redis 设置失败
	CodeErrorRedisSetFailed = NewAddApiCodePanic(2102, "redis set record failed")
	//&ApiCode{
	//	Code:    2102,
	//	Message: "redis set record failed",
	//}

	// redis 删除失败
	CodeErrorRedisDeleteFailed = NewAddApiCodePanic(2103, "redis delete record failed")
	//&ApiCode{
	//	Code:    2103,
	//	Message: "redis delete record failed",
	//}

	// dynamodb 错误
	CodeErrorDynamodb = NewAddApiCodePanic(2200, "dynamodb error")
	//&ApiCode{
	//	Code:    2200,
	//	Message: "dynamodb error",
	//}

	// dynamodb 获取失败
	CodeErrorDynamodbGetFailed = NewAddApiCodePanic(2201, "dynamodb get record failed")
	//&ApiCode{
	//	Code:    2201,
	//	Message: "dynamodb get record failed",
	//}

	// dynamodb 设置失败
	CodeErrorDynamodbSetFailed = NewAddApiCodePanic(2202, "dynamodb set record failed")
	//&ApiCode{
	//	Code:    2202,
	//	Message: "dynamodb set record failed",
	//}

	// dynamodb 删除失败
	CodeErrorDynamodbDeleteFailed = NewAddApiCodePanic(2203, "dynamodb delete record failed")
	//&ApiCode{
	//	Code:    2203,
	//	Message: "dynamodb delete record failed",
	//}

	// dynamodb unmarshal失败
	CodeErrorDynamodbAttributeUnmarshalFailed = NewAddApiCodePanic(2204, "dynamodb unmarshal failed")
	//&ApiCode{
	//	Code:    2204,
	//	Message: "dynamodb unmarshal failed",
	//}

	// es
	CodeErrorElasticsearch = NewAddApiCodePanic(2300, "elasticserach error")
	//&ApiCode{
	//	Code:    2300,
	//	Message: "elasticserach error",
	//}

	CodeErrorElasticsearchCreate = NewAddApiCodePanic(2301, "elasticsearch create error")
	//&ApiCode{
	//	Code:    2301,
	//	Message: "elasticsearch create error",
	//}

	CodeErrorElasticsearchSearch = NewAddApiCodePanic(2302, "elasticsearch search error")
	//&ApiCode{
	//	Code:    2302,
	//	Message: "elasticsearch search error",
	//}
)
