package builder

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/debugTools"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"github.com/aenjoy/BuilderX-go/utils/macro"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// BuildConfig
//
//	构建配置 本地项目
type BuildConfig struct {
	InputFile   string    //单个文件名或路径
	OutputFile  string    //输出文件名或路径
	outName     string    //命令提示中输出的文件名(程序包名)
	packageName string    //相对程序包名(最后包名)
	Ldflags     []string  //传递给链接参数的值
	VarFlags    []VarFlag //传递给main.go 中属性参数的值
	V           bool      // -v 打印编译的包名和文件名
	Cgo         bool      //enable cgo default false
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
	Targets      []BuildArch //构建的目标架构
	command      string
	command2     []string
	status       bool //true:就绪 false:未就绪
	MacroContext macro.Macro
	//
	HaveMacroBeforeCompile bool
	//
	Before
	Checksum
	Archives
	After
}

func (c *BuildConfig) ParseConfig() bool {
	if c.status { //防止重复处理
		return true
	}
	_, err := exec.LookPath(global.GoExe)
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
	c.parseBefore()
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
	} else {
		c.OutputFile = "./"
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
	c.packageName, err = getGoPackageName2()
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
		tOut := c.OutputFile
		os.Setenv("GOOS", target.GOOS)
		os.Setenv("GOARCH", target.GOARCH)
		logrus.Infoln("编译平台：", target.GOOS, "/", target.GOARCH)
		if c.HaveMacroBeforeCompile {
			//todo 编译前执行命令
		}
		ioTools.GetOutputContinually2(global.GoExe, c.command2...)
		info, err := os.Stat(tOut)
		if OutFileNameFmt == "a" {
			if err == nil {
				if info.IsDir() {
					if tOut[len(tOut)-1] != '/' {
						tOut += "/"
					}
					if target.GOOS == "windows" {
						//windows的程序.exe结尾
						fileName := tOut + c.packageName + ".exe"
						//os.Remove(c.OutputFile + c.packageName + "-" + target.GOOS + "-" + target.GOARCH + ".exe")
						err := os.Rename(fileName, tOut+c.packageName+"-"+target.GOOS+"-"+target.GOARCH+".exe")
						if err != nil {
							logrus.Errorln(err)
						}
						tOut = tOut + c.packageName + "-" + target.GOOS + "-" + target.GOARCH + ".exe"
					} else {
						fileName := tOut + c.packageName
						//os.Remove(c.OutputFile + c.packageName + "-" + target.GOOS + "-" + target.GOARCH)
						err := os.Rename(fileName, tOut+c.packageName+"-"+target.GOOS+"-"+target.GOARCH)
						if err != nil {
							logrus.Errorln(err)
						}
						tOut = tOut + c.packageName + "-" + target.GOOS + "-" + target.GOARCH
					}
				} else {
					//是文件 用户自己指定的,不需要修改
				}
			}
			_, err = os.Stat(tOut)
			if err != nil {
				logrus.Errorln("平台:"+target.GOOS, "/", target.GOARCH+"编译失败,找不到编译后的文件:", tOut)
			} else {
				logrus.Infoln("平台:"+target.GOOS, "/", target.GOARCH+"编译成功,文件路径:", tOut)
			}
		} else {
			//没有指定-F=a
			if info.IsDir() {
				if tOut[len(tOut)-1] != '/' {
					tOut += "/"
				}
				if target.GOOS == "windows" {
					tOut = tOut + c.packageName + ".exe"
				} else {
					tOut = tOut + c.packageName
				}
			}
			_, err = os.Stat(tOut)
			if err != nil {
				logrus.Errorln("平台:"+target.GOOS, "/", target.GOARCH+"编译失败,找不到编译后的文件:", tOut)
			} else {
				logrus.Infoln("平台:"+target.GOOS, "/", target.GOARCH+"编译成功,文件路径:", tOut)
			}
		}
	}
	os.Chdir(global.RootDir)
	//<-ioTools.GetOutputContinually("go", "build", c.command)
	return true
}

func (c *BuildConfig) parseBefore() bool {
	for _, s := range c.Before.Command {
		command := strings.Split(s, " ")
		ioTools.GetOutputDirectly(command[0], command[1:]...)
	}
	return true
}

func (c *BuildConfig) parseAfter() bool {
	for _, s := range c.After.Command {
		command := strings.Split(s, " ")
		ioTools.GetOutputDirectly(command[0], command[1:]...)
	}
	return true
}

