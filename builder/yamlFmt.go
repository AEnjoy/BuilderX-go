package builder

import (
	"fmt"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/debugTools"
	"github.com/aenjoy/BuilderX-go/utils/hashtool"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	ConfigType          string   `yaml:"configType"`
	ConfigFile          []string `yaml:"configFile"`
	ConfigApiVersion    int      `yaml:"configApiVersion"`
	ConfigMinApiVersion int      `yaml:"configMinApiVersion"`
	BaseConfig          struct {
		RemoteConfig struct {
			RemoteStore    string `yaml:"remoteStore"`
			LocalStoreTemp string `yaml:"localStoreTemp"`
			RemoteCloneWay string `yaml:"remote-clone-way"`
		} `yaml:"remoteConfig"`
		InputFile  string   `yaml:"inputFile"`
		OutputFile string   `yaml:"outputFile"`
		Ldflags    []string `yaml:"ldflags"`
		VarFlags   []string `yaml:"varFlags"` //支持{}宏
		V          bool     `yaml:"v"`        //verbose
		Cgo        bool     `yaml:"cgo"`
	} `yaml:"baseConfig"`
	OtherFlags struct {
		ForceBuildPackage bool     `yaml:"forceBuildPackage"`
		BuildProcess      int      `yaml:"buildProcess"`
		Race              bool     `yaml:"race"`
		Msan              bool     `yaml:"msan"`
		Cover             bool     `yaml:"cover"`
		CoverMode         string   `yaml:"coverMode"`
		Work              bool     `yaml:"work"`
		AsmFlags          []string `yaml:"asmFlags"`
		BuildMode         string   `yaml:"buildMode"`
		BuildVcs          string   `yaml:"buildVcs"`
		Compiler          string   `yaml:"compiler"`
		Gccgoflags        []string `yaml:"gccgoflags"`
		Gcflags           []string `yaml:"gcflags"`
		Linkshared        bool     `yaml:"linkshared"`
		Mod               string   `yaml:"mod"`
		Modcacherw        bool     `yaml:"modcacherw"`
		Modfile           string   `yaml:"modfile"`
		Overlay           string   `yaml:"overlay"`
		Pgo               string   `yaml:"pgo"`
		Pkgdir            string   `yaml:"pkgdir"`
		Tags              string   `yaml:"tags"`
		Targets           []string `yaml:"targets"`
	} `yaml:"otherFlags"`
}

var defaultConfigY = yamlConfig{
	ConfigType:          "build-config-local",
	ConfigApiVersion:    global.ConfigApiVersion,
	ConfigMinApiVersion: 1,
}

