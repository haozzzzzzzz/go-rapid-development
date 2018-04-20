package ginbuilder

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Context struct {
	GinContext *gin.Context
	Logger     *logrus.Logger
	Session    *Session
}

func NewContext(ginContext *gin.Context, ask interface{}) (ctx *Context) {
	ctx = &Context{
		GinContext: ginContext,
	}

	return
}
