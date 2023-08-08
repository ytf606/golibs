package builders

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"git.100tal.com/wangxiao_monkey_tech/lib/logx/log4go"
	"git.100tal.com/wangxiao_monkey_tech/lib/logx/logutils"
	"git.100tal.com/wangxiao_monkey_tech/lib/logx/stackerr"
	"github.com/spf13/cast"
)

type DefaultBuilder struct {
}

var hostname string

func (this *DefaultBuilder) LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{}) {
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

	tag = logutils.Filter(tag)
	position, message := this.Build(ctx, args, v...)

	if startValue := ctx.Value("start"); startValue != nil {
		if start, ok := startValue.(time.Time); ok {
			cost := time.Now().Sub(start)
			message = message + " COST:" + fmt.Sprintf("%.2f", cost.Seconds()*1e3)
		}
	}

	switch lvl {
	case "DEBUG":
		log4go.Log(log4go.DEBUG, position, tag+"\t"+message)
	case "TRACE":
		log4go.Log(log4go.TRACE, position, tag+"\t"+message)
	case "INFO":
		log4go.Log(log4go.INFO, position, tag+"\t"+message)
	case "WARNING":
		log4go.Log(log4go.WARNING, position, tag+"\t"+message)
	case "ERROR":
		log4go.Log(log4go.ERROR, position, tag+"\t"+message)
	case "CRITICAL":
		log4go.Log(log4go.CRITICAL, position, tag+"\t"+message)
	case "FATAL":
		log4go.Log(log4go.CRITICAL, position, tag+"\t"+message)
		panic(message)
	}
}

func (this *DefaultBuilder) Build(ctx context.Context, args interface{}, v ...interface{}) (position string, message string) {
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
