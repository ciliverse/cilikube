// main.go (回顾)
package main

import (
	"log"

	"github.com/ciliverse/cilikube/internal/initialization"
	// "flag" // 如果 InitAll 内部处理 flag，main 中可以不直接导入
)

func main() {
	log.Println("应用程序主程序启动...")
	// flag.Parse() // 如果 InitAll 不处理 flag 解析，则应在此处解析一次

	cfg, err := initialization.InitAll()
	if err != nil {
		log.Fatalf("应用程序初始化失败: %v", err)
	}
	log.Println("核心组件初始化完成。")

	services := initialization.InitializeServices(cfg)
	appHandlers := initialization.InitializeHandlers(services) // ManagementHandler 也在其中
	log.Println("应用服务和处理器初始化完成。")

	router := initialization.SetupRouter(cfg, appHandlers) // ManagementRoutes 也在此注册
	log.Println("Gin 路由设置完成。")

	// 使用 initialization 包中的 StartServer (来自你的 start.go)
	initialization.StartServer(cfg, router) // router 是 *gin.Engine，它实现了 http.Handler

	log.Println("服务器已关闭。")
}