func UsingYaml(f string, taskName string) []Task {
	logrus.Infoln("Using YAML: ", f, " parse...")
	file, err := os.Open(f)
	if err != nil {
		logrus.Errorln("Error opening file: ", f, err)
		return nil
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var config yamlConfig
	err = decoder.Decode(&config)
	if err != nil {
		logrus.Errorln("Error decoding YAML:", err)
		return nil
	}
	if global.ConfigApiVersion < config.ConfigMinApiVersion {
		logrus.Errorln("The current configuration version supported by BuilderX is too low to load the configuration file, and you should upgrade BuilderX.: SupportVersion:", global.ConfigApiVersion, " ConfigVersion:", config.ConfigMinApiVersion)
		return nil
	}
	var returnVal []Task
	if config.ConfigType == "multiple-config" {
		logrus.Infoln("Using multiple configs mode: ")
		for _, v := range config.ConfigFile {
			returnVal = append(returnVal, UsingYaml(v, taskName)...)
		}
	} else if config.ConfigType == "build-config-remote" {
		logrus.Infoln("Using remote config mode: ")
		_, err = exec.LookPath("git")
		if err != nil {
			logrus.Errorln("Error checking for git - 请先安装git工具: ", err)
			return nil
		}
		config2 := config
		var url []string // github.com/aenjoy/BuilderX host/username/project
		url = strings.Split(config2.BaseConfig.RemoteConfig.RemoteStore, "/")
		debugTools.PrintlnOnlyInDebugMode("解析的URL:", url)
		if len(url) == 3 {
			logrus.Info("拉取远程仓库: ", config2.BaseConfig.RemoteConfig.RemoteStore, " 使用方法:", config2.BaseConfig.RemoteConfig.RemoteCloneWay)
			err = os.Chdir(config2.BaseConfig.RemoteConfig.LocalStoreTemp)
			if err != nil {
				//config2.BaseConfig.RemoteConfig.LocalStoreTemp = "./project/"
				//config2.BaseConfig.InputFile = "./project/" + url[2]
				os.Chdir("./project")
			}
			config2.BaseConfig.InputFile = ""
			if config2.BaseConfig.RemoteConfig.RemoteCloneWay == "https" {
				ioTools.GetOutputContinually2("git", "clone", "https://"+url[0]+"/"+url[1]+"/"+url[2])
			} else if config2.BaseConfig.RemoteConfig.RemoteCloneWay == "ssh" {
				s := "ssh://git@ssh." + url[0] + ":443/" + url[1] + "/" + url[2]
				ioTools.GetOutputContinually2("git", "clone", s)
			} else if config2.BaseConfig.RemoteConfig.RemoteCloneWay == "git" {
				s := "git@" + url[0] + ":" + url[1] + "/" + url[2]
				ioTools.GetOutputContinually2("git", "clone", s)
			}
			os.Chdir(url[2])
			debugTools.PrintlnOnlyInDebugMode("Debug:config2.BaseConfig.InputFile:", config2.BaseConfig.InputFile)
			var task Task
			task.CreatTime = time.Now()
			global.BuildedTask++
			task.Config = yamlConfig2BuildConfig(config2)
			task.Config.OutputFile = "./bin/"
			task.Config.outName = url[2]
			task.TaskID = hashtool.MD5(task.CreatTime.Format("2006-01-02-15:04:05") + strconv.Itoa(global.BuildedTask) + taskName)
			returnVal = append(returnVal, task)
		} else {
			logrus.Error("Error with remote config fmt.")
			return nil
		}
	} else if config.ConfigType == "build-config-local" {
		logrus.Infoln("Using local config mode: ")
		var task Task
		task.CreatTime = time.Now()
		global.BuildedTask++
		task.TaskID = hashtool.MD5(task.CreatTime.Format("2006-01-02-15:04:05") + strconv.Itoa(global.BuildedTask) + taskName)
		task.Config = yamlConfig2BuildConfig(config)
		returnVal = append(returnVal, task)
		// todo
	}
	logrus.Infoln("Config parsed successfully. ", f)
	return returnVal
}

// yamlConfig2BuildConfig
// yamlConfig to BuildConfig
func yamlConfig2BuildConfig(config yamlConfig) (returnVal BuildConfig) {
	for _, v := range config.BaseConfig.VarFlags {
		var varFlag VarFlag
		a := strings.Split(v, "=")
		if len(a) == 2 {
			varFlag.Key = a[0]
			varFlag.Value = ParserMacro(a[1])
		} else {
			continue
		}
		returnVal.VarFlags = append(returnVal.VarFlags, varFlag)
	}
	returnVal.Ldflags = config.BaseConfig.Ldflags
	returnVal.V = config.BaseConfig.V
	returnVal.Cgo = config.BaseConfig.Cgo
	returnVal.InputFile = config.BaseConfig.InputFile
	returnVal.OutputFile = config.BaseConfig.OutputFile
	returnVal.ForceBuildPackage = config.OtherFlags.ForceBuildPackage
	returnVal.BuildProcess = config.OtherFlags.BuildProcess
	returnVal.Race = config.OtherFlags.Race
	returnVal.Msan = config.OtherFlags.Msan
	returnVal.Cover = config.OtherFlags.Cover
	if config.OtherFlags.CoverMode == "set" {
		returnVal.CoverMode = Set
	} else if config.OtherFlags.CoverMode == "count" {
		returnVal.CoverMode = Count
	} else if config.OtherFlags.CoverMode == "atomic" {
		returnVal.CoverMode = Atomic
	}
	returnVal.Gcflags = config.OtherFlags.Gcflags
	returnVal.Linkshared = config.OtherFlags.Linkshared
	returnVal.Mod = config.OtherFlags.Mod
	returnVal.Modfile = config.OtherFlags.Modfile
	returnVal.Modcacherw = config.OtherFlags.Modcacherw
	returnVal.Overlay = config.OtherFlags.Overlay
	returnVal.Pgo = config.OtherFlags.Pgo
	returnVal.Pkgdir = config.OtherFlags.Pkgdir
	returnVal.Tags = config.OtherFlags.Tags
	for _, v := range config.OtherFlags.Targets {
		a := strings.Split(v, "/")
		if len(a) == 2 {
			logrus.Info("Founding and adding build target ", a[0], "/", a[1])
			returnVal.Targets = append(returnVal.Targets, BuildArch{GOOS: a[0], GOARCH: a[1]})
		}
	}
	return
}

func ExportDefaultConfigYaml(f string) {
	file, err := os.Create(f)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(&defaultConfigY)
	if err != nil {
		logrus.Errorln("Error encoding YAML:", err)
		return
	}
	logrus.Infoln("YAML data saved to ", f)
}

func loadConfigYaml(f string) (defaultConfig yamlConfig) {
	file, err := os.Open(f)
	if err != nil {
		fmt.Println("Error opening file, using default:", err)
		return defaultConfigY
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&defaultConfig)
	return
}

func LoadDefault() {
	logrus.Infoln("Loading default config:(如果你不想加载默认的文件配置而使用内置配置,请使用--not-load-temple-default选项)")
	_, err := os.Stat("config.yaml")
	if err != nil {
		_, err = os.Stat("/etc/BuilderX/config.yaml")
		if err == nil {
			defaultConfigY = loadConfigYaml("/etc/BuilderX/config.yaml")
			logrus.Infoln("Loaded config from /etc/BuilderX/config.yaml")
			return
		}
	} else {
		defaultConfigY = loadConfigYaml("config.yaml")
		logrus.Infoln("Loaded config from config.yaml")
		return
	}
	logrus.Warningln("No config file found, using built-in-default config.")
}
