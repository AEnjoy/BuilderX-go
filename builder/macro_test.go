package builder

import (
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"regexp"
	"strings"
	"testing"
)

func TestParserMacro(t *testing.T) {
	str := "${command `git version`} and ${env `envName`} "
	if !isMacro(str) {
		return
	}
	matches := re.FindAllStringSubmatch(str, -1)
	for i, match := range matches {
		t.Logf("Match %d: %s\n", i, match[1])
		command := strings.Split(match[1], " ")
		instruct := command[0]
		args := strings.Replace(strings.Join(command[1:], " "), "`", "", -1)
		//args = strings.Replace(args, "`", "", 2)
		command2 := strings.Split(args, " ")
		t.Logf("instruct:%s, args:%s\n", instruct, args)
		if instruct == "command" {
			t.Logf("command\n")
			t.Logf("%s %s\n", command2[0], command2[1:])
			t.Logf("%s\n", ioTools.GetOutputDirectly(command2[0], command2[1:]...))
		}
		if instruct == "env" {
			t.Logf("env\n")
			t.Logf("%s\n", args)
		}
	}
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
