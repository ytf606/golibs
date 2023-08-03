package logx

import (
	"fmt"
	"os"
	"strings"

	"github.com/ytf606/golibs/logx/log4go"
	"github.com/ytf606/golibs/logx/logutils"
)

//使用自定义option
func NewLogConfig() *log4go.LogConfig {
	config := new(log4go.LogConfig)
	//存储路径
	config.LogPath = "/tmp/default.log"
	//日志级别
	config.Level = "INFO"
	//日志标签 多日志时使用
	config.Tag = "default"
	//日志格式
	config.Format = "%G %L %S %M"
	//最大行数切割
	config.RotateLines = "0K"
	//最大容量切割
	config.RotateSize = "0M"
	//按日期切割
	config.RotateHourly = true
	//是否启用切割
	config.Rotate = true
	//日志保留时间，day
	config.Retention = "0"
	return config
}

//自定义config Init
func InitLogWithConfig(config *log4go.LogConfig) {
	if config.LogPath == "" {
		fmt.Fprintf(os.Stderr, "InitLoggerConfig: Error: config could not found logpath %s\n", config.LogPath)
		os.Exit(1)
	}
	checkLogConfig(config)
	log4go.LoadLogConfig(config)
	logutils.Inited = true
}

func Close() {
	log4go.Close()
}
func checkLogConfig(config *log4go.LogConfig) {
	if _, ok := logutils.LevelMap[config.Level]; ok {
		if logutils.LevelMap[config.Level] < logutils.SortLevel {
			logutils.SortLevel = logutils.LevelMap[config.Level]
			logutils.Level = config.Level
		}
	}
	paths := strings.Split(config.LogPath, "/")
	if len(paths) > 1 {
		dir := strings.Join(paths[0:len(paths)-1], "/")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not create directory %s, err:%s\n", dir, err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: log directory invalid %s\n", config.LogPath)
		os.Exit(1)
	}
}
