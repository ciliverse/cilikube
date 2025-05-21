// internal/initialization/init.go
package initialization

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// API and internal packages
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/api/v1/routes"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s" // For ClientManager

	// External packages
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AppServices 和 AppHandlers 结构体定义保持不变 (如你所提供)

// AppServices holds all initialized services
type AppServices struct {
	PodService           *service.PodService
	DeploymentService    *service.DeploymentService
	DaemonSetService     *service.DaemonSetService
	ServiceService       *service.ServiceService // 之前叫 ServiceService
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
	InstallerService     service.InstallerService
	EventsService        *service.EventsService
	RbacService          *service.RbacService
}

// AppHandlers holds all initialized handlers
type AppHandlers struct {
	PodHandler           *handlers.PodHandler
	DeploymentHandler    *handlers.DeploymentHandler
	DaemonSetHandler     *handlers.DaemonSetHandler
	ServiceHandler       *handlers.ServiceHandler // 之前叫 ServiceHandler
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
	InstallerHandler     *handlers.InstallerHandler
	EventsHandler        *handlers.EventsHandler
	RbacHandler          *handlers.RbacHandler
	ManagementHandler    *handlers.ManagementHandler // 新增: 管理集群列表的处理器
}

// InitAll 是新的主初始化函数，由 main.go 调用
func InitAll() (*configs.Config, error) {
	log.Println("信息: 开始执行 InitAll 进行核心应用初始化...")

	// 1. 加载应用配置
	var configPath string
	// 确保 flag 只被解析一次。通常在程序入口处（main）解析。
	// 如果 main.go 已经调用 flag.Parse()，这里就不需要了。
	// 假设 flag 尚未解析 (或者将 flag 解析移到 main 中，这里只读取值)
	if !flag.Parsed() {
		// 从命令行参数 -config 读取配置文件路径
		flag.StringVar(&configPath, "config", "", "配置文件路径 (例如: configs/config.yaml)")
		flag.Parse() // 解析定义的 flag
	} else {
		// 如果 flag 已解析，尝试获取值 (这要求 main 已经定义了该 flag)
		// 或者，更好的方式是 main 解析 flag 后将 configPath 传递给 InitAll
		// 为简单起见，这里重复了 flag 定义，但理想情况下应该避免。
		// 另一种方式是去掉这里的 flag 解析，完全依赖 main 传递路径或环境变量。
		// configPath = flag.Lookup("config").Value.(flag.Getter).Get().(string) // 比较复杂
	}

	if configPath == "" {
		configPath = os.Getenv("CILIKUBE_CONFIG_PATH")
		if configPath != "" {
			log.Printf("信息: 使用环境变量 CILIKUBE_CONFIG_PATH 指定的配置文件路径: %s", configPath)
		}
	}
	if configPath == "" {
		configPath = "configs/config.yaml" // 默认路径
		log.Printf("信息: 未通过命令行或环境变量指定配置文件路径，使用默认路径: %s", configPath)
	}

	cfg, err := configs.Load(configPath) // configs.Load 会设置全局的 configs.GlobalAppConfig
	if err != nil {
		return nil, fmt.Errorf("加载应用配置 '%s' 失败: %w", configPath, err)
	}
	log.Printf("应用配置 '%s' 加载成功。", configPath)

	// 2. 初始化 Kubernetes ClientManager
	// configs.GlobalAppConfig.Kubernetes.Clusters 包含了所有集群的配置信息
	if err := k8s.InitClientManager(configs.GlobalAppConfig.Kubernetes.Clusters); err != nil {
		log.Printf("警告: 初始化 Kubernetes ClientManager 时遇到一个或多个错误: %v", err)
		// 即使这里有错误，也可能是部分集群初始化失败，ClientManager 实例本身应该还是创建了的。
	}

	// 检查 ClientManager 初始化后的状态
	manager, getMgrErr := k8s.GetClientManager()
	if getMgrErr != nil {
		log.Printf("严重错误: 获取 ClientManager 实例失败: %v。Kubernetes 功能可能完全不可用。", getMgrErr)
		// 这通常意味着 InitClientManager 内部的 once.Do 没有正确设置 globalClientManager
		// 或者是 InitClientManager 根本没被调用或 panic 了。
		// 根据应用需求，这可能是一个需要 panic 的致命错误。
		// return cfg, fmt.Errorf("获取 ClientManager 实例失败: %w", getMgrErr)
	} else {
		availableClusters := manager.ListClusterNames()
		if len(availableClusters) == 0 {
			log.Println("警告: ClientManager 初始化完成，但当前没有可用的 Kubernetes 集群客户端。请检查集群配置和连接性。")
		} else {
			log.Printf("Kubernetes ClientManager 初始化成功。当前可用集群: %v", availableClusters)
		}
	}

	// 3. 其他全局初始化 (例如：深化日志配置，数据库连接等)
	// log.SetupGlobalLogger(cfg.Server.LogLevel) // 假设你有这样的函数

	log.Println("核心应用初始化 (InitAll) 完成。")
	return cfg, nil // 返回加载的配置，main.go 中的 StartServer 可能需要它
}

