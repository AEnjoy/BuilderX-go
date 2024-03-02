# yaml配置中BaseConfig-VarFlags宏指令说明

## 宏在各种自动化工具中是非常有用的,BuilderX支持使用特定的宏指令来生成配置文件,下面是BuilderX中VarFlags宏指令的说明文档:

Key=Value Value部分支持多条宏指令,单条指令存放在字段 ${ } 中.一条Value可以包含多条指令.

${ } 中的内容会被替换为对应的值,如果${ }中的内容是空的,则不会进行替换,格式为${macro,args...}


目前在varFlags的value中定义了以下宏指令:(未打勾的表示还未支持,将在后续版本中支持)

- [x] ${command,\`cmd\`} 执行cmd命令,并将执行结果作为value的值
- [x] ${env,\`envName\`} 获取环境变量envName的值,并将结果作为value的值
- [x] ${file,\`filePath\`} 读取filePath文件的内容,并将结果作为value的值
- [ ] ${json,\`jsonFile\`,\`Config\`} 解析jsonFile为json对象,并将Config结果作为value的值
- [ ] ${yaml,\`yamlFile\`,\`Config\`} 解析yamlFile为yaml对象,并将Config结果作为value的值
- [x] ${base64,\`base64String\`} 将base64String解码为原始字符串,并将结果作为value的值
- [x] ${date,\`format\`} 获取当前时间,并按照format格式化,并将结果作为value的值 例如:2006-01-02--15:04:05

### 其它:

在宏指令前增加"!",表示在编译前执行该宏指令,例如:${!command,\`cmd\`}

## 示例:

```yaml
configType: build-config-local
configApiVision: 1
configMinApiVision: 1
baseConfig:
  inputFile: "."
  outputFile: ""
  ldflags:
    -  "-s"
    -  "-w"
  varFlags:
    -  "main.Version=${file,`version`}"
    -  "main.BuildTime=${date,`2006-01-02--15:04:05`}"
    -  "main.GoVersion=${command,`go env GOVERSION`}"
    -  "main.GitTag=${command,`git rev-parse --short HEAD`}"
    -  "main.Features=NoWeb ${env,`Features`}"
    -  "main.GOOS=${!env,`GOOS`}"
    -  "main.GOARCH=${!env,`GOARCH`}"
```

获取项目根目录下version文件的内容至main.Version

获取当前时间并按照2006-01-02--15:04:05格式化至main.BuildTime

获取系统中go的版本至main.GoVersion

将当前版本库的git  commit id作为main.GitTag的值

获取环境变量中Features作为main.Features的值

在编译前执行命令获取系统环境变量GOOS和GOARCH,并设置至main.GOOS和main.GOARCH的值