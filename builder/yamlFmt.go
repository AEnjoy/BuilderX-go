package builder

import (
	"BuilderX/global"
	"BuilderX/utils/debugTools"
	"BuilderX/utils/hashtool"
	"BuilderX/utils/iotools"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	ConfigType string   `yaml:"configType"`
	ConfigFile []string `yaml:"configFile"`
	BaseConfig struct {
		RemoteConfig struct {
			RemoteStore    string `yaml:"remoteStore"`
			LocalStoreTemp string `yaml:"localStoreTemp"`
			RemoteCloneWay string `yaml:"remote-clone-way"`
		} `yaml:"remoteConfig"`
		InputFile  string   `yaml:"inputFile"`
		OutputFile string   `yaml:"outputFile"`
		Ldflags    []string `yaml:"ldflags"`
		VarFlags   []string `yaml:"varFlags"`
		V          bool     `yaml:"v"` //verbose
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

func UsingYaml(f string, taskName string) []Task {
	logrus.Info("Using YAML: ", f, " parse...")
	file, err := os.Open(f)
	if err != nil {
		logrus.Error("Error opening file: ", f, err)
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
	var returnVal []Task
	if config.ConfigType == "multiple-config" {
		logrus.Info("Using multiple configs mode: ")
		for _, v := range config.ConfigFile {
			returnVal = append(returnVal, UsingYaml(v, taskName)...)
		}
	} else if config.ConfigType == "build-config-remote" {
		logrus.Info("Using remote config mode: ")
		_, err = os.Stat("git")
		if err != nil {
			logrus.Error("Error checking for git: 请先安装git工具", err)
			return nil
		}
		config2 := config
		var url []string // github.com/aenjoy/BuilderX host/username/project
		url = strings.Split(config2.BaseConfig.RemoteConfig.RemoteStore, "/")
		debugTools.PrintLogsOnlyInDebugMode("解析的URL:", url)
		if len(url) == 3 {
			logrus.Info("拉取远程仓库. ", "使用方法:", config2.BaseConfig.RemoteConfig.RemoteCloneWay)
			if config2.BaseConfig.RemoteConfig.RemoteCloneWay == "https" {
				<-iotools.GetOutputContinually("git", "clone", "https://"+url[0]+"/"+url[1]+"/"+url[2])
			} else if config2.BaseConfig.RemoteConfig.RemoteCloneWay == "ssh" {
				s := "ssh://git@ssh." + url[0] + ":443/" + url[1] + "/" + url[2]
				<-iotools.GetOutputContinually("git", "clone", s)
			} else if config2.BaseConfig.RemoteConfig.RemoteCloneWay == "git" {
				s := "git@" + url[0] + ":" + url[1] + "/" + url[2]
				<-iotools.GetOutputContinually("git", "clone", s)
			}
			config2.BaseConfig.InputFile = config2.BaseConfig.RemoteConfig.LocalStoreTemp + url[2] + "/."
			var task Task
			task.CreatTime = time.Now()
			global.BuildedTask++
			task.Config = localParse(config2)
			task.TaskID = hashtool.MD5(task.CreatTime.Format("2006-01-02-15:04:05") + strconv.Itoa(global.BuildedTask) + taskName)
			returnVal = append(returnVal, task)
		} else {
			logrus.Error("Error with remote config")
			return nil
		}
	} else if config.ConfigType == "build-config-local" {
		logrus.Info("Using local config mode: ")
		var task Task
		task.CreatTime = time.Now()
		global.BuildedTask++
		task.TaskID = hashtool.MD5(task.CreatTime.Format("2006-01-02-15:04:05") + strconv.Itoa(global.BuildedTask) + taskName)
		task.Config = localParse(config)
		returnVal = append(returnVal, task)
		// todo
	}
	logrus.Info("Config parsed successfully. ", f)
	return returnVal
}
func localParse(config yamlConfig) BuildConfig {
	var returnVal BuildConfig
	for _, v := range config.BaseConfig.VarFlags {
		var varFlag VarFlag
		a := strings.Split(v, "=")
		if len(a) == 2 {
			varFlag.Key = a[0]
			varFlag.Value = a[1]
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
	return returnVal
}
func ExportDefaultConfig(f string) {
	i := yamlConfig{
		ConfigType: "build-config-local",
	}

	file, err := os.Create(f)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(&i)
	if err != nil {
		fmt.Println("Error encoding YAML:", err)
		return
	}
	fmt.Println("YAML data saved to ", f)
}
