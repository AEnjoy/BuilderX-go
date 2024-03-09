# BuilderX-Go

## 这是什么？

BuilderX-Go 是一个Go项目的命令行工具，用于构建生成目标为多平台的Go程序。它提供一个统一的构建命令，能够使用yaml格式的配置文件或json格式的配置文件描述自动化编译的各项功能和参数来进行构建。

同时，我们还提供了一个可视化的界面(正在开发中)，用户可以通过图形界面选择要构建的目标平台，并指定构建参数。

支持:远程仓库拉库自动构建,本地路径构建,预配置属性参数,高级命令,定时编译支持.


## 功能

- 支持多平台编译
- 支持远程仓库拉库自动构建
- 支持本地路径构建
- 支持预配置属性参数
- 支持高级命令
- 生成成功后自动打包
- 生成摘要文件

## 编译本项目:

在项目目录下运行命令 go run main.go var.go -l --out-file-name-fmt

程序将输出在./bin下

## 如何使用本工具

### 1.一次性使用

如果你是要编译远程仓库上的项目,

使用命令

```
BuilderX-Go -r="远程仓库地址" #如 BuilderX-Go -r="github.com/aenjoy/BuilderX-go" 
```

目标程序将生成在./project/项目名(BuilderX-go)/bin中

本地使用

在你要编译的项目路径,使用

```
BuilderX-Go -l .
#或
BuilderX-Go -l
#或
BuilderX-Go -l <targetDir>
```

目标将生成在目标路径的/bin目录下

### 2.多次使用

将自己平台的程序下载下来,解压后添加至环境变量,或使用命令

```
(sudo) BuilderX-Go --install (别名)
```

程序将自动安装至系统中

你可以使用命令build-go或自行设定的别名来启动程序.

### 3.命令行用法

[参考](doc/command.md)

### 4.配置文件的编写

[参考](doc/configFile.md)

如

```yaml
configType: build-config-local
configFile: []
configApiVersion: 3
configMinApiVersion: 1
name: "${define,`name`}"
define:
    -  "version=${file,`version`}"
    -  "name=BuilderX-Go"
baseConfig :
    inputFile : "."
    outputFile : "./bin/"
    ldflags :
        -  "-s"
        -  "-w"
        -  "-X main.Version=${using,`version`}"
        -  "-X main.BuildTime=${date,`2006-01-02--15:04:05`}"
        -  "-X main.GoVersion=${command,`go env GOVERSION`}"
        -  "-X main.GitTag=${command,`git rev-parse --short HEAD`}"
        -  "-X main.Features=NoWeb,CGO_ENABLED=${command,`go env CGO_ENABLED`},${env,`Features`}"
        -  "-X main.Platform=${command,`uname -s`}"
        -  "-X main.PlatformVersion=${command,`uname -r`}"
    v: true
    cgo: true
otherFlags:
    targets:
        - "windows/386"
        - "windows/amd64"
        - "windows/arm64"
        - "linux/386"
        - "linux/arm"
        - "linux/arm64"
        - "linux/amd64"
        - "darwin/amd64"
        - "darwin/arm64"
archives:
    enable: true
    name: "./bin/BuilderX-${using,`os`}-${using,`arch`}.zip"
    format: "zip"
    files:
        - "./readme.md:README.md"
        - "./LICENSE:LICENSE"
        - "./doc:doc"
        - "./res:res"
checksum:
    file:
        - "./bin/BuilderX-${using,`os`}-${using,`arch`}.zip"
```

