package ginbuilder

type SessionBuilderFunc func(ctx *Context) error

type Session interface {
	BeforeHandle(ctx *Context)
	AfterHandle(err error)
	Panic(err interface{})
}
