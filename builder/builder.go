package builder

import (
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/debugTools"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
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

func getGoPackageName() (string, error) {
	var retval string
	file, err := os.Open("go.mod")
	if err == nil {
		defer file.Close()
		r := bufio.NewReader(file)
		l, _, _ := r.ReadLine()
		retval = strings.Replace(string(l), "module", "", 1)
	} else {
		logrus.Errorln("go.mod file not found in current directory.")
		return "", err
	}
	return retval, nil
}

// BuildConfig
//
//	构建配置 本地项目
type BuildConfig struct {
	InputFile  string    //单个文件名或路径
	OutputFile string    //输出文件名或路径
	outName    string    //命令提示中输出的文件名(程序包名)
	Ldflags    []string  //传递给链接参数的值
	VarFlags   []VarFlag //传递给main.go 中属性参数的值
	V          bool      // -v 打印编译的包名和文件名
	Cgo        bool      //enable cgo default false
	//不常用参数
	ForceBuildPackage bool     // -a 强制重新生成已更新的包。
	BuildProcess      int      //-p n 编译线程数
	Race              bool     // -race 启用数据竞争检测（only in 64bit ）
	Msan              bool     // -msan 启用内存分配扫描 linux/amd64, linux/arm64, freebsd/amd64
	Cover             bool     // -cover 启用覆盖率分析
	CoverMode         int      //设置覆盖率分析的模式 默认set
	Work              bool     //打印临时工作目录的名称，退出时不要将其删除。
	AsmFlags          []string //传递给汇编器的值 [pattern=]arg
	BuildMode         string   //编译模式
	BuildVcs          string   //是否用版本控制信息标记二进制文件（auto，true，false）
	Compiler          string   //gccgo or gc
	Gccgoflags        []string //arguments to pass on each gccgo compiler/linker invocation.
	Gcflags           []string //arguments to pass on each go tool compile invocation.
	Linkshared        bool     //创建共享库 BuildMode=shared
	Mod               string   //readonly, vendor, or mod.
	Modcacherw        bool     //使模块缓存中新创建的目录保持读写状态，而不是使其只读。
	Modfile           string   //fileName 使用备用的go.mod
	Overlay           string   //fileName 读取为构建操作提供overlay的JSON配置文件。
	Pgo               string   //fileName 指定一个profile-guided optimization文件
	Pkgdir            string   //dir 使用额外的包目录
	Tags              string   //指定额外的tags
	//
	Targets  []BuildArch //构建的目标架构
	command  string
	command2 []string
	status   bool //true:就绪 false:未就绪
}

func (c *BuildConfig) ParseConfig() bool {
	if c.status { //防止重复处理
		return true
	}
	_, err := exec.LookPath("go")
	c.command2 = append(c.command2, "build")
	if err != nil {
		logrus.Errorln("未找到go，请先安装golang", err)
		c.command2 = make([]string, 0)
		return false
	}
	if c.Cgo {
		_, err = exec.LookPath("gcc")
		if err != nil {
			logrus.Errorln("未找到gcc，请先安装gcc,或关闭cgo选项。", err)
			c.command2 = make([]string, 0)
			return false
		}
	}
	if c.InputFile != "" && c.InputFile != "." {
		//文件判断
		f, err := os.Stat(c.InputFile)
		if err != nil { //不是目录或文件 (但还有可能是包名,待验证)
			logrus.Errorln("未找到输入包（文件），请检查输入的项目文件（包）路径是否正确.", err)
			t, _ := os.Getwd()
			logrus.Infoln("E:当前输入的包信息:", c.InputFile, "当前路径：", t)
			c.command2 = make([]string, 0)
			return false
			// todo 如果是包名,还需判断
		}
		//目录或文件判断
		if f.IsDir() {
			os.Chdir(c.InputFile)
			c.outName, err = getGoPackageName()
			if err != nil {
				os.Chdir(global.RootDir)
				return false
			}
		} else {
			//go.mod or main.go
			if ioTools.IsStrAInStrB("go.mod", c.InputFile) {
				dir := ioTools.GetFileDir("go.mod")
				err = os.Chdir(dir)
				if err != nil {
					logrus.Errorln("未找到输入包（文件），请检查输入的项目文件（包）路径是否正确.", err)
					return false
				}
				c.outName, err = getGoPackageName()
				if err != nil {
					logrus.Errorln("不是有效的go项目，请检查输入的项目文件（包）路径是否正确.")
					os.Chdir(global.RootDir)
					return false
				}
				c.InputFile = "." //
			} else if !isGoFile(c.InputFile) {
				logrus.Error("输入文件格式错误，请检查输入的项目文件（包）路径是否正确.")
				return false
			}
		}
	} else {
		c.outName, err = getGoPackageName()
		if err != nil {
			return false
		}
	}
	if c.OutputFile != "" {
		logrus.Infoln("输出文件路径为:", c.OutputFile)
		c.command += "-o " + `"` + c.OutputFile + `" `
		c.command2 = append(c.command2, "-o", c.OutputFile)
	}
	if c.ForceBuildPackage {
		c.command += "-a "
		c.command2 = append(c.command2, "-a")
	}
	if c.Cgo {
		err = os.Setenv("CGO_ENABLED", "1")
		if err != nil {
			logrus.Errorln("设置CGO_ENABLED失败")
			c.command2 = make([]string, 0)
			return false
		}
	}
	if c.BuildProcess != 0 {
		c.command += "-p " + strconv.Itoa(c.BuildProcess) + " "
		c.command2 = append(c.command2, "-p", strconv.Itoa(c.BuildProcess))
	}
	if c.V {
		c.command += "-v "
		c.command2 = append(c.command2, "-v")
	}
	if c.Race {
		c.command += "-race "
		c.command2 = append(c.command2, "-race")
	}
	if c.Msan {
		c.command += "-msan "
		c.command2 = append(c.command2, "-msan")
	}
	if c.Cover {
		c.command += "-cover "
		c.command2 = append(c.command2, "-cover")
	}
	if c.CoverMode == Count {
		c.command += "-covermode count "
		c.command2 = append(c.command2, "-covermode", "count")
	} else if c.CoverMode == Atomic {
		c.command += "-covermode atomic "
		c.command2 = append(c.command2, "-covermode", "atomic")
	}
	if c.Work {
		c.command += "-work "
		c.command2 = append(c.command2, "-work")
	}
	if c.Modcacherw {
		c.command += "-modcacherw "
		c.command2 = append(c.command2, "-modcacherw")
	}
	if c.BuildMode != "" {
		if ioTools.IsStrAInStrB(c.BuildMode, global.BuildModeSupported) {
			c.command += "-buildmode " + c.BuildMode + " "
			c.command2 = append(c.command2, "-buildmode", c.BuildMode)
		} else {
			logrus.Errorln("未支持的构建模式，请检查构建模式是否正确. Return")
			c.command2 = make([]string, 0)
			return false
		}
	}
	if c.BuildVcs != "" {
		if ioTools.IsStrAInStrB(c.BuildVcs, global.BuildVcsSupported) {
			c.command += "-buildvcs " + c.BuildVcs + " "
			c.command2 = append(c.command2, "-buildvcs", c.BuildVcs)
		} else {
			logrus.Warningln("未支持的buildvcs选项，请检查构建模式是否正确。默认auto")
		}
	}
	if c.Compiler != "" {
		if ioTools.IsStrAInStrB(c.Compiler, global.CompileSupported) {
			c.command += "-compiler " + c.Compiler + " "
			c.command2 = append(c.command2, "-compiler", c.Compiler)
		} else {
			logrus.Warningln("未支持的编译器，请检查构建模式是否正确。")
		}
	}
	if c.Linkshared {
		c.command += "-linkshared -buildmode=shared "
		c.command2 = append(c.command2, "-linkshared", "-buildmode=shared")
	}
	if c.Mod != "" {
		if ioTools.IsStrAInStrB(c.Mod, global.ModSupported) {
			c.command += "-mod " + c.Mod + " "
			c.command2 = append(c.command2, "-mod", c.Mod)
		} else {
			logrus.Warningln("未支持的Mod.")
		}
	}
	if c.Modfile != "" {
		_, err = os.Stat(c.Modfile)
		if err != nil {
			logrus.Warningln("未找到输入modfile，请检查输入文件路径是否正确")
		} else {
			c.command += "-modfile " + c.Modfile + " "
			c.command2 = append(c.command2, "-modfile", c.Modfile)
		}
	}
	if c.Overlay != "" {
		_, err = os.Stat(c.Overlay)
		if err != nil {
			logrus.Warningln("未找到输入overlay file，请检查输入文件路径是否正确")
		} else {
			c.command += "-overlay " + c.Overlay + " "
			c.command2 = append(c.command2, "-overlay", c.Overlay)
		}
	}
	if c.Pgo != "" {
		_, err = os.Stat(c.Pgo)
		if err != nil {
			logrus.Warningln("未找到输入pgo file，请检查输入文件路径是否正确")
		} else {
			c.command += "-pgo " + c.Pgo + " "
			c.command2 = append(c.command2, "-pgo", c.Pgo)
		}
	}
	if c.Pkgdir != "" {
		_, err = os.Stat(c.Pkgdir)
		if err != nil {
			logrus.Warningln("不存在Pkg dir，请检查输入路径是否正确")
		} else {
			c.command += "-pkgdir " + c.Pkgdir + " "
			c.command2 = append(c.command2, "-pkgdir", c.Pkgdir)
		}
	}
	if c.Tags != "" {
		c.command += "-tags " + c.Tags + " "
		c.command2 = append(c.command2, "-tags", c.Tags)
	}
	var command = ""
	if len(c.Ldflags) != 0 && c.VarFlags != nil {
		c.command += `-ldflags "`
		c.command2 = append(c.command2, "-ldflags")
		for _, i := range c.Ldflags {
			c.command += i + " "
			command += i + " "
		}
		for _, i := range c.VarFlags {
			c.command += "-X " + i.Key + "=" + i.Value + " "
			command += "-X " + i.Key + "=" + i.Value + " "
		}
		c.command2 = append(c.command2, command)
		c.command += `" `
	} else if len(c.Ldflags) != 0 && c.VarFlags == nil {
		c.command += `-ldflags "`
		c.command2 = append(c.command2, "-ldflags")
		for _, i := range c.Ldflags {
			c.command += i + " "
			command += i + " "
		}
		c.command2 = append(c.command2, command)
		c.command += `" `
	} else if len(c.Ldflags) == 0 && c.VarFlags != nil {
		c.command += `-ldflags "`
		c.command2 = append(c.command2, "-ldflags")
		for _, i := range c.VarFlags {
			c.command += "-X " + i.Key + "=" + i.Value + " "
			command += "-X " + i.Key + "=" + i.Value + " "
		}
		c.command2 = append(c.command2, command)
		c.command += `" `
	}
	command = ""
	if len(c.AsmFlags) != 0 {
		c.command += `-asmflags "`
		c.command2 = append(c.command2, "-asmflags")
		for _, i := range c.AsmFlags {
			c.command += i + " "
			command += i + " "
		}
		c.command2 = append(c.command2, command)
		c.command += `" `
	}
	command = ""
	if len(c.Gccgoflags) != 0 {
		c.command += `-gccgoflags "`
		c.command2 = append(c.command2, "-gccgoflags")
		for _, i := range c.Gccgoflags {
			c.command += i + " "
			command += i + " "
		}
		c.command2 = append(c.command2, command)
		c.command += `" `
	}
	command = ""
	if len(c.Gcflags) != 0 {
		c.command += `-gcflags "`
		c.command2 = append(c.command2, "-gcflags")
		for _, i := range c.Gcflags {
			c.command += i + " "
			command += i + " "
		}
		c.command2 = append(c.command2, command)
		c.command += `" `
	}
	debugTools.PrintlnOnlyInDebugMode("Command:", c.command)
	c.status = true
	return true
}

func (c *BuildConfig) Build() bool {
	if !c.status {
		logrus.Errorln("编译状态未就绪，请先执行ParseConfig再进行编译")
		return false
	}
	logrus.Infoln("开始编译")
	t, _ := os.Getwd()
	debugTools.PrintlnOnlyInDebugMode("编译命令:", c.command2, "当前路径:", t)
	if len(c.Targets) == 0 {
		logrus.Warningln("未设置编译目标，将默认编译目标为当前平台")
		c.Targets = append(c.Targets, getNowArch())
	}
	for _, target := range c.Targets {
		os.Setenv("GOOS", target.GOOS)
		os.Setenv("GOARCH", target.GOARCH)
		logrus.Infoln("编译平台：", target.GOOS, "/", target.GOARCH)
		if OutFileNameFmt == "a" {
			var flag = false //有没有-o
			for i := 0; i < len(c.command2); i++ {
				if c.command2[i] == "-o" {
					flag = true
					if ioTools.IsStrAInStrB("./bin/", c.command2[i+1]) {
						if runtime.GOOS == "windows" {
							c.command2[i+1] += c.outName + "-" + target.GOOS + "-" + target.GOARCH + ".exe"
						} else {
							c.command2[i+1] += c.outName + "-" + target.GOOS + "-" + target.GOARCH
						}
					}
				}
			}
			if !flag {
				if runtime.GOOS == "windows" {
					c.command2 = append(c.command2, "-o", c.outName+"-"+target.GOOS+"-"+target.GOARCH+".exe")
				} else {
					c.command2 = append(c.command2, "-o", c.outName+"-"+target.GOOS+"-"+target.GOARCH)
				}
			}
		}
		ioTools.GetOutputContinually2("go", c.command2...)
	}
	os.Chdir(global.RootDir)
	//<-ioTools.GetOutputContinually("go", "build", c.command)
	return true
}

func EnableCGO() {
	defaultConfig.BaseConfig.Cgo = true
}
func init() {
	os.Mkdir("project", 0750)
	var d = yamlConfig{}.BaseConfig
	d.RemoteConfig.LocalStoreTemp = "./project"
	defaultConfig.BaseConfig = d
	var d2 = jsonConfig{}.BaseConfig
	d2.RemoteConfig.LocalStoreTemp = "./project"
	defaultConfigJ.BaseConfig = d2
}
