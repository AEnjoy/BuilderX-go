package builder

import (
	"github.com/aenjoy/BuilderX-go/utils/debugTools"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Task struct {
	TaskID    string
	TaskName  string
	Config    BuildConfig
	CreatTime time.Time
}

func (t *Task) Build() {
	logrus.Info("开始构建任务:", t.TaskName, " 任务ID:", t.TaskID)
	if t.Config.ParseConfig() {
		logrus.Info("初始化编译配置成功。")
		debugTools.PrintlnOnlyInDebugMode("编译命令:" + strings.Join(t.Config.command2, " "))
		if t.Config.Build() {
			logrus.Info("编译成功。输出:", t.Config.OutputFile, " Package:", t.Config.outName)
		} else {
			logrus.Info("编译失败。")
		}
	} else {
		logrus.Info("初始化编译配置失败。构建失败。")
	}
}
