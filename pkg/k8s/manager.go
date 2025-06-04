package k8s

import (
	"fmt"
	"log"
	"sync"

	"github.com/ciliverse/cilikube/configs" // 确保路径正确
	"k8s.io/client-go/kubernetes"           // 为 GetActiveClientset 引入
	"k8s.io/client-go/rest"                 // 为 GetActiveConfig 引入
)

// ClusterManager 负责管理多个 Kubernetes 集群的客户端实例。
type ClusterManager struct {
	clients          map[string]*Client // 键是集群名称 (来自 ClusterInfo.Name)，值是对应的 *k8s.Client
	lock             sync.RWMutex       // 用于保护 clients map 和 GlobalConfig 相关操作的并发访问
	activeClientName string             // 当前活动集群的名称
	activeClient     *Client            // 当前活动集群的客户端实例
}

// NewClusterManager 根据提供的配置初始化所有已定义的、活动的集群客户端。
func NewClusterManager() (*ClusterManager, map[string]bool) {
	if configs.GlobalConfig == nil {
		log.Panicln("错误: 全局配置 (configs.GlobalConfig) 尚未加载。请先调用 configs.Load()")
	}

	manager := &ClusterManager{
		clients: make(map[string]*Client),
	}
	clusterAvailability := make(map[string]bool)

	if len(configs.GlobalConfig.Clusters) == 0 {
		log.Println("警告: 配置文件中没有定义任何集群 (config.Clusters 列表为空)。")
	}

	log.Println("开始初始化配置文件中定义的 Kubernetes 集群客户端...")
	for _, clusterInfo := range configs.GlobalConfig.Clusters {
		if clusterInfo.Name == "" {
			log.Printf("警告: 配置中发现一个未命名的集群条目 (ConfigPath: '%s')，已跳过。", clusterInfo.ConfigPath)
			continue
		}

		if !clusterInfo.IsActive {
			log.Printf("集群 '%s' (ConfigPath: '%s') 在配置中标记为非活动 (IsActive: false)，跳过初始化。", clusterInfo.Name, clusterInfo.ConfigPath)
			clusterAvailability[clusterInfo.Name] = false
			continue
		}

		if _, exists := manager.clients[clusterInfo.Name]; exists {
			log.Printf("警告: 集群名称 '%s' 重复。将使用第一个定义。", clusterInfo.Name)
			continue
		}

		log.Printf("开始初始化集群: '%s', Kubeconfig路径: '%s'", clusterInfo.Name, clusterInfo.ConfigPath)
		client, available := initializeSingleK8sClient(clusterInfo.ConfigPath, clusterInfo.Name)
		clusterAvailability[clusterInfo.Name] = available
		if available {
			manager.clients[clusterInfo.Name] = client
			log.Printf("集群 '%s' 的 Kubernetes 客户端创建并连接成功。", clusterInfo.Name)
			if client.Config != nil {
				log.Printf("集群 '%s' 连接到 API Server: %s", clusterInfo.Name, client.Config.Host)
			}
		} else {
			log.Printf("警告: 创建或连接集群 '%s' (Kubeconfig: '%s') 的 Kubernetes 客户端失败。", clusterInfo.Name, clusterInfo.ConfigPath)
		}
	}

	if configs.GlobalConfig.Server.ActiveCluster != "" {
		activeClusterNameFromConfig := configs.GlobalConfig.Server.ActiveCluster
		if client, exists := manager.clients[activeClusterNameFromConfig]; exists {
			if client != nil && client.Clientset != nil {
				manager.activeClientName = activeClusterNameFromConfig
				manager.activeClient = client
				log.Printf("集群 '%s' 已根据配置文件设置为活动集群。", manager.activeClientName)
			} else {
				log.Printf("警告: 配置文件中指定的活动集群 '%s' 的客户端无效, 未能设置为活动集群。", activeClusterNameFromConfig)
			}
		} else {
			log.Printf("警告: 配置文件中指定的活动集群 '%s' 未找到或未成功初始化, 未能设置为活动集群。", activeClusterNameFromConfig)
		}
	} else if len(manager.clients) > 0 && manager.activeClient == nil {
		log.Println("提示: 配置文件未指定活动集群，当前无活动集群。")
	}

	if len(manager.clients) == 0 && len(configs.GlobalConfig.Clusters) > 0 {
		log.Println("警告: 未能成功初始化任何在 config.Clusters 中定义的活动集群客户端。")
	} else if len(manager.clients) > 0 {
		log.Printf("成功初始化 %d 个 Kubernetes 集群客户端。", len(manager.clients))
	}

	return manager, clusterAvailability
}

