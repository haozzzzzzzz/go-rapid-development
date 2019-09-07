package ginbuilder

import (
	"github.com/gin-gonic/gin"
)

// engine
func DefaultEngine() (engine *gin.Engine) {
	engine = gin.New()
	engine.Use(LogAccessMiddleware(), gin.Recovery())
	return
}
