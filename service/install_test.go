package service

import (
	"os"
	"testing"
)

func TestPath(t *testing.T) {
	path, _ := os.Executable()
	t.Logf("%s", path)
}
