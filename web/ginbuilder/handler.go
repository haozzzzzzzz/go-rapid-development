package ginbuilder

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HandleFunc struct {
	HttpMethod   string
	RelativePath string
	Handle       func(ctx *Context) (err error)
}

func (m *HandleFunc) GinHandler(ginCtx *gin.Context) {
	var err error
	ctx, err := NewContext(ginCtx)
	if nil != err {
		ctx.Logger.Errorf("new context failed. %s", err)
		return
	}

	if ctx.Session != nil {
		ctx.Session.BeforeHandle(ctx, m.HttpMethod, m.RelativePath)
		defer func() {
			var errPanic interface{}
			if errPanic = recover(); errPanic != nil {
				logrus.Panic(errPanic)
				ctx.Session.Panic(errPanic)
			}

			ctx.Session.AfterHandle(err)

		}()

	}

	err = m.Handle(ctx)
	if nil != err {
		ctx.Logger.Errorf("handler %q failed. %s.", ginCtx.Request.RequestURI, err)
	}
}
