package ioTools

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	filename := "../../version"
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
	t.Logf("%s", str)
}
