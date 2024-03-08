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

## define

[ ]类型:用于定义一些常量,在后续的配置处理过程中可以使用到.

定义在前面的define可以被后面的define覆盖或引用.

格式: 支持key=value 或使用指令 ${define,\`defineName\`,\`defineValue\`}

如

```yaml
define:
  - "MY_DEFINE_CONSTANT=1"
  - "MY_DEFINE_CONSTANT2=2"
  - "version=${file,`version`}" 
  - "${define,`defineName`,`defineValue`}" 
  - "v2=${using,`version`}" 
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

以下字段支持使用宏

----


## name

项目的名字


## before

该字段用于在打开项目,编译前执行指令.一般用于项目依赖的更新(在新版本的golang中,项目的编译将会自动执行依赖的更新).

### command

[ ]类型:,用于执行命令 

如

```yaml
before:
  command:
    - "go mod tidy -v"
```



## checksum

用于导出文件摘要

### file

需要生成摘要的文件文件名 当file不为空时,checksum启用

如

```yaml
checksum:
  file: 
    - "checksums.txt"
```



## archives

在编译后打包项目

### enable

bool 当值为true时,表示启用打包

### name

输出的压缩包文件名 支持使用宏

### format

压缩包的格式.支持 zip, tar, ~~tar.bz2,~~ tar.gz,

### files:

[ ]类型:额外打包的文件 支持使用宏

格式:欲添加的文件路径:文件在压缩包中的路径(可省略)

注意:文件路径末尾不能有斜杠. 如: ../doc/:docs是错误的, 应该写成:../doc:docs; ../doc:docs/也是错误的, 应该写成:../doc:docs .

例子:

```yaml
archives:
  enable: true
  name: "builderX-${targets `os`}-${targets `arch`}"
  format: "zip"
  files:
      - "./readme.md:README.md"
      - "./LICENSE:"
      - "./doc:docs"
```



## after

在打包后执行的操作

### command

[ ]类型:执行命令

例子:

```yaml
after:
  command:
    - "echo 编译成功"
```

