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
			logrus.Errorln("BuilderX:创建Lock失败", err)
		}
		locks.WriteString(strconv.Itoa(os.Getpid()))
		//lock.Close()
	} else {
		logrus.Warningln("BuilderXLock存在，尝试取得lock: 你可以使用--not-running-check忽略重复执行检测")
		e := os.Remove(lockFile)
		if e != nil {
			logrus.Errorln("BuilderX:获取Lock失败")
			Exit(1, "BuilderX重复执行，当前进程ID为:", os.Getpid(), ", 退出.")
		}
		logrus.Infoln("BuilderX:获取Lock成功")
		Lock()
		return
	}
}
func UnLock() {
	locks.Close()
	os.Remove(lockFile)
}
