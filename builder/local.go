package builder

import (
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/debugTools"
	"github.com/aenjoy/BuilderX-go/utils/hashtool"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func UsingLocal(path string) Task {
	var task Task
	config := defaultConfigY
	logrus.Infoln("Using local path:", path, " parse ...")
	if path == "" || path == "." {
		task.Config.InputFile = "."
		// 使用当前路径
	} else {
		_, err := os.Stat(path)
		if err != nil {
			return Task{}
		}
	}
	//判断当前路径下有无配置文件,有则加载该配置
	_, err := os.Stat("builderX.yaml")
	if err == nil {
		logrus.Infoln("Using builderX.yaml config.")
		config = loadConfigYaml("builderX.yaml")
	}
	task.Config = yamlConfig2BuildConfig(config)
	task.CreatTime = time.Now()
	task.TaskName = "localBuild"
	task.TaskID = hashtool.MD5(task.CreatTime.Format("2006-01-02-15:04:05") + strconv.Itoa(global.BuildedTask) + task.TaskName)
	global.BuildedTask++
	debugTools.PrintlnOnlyInDebugMode("debug:", task.Config.command)
	//task.Config.Targets = make([]BuildArch, 0)
	task.Config.Targets = append(task.Config.Targets, getNowArch())
	if !task.Config.ParseConfig() {
		logrus.Error("parse config failed.")
	}
	return task
}
