package initialization

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/api/v1/routes"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/database"
	"github.com/ciliverse/cilikube/pkg/k8s" // 引入 k8s 包
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "k8s.io/client-go/kubernetes" // 通常由 k8s.Client 内部管理
	// "k8s.io/client-go/rest"       // 通常由 k8s.Client 内部管理
)

type AppServices struct {
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
	InstallerService     service.InstallerService // 非 K8s 服务
	AuthService          *service.AuthService
	ProxyService         *service.ProxyService // ProxyService 可能仍需 rest.Config，但会动态获取
}

type AppHandlers struct {
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
	InstallerHandler     *handlers.InstallerHandler // 非 K8s 处理器
	AuthHandler          *handlers.AuthHandler
	ProxyHandler         *handlers.ProxyHandler
}

func InitializeServices(
	k8sClusterManager *k8s.ClusterManager, // 主要的集群客户端管理器
	initialK8sClient *k8s.Client, // 一个可选的、用于特定场景的初始/默认 k8s 客户端
	initialK8sClientAvailable bool,
	cfg *configs.Config,
) *AppServices {
	log.Println("初始化服务层...")
	services := &AppServices{}

	// 初始化非 Kubernetes 依赖的服务 (例如 InstallerService, AuthService)
	services.InstallerService = service.NewInstallerService(cfg)
	log.Println("Installer 服务初始化完成。")

	if cfg.Database.Enabled {
		if database.DB != nil { // 确保数据库已连接
			// services.AuthService = service.NewAuthService(database.DB, cfg) // 示例，根据你的AuthService构造函数调整
			log.Println("AuthService (如果依赖数据库) 已初始化或准备就绪。")
		} else {
			log.Println("警告: 数据库已启用但连接失败，依赖数据库的 AuthService 可能无法正常工作。")
		}
	} else {
		// services.AuthService = service.NewAuthService(nil, cfg) // 如果AuthService可以无数据库运行
		log.Println("数据库未启用，AuthService (如果依赖数据库) 将以受限模式运行或不运行。")
	}

	log.Println("准备初始化 Kubernetes 相关服务...")

	// 检查是否有任何 K8s 集群可用（通过 ClusterManager）
	availableClusters := k8sClusterManager.GetAvailableClientNames()
	if len(availableClusters) > 0 {
		log.Printf("检测到 %d 个可用的 Kubernetes 集群。Kubernetes 相关服务将被实例化。", len(availableClusters))

		// 示例：PodService。假设 NewPodService() 现在不需要参数或只需要非 K8s 配置。
		services.PodService = service.NewPodService() // 构造函数将在第4步修改
		services.DeploymentService = service.NewDeploymentService()
		services.DaemonSetService = service.NewDaemonSetService()
		// services.ServiceService = service.NewServiceService()
		// services.IngressService = service.NewIngressService()
		// services.NetworkPolicyService = service.NewNetworkPolicyService()
		// services.ConfigMapService = service.NewConfigMapService()
		// services.SecretService = service.NewSecretService()
		services.PVCService = service.NewPVCService()
		services.PVService = service.NewPVService()
		services.StatefulSetService = service.NewStatefulSetService()
		services.NodeService = service.NewNodeService()
		services.NamespaceService = service.NewNamespaceService()
		// services.SummaryService = service.NewSummaryService()
		// services.EventsService = service.NewEventsService()
		// services.RbacService = service.NewRbacService()

		if initialK8sClientAvailable && initialK8sClient.Config != nil {
			log.Printf("为 ProxyService 使用初始客户端的 rest.Config (来自集群: %s)", initialK8sClient.Config.Host)
			services.ProxyService = service.NewProxyService(initialK8sClient.Config) // 假设 ProxyService 只需 config
		} else {

			log.Println("警告: 初始 Kubernetes 客户端或其配置不可用，ProxyService 可能无法按预期初始化或将以受限模式运行。")

		}
		log.Println("核心 Kubernetes 服务结构已实例化。它们将在请求时动态使用特定集群的客户端。")

	} else {
		log.Println("警告: 没有可用的 Kubernetes 集群 (通过 ClusterManager 检测)。所有 Kubernetes 相关服务将不可用或以空操作模式运行。")

	}

	log.Println("服务层初始化尝试完成。")
	return services
}

