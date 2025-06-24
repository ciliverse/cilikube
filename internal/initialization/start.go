package initialization

import (
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Version 允许通过编译参数注入（go build -ldflags "-X 'github.com/ciliverse/cilikube/internal/initialization.Version=v0.3.1'")
var Version = ""

// DisplayServerInfo 打印服务启动信息，包括本地/局域网地址、模式、版本号、Go版本、启动时间等
func DisplayServerInfo(serverAddr, mode string) {
	version := getVersion()
	goVersion := runtime.Version()
	buildTime := getBuildTime()
	hostname, _ := os.Hostname()
	// 设置为北京时间（东八区）
	loc, err := time.LoadLocation("Asia/Shanghai")
	var startTime string
	if err == nil {
		startTime = time.Now().In(loc).Format("2006-01-02 15:04:05 MST")
	} else {
		startTime = time.Now().Format("2006-01-02 15:04:05")
	}
	color.Cyan("🚀 CiliKube Server is running!")
	color.Green("   ➜  Local:       http://127.0.0.1%s", serverAddr)
	color.Green("   ➜  Network:     http://%s%s", getLocalIP(), serverAddr)
	color.Yellow("  ➜  Mode:        %s", mode)
	color.Magenta("  ➜  Version:     %s", version)
	color.Cyan("   ➜  Go Version:   %s", goVersion)
	color.Cyan("   ➜  Hostname:     %s", hostname)
	color.Cyan("   ➜  Start Time:   %s", startTime)
	if buildTime != "" {
		color.Cyan("   ➜  Build Time:   %s", buildTime)
	}
	color.White("-------------------------------------------------")
}

// getVersion 获取版本号，优先级：环境变量 > 编译变量 > VERSION 文件 > 默认值
func getVersion() string {
	if v := os.Getenv("CILIKUBE_VERSION"); v != "" {
		return v
	}
	if Version != "" {
		return Version
	}
	data, err := os.ReadFile("VERSION")
	if err == nil {
		return strings.TrimSpace(string(data))
	}
	log.Printf("[WARN] 获取版本号失败（环境变量、编译变量、VERSION 文件均无效），使用默认版本号: %v", err)
	return "v0.3.1"
}

// getBuildTime 支持通过编译参数注入构建时间（go build -ldflags "-X 'github.com/ciliverse/cilikube/internal/initialization.BuildTime=2025-06-24T12:00:00Z'")
var BuildTime = ""

func getBuildTime() string {
	return BuildTime
}

// getLocalIP 获取本机第一个非回环 IPv4 地址，常用于局域网访问
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "unknown"
}
