package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/api/v1/routes"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s" // Your custom k8s client package
	// Assuming metrics package exists
	"github.com/ciliverse/cilikube/pkg/utils" // Assuming utils package exists
	"github.com/gin-gonic/gin"
	// Import necessary k8s types if needed directly here, though likely handled in service/handler
)

func main() {
	// --- Configuration Loading ---
	log.Println("加载配置文件...")
	cfg, err := configs.Load("config.yaml") // Consider making path configurable
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Println("配置文件加载成功。")

	// --- Kubernetes Client Initialization ---
	// Determine kubeconfig path: Use value from config, or empty string for in-cluster
	kubeconfigPath := cfg.Kubernetes.Kubeconfig
	// If config value is explicitly "in-cluster" or empty, treat as in-cluster attempt
	if kubeconfigPath == "in-cluster" {
		kubeconfigPath = "" // NewClient treats "" as in-cluster attempt
		log.Println("配置指定使用 in-cluster Kubernetes 配置。")
	} else if kubeconfigPath == "" {
		log.Println("配置文件中未指定 kubeconfig 路径，将尝试 in-cluster 配置。")
	} else {
		log.Printf("从配置文件获取 kubeconfig 路径: %s\n", kubeconfigPath)
	}

	var k8sClient *k8s.Client // Use your custom client type (now holds Config)
	var k8sAvailable bool     // Flag to track if K8s is usable

	// Initialize your custom Kubernetes client which should hold both clientset and config
	k8sClient, err = k8s.NewClient(kubeconfigPath) // Pass resolved path or ""
	if err != nil {
		// NewClient now returns more specific errors
		log.Printf("警告: 创建 Kubernetes 客户端失败: %v。Kubernetes 相关功能将不可用。", err)
		k8sAvailable = false
	} else {
		log.Println("Kubernetes 客户端创建成功。检查集群连接...")
		// Perform a connection check using your client's method
		if err := k8sClient.CheckConnection(); err != nil {
			log.Printf("警告: 无法连接 Kubernetes 集群: %v。Kubernetes 相关功能将不可用。", err)
			k8sAvailable = false
		} else {
			log.Println("成功连接到 Kubernetes 集群。")
			k8sAvailable = true
			// Log API Server URL for confirmation (optional)
			if k8sClient.Config != nil {
				log.Printf("连接到 API Server: %s", k8sClient.Config.Host)
			}
		}
	}

	// --- Service Layer Initialization ---
	log.Println("初始化服务层...")
	// Declare all service variables
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

	// Initialize services only if K8s client is available and connected
	if k8sAvailable {
		log.Println("Kubernetes 可用，初始化 Kubernetes 相关服务...")
		// ==================================================================
		// Corrected PodService Initialization: Access k8sClient.Config
		// ==================================================================
		// Pass BOTH k8sClient.Clientset AND k8sClient.Config
		if k8sClient.Config == nil {
			// This should not happen if k8sAvailable is true and NewClient worked correctly
			log.Fatalf("严重错误: k8sClient.Config 为 nil，无法初始化需要 rest.Config 的服务！")
		}
		podService = service.NewPodService(k8sClient.Clientset, k8sClient.Config) // CORRECTED

		// Initialize other services that likely only need the clientset
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
		log.Println("Kubernetes 相关服务初始化完成。")
	} else {
		log.Println("Kubernetes 不可用，跳过相关服务初始化。")
	}

	// Initialize services that don't depend on Kubernetes (always run)
	installerService := service.NewInstallerService(cfg) // Assuming this doesn't need k8s client
	log.Println("Installer 服务初始化完成。")

	// --- Handler Initialization ---
	log.Println("初始化处理器层...")
	// Declare all handler variables
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

	// Initialize handlers only if their corresponding services were initialized (i.e., k8sAvailable)
	if k8sAvailable {
		log.Println("Kubernetes 可用，初始化 Kubernetes 相关处理器...")
		// Ensure podService is not nil before creating handler
		if podService == nil {
			log.Fatalf("严重错误: podService 未初始化，无法创建 podHandler！")
		}
		podHandler = handlers.NewPodHandler(podService) // podService is now correctly initialized

		// ... initialize other handlers ...
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
		log.Println("Kubernetes 相关处理器初始化完成。")
	} else {
		log.Println("Kubernetes 不可用，跳过相关处理器初始化。")
	}

	// Initialize handlers for non-k8s services
	installerHandler := handlers.NewInstallerHandler(installerService)
	log.Println("Installer 处理器初始化完成。")

	// --- Gin Router Setup ---
	log.Println("设置 Gin 路由器...")
	// gin.SetMode(gin.ReleaseMode) // Uncomment for production
	router := gin.Default()

	// --- Middlewares ---
	// CORS Middleware (using your utility function)
	log.Println("应用 CORS 中间件...")
	router.Use(utils.Cors()) // Ensure utils.Cors() is correctly configured

	// Prometheus Metrics Middleware (if enabled)
	// if cfg.Server.EnableMetrics { // Example: Check config to enable metrics
	// 	log.Println("应用 Prometheus 中间件...")
	// 	router.Use(metrics.PromMiddleware()) // Ensure PromMiddleware exists and works
	// } else {
	// 	log.Println("Prometheus 指标未启用。")
	// }

	// --- Routes ---
	// Health Check Endpoint
	router.GET("/healthz", func(c *gin.Context) { // Use /healthz convention
		// More detailed health check could check K8s connection status too
		healthStatus := gin.H{"status": "ok", "timestamp": time.Now().UTC()}
		if k8sAvailable {
			healthStatus["kubernetes"] = "connected"
		} else {
			healthStatus["kubernetes"] = "disconnected"
		}
		c.JSON(http.StatusOK, healthStatus)
	})

	// Prometheus Metrics Endpoint (if enabled)
	// if cfg.Server.EnableMetrics {
	// 	log.Println("注册 /metrics 指标端点...")
	// 	router.GET("/metrics", metrics.PromHandler()) // Ensure PromHandler exists
	// }

	// API v1 Routes
	log.Println("注册 API v1 路由...")
	v1 := router.Group("/api/v1")
	{
		// Register K8s related routes only if handlers were initialized
		if k8sAvailable {
			log.Println("注册 Kubernetes API 路由...")
			// Ensure handlers are not nil before registering routes
			if podHandler != nil {
				routes.RegisterPodRoutes(v1, podHandler)
			}
			if deploymentHandler != nil {
				routes.RegisterDeploymentRoutes(v1, deploymentHandler)
			}
			if daemonsetHandler != nil {
				routes.RegisterDaemonSetRoutes(v1, daemonsetHandler)
			}
			if serviceHandler != nil {
				routes.RegisterServiceRoutes(v1, serviceHandler)
			}
			if ingressHandler != nil {
				routes.RegisterIngressRoutes(v1, ingressHandler)
			}
			if networkpolicyHandler != nil {
				routes.RegisterNetworkPolicyRoutes(v1, networkpolicyHandler)
			}
			if configMapHandler != nil {
				routes.RegisterConfigMapRoutes(v1, configMapHandler)
			}
			if secretHandler != nil {
				routes.RegisterSecretRoutes(v1, secretHandler)
			}
			if pvcHandler != nil {
				routes.RegisterPVCRoutes(v1, pvcHandler)
			}
			if pvHandler != nil {
				routes.RegisterPVRoutes(v1, pvHandler)
			}
			if statefulsetHandler != nil {
				routes.RegisterStatefulSetRoutes(v1, statefulsetHandler)
			}
			if nodeHandler != nil {
				routes.RegisterNodeRoutes(v1, nodeHandler)
			}
			if namespaceHandler != nil {
				routes.RegisterNamespaceRoutes(v1, namespaceHandler)
			}
			if summaryHandler != nil {
				routes.RegisterSummaryRoutes(v1, summaryHandler)
			}
		} else {
			log.Println("跳过 Kubernetes API 路由注册。")
			// Optional: Register dummy endpoints or status endpoints if K8s is unavailable
			v1.GET("/kubernetes-status", func(c *gin.Context) {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "Kubernetes service unavailable"})
			})
		}

		// Always register non-k8s routes
		log.Println("注册 Installer API 路由...")
		if installerHandler != nil {
			routes.RegisterInstallerRoutes(v1, installerHandler)
		} else {
			log.Println("警告: Installer handler 未初始化，无法注册相关路由。")
		}
	}
	log.Println("API 路由注册完成。")

	// --- Start Server ---
	serverAddr := ":" + cfg.Server.Port
	log.Printf("服务器准备启动，监听地址 %s", serverAddr)
	// Use ListenAndServe for better error handling feedback
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
		// Add timeouts for production hardening
		// ReadTimeout:  15 * time.Second,
		// WriteTimeout: 30 * time.Second, // Longer for potential uploads/long ops
		// IdleTimeout:  120 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("启动服务器失败: %v", err)
	}

	log.Println("服务器已关闭。") // Log when server stops gracefully
}
