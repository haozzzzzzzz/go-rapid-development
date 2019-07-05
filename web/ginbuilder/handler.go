package ginbuilder

import (
	"encoding/json"

	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uruntime"
)

type HandleFunc struct {
	HttpMethod    string
	RelativePath  string   // 单个。废弃
	RelativePaths []string // 多个
	Handle        func(ctx *Context) (err error)
}

func (m *HandleFunc) GinHandler(ginCtx *gin.Context) {
	var err error
	ctx, err := NewContext(ginCtx)
	if nil != err {
		ctx.Logger.Errorf("new context failed. %s", err)
		return
	}

	defer func() {
		if errPanic := recover(); errPanic != nil {
			if ctx.Session != nil {
				ctx.Session.Panic(errPanic)
			}

			// print panic stack
			stack := uruntime.Stack(3)
			ctx.Logger.Error(string(stack))

			if err == nil {
				err = uerrors.Newf("panic occurs: %#v", errPanic)
			}

			ginCtx.AbortWithStatus(http.StatusInternalServerError)
		}

		// print error
		if err != nil {
			ctx.Logger.Errorf("error: %s", err)

			// dump request
			dumpRequest := &struct {
				Response interface{} `json:"response"`
			}{
				Response: ctx.ResponseData,
			}

			byteDumpInfo, errMarshal := json.Marshal(dumpRequest)
			if errMarshal != nil {
				ctx.Logger.Errorf("%#v", dumpRequest)
			}

			ctx.Logger.Error(string(byteDumpInfo))

		}

		if ctx.Session != nil {
			ctx.Session.AfterHandle(err)
		}

	}()

	if ctx.Session != nil {
		uri := ginCtx.Request.URL.Path
		for _, param := range ginCtx.Params {
			uri = strings.Replace(uri, param.Value, ":"+param.Key, 1)
		}
		ctx.Session.BeforeHandle(ctx, m.HttpMethod, uri)
	}

	err = m.Handle(ctx)
	if nil != err {
		ctx.Logger.Errorf("handler %q failed. %s.", ginCtx.Request.RequestURI, err)
	}
}
