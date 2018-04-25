package ginbuilder

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type Context struct {
	GinContext *gin.Context
	Logger     *logrus.Logger
	Session    *Session
}

func NewContext(ginContext *gin.Context) (ctx *Context) {
	ctx = &Context{
		GinContext: ginContext,
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
