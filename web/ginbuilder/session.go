package ginbuilder

import (
	"github.com/haozzzzzzzz/go-rapid-development/api/code"
)

type SessionBuilderFunc func(ctx *Context) error

type Session interface {
	BeforeHandle(ctx *Context, method string, uri string)
	AfterHandle(err error)
	Panic(err interface{})
	SetReturnCode(code *code.ApiCode)
}
