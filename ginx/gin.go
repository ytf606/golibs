package ginx

import (
	"github.com/gin-gonic/gin"
)

type (
	Context     = gin.Context
	HandlerFunc = gin.HandlerFunc
	Engine      = gin.Engine
	RouterGroup = gin.RouterGroup
)

func New(mode string, mws ...HandlerFunc) *Engine {
	gin.SetMode(mode)

	app := gin.New()
	app.HandleMethodNotAllowed = true
	app.Use(mws...)

	return app
}
