package osx

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

var defaultRLimit uint64 = 1024000

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

// 设置进程最大打开文件描述符个数，防止出现too many open files
func SetRLimit() {
	var rlim syscall.Rlimit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if rlim.Cur < defaultRLimit {
		rlim.Cur = defaultRLimit
		rlim.Max = defaultRLimit
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim)
		syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
}

func GetHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

// 最终方案-全兼容
func getCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
