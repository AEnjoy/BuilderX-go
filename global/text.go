package global

const Help = `WebServerOption:
	-p, --port  (default "18088") web Port
	-w, --web    Start Dashboard
	参数-w与LocalToolsOption冲突.
LocalToolsOption:
	-F, --out-file-name-fmt[="a"] or (default "default")
	-l, --local [="."]
	-r, --remote [="github.com/aenjoy/BuilderX-go"]
		--remote-clone-way [way] way:https,git,ssh。 (default "https")
	-f, --file-auto fileName
	    --file-json fileName
	    --file-yaml fileName
	    --file-zip  fileName
	-e, --export-conf fileName
		--export-conf-type yaml or json (default "yaml")
	-n, --project-name string
GlobalOption:
	--not-running-check
	--not-load-temple-default
	-c, --cgo		Enable CGO
	-d, --debug
	-g, --go-exe string[="go"]
------------
	-h, --help
	-v, --version
	指定-h或-v将忽略其它参数
`
const BuildMode = `-buildmode=archive
                将列出的非主程序包构建到.a文件中。名为main的包将被忽略。
        -buildmode=c-archive
                构建列出的主程序包，以及它导入的所有程序包，并转换为c-archive(.a)文件。唯一可以调用的符号将是
				使用cgo//export 注释导出的函数。注意:只能列出一个主要包。

        -buildmode=c-shared
                构建列出的主程序包，以及它导入的所有程序包，并转换为c-shared(共享库.so/dll)文件。唯一可以调用的符号将是
				使用cgo//export 注释导出的函数。注意:只能列出一个主要包。

        -buildmode=default
                列出的主程序包内置在可执行文件中,并将非主包内置在.a文件中.（默认行为）

        -buildmode=shared
				将所有列出的非主包组合到一个共享库中，该选项将在使用-linkshared选项构建时使用。main包中的函数将被忽略。

        -buildmode=exe
                构建列出的主包以及它们导入到可执行文件中的所有内容。main以外的包将被忽略。

        -buildmode=pie
				构建列出的主包以及它们导入到位置独立可执行文件（PIE）中的所有内容。main以外的包将被忽略。

        -buildmode=plugin
                将列出的主包以及它们导入的所有包构建到Go插件中。main以外的包将被忽略。
`
const BuildModeSupported = "c-archive|c-shared|default|exe|pie|plugin|shared|archive"
const BuildVcsSupported = "auto|true|false"
const CompileSupported = "gccgo|gc"
const ModSupported = "readonly|vendor|or|mod"

// GoOSAndGoArchSupported
// 使用go tool dist list获取
var GoOSAndGoArchSupported = `aix/ppc64
android/386
android/amd64
android/arm
android/arm64
darwin/amd64
darwin/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
freebsd/arm64
freebsd/riscv64
illumos/amd64
ios/amd64
ios/arm64
js/wasm
linux/386
linux/amd64
linux/arm
linux/arm64
linux/loong64
linux/mips
linux/mips64
linux/mips64le
linux/mipsle
linux/ppc64
linux/ppc64le
linux/riscv64
linux/s390x
netbsd/386
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/386
openbsd/amd64
openbsd/arm
openbsd/arm64
plan9/386
plan9/amd64
plan9/arm
solaris/amd64
wasip1/wasm
windows/386
windows/amd64
windows/arm
windows/arm64`
