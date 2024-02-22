package main

import (
	"github.com/aenjoy/BuilderX-go/builder"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/router"
	"github.com/aenjoy/BuilderX-go/utils/debugTools"
	"github.com/aenjoy/BuilderX-go/utils/lock"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// 预编译配置
var (
	Version   string
	BuildTime string
	GoVersion string
	GitTag    string
	GOOS      string
	GOARCH    string
)

// 命令行参数解析
var (
	server          = true
	version         = true
	notLoadDefault  = true
	cgo             = true
	local           string
	remote          string
	remoteCloneWay  string
	yaml            string
	json            string
	export          string
	exportType      string
	notRunningCheck bool
)

func init() {
	pflag.BoolVar(&notRunningCheck, "not-running-check", false, "不检查是否正在运行。如果指定了该参数，则 BuilderX将不检查服务是否正在运行,并可以运行多个实例。默认检查是否正在运行。")
	pflag.BoolVar(&notLoadDefault, "not-load-temple-default", false, "不加载默认模板配置文件(仅yaml,json不支持导入模板)。如果指定了该参数，则 BuilderX将不加载模板配置文件而使用内置配置。默认外部模板配置文件路径：当前目录下的config.yaml或/etc/BuilderX/config.yaml。")
	//var remoteBranch = pflag.StringP("remoteBranch", "b", "master", "远程项目分支。如果指定了该参数，则 BuilderX将使用该分支构建。默认master")
	global.WebPort = *pflag.StringP("port", "p", "18088", "Web管理面板的端口号。")
	pflag.BoolVarP(&server, "web", "w", false, "启动 BuilderX Web管理面板。如果指定了该参数，则 BuilderX 将启动 Web 管理面板，并忽略除--not-load-temple-default,--cgo,--port和--debug外其它命令行参数。")
	pflag.BoolVarP(&version, "version", "v", false, "Show BuilderX version and exit.")
	pflag.BoolVarP(&debugTools.DebugFlag, "debug", "d", false, "Debug mode 将显示一些额外信息，并忽略一些错误，可能会泄露某些数据。")
	pflag.BoolVarP(&cgo, "cgo", "c", false, "全局是否启用cgo。")
	pflag.StringVarP(&local, "local", "l", "null", "本地项目路径(或欲编译的文件路径)。如果指定了该参数，则 BuilderX将使用该目录构建。如果指定了但不选择地址，则使用当前目录。")
	pflag.Lookup("local").NoOptDefVal = "."
	pflag.StringVarP(&builder.OutFileNameFmt, "out-file-name-fmt", "F", "default", "输出文件名格式。如果指定了该参数，则 BuilderX将使用该格式构建,否则使用go默认输出格式。(default:使用go默认输出格式(packageName[.exe]),a:{package-name}-{os}-{arch}[.exe])")
	pflag.Lookup("out-file-name-fmt").NoOptDefVal = "a"

	pflag.StringVarP(&remote, "remote", "r", "null", "远程项目地址。如果指定了该参数，则 BuilderX将使用该地址构建。如果指定了但不选择地址，则使用BuilderX项目地址。格式：主机名[:端口]/用户名/项目名。")
	pflag.Lookup("remote").NoOptDefVal = "github.com/aenjoy/BuilderX-go"
	pflag.StringVar(&remoteCloneWay, "remote-clone-way", "https", "远程项目拉取方式。如果指定了该参数，则 BuilderX将使用该方式克隆远程项目。可选择的方式：https,git,ssh。")
	pflag.StringVarP(&yaml, "file-yaml", "Y", "null", "BuilderX-自动构建配置文件路径。如果指定了该参数，则 BuilderX将使用该文件(.yaml)进行构建。")
	pflag.StringVarP(&json, "file-json", "J", "null", "BuilderX-自动构建配置文件路径。如果指定了该参数，则 BuilderX将使用该文件(.json)进行构建。")
	pflag.StringVarP(&export, "export-conf", "e", "null", "导出一个配置文件示例。")
	pflag.StringVar(&exportType, "export-conf-type", "yaml", "默认使用yaml导出一个配置文件。支持yaml，json。")
	pflag.Parse()
}
func main() {
	if version {
		println("BuilderX Version:", Version)
		println("Build Time:", BuildTime)
		println("Go Version:", GoVersion)
		println("Go OS/Arch:", GOOS+"/"+GOARCH)
		println("Git Tag:", GitTag)
		os.Exit(0)
	}
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go lock.ExitHandle(exitChan)

	if !notLoadDefault {
		builder.LoadDefault()
	}
	if !notRunningCheck {
		lock.Lock()
	}
	if cgo {
		logrus.Infoln("全局配置:Enable CGO")
		builder.EnableCGO()
	}
	if server {
		// 启动 BuilderX Web管理面板
		// ...
		_, err := strconv.Atoi(global.WebPort)
		if err != nil {
			logrus.Errorln("WebPort 配置错误,将使用默认端口18088:", err)
			global.WebPort = "18088"
		}
		router.InitRouter()
		select {}
	}
	if export != "null" {
		if exportType == "yaml" {
			builder.ExportDefaultConfigYaml(export)
		} else if exportType == "json" {
			builder.ExportDefaultConfigJson(export)
		}
		lock.Exit(0)
	}
	if local != "null" {
		//使用本地路径构建
		task := builder.UsingLocal(local)
		if task.TaskID != "" {
			task.Build()
		} else {
			logrus.Errorln("No build task found in local path. Exit.")
			lock.Exit(1, "No build task found in local path. Exit.")
		}
		lock.Exit(0)
	}
	if remote != "null" {
		task := builder.UsingRemote(remote, remoteCloneWay)
		if len(task) == 0 {
			logrus.Errorln("No task found in yaml file. Exit.")
			lock.Exit(1, "No task found in yaml file. Exit.")
		}
		for _, t := range task {
			//t.Config.ParseConfig()
			t.Build()
		}
		lock.Exit(0)
	}
	if yaml != "null" {
		task := builder.UsingYaml(yaml, "Build from console.")
		if len(task) == 0 {
			logrus.Errorln("No task found in yaml file. Exit.")
			lock.Exit(1, "No task found in yaml file. Exit.")
		}
		for _, t := range task {
			t.Build()
		}
		lock.Exit(0)
	}
	if json != "null" {
		task := builder.UsingJson(json, "Build from console.")
		if len(task) == 0 {
			logrus.Errorln("No task found in json file. Exit.")
			lock.Exit(1, "No task found in json file. Exit.")
		}
		for _, t := range task {
			t.Build()
		}
		lock.Exit(0)
	}
	println("需要指定参数")
	pflag.Usage()
	lock.Exit(0)
}
