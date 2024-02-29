# yaml配置中BaseConfig-VarFlags宏指令说明

## 宏在各种自动化工具中是非常有用的,BuilderX支持使用特定的宏指令来生成配置文件,下面是BuilderX中VarFlags宏指令的说明文档:

目前在varFlags的value中定义了以下宏指令:

- [x] ${command \`cmd\`} 执行cmd命令,并将执行结果作为value的值

- [x] ${env \`envName\`} 获取环境变量envName的值,并将结果作为value的值

- [x] ${file \`filePath\`} 读取filePath文件的内容,并将结果作为value的值
- [ ] ${json "jsonFile"  \`Config\`} 解析jsonFile为json对象,并将Config结果作为value的值
- [ ] ${yaml "yamlFile"  \`Config\`} 解析yamlFile为yaml对象,并将Config结果作为value的值
- [x] ${base64 \`base64String\`} 将base64String解码为原始字符串,并将结果作为value的值


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
    - "main.Version=${file `version`}"
    - "main.BuildTime=2024-2-21-18:25:19"
    - "main.GoVersion=${command `go env GOVERSION`}"
    - "main.GitTag=${command `git rev-parse --short HEAD`}"
    - "main.Features=${env `Features`}"
```

获取项目根目录下version文件的内容至main.Version

获取系统中go的版本至main.GoVersion

将当前版本库的git  commit id作为main.GitTag的值

获取环境变量中Features作为main.Features的值