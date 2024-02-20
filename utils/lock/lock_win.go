//go:build windows

package lock

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func Lock() {
	var err error
	_, err = os.Stat(lockFile)
	if err != nil && os.IsNotExist(err) {
		locks, err = os.Create(lockFile)
		if err != nil {
			logrus.Errorln("BuildGoX:创建Lock失败", err)
		}
		locks.WriteString(strconv.Itoa(os.Getpid()))
		//lock.Close()
	} else {
		logrus.Warning("BuildGoXLock存在，尝试取得lock:")
		e := os.Remove(lockFile)
		if e != nil {
			logrus.Errorln("BuildGoX:获取Lock失败")
			Exit(1, "BuildGoX重复执行，当前进程ID为:", os.Getpid(), ", 退出.")
		}
		logrus.Info("IOM:获取Lock成功")
		Lock()
		return
	}
}
func UnLock() {
	locks.Close()
	os.Remove(lockFile)
}
