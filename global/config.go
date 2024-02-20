package global

import (
	"BuilderX/utils/debugTools"
	"BuilderX/utils/iotools"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine
var WebPort = "18088"
var StartTime time.Time
var BuildedTask = 0

func init() {
	StartTime = time.Now()
	debugTools.StartTime = StartTime
	_, err := exec.LookPath("go")
	if err != nil {
		debugTools.PrintLogs("未找到golang，请先安装golang")
		return
	}
	GoOSAndGoArchSupported = string(iotools.GetOutputDirectly("go", "tool", "dist", "list"))
	debugTools.PrintLogs("GoOSAndGoArchSupported:", GoOSAndGoArchSupported)
}