func initializeSingleK8sClient(kubeconfigPath string, clusterNameForLog string) (*Client, bool) {
	var clientLogName string
	if clusterNameForLog == "" {
		clientLogName = "未命名集群"
	} else {
		clientLogName = fmt.Sprintf("集群 '%s'", clusterNameForLog)
	}

	effectiveKubeconfigPath := kubeconfigPath
	if kubeconfigPath == "in-cluster" {
		effectiveKubeconfigPath = ""
		log.Printf("%s 配置指定使用 'in-cluster' Kubernetes 配置。", clientLogName)
	} else if kubeconfigPath == "" {
		log.Printf("警告: %s 的 kubeconfig 路径为空。k8s.NewClient 可能会尝试 'in-cluster' 或其他默认行为。", clientLogName)
	} else if kubeconfigPath == "default" {
		log.Printf("%s 使用 'default' kubeconfig 路径。", clientLogName)
	} else {
		log.Printf("%s 使用配置文件路径: '%s'", clientLogName, kubeconfigPath)
	}

	k8sClientInstance, err := NewClient(effectiveKubeconfigPath)
	if err != nil {
		log.Printf("警告: 为 %s 创建 Kubernetes 客户端失败: %v。", clientLogName, err)
		return nil, false
	}

	log.Printf("%s 的 Kubernetes 客户端核心创建成功。正在检查到 API Server 的连接...", clientLogName)
	if err := k8sClientInstance.CheckConnection(); err != nil {
		log.Printf("警告: 无法连接到 %s 的 Kubernetes API Server: %v。", clientLogName, err)
		return k8sClientInstance, false
	}

	log.Printf("%s 成功连接到 Kubernetes API Server。", clientLogName)
	return k8sClientInstance, true
}

// --- 活动客户端管理方法 ---

func (cm *ClusterManager) SetActiveCluster(clusterName string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	client, exists := cm.clients[clusterName]
	if !exists {
		for _, cfg := range configs.GlobalConfig.Clusters {
			if cfg.Name == clusterName {
				if !cfg.IsActive {
					return fmt.Errorf("无法设置活动集群: 集群 '%s' 未激活", clusterName)
				}
				return fmt.Errorf("无法设置活动集群: 集群 '%s' 客户端初始化失败或当前不可用", clusterName)
			}
		}
		return fmt.Errorf("无法设置活动集群: 集群 '%s' 未找到", clusterName)
	}

	if client == nil || client.Clientset == nil {
		return fmt.Errorf("无法设置活动集群: 集群 '%s' 的客户端实例无效", clusterName)
	}

	cm.activeClient = client
	cm.activeClientName = clusterName
	configs.GlobalConfig.Server.ActiveCluster = clusterName

	if err := configs.SaveGlobalConfig(); err != nil {
		return fmt.Errorf("设置活动集群为 '%s' 成功，但保存配置文件失败: %w。请手动检查配置文件。", clusterName, err)
	}

	log.Printf("活动集群已切换为: %s，并已保存到配置文件。", clusterName)
	return nil
}

func (cm *ClusterManager) GetActiveClient() (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	if cm.activeClient == nil {
		if configs.GlobalConfig.Server.ActiveCluster != "" {
			return nil, fmt.Errorf("无活动 Kubernetes 客户端。配置文件期望 '%s' 为活动集群，但该集群可能未成功初始化或已被移除/设为非活动", configs.GlobalConfig.Server.ActiveCluster)
		}
		return nil, fmt.Errorf("无活动 Kubernetes 客户端配置")
	}
	return cm.activeClient, nil
}

func (cm *ClusterManager) GetActiveClientset() (kubernetes.Interface, error) {
	client, err := cm.GetActiveClient()
	if err != nil {
		return nil, err
	}
	if client.Clientset == nil {
		return nil, fmt.Errorf("活动集群 '%s' 的 Clientset 为空", cm.activeClientName)
	}
	return client.Clientset, nil
}

func (cm *ClusterManager) GetActiveConfig() (*rest.Config, error) {
	client, err := cm.GetActiveClient()
	if err != nil {
		return nil, err
	}
	if client.Config == nil {
		return nil, fmt.Errorf("活动集群 '%s' 的 Config 为空", cm.activeClientName)
	}
	return client.Config, nil
}

func (cm *ClusterManager) GetActiveClusterName() string {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	return cm.activeClientName
}

// --- CRUD 方法 ---

