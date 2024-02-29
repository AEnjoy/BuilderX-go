# 命令说明:

BuilderX-go支持在本地使用命令行做任务构建.

> Usage of BuilderX:
> 
> WebServerOption:
> 
>         -p, --port  (default "18088") web Port
> 
>         -w, --web    Start Dashboard
> 
>         参数-w与LocalToolsOption冲突.
> 
> LocalToolsOption:
> 
>         -F, --out-file-name-fmt[="a"] or (default "default")
> 
>         -l, --local [="."]
> 
>         -r, --remote [="github.com/aenjoy/BuilderX-go"]
> 
>                 --remote-clone-way [way] way:https,git,ssh。 (default "https")
> 
>         -J, --file-json fileName
> 
>         -Y, --file-yaml fileName
>  
>         -Z, --file-zip  fileName
> 
>         -e, --export-conf fileName
> 
>                 --export-conf-type yaml or json (default "yaml")
> 
> GlobalOption:
> 
>         --not-running-check
> 
>         --not-load-temple-default
> 
>         -c, --cgo               Enable CGO
> 
>         -d, --debug
>
> ----------------------
>
> ​		-h,--help
>
> ​		-v, --version

LocalToolsOption参数组与LocalToolsOption存在冲突,当指定了不同参数组时程序将出现与预期不同的执行结果.

## WebServerOption:

与DashBoard相关的命令组

### -w, --web

启动webDashboard程序,我们可以通过浏览器可视化的构建一个或多个,本地或远程的项目.

暴露的服务运行在-p port上,默认使用18088端口.

### -p, --port

 指定web服务监听的端口

## LocalToolsOption:

在本地构建项目而不启动web相关的命令组:

### -F, --out-file-name-fmt[="a"]

本地构建时输出的文件名格式 默认指定为-F default

#### 指定-F default 或不指定-F

 使用go默认的输出格式:即输出为包名或文件名[.exe]

-F=default,--out-file-name-fmt=default 同-F default

例如,当我们在Windows环境下,在builderX根目录下编译程序

使用go  run main.go -l . ,程序将输出在./BuilderX-go.exe

#### 指定-F a

-F a 使用内置的模板A作为输出格式: {package-name(dir)}-{os}-{arch}[.exe]

--out-file-name-fmt=a和-F a和-F

例如,当我们在Windows环境下,在builderX根目录下编译程序

使用go  run main.go -l . -F ,程序将输出在./github.com/aenjoy/BuilderX-go-windows-amd64.exe

###  -l, --local [="."]

使用系统指定路径下的项目或文件编译.默认为当前路径

支持目录(存在go.mod),以go结尾的go文件,或go.mod文件

示例:

-l

-l .

-l=.

-l=~/builderX/main.go

###  -r, --remote [="github.com/aenjoy/BuilderX-go"]

使用远程项目构建

若指定本参数,将自动从远程仓库地址拉取到本地进行构建.使用--remote-clone-way 指定的方式克隆到本地后再进行构建

本地临时文件夹将存放在./project或默认配置指定的其它路径中

例如,将github.com/aenjoy/BuilderX-go拉取到本地进行构建.

go run main.go -r="github.com/aenjoy/BuilderX-go"

项目将输出在./project/BuilderX-go/bin/BuilderX-go

如果指定了-F,即go run main.go -r="github.com/aenjoy/BuilderX-go" -F

项目将输出在./project./BuilderX-go./bin./github.com./aenjoy/BuilderX-go-{os}-{arch}[.exe]

子参数:

####  --remote-clone-way [way] way:https,git,ssh。 (default "https")

在克隆项目时使用的克隆方式

默认使用https方式,如果使用git和ssh方式,请先在本地配置授权

###  -J, --file-json fileName

使用fileName json配置进行构建

解析配置并进行构建.

###  -Y, --file-yaml fileName

使用fileName yaml配置进行构建

解析配置并进行构建.

[具体参考configFile.md](configFile.md)

### -Z, --file-zzip fileName

使用zip包进行构建,zip来自Github或gitlab等直接下载获得

###  -e, --export-conf fileName

从内存中导出当前的配置文件模板,并将配置文件导出到fileName

###  --export-conf-type yaml or json (default "yaml")

导出的文件类型,默认为yaml

## GlobalOption:

全局选项:全局选项不与LocalToolsOption和WebServerOption冲突,可混合使用.

### --not-running-check

不进行运行时冲突检查:例如,当你要运行两个或更多builderX实例时,第二个builderX将拒绝执行,除非指定了该选项

### --not-load-temple-default

不加载外部的模板配置文件,而使用内存中默认的配置文件

[具体参考configFile.md](configFile.md)

### -c, --cgo

启用cgo编译项目

###  -d, --debug

显示一些额外的调试信息,有助于在你遇到错误时将错误日志发生给我们以进行调试追踪.

调试模式收集的数据:全部的输入与输出,系统环境参数等(不包括个人隐私信息)

### -g, --go-exe string[="go"]

指定go程序的路径 适用于安装了多个go程序的情况.默认值为go,即使用系统环境变量中的第一go程序.

## 其它

###  -h, --help

显示程序的帮助信息

###  -v, --version

显示BuilderX的版本和一些编译信息.