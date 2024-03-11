package builder

import (
	"gopkg.in/yaml.v3"
	"os"
)

func UsingRemote(url string, types string) []Task {
	tConfig := defaultConfigY
	tConfig.ConfigType = C_Type_Remote
	tConfig.BaseConfig.RemoteConfig.RemoteStore = url
	tConfig.BaseConfig.RemoteConfig.RemoteCloneWay = types
	file, _ := os.Create("logs/f")
	encoders := yaml.NewEncoder(file)
	encoders.Encode(&tConfig)
	file.Close()
	retval := UsingYaml("logs/f", "Remote")
	os.Remove("logs/f")
	return retval
}
