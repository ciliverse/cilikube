package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// time is still needed for healthz in main
	"github.com/casbin/casbin/v2"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/initialization" // Import the new package
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/database"
	"github.com/ciliverse/cilikube/pkg/k8s" // Your custom k8s client package
)

func main() {
	// --- Configuration Loading ---
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("初始化失败: 加载配置失败: %v", err)
	}
	log.Println("配置加载成功。")

	// --- 数据库初始化 ---
	// Initialize the database connection
	// if err := database.InitDatabase(); err != nil {
	// 	log.Fatalf("初始化失败: 数据库连接失败: %v", err)
	// }
	// log.Println("数据库连接成功。")
	// // --- 数据库自动迁移 ---
	// if err := database.AutoMigrate(); err != nil {
	// 	log.Fatalf("初始化失败: 数据库自动迁移失败: %v", err)
	// }
	// log.Println("数据库自动迁移成功。")

	// --- Kubernetes Client Initialization ---
	// initializeK8sClient remains in main as it's the connection point for k8s
	k8sClient, k8sAvailable := initializeK8sClient(cfg)

	// --- Application Initialization (Services & Handlers) ---
	// Call functions from the new initialization package
	// repositories := initialization.InitializeRepositories(DB)
	services := initialization.InitializeServices(k8sClient, k8sAvailable, cfg)
	appHandlers := initialization.InitializeHandlers(services)

	// <--- ADDED: Casbin Initialization ---
	var e *casbin.Enforcer // 声明 enforcer 变量
	if cfg.Database.Enabled {
		// 确保数据库已成功初始化
		if database.DB != nil {
			e, err = auth.InitCasbin(database.DB) // 调用 InitCasbin
			if err != nil {
				log.Fatalf("初始化 Casbin 失败: %v", err)
			}
			log.Println("Casbin 初始化成功。")
			// (可选) 在这里调用初始化默认用户和 Casbin 角色绑定
			// initialization.InitializeDefaultUser(e, database.DB)
		} else {
			log.Println("警告: 数据库已启用但未成功初始化，跳过 Casbin 初始化。")
		}
	} else {
		log.Println("警告: 数据库未启用，跳过 Casbin 初始化。")
	}
	// <--- ADDED: Casbin Initialization End ---

	// --- Gin Router Setup ---
	// Call function from the new initialization package
	router := initialization.SetupRouter(cfg, appHandlers, k8sAvailable, e)

	// --- Start Server ---
	// startServer remains in main as it's the server lifecycle management
	initialization.StartServer(cfg, router)

	log.Println("服务器已关闭。")
}

// loadConfig loads the application configuration.
// Kept in main for program entry point setup.
func loadConfig() (*configs.Config, error) {
	log.Println("加载配置文件...")

	// 优先从命令行参数获取配置文件路径
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.Parse()

	// 如果命令行参数未指定，则从环境变量获取
	if configPath == "" {
		configPath = os.Getenv("CILIKUBE_CONFIG_PATH")
	}

	// 如果环境变量也未指定，则使用默认路径
	if configPath == "" {
		//获取工作目录
		if wd, err := os.Getwd(); err == nil {
			// 当工作路径为项目根路径时，使用绝对路径
			configPath = wd + "/configs/config.yaml"
			if _, err := os.Stat(configPath); err == nil {
				log.Printf("使用默认配置文件路径: %s\n", configPath)
				return configs.Load(configPath)
			}
		}
		// 当工作路径为当前路径时，使用相对路径
		configPath = "../../configs/config.yaml"
	}

	// 调用 configs.Load 方法加载配置
	cfg, err := configs.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	log.Println("配置文件加载成功。")
	return cfg, nil
}

// initializeK8sClient initializes the Kubernetes client and checks connectivity.
// Kept in main as it's a core dependency check for the application startup.
func initializeK8sClient(cfg *configs.Config) (*k8s.Client, bool) {
	// Determine kubeconfig path
	kubeconfigPath := cfg.Kubernetes.Kubeconfig
	if kubeconfigPath == "in-cluster" {
		kubeconfigPath = "" // NewClient treats "" as in-cluster attempt
		log.Println("配置指定使用 in-cluster Kubernetes 配置。")
	} else if kubeconfigPath == "" {
		log.Println("配置文件中未指定 kubeconfig 路径，将尝试 in-cluster 配置。")
	} else {
		log.Printf("从配置文件获取 kubeconfig 路径: %s\n", kubeconfigPath)
	}

	// Initialize your custom Kubernetes client
	k8sClient, err := k8s.NewClient(kubeconfigPath)
	if err != nil {
		log.Printf("警告: 创建 Kubernetes 客户端失败: %v。Kubernetes 相关功能将不可用。", err)
		return nil, false
	}

	log.Println("Kubernetes 客户端创建成功。检查集群连接...")
	// Perform a connection check
	if err := k8sClient.CheckConnection(); err != nil {
		log.Printf("警告: 无法连接 Kubernetes 集群: %v。Kubernetes 相关功能将不可用。", err)
		// Return client even if check fails, some limited info might be possible,
		// but k8sAvailable flag controls feature availability.
		return k8sClient, false
	}

	log.Println("成功连接到 Kubernetes 集群。")
	// Log API Server URL for confirmation (optional)
	if k8sClient.Config != nil {
		log.Printf("连接到 API Server: %s", k8sClient.Config.Host)
	}
	return k8sClient, true
}
