package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath" // 确保导入 path/filepath

	"github.com/casbin/casbin/v2"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/initialization"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/database"
	"github.com/ciliverse/cilikube/pkg/k8s" // 引入我们的 k8s 包
)

func main() {
	// --- 配置加载 ---
	cfg, err := loadConfig() // loadConfig 函数保持不变
	if err != nil {
		log.Fatalf("初始化失败: 加载配置失败: %v", err)
	}
	log.Println("配置加载成功。")

	// --- 数据库初始化 ---
	// (数据库初始化逻辑保持不变)
	if cfg.Database.Enabled {
		if err := database.InitDatabase(); err != nil { // 假设 InitDatabase 使用 cfg.GetDSN() 或 GlobalConfig
			log.Fatalf("初始化失败: 数据库连接失败: %v", err)
		}
		log.Println("数据库连接成功。")
		if err := database.AutoMigrate(); err != nil {
			log.Fatalf("初始化失败: 数据库自动迁移失败: %v", err)
		}
		log.Println("数据库自动迁移成功。")
	} else {
		log.Println("数据库未启用，跳过数据库初始化。")
	}

	// --- Kubernetes 集群管理器初始化 ---
	// 旧的 initializeK8sClient 函数已被移除
	// 我们现在使用 NewClusterManager 来初始化所有在 cfg.Clusters 中定义的集群
	log.Println("初始化 Kubernetes 集群管理器...")
	k8sClusterManager, k8sClusterAvailabilityStatus := k8s.NewClusterManager()
	// k8sClusterAvailabilityStatus 是一个 map[string]bool，告诉每个配置的集群是否成功初始化

	// 判断是否至少有一个 Kubernetes 集群可用。
	// 这个 anyK8sAvailable 标志可以用于全局健康检查或作为 SetupRouter 的参数。
	var anyK8sClientActuallyAvailable bool
	for clusterName, available := range k8sClusterAvailabilityStatus {
		if available {
			anyK8sClientActuallyAvailable = true
			log.Printf("集群 '%s' 初始化成功且可用。", clusterName)
		} else {
			log.Printf("集群 '%s' 初始化失败或不可用。", clusterName)
		}
	}

	if anyK8sClientActuallyAvailable {
		log.Println("至少有一个 Kubernetes 集群成功初始化并可用。")
	} else {
		log.Println("警告: 没有 Kubernetes 集群成功初始化。Kubernetes 相关功能将严重受限或不可用。")
	}

	// --- 应用初始化 (服务 & 处理器) ---
	// InitializeServices 和 InitializeHandlers 将需要接收 k8sClusterManager。
	// 这是一个重要的变化

	// 为了平滑过渡，可以尝试获取一个 "初始" 或 "默认活动" 的客户端实例
	// 这个客户端可以用于那些尚未完全适配多集群、或者需要一个默认上下文的服务
	// 注意：这仍然是过渡性措施，最终目标是服务和处理器能够按需从 k8sClusterManager 中获取指定集群的客户端。
	var initialK8sClientForServices *k8s.Client // 这是 *k8s.Client 类型
	var initialK8sClientIsAvailableForServices bool

	if cfg.Server.ActiveCluster != "" {
		log.Printf("尝试使用配置文件中指定的活动集群 '%s' 作为服务层的初始客户端。", cfg.Server.ActiveCluster)
		client, err := k8sClusterManager.GetClient(cfg.Server.ActiveCluster)
		if err == nil {
			// 再次确认这个从Manager获取的Client在初始化时是真的available
			if status, ok := k8sClusterAvailabilityStatus[cfg.Server.ActiveCluster]; ok && status {
				initialK8sClientForServices = client
				initialK8sClientIsAvailableForServices = true
				log.Printf("成功获取活动集群 '%s' 的客户端 (%s) 作为服务层初始客户端。", cfg.Server.ActiveCluster, client.Config.Host)
			} else {
				log.Printf("警告: 配置的活动集群 '%s' 在启动时未能成功初始化或标记为不可用。将尝试查找其他可用集群。", cfg.Server.ActiveCluster)
			}
		} else {
			log.Printf("警告: 无法从集群管理器获取名为 '%s' 的客户端 (错误: %v)。可能是名称配置错误或该集群未成功初始化。将尝试查找其他可用集群。", cfg.Server.ActiveCluster, err)
		}
	}

	if !initialK8sClientIsAvailableForServices && anyK8sClientActuallyAvailable {
		log.Println("未指定有效活动集群或活动集群不可用，尝试使用第一个成功初始化的集群作为服务层的初始客户端。")
		availableNames := k8sClusterManager.GetAvailableClientNames() // 获取所有可用集群名称
		if len(availableNames) > 0 {
			firstAvailableName := availableNames[0]
			client, err := k8sClusterManager.GetClient(firstAvailableName) // 应该不会出错，因为刚从可用列表获取
			if err == nil {
				initialK8sClientForServices = client
				initialK8sClientIsAvailableForServices = true
				log.Printf("已选择第一个可用集群 '%s' (%s) 作为服务层的初始客户端。", firstAvailableName, client.Config.Host)
			}
		}
	}

	if !initialK8sClientIsAvailableForServices {
		log.Println("警告: 未能为服务层提供一个初始的 Kubernetes 客户端。依赖初始客户端的 Kubernetes 服务可能无法正常工作。")
	}

	// 注意：InitializeServices 和 InitializeHandlers 的函数签名将在下一步中修改，
	// 以便它们可以接收并使用 k8sClusterManager。
	// `initialK8sClientForServices` 和 `initialK8sClientIsAvailableForServices` 是过渡参数。
	services := initialization.InitializeServices(
		k8sClusterManager,                      // 传递管理器本身
		initialK8sClientForServices,            // 传递选定的初始客户端 (可能为 nil)
		initialK8sClientIsAvailableForServices, // 初始客户端的可用状态
		cfg,
	)
	appHandlers := initialization.InitializeHandlers(services, k8sClusterManager, cfg) // Handler 也需要 Manager

	// --- Casbin 初始化 ---
	// (Casbin 初始化逻辑保持不变)
	var e *casbin.Enforcer
	if cfg.Database.Enabled && database.DB != nil { // 确保数据库已初始化
		var casbinErr error
		e, casbinErr = auth.InitCasbin(database.DB)
		if casbinErr != nil {
			log.Fatalf("初始化 Casbin 失败: %v", casbinErr)
		}
		log.Println("Casbin 初始化成功。")
		// 可选: initialization.InitializeDefaultUser(e, database.DB)
	} else if cfg.Database.Enabled && database.DB == nil {
		log.Println("警告: 数据库已启用但未成功初始化，跳过 Casbin 初始化。")
	} else {
		log.Println("警告: 数据库未启用，跳过 Casbin 初始化。")
	}

	// --- Gin 路由器设置 ---
	// SetupRouter 也需要 k8sClusterManager 来正确地将请求路由到目标集群的处理器
	// anyK8sClientActuallyAvailable 用于指示是否有任何 K8s 功能可用（例如，影响 /healthz 端点）。
	router := initialization.SetupRouter(cfg, appHandlers, anyK8sClientActuallyAvailable, e, k8sClusterManager) // 传递 k8sClusterManager

	// --- 启动服务器 ---
	// (startServer 逻辑保持不变)
	initialization.StartServer(cfg, router)

	log.Println("服务器已关闭。")
}

