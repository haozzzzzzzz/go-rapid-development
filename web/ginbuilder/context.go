package ginbuilder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type Context struct {
	GinContext *gin.Context
	Logger     *logrus.Entry
	Session    *Session
}

func NewContext(ginContext *gin.Context) (ctx *Context) {
	ctx = &Context{
		GinContext: ginContext,
		Logger:     logrus.WithFields(logrus.Fields{}),
		Session:    &Session{},
	}
	return
}

func (m *Context) BindQueryData(queryData interface{}) (err error) {
	err = m.GinContext.ShouldBindQuery(queryData)
	return
}

func (m *Context) BindPostData(postData interface{}) (err error) {
	err = m.GinContext.MustBindWith(postData, binding.JSON)
	return
}

func (m *Context) BindPathData(pathData interface{}) (err error) {
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
	return
}

func (m *Context) Success() {
	m.Send(CodeSuccess, nil)
}

func (m *Context) SuccessReturn(obj interface{}) {
	m.Send(CodeSuccess, obj)
}

func (m *Context) Error(code *ReturnCode, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Error(logArgs...)
}

func (m *Context) Errorf(code *ReturnCode, logFormat string, logArgs ...interface{}) {
	m.Send(code, nil)
	m.Logger.Errorf(logFormat, logArgs...)
}
