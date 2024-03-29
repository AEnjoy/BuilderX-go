package builder

import (
	"archive/zip"
	"fmt"
	"github.com/aenjoy/BuilderX-go/global"
	"testing"
)

// zip url: https://github.com/casdoor/casdoor
func TestZipDir(t *testing.T) {
	zr, err := zip.OpenReader("test.zip")
	if err != nil {
		fmt.Println("Error creating zip reader:", err)
		return
	}
	defer zr.Close()
	files := zr.File
	for i, file := range files {
		if file.FileInfo().IsDir() {
			wantI := 0
			wantS := "casdoor-master/"
			if wantI != i || wantS != file.Name {
				t.Errorf("Error: want %d, %s, but got %d, %s", wantI, wantS, i, file.Name)
			}
			break
		}
	}
}
func TestUnzipFile(t *testing.T) {
	archive, err := zip.OpenReader("test.zip")
	defer archive.Close()
	tConfig := defaultConfigY
	for _, file := range archive.File {
		if err = unzipFile(file, tConfig.BaseConfig.RemoteConfig.LocalStoreTemp); err != nil {
			t.Errorf("unzip file error: %s\n", err)
		}
	}
}
func TestUsingZip(t *testing.T) {
	global.GoExe = "go"
	task := UsingZip("test.zip", "test")
	if task.TaskID == "" {
		t.Errorf("Error: task id is empty")
	}
	task.Build()
}

func TestBuildConfig_parseArchives(t *testing.T) {
	var config BuildConfig
	config.Format = "tar.gz" //zip,tar,tar.gz测试都通过
	config.Archives.Enable = true
	config.Archives.Files = []string{"apis.go:", "local.go:local_test_name.go", "../doc:docs", "../res:res", "../readme.md:readme.md"}
	config.parseArchives("logs/test.tar.gz", "../bin/BuilderX", "bin/BuilderX")
}
