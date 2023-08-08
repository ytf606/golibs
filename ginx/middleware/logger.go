package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"git.100tal.com/wangxiao_monkey_tech/lib/logx"
	"git.100tal.com/wangxiao_monkey_tech/lib/logx/logtrace"
	"github.com/gin-gonic/gin"
)

var (
	hostname, _ = os.Hostname()
)

// 创建logger相关
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("logid", strconv.FormatInt(logx.Id(), 10))
		ctx.Set("hostname", hostname)
		ctx.Set("start", time.Now())
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		var body []byte
		if ctx.Request.Body != nil {
			body, _ = ioutil.ReadAll(ctx.Request.Body)
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if raw != "" {
			path = path + "?" + raw
		}

		logtraceMap := logtrace.GenLogTraceMetadata()
		logtraceMap.Set("request_uri", fmt.Sprintf("\"%s\"", path))
		if len(body) > 0 && bytes.HasPrefix(body, []byte("{")) {
			logtraceMap.Set("x_request_param", fmt.Sprintf("%s", body))
		} else {
			logtraceMap.Set("request_param", fmt.Sprintf("\"%s\"", body))
		}
		logtraceMap.Set("request_method", fmt.Sprintf("\"%s\"", ctx.Request.Method))
		logtraceMap.Set("request_client_ip", fmt.Sprintf("\"%s\"", ctx.ClientIP()))
		if traceId := ctx.GetHeader("traceid"); traceId != "" {
			logtraceMap.Set("x_trace_id", "\""+traceId+"\"")
			if strings.HasPrefix(traceId, "pts_") {
				ctx.Set("IS_BENCHMARK", "1")
			}
		}
		if traceId := ctx.GetHeader("traceId"); traceId != "" {
			logtraceMap.Set("x_trace_id", "\""+traceId+"\"")
		}
		if rpcId := ctx.GetHeader("rpcid"); rpcId != "" {
			rpcId = rpcId + ".0"
			logtraceMap.Set("x_rpcid", "\""+rpcId+"\"")
		}
		if rpcId := ctx.GetHeader("rpcId"); rpcId != "" {
			rpcId = rpcId + ".0"
			logtraceMap.Set("x_rpcid", "\""+rpcId+"\"")
		}
		ctx.Set(logtrace.GetMetadataKey(), logtraceMap)
	}
}

// Logger -
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Response
		blw := &bodyLogWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = blw

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		var body []byte
		if c.Request.Body != nil {
			body, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Process request
		c.Next()

		_, skip := c.Get("SKIPLOG")
		if skip {
			return
		}

		var userID string
		if val, exists := c.Get("uid"); exists {
			if _, ok := val.(int); ok {
				userID = fmt.Sprintf("%d", val)
			} else if _, ok := val.(int64); ok {
				userID = fmt.Sprintf("%d", val)
			} else {
				userID = fmt.Sprintf("%s", val)
			}
		}

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		buf := blw.body.Bytes()

		_, cut := c.Get("CUTREQBODY")

		if cut {
			logx.Ix(c, "[GIN]", "%v | %3d | %13v | %15s | {userId:%v} | %-7s %s  %s  %s  %s",
				end.Format("2006/01/02 - 15:04:05"), statusCode, latency, clientIP, userID, method, path, comment, body[0:256], buf,
			)
			return
		}

		logx.Ix(c, "[GIN]", "%v | %3d | %13v | %15s | {userId:%v} | %-7s %s  %s  %s  %s",
			end.Format("2006/01/02 - 15:04:05"), statusCode, latency, clientIP, userID, method, path, comment, body, buf,
		)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if _, err := w.body.Write(b); err != nil {
		fmt.Printf("bodyLogWriter err:%v", err)
	}
	return w.ResponseWriter.Write(b)
}
