# yaml配置中BaseConfig-VarFlags宏指令说明

## 宏在各种自动化工具中是非常有用的,BuilderX支持使用特定的宏指令来生成配置文件,下面是BuilderX中VarFlags宏指令的说明文档:

目前在varFlags的value中定义了以下宏指令:

{command "cmd"} 执行cmd命令,并将执行结果作为value的值

{env "envName"} 获取环境变量envName的值,并将结果作为value的值

{file "filePath"} 读取filePath文件的内容,并将结果作为value的值

{json "jsonString"} 解析jsonString为json对象,并将结果作为value的值

{yaml "yamlString"} 解析yamlString为yaml对象,并将结果作为value的值

{base64 "base64String"} 将base64String解码为原始字符串,并将结果作为value的值

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
    - "main.Version=0.0.1-Dev"
    - "main.BuildTime=2024-2-21-18:25:19"
    - "main.GoVersion=1.21.5"
    - "main.GitTag={command "git rev-parse --short HEAD"}"
```

将当前版本库的git  commit id作为main.GitTag的值