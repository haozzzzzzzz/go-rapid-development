package wgin

import (
	"github.com/gin-gonic/gin"
	"github.com/xiam/to"
)

const (
	DataKeyResponseRetCode string = "__resp_code__"
)

func SetResponseRetCode(ginCtx *gin.Context, code int) {
	ginCtx.Set(DataKeyResponseRetCode, code)
}

func GetResponseRetCode(ginCtx *gin.Context) (code int, ok bool) {
	iCode, ok := ginCtx.Get(DataKeyResponseRetCode)
	if !ok {
		return
	}

	code = to.Int(iCode)
	return
}
