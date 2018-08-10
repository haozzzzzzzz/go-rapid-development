package ginbuilder

type SessionBuilderFunc func(ctx *Context) error

type Session interface {
	BeforeHandle(ctx *Context, method string, uri string)
	AfterHandle(err error)
	Panic(err interface{})
}
