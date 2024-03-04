package main

import (
	"github.com/aenjoy/BuilderX-go/builder"
	"os"
	"testing"
)

func TestBuild(t *testing.T) {
	t.Logf("测试构建")
	task := builder.UsingYaml("builderX.yaml", "BuilderX-Go")
	if len(task) == 0 {
		t.Errorf("No task found in yaml file. Exit.\n")
		//lock.Exit(1, "No task found in yaml file. Exit.")
	}
	for _, ta := range task {
		ta.Build()
		_, err := os.Stat(ta.Config.OutputFile)
		if err != nil {
			t.Errorf("Build failed.) %s\n", err.Error())
		}
	}

}
