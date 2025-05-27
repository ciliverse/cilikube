package initialization

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/ciliverse/cilikube/configs"
	"github.com/fatih/color"
)

// StartServer 启动 HTTP 服务器
func StartServer(cfg *configs.Config, router http.Handler) {
	serverAddr := ":" + cfg.Server.Port
	version := getVersion()

	// 动态获取运行模式
	mode := os.Getenv("CILIKUBE_MODE")
	if mode == "" {
		mode = "development" // 默认模式
	}

	// 显示服务器运行信息
	displayServerInfo(serverAddr, mode, version)

	// 启动服务器
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// displayServerInfo 显示服务器运行信息
func displayServerInfo(serverAddr, mode, version string) {
	color.Cyan("🚀 CiliKube Server is running!")
	color.Green("   ➜  Local:      http://127.0.0.1%s", serverAddr)
	color.Green("   ➜  Network:    http://%s%s", getLocalIP(), serverAddr)
	color.Yellow("   ➜  Mode:       %s", mode)
	color.Magenta("   ➜  Version:    %s", version)
}

// getVersion 获取版本号
func getVersion() string {
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return "v0.1.3" // 如果读取失败，返回默认版本号
	}
	return string(data)
}

// getLocalIP 获取本地 IP 地址（用于 Network 信息）
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
