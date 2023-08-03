package ginx

/// utils package for gin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContextKey interface{}

// StdCtx - convert gin.Context to context.Context
func StdCtx(ctx *gin.Context) context.Context {
	c := context.Background()

	for k, v := range ctx.Keys {
		c = context.WithValue(c, k, v)
	}

	return c
}

// GinHandler 将http.HandlerFunc转为gin.HandlerFunc
func GinHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
