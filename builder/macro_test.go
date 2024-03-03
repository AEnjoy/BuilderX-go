package builder

import (
	"encoding/base64"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	tools "github.com/aenjoy/BuilderX-go/utils/jsonYamlTools"
	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParserMacro(t *testing.T) {
	str := "NoWeb,CGO_ENABLED=${command,`go env CGO_ENABLED`}"
	if !isMacro(str) {
		return
	}
	matches := re.FindAllStringSubmatch(str, -1)
	for i, match := range matches {
		t.Logf("Match %d: %s\n", i, match[1])
		command := strings.Split(match[1], global.MacroSplit)
		instruct := command[0]
		args := strings.Replace(strings.Join(command[1:], " "), "`", "", -1)
		commandArgs := strings.Split(args, " ")
		//t.Logf("instruct:%s, args:%s\n", instruct, args)
		switch instruct {
		case "command":
			value := string(ioTools.GetOutputDirectly(commandArgs[0], commandArgs[1:]...))
			str = strings.Replace(str, match[1], value, 1)
		case "env":
			value := os.Getenv(args)
			str = strings.Replace(str, match[1], value, 1)
		case "file":
			value := ioTools.FileReadAll(args)
			str = strings.Replace(str, match[1], value, 1)
		case "date":
			value := time.Now().Format(args)
			str = strings.Replace(str, match[1], value, 1)
		case "base64":
			value := base64.StdEncoding.EncodeToString([]byte(args))
			str = strings.Replace(str, match[1], value, 1)
		}
	}
	//在最后,要去掉所有的 "${"和"}"
	str = strings.Replace(str, "${", "", -1)
	str = strings.Replace(str, "}", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	t.Logf("%s\n", str)
}

func TestReg(t *testing.T) {
	input := "这是一个包含\"反单引号\"\\的字符串。(`123` 456)"

	// 创建一个正则表达式对象，匹配反单引号\\中的内容
	re2, _ := regexp.Compile(`"[\s\S]*?"`)

	// 在字符串中查找匹配项
	matched := re2.MatchString(input)
	if matched {

		// 提取匹配项中的内容
		matches := re2.FindAllStringSubmatch(input, -1)
		for _, match := range matches {
			t.Logf("Match: %s\n", match[0])
		}
	}
}

func TestParserMacroF(t *testing.T) {
	os.Chdir("..")
	config2 := loadConfigYaml("builderX.yaml")
	t.Logf("%s\n len:%d\n", config2.BaseConfig.VarFlags, len(config2.BaseConfig.VarFlags))
	var varFlags []VarFlag
	for i, v := range config2.BaseConfig.VarFlags {
		var varFlag VarFlag
		t.Logf("config2.BaseConfig.VarFlags[%d], v is :%s\n", i, v)
		a := strings.Split(v, "=")
		//t.Logf("config2.BaseConfig.VarFlags[%d], a is :%s\n", i, a)
		if len(a) >= 2 {
			t.Logf("Found varFlag: %s%s%s", a[0], "=", a[1])
			varFlag.Key = a[0]
			str := a[1]
			if !isMacro(str) {
				t.Logf("%s is not macro.\n\n", str)
				varFlag.Value = str
				varFlags = append(varFlags, varFlag)
				continue
			}
			matches := re.FindAllStringSubmatch(str, -1)
			for i, match := range matches {
				t.Logf("Match %d: %s\n", i, match[1])
				command := strings.Split(match[1], global.MacroSplit)
				instruct := command[0]
				args := strings.Replace(strings.Join(command[1:], " "), "`", "", -1)
				//args = strings.Replace(args, "`", "", 2)
				command2 := strings.Split(args, " ")
				t.Logf("instruct:%s, args:%s\n", instruct, args)
				if instruct == "command" {
					t.Logf("command\n")
					t.Logf("%s %s\n", command2[0], command2[1:])
					t.Logf("%s\n", ioTools.GetOutputDirectly(command2[0], command2[1:]...))
					varFlag.Value = string(ioTools.GetOutputDirectly(command2[0], command2[1:]...))
				}
				if instruct == "env" {
					t.Logf("env\n")
					t.Logf("%s\n", os.Getenv(args))
					varFlag.Value = os.Getenv(args)
				}
				if instruct == "file" {
					t.Logf("file\n")
					t.Logf("%s\n", ioTools.FileReadAll(args))
					varFlag.Value = ioTools.FileReadAll(args)
				}
				if instruct == "date" {
					t.Logf("date\n")
					t.Logf("%s\n", time.Now().Format(args))
					varFlag.Value = time.Now().Format(args)
				}
				//time.Now().Format("2006-01-02--15:04:05")
			}
			varFlags = append(varFlags, varFlag)
		} else {
			continue
		}
	}
	for i, flag := range varFlags {
		t.Logf("varFlags[%d],%s=%s", i, flag.Key, flag.Value)
	}
}

func TestHaveMacroBeforeCompileRe(t *testing.T) {
	input := "123 ${!abc000 cc} ${!eee000 cc}"
	input2 := "abc000"
	// 创建一个正则表达式对象，匹配 ${! }中的内容
	regex := regexp.MustCompile("\\${!([^}]+)}")

	// 使用FindAllString方法查找所有匹配项
	matches := regex.FindAllString(input, -1)
	matches2 := regex.FindAllString(input2, -1)
	t.Logf("matches:%v\n", matches)
	t.Logf("matches2:%v\n", matches2)
}

func TestGetFieldY(t *testing.T) {
	str := "${yaml,`../logs/yaml_test.yaml`,`test.name`}"
	/*
		test:
		  name: "John Doe"
		  dob: "1990-01-01"
	*/
	matches := re.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		// match[1] : "指令,`arg1`,`arg2`,..."
		command := strings.Split(match[1], global.MacroSplit)
		//instruct := command[0]
		file := strings.Replace(command[1], "`", "", -1)
		args := strings.Replace(command[2], "`", "", -1)
		t.Logf("file=%s,args=%s\n", file, args)
		var data map[string]interface{}
		f, err := os.ReadFile(file)
		if err != nil {
			t.Error("YamlFile load error:", err)
			continue
		}
		err = yaml.Unmarshal(f, &data)
		if err != nil {
			t.Error("YamlFile load error:", err, " Ignore this macro.")
			continue
		}
		value, flag := tools.GetFieldValue(strings.Split(args, "."), data)
		var v string
		if flag {
			switch t1 := value.(type) {
			case int:
				v = strconv.Itoa(t1)
			case string:
				v = t1
			default:
				t.Error("Unsupported type:", t1, " Ignore this macro.")
				continue
			}
		}
		t.Log("value:", v)
	}
}

func TestGetFieldJ(t *testing.T) {
	str := "${json,`../logs/json_test.json`,`test.name`}"
	/*
		{
		  "test": {
		    "name": "John Doe",
		    "dob": "1990-01-01"
		  }
		}
	*/
	matches := re.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		// match[1] : "指令,`arg1`,`arg2`,..."
		command := strings.Split(match[1], global.MacroSplit)
		//instruct := command[0]
		file := strings.Replace(command[1], "`", "", -1)
		args := strings.Replace(command[2], "`", "", -1)
		t.Logf("file=%s,args=%s\n", file, args)
		f, err := os.ReadFile(file)
		if err != nil {
			t.Error("jsonFile load error:", err)
			continue
		}
		var data map[string]interface{}
		err = sonic.Unmarshal(f, &data)
		value, flag := tools.GetFieldValue(strings.Split(args, "."), data)
		var v string
		if flag {
			switch t1 := value.(type) {
			case int:
				v = strconv.Itoa(t1)
			case string:
				v = t1
			default:
				t.Error("Unsupported type:", t1, " Ignore this macro.")
				continue
			}
		}
		t.Log("value:", v)
	}
}
