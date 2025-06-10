package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ciliverse/cilikube/internal/app"
)

func main() {
	// --- 1. 获取配置路径 ---
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件路径 (例如：configs/config.yaml)")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CILIKUBE_CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	// --- 2. 创建并运行应用 ---
	application, err := app.New(configPath)
	if err != nil {
		// 使用 slog 来记录致命错误
		slog.Error("应用初始化失败", "error", err)
		os.Exit(1)
	}

	// 启动应用
	application.Run()
}
