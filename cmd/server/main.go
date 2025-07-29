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

	// 支持环境变量配置路径
	if configPath == "" {
		configPath = os.Getenv("CILIKUBE_CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	// --- 2. 创建并运行应用 ---
	// 现在使用集成了viper的配置加载
	application, err := app.New(configPath)
	if err != nil {
		slog.Error("应用初始化失败", "error", err)
		os.Exit(1)
	}
	slog.Info("应用初始化成功", "configPath", configPath)
	
	// 启动应用
	application.Run()
}
