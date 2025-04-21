package main

import (
    "log"
    "net/http"
    "os"
    "time"

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

    // 检查 kubeconfig 是否存在
    kubeconfigPath := cfg.Kubernetes.Kubeconfig
    _, err = os.Stat(kubeconfigPath)
    k8sAvailable := err == nil

    var k8sClient *k8s.Client
    if k8sAvailable {
        // 初始化 Kubernetes 客户端
        k8sClient, err = k8s.NewClient(kubeconfigPath)
        if err != nil {
            log.Printf("警告: 创建 Kubernetes 客户端失败 (可能是可选的): %v", err)
            k8sAvailable = false
        } else {
            // 服务健康检查
            if err := k8sClient.CheckConnection(); err != nil {
                log.Printf("警告: 无法连接 Kubernetes 集群 (可能是可选的): %v", err)
                k8sAvailable = false
            }
        }
    } else {
        log.Println("Kubernetes 配置不存在，跳过 Kubernetes 相关服务初始化")
    }

    // 初始化服务层
    var podService *service.PodService
    var deploymentService *service.DeploymentService
    var daemonsetService *service.DaemonSetService
    var serviceService *service.ServiceService
    var ingressService *service.IngressService
    var networkpolicyService *service.NetworkPolicyService
    var configMapService *service.ConfigMapService
    var secretService *service.SecretService
    var pvcService *service.PVCService
    var pvService *service.PVService
    var statefulsetService *service.StatefulSetService
    var nodeService *service.NodeService
    var namespaceService *service.NamespaceService
    var summaryService *service.SummaryService

    if k8sAvailable {
        podService = service.NewPodService(k8sClient.Clientset)
        deploymentService = service.NewDeploymentService(k8sClient.Clientset)
        daemonsetService = service.NewDaemonSetService(k8sClient.Clientset)
        serviceService = service.NewServiceService(k8sClient.Clientset)
        ingressService = service.NewIngressService(k8sClient.Clientset)
        networkpolicyService = service.NewNetworkPolicyService(k8sClient.Clientset)
        configMapService = service.NewConfigMapService(k8sClient.Clientset)
        secretService = service.NewSecretService(k8sClient.Clientset)
        pvcService = service.NewPVCService(k8sClient.Clientset)
        pvService = service.NewPVService(k8sClient.Clientset)
        statefulsetService = service.NewStatefulSetService(k8sClient.Clientset)
        nodeService = service.NewNodeService(k8sClient.Clientset)
        namespaceService = service.NewNamespaceService(k8sClient.Clientset)
        summaryService = service.NewSummaryService(k8sClient.Clientset)
    }

    // 始终初始化安装程序服务
    installerService := service.NewInstallerService(cfg)

    // 初始化处理器
    var podHandler *handlers.PodHandler
    var deploymentHandler *handlers.DeploymentHandler
    var daemonsetHandler *handlers.DaemonSetHandler
    var serviceHandler *handlers.ServiceHandler
    var ingressHandler *handlers.IngressHandler
    var networkpolicyHandler *handlers.NetworkPolicyHandler
    var configMapHandler *handlers.ConfigMapHandler
    var secretHandler *handlers.SecretHandler
    var pvcHandler *handlers.PVCHandler
    var pvHandler *handlers.PVHandler
    var statefulsetHandler *handlers.StatefulSetHandler
    var nodeHandler *handlers.NodeHandler
    var namespaceHandler *handlers.NamespaceHandler
    var summaryHandler *handlers.SummaryHandler

    if k8sAvailable {
        podHandler = handlers.NewPodHandler(podService)
        deploymentHandler = handlers.NewDeploymentHandler(deploymentService)
        daemonsetHandler = handlers.NewDaemonSetHandler(daemonsetService)
        serviceHandler = handlers.NewServiceHandler(serviceService)
        ingressHandler = handlers.NewIngressHandler(ingressService)
        networkpolicyHandler = handlers.NewNetworkPolicyHandler(networkpolicyService)
        configMapHandler = handlers.NewConfigMapHandler(configMapService)
        secretHandler = handlers.NewSecretHandler(secretService)
        pvcHandler = handlers.NewPVCHandler(pvcService)
        pvHandler = handlers.NewPVHandler(pvService)
        statefulsetHandler = handlers.NewStatefulSetHandler(statefulsetService)
        nodeHandler = handlers.NewNodeHandler(nodeService)
        namespaceHandler = handlers.NewNamespaceHandler(namespaceService)
        summaryHandler = handlers.NewSummaryHandler(summaryService)
    }

    // 始终初始化安装程序处理器
    installerHandler := handlers.NewInstallerHandler(installerService)

    // 创建 Gin 实例
    router := gin.Default()
    router.Use(utils.Cors()) // 跨域请求处理

    // 添加健康检查路由
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().UTC()})
    })

    // 注册 API v1 路由
    v1 := router.Group("/api/v1")
    {
        if k8sAvailable {
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
        }
        // 始终注册安装程序路由
        routes.RegisterInstallerRoutes(v1, installerHandler)
    }

    // 监控路由
    router.Use(metrics.PromMiddleware())
    router.GET("/metrics", metrics.PromHandler())

    // 启动服务器
    log.Printf("服务器启动，监听端口 %s", cfg.Server.Port)
    if err := router.Run(":" + cfg.Server.Port); err != nil {
        log.Fatalf("启动服务器失败: %v", err)
    }
}