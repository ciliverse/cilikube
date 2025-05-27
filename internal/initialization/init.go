package initialization

import (
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
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// Add any other necessary imports that were in main.go functions moved here
	// For k8sClient.Config (type *rest.Config), you might need "k8s.io/client-go/rest"
	// depending on how the k8s.Client struct is defined and used.
)

// AppConfig holds all initialized Repositories
// type AppRepository struct {
// 	AuthRepository *repository.AuthRepository
// }

// AppServices holds all initialized services
// Moved from main.go
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
	InstallerService     service.InstallerService // Non-k8s service
	AuthService          *service.AuthService     // auth service
}

// AppHandlers holds all initialized handlers
// Moved from main.go
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
	InstallerHandler     *handlers.InstallerHandler // Non-k8s handlers
	AuthHandler          *handlers.AuthHandler      // auth handler
}

// InitializeRepository initializes the database repository.
// func InitializeRepositories(DB *gorm.DB) *AppRepository {

// 	// 初始化 Repository
// 	appRepository := &AppRepository{}
// 	appRepository.AuthRepository = repository.NewAuthRepository(DB)

// 	return appRepository
// }

// InitializeServices initializes all application services.
// K8s-dependent services are only initialized if k8sAvailable is true.
// Moved from main.go
func InitializeServices(k8sClient *k8s.Client, k8sAvailable bool, cfg *configs.Config) *AppServices {
	log.Println("初始化服务层...")
	services := &AppServices{}

	// Initialize non-k8s services (always)
	services.InstallerService = service.NewInstallerService(cfg)
	log.Println("Installer 服务初始化完成。")

	if cfg.Database.Enabled {
		log.Println("数据库已启用，开始初始化...")
		if err := database.InitDatabase(); err != nil {
			log.Fatalf("初始化失败: 数据库连接失败: %v", err)
		}
		// 只有 InitDatabase 成功 (DB != nil) 才会继续
		if database.DB != nil {
			// log.Println("数据库连接成功。") // 这条日志现在只会在真正连接后打印
			if err := database.AutoMigrate(); err != nil {
				log.Fatalf("初始化失败: 数据库自动迁移失败: %v", err)
			}

		} else {
			// 这种情况理论上不应该发生，除非 InitDatabase 内部逻辑有误
			log.Println("警告: 数据库已启用但初始化失败，相关服务将无法使用。")
		}
	} else {
		log.Println("警告: 数据库未启用，相关服务将无法使用。")
	}
	// Initialize AuthService

	// --- Auth Initialization ---
	// - 这里初始化 AuthService
	// - 需要确保 database.DB 已经成功连接

	// Initialize K8s-dependent services (conditionally)
	if k8sAvailable && k8sClient != nil && k8sClient.Clientset != nil {
		log.Println("Kubernetes 可用，初始化 Kubernetes 相关服务...")
		// Services requiring rest.Config need special handling
		if k8sClient.Config == nil {
			log.Printf("警告: k8sClient.Config 为 nil。需要 rest.Config 的 Kubernetes 服务将无法完全初始化！")
			// Proceed, as many services only need Clientset
		}

		// Pass k8sClient.Clientset and potentially k8sClient.Config where needed
		services.PodService = service.NewPodService(k8sClient.Clientset, k8sClient.Config) // Assuming PodService needs Config
		services.DeploymentService = service.NewDeploymentService(k8sClient.Clientset)
		services.DaemonSetService = service.NewDaemonSetService(k8sClient.Clientset)
		services.ServiceService = service.NewServiceService(k8sClient.Clientset)
		services.IngressService = service.NewIngressService(k8sClient.Clientset)
		services.NetworkPolicyService = service.NewNetworkPolicyService(k8sClient.Clientset)
		services.ConfigMapService = service.NewConfigMapService(k8sClient.Clientset)
		services.SecretService = service.NewSecretService(k8sClient.Clientset)
		services.PVCService = service.NewPVCService(k8sClient.Clientset)
		services.PVService = service.NewPVService(k8sClient.Clientset)
		services.StatefulSetService = service.NewStatefulSetService(k8sClient.Clientset)
		services.NodeService = service.NewNodeService(k8sClient.Clientset)
		services.NamespaceService = service.NewNamespaceService(k8sClient.Clientset)
		services.SummaryService = service.NewSummaryService(k8sClient.Clientset)
		services.EventsService = service.NewEventsService(k8sClient.Clientset)
		services.RbacService = service.NewRbacService(k8sClient.Clientset)
		log.Println("Kubernetes 相关服务初始化完成。")
	} else {
		log.Println("Kubernetes 不可用，跳过相关服务初始化。")
		// K8s service pointers in 'services' struct remain nil
	}

	log.Println("服务初始化尝试完成。")
	return services
}

