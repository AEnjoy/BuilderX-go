package builder

import (
	"encoding/json"
	"fmt"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/debugTools"
	"github.com/aenjoy/BuilderX-go/utils/hashtool"
	"github.com/aenjoy/BuilderX-go/utils/ioTools"
	"github.com/aenjoy/BuilderX-go/utils/macro"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type jsonConfig struct {
	ConfigFile []string `json:"configFile"`
	BaseConfig struct {
		Cgo          bool     `json:"cgo"`
		InputFile    string   `json:"inputFile"`
		Ldflags      []string `json:"ldflags"`
		OutputFile   string   `json:"outputFile"`
		V            bool     `json:"v"`
		VarFlags     []string `json:"varFlags"`
		RemoteConfig struct {
			RemoteStore    string `json:"remoteStore"`
			LocalStoreTemp string `json:"localStoreTemp"`
			RemoteCloneWay string `json:"remote-clone-way"`
		} `json:"remoteConfig"`
	} `json:"baseConfig"`
	ConfigAPIVersion    int    `json:"configApiVersion"`
	ConfigMinAPIVersion int    `json:"configMinApiVersion"`
	ConfigType          string `json:"configType"`
	OtherFlags          struct {
		AsmFlags          []string `json:"asmFlags"`
		BuildMode         string   `json:"buildMode"`
		BuildProcess      int      `json:"buildProcess"`
		BuildVcs          string   `json:"buildVcs"`
		Compiler          string   `json:"compiler"`
		Cover             bool     `json:"cover"`
		CoverMode         string   `json:"coverMode"`
		ForceBuildPackage bool     `json:"forceBuildPackage"`
		Gccgoflags        []string `json:"gccgoflags"`
		Gcflags           []string `json:"gcflags"`
		Linkshared        bool     `json:"linkshared"`
		Mod               string   `json:"mod"`
		Modcacherw        bool     `json:"modcacherw"`
		Modfile           string   `json:"modfile"`
		Msan              bool     `json:"msan"`
		Overlay           string   `json:"overlay"`
		Pgo               string   `json:"pgo"`
		Pkgdir            string   `json:"pkgdir"`
		Race              bool     `json:"race"`
		Tags              string   `json:"tags"`
		Targets           []string `json:"targets"`
		Work              bool     `json:"work"`
	} `json:"otherFlags"`
	Before struct {
		Command []string `json:"command"`
	} `json:"before"`
	Checksum struct {
		File []string `json:"file"`
	} `json:"checksum"`
	Archives struct {
		Enable bool     `json:"enable"`
		Name   string   `json:"name"`
		Format string   `json:"format"`
		Files  []string `json:"files"`
	} `json:"archives"`
	After struct {
		Command []string `json:"command"`
	} `json:"after"`
	ProjectName string   `json:"name"`
	Define      []string `json:"define"`
}

var defaultConfigJ = jsonConfig{
	ConfigType:          C_Type_Local,
	ConfigAPIVersion:    global.ConfigApiVersion,
	ConfigMinAPIVersion: 1,
}

func UsingJson(f string, taskName string) (returnVal []Task) {
	var config jsonConfig
	logrus.Infoln("Using JSON: ", f, " parse...")
	file, err := os.ReadFile(f)
	if err != nil {
		logrus.Errorln("Error opening file: ", f, err)
		return nil
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		logrus.Errorln("Error decoding JSON:", err)
		return nil
	}
	if global.ConfigApiVersion < config.ConfigMinAPIVersion {
		logrus.Errorln("The current configuration version supported by BuilderX is too low to load the configuration file, and you should upgrade BuilderX.: SupportVersion:", global.ConfigApiVersion, " ConfigVersion:", config.ConfigMinAPIVersion)
		return nil
	}
	if global.ConfigApiVersion < config.ConfigAPIVersion {
		logrus.Warningln("BuilderX does not fully support some directives or fields in the configuration file, this problem is caused by the fact that the api version supported by BuilderX is lower than the api version of the configuration file. You should upgrade your BuilderX to ensure that the configuration file can be fully supported. SupportVersion:", global.ConfigApiVersion, "ConfigVersion:", config.ConfigAPIVersion)
	}
	os.Chdir(filepath.Dir(f))
	if config.ConfigType == C_Type_Multi {
		logrus.Infoln("Using multiple configs mode: ")
		for _, v := range config.ConfigFile {
			returnVal = append(returnVal, UsingJson(v, taskName)...)
		}
	} else if config.ConfigType == C_Type_Remote {
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
			global.BuiltTask++
			task.Config = jsonConfig2BuildConfig(config2)
			task.Config.OutputFile = "./bin/"
			task.Config.outName = url[2]
			if taskName == "" {
				taskName = config2.ProjectName
			}
			task.TaskID = hashtool.MD5(task.CreatTime.Format("2006-01-02-15:04:05") + strconv.Itoa(global.BuiltTask) + taskName)
			returnVal = append(returnVal, task)
		} else {
			logrus.Error("Error with remote config fmt.")
			return nil
		}
	} else if config.ConfigType == C_Type_Local {
		logrus.Infoln("Using local config mode: ")
		var task Task
		task.CreatTime = time.Now()
		global.BuiltTask++
		if taskName == "" {
			taskName = config.ProjectName
		}
		task.TaskID = hashtool.MD5(task.CreatTime.Format("2006-01-02-15:04:05") + strconv.Itoa(global.BuiltTask) + taskName)
		task.Config = jsonConfig2BuildConfig(config)
		returnVal = append(returnVal, task)
		// todo
	}
	logrus.Infoln("Config parsed successfully. ", f)
	return
}
func jsonConfig2BuildConfig(config jsonConfig) (returnVal BuildConfig) {
	var tMacro macro.Macro
	tMacro.SetMacroSplit(global.MacroSplit)
	tMacro.SetDefineContext(config.Define)
	returnVal.MacroContext = tMacro
	for _, v := range config.BaseConfig.VarFlags {
		var varFlag VarFlag
		a := strings.Split(v, "=")
		if len(a) >= 2 {
			varFlag.Key = tMacro.ParserMacro(a[0])
			varFlag.Value = tMacro.ParserMacro(strings.Join(a[1:], "="))
			returnVal.HaveMacroBeforeCompile = macro.HaveMacroBeforeCompile(strings.Join(a[1:], "="))
		} else {
			continue
		}
		returnVal.VarFlags = append(returnVal.VarFlags, varFlag)
	}
	parserArr := func(str []string) (retVal []string) {
		for _, s := range str {
			retVal = append(retVal, tMacro.ParserMacro(s))
		}
		return
	}
	returnVal.Ldflags = parserArr(config.BaseConfig.Ldflags)
	returnVal.V = config.BaseConfig.V
	returnVal.Cgo = config.BaseConfig.Cgo
	returnVal.InputFile = tMacro.ParserMacro(config.BaseConfig.InputFile)
	returnVal.OutputFile = tMacro.ParserMacro(config.BaseConfig.OutputFile)
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
	returnVal.Gcflags = parserArr(config.OtherFlags.Gcflags)
	returnVal.Linkshared = config.OtherFlags.Linkshared
	returnVal.Mod = config.OtherFlags.Mod
	returnVal.Modfile = tMacro.ParserMacro(config.OtherFlags.Modfile)
	returnVal.Modcacherw = config.OtherFlags.Modcacherw
	returnVal.Overlay = tMacro.ParserMacro(config.OtherFlags.Overlay)
	returnVal.Pgo = tMacro.ParserMacro(config.OtherFlags.Pgo)
	returnVal.Pkgdir = tMacro.ParserMacro(config.OtherFlags.Pkgdir)
	returnVal.Tags = tMacro.ParserMacro(config.OtherFlags.Tags)
	returnVal.After = After(config.After)
	returnVal.Before = Before(config.Before)
	returnVal.Checksum = Checksum(config.Checksum)
	returnVal.Archives = Archives(config.Archives)
	for _, v := range config.OtherFlags.Targets {
		a := strings.Split(v, "/")
		if len(a) == 2 {
			logrus.Info("Founding and adding build target ", a[0], "/", a[1])
			returnVal.Targets = append(returnVal.Targets, BuildArch{GOOS: a[0], GOARCH: a[1]})
		}
	}
	return
}
func ExportDefaultConfigJson(f string) {
	file, err := os.Create(f)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(&defaultConfigJ)
	if err != nil {
		logrus.Errorln("Error encoding JSON:", err)
		return
	}
	logrus.Infoln("JSON data saved to ", f)
}
