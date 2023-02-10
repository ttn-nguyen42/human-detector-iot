package routes

import (
	"github.com/gin-gonic/gin"
)

func Create(engine *gin.Engine) {
	// Enables logging for HTTP server
	engine.Use(gin.Logger())
	// Returns 500 on panic()
	engine.Use(gin.Recovery())
	
}