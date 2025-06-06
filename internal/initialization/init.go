package initialization

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/api/v1/routes"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AppServices 结构体需要包含新增的 ClusterService。
type AppServices struct {
	ClusterService       *service.ClusterService
	PodService           *service.PodService
	DeploymentService    *service.DeploymentService
	DaemonSetService     *service.DaemonSetService
	ServiceService       *service.ServiceService
	IngressService       *service.IngressService
	NetworkPolicyService *service.NetworkPolicyService
	ConfigMapService     *service.ConfigMapService
	SecretService        *service.SecretService
	PVCService           *service.PVCService
	PVService            *service.PVService
	StatefulSetService   *service.StatefulSetService
	NodeService          *service.NodeService
	NamespaceService     *service.NamespaceService
	SummaryService       *service.SummaryService
	EventsService        *service.EventsService
	RbacService          *service.RbacService
	InstallerService     service.InstallerService
	AuthService          *service.AuthService
	ProxyService         *service.ProxyService
}

// AppHandlers 结构体需要包含新增的 ClusterHandler。
type AppHandlers struct {
	ClusterHandler       *handlers.ClusterHandler
	PodHandler           *handlers.PodHandler
	DeploymentHandler    *handlers.DeploymentHandler
	DaemonSetHandler     *handlers.DaemonSetHandler
	ServiceHandler       *handlers.ServiceHandler
	IngressHandler       *handlers.IngressHandler
	NetworkPolicyHandler *handlers.NetworkPolicyHandler
	ConfigMapHandler     *handlers.ConfigMapHandler
	SecretHandler        *handlers.SecretHandler
	PVCHandler           *handlers.PVCHandler
	PVHandler            *handlers.PVHandler
	StatefulSetHandler   *handlers.StatefulSetHandler
	NodeHandler          *handlers.NodeHandler
	NamespaceHandler     *handlers.NamespaceHandler
	SummaryHandler       *handlers.SummaryHandler
	EventsHandler        *handlers.EventsHandler
	RbacHandler          *handlers.RbacHandler
	InstallerHandler     *handlers.InstallerHandler
	AuthHandler          *handlers.AuthHandler
	ProxyHandler         *handlers.ProxyHandler
}

// InitializeServices 的职责是创建所有 Service 层的实例。
func InitializeServices(k8sManager *k8s.ClusterManager, cfg *configs.Config) *AppServices {
	log.Println("正在初始化服务层...")

	return &AppServices{
		// 新增的集群管理服务
		ClusterService: service.NewClusterService(k8sManager),

		// 现有的 Kubernetes 资源服务，它们的构造函数现在是无参的
		PodService:           service.NewPodService(),
		DeploymentService:    service.NewDeploymentService(),
		DaemonSetService:     service.NewDaemonSetService(),
		ServiceService:       service.NewServiceService(),
		IngressService:       service.NewIngressService(),
		NetworkPolicyService: service.NewNetworkPolicyService(),
		ConfigMapService:     service.NewConfigMapService(),
		SecretService:        service.NewSecretService(),
		PVCService:           service.NewPVCService(),
		PVService:            service.NewPVService(),
		StatefulSetService:   service.NewStatefulSetService(),
		NodeService:          service.NewNodeService(),
		NamespaceService:     service.NewNamespaceService(),
		// SummaryService:       service.NewSummaryService(),
		// EventsService:        service.NewEventsService(),
		// RbacService:          service.NewRbacService(),
		// ProxyService:         service.NewProxyService(k8sManager), // ProxyService也需要manager来动态获取配置

		// 非 K8s 服务
		InstallerService: service.NewInstallerService(cfg),
		// AuthService:       service.NewAuthService(), // 根据您的实现调整
	}
}

// InitializeHandlers 的职责是创建所有 Handler 层的实例。
func InitializeHandlers(services *AppServices, k8sManager *k8s.ClusterManager) *AppHandlers {
	log.Println("正在初始化处理器层...")

	return &AppHandlers{
		// 新增的集群管理处理器
		ClusterHandler: handlers.NewClusterHandler(services.ClusterService),

		// 现有的 Kubernetes 资源处理器，它们的构造函数现在需要传入 k8sManager
		PodHandler:           handlers.NewPodHandler(services.PodService, k8sManager),
		DeploymentHandler:    handlers.NewDeploymentHandler(services.DeploymentService, k8sManager),
		DaemonSetHandler:     handlers.NewDaemonSetHandler(services.DaemonSetService, k8sManager),
		ServiceHandler:       handlers.NewServiceHandler(services.ServiceService, k8sManager),
		IngressHandler:       handlers.NewIngressHandler(services.IngressService, k8sManager),
		NetworkPolicyHandler: handlers.NewNetworkPolicyHandler(services.NetworkPolicyService, k8sManager),
		ConfigMapHandler:     handlers.NewConfigMapHandler(services.ConfigMapService, k8sManager),
		SecretHandler:        handlers.NewSecretHandler(services.SecretService, k8sManager),
		PVCHandler:           handlers.NewPVCHandler(services.PVCService, k8sManager),
		PVHandler:            handlers.NewPVHandler(services.PVService, k8sManager),
		StatefulSetHandler:   handlers.NewStatefulSetHandler(services.StatefulSetService, k8sManager),
		NodeHandler:          handlers.NewNodeHandler(services.NodeService, k8sManager),
		NamespaceHandler:     handlers.NewNamespaceHandler(services.NamespaceService, k8sManager),
		// SummaryHandler:       handlers.NewSummaryHandler(services.SummaryService, k8sManager),
		// EventsHandler:        handlers.NewEventsHandler(services.EventsService, k8sManager),
		// RbacHandler:          handlers.NewRbacHandler(services.RbacService, k8sManager),
		// ProxyHandler:         handlers.NewProxyHandler(services.ProxyService, k8sManager),

		// 非 K8s 处理器
		// InstallerHandler: handlers.NewInstallerHandler(&services.InstallerService),
		// AuthHandler:       handlers.NewAuthHandler(services.AuthService), // 根据您的实现调整
	}
}

