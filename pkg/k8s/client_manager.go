// pkg/k8s/client_manager.go
package k8s

import (
	"fmt"
	"log"  // 使用 log 包进行日志输出
	"sort" // 用于对集群名称列表进行排序
	"sync"

	// 确保路径正确指向你的 configs 包
	"github.com/ciliverse/cilikube/configs"
)

// ClientManager 管理多个 Kubernetes 客户端实例
type ClientManager struct {
	clients           map[string]*Client // 存储集群名称到 K8s 客户端的映射
	defaultClientName string             // 默认集群的名称
	mu                sync.RWMutex       // 读写锁，保护 clients 映射的并发访问
}

var (
	globalClientManager *ClientManager // 全局客户端管理器实例
	initOnce            sync.Once      // 确保 ClientManager 初始化只执行一次
)

// InitClientManager 根据提供的集群配置列表初始化全局的 ClientManager。
// clusterCfgs 参数通常来自 configs.GlobalAppConfig.Kubernetes.Clusters。
// 这个函数应该在应用启动时被调用一次。
func InitClientManager(clusterCfgs []configs.IndividualClusterConfig) error {
	var errGlobal error
	initOnce.Do(func() {
		log.Println("信息: 开始初始化 Kubernetes ClientManager...")
		manager := &ClientManager{
			clients: make(map[string]*Client),
		}
		var defaultExplicitlySet bool
		var firstSuccessfullyInitializedClusterName string

		if len(clusterCfgs) == 0 {
			log.Println("警告: ClientManager 初始化时未提供任何集群配置。")
			globalClientManager = manager // 即使没有配置，也初始化一个空的 manager
			return
		}

		for _, cfg := range clusterCfgs {
			if cfg.Name == "" {
				log.Printf("警告: 跳过一个未命名集群的配置 (kubeconfig: %s)\n", cfg.Kubeconfig)
				continue
			}
			log.Printf("信息: 正在为集群 '%s' 初始化客户端 (kubeconfig: '%s')...\n", cfg.Name, cfg.Kubeconfig)
			client, err := NewClient(cfg.Kubeconfig) // 使用你已有的 NewClient 函数
			if err != nil {
				errMessage := fmt.Sprintf("错误: 初始化集群 '%s' 客户端失败: %v\n", cfg.Name, err)
				log.Print(errMessage)
				if errGlobal == nil {
					errGlobal = fmt.Errorf(errMessage) // 只记录第一个遇到的错误
				}
				continue // 跳过此失败的集群
			}

			// 可选但推荐: 检查连接性
			if err := client.CheckConnection(); err != nil {
				errMessage := fmt.Sprintf("错误: 集群 '%s' 连接检查失败: %v\n", cfg.Name, err)
				log.Print(errMessage)
				if errGlobal == nil {
					errGlobal = fmt.Errorf(errMessage)
				}
				// 根据策略，连接检查失败也可能意味着不将此客户端添加到管理器
				// continue
			}

			manager.clients[cfg.Name] = client
			log.Printf("信息: 集群 '%s' 客户端已成功初始化。\n", cfg.Name)
			if firstSuccessfullyInitializedClusterName == "" {
				firstSuccessfullyInitializedClusterName = cfg.Name
			}

			if cfg.IsDefault {
				if defaultExplicitlySet {
					log.Printf("警告: 找到多个 'isDefault: true' 的集群配置。将使用第一个标记为默认的 '%s'。集群 '%s' 的 isDefault 标志将被忽略。\n", manager.defaultClientName, cfg.Name)
				} else {
					manager.defaultClientName = cfg.Name
					defaultExplicitlySet = true
					log.Printf("信息: 集群 '%s' 已被显式设置为默认集群。\n", cfg.Name)
				}
			}
		}

		// 如果没有通过 isDefault 显式设置默认集群，
		// 并且 configs.GlobalAppConfig.Server.ActiveCluster (如果存在并被使用) 也未指定，
		// 则选择第一个成功初始化的集群作为默认。
		if !defaultExplicitlySet && firstSuccessfullyInitializedClusterName != "" {
			// 另一种方式是通过 configs.GlobalAppConfig.Server.ActiveCluster (如果你的 ServerConfig 中有这个字段并希望使用它)
			// activeClusterNameFromConfig := configs.GlobalAppConfig.Server.ActiveCluster // 假设可以访问
			// if activeClusterNameFromConfig != "" {
			// 	if _, ok := manager.clients[activeClusterNameFromConfig]; ok {
			// 		manager.defaultClientName = activeClusterNameFromConfig
			// 		log.Printf("信息: 根据 server.activeCluster 配置, 将 '%s' 设置为默认集群。\n", activeClusterNameFromConfig)
			//      defaultExplicitlySet = true // 标记一下，避免被下面的逻辑覆盖
			// 	} else {
			// 		log.Printf("警告: server.activeCluster 指定的集群 '%s' 未找到或初始化失败，将尝试其他默认逻辑。\n", activeClusterNameFromConfig)
			// 	}
			// }
			// // 如果上面的 activeCluster 逻辑没有设置成功，再使用 firstSuccessfullyInitializedClusterName
			// if !defaultExplicitlySet && firstSuccessfullyInitializedClusterName != "" {
			manager.defaultClientName = firstSuccessfullyInitializedClusterName
			log.Printf("信息: 未显式指定默认集群，自动将第一个成功初始化的集群 '%s' 设置为默认。\n", manager.defaultClientName)
			// }
		}
		globalClientManager = manager
		log.Println("信息: Kubernetes ClientManager 初始化过程完成。")
	})

	if globalClientManager != nil && len(globalClientManager.clients) == 0 && errGlobal == nil && len(clusterCfgs) > 0 {
		return fmt.Errorf("ClientManager 初始化完成，但所有集群客户端都未能成功加载")
	}
	return errGlobal // 返回初始化过程中遇到的第一个错误（如果有）
}

