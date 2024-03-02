package builder

import (
	"encoding/base64"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
	"time"
)

var re, _ = regexp.Compile(`\${([^}]+)}`)

func isMacro(str string) bool {
	matches := re.FindAllStringSubmatch(str, -1)
	if matches == nil {
		return false
	} else {
		return true
	}
}

func ParserMacro(str string) (retVal string) {
	retVal = str
	//该字段是否存在一条或多条宏指令
	if !isMacro(str) {
		return
	}
	matches := re.FindAllStringSubmatch(str, -1)
	for i, match := range matches {
		// match[1] : "指令,`arg1`,`arg2`,..."
		command := strings.Split(match[1], global.MacroSplit)
		instruct := command[0]
		if len(command) == 2 {
			args := strings.Replace(strings.Join(command[1:], " "), "`", "", -1) //所有的参数 带空格 如"exec a b c" 或"fileName"
			commandArgs := strings.Split(args, " ")                              //所有按" "空格分隔的参数 如 exec,a,b,c
			var value string
			switch instruct {
			case "command":
				value = string(ioTools.GetOutputDirectly(commandArgs[0], commandArgs[1:]...))
				retVal = strings.Replace(retVal, match[1], value, 1)
			case "env":
				value = os.Getenv(args)
				retVal = strings.Replace(retVal, match[1], value, 1)
			case "file":
				value = ioTools.FileReadAll(args)
				retVal = strings.Replace(retVal, match[1], value, 1)
			case "date":
				value = time.Now().Format(args)
				retVal = strings.Replace(retVal, match[1], value, 1)
			case "base64":
				value = base64.StdEncoding.EncodeToString([]byte(args))
				retVal = strings.Replace(retVal, match[1], value, 1)
			}
		} else if len(command) == 3 {
			//instruct is json or yaml
			switch instruct {
			case "json":
				//todo
			case "yaml":
				//todo
			}
		} else {
			logrus.Warningf("command[%d] format error: %s. \n", i, match[1])
			logrus.Infoln("ignore this macro.")
		}
	}
	//在最后,要去掉所有的 "${"和"}"
	retVal = strings.Replace(retVal, "${", "", -1)
	retVal = strings.Replace(retVal, "}", "", -1)
	return
}

func HaveMacroBeforeCompile(str string) bool {
	regex := regexp.MustCompile("\\${!([^}]+)}")
	matches := regex.FindAllString(str, -1)
	return len(matches) != 0
}
