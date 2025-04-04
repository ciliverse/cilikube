package main

import (
	"log"

	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/api/v1/routes"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/metrics"
	"github.com/ciliverse/cilikube/pkg/utils"
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

	// 服务健康检查
	if err := k8sClient.CheckConnection(); err != nil {
		log.Fatalf("无法连接Kubernetes集群: %v", err)
	}
	// 初始化服务层
	podService := service.NewPodService(k8sClient.Clientset)
	deploymentService := service.NewDeploymentService(k8sClient.Clientset)
	daemonsetService := service.NewDaemonSetService(k8sClient.Clientset)
	serviceService := service.NewServiceService(k8sClient.Clientset)
	ingressService := service.NewIngressService(k8sClient.Clientset)
	networkpolicyService := service.NewNetworkPolicyService(k8sClient.Clientset)
	configMapService := service.NewConfigMapService(k8sClient.Clientset)
	secretService := service.NewSecretService(k8sClient.Clientset)
	pvcService := service.NewPVCService(k8sClient.Clientset)
	pvService := service.NewPVService(k8sClient.Clientset)
	statefulsetService := service.NewStatefulSetService(k8sClient.Clientset)
	nodeService := service.NewNodeService(k8sClient.Clientset)
	namespaceService := service.NewNamespaceService(k8sClient.Clientset)
	summaryService := service.NewSummaryService(k8sClient.Clientset)

	// 初始化处理器
	podHandler := handlers.NewPodHandler(podService)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService)
	daemonsetHandler := handlers.NewDaemonSetHandler(daemonsetService)
	serviceHandler := handlers.NewServiceHandler(serviceService)
	ingressHandler := handlers.NewIngressHandler(ingressService)
	networkpolicyHandler := handlers.NewNetworkPolicyHandler(networkpolicyService)
	configMapHandler := handlers.NewConfigMapHandler(configMapService)
	secretHandler := handlers.NewSecretHandler(secretService)
	pvcHandler := handlers.NewPVCHandler(pvcService)
	pvHandler := handlers.NewPVHandler(pvService)
	statefulsetHandler := handlers.NewStatefulSetHandler(statefulsetService)
	nodeHandler := handlers.NewNodeHandler(nodeService)
	namespaceHandler := handlers.NewNamespaceHandler(namespaceService)
	summaryHandler := handlers.NewSummaryHandler(summaryService)

	// 创建Gin实例
	router := gin.Default()
	router.Use(utils.Cors()) // 跨域请求处理

	// 注册路由
	v1 := router.Group("/api/v1")
	routes.RegisterPodRoutes(v1, podHandler)
	routes.RegisterDeploymentRoutes(v1, deploymentHandler)
	routes.RegisterDaemonSetRoutes(v1, daemonsetHandler)
	routes.RegisterServiceRoutes(v1, serviceHandler)
	routes.RegisterIngressRoutes(v1, ingressHandler)
	routes.RegisterNetworkPolicyRoutes(v1, networkpolicyHandler)
	routes.RegisterConfigMapRoutes(v1, configMapHandler)
	routes.RegisterSecretRoutes(v1, secretHandler)
	routes.RegisterPVCRoutes(v1, pvcHandler)
	routes.RegisterPVRoutes(v1, pvHandler)
	routes.RegisterStatefulSetRoutes(v1, statefulsetHandler)
	routes.RegisterNodeRoutes(v1, nodeHandler)
	routes.RegisterNamespaceRoutes(v1, namespaceHandler)
	routes.RegisterSummaryRoutes(v1, summaryHandler)

	// 监控路由
	router.Use(metrics.PromMiddleware())
	router.GET("/metrics", metrics.PromHandler())

	//增强日志
	// log.Init(false) // 开发模式
	// router.Use(log.GinLogger())

	// 启动服务器
	log.Printf("服务器启动，监听端口 %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// 自定义请求日志中间件
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("收到请求: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}
