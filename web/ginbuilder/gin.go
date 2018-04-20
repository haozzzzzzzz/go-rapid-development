package ginbuilder

import (
	"github.com/gin-gonic/gin"
)

var engine *gin.Engine

func GetEngine() *gin.Engine {
	if engine == nil {
		engine = gin.Default()
	}

	return engine
}