func (cm *ClusterManager) AddCluster(newInfo configs.ClusterInfo) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	for _, c := range configs.GlobalConfig.Clusters {
		if c.Name == newInfo.Name {
			return fmt.Errorf("配置错误: 集群 '%s' 已存在于配置文件中", newInfo.Name)
		}
	}
	if _, exists := cm.clients[newInfo.Name]; exists { // 检查内存中的活跃客户端
		return fmt.Errorf("内部错误: 集群 '%s' 已作为活跃客户端存在", newInfo.Name)
	}

	originalClusters := make([]configs.ClusterInfo, len(configs.GlobalConfig.Clusters))
	copy(originalClusters, configs.GlobalConfig.Clusters)
	configs.GlobalConfig.Clusters = append(configs.GlobalConfig.Clusters, newInfo)

	if err := configs.SaveGlobalConfig(); err != nil {
		configs.GlobalConfig.Clusters = originalClusters
		return fmt.Errorf("添加集群 '%s' 失败: 保存配置文件错误: %w", newInfo.Name, err)
	}

	if newInfo.IsActive {
		log.Printf("集群 '%s' 配置已添加并保存，尝试初始化客户端...", newInfo.Name)
		client, available := initializeSingleK8sClient(newInfo.ConfigPath, newInfo.Name)
		if !available {
			log.Printf("警告: 集群 '%s' 配置已添加，但客户端初始化或连接失败。", newInfo.Name)
		} else {
			cm.clients[newInfo.Name] = client
			log.Printf("集群 '%s' 已成功添加并激活。", newInfo.Name)
			if cm.activeClient == nil && newInfo.Name == configs.GlobalConfig.Server.ActiveCluster {
				cm.activeClient = client
				cm.activeClientName = newInfo.Name
				log.Printf("新添加的集群 '%s' 已根据配置文件期望设为活动集群。", newInfo.Name)
			}
		}
	} else {
		log.Printf("集群 '%s' 已添加但标记为非活动, 跳过客户端初始化。", newInfo.Name)
	}
	return nil
}

func (cm *ClusterManager) UpdateCluster(clusterName string, updatedInfo configs.ClusterInfo) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	if clusterName != updatedInfo.Name && updatedInfo.Name != "" {
		return fmt.Errorf("不允许更改集群名称 (从 '%s' 到 '%s')", clusterName, updatedInfo.Name)
	}
	updatedInfo.Name = clusterName

	found := false
	var originalCluster configs.ClusterInfo
	var clusterIndex int = -1
	wasActive := (cm.activeClientName == clusterName)

	for i, c := range configs.GlobalConfig.Clusters {
		if c.Name == clusterName {
			originalCluster = c
			clusterIndex = i
			configs.GlobalConfig.Clusters[i] = updatedInfo
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("集群 '%s' 在配置文件中未找到，无法更新", clusterName)
	}

	if wasActive && !updatedInfo.IsActive {
		configs.GlobalConfig.Server.ActiveCluster = ""
		log.Printf("活动集群 '%s' 被更新为非活动状态，已清除配置文件中的活动集群设置。", clusterName)
	}

	if err := configs.SaveGlobalConfig(); err != nil {
		if clusterIndex != -1 {
			configs.GlobalConfig.Clusters[clusterIndex] = originalCluster
			if wasActive && !updatedInfo.IsActive {
				configs.GlobalConfig.Server.ActiveCluster = clusterName
			}
		}
		return fmt.Errorf("更新集群 '%s' 失败: 保存配置文件错误: %w", clusterName, err)
	}

	log.Printf("集群 '%s' 配置已更新并保存。", clusterName)

	// 使用空白标识符 "_" 忽略 oldClient，因为我们不直接使用它 (例如调用 Close 方法)
	if _, exists := cm.clients[clusterName]; exists {
		delete(cm.clients, clusterName)
		log.Printf("已移除集群 '%s' 的旧客户端实例。", clusterName)
		// 如果 Client 实例需要显式关闭资源，应在此处调用:
		// oldClient.Close() // 并确保 Client 类型有 Close() 方法
	}

	if wasActive {
		cm.activeClient = nil
		cm.activeClientName = ""
		log.Printf("原活动集群 '%s' 已更新，暂时清除活动状态。", clusterName)
	}

	if updatedInfo.IsActive {
		log.Printf("尝试为更新后的集群 '%s' 初始化客户端...", updatedInfo.Name)
		client, available := initializeSingleK8sClient(updatedInfo.ConfigPath, updatedInfo.Name)
		if !available {
			log.Printf("警告: 更新后的集群 '%s' 客户端初始化或连接失败。", updatedInfo.Name)
		} else {
			cm.clients[updatedInfo.Name] = client
			log.Printf("集群 '%s' 已成功更新并激活。", updatedInfo.Name)
			if updatedInfo.Name == configs.GlobalConfig.Server.ActiveCluster {
				cm.activeClient = client
				cm.activeClientName = updatedInfo.Name
				log.Printf("更新后的集群 '%s' 已根据配置文件期望设为/保持为活动集群。", updatedInfo.Name)
			} else if wasActive {
				log.Printf("集群 '%s' 已更新并激活，但不再是配置文件中指定的活动集群 (当前配置文件期望: '%s')。", updatedInfo.Name, configs.GlobalConfig.Server.ActiveCluster)
			}
		}
	} else {
		log.Printf("集群 '%s' 已更新并标记为非活动。", updatedInfo.Name)
	}
	return nil
}

