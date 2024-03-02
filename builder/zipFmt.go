package builder

import (
	"archive/zip"
	"github.com/aenjoy/BuilderX-go/global"
	"github.com/aenjoy/BuilderX-go/utils/hashtool"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Zip格式支持从Github等代码仓库直接下载的zip文件

func unzipFile(file *zip.File, dir string) error {
	// Prevent path traversal vulnerability.
	// Such as if the file name is "../../../path/to/file.txt" which will be cleaned to "path/to/file.txt".
	name := strings.TrimPrefix(filepath.Join(string(filepath.Separator), file.Name), string(filepath.Separator))
	filePath := path.Join(dir, name)

	// Create the directory of file.
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Open the file.
	r, err := file.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	// Create the file.
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// Save the decompressed file content.
	_, err = io.Copy(w, r)
	return err
}

func UsingZip(f string, taskName string) (task Task) {
	logrus.Infoln("Using zip: ", f, " parse...")
	archive, err := zip.OpenReader(f)
	if err != nil {
		logrus.Errorln("zip open error: ", err)
		return
	}
	defer archive.Close()
	tConfig := defaultConfigY
	var packName string
	for i, file := range archive.File {
		if i == 0 && file.FileInfo().IsDir() {
			packName = strings.Replace(file.Name, "/", "", -1)
		}
		if err = unzipFile(file, tConfig.BaseConfig.RemoteConfig.LocalStoreTemp); err != nil {
			logrus.Errorln("unzip file error: ", err)
			return
		}
	}
	err = os.Chdir(tConfig.BaseConfig.RemoteConfig.LocalStoreTemp)
	if err != nil {
		os.Chdir("./project")
	}
	err = os.Chdir(packName)
	if err != nil {
		logrus.Errorln("change dir error: ", err)
		return
	}
	//todo
	tConfig.BaseConfig.InputFile = ""
	task.CreatTime = time.Now()
	global.BuildedTask++
	task.Config = yamlConfig2BuildConfig(tConfig)
	task.Config.OutputFile = "./bin/"
	task.Config.outName = packName
	task.TaskID = hashtool.MD5(task.CreatTime.Format("2006-01-02-15:04:05") + strconv.Itoa(global.BuildedTask) + taskName)
	return
}