func InitializeHandlers(
	services *AppServices,
	k8sClusterManager *k8s.ClusterManager, // 引入 ClusterManager
	cfg *configs.Config, // 如果有处理器需要直接访问配置
) *AppHandlers {
	log.Println("初始化处理器层...")
	appHandlers := &AppHandlers{}

	// 初始化非 Kubernetes 依赖的处理器
	if services.InstallerService != nil {
		appHandlers.InstallerHandler = handlers.NewInstallerHandler(services.InstallerService)
		log.Println("Installer 处理器初始化完成。")
	}
	// if services.AuthService != nil {
	// appHandlers.AuthHandler = handlers.NewAuthHandler(services.AuthService, cfg, k8sClusterManager) // AuthService 可能也需要 clusterManager
	// log.Println("Auth 处理器初始化完成。")
	// }

	// 初始化 Kubernetes 相关的处理器。

	// (NewXxxHandler 构造函数签名将在后续步骤中修改 Handler 代码时调整)
	if services.PodService != nil { // 检查服务是否已实例化
		appHandlers.PodHandler = handlers.NewPodHandler(services.PodService, k8sClusterManager)
	}
	if services.DeploymentService != nil {
		appHandlers.DeploymentHandler = handlers.NewDeploymentHandler(services.DeploymentService, k8sClusterManager)
	}
	if services.DaemonSetService != nil {
		appHandlers.DaemonSetHandler = handlers.NewDaemonSetHandler(services.DaemonSetService, k8sClusterManager)
	}
	// if services.ServiceService != nil {
	// 	appHandlers.ServiceHandler = handlers.NewServiceHandler(services.ServiceService, k8sClusterManager)
	// }
	// if services.IngressService != nil {
	// 	appHandlers.IngressHandler = handlers.NewIngressHandler(services.IngressService, k8sClusterManager)
	// }
	// if services.NetworkPolicyService != nil {
	// 	appHandlers.NetworkPolicyHandler = handlers.NewNetworkPolicyHandler(services.NetworkPolicyService, k8sClusterManager)
	// }
	// if services.ConfigMapService != nil {
	// 	appHandlers.ConfigMapHandler = handlers.NewConfigMapHandler(services.ConfigMapService, k8sClusterManager)
	// }
	// if services.SecretService != nil {
	// 	appHandlers.SecretHandler = handlers.NewSecretHandler(services.SecretService, k8sClusterManager)
	// }
	if services.PVCService != nil {
		appHandlers.PVCHandler = handlers.NewPVCHandler(services.PVCService, k8sClusterManager)
	}
	if services.PVService != nil {
		appHandlers.PVHandler = handlers.NewPVHandler(services.PVService, k8sClusterManager)
	}
	if services.StatefulSetService != nil {
		appHandlers.StatefulSetHandler = handlers.NewStatefulSetHandler(services.StatefulSetService, k8sClusterManager)
	}
	if services.NodeService != nil {
		appHandlers.NodeHandler = handlers.NewNodeHandler(services.NodeService, k8sClusterManager)
	}
	if services.NamespaceService != nil {
		appHandlers.NamespaceHandler = handlers.NewNamespaceHandler(services.NamespaceService, k8sClusterManager)
	}
	// if services.SummaryService != nil {
	// 	appHandlers.SummaryHandler = handlers.NewSummaryHandler(services.SummaryService, k8sClusterManager)
	// }
	// if services.EventsService != nil {
	// 	appHandlers.EventsHandler = handlers.NewEventsHandler(services.EventsService, k8sClusterManager)
	// }
	// if services.RbacService != nil {
	// 	appHandlers.RbacHandler = handlers.NewRbacHandler(services.RbacService, k8sClusterManager)
	// }
	// if services.ProxyService != nil {
	// 	// ProxyHandler 的构造函数也需要 ClusterManager，因为它需要动态选择目标集群的 rest.Config
	// 	appHandlers.ProxyHandler = handlers.NewProxyHandler(services.ProxyService, k8sClusterManager)
	// }

	log.Println("处理器层初始化尝试完成 (部分可能因对应服务未初始化而被跳过)。")
	return appHandlers
}