func (cm *ClusterManager) RemoveCluster(clusterName string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	var originalClusters []configs.ClusterInfo
	var newClusters []configs.ClusterInfo
	found := false
	wasActive := (cm.activeClientName == clusterName)

	originalClusters = make([]configs.ClusterInfo, len(configs.GlobalConfig.Clusters))
	copy(originalClusters, configs.GlobalConfig.Clusters)

	for _, c := range configs.GlobalConfig.Clusters {
		if c.Name == clusterName {
			found = true
		} else {
			newClusters = append(newClusters, c)
		}
	}

	if !found {
		return fmt.Errorf("集群 '%s' 在配置文件中未找到，无法删除", clusterName)
	}
	configs.GlobalConfig.Clusters = newClusters

	originalActiveClusterNameInConfig := configs.GlobalConfig.Server.ActiveCluster
	if wasActive {
		configs.GlobalConfig.Server.ActiveCluster = ""
		log.Printf("被删除的集群 '%s' 是活动集群，已清除配置文件中的活动集群设置。", clusterName)
	}

	if err := configs.SaveGlobalConfig(); err != nil {
		configs.GlobalConfig.Clusters = originalClusters
		if wasActive {
			configs.GlobalConfig.Server.ActiveCluster = originalActiveClusterNameInConfig
		}
		return fmt.Errorf("删除集群 '%s' 失败: 保存配置文件错误: %w", clusterName, err)
	}

	// 使用空白标识符 "_" 忽略 oldClient
	if _, exists := cm.clients[clusterName]; exists {
		delete(cm.clients, clusterName)
		log.Printf("已移除集群 '%s' 的客户端实例。", clusterName)
		// 如果 Client 实例需要显式关闭资源，应在此处调用:
		// oldClient.Close() // 并确保 Client 类型有 Close() 方法
	}

	if wasActive {
		cm.activeClient = nil
		cm.activeClientName = ""
		log.Printf("原活动集群 '%s' 已删除，清除活动状态。", clusterName)
	}
	log.Printf("集群 '%s' 配置已成功删除。", clusterName)
	return nil
}

func (cm *ClusterManager) GetClusterConfiguration(clusterName string) (configs.ClusterInfo, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	for _, c := range configs.GlobalConfig.Clusters {
		if c.Name == clusterName {
			return c, nil
		}
	}
	return configs.ClusterInfo{}, fmt.Errorf("集群 '%s' 的配置未找到", clusterName)
}

func (cm *ClusterManager) ListClusterConfigurations() []configs.ClusterInfo {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	list := make([]configs.ClusterInfo, len(configs.GlobalConfig.Clusters))
	copy(list, configs.GlobalConfig.Clusters)
	return list
}

func (cm *ClusterManager) GetClient(clusterName string) (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	client, exists := cm.clients[clusterName]
	if !exists {
		for _, cfg := range configs.GlobalConfig.Clusters {
			if cfg.Name == clusterName {
				if !cfg.IsActive {
					return nil, fmt.Errorf("名为 '%s' 的集群客户端未激活", clusterName)
				}
				return nil, fmt.Errorf("名为 '%s' 的集群客户端初始化失败或当前不可用", clusterName)
			}
		}
		return nil, fmt.Errorf("未找到名为 '%s' 的集群客户端或配置", clusterName)
	}
	return client, nil
}

func (cm *ClusterManager) GetAllClients() map[string]*Client {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	clientsCopy := make(map[string]*Client, len(cm.clients))
	for name, client := range cm.clients {
		clientsCopy[name] = client
	}
	return clientsCopy
}

func (cm *ClusterManager) GetAvailableClientNames() []string {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	names := make([]string, 0, len(cm.clients))
	for name := range cm.clients {
		names = append(names, name)
	}
	return names
}
