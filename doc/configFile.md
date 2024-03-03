# 配置项说明
    配置项说明以yaml格式举例,json格式类似

## configType

配置文件类型

> 支持的类型有:
> build-config-remote  使用远程仓库的项目,克隆至本地后编译
>
> build-config-local 使用本地的项目编译
>
> multiple-config  用于将多个配置文件(多个项目)聚合起来一次性编译

如

```yaml
configType: multiple-config
```

## configApiVersion 和 configMinApiVersion

用于描述配置文件的版本.

configApiVersion:用于描述当前配置项的api版本,当builderX版本大于或等于configApiVersion时可以使用到其中的所有功能.

configMinApiVersion:用于描述打开该配置项所需的最低builderX版本.当builderX的api版本低于configMinApiVersion的版本,builderX将拒绝读取该配置项.

如

```yaml
configApiVersion: 1
configMinApiVersion: 1
```



## configFile

该字段仅当configType: multiple-config时生效,用于加载多个编译配置.类型为[]

如

```yaml
configType: multiple-config
configApiVersion: 1
configMinApiVersion: 1
configFile:
  - "config1.yaml"
  - "config2.yaml"
```



## before

该字段用于在打开项目,编译前执行指令.一般用于项目依赖的更新(在新版本的golang中,项目的编译将会自动执行依赖的更新).

### command

[]类型,用于执行命令 

如

```yaml
before:
  command:
    - "go mod tidy -v"
```



## checksum

用于导出文件摘要

### file

输出的摘要文件名 当file不为空时,checksum启用

如

```yaml
checksum:
  file: "checksums.txt"
```



## archives

在编译后打包项目

### enable

bool 当值为true时,表示启用打包