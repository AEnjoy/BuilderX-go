package macro

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	tools "github.com/aenjoy/BuilderX-go/utils/jsonYamlTools"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var re, _ = regexp.Compile(`\${([^}]+)}`)

func IsMacro(str string) bool {
	matches := re.FindAllStringSubmatch(str, -1)
	if matches == nil {
		return false
	} else {
		return true
	}
}

func (m *Macro) ParserMacro(str string) (retVal string) {
	retVal = str
	//该字段是否存在一条或多条宏指令
	if !IsMacro(str) {
		return
	}
	if m.IsDefineMacro(str) {
		retVal = m.ParserDefineMacro(str)
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
			switch instruct {
			case "command":
				value := string(ioTools.GetOutputDirectly(commandArgs[0], commandArgs[1:]...))
				retVal = strings.Replace(retVal, match[1], value, 1)
			case "env":
				value := os.Getenv(args)
				retVal = strings.Replace(retVal, match[1], value, 1)
			case "file":
				value := ioTools.FileReadAll(args)
				retVal = strings.Replace(retVal, match[1], value, 1)
			case "date":
				value := time.Now().Format(args)
				retVal = strings.Replace(retVal, match[1], value, 1)
			case "base64":
				value := base64.StdEncoding.EncodeToString([]byte(args))
				retVal = strings.Replace(retVal, match[1], value, 1)
			}
		} else if len(command) == 3 {
			//instruct is json or yaml
			file := strings.Replace(command[1], "`", "", -1)
			args := strings.Replace(command[2], "`", "", -1)
			f, err := os.ReadFile(file)
			if err != nil {
				logrus.Errorln("File load error:", err, " Ignore this macro.")
				continue
			}
			var data map[string]interface{}
			switch instruct {
			case "json":
				err = json.Unmarshal(f, &data)
				if err != nil {
					logrus.Errorln("jsonFile load error:", err, " Ignore this macro.")
					continue
				}
			case "yaml":
				err = yaml.Unmarshal(f, &data)
				if err != nil {
					logrus.Errorln("YamlFile load error:", err, " Ignore this macro.")
					continue
				}
			}
			value, flag := tools.GetFieldValue(strings.Split(args, "."), data)
			var v string
			if flag {
				switch t := value.(type) {
				case int:
					v = strconv.Itoa(t)
				case string:
					v = t
				default:
					logrus.Errorln("Unsupported type:", t, " Ignore this macro.")
					continue
				}
			} else {
				logrus.Errorln("Can't find field:", args, " Ignore this macro.")
				continue
			}
			retVal = strings.Replace(retVal, match[1], v, 1)
		} else {
			logrus.Warningf("command[%d] format error: %s. \n", i, match[1])
			logrus.Infoln("ignore this macro.")
		}
	}
	//在最后,要去掉所有的 "${"和"}"
	//这里其实有一个bug，如果这条宏不存在,也会被替换掉
	retVal = strings.Replace(retVal, "${", "", -1)
	retVal = strings.Replace(retVal, "}", "", -1)
	retVal = strings.Replace(retVal, "\n", "", -1)
	return
}

func HaveMacroBeforeCompile(str string) bool {
	regex := regexp.MustCompile("\\${!([^}]+)}")
	matches := regex.FindAllString(str, -1)
	return len(matches) != 0
}

type Macro struct {
	defineContext map[string]string
	init          bool
}

func (m *Macro) SetDefineContext(str []string) {
	if !m.init {
		m.defineContext = make(map[string]string)
		m.init = true
	}

	for _, s := range str {
		v := strings.Split(s, "=")
		if len(v) >= 2 {
			m.defineContext[v[0]] = m.ParserMacro(strings.Join(v[1:], "="))
		} else if m.IsDefineMacro(s) {
			m.ParserDefineMacro(s)
		}
	}
}

func (m *Macro) IsDefineMacro(str string) bool {
	var r = regexp.MustCompile("\\${define,(.*?)}")
	matches := r.FindAllStringSubmatch(str, -1)
	var r2 = regexp.MustCompile("\\${using,(.*?)}")
	matches2 := r2.FindAllStringSubmatch(str, -1)
	return len(matches) != 0 || len(matches2) != 0
}

func (m *Macro) ParserDefineMacro(str string) (retVal string) {
	retVal = str
	if !m.IsDefineMacro(str) {
		return
	}
	if !m.init {
		m.defineContext = make(map[string]string)
		m.init = true
	}
	matches := re.FindAllStringSubmatch(str, -1)
	for i, match := range matches {
		command := strings.Split(match[1], global.MacroSplit)
		if len(command) < 2 {
			logrus.Warningf("command[%d] format error: %s. \n", i, match[1])
			logrus.Infoln("ignore this macro.")
			continue
		} else if len(command) > 3 {
			logrus.Warningf("command[%d] format error: %s. \n", i, match[1])
			logrus.Infoln("ignore this macro.")
			continue
		}
		define := strings.Replace(command[1], "`", "", -1)
		if len(command) == 3 && command[0] == "define" {
			//Set define
			value := strings.Replace(command[2], "`", "", -1)
			m.defineContext[define] = value
		} else if len(command) == 2 && command[0] == "using" {
			v, ok := m.defineContext[define]
			if ok {
				retVal = strings.Replace(retVal, match[1], v, -1)
			} else {
				logrus.Errorln("Define macro not found:", command[0])
			}
		}
	}
	//这里其实有一个bug，如果这条宏不存在,也会被替换掉
	retVal = strings.Replace(retVal, "${", "", -1)
	retVal = strings.Replace(retVal, "}", "", -1)
	retVal = strings.Replace(retVal, "\n", "", -1)
	return
}

func (m *Macro) GetDefine(define string) string {
	v, ok := m.defineContext[define]
	if ok {
		return v
	} else {
		return ""
	}
}