// GetClientManager 返回全局的 ClientManager 实例。
func GetClientManager() (*ClientManager, error) {
	if globalClientManager == nil {
		// 这通常意味着 InitClientManager 还未被调用，或在调用中 panic 了。
		return nil, fmt.Errorf("ClientManager 尚未初始化。请确保在应用启动时调用了 InitClientManager")
	}
	return globalClientManager, nil
}

// GetClient 根据集群名称检索一个 Kubernetes 客户端。
// 如果 clusterName 为空字符串，则尝试返回默认集群的客户端。
func (m *ClientManager) GetClient(clusterName string) (*Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	targetName := clusterName
	if targetName == "" {
		if m.defaultClientName == "" {
			if len(m.clients) == 1 { // 如果只有一个集群，它就是事实上的默认
				for name := range m.clients {
					targetName = name
					log.Printf("警告: 未指定集群名称，且无明确默认集群。自动选择唯一可用集群 '%s'。\n", targetName)
					break
				}
			} else if len(m.clients) == 0 {
				return nil, fmt.Errorf("请求了默认客户端，但 ClientManager 中没有可用的集群客户端")
			} else {
				return nil, fmt.Errorf("请求了默认客户端，但没有配置/确定默认集群，且存在多个 (%d) 集群", len(m.clients))
			}
		} else {
			targetName = m.defaultClientName
			// log.Printf("信息: 未指定集群名称，使用默认集群: '%s'\n", targetName) // 这个日志可能过于频繁
		}
	}

	client, ok := m.clients[targetName]
	if !ok {
		// 需要区分是集群未配置，还是配置了但初始化失败
		// 可以通过检查原始配置列表来确定是否“未配置”
		originalCfgExists := false
		if configs.GlobalAppConfig != nil { // 确保全局配置已加载
			for _, cfg := range configs.GlobalAppConfig.Kubernetes.Clusters {
				if cfg.Name == targetName {
					originalCfgExists = true
					break
				}
			}
		}
		if originalCfgExists {
			return nil, fmt.Errorf("集群 '%s' 已配置但其客户端初始化失败或不可用", targetName)
		}
		return nil, fmt.Errorf("未找到名为 '%s' 的集群的配置或其客户端", targetName)
	}
	return client, nil
}

// GetDefaultClient 显式获取默认集群的客户端。
func (m *ClientManager) GetDefaultClient() (*Client, error) {
	return m.GetClient("") // GetClient 内部处理空字符串即为获取默认
}

// ListClusterNames 返回所有已成功初始化的集群的名称列表 (已排序)。
func (m *ClientManager) ListClusterNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.clients) == 0 {
		return []string{}
	}
	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	sort.Strings(names) // 对名称进行排序
	return names
}
