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

// Version allows injection via build parameters (go build -ldflags "-X 'github.com/ciliverse/cilikube/internal/initialization.Version=v0.3.1'")
var Version = ""

// DisplayServerInfo prints service startup information, including local/LAN addresses, mode, version, Go version, startup time, etc.
func DisplayServerInfo(serverAddr, mode string) {
	version := getVersion()
	goVersion := runtime.Version()
	buildTime := getBuildTime()
	hostname, _ := os.Hostname()
	// Set to Beijing time (UTC+8)
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

// getVersion gets version number, priority: environment variable > build variable > VERSION file > default value
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
	log.Printf("[WARN] Failed to get version number (environment variable, build variable, VERSION file all invalid), using default version: %v", err)
	return "v0.0.1"
}

// getBuildTime supports injecting build time via build parameters (go build -ldflags "-X 'github.com/ciliverse/cilikube/internal/initialization.BuildTime=2025-06-24T12:00:00Z'")
var BuildTime = ""

func getBuildTime() string {
	return BuildTime
}

// getLocalIP gets the first non-loopback IPv4 address of the local machine, commonly used for LAN access
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
