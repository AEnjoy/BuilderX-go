//go:build !windows
// +build !windows

package lock

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func Lock() {
	var err error
	locks, err = os.Create(lockFile)
	if err != nil {
		logrus.Errorln("BuildGoX:创建Lock失败", err)
	}
	err = syscall.Flock(int(locks.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		Exit(1, "BuildGoX重复执行，当前进程ID为:", os.Getpid(), ",退出.")
	}
}
func UnLock() {
	syscall.Flock(int(locks.Fd()), syscall.LOCK_UN)
	locks.Close()
	os.Remove(lockFile)
}