// 这里为了完整性再次列出，并确保 filepath 被导入。
func loadConfig() (*configs.Config, error) {
	log.Println("加载配置文件...")

	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CILIKUBE_CONFIG_PATH")
	}

	if configPath == "" {
		// 尝试基于工作目录的路径
		wd, err := os.Getwd()
		if err == nil {
			potentialPath := filepath.Join(wd, "configs", "config.yaml")
			if _, statErr := os.Stat(potentialPath); statErr == nil {
				configPath = potentialPath
				log.Printf("使用默认配置文件路径 (基于工作目录): %s\n", configPath)
				return configs.Load(configPath) // 直接加载并返回
			}
		}
		// 如果上述未找到，尝试上一级目录的 configs (适用于 cmd/appname/main.go 结构)
		// 这是一个常见的项目结构，但可能需要根据您的实际情况调整
		// ../../configs/config.yaml 假设 main.go 在两级子目录下，例如 cmd/your_app/main.go
		// 如果 main.go 在项目根目录，应该是 ./configs/config.yaml (上面已处理)
		// 如果 main.go 在一级子目录 (例如 cmd/main.go), 应该是 ../configs/config.yaml
		// 为了更普适继续之前的逻辑但优先使用绝对路径或明确的环境变量/命令行参数。
		// 这里的路径查找逻辑可能需要根据您的项目实际结构进行微调。
		// 如果上面的 wd/configs/config.yaml 不存在，则执行到这里。
		// 维持您之前代码中的相对路径作为最后的备选方案
		if wd != "" { // 确保 wd 获取成功
			// 示例: 假设 main.go 在项目根下的 cmd/cilikube/ 目录，而 configs 在项目根目录
			// 则需要 ../../configs/config.yaml
			// 但如果是在项目根目录直接 go run main.go，则是 configs/config.yaml (上面已处理)
			// 此处的相对路径非常依赖于执行命令时的工作目录。
			// 最好是通过命令行参数或环境变量指定。
			// 作为最后的备选，使用您提供的 "../../configs/config.yaml"
			// （如果之前的 wd + "/configs/config.yaml" 逻辑已满足，这里可能不需要了）
		}
		// 如果到这里 configPath 仍然为空，使用您之前代码中的最终备选路径
		if configPath == "" {
			configPath = "../../configs/config.yaml" // 这个路径非常依赖执行上下文
			log.Printf("尝试使用备选相对路径: %s\n", configPath)
		}
	}
	// 确保在调用 configs.Load 之前，configPath 是确定的。
	if configPath == "" {
		return nil, fmt.Errorf("无法确定配置文件路径。请使用 -config 命令行参数或 CILIKUBE_CONFIG_PATH 环境变量指定。")
	}

	log.Printf("最终尝试加载配置文件: %s\n", configPath)
	cfg, err := configs.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置 '%s' 失败: %w", configPath, err)
	}

	log.Println("配置文件加载成功。")
	return cfg, nil
}

// 旧的 initializeK8sClient 函数已从 main.go 中移除，
// 其功能被 pkg/k8s/manager.go 中的 initializeSingleK8sClient 和 NewClusterManager 替代
