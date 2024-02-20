package iotools

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
	"sync"
)

func GetOutputDirectly(name string, args ...string) (output []byte) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output() // 等到命令执行完, 一次性获取输出
	if err != nil {
		panic(err)
	}
	output, err = simplifiedchinese.GB18030.NewDecoder().Bytes(output)
	if err != nil {
		panic(err)
	}
	return
}

// GetOutputContinually
// 不断输出到stdout, 直到结束
//
//	<-GetOutputContinually("tree")
func GetOutputContinually(name string, args ...string) <-chan struct{} {
	cmd := exec.Command(name, args...)
	closed := make(chan struct{})
	defer close(closed)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer func(stdoutPipe io.ReadCloser) {
		err := stdoutPipe.Close()
		if err != nil {

		}
	}(stdoutPipe)

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() { // 命令在执行的过程中, 实时地获取其输出
			data, err := simplifiedchinese.GB18030.NewDecoder().Bytes(scanner.Bytes()) // 防止乱码
			if err != nil {
				fmt.Println("transfer error with bytes:", scanner.Bytes())
				continue
			}

			fmt.Printf("%s\n", string(data))
		}
	}()

	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return closed
}
func GetOutputContinually2(name string, args ...string) {
	getOutput := func(reader *bufio.Reader) {
		var sumOutput string //统计屏幕的全部输出内容
		outputBytes := make([]byte, 200)
		for {
			n, err := reader.Read(outputBytes) //获取屏幕的实时输出(并不是按照回车分割，所以要结合sumOutput)
			if err != nil {
				if err == io.EOF {
					break
				}
				logrus.Errorln(err)
				sumOutput += err.Error()
			}
			output := string(outputBytes[:n])
			fmt.Print(output) //输出屏幕内容
			sumOutput += output
		}
		return
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin

	var wg sync.WaitGroup
	wg.Add(2)
	//捕获标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorln("ERROR:", err)
		os.Exit(1)
	}
	readout := bufio.NewReader(stdout)
	go func() {
		defer wg.Done()
		getOutput(readout)
	}()

	//捕获标准错误
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Errorln("ERROR:", err)
		os.Exit(1)
	}
	readerr := bufio.NewReader(stderr)
	go func() {
		defer wg.Done()
		getOutput(readerr)
	}()

	//执行命令
	err = cmd.Run()
	if err != nil {
		logrus.Errorln("ERROR:", err)
		os.Exit(1)
	}
	wg.Wait()
	return
}
