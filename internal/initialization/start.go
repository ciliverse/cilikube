package initialization

import (
	"net"
	"os"

	"github.com/fatih/color"
)

// DisplayServerInfo 在控制台显示彩色的服务器运行信息。
// 它现在是一个可导出的公共函数。
func DisplayServerInfo(serverAddr, mode string) {
	version := getVersion() // 版本信息在内部获取

	color.Cyan("🚀 CiliKube Server is running!")
	color.Green("   ➜  Local:       http://127.0.0.1%s", serverAddr)
	color.Green("   ➜  Network:     http://%s%s", getLocalIP(), serverAddr)
	color.Yellow("  ➜  Mode:        %s", mode)
	color.Magenta("  ➜  Version:     %s", version)
	color.White("-------------------------------------------------")
}

// getVersion 从项目根目录的 VERSION 文件获取版本号
func getVersion() string {
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return "v0.2.4" // 如果读取失败，返回默认版本号
	}
	return string(data)
}

// getLocalIP 获取本机的局域网 IP 地址
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
