package ginbuilder

import "github.com/gin-gonic/gin"

type HandleFunc struct {
	gin.HandlerFunc
	HttpMethod   string
	RelativePath string
}
