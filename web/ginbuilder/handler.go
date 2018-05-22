package ginbuilder

import (
	"github.com/gin-gonic/gin"
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
		ctx.Session.BeforeHandle(ctx)
		defer func() {
			if errPanic := recover(); errPanic != nil {
				ctx.Session.Panic(errPanic)
			}
		}()
		defer func() {
			ctx.Session.AfterHandle(err)
		}()
	}

	err = m.Handle(ctx)
	if nil != err {
		ctx.Logger.Errorf("handler %q failed. %s.", ginCtx.Request.RequestURI, err)
	}
}