// InitializeHandlers initializes all application handlers.
// Handlers are only initialized if their corresponding service is available (non-nil).
// Moved from main.go
func InitializeHandlers(services *AppServices) *AppHandlers {
	log.Println("初始化处理器层...")
	appHandlers := &AppHandlers{}

	// Initialize non-k8s handlers (if service exists)
	if services.InstallerService != nil { // Should always be true if initializeServices ran
		appHandlers.InstallerHandler = handlers.NewInstallerHandler(services.InstallerService)
		log.Println("Installer 处理器初始化完成。")
	} else {
		log.Println("警告: Installer 服务未初始化，跳过 Installer 处理器初始化。")
	}

	// Initialize K8s-dependent handlers (conditionally based on service)
	// Check if the specific service pointer is non-nil
	if services.PodService != nil {
		appHandlers.PodHandler = handlers.NewPodHandler(services.PodService)
	}
	if services.DeploymentService != nil {
		appHandlers.DeploymentHandler = handlers.NewDeploymentHandler(services.DeploymentService)
	}
	if services.DaemonSetService != nil {
		appHandlers.DaemonSetHandler = handlers.NewDaemonSetHandler(services.DaemonSetService)
	}
	if services.ServiceService != nil {
		appHandlers.ServiceHandler = handlers.NewServiceHandler(services.ServiceService)
	}
	if services.IngressService != nil {
		appHandlers.IngressHandler = handlers.NewIngressHandler(services.IngressService)
	}
	if services.NetworkPolicyService != nil {
		appHandlers.NetworkPolicyHandler = handlers.NewNetworkPolicyHandler(services.NetworkPolicyService)
	}
	if services.ConfigMapService != nil {
		appHandlers.ConfigMapHandler = handlers.NewConfigMapHandler(services.ConfigMapService)
	}
	if services.SecretService != nil {
		appHandlers.SecretHandler = handlers.NewSecretHandler(services.SecretService)
	}
	if services.PVCService != nil {
		appHandlers.PVCHandler = handlers.NewPVCHandler(services.PVCService)
	}
	if services.PVService != nil {
		appHandlers.PVHandler = handlers.NewPVHandler(services.PVService)
	}
	if services.StatefulSetService != nil {
		appHandlers.StatefulSetHandler = handlers.NewStatefulSetHandler(services.StatefulSetService)
	}
	if services.NodeService != nil {
		appHandlers.NodeHandler = handlers.NewNodeHandler(services.NodeService)
	}
	if services.NamespaceService != nil {
		appHandlers.NamespaceHandler = handlers.NewNamespaceHandler(services.NamespaceService)
	}
	if services.SummaryService != nil {
		appHandlers.SummaryHandler = handlers.NewSummaryHandler(services.SummaryService)
	}
	if services.EventsService != nil {
		appHandlers.EventsHandler = handlers.NewEventsHandler(services.EventsService)
	}
	if services.RbacService != nil {
		appHandlers.RbacHandler = handlers.NewRbacHandler(services.RbacService)
	}

	log.Println("处理器初始化尝试完成 (部分可能因服务未初始化而跳过)。")
	return appHandlers
}

