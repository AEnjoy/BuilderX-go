package iotools

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

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
