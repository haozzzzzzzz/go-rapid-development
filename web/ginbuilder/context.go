package ginbuilder

import (
	"fmt"
	"net/http"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
	"github.com/haozzzzzzzz/go-rapid-development/api/code"
)

var sessionBuilder SessionBuilderFunc

func BindSessionBuilder(sesBuilder SessionBuilderFunc) {
	if sessionBuilder != nil {
		logrus.Fatalf("session builder was bound")
		return
	}

	sessionBuilder = sesBuilder
}

type Context struct {
	GinContext *gin.Context
	RequestCtx context.Context
	Logger     *logrus.Entry
	Session    Session
}

func NewContext(ginContext *gin.Context) (ctx *Context, err error) {
	ctx = &Context{
		GinContext: ginContext,
		RequestCtx: ginContext.Request.Context(),
		Logger:     logrus.WithFields(logrus.Fields{}),
	}

	if sessionBuilder != nil {
		err = sessionBuilder(ctx)
		if nil != err {
			logrus.Errorf("session builder error. %s.", err)
			return
		}
	}

	return
}

func (m *Context) BindQueryData(queryData interface{}) (code *code.ApiCode, err error) {
	err = m.GinContext.ShouldBindQuery(queryData)
	if err != nil {
		code = code.CodeErrorQueryParams.Clone()
		validateErrors, ok := err.(validator.ValidationErrors)
		if ok {
			for _, fieldError := range validateErrors {
				code.Message = fmt.Sprintf("%s. %q:%s", code.Message, fieldError.Name, fieldError.Tag)
				break
			}
		}
	}
	return
}

func (m *Context) BindPostData(postData interface{}) (code *code.ApiCode, err error) {
	err = m.GinContext.MustBindWith(postData, binding.JSON)
	if err != nil {
		code = code.CodeErrorPostParams.Clone()
		validateErrors, ok := err.(validator.ValidationErrors)
		if ok {
			for _, fieldError := range validateErrors {
				code.Message = fmt.Sprintf("%s. %q:%s", code.Message, fieldError.Name, fieldError.Tag)
				break
			}
		}
	}
	return
}

func (m *Context) BindPathData(pathData interface{}) (code *code.ApiCode, err error) {
	defer func() {
		if err != nil {
			code = code.CodeErrorPathParams.Clone()
			validateErrors, ok := err.(validator.ValidationErrors)
			if ok {
				for _, fieldError := range validateErrors {
					code.Message = fmt.Sprintf("%s. %q:%s", code.Message, fieldError.Name, fieldError.Tag)
					break
				}
			}
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

func (m *Context) Send(code *code.ApiCode, obj interface{}) {
	response := NewResponse(code, obj)
	m.GinContext.JSON(http.StatusOK, response)
	if m.Session != nil {
		m.Session.SetReturnCode(code)
	}
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

func (m *Context) Warnf(code *code.ApiCode, logFormat string, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Warnf(logFormat, logArgs...)
}

func (m *Context) WarnfReturn(code *code.ApiCode, obj interface{}, logFormat string, logArgs ...interface{}) {
	m.Send(code, obj)
	m.Logger.Warnf(logFormat, logArgs...)
}

func (m *Context) TemporaryRedirect(location string) {
	m.GinContext.Redirect(http.StatusTemporaryRedirect, location)
}

func (m *Context) PermanentRedirect(location string) {
	m.GinContext.Redirect(http.StatusPermanentRedirect, location)
}
