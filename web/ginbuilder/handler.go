package ginbuilder

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
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
				ctx.Session.Panic(errPanic)
				if err == nil {
					err = uerrors.Newf("panic occurs: %#v", errPanic)
				}
			}
			defer func() { // 先处理，后抛出
				if errPanic != nil {
					logrus.Panic(errPanic)
				}
			}()

			ctx.Session.AfterHandle(err)

			// print error
			if err != nil {
				ctx.Logger.Errorf("error occurs: %#v", err)

				if ctx.ResponseData != nil {
					byteResponseData, errMarshal := json.Marshal(ctx.ResponseData)
					if errMarshal != nil {
						log.Printf("marshal response data failed.\n%#v\n", ctx.ResponseData)
					}

					log.Println(string(byteResponseData))
				}
			}

		}()

	}

	err = m.Handle(ctx)
	if nil != err {
		ctx.Logger.Errorf("handler %q failed. %s.", ginCtx.Request.RequestURI, err)
	}
}
