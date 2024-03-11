package service

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func Install(name string) bool {
	if runtime.GOOS == "windows" {
		println("Error:Not supported on Windows.")
		return false
	}
	path, err := os.Executable()
	if err != nil {
		logrus.Errorln(err)
		return false
	}
	err = os.Symlink(path, "/usr/bin/"+name)
	if err != nil {
		logrus.Errorln(err)
		return false
	}
	return true
}