// SetupRouter 配置 Gin 路由器。
// anyK8sAvailable 指示是否有任何 K8s 集群成功初始化
// k8sClusterManager 虽然在此函数中不直接使用，但它是初始化 Handlers 所需的，而 Handlers 由此函数使用。
func SetupRouter(
	cfg *configs.Config,
	appHandlers *AppHandlers,
	anyK8sAvailable bool, // 来自 main.go 的 anyK8sClientActuallyAvailable
	e *casbin.Enforcer,
	k8sClusterManager *k8s.ClusterManager, // 虽然 SetupRouter 不直接用，但传递性依赖可能需要
) *gin.Engine {
	log.Println("设置 Gin 路由器...")
	// gin.SetMode(cfg.Server.Mode) // 根据配置设置模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New() // 使用 gin.New() 以便自定义日志和恢复中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS 中间件配置 (与原配置保持一致)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Cluster-Name"}, // 考虑添加 X-Cluster-Name 如果用请求头传递集群
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	log.Println("应用 CORS 中间件。")

	// --- 健康检查路由 ---
	router.GET("/healthz", func(c *gin.Context) {
		healthStatus := gin.H{"status": "ok", "timestamp": time.Now().UTC()}
		// anyK8sAvailable 反映的是是否有任何集群在启动时连接成功
		if anyK8sAvailable {
			healthStatus["kubernetes_connectivity"] = "至少一个集群在启动时连接成功"
			healthStatus["available_clusters_at_startup"] = k8sClusterManager.GetAvailableClientNames()
		} else {
			healthStatus["kubernetes_connectivity"] = "没有集群在启动时连接成功 (Kubernetes 功能可能完全不可用)"
		}
		// 可以考虑在这里添加数据库连接状态的检查
		if cfg.Database.Enabled {
			if database.DB != nil {
				sqlDB, err := database.DB.DB() // 正确接收两个返回值
				if err != nil {
					// 获取底层 sql.DB 实例时发生错误
					log.Printf("警告: /healthz 无法获取数据库实例: %v", err)
					healthStatus["database"] = "error (获取DB实例失败)"
				} else if sqlDB == nil {
					// 理论上，如果 database.DB 不是 nil 且 err 是 nil，sqlDB 不应该为 nil，但作为健壮性检查
					healthStatus["database"] = "error (DB实例为nil)"
				} else {
					// 现在我们有 sqlDB，可以安全地调用 Ping()
					if pingErr := sqlDB.Ping(); pingErr == nil {
						healthStatus["database"] = "connected"
					} else {
						log.Printf("警告: /healthz 数据库ping失败: %v", pingErr)
						healthStatus["database"] = "disconnected (ping失败)"
					}
				}
			} else {
				// database.DB 本身就是 nil，意味着数据库未初始化或连接失败
				healthStatus["database"] = "disconnected (未初始化)"
			}
		} else {
			healthStatus["database"] = "not_enabled (配置中未启用)"
		}
		c.JSON(http.StatusOK, healthStatus)
	})

	// --- API v1 路由组 ---
	apiV1 := router.Group("/api/v1")
	{
		// --- 非 Kubernetes 依赖的路由 (例如 Auth, Installer) ---
		// Auth 路由 (例如登录)
		// routes.RegisterAuthRoutes(apiV1, appHandlers.AuthHandler)

		// Installer 路由
		if appHandlers.InstallerHandler != nil {
			routes.RegisterInstallerRoutes(apiV1, appHandlers.InstallerHandler)
			log.Println("注册 Installer 相关路由。")
		} else {
			log.Println("InstallerHandler 未初始化，跳过其路由注册。")
		}

		// --- 中间件 ---
		// JWT 中间件 (如果需要) - 应用于需要认证的路由组
		// authenticatedRoutes := apiV1.Group("") // 或者特定的子组
		// authenticatedRoutes.Use(auth.JWTAuthMiddleware(cfg.JWT.SecretKey)) // 假设的中间件

		// Casbin RBAC 中间件
		if e != nil {
			log.Println("为 /api/v1 应用 Casbin RBAC 中间件 (忽略 /auth/login)...")
			// 注意：CasbinMiddleware 需要能正确处理新的带 :cluster_name 的路径
			// 或者应用到更细分的路由组上。
			// IgnorePath 需要确保与实际登录路径匹配。
			// 如果所有受保护的 K8s 操作都在 /clusters/:cluster_name/ 下，
			// Casbin 策略也需要适配这种路径格式。
			apiV1.Use(auth.NewCasbinBuilder().
				IgnorePath("/api/v1/auth/login"). // 假设登录路径是这个
				// IgnorePathPrefixes("/api/v1/install") // 如果安装器路径也不需要认证
				CasbinMiddleware(e))
		} else {
			log.Println("Casbin 未初始化，跳过 RBAC 中间件的注册。")
		}

		// --- Kubernetes 相关路由 (现在需要包含集群名称) ---
		// 创建一个新的子路由组，用于处理所有针对特定集群的操作
		// URL 结构: /api/v1/clusters/{cluster_name}/resource...
		clusterSpecificRoutes := apiV1.Group("/clusters/:cluster_name")
		{
			// 在这里注册所有需要指定集群的 Kubernetes 资源路由
			// 路由注册函数 (如 RegisterPodRoutes) 仍然接收原始的 appHandlers.*Handler 实例
			// Handler 内部会使用 cluster_name 参数来获取正确的客户端
			if anyK8sAvailable { // 仅当至少一个K8s集群可能可用时才注册这些路由
				log.Println("准备注册特定集群的 Kubernetes API 路由...")
				if appHandlers.PodHandler != nil {
					routes.RegisterPodRoutes(clusterSpecificRoutes, appHandlers.PodHandler)
				}
				if appHandlers.DeploymentHandler != nil {
					routes.RegisterDeploymentRoutes(clusterSpecificRoutes, appHandlers.DeploymentHandler)
				}
				if appHandlers.DaemonSetHandler != nil {
					routes.RegisterDaemonSetRoutes(clusterSpecificRoutes, appHandlers.DaemonSetHandler)
				}
				if appHandlers.ServiceHandler != nil {
					routes.RegisterServiceRoutes(clusterSpecificRoutes, appHandlers.ServiceHandler)
				}
				if appHandlers.IngressHandler != nil {
					routes.RegisterIngressRoutes(clusterSpecificRoutes, appHandlers.IngressHandler)
				}
				if appHandlers.NetworkPolicyHandler != nil {
					routes.RegisterNetworkPolicyRoutes(clusterSpecificRoutes, appHandlers.NetworkPolicyHandler)
				}
				if appHandlers.ConfigMapHandler != nil {
					routes.RegisterConfigMapRoutes(clusterSpecificRoutes, appHandlers.ConfigMapHandler)
				}
				if appHandlers.SecretHandler != nil {
					routes.RegisterSecretRoutes(clusterSpecificRoutes, appHandlers.SecretHandler)
				}
				if appHandlers.PVCHandler != nil {
					routes.RegisterPVCRoutes(clusterSpecificRoutes, appHandlers.PVCHandler)
				}
				if appHandlers.PVHandler != nil {
					routes.RegisterPVRoutes(clusterSpecificRoutes, appHandlers.PVHandler)
				}
				if appHandlers.StatefulSetHandler != nil {
					routes.RegisterStatefulSetRoutes(clusterSpecificRoutes, appHandlers.StatefulSetHandler)
				}
				if appHandlers.NodeHandler != nil {
					// 例如：RegisterNodeRoutes 现在应该操作在 clusterSpecificRoutes 这个 group 上
					// routes.RegisterNodeRoutes(router *gin.RouterGroup, handler *handlers.NodeHandler)
					routes.RegisterNodeRoutes(clusterSpecificRoutes, appHandlers.NodeHandler)
				}
				if appHandlers.NamespaceHandler != nil {
					routes.RegisterNamespaceRoutes(clusterSpecificRoutes, appHandlers.NamespaceHandler)
				}
				if appHandlers.SummaryHandler != nil {
					routes.RegisterSummaryRoutes(clusterSpecificRoutes, appHandlers.SummaryHandler)
				}
				if appHandlers.EventsHandler != nil {
					routes.RegisterEventsRoutes(clusterSpecificRoutes, appHandlers.EventsHandler)
				}
				if appHandlers.RbacHandler != nil {
					routes.RegisterRbacRoutes(clusterSpecificRoutes, appHandlers.RbacHandler)
				}
				// Kubernetes Proxy 路由也需要适配集群上下文
				if appHandlers.ProxyHandler != nil {
					// KubernetesProxyRoutes 应该在 clusterSpecificRoutes 下注册
					// 例如: /api/v1/clusters/{cluster_name}/proxy/...
					routes.KubernetesProxyRoutes(clusterSpecificRoutes, appHandlers.ProxyHandler)
				}
				log.Println("特定集群的 Kubernetes API 路由注册尝试完成。")
			} else {
				log.Println("由于没有可用的 Kubernetes 集群，跳过特定集群的 Kubernetes API 路由注册。")
				// 可以选择性地为 /api/v1/clusters/:cluster_name 返回一个错误信息
				clusterSpecificRoutes.GET("/*any", func(c *gin.Context) {
					clusterName := c.Param("cluster_name")
					c.JSON(http.StatusServiceUnavailable, gin.H{
						"error":       fmt.Sprintf("Kubernetes cluster '%s' is targeted, but no Kubernetes services are available globally.", clusterName),
						"message":     "No Kubernetes clusters were successfully initialized at startup.",
						"clusterName": clusterName,
					})
				})
			}
		}
	}
	log.Println("API 路由注册完成。")
	return router
}

// Cleanup 函数 (如果有的话，保持不变)
func Cleanup() {
	if configs.GlobalConfig != nil && configs.GlobalConfig.Database.Enabled {
		if err := database.CloseDatabase(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		} else {
			log.Println("数据库连接已关闭。")
		}
	}
	// 这里也可以添加关闭 ClusterManager 中所有客户端的逻辑（如果它们维护了需要显式关闭的连接或 watch）
	// 但通常 client-go 的客户端不需要显式关闭。
}
