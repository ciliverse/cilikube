package main

import (
	"log"

	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/api/v1/routes"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := configs.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化Kubernetes客户端
	k8sClient, err := k8s.NewClient(cfg.Kubernetes.Kubeconfig)
	if err != nil {
		log.Fatalf("创建Kubernetes客户端失败: %v", err)
	}

	// 初始化服务层
	podService := service.NewPodService(k8sClient.Clientset)

	// 初始化处理器
	podHandler := handlers.NewPodHandler(podService)

	// 创建Gin实例
	router := gin.Default()

	// 注册路由
	v1 := router.Group("/api/v1")
	routes.RegisterPodRoutes(v1, podHandler)

	// 启动服务器
	log.Printf("服务器启动，监听端口 %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
