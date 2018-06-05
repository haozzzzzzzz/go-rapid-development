package ginbuilder

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
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
	Logger     *logrus.Entry
	Session    Session
}

func NewContext(ginContext *gin.Context) (ctx *Context, err error) {
	ctx = &Context{
		GinContext: ginContext,
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

func (m *Context) BindQueryData(queryData interface{}) (code *ReturnCode, err error) {
	err = m.GinContext.ShouldBindQuery(queryData)
	if err != nil {
		code = CodeErrorQueryParams.Clone()
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

func (m *Context) BindPostData(postData interface{}) (code *ReturnCode, err error) {
	err = m.GinContext.MustBindWith(postData, binding.JSON)
	if err != nil {
		code = CodeErrorPostParams.Clone()
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

func (m *Context) BindPathData(pathData interface{}) (code *ReturnCode, err error) {
	defer func() {
		if err != nil {
			code = CodeErrorPathParams.Clone()
			validateErrors, ok := err.(validator.ValidationErrors)
			if ok {
				for _, fieldError := range validateErrors {
					code.Message = fmt.Sprintf("%s. %q:%s", code.Message, fieldError.Name, fieldError.Tag)
					break
				}
			}
		}
	}()
	err = bindParams(m.GinContext.Params, pathData)
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

func (m *Context) Send(code *ReturnCode, obj interface{}) {
	response := NewReponse(code, obj)
	m.GinContext.JSON(http.StatusOK, response)
	if m.Session != nil {
		m.Session.SetReturnCode(code)
	}
	return
}

func (m *Context) Success() {
	m.Send(CodeSuccess.Clone(), nil)
}

func (m *Context) SuccessReturn(obj interface{}) {
	m.Send(CodeSuccess.Clone(), obj)
}

func (m *Context) Error(code *ReturnCode, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Error(logArgs...)
}

func (m *Context) Warn(code *ReturnCode, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Warn(logArgs...)
}

func (m *Context) Errorf(code *ReturnCode, logFormat string, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Errorf(logFormat, logArgs...)
}

func (m *Context) Warnf(code *ReturnCode, logFormat string, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Warnf(logFormat, logArgs...)
}
