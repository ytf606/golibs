package ginx

import (
	"git.100tal.com/wangxiao_monkey_tech/lib/errorx"
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
	app.NoMethod(NoMethodHandler())
	app.NoRoute(NoRouteHandler())
	app.Use(mws...)

	return app
}

func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ErrResponse(c, errorx.ErrMethodNotAllow)
	}
}

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ErrResponse(c, errorx.ErrNotFound)
	}
}
