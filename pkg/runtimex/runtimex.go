package runtimex

import (
	"path"
	"runtime"
	"strings"
)

//获取调用处的包名_文件名_函数名
func PF() string {
	pc, file, _, _ := runtime.Caller(1)
	pcName := strings.Split(path.Base(runtime.FuncForPC(pc).Name()), ".")
	fileName, fileSuffix := path.Base(file), path.Ext(file)
	file = fileName[:len(fileName)-len(fileSuffix)]
	return "[" + pcName[0] + "_" + file + "_" + pcName[len(pcName)-1] + "]"
}

func SetCPU(num int) {
	if runtime.NumCPU() <= num {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	}
}
