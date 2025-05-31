package k8s

import (
	"fmt"
	"log"
	"sync"

	"github.com/ciliverse/cilikube/configs" // 确保路径正确
	// "k8s.io/client-go/kubernetes" // Client struct already holds kubernetes.Interface
	// "k8s.io/client-go/rest"       // Client struct already holds *rest.Config
)

// ClusterManager 负责管理多个 Kubernetes 集群的客户端实例。
type ClusterManager struct {
	clients map[string]*Client // 键是集群名称 (来自 ClusterInfo.Name)，值是对应的 *k8s.Client
	lock    sync.RWMutex       // 用于保护 clients map 的并发访问
}

// NewClusterManager 根据提供的配置初始化所有已定义的、活动的集群客户端。
// 返回 ClusterManager 实例和一个记录各集群初始化状态的 map (clusterName -> bool)。
func NewClusterManager(appConfig *configs.Config) (*ClusterManager, map[string]bool) {
	manager := &ClusterManager{
		clients: make(map[string]*Client),
	}
	clusterAvailability := make(map[string]bool)

	if len(appConfig.Clusters) == 0 {
		log.Println("警告: 配置文件中没有定义任何集群 (config.Clusters 列表为空)。")
		// 此时可以考虑是否尝试使用 appConfig.Kubernetes.Kubeconfig 作为备用或默认集群
		// 但为了清晰的多集群管理，我们主要依赖 Clusters 列表。
		// 如果需要备用逻辑，可以在这里添加。例如：
		// if appConfig.Kubernetes.Kubeconfig != "" {
		//     log.Printf("尝试使用顶层 kubernetes.kubeconfig (%s) 初始化一个名为 'default_fallback' 的集群", appConfig.Kubernetes.Kubeconfig)
		//     client, available := initializeSingleK8sClient(appConfig.Kubernetes.Kubeconfig, "default_fallback")
		//     if available {
		//         manager.clients["default_fallback"] = client
		//         clusterAvailability["default_fallback"] = true
		//         log.Println("已使用顶层 kubernetes.kubeconfig 初始化 'default_fallback' 集群。")
		//     }
		// }
	}

	log.Println("开始初始化配置文件中定义的 Kubernetes 集群客户端...")
	for _, clusterInfo := range appConfig.Clusters {
		if clusterInfo.Name == "" {
			log.Printf("警告: 配置中发现一个未命名的集群条目 (ConfigPath: '%s')，已跳过。请为所有集群提供一个唯一的 'name'。", clusterInfo.ConfigPath)
			continue
		}

		if !clusterInfo.IsActive {
			log.Printf("集群 '%s' (ConfigPath: '%s') 在配置中标记为非活动 (IsActive: false)，跳过初始化。", clusterInfo.Name, clusterInfo.ConfigPath)
			clusterAvailability[clusterInfo.Name] = false
			continue
		}

		if _, exists := manager.clients[clusterInfo.Name]; exists {
			log.Printf("警告: 集群名称 '%s' 重复。将使用第一个定义，后续同名配置 '%s' 将被忽略。", clusterInfo.Name, clusterInfo.ConfigPath)
			continue
		}

		log.Printf("开始初始化集群: '%s', Kubeconfig路径: '%s'", clusterInfo.Name, clusterInfo.ConfigPath)
		client, available := initializeSingleK8sClient(clusterInfo.ConfigPath, clusterInfo.Name)
		clusterAvailability[clusterInfo.Name] = available
		if available {
			manager.clients[clusterInfo.Name] = client
			log.Printf("集群 '%s' 的 Kubernetes 客户端创建并连接成功。", clusterInfo.Name)
			if client.Config != nil { // client 是 *k8s.Client 类型
				log.Printf("集群 '%s' 连接到 API Server: %s", clusterInfo.Name, client.Config.Host)
			}
		} else {
			log.Printf("警告: 创建或连接集群 '%s' (Kubeconfig: '%s') 的 Kubernetes 客户端失败。该集群相关功能将不可用。", clusterInfo.Name, clusterInfo.ConfigPath)
		}
	}

	if len(manager.clients) == 0 {
		log.Println("警告: 未能成功初始化任何在 config.Clusters 中定义的活动集群客户端。")
	} else {
		log.Printf("成功初始化 %d 个 Kubernetes 集群客户端。", len(manager.clients))
	}

	return manager, clusterAvailability
}

