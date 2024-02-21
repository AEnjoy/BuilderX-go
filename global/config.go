package global

import (
	"BuilderX/utils/debugTools"
	"BuilderX/utils/ioTools"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine
var WebPort = "18088"
var StartTime time.Time
var BuildedTask = 0
var RootDir string
var ConfigApiVersion = 1 // BuilderX-go支持的当前配置文件api版本

func init() {
	StartTime = time.Now()
	debugTools.StartTime = StartTime
	_, err := exec.LookPath("go")
	if err != nil {
		logrus.Errorln("未找到golang，请先安装golang")
		return
	}
	GoOSAndGoArchSupported = string(ioTools.GetOutputDirectly("go", "tool", "dist", "list"))
	debugTools.PrintlnOnlyInDebugMode("GoOSAndGoArchSupported:", GoOSAndGoArchSupported)
	RootDir, _ = os.Getwd()
}
