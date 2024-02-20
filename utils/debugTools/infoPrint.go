package debugTools

import (
	"github.com/sirupsen/logrus"
)

var DebugFlag bool

func PrintLogs(a ...interface{}) {
	if DebugFlag {
		logrus.Debugln(a)
	} else {
		logrus.Infoln(a)
	}
}

func Println(a ...interface{}) {
	if DebugFlag {
		println("Debug:", a)
	} else {
		println(a)
	}
}
func PrintlnOnlyInDebugMode(a ...interface{}) bool {
	if DebugFlag {
		Println(a)
	}
	return DebugFlag
}
func PrintLogsOnlyInDebugMode(a ...interface{}) bool {
	if DebugFlag {
		PrintLogs(a)
	}
	return DebugFlag // 返回debug flag
}
