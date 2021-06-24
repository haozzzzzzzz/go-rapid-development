# ginbuilder
An gin api builder with another api style, it's light-extended and 100% compatible with [gin](https://github.com/gin-gonic/gin).

## ginbuilder api style
```go
// SayHiApi Greeting
var SayHiApi ginbuilder.HandleFunc = ginbuilder.HandleFunc{
	HttpMethod: "GET",
	RelativePaths: []string{
		"/api/app/say_hi",
		"/api/web/say_hi",
	},
	Handle: func(ctx *ginbuilder.Context) (err error) {
		ctx.Success()
		return
	},
}
```

## api tool
[Api tool](https://github.com/haozzzzzzzz/go-tool) can auto bind your api to gin engine, and it also can generate swagger doc.

```shell
api --help # see help
api compile --help # see compile help
api swagger --help # see generate doc help
```