// SetupRouter configures the Gin router with middleware and routes.
// Moved from main.go
func SetupRouter(cfg *configs.Config, handlers *AppHandlers, k8sAvailable bool, e *casbin.Enforcer) *gin.Engine {
	log.Println("设置 Gin 路由器...")
	// gin.SetMode(gin.ReleaseMode) // Uncomment for production
	router := gin.Default()

	// 从配置或环境变量加载允许的源
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// --- Middlewares ---
	log.Println("应用 CORS 中间件...")
	//router.Use(utils.Cors(origins)) // Ensure utils.Cors() is correctly configured

	// Prometheus Metrics Middleware (if enabled) - Example
	// if cfg.Server.EnableMetrics {
	// 	log.Println("应用 Prometheus 中间件...")
	// 	// router.Use(metrics.PromMiddleware()) // Assuming metrics package exists
	// } else {
	// 	log.Println("Prometheus 指标未启用。")
	// }

	// --- Routes ---
	// Health Check Endpoint
	router.GET("/healthz", func(c *gin.Context) {
		healthStatus := gin.H{"status": "ok", "timestamp": time.Now().UTC()}
		if k8sAvailable {
			healthStatus["kubernetes"] = "connected"
		} else {
			healthStatus["kubernetes"] = "disconnected (features disabled)"
		}
		c.JSON(http.StatusOK, healthStatus)
	})

	// Prometheus Metrics Endpoint (if enabled) - Example
	// if cfg.Server.EnableMetrics {
	// 	log.Println("注册 /metrics 指标端点...")
	// 	// router.GET("/metrics", metrics.PromHandler()) // Assuming metrics package exists
	// }

	// API v1 Routes Group
	log.Println("注册 API v1 路由...")
	v1 := router.Group("/api/v1")
	{
		// Register non-k8s routes first (if handlers are initialized)
		// --- Auth Routes ---
		// 首先先从环境变量中获取 secret key
		// secretKey := os.Getenv("CILIKUBE_JWT_SECRET")

		// --- Protected Routes ---
		// 将与 k8s 相关操作放进 v1 中，将 login
		// 或者其他不需要受保护的路由 Ignore (这里需要注意的事一定要用v1，否则不会生效)
		log.Println("初始化 JWT 中间件...")
		// 添加后会用 JWT 中间件来保护 API 路由
		// v1.Use(auth.JWTAuthMiddleware())

		// --- RBAC Middleware ---
		// 注册 RBAC 中间件，这个中间件会在每个受保护的路由上应用
		// log.Println("初始化 RBAC 中间件...")
		if e != nil {
			log.Println("应用 RBAC 中间件...")
			v1.Use(auth.NewCasbinBuilder().
				IgnorePath("/api/v1/auth/login"). // 忽略登录路径
				CasbinMiddleware(e))              // 应用 Casbin
		} else {
			log.Println("Casbin 未初始化，跳过 RBAC 中间件。")
		}

		// --- Auth Routes ---

		// Register K8s related routes only if handlers were initialized
		// We check if the specific handlers pointer is non-nil
		if k8sAvailable { // Optional log: k8sAvailable check here gives context
			log.Println("注册 Kubernetes API 路由...")
			if handlers.PodHandler != nil {
				routes.RegisterPodRoutes(v1, handlers.PodHandler)
			} else {
				log.Println("跳过 Pod 路由注册: Handler 未初始化。")
			} // Optional detailed logs
			if handlers.DeploymentHandler != nil {
				routes.RegisterDeploymentRoutes(v1, handlers.DeploymentHandler)
			} else {
				log.Println("跳过 Deployment 路由注册: Handler 未初始化。")
			}
			if handlers.DaemonSetHandler != nil {
				routes.RegisterDaemonSetRoutes(v1, handlers.DaemonSetHandler)
			} else {
				log.Println("跳过 DaemonSet 路由注册: Handler 未初始化。")
			}
			if handlers.ServiceHandler != nil {
				routes.RegisterServiceRoutes(v1, handlers.ServiceHandler)
			} else {
				log.Println("跳过 Service 路由注册: Handler 未初始化。")
			}
			if handlers.IngressHandler != nil {
				routes.RegisterIngressRoutes(v1, handlers.IngressHandler)
			} else {
				log.Println("跳过 Ingress 路由注册: Handler 未初始化。")
			}
			if handlers.NetworkPolicyHandler != nil {
				routes.RegisterNetworkPolicyRoutes(v1, handlers.NetworkPolicyHandler)
			} else {
				log.Println("跳过 NetworkPolicy 路由注册: Handler 未初始化。")
			}
			if handlers.ConfigMapHandler != nil {
				routes.RegisterConfigMapRoutes(v1, handlers.ConfigMapHandler)
			} else {
				log.Println("跳过 ConfigMap 路由注册: Handler 未初始化。")
			}
			if handlers.SecretHandler != nil {
				routes.RegisterSecretRoutes(v1, handlers.SecretHandler)
			} else {
				log.Println("跳过 Secret 路由注册: Handler 未初始化。")
			}
			if handlers.PVCHandler != nil {
				routes.RegisterPVCRoutes(v1, handlers.PVCHandler)
			} else {
				log.Println("跳过 PVC 路由注册: Handler 未初始化。")
			}
			if handlers.PVHandler != nil {
				routes.RegisterPVRoutes(v1, handlers.PVHandler)
			} else {
				log.Println("跳过 PV 路由注册: Handler 未初始化。")
			}
			if handlers.StatefulSetHandler != nil {
				routes.RegisterStatefulSetRoutes(v1, handlers.StatefulSetHandler)
			} else {
				log.Println("跳过 StatefulSet 路由注册: Handler 未初始化。")
			}
			if handlers.NodeHandler != nil {
				routes.RegisterNodeRoutes(v1, handlers.NodeHandler)
			} else {
				log.Println("跳过 Node 路由注册: Handler 未初始化。")
			}
			if handlers.NamespaceHandler != nil {
				routes.RegisterNamespaceRoutes(v1, handlers.NamespaceHandler)
			} else {
				log.Println("跳过 Namespace 路由注册: Handler 未初始化。")
			}
			if handlers.SummaryHandler != nil {
				routes.RegisterSummaryRoutes(v1, handlers.SummaryHandler)
			} else {
				log.Println("跳过 Summary 路由注册: Handler 未初始化。")
			}
			if handlers.EventsHandler != nil {
				routes.RegisterEventsRoutes(v1, handlers.EventsHandler)
			} else {
				log.Println("跳过 Events 路由注册: Handler 未初始化。")
			}
			if handlers.RbacHandler != nil {
				routes.RegisterRbacRoutes(v1, handlers.RbacHandler)
			} else {
				log.Println("跳过 Rbac 路由注册: Handler 未初始化。")
			}

			// Optional check if any K8s routes were registered
			// This check is still a bit manual, could be more abstract, but works.
			if handlers.PodHandler == nil && handlers.DeploymentHandler == nil && // ... check all k8s handlers ...
				handlers.DaemonSetHandler == nil && handlers.ServiceHandler == nil && handlers.IngressHandler == nil &&
				handlers.NetworkPolicyHandler == nil && handlers.ConfigMapHandler == nil && handlers.SecretHandler == nil &&
				handlers.PVCHandler == nil && handlers.PVHandler == nil && handlers.StatefulSetHandler == nil &&
				handlers.NodeHandler == nil && handlers.NamespaceHandler == nil && handlers.SummaryHandler == nil &&
				handlers.EventsHandler == nil && handlers.RbacHandler == nil {
				log.Println("警告: Kubernetes 似乎可用，但没有注册任何 Kubernetes API 路由。")
			} else {
				log.Println("Kubernetes API 路由注册完成。")
			}

		} else {
			log.Println("Kubernetes 不可用，跳过相关 API 路由注册。")
			// Register a status endpoint if K8s is unavailable
			v1.GET("/kubernetes-status", func(c *gin.Context) {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "Kubernetes service unavailable", "details": "Kubernetes client initialization or connection failed"})
			})
		}

		// Always register non-k8s routes if handlers exists
		log.Println("注册非 Kubernetes API 路由...")
		if handlers.InstallerHandler != nil {
			routes.RegisterInstallerRoutes(v1, handlers.InstallerHandler)
		} else {
			log.Println("警告: Installer handlers 未初始化，无法注册相关路由。")
			v1.GET("/installer-status", func(c *gin.Context) {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "Installer service unavailable", "details": "Installer handlers not initialized"})
			})
		}
	}
	log.Println("API 路由注册完成。")
	return router
}

