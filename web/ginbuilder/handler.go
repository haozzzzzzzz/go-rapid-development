package ginbuilder

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"net/http"

	"fmt"
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
		logrus.Errorf("new context failed. %s", err)
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
		relPaths := m.RelativePaths
		if len(relPaths) == 0 && m.RelativePath != "" {
			relPaths = []string{m.RelativePath}
		}

		for _, relPath := range relPaths {
			tempPath := relPath
			for _, param := range ginCtx.Params {
				tempPath = strings.Replace(tempPath, fmt.Sprintf(":%s", param.Key), param.Value, -1)
			}

			if strings.Contains(uri, tempPath) {
				uri = relPath
				break
			}
		}

		ctx.Session.BeforeHandle(ctx, m.HttpMethod, uri)
	}

	err = m.Handle(ctx)
	if nil != err {
		ctx.Logger.Errorf("handler %q failed. %s.", ginCtx.Request.RequestURI, err)
	}
}
