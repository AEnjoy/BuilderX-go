package ioTools

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
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
func TestChdir(t *testing.T) {
	getwd, _ := os.Getwd()
	t.Logf("Before Dir: %s", getwd)
	file, err := os.Open("../../version")
	defer file.Close()
	if err != nil {
		t.Errorf("%s", err)
	}
	err = os.Chdir(filepath.Dir("../../version"))
	if err != nil {
		t.Errorf("%s", err)
	}
	getwd, _ = os.Getwd()
	if err != nil {
		return
	}
	t.Logf("Now Dir: %s", getwd)
}