// InitializeServices 初始化所有应用服务。
// Kubernetes 相关的服务构造函数不再接收 k8sClient 或 k8sAvailable。
// 它们内部将使用 k8s.GetClientManager() 来获取客户端。
func InitializeServices(cfg *configs.Config) *AppServices {
	log.Println("信息: 开始初始化应用服务层...")
	services := &AppServices{}

	// 1. 初始化非 Kubernetes 依赖的服务
	// 假设 NewInstallerService 的签名是 func(cfg *configs.Config) service.InstallerService
	services.InstallerService = service.NewInstallerService(cfg)
	log.Println("信息: Installer 服务初始化完成。")

	// 2. 初始化 Kubernetes 依赖的服务
	// 确保你的 Service 构造函数已修改为不接收 k8s client 参数。
	// 例如: func NewPodService() *service.PodService
	// services.PodService = service.NewPodService()
	// services.DeploymentService = service.NewDeploymentService()
	// services.DaemonSetService = service.NewDaemonSetService()
	// services.ServiceService = service.NewServiceService()
	// services.IngressService = service.NewIngressService()
	// services.NetworkPolicyService = service.NewNetworkPolicyService()
	// services.ConfigMapService = service.NewConfigMapService()
	// services.SecretService = service.NewSecretService()
	// services.PVCService = service.NewPVCService()
	// services.PVService = service.NewPVService()
	// services.StatefulSetService = service.NewStatefulSetService()
	services.NodeService = service.NewNodeService() // NodeService 将使用 ClientManager
	// services.NamespaceService = service.NewNamespaceService()
	// services.SummaryService = service.NewSummaryService()
	// services.EventsService = service.NewEventsService()
	// services.RbacService = service.NewRbacService()
	log.Println("信息: 所有 Kubernetes 相关服务已实例化 (它们将在运行时按需获取客户端)。")

	log.Println("信息: 应用服务层初始化完成。")
	return services
}