func (c *BuildConfig) parseChecksum() bool {
	for _, s := range c.Checksum.File {
		output := ioTools.GetOutputDirectly("sha256sum", s)
		err := os.WriteFile(s+".sha256", []byte(output), 0644)
		if err != nil {
			logrus.Errorln("生成摘要文件:" + s + ".sha256失败")
		}
	}
	return true
}

func (c *BuildConfig) parseArchives(outArchivesFile, projectFile, projectFileOutFmt string) bool {
	if !c.Archives.Enable {
		return false
	}
	type fileInfo struct {
		file string
		path string
	}
	var files []fileInfo
	for _, t := range c.Archives.Files {
		t = c.MacroContext.ParserMacro(t)
		t1 := strings.Split(t, ":")
		//欲添加的文件路径:文件在压缩包中的路径
		var file, path string
		switch len(t1) {
		case 1:
			file = t1[0]
			path = t1[0]
		case 2:
			file = t1[0]
			if t1[1] != "" {
				path = t1[1]
			} else {
				path = t1[0]
			}
		default:
			logrus.Errorln("文件路径格式错误,跳过处理该文件")
			continue
		}
		files = append(files, fileInfo{file: file, path: path})
	}
	if c.Archives.Format == "zip" {
		archive, _ := os.Create(outArchivesFile)
		defer archive.Close()
		zipWriter := zip.NewWriter(archive)
		defer zipWriter.Close()
		for _, f := range files {
			//额外添加文件
			info, err := os.Stat(f.file)
			if err != nil {
				logrus.Errorln("打开文件失败:" + f.file)
				continue
			}
			if info.IsDir() {
				var dir []string
				dir, _ = ioTools.GetAllFile(f.file, dir)
				for _, s := range dir {
					file, err := os.Open(s)
					if err != nil {
						logrus.Errorln("打开文件失败:" + s)
						continue
					}
					defer file.Close()
					t2 := strings.Split(strings.Replace(s, "../", "", -1), "/")
					t2[0] = f.path
					//println("TargetFile:", s) "TargetFile: ../doc/command.md"
					//println("GetTargetDir:", strings.Join(t2, "/")) "GetTargetDir: docs/command.md"
					create, _ := zipWriter.Create(strings.Join(t2, "/"))
					io.Copy(create, file)
				}
				continue
			}
			file, err := os.Open(f.file)
			if err != nil {
				logrus.Errorln("打开文件失败:" + f.file)
				continue
			}
			defer file.Close()
			create, _ := zipWriter.Create(f.path)
			io.Copy(create, file)
		}
		//工程文件
		file, err := os.Open(projectFile)
		if err != nil {
			logrus.Errorln("打开生成文件失败:" + projectFile)
			return false
		}
		defer file.Close()
		create, _ := zipWriter.Create(projectFileOutFmt)
		io.Copy(create, file)
		return true
	} else { //tar家族
		archive, _ := os.Create(outArchivesFile)
		var tw *tar.Writer
		defer archive.Close()
		if c.Archives.Format == "tar.gz" {
			gw := gzip.NewWriter(archive)
			tw = tar.NewWriter(gw)
		} else if c.Archives.Format == "tar" {
			tw = tar.NewWriter(archive)
		}
		//todo tar.bzip2
		defer tw.Close()
		//工程文件
		files = append(files, fileInfo{file: projectFile, path: projectFileOutFmt})
		for _, f := range files {
			//额外添加文件
			info, err := os.Stat(f.file)
			if err != nil {
				logrus.Errorln("打开文件失败:" + f.file)
				continue
			}
			if info.IsDir() {
				var dir []string
				dir, _ = ioTools.GetAllFile(f.file, dir)
				for _, s := range dir {
					file, err := os.Open(s)
					if err != nil {
						logrus.Errorln("打开文件失败:" + s)
						continue
					}
					defer file.Close()
					t2 := strings.Split(strings.Replace(s, "../", "", -1), "/")
					t2[0] = f.path
					info, _ = os.Stat(s)
					header, err := tar.FileInfoHeader(info, "")
					header.Name = strings.Join(t2, "/")
					err = tw.WriteHeader(header)
					if err != nil {
						logrus.Errorln("写入文件头失败:" + s)
						continue
					}
					io.Copy(tw, file)
				}
				continue
			}
			file, err := os.Open(f.file)
			if err != nil {
				logrus.Errorln("打开文件失败:" + f.file)
				continue
			}
			defer file.Close()
			info, _ = os.Stat(f.file)
			header, err := tar.FileInfoHeader(info, "")
			header.Name = f.path
			err = tw.WriteHeader(header)
			if err != nil {
				logrus.Errorln("写入文件头失败:" + f.file)
				continue
			}
			io.Copy(tw, file)
		}
		return true
	}
}
