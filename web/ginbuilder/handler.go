package ginbuilder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HandleFunc struct {
	HttpMethod   string
	RelativePath string
	Handle       func(ctx *Context) (response interface{}, err error)
}

func (m *HandleFunc) GinHandler(ginCtx *gin.Context) {
	ctx := NewContext(ginCtx)
	response, err := m.Handle(ctx)
	if nil != err {
		logrus.Errorf("handler %q failed. \n%s.", ginCtx.Request.RequestURI, err)
	}

	if response != nil {
		ctx.GinContext.JSON(http.StatusOK, response)
	}
}
