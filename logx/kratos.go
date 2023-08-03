package logx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

var ins *xlog

type xlog struct {
	log MessageBuilder
}

func NewLogger(build MessageBuilder) *xlog {
	if ins != nil {
		return ins
	}
	ins = &xlog{
		log: build,
	}
	return ins
}

func GetLogger() *xlog {
	return ins
}

func WithCaller() log.Logger {
	return With("caller", log.DefaultCaller)
}

func With(kv ...interface{}) log.Logger {
	return log.With(
		GetLogger(),
		kv...,
	)
}

func (l *xlog) Log(level log.Level, keyvals ...interface{}) error {
	tag := "[Logx]"
	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		l.log.LoggerX(nil, "WARNING", tag, fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	p1, p2 := make([]string, 0), make([]interface{}, 0)
	for i := 0; i < len(keyvals); i += 2 {
		k1, k2 := keyvals[i], keyvals[i+1]
		k1Val, _ := toString(k1)
		_, k2Label := toString(k2)
		p1 = append(p1, fmt.Sprintf("%s: %s", k1Val, k2Label))
		p2 = append(p2, k2)
	}
	args := fmt.Sprintf("_logx_info %s", strings.Join(p1, ", "))
	switch level {
	case log.LevelDebug:
		l.log.LoggerX(nil, "DEBUG", tag, args, p2...)
	case log.LevelInfo:
		l.log.LoggerX(nil, "INFO", tag, args, p2...)
	case log.LevelWarn:
		l.log.LoggerX(nil, "WARNING", tag, args, p2...)
	case log.LevelError:
		l.log.LoggerX(nil, "ERROR", tag, args, p2...)
	case log.LevelFatal:
		l.log.LoggerX(nil, "FATAL", tag, args, p2...)
	}
	return nil
}

func toString(v interface{}) (string, string) {
	var key string
	label := "%+v"
	if v == nil {
		return key, label
	}
	switch v := v.(type) {
	case float64:
		key = strconv.FormatFloat(v, 'f', -1, 64)
		label = "%d"
	case float32:
		key = strconv.FormatFloat(float64(v), 'f', -1, 32)
		label = "%d"
	case int:
		key = strconv.Itoa(v)
		label = "%d"
	case uint:
		key = strconv.FormatUint(uint64(v), 10)
		label = "%d"
	case int8:
		key = strconv.Itoa(int(v))
		label = "%d"
	case uint8:
		key = strconv.FormatUint(uint64(v), 10)
		label = "%d"
	case int16:
		key = strconv.Itoa(int(v))
		label = "%d"
	case uint16:
		key = strconv.FormatUint(uint64(v), 10)
		label = "%d"
	case int32:
		key = strconv.Itoa(int(v))
		label = "%d"
	case uint32:
		key = strconv.FormatUint(uint64(v), 10)
		label = "%d"
	case int64:
		key = strconv.FormatInt(v, 10)
		label = "%d"
	case uint64:
		key = strconv.FormatUint(v, 10)
		label = "%d"
	case string:
		key = v
		label = "%s"
	case bool:
		key = strconv.FormatBool(v)
	case []byte:
		key = string(v)
		label = "%s"
	case fmt.Stringer:
		key = v.String()
		label = "%s"
	default:
		newValue, _ := json.Marshal(v)
		key = string(newValue)
	}
	return key, label
}
