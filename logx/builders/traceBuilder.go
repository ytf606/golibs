package builders

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/ytf606/golibs/logx/log4go"
	"github.com/ytf606/golibs/logx/logtrace"
	"github.com/ytf606/golibs/logx/logutils"
	"github.com/ytf606/golibs/logx/stackerr"
	"github.com/spf13/cast"
)

type TraceBuilder struct {
	department string
	version    string
}

func (this *TraceBuilder) LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{}) {
	if !logutils.ValidLevel(lvl) {
		return
	}
	if tag == "" {
		tag = "NoTagError"
	}

	if ctx == nil {
		id := strconv.FormatInt(logutils.GenLoggerId(), 10)
		ctx = context.WithValue(context.Background(), "logid", id)
	}

	if logutils.LogLevel(cast.ToString(ctx.Value("logLevel"))) > logutils.LogLevel(lvl) {
		return
	}

	tag = logutils.Filter(tag)
	position, message := this.Build(ctx, args, v...)

	if startValue := ctx.Value("start"); startValue != nil {
		if start, ok := startValue.(time.Time); ok {
			cost := time.Now().Sub(start)
			message = message + " COST:" + fmt.Sprintf("%.6f", cost.Seconds())
		}
	}
	traceNode := logtrace.ExtractTraceNodeFromContext(ctx)
	metadata := traceNode.ForkMap()
	metadata["x_department"] = this.department
	metadata["x_version"] = this.version
	metadata["x_module"] = `"` + tag + `"`
	metadata["x_org_id"] = `"` + cast.ToString(ctx.Value("org_id")) + `"`
	logtrace.IncrementRpcId(ctx)
	switch lvl {
	case "DEBUG":
		log4go.LogTraceMap(log4go.DEBUG, position, tag+"\t"+message, metadata)
	case "TRACE":
		log4go.LogTraceMap(log4go.TRACE, position, tag+"\t"+message, metadata)
	case "INFO":
		log4go.LogTraceMap(log4go.INFO, position, tag+"\t"+message, metadata)
	case "WARNING":
		log4go.LogTraceMap(log4go.WARNING, position, tag+"\t"+message, metadata)
	case "ERROR":
		log4go.LogTraceMap(log4go.ERROR, position, tag+"\t"+message, metadata)
	case "CRITICAL":
		log4go.LogTraceMap(log4go.CRITICAL, position, tag+"\t"+message, metadata)
	case "FATAL":
		log4go.LogTraceMap(log4go.CRITICAL, position, tag+"\t"+message, metadata)
		panic(message)
	}

}

func (this *TraceBuilder) Build(ctx context.Context, args interface{}, v ...interface{}) (position string, message string) {
	id := ctx.Value("logid")
	logid := cast.ToString(id)
	switch t := args.(type) {
	case *stackerr.StackErr:
		message = t.Info
		position = t.Position
	case error:
		message = t.Error()
	case string:
		if len(v) > 0 {
			message = fmt.Sprintf(t, v...)
		} else {
			message = t
		}
	default:
		message = fmt.Sprint(t)
	}

	if position == "" {
		_, file, line, ok := runtime.Caller(5)
		if ok {
			position = filepath.Base(file) + ":" + strconv.Itoa(line) + ":" + logid
		} else {
			position = "EMPTY"
		}
	}

	if hostname == "" {
		hostname, _ = os.Hostname()
	}

	position = position + "\t" + hostname
	message = logutils.Filter(message)
	return
}

func (this *TraceBuilder) SetTraceDepartment(department string) {
	this.department = `"` + department + `"`
}

func (this *TraceBuilder) SetTraceVersion(version string) {
	this.version = `"` + version + `"`
}
