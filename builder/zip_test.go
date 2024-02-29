package builder

import (
	"archive/zip"
	"fmt"
	"testing"
)

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
