package ioTools

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

//export  FileReadAll
func FileReadAll(filename string) string {
	var str = ""
	file, _ := os.Open(filename)
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else {
			str = str + s
		}
	}
	logrus.Info("FileReadAll:", str)
	return str
}

func GetFileDir(filePath string) string {
	return filepath.Dir(filePath)
}
func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := os.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}
