package main

import (
	"flag"
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/initialization"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/database"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

func main() {
	// --- 1. 加载配置 ---
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CILIKUBE_CONFIG_PATH")
	}
	if configPath == "" {
		configPath = "configs/config.yaml" // 默认值
	}

	if _, err := configs.Load(configPath); err != nil {
		log.Fatalf("FATAL: 加载配置失败: %v", err)
	}
	log.Println("配置加载成功。")

	// --- 2. 数据库和 Store 初始化 ---
	var clusterStore store.ClusterStore

	if configs.GlobalConfig.Database.Enabled {
		if err := database.InitDatabase(); err != nil {
			log.Fatalf("FATAL: 数据库连接失败: %v", err)
		}
		defer database.CloseDatabase()

		// 数据库迁移
		if err := database.AutoMigrate(); err != nil {
			log.Fatalf("FATAL: 数据库自动迁移失败: %v", err)
		}

		// 创建默认用户
		if err := database.CreateDefaultAdmin(); err != nil {
			log.Fatalf("FATAL: 创建默认管理员失败: %v", err)
		}

		// 只有在数据库启用并成功初始化后，才创建 ClusterStore
		if configs.GlobalConfig.Server.EncryptionKey == "" {
			log.Fatalf("FATAL: 数据库已启用，但配置中未设置 EncryptionKey。")
		}
		encryptionKey := []byte(configs.GlobalConfig.Server.EncryptionKey)
		var err error
		clusterStore, err = store.NewGormClusterStore(database.DB, encryptionKey)
		if err != nil {
			log.Fatalf("FATAL: 初始化 Cluster Store 失败: %v", err)
		}
	} else {
		log.Println("数据库未启用，跳过相关初始化。ClusterStore 将为 nil。")
	}

	// --- 3. 初始化 ClusterManager ---
	k8sManager, err := k8s.NewClusterManager(clusterStore, configs.GlobalConfig)
	if err != nil {
		log.Fatalf("FATAL: 初始化 Kubernetes 集群管理器失败: %v", err)
	}

	// --- 4. 初始化应用服务和处理器 ---
	services := initialization.InitializeServices(k8sManager, configs.GlobalConfig)
	appHandlers := initialization.InitializeHandlers(services, k8sManager)

	// --- 5. Casbin 初始化 ---
	var e *casbin.Enforcer
	if configs.GlobalConfig.Database.Enabled {
		var casbinErr error
		e, casbinErr = auth.InitCasbin(database.DB)
		if casbinErr != nil {
			log.Fatalf("初始化 Casbin 失败: %v", casbinErr)
		}
	}

	// --- 6. Gin 路由器设置 ---
	router := initialization.SetupRouter(configs.GlobalConfig, appHandlers, e)

	// --- 7. 启动服务器 ---
	initialization.StartServer(configs.GlobalConfig, router)
}