// initializeSingleK8sClient 是一个辅助函数，用于为单个集群配置初始化 k8s.Client。
// kubeconfigPath 可以是文件路径，也可以是 "in-cluster"。
// clusterName 用于日志记录。
func initializeSingleK8sClient(kubeconfigPath string, clusterNameForLog string) (*Client, bool) {
	var clientLogName string
	if clusterNameForLog == "" {
		clientLogName = "未命名集群"
	} else {
		clientLogName = fmt.Sprintf("集群 '%s'", clusterNameForLog)
	}

	effectiveKubeconfigPath := kubeconfigPath
	if kubeconfigPath == "in-cluster" {
		effectiveKubeconfigPath = "" // k8s.NewClient 将空字符串视作尝试 InClusterConfig
		log.Printf("%s 配置指定使用 'in-cluster' Kubernetes 配置。", clientLogName)
	} else if kubeconfigPath == "" {
		// 如果允许ConfigPath为空并也解释为in-cluster，可以在这里处理或依赖NewClient的行为
		log.Printf("警告: %s 的 kubeconfig 路径为空。k8s.NewClient 可能会尝试 'in-cluster' 或其他默认行为。", clientLogName)
		// 如果希望空路径也明确表示in-cluster，则 effectiveKubeconfigPath = ""
	} else {
		log.Printf("%s 使用配置文件路径: '%s'", clientLogName, kubeconfigPath)
	}

	// 调用你现有的 pkg/k8s/client.go 中的 NewClient 工厂函数
	k8sClientInstance, err := NewClient(effectiveKubeconfigPath) // NewClient 来自你的 pkg/k8s/client.go
	if err != nil {
		log.Printf("警告: 为 %s 创建 Kubernetes 客户端失败: %v。Kubernetes 相关功能将不可用。", clientLogName, err)
		return nil, false
	}

	// k8s.NewClient 成功不代表连接一定成功，CheckConnection 会做进一步检查
	log.Printf("%s 的 Kubernetes 客户端核心创建成功。正在检查到 API Server 的连接...", clientLogName)
	if err := k8sClientInstance.CheckConnection(); err != nil {
		log.Printf("警告: 无法连接到 %s 的 Kubernetes API Server: %v。Kubernetes 相关功能将受限。", clientLogName, err)
		// 即使连接检查失败，也可能返回客户端实例，但标记为不可用。
		// 取决于你的策略，如果连接失败则不应使用该客户端，可以返回 nil, false
		return k8sClientInstance, false // 或者 return nil, false，如果严格要求连接成功
	}

	log.Printf("%s 成功连接到 Kubernetes API Server。", clientLogName)
	return k8sClientInstance, true
}

// GetClient 根据集群名称检索已初始化的 Kubernetes 客户端。
// 如果找不到或客户端未成功初始化，将返回错误。
func (cm *ClusterManager) GetClient(clusterName string) (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	client, exists := cm.clients[clusterName]
	if !exists {
		return nil, fmt.Errorf("未找到名为 '%s' 的集群客户端，或者该客户端未在启动时成功初始化", clusterName)
	}
	// 这里可以根据需要添加额外的健康检查逻辑，但通常在获取时假定初始化成功的客户端是可用的。
	// 调用方应处理使用客户端时可能发生的网络错误等。
	return client, nil
}

// GetAllClients 返回所有成功初始化的客户端的映射副本。
// 主要用于需要遍历所有集群的场景（请谨慎使用）。
func (cm *ClusterManager) GetAllClients() map[string]*Client {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	clientsCopy := make(map[string]*Client, len(cm.clients))
	for name, client := range cm.clients {
		clientsCopy[name] = client
	}
	return clientsCopy
}

// GetAvailableClientNames 返回所有成功初始化并被认为可用的集群名称列表。
func (cm *ClusterManager) GetAvailableClientNames() []string {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	names := make([]string, 0, len(cm.clients))
	for name := range cm.clients {
		// 这里的 "可用" 指的是初始化时成功。可以结合 clusterAvailability 状态图。
		// 为简单起见，如果它在 cm.clients 中，就认为它在启动时是可用的。
		names = append(names, name)
	}
	return names
}
