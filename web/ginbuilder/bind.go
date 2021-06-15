package ginbuilder

import (
	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-rapid-development/v2/api/request"
)

// Deprecated: new gin has full support of bind
func BindParams(params gin.Params, v interface{}) error {

	form := make(map[string][]string)

	for _, param := range params {
		form[param.Key] = []string{param.Value}
	}

	return request.StructMapForm(v, form, "form")

}
