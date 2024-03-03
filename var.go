package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
)

// 预编译配置
var (
	Version         string
	BuildTime       string
	GoVersion       string
	GitTag          string
	Features        string
	Platform        string
	PlatformVersion string
)

func printVar() {
	println("BuilderX Version:", Version)
	println("BuildHost Build Time:", BuildTime)
	println("BuildHost Go Version:", GoVersion)
	println("BuildHost Platform:", Platform)
	println("BuildHost Platform Version:", PlatformVersion)
	println("BuilderX-Go Platform OS/Arch:", runtime.GOOS+"/"+runtime.GOARCH)
	println("Git Tag:", GitTag)
	println("Features:", Features)
	println("-----------------")
	printHardware()
	println("-----------------")

}
func printHardware() {
	info, err := cpu.Info()
	if err == nil && len(info) > 0 {
		cpus := info[0]
		println("CPUModel:", cpus.ModelName)
	}
	hi, err := host.Info()
	if err == nil {
		println("Platform:", hi.Platform)
		println("PlatformVersion:", hi.PlatformVersion)
		println("SystemArch:", hi.KernelArch)
	}
	vm, err := mem.VirtualMemory()
	if err == nil {
		fmt.Printf("MemoryTotal: %0.2f MB\n", float64(vm.Total)/1024/1024)
	}
}
