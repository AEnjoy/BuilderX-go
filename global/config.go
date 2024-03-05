package global

import (
	"github.com/aenjoy/BuilderX-go/utils/debugTools"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine
var WebPort = "18088"
var StartTime time.Time
var BuiltTask = 0
var RootDir string

const ConfigApiVersion = 3 // BuilderX-go支持的当前配置文件api版本
var GoExe string

const MacroSplit = "," //宏指令的分隔符

func init() {
	StartTime = time.Now()
	debugTools.StartTime = StartTime
	_, err := exec.LookPath("go")
	if err != nil {
		logrus.Warningln("未找到系统内安装的golang，请先安装golang或设置--go-exe参数")
	}
	GoOSAndGoArchSupported = string(ioTools.GetOutputDirectly("go", "tool", "dist", "list"))
	debugTools.PrintlnOnlyInDebugMode("GoOSAndGoArchSupported:", GoOSAndGoArchSupported)
	RootDir, _ = os.Getwd()
}
