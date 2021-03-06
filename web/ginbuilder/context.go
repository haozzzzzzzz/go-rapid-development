package ginbuilder

import (
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/v2/web/wgin"
	"net/http"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/haozzzzzzzz/go-rapid-development/v2/api/code"
	"github.com/sirupsen/logrus"
)

const TRACE_REQUEST_KEY = "trace_fields"
const NO_ACCESS_LOG_PRINT = "no_access_log_print" // key for not print access log

// Context context
type Context struct {
	GinContext   *gin.Context
	RequestCtx   context.Context
	Logger       *logrus.Entry
	Session      Session
	ResponseData *Response
}

func NewContext(ginContext *gin.Context) (ctx *Context, err error) {
	ctx = &Context{
		GinContext: ginContext,
		RequestCtx: ginContext.Request.Context(),
		Logger:     logrus.WithFields(logrus.Fields{}),
	}
	return
}

func (m *Context) SetTraceFields(fields logrus.Fields) {
	m.Logger = m.Logger.WithFields(fields)
	m.GinContext.Set(TRACE_REQUEST_KEY, m.Logger.Data)
}

func (m *Context) BindQueryData(queryData interface{}) (retCode *code.ApiCode, err error) {
	err = m.GinContext.ShouldBindQuery(queryData)
	if err != nil {
		retCode = code.CodeErrorQueryParams.Clone()
		retCode.Message = fmt.Sprintf("%s. %s", retCode.Message, err.Error())
	}
	return
}

func (m *Context) BindPostData(postData interface{}) (retCode *code.ApiCode, err error) {
	err = m.GinContext.ShouldBindWith(postData, binding.JSON) // json only
	if err != nil {
		retCode = code.CodeErrorPostParams.Clone()
		retCode.Message = fmt.Sprintf("%s. %s", retCode.Message, err.Error())
	}
	return
}

func (m *Context) BindBodyData(bodyData interface{}) (retCode *code.ApiCode, err error) {
	err = m.GinContext.ShouldBindWith(bodyData, binding.JSON) // json only
	if err != nil {
		retCode = code.CodeErrorPostParams.Clone()
		retCode.Message = fmt.Sprintf("%s. %s", retCode.Message, err.Error())
	}
	return
}

// BindCacheBodyData read body, and reuse it in request
func (m *Context) BindCacheBodyData(bodyData interface{}) (retCode *code.ApiCode, err error) {
	err = m.GinContext.ShouldBindBodyWith(bodyData, binding.JSON) // json only
	if err != nil {
		retCode = code.CodeErrorPostParams.Clone()
		retCode.Message = fmt.Sprintf("%s. %s", retCode.Message, err.Error())
	}
	return
}

func (m *Context) GetCachedBody() (body []byte, exists bool) {
	var iBody interface{}
	iBody, exists = m.GinContext.Get(gin.BodyBytesKey)
	if exists {
		body, exists = iBody.([]byte)
	}
	return
}

// BindPostForm form-urlencoded
func (m *Context) BindPostForm(postData interface{}) (err error) {
	err = m.GinContext.ShouldBindWith(postData, binding.FormPost)
	if nil != err {
		logrus.Errorf("gin should bind with form post failed. error: %s.", err)
		return
	}
	return
}

// Deprecated: use BindUriData which uses gin's bind.
func (m *Context) BindPathData(pathData interface{}) (retCode *code.ApiCode, err error) {
	defer func() {
		if err != nil {
			retCode = code.CodeErrorUriParams.Clone()
			retCode.Message = fmt.Sprintf("%s. %s", retCode.Message, err.Error())
		}
	}()

	err = BindParams(m.GinContext.Params, pathData)
	if nil != err {
		logrus.Errorf("bind path data failed. \n%s.", err)
		return
	}

	err = binding.Validator.ValidateStruct(pathData)
	if nil != err {
		logrus.Errorf("validate path data failed. \n%s.", err)
		return
	}

	return
}

