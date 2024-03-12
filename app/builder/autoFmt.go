package builder

import (
	"archive/zip"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

func UsingAuto(f string, taskName string) (task []Task) {
	stat, err := os.Stat(f)
	if err != nil {
		logrus.Errorln("Error opening file: ", f, err)
		return nil
	}
	if stat.IsDir() {
		return append(task, UsingLocal(f))
	}
	file, err := os.Open(f)
	if err != nil {
		logrus.Errorln("Error opening file: ", f, err)
		return nil
	}
	decoder := yaml.NewDecoder(file)
	var configY yamlConfig
	var configJ jsonConfig
	err = decoder.Decode(&configY)
	if err == nil {
		logrus.Infoln("Found YAML file.")
		return UsingYaml(f, taskName)
	}
	logrus.Infoln("Not YAML.", err)
	file.Close()
	file2, err := os.ReadFile(f)
	err = json.Unmarshal(file2, &configJ)
	if err == nil {
		logrus.Infoln("Found JSON file.")
		return UsingJson(f, taskName)
	}
	logrus.Infoln("Not JSON.", err)
	archive, err := zip.OpenReader(f)
	if err == nil {
		logrus.Infoln("Found ZIP file.")
		archive.Close()
		return append(task, UsingZip(f, taskName))
	}
	logrus.Error("Unknown file format.")
	return nil
}
