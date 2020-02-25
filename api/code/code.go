package code

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type ApiCode struct {
	Code    uint32
	Message string
	Err     error
}

func (m *ApiCode) Clone() *ApiCode {
	return &ApiCode{
		Code:    m.Code,
		Message: m.Message,
	}
}

func (m *ApiCode) WithMessage(msg string) *ApiCode {
	apiCode := m.Clone()
	apiCode.Message = msg
	return apiCode
}

func (m ApiCode) String() string {
	strErr := ""
	if m.Err != nil {
		strErr = m.Err.Error()
	}
	return fmt.Sprintf("Code: %d; Message: %s; Err: %s", m.Code, m.Message, strErr)
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

	// 服务器错误
	CodeErrorServer = NewAddApiCodePanic(1000, "server error")

	// 请求失败
	CodeErrorRequestFailed = NewAddApiCodePanic(1001, "request failed")

	// 校验query参数失败
	CodeErrorQueryParams = NewAddApiCodePanic(1002, "verify query params failed")

	// 校验path参数失败
	CodeErrorPathParams = NewAddApiCodePanic(1003, "verify uri params failed")

	CodeErrorUriParams = CodeErrorPathParams

	// 校验post参数失败
	CodeErrorPostParams = NewAddApiCodePanic(1004, "verify post params failed")

	// 校验header参数失败
	CodeErrorHeaderParams = NewAddApiCodePanic(1005, "verify header params failed")

	// token无效
	CodeErrorToken = NewAddApiCodePanic(1006, "verify token failed")

	// 请求被拒绝
	CodeErrorRequestRejected = NewAddApiCodePanic(1007, "request was rejected")

	// session无效
	CodeErrorSession = NewAddApiCodePanic(1008, "invalid session")

	// 请求参数错误
	CodeErrorRequestParams = NewAddApiCodePanic(1009, "verify request params failed")

	// 请求过于频繁
	CodeErrorRequestFrequently = NewAddApiCodePanic(1010, "request frequently")

	// 请求条件不符
	CodeErrorConditionCheckFailed = NewAddApiCodePanic(1011, "condition check failed")

	// 请求资源不存在
	CodeErrorResourceNotExists = NewAddApiCodePanic(1012, "resource not exists")

	CodeErrorPostFormParams = NewAddApiCodePanic(1013, "verify post form params failed")

	// 请求参数不匹配
	CodeErrorRequestParamNotMatch = NewAddApiCodePanic(1014, "request params not match")

	// 需要授权
	CodeErrorAuthRequired = NewAddApiCodePanic(1015, "require auth")

	// json body参数
	CodeErrorBodyParams = NewAddApiCodePanic(1016, "verify body params failed")

	// 用户信息存在
	CodeErrorUserNotExists = NewAddApiCodePanic(1017, "user not exists")

	// 用户信息访问失败
	CodeErrorGetUserInfoFailed = NewAddApiCodePanic(1018, "get user info failed")

	// proxy
	CodeErrorProxyFailed = NewAddApiCodePanic(1100, "proxy failed")

	// 请求api失败
	CodeErrorProxyRequestFailed = NewAddApiCodePanic(1101, "request porxy failed")

	// 数据库错误
	CodeErrorDB = NewAddApiCodePanic(2000, "db error")

	// 数据库查询失败
	CodeErrorDBQueryFailed = NewAddApiCodePanic(2001, "db query failed")

	// 数据库插入失败
	CodeErrorDBInsertFailed = NewAddApiCodePanic(2002, "db insert failed")

	// 数据库更新失败
	CodeErrorDBUpdateFailed = NewAddApiCodePanic(2003, "db update failed")

	// 数据库删除记录失败
	CodeErrorDBDeleteFailed = NewAddApiCodePanic(2004, "db delete failed")

	// 没有符合条件的记录
	CodeErrorDBQueryNoRecords = NewAddApiCodePanic(2005, "db query no records")

	// 记录已存在
	CodeErrorDBRecordExists = NewAddApiCodePanic(2006, "db record exists")

	// redis错误
	CodeErrorRedis = NewAddApiCodePanic(2100, "redis error")

	// redis 获取失败
	CodeErrorRedisGetFailed = NewAddApiCodePanic(2101, "redis get record failed")

	// redis 设置失败
	CodeErrorRedisSetFailed = NewAddApiCodePanic(2102, "redis set record failed")

	// redis 删除失败
	CodeErrorRedisDeleteFailed = NewAddApiCodePanic(2103, "redis delete record failed")

	// dynamodb 错误
	CodeErrorDynamodb = NewAddApiCodePanic(2200, "dynamodb error")

	// dynamodb 获取失败
	CodeErrorDynamodbGetFailed = NewAddApiCodePanic(2201, "dynamodb get record failed")

	// dynamodb 设置失败
	CodeErrorDynamodbSetFailed = NewAddApiCodePanic(2202, "dynamodb set record failed")

	// dynamodb 删除失败
	CodeErrorDynamodbDeleteFailed = NewAddApiCodePanic(2203, "dynamodb delete record failed")

	// dynamodb unmarshal失败
	CodeErrorDynamodbAttributeUnmarshalFailed = NewAddApiCodePanic(2204, "dynamodb unmarshal failed")

	// dynamodb 没有符合条件的记录
	CodeErrorDynamodbQueryNoRecords = NewAddApiCodePanic(2205, "dynamodb query no records")

	// es
	CodeErrorElasticsearch = NewAddApiCodePanic(2300, "elasticserach error")

	CodeErrorElasticsearchCreate = NewAddApiCodePanic(2301, "elasticsearch create error")

	CodeErrorElasticsearchSearch = NewAddApiCodePanic(2302, "elasticsearch search error")
)
