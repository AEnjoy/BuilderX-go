package lock

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Exit(code int, msg ...interface{}) {
	if code != 0 {
		logrus.Errorln("ServerExitWithMessage:", msg)
	} else {
		logrus.Info("ServerExitWithMessage:", msg)
	}
	logrus.Info("ServerExitWithCode:", code)
	UnLock()
	os.Exit(code)
}
func ExitHandle(exitChan chan os.Signal) {
	for {
		select {
		case sig := <-exitChan:
			logrus.Info("收到来自系统的信号：", sig)
			Exit(2, sig.String())
		}
	}

}
