package main

import (
	"BuilderX/builder"
	"BuilderX/global"
	"BuilderX/router"
	"BuilderX/utils/debugTools"
	"BuilderX/utils/lock"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	Version   string
	BuildTime string
	GoVersion string
	GitTag    string
	GOOS      string
	GOARCH    string
)

func main() {
	var server = true
	var version = true
	var local = pflag.StringP("local", "l", "null", "本地项目路径(或欲编译的文件路径)。如果指定了该参数，则 BuilderX将使用该目录构建。如果指定了但不选择地址，则使用当前目录。")
	pflag.Lookup("local").NoOptDefVal = "."
	var remote = pflag.StringP("remote", "r", "null", "远程项目地址。如果指定了该参数，则 BuilderX将使用该地址构建。如果指定了但不选择地址，则使用BuilderX项目地址。格式：主机名[:端口]/用户名/项目名。")
	pflag.Lookup("remote").NoOptDefVal = "github.com/aenjoy/BuilderX"
	var remoteCloneWay = pflag.String("remote-clone-way", "https", "远程项目拉取方式。如果指定了该参数，则 BuilderX将使用该方式克隆远程项目。可选择的方式：https,git,ssh。")
	var yaml = pflag.StringP("file-yaml", "f", "null", "BuilderX-自动构建配置文件路径。如果指定了该参数，则 BuilderX将使用该文件(.yaml)进行构建。")
	var json = pflag.String("file-json", "null", "BuilderX-自动构建配置文件路径。如果指定了该参数，则 BuilderX将使用该文件(.json)进行构建。")
	var export = pflag.StringP("export-conf", "e", "null", "导出一个配置文件示例。")
	var exportType = pflag.String("export-conf-type", "yaml", "默认使用yaml导出一个配置文件。支持yaml，json。")
	//

	//var remoteBranch = pflag.StringP("remoteBranch", "b", "master", "远程项目分支。如果指定了该参数，则 BuilderX将使用该分支构建。默认master")
	global.WebPort = *pflag.StringP("port", "p", "18088", "Web管理面板的端口号。")
	pflag.BoolVarP(&server, "web", "w", false, "启动 BuilderX Web管理面板。如果指定了该参数，则 BuilderX 将启动 Web 管理面板，并忽略其它命令行参数。")
	pflag.BoolVarP(&version, "version", "V", false, "Show BuilderX version and exit.")
	pflag.BoolVarP(&debugTools.DebugFlag, "debug", "d", false, "Debug mode 将显示一些额外信息，并忽略一些错误，可能会泄露某些数据。")
	pflag.Parse()
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
	if *export != "null" {
		if *exportType == "yaml" {
			builder.ExportDefaultConfig(*export)
		}
		return
	}
	if server {
		// 启动 BuilderX Web管理面板
		// ...
		router.InitRouter()
		select {}
	}
	if *local != "null" {
		//使用本地路径构建
	}
	if *remote != "null" {
		println(*remote, *remoteCloneWay)
		//使用远程地址构建
	}
	if *yaml != "null" {
		task := builder.UsingYaml(*yaml, "Build from console.")
		if len(task) == 0 {
			logrus.Errorln("No task found in yaml file. Exit.")
			return
		}
		for _, t := range task {
			//t.Config.ParseConfig()
			t.Build()
		}
	}
	if *json != "null" {
		//使用json文件构建
	}
}