// SetupRouter 配置 Gin 引擎和所有 API 路由。
func SetupRouter(cfg *configs.Config, appHandlers *AppHandlers, e *casbin.Enforcer) *gin.Engine {
	log.Println("正在设置 Gin 路由器...")
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	apiV1 := router.Group("/api/v1")
	{
		// 注册非 K8s 资源和集群管理相关的路由
		routes.RegisterInstallerRoutes(apiV1, appHandlers.InstallerHandler)
		routes.RegisterClusterRoutes(apiV1, appHandlers.ClusterHandler)
		// routes.RegisterAuthRoutes(apiV1, appHandlers.AuthHandler)

		// 在这里可以应用 JWT 和 Casbin 等中间件
		if e != nil {
			apiV1.Use(auth.NewCasbinBuilder().IgnorePath("/api/v1/auth/login").CasbinMiddleware(e))
		}

		// 创建一个新的子路由组，用于处理所有需要指定集群的资源操作
		// URL 结构: /api/v1/clusters/{cluster_name}/...
		clusterSpecificRoutes := apiV1.Group("/clusters/:cluster_name")
		{
			routes.RegisterPodRoutes(clusterSpecificRoutes, appHandlers.PodHandler)
			routes.RegisterDeploymentRoutes(clusterSpecificRoutes, appHandlers.DeploymentHandler)
			routes.RegisterDaemonSetRoutes(clusterSpecificRoutes, appHandlers.DaemonSetHandler)
			routes.RegisterServiceRoutes(clusterSpecificRoutes, appHandlers.ServiceHandler)
			routes.RegisterIngressRoutes(clusterSpecificRoutes, appHandlers.IngressHandler)
			routes.RegisterNetworkPolicyRoutes(clusterSpecificRoutes, appHandlers.NetworkPolicyHandler)
			routes.RegisterConfigMapRoutes(clusterSpecificRoutes, appHandlers.ConfigMapHandler)
			routes.RegisterSecretRoutes(clusterSpecificRoutes, appHandlers.SecretHandler)
			routes.RegisterPVCRoutes(clusterSpecificRoutes, appHandlers.PVCHandler)
			routes.RegisterPVRoutes(clusterSpecificRoutes, appHandlers.PVHandler)
			routes.RegisterStatefulSetRoutes(clusterSpecificRoutes, appHandlers.StatefulSetHandler)
			routes.RegisterNodeRoutes(clusterSpecificRoutes, appHandlers.NodeHandler)
			routes.RegisterNamespaceRoutes(clusterSpecificRoutes, appHandlers.NamespaceHandler)
			routes.RegisterSummaryRoutes(clusterSpecificRoutes, appHandlers.SummaryHandler)
			routes.RegisterEventsRoutes(clusterSpecificRoutes, appHandlers.EventsHandler)
			routes.RegisterRbacRoutes(clusterSpecificRoutes, appHandlers.RbacHandler)
			routes.KubernetesProxyRoutes(clusterSpecificRoutes, appHandlers.ProxyHandler)
		}
	}

	log.Println("API 路由注册完成。")
	return router
}

// StartServer 启动 HTTP 服务器（这是您提供的代码，保持不变）。
func StartServer(cfg *configs.Config, router http.Handler) {
	serverAddr := ":" + cfg.Server.Port
	version := getVersion()
	mode := os.Getenv("CILIKUBE_MODE")
	if mode == "" {
		mode = "development"
	}
	displayServerInfo(serverAddr, mode, version)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// displayServerInfo, getVersion, getLocalIP 等辅助函数也保持不变
func displayServerInfo(serverAddr, mode, version string) {
	color.Cyan("🚀 CiliKube Server is running!")
	color.Green("   ➜  Local:       http://127.0.0.1%s", serverAddr)
	color.Green("   ➜  Network:     http://%s%s", getLocalIP(), serverAddr)
	color.Yellow("  ➜  Mode:        %s", mode)
	color.Magenta("  ➜  Version:     %s", version)
}

func getVersion() string {
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return "v0.2.2" // default version
	}
	return string(data)
}

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
