package ginbuilder

import (
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/v2/utils/uerrors"
	"github.com/haozzzzzzzz/go-rapid-development/v2/utils/uruntime"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HandleFunc struct {
	HttpMethod string
	// Deprecated: please use RelativePaths for multi url
	RelativePath  string   // 单个。废弃
	RelativePaths []string // 多个
	BeforeHandle  func(ctx *Context) (stop bool, err error)
	Handle        func(ctx *Context) (err error)
	AfterHandle   func(ctx *Context, handleErr error) (err error)
}

func (m *HandleFunc) GinHandler(ginCtx *gin.Context) {
	defer func() {
		if iRec := recover(); iRec != nil {
			logrus.Errorf("gin handler panic : %s.", iRec)
		}
	}()

	var err error
	ginReq := ginCtx.Request
	ctx, err := NewContext(ginCtx)
	if nil != err {
		logrus.Errorf("new context failed. %s", err)
		return
	}

	if m.BeforeHandle != nil {
		stop, errBef := m.BeforeHandle(ctx)
		err = errBef
		if err != nil {
			ctx.Logger.Errorf("before handle %s failed. %s", ginReq.URL, err)
		}

		if stop {
			return
		}
	}

	err = buildSessionForCtx(ctx)
	if nil != err {
		ctx.Logger.Errorf("build session for ctx failed. %s", err)
		return
	}

	defer func() {
		if iRec := recover(); iRec != nil {
			if ctx.Session != nil {
				ctx.Session.Panic(iRec)
			}

			if err == nil {
				err = uerrors.Newf("gin handler panic: %s", iRec)
				ctx.Logger.Error(err)
			}

			// print panic stack
			stack := uruntime.Stack(3)
			ctx.Logger.Error(string(stack))

			ginCtx.AbortWithStatus(http.StatusInternalServerError)
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

	if m.AfterHandle != nil {
		err = m.AfterHandle(ctx, err)
		if nil != err {
			ctx.Logger.Errorf("after handle %s failed. %s", ginReq.RequestURI, err)
		}
	}
}
