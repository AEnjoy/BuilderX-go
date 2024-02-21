package lock

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Exit(code int, msg ...interface{}) {
	if code != 0 {
		logrus.Errorln("ServerExitWithMessage:", msg)
	} else {
		logrus.Infoln("ServerExitWithMessage:", msg)
	}
	logrus.Infoln("ServerExitWithCode:", code)
	UnLock()
	os.Exit(code)
}
func ExitHandle(exitChan chan os.Signal) {
	for {
		select {
		case sig := <-exitChan:
			logrus.Infoln("收到来自系统的信号：", sig)
			Exit(2, sig.String())
		}
	}

}