// InitializeDefaultConfig 初始化默认配置

// InitializeDefaultUser 创建超级管理员账户和游客账户
// func InitializeDefaultUser(e *casbin.Enforcer, db *gorm.DB) {
// 	// --- 检查是否已经存在超级管理员 ---
// 	if _, err := db.Where("username = ?", "admin").First(&models.User{}).Rows(); err == nil {
// 		// 用户已存在，不做任何操作
// 		log.Println("默认用户 'admin' 已存在，跳过创建。")
// 	} else {
// 		// 加密初始化密码
// 		password, _ := bcrypt.GenerateFromPassword([]byte("cilikube"), bcrypt.DefaultCost)
// 		user := &models.User{
// 			Username: "admin",
// 			Password: string(password),
// 			Role:     "super_admin",
// 		}
// 		_, err := db.Create(user).Rows()
// 		if err != nil {
// 			log.Fatalf("创建默认用户 'admin' 失败: %v", err)
// 		}
// 		// 绑定用户的角色权限
// 		// --- 默认用户绑定默认权限 ---
// 		// admin <- super_admin
// 		if _, err := e.AddRoleForUser("admin", "super_admin"); err != nil {
// 			log.Fatalf("绑定用户 'admin' 的角色 'super_admin' 失败: %v", err)
// 		}
// 		log.Println("超级管理员用户初始化完成。")
// 	}

// 	// --- 添加游客账户 ---
// 	if err := db.Where("username = ?", "visitor").First(&models.User{}).Row(); err != nil {
// 		log.Println("默认游客账号已存在，不再创建。")
// 	} else {
// 		password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
// 		user := &models.User{
// 			Username: "visitor",
// 			Password: string(password),
// 			Role:     "normal_user",
// 		}
// 		_, err := db.Create(user).Rows()
// 		if err != nil {
// 			log.Fatalf("创建默认游客账号 'visitor' 失败: %v", err)
// 		}
// 		// 给游客绑定 normal_user 的权限，全部都是只给GET
// 		if _, err := e.AddRoleForUser("visitor", "normal_user"); err != nil {
// 			log.Fatalf("绑定用户 'visitor' 的角色'normal_user' 失败: %v", err)
// 		}
// 	}
// 	log.Println("游客账户初始化完成。")
// }

func Cleanup() {
	// 关闭数据库连接
	if err := database.CloseDatabase(); err != nil {
		log.Fatalf("关闭数据库连接失败: %v", err)
	}
	log.Println("数据库连接已关闭。")
}
