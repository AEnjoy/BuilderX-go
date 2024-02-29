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
	Features  string
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
	zip             string
	export          string
	exportType      string
	notRunningCheck bool
)

var flagsetW *pflag.FlagSet

func init() {
	flagsetW = pflag.NewFlagSet("BuilderX-Go", pflag.ExitOnError)
	global.WebPort = *flagsetW.StringP("port", "p", "18088", "Web管理面板的端口号。")
	flagsetW.BoolVarP(&server, "web", "w", false, "启动 BuilderX Web管理面板。如果指定了该参数，则 BuilderX 将启动 Web 管理面板，并忽略--local及以下的其它命令行参数。")
	flagsetW.BoolVar(&notRunningCheck, "not-running-check", false, "不检查是否正在运行。如果指定了该参数，则 BuilderX将不检查服务是否正在运行,并可以运行多个实例。默认检查是否正在运行。")
	flagsetW.BoolVar(&notLoadDefault, "not-load-temple-default", false, "不加载默认模板配置文件(仅yaml,json不支持导入模板)。如果指定了该参数，则 BuilderX将不加载模板配置文件而使用内置配置。默认外部模板配置文件路径：当前目录下的config.yaml或/etc/BuilderX/config.yaml。")
	//var remoteBranch = pflag.StringP("remoteBranch", "b", "master", "远程项目分支。如果指定了该参数，则 BuilderX将使用该分支构建。默认master")

	flagsetW.BoolVarP(&version, "version", "v", false, "Show BuilderX version and exit.")
	flagsetW.BoolVarP(&debugTools.DebugFlag, "debug", "d", false, "Debug mode 将显示一些额外信息，并忽略一些错误，可能会泄露某些数据。")
	flagsetW.BoolVarP(&cgo, "cgo", "c", false, "全局是否启用cgo。")
	flagsetW.StringVarP(&global.GoExe, "go-exe", "g", "go", "Go编译器路径。如果指定了该参数,这使用指定路径的Go编译器进行构建,否则使用$path中的go编译。")
	flagsetW.Lookup("go-exe").NoOptDefVal = "go"

	flagsetW.StringVarP(&local, "local", "l", "null", "本地项目路径(或欲编译的文件路径)。如果指定了该参数，则 BuilderX将使用该目录构建。如果指定了但不选择地址，则使用当前目录。")
	flagsetW.Lookup("local").NoOptDefVal = "."
	flagsetW.StringVarP(&builder.OutFileNameFmt, "out-file-name-fmt", "F", "default", "输出文件名格式。如果指定了该参数，则 BuilderX将使用该格式构建,否则使用go默认输出格式。(default:使用go默认输出格式(packageName[.exe]),a:{package-name}-{os}-{arch}[.exe])")
	flagsetW.Lookup("out-file-name-fmt").NoOptDefVal = "a"
	flagsetW.StringVarP(&remote, "remote", "r", "null", "远程项目地址。如果指定了该参数，则 BuilderX将使用该地址构建。如果指定了但不选择地址，则使用BuilderX项目地址。格式：主机名[:端口]/用户名/项目名。")
	flagsetW.Lookup("remote").NoOptDefVal = "github.com/aenjoy/BuilderX-go"
	flagsetW.StringVar(&remoteCloneWay, "remote-clone-way", "https", "远程项目拉取方式。如果指定了该参数，则 BuilderX将使用该方式克隆远程项目。可选择的方式：https,git,ssh。")
	flagsetW.StringVarP(&yaml, "file-yaml", "Y", "null", "BuilderX-自动构建配置文件路径。如果指定了该参数，则 BuilderX将使用该文件(.yaml)进行构建。")
	flagsetW.StringVarP(&json, "file-json", "J", "null", "BuilderX-自动构建配置文件路径。如果指定了该参数，则 BuilderX将使用该文件(.json)进行构建。")
	flagsetW.StringVarP(&zip, "file-zip", "Z", "null", "BuilderX-自动构建的仓库zip包。如果指定了该参数，则 BuilderX将使用该文件(.zip)进行构建。")
	flagsetW.StringVarP(&export, "export-conf", "e", "null", "导出一个配置文件示例。")
	flagsetW.StringVar(&exportType, "export-conf-type", "yaml", "默认使用yaml导出一个配置文件。支持yaml，json。")
	//pflag.Parse()
	flagsetW.SortFlags = false
	flagsetW.Parse(os.Args[1:])
	//pflag.Parse()

}
func main() {
	if version {
		println("BuilderX Version:", Version)
		println("Build Time:", BuildTime)
		println("Go Version:", GoVersion)
		println("Features:", Features)
		println("Go OS/Arch:", GOOS+"/"+GOARCH)
		println("Git Tag:", GitTag)
		os.Exit(0)
	}
	var hasBuildTask bool
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
		hasBuildTask = true
		router.InitRouter()
		select {}
	}
	if export != "null" {
		if exportType == "yaml" {
			builder.ExportDefaultConfigYaml(export)
		} else if exportType == "json" {
			builder.ExportDefaultConfigJson(export)
		}
		hasBuildTask = true
	}
	if local != "null" {
		//使用本地路径构建
		task := builder.UsingLocal(local)
		if task.TaskID != "" {
			task.Build()
		} else {
			logrus.Errorln("No build task found in local path. Exit.")
			//lock.Exit(1, "No build task found in local path. Exit.")
		}
		hasBuildTask = true
	}
	if remote != "null" {
		task := builder.UsingRemote(remote, remoteCloneWay)
		if len(task) == 0 {
			logrus.Errorln("No task found in yaml file. Exit.")
			//lock.Exit(1, "No task found in yaml file. Exit.")
		}
		for _, t := range task {
			//t.Config.ParseConfig()
			t.Build()
		}
		hasBuildTask = true
	}
	if yaml != "null" {
		task := builder.UsingYaml(yaml, "Build from console.")
		if len(task) == 0 {
			logrus.Errorln("No task found in yaml file. Exit.")
			//lock.Exit(1, "No task found in yaml file. Exit.")
		}
		for _, t := range task {
			t.Build()
		}
		hasBuildTask = true
	}
	if json != "null" {
		task := builder.UsingJson(json, "Build from console.")
		if len(task) == 0 {
			logrus.Errorln("No task found in json file. Exit.")
			//lock.Exit(1, "No task found in json file. Exit.")
		}
		for _, t := range task {
			t.Build()
		}
		hasBuildTask = true
	}
	if zip != "null" {
		task := builder.UsingZip(zip, "Build from console.")
		if task.TaskID != "" {
			task.Build()
		} else {
			logrus.Errorln("No build task found in zip path. Exit.")
		}
		hasBuildTask = true
	}
	if !hasBuildTask {
		println("需要指定参数. 请使用 -h 或者 --help 获取更多帮助信息.")
		pflag.Usage()
		println(global.Help)
		//flagsetW.Usage()
		lock.UnLock()
		os.Exit(1)
	}
	lock.UnLock()
	lock.Exit(0)
}
