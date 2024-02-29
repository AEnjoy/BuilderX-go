package builder

import (
	"encoding/base64"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
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
	if !isMacro(str) {
		return
	}
	matches := re.FindAllStringSubmatch(str, -1)
	for i, match := range matches {
		command := strings.Split(match[1], " ")
		if len(command) >= 2 {
			var args string
			var command2 []string
			//var file string
			//var config string
			instruct := command[0]
			if instruct != "json" && instruct != "yaml" {
				args = strings.Replace(strings.Join(command[1:], " "), "`", "", -1)
				command2 = strings.Split(args, " ")
			} else {
				// todo yaml,json支持暂未实现
				/*re2, _ := regexp.Compile(`"[\s\S]*?"`)
				matched := re2.MatchString(strings.Join(command[1:], " "))
				if matched {
					// 提取匹配项中的内容
					matches2 := re2.FindAllStringSubmatch(strings.Join(command[1:], " "), -1)
					for j, match2 := range matches2 {
						if j == 0 {
							file = match2[0]
						}
						if j == 1 {
							config = match2[0]
						}
					}
				}*/
			}
			//args = strings.Replace(args, "`", "", 2)
			switch instruct {
			case "command":
				retVal = strings.Replace(retVal, match[1], string(ioTools.GetOutputDirectly(command2[1], command2[2:]...)), -1)
			case "env":
				retVal = strings.Replace(retVal, match[1], os.Getenv(args), -1)
			case "file":
				retVal = strings.Replace(retVal, match[1], ioTools.FileReadAll(args), -1)
			case "json":
				// todo yaml,json支持暂未实现
				continue
				/*
					f, err := os.ReadFile(file)
					if err != nil {
						logrus.Warningln("file not exist:", args)
						continue
					}
					var payload map[string]interface{}
					err = sonic.Unmarshal(f, &payload)
					fromString, err := sonic.GetFromString(config, &payload)
					if err != nil {
						logrus.Warningln("Error parsing json:", err)
						continue
					} else {
						s, err := fromString.Get(config).String()
						if err != nil {
							logrus.Warningln("Error parsing json:", err)
							continue
						}
						retVal = strings.Replace(retVal, match[1], s, -1)
					}*/
			case "yaml":
				// todo yaml,json支持暂未实现
				continue
			case "base64":
				decodedStr, err := base64.StdEncoding.DecodeString(args)
				if err != nil {
					logrus.Warningln("Error decoding base64:", err)
					continue
				}
				retVal = strings.Replace(retVal, match[1], string(decodedStr), -1)
			}
		} else {
			logrus.Warningf("command[%d] format error: %s. \n", i, match[1])
			continue
		}
	}
	return
}
