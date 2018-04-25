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
	ctx := NewContext(ginCtx)
	err := m.Handle(ctx)
	if nil != err {
		ctx.Logger.Errorf("handler %q failed. \n%s.", ginCtx.Request.RequestURI, err)
	}
}