// BindUriData use form tag `uri`
func (m *Context) BindUriData(uriData interface{}) (retCode *code.ApiCode, err error) {
	err = m.GinContext.ShouldBindUri(uriData)
	if nil != err {
		retCode = code.CodeErrorUriParams.Clone()
		retCode.Message = fmt.Sprintf("%s. %s", retCode.Message, err.Error())
	}
	return
}

func (m *Context) BindHeaderData(
	headerData interface{},
) (retCode *code.ApiCode, err error) {
	err = m.GinContext.ShouldBindHeader(headerData)
	if nil != err {
		retCode = code.CodeErrorHeaderParams.Clone()
		retCode.Message = fmt.Sprintf("%s. %s", retCode.Message, err.Error())
	}
	return
}

func (m *Context) Send(code *code.ApiCode, obj interface{}) {
	m.ResponseData = NewResponse(code, obj)
	m.GinContext.JSON(http.StatusOK, m.ResponseData)
	wgin.SetResponseRetCode(m.GinContext, int(code.Code))
	return
}

func (m *Context) Success() {
	m.Send(code.CodeSuccess.Clone(), nil)
}

func (m *Context) SuccessReturn(obj interface{}) {
	m.Send(code.CodeSuccess.Clone(), obj)
}

func (m *Context) Error(code *code.ApiCode, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Error(logArgs...)
}

func (m *Context) Warn(code *code.ApiCode, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Warn(logArgs...)
}

func (m *Context) Errorf(code *code.ApiCode, logFormat string, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Errorf(logFormat, logArgs...)
}

func (m *Context) ErrorfReturn(code *code.ApiCode, obj interface{}, logFormat string, logArgs ...interface{}) {
	m.Send(code, obj)
	m.Logger.Errorf(logFormat, logArgs...)
}

func (m *Context) Warnf(code *code.ApiCode, logFormat string, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Warnf(logFormat, logArgs...)
}

func (m *Context) WarnfReturn(code *code.ApiCode, obj interface{}, logFormat string, logArgs ...interface{}) {
	m.Send(code, obj)
	m.Logger.Warnf(logFormat, logArgs...)
}

// TemporaryRedirect 永久性定向. 301
// 该状态码表示请求的资源已被分配了新的URI，以后应使用资源现在所指的URI
func (m *Context) TemporaryRedirect(location string) {
	m.GinContext.Redirect(http.StatusTemporaryRedirect, location)
}

// PermanentRedirect 临时重定向.307
// 该状态码与302有相同的含义。307会遵照浏览器标准，不会从post变为get。但是对于处理响应时的行为，各种浏览器有可能出现不同的情况。
func (m *Context) PermanentRedirect(location string) {
	m.GinContext.Redirect(http.StatusPermanentRedirect, location)
}

// StatusFoundRedirect 临时性重定向. 302
// 该状态码表示请求的资源已被分配了新的URI，希望用户（本次）能使用新的URI访问。
func (m *Context) StatusFoundRedirect(location string) {
	m.GinContext.Redirect(http.StatusFound, location)
}

func (m *Context) StatusNotFoundWarnf(logFormat string, logArgs ...interface{}) {
	m.GinContext.Status(http.StatusNotFound)
	m.Logger.Errorf(logFormat, logArgs...)
}

func (m *Context) String(result string) {
	m.GinContext.String(http.StatusOK, result)
}

func (m *Context) StringErrorf(result string, logFormat string, logArgs ...interface{}) {
	m.GinContext.String(http.StatusOK, result)
	m.Logger.Errorf(logFormat, logArgs...)
}

func (m *Context) JSON(result interface{}) {
	m.GinContext.JSON(http.StatusOK, result)
}

func (m *Context) JSONErrorf(result interface{}, logFormat string, logArgs ...interface{}) {
	m.GinContext.JSON(http.StatusOK, result)
	m.Logger.Errorf(logFormat, logArgs...)
}