// InitializeHandlers 初始化所有应用处理器。
// 处理器依赖于已初始化的服务。
func InitializeHandlers(services *AppServices) *AppHandlers {
	log.Println("信息: 开始初始化 API 处理器层...")
	appHandlers := &AppHandlers{}

	// 非 K8s 处理器
	// 假设 InstallerService 是值类型或接口，其初始化不涉及指针检查
	appHandlers.InstallerHandler = handlers.NewInstallerHandler(services.InstallerService)
	log.Println("信息: Installer 处理器初始化完成。")

	// K8s 依赖的处理器
	// 确保 NewXxxHandler 接收对应的 Service 实例
	appHandlers.PodHandler = handlers.NewPodHandler(services.PodService)
	appHandlers.DeploymentHandler = handlers.NewDeploymentHandler(services.DeploymentService)
	appHandlers.DaemonSetHandler = handlers.NewDaemonSetHandler(services.DaemonSetService)
	appHandlers.ServiceHandler = handlers.NewServiceHandler(services.ServiceService)
	appHandlers.IngressHandler = handlers.NewIngressHandler(services.IngressService)
	appHandlers.NetworkPolicyHandler = handlers.NewNetworkPolicyHandler(services.NetworkPolicyService)
	appHandlers.ConfigMapHandler = handlers.NewConfigMapHandler(services.ConfigMapService)
	appHandlers.SecretHandler = handlers.NewSecretHandler(services.SecretService)
	appHandlers.PVCHandler = handlers.NewPVCHandler(services.PVCService)
	appHandlers.PVHandler = handlers.NewPVHandler(services.PVService)
	appHandlers.StatefulSetHandler = handlers.NewStatefulSetHandler(services.StatefulSetService)
	appHandlers.NodeHandler = handlers.NewNodeHandler(services.NodeService) // NodeHandler 将处理多集群请求
	appHandlers.NamespaceHandler = handlers.NewNamespaceHandler(services.NamespaceService)
	appHandlers.SummaryHandler = handlers.NewSummaryHandler(services.SummaryService)
	appHandlers.EventsHandler = handlers.NewEventsHandler(services.EventsService)
	appHandlers.RbacHandler = handlers.NewRbacHandler(services.RbacService)
	log.Println("信息: Kubernetes 相关处理器初始化完成。")

	// 新增: 初始化 ManagementHandler
	appHandlers.ManagementHandler = handlers.NewManagementHandler()
	log.Println("信息: Management 处理器初始化完成。")

	log.Println("信息: API 处理器层初始化完成。")
	return appHandlers
}

