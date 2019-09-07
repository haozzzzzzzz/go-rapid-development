package ginbuilder

import (
	"github.com/sirupsen/logrus"
)

type SessionBuilderFunc func(ctx *Context) error

type Session interface {
	BeforeHandle(ctx *Context, method string, uri string)
	AfterHandle(err error)
	Panic(err interface{})
}

var sessionBuilder SessionBuilderFunc

func BindSessionBuilder(sesBuilder SessionBuilderFunc) {
	if sessionBuilder != nil {
		logrus.Fatalf("session builder was bound")
		return
	}

	sessionBuilder = sesBuilder
}

func buildSessionForCtx(ctx *Context) (err error) {
	if sessionBuilder != nil {
		err = sessionBuilder(ctx)
		if nil != err {
			logrus.Errorf("session builder error. %s.", err)
			return
		}
	}

	return
}
