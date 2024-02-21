package debugTools

import (
	"github.com/sirupsen/logrus"
)

var DebugFlag bool

// PrintLogs
// 打印日志，Debug模式下打印Debug级别日志到std-debug，否则打印标准Info日志到stdout
func PrintLogs(a ...interface{}) {
	if DebugFlag {
		logrus.Debugln(a)
	} else {
		logrus.Infoln(a)
	}
}

// Println
// 打印日志，Debug模式下打印日志到stdout(增加"Debug:"字符)，否则正常打印日志到stdout
func Println(a ...interface{}) {
	if DebugFlag {
		println("Debug:", a)
	} else {
		println(a)
	}
}

// PrintlnOnlyInDebugMode
//
//	打印日志，Debug模式下打印日志到stdout(增加"Debug:"字符)，否则不打印
func PrintlnOnlyInDebugMode(a ...interface{}) bool {
	if DebugFlag {
		Println(a)
	}
	return DebugFlag
}

// PrintLogsOnlyInDebugMode
//
//	打印日志，Debug模式下打印Debug日志，否则不打印
func PrintLogsOnlyInDebugMode(a ...interface{}) bool {
	if DebugFlag {
		PrintLogs(a)
	}
	return DebugFlag // 返回debug flag
}