// SetupRouter 配置 Gin 路由器。
// 移除了 k8sAvailable 参数。路由总是注册，可用性在运行时处理。
func SetupRouter(cfg *configs.Config, appHandlers *AppHandlers) *gin.Engine {
	log.Println("信息: 开始设置 Gin 路由器...")

	// 根据配置设置 Gin 模式 (例如，从 cfg.Server.Mode 读取)
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
		log.Println("信息: Gin 运行在 Release 模式。")
	} else {
		gin.SetMode(gin.DebugMode) // 默认为 Debug 模式
		log.Println("信息: Gin 运行在 Debug 模式。")
	}

	router := gin.New()        // 使用 New() 以便更好地控制中间件
	router.Use(gin.Logger())   // Gin 的日志中间件 (可以替换为你自定义的日志中间件)
	router.Use(gin.Recovery()) // Gin 的 Panic 恢复中间件

	// CORS 中间件配置
	corsConfig := cors.Config{
		AllowAllOrigins:  true, // 生产环境建议配置具体的、受信任的源列表
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "X-Cluster-Name"}, // 如果你用 Header 传集群名
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))
	log.Println("信息: CORS 中间件已应用。")

	// 健康检查端点
	router.GET("/healthz", func(c *gin.Context) {
		status := gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		manager, err := k8s.GetClientManager() // 尝试获取 ClientManager
		if err == nil && manager != nil {
			availableClusters := manager.ListClusterNames()
			status["cluster_manager_status"] = "active"
			status["available_clusters"] = availableClusters
			if len(availableClusters) == 0 {
				status["kubernetes_notice"] = "ClientManager 已激活，但当前无可用或已配置的集群。"
			}
		} else {
			status["cluster_manager_status"] = "inactive_or_error"
			status["kubernetes_notice"] = fmt.Sprintf("无法获取 ClientManager: %v", err)
		}
		c.JSON(http.StatusOK, status)
	})

	// API v1 路由组
	apiV1 := router.Group("/api/v1")
	{
		log.Println("信息: 注册 API v1 路由...")

		// 1. 注册管理接口 (例如列出集群)
		if appHandlers.ManagementHandler != nil {
			routes.RegisterManagementRoutes(apiV1, appHandlers.ManagementHandler)
			log.Println("信息: 管理相关路由已注册。")
		} else {
			log.Println("警告: ManagementHandler 未初始化，跳过管理路由注册。")
		}

		// 2. 注册特定资源的路由 (这些路由现在应该支持 /clusters/{clusterName}/... 模式)
		// 确保你的 routes.RegisterXxxRoutes 函数内部适配了新的 Handler 和 Service 逻辑
		if appHandlers.NodeHandler != nil {
			routes.RegisterNodeRoutes(apiV1, appHandlers.NodeHandler) // Node 路由将包含集群选择
			log.Println("信息: Node 相关路由已注册。")
		} else {
			log.Println("警告: NodeHandler 未初始化，跳过 Node 路由注册。")
		}

		// 为其他所有 Kubernetes 资源类型注册路由
		// 假设它们的 RegisterXxxRoutes 函数也已适配多集群路径
		if appHandlers.PodHandler != nil {
			routes.RegisterPodRoutes(apiV1, appHandlers.PodHandler)
		}
		if appHandlers.DeploymentHandler != nil {
			routes.RegisterDeploymentRoutes(apiV1, appHandlers.DeploymentHandler)
		}
		if appHandlers.DaemonSetHandler != nil {
			routes.RegisterDaemonSetRoutes(apiV1, appHandlers.DaemonSetHandler)
		}
		if appHandlers.ServiceHandler != nil {
			routes.RegisterServiceRoutes(apiV1, appHandlers.ServiceHandler)
		}
		if appHandlers.IngressHandler != nil {
			routes.RegisterIngressRoutes(apiV1, appHandlers.IngressHandler)
		}
		if appHandlers.NetworkPolicyHandler != nil {
			routes.RegisterNetworkPolicyRoutes(apiV1, appHandlers.NetworkPolicyHandler)
		}
		if appHandlers.ConfigMapHandler != nil {
			routes.RegisterConfigMapRoutes(apiV1, appHandlers.ConfigMapHandler)
		}
		if appHandlers.SecretHandler != nil {
			routes.RegisterSecretRoutes(apiV1, appHandlers.SecretHandler)
		}
		if appHandlers.PVCHandler != nil {
			routes.RegisterPVCRoutes(apiV1, appHandlers.PVCHandler)
		}
		if appHandlers.PVHandler != nil {
			routes.RegisterPVRoutes(apiV1, appHandlers.PVHandler)
		}
		if appHandlers.StatefulSetHandler != nil {
			routes.RegisterStatefulSetRoutes(apiV1, appHandlers.StatefulSetHandler)
		}
		if appHandlers.NamespaceHandler != nil {
			routes.RegisterNamespaceRoutes(apiV1, appHandlers.NamespaceHandler)
		}
		if appHandlers.SummaryHandler != nil {
			routes.RegisterSummaryRoutes(apiV1, appHandlers.SummaryHandler)
		}
		if appHandlers.EventsHandler != nil {
			routes.RegisterEventsRoutes(apiV1, appHandlers.EventsHandler)
		}
		if appHandlers.RbacHandler != nil {
			routes.RegisterRbacRoutes(apiV1, appHandlers.RbacHandler)
		}
		log.Println("信息: Kubernetes 资源相关路由已尝试注册。")

		// 3. 注册非 Kubernetes 相关的路由 (例如 Installer)
		if appHandlers.InstallerHandler != nil {
			routes.RegisterInstallerRoutes(apiV1, appHandlers.InstallerHandler)
			log.Println("信息: Installer 相关路由已注册。")
		} else {
			log.Println("警告: InstallerHandler 未初始化，跳过 Installer 路由注册。")
		}
	}
	log.Println("信息: Gin 路由器设置完成。")
	return router // 返回配置好的 Gin 引擎
}

// StartServer 函数 (从你的 start.go 文件移入，或者保持分离并通过 main.go 调用)
// func StartServer(cfg *configs.Config, router http.Handler) { ... } // 已在 start.go 定义，无需重复
// 注意：你的 start.go 中的 StartServer 期望 http.Handler，而 gin.Engine 实现了 http.Handler 接口。
