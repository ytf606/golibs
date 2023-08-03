package osx

import (
	"os"
	"syscall"
)

var defaultRLimit uint64 = 1024000

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
