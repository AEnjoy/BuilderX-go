package builder

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var OutFileNameFmt string // 输出文件名格式: default/a default:使用go默认输出格式(packageName[.exe]),a:{package-name}-{os}-{arch}[.exe]

// VarFlag
// 传递给main.go 中属性参数的值
type VarFlag struct {
	Key   string
	Value string
}

const (
	Set = iota
	Count
	Atomic
)

type BuildArch struct {
	GOOS   string
	GOARCH string
}

func isGoFile(fileName string) bool {
	var retval1 = func() bool { return strings.HasSuffix(fileName, ".go") }
	var retval2 = func() bool {
		var retval string
		file, err := os.Open("go.mod")
		if err == nil {
			defer file.Close()
			r := bufio.NewReader(file)
			l, _, _ := r.ReadLine()
			retval = strings.Replace(string(l), "package", "", 1)
		}
		if retval != "" {
			return true
		}
		return false
	}
	return retval1() && retval2()
}

func getNowArch() BuildArch {
	return BuildArch{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH}
}

// getGoPackageName 获取当前项目go.mod文件中的module值(完整包名)
func getGoPackageName() (string, error) {
	var retval string
	file, err := os.Open("go.mod")
	if err == nil {
		defer file.Close()
		r := bufio.NewReader(file)
		l, _, _ := r.ReadLine()
		retval = strings.Replace(string(l), "module", "", 1)
		retval = strings.TrimSpace(retval)
	} else {
		logrus.Errorln("go.mod file not found in current directory.")
		return "", err
	}
	return retval, nil
}

// getGoPackageName2 返回相对路径下的包名
func getGoPackageName2() (string, error) {
	var retval string
	file, err := os.Open("go.mod")
	if err == nil {
		defer file.Close()
		r := bufio.NewReader(file)
		l, _, _ := r.ReadLine()
		retval = strings.Replace(string(l), "module", "", 1)
		retval = strings.TrimSpace(retval)
	} else {
		logrus.Errorln("go.mod file not found in current directory.")
		return "", err
	}
	ret := strings.Split(retval, "/")
	if len(ret)-1 != 0 {
		retval = ret[len(ret)-1]
		return retval, nil
	}
	return "", nil
}

func EnableCGO() {
	defaultConfigY.BaseConfig.Cgo = true
}
func init() {
	os.Mkdir("project", 0750)
	var d = yamlConfig{}.BaseConfig
	d.RemoteConfig.LocalStoreTemp = "./project"
	defaultConfigY.BaseConfig = d
	var d2 = jsonConfig{}.BaseConfig
	d2.RemoteConfig.LocalStoreTemp = "./project"
	defaultConfigJ.BaseConfig = d2
}
