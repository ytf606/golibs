package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	AllowAllOrigins bool `ini:"allowAllOrigins"`

	// AllowOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// Default value is []
	AllowOrigins []string `ini:"allowOrigins"`

	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowOrigins is ignored.
	AllowOriginFunc func(origin string) bool

	// AllowMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET and POST)
	AllowMethods []string `ini:"allowMethods"`

	// AllowHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	AllowHeaders []string `ini:"allowHeaders"`

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool `ini:"allowCredentials"`

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposeHeaders []string `ini:"exposeHeaders"`

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge time.Duration `ini:"maxAge"`

	// Allows to add origins like http://some-domain/*, https://api.* or http://some.*.subdomain.com
	AllowWildcard bool `ini:"allowWildcard"`

	// Allows usage of popular browser extensions schemas
	AllowBrowserExtensions bool `ini:"allowBrowserExtensions"`

	// Allows usage of WebSocket protocol
	AllowWebSockets bool `ini:"allowWebSockets"`

	// Allows usage of file:// schema (dangerous!) use it only when you 100% sure it's needed
	AllowFiles bool `ini:"allowFiles"`
}

func CorsMiddleware(config Config) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
	}
	return cors.New(corsConfig)
}

func NoCacheMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-store")
	}
}

// Options -
func OptionsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToLower(ctx.Request.Method) == "options" {
			ctx.String(204, "ok")
			ctx.Abort()
		}
	}
}

func ResponseHeaderMiddleware(label string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", label+"/"+hostname)
	}
}
