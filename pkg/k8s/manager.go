package k8s

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
)

// ClusterInfoResponse 定义了 ListClusterInfo 方法的返回结构体，以增强类型安全。
type ClusterInfoResponse struct {
	Name        string `json:"name"`
	Server      string `json:"server"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	Source      string `json:"source"` // 集群来源: "database" 或 "file"
	Environment string `json:"environment"`
}

// ClusterManager 负责管理所有 Kubernetes 集群的客户端实例和状态。
type ClusterManager struct {
	clients          map[string]*Client             // 存储所有激活的客户端实例，键为集群名称。
	clientInfo       map[string]store.Cluster       // 存储从数据库加载的集群元数据，用于快速访问。
	store            store.ClusterStore             // 数据库存储接口，用于持久化集群信息。
	statusCache      map[string]ClusterInfoResponse // 缓存各集群的健康状态和信息。
	lock             sync.RWMutex                   // 保护对 manager 内部共享资源的并发访问。
	activeClientName string                         // 当前活动集群的名称。
	activeClient     *Client                        // 当前活动集群的客户端实例。
}

// NewClusterManager 创建并初始化一个新的 ClusterManager 实例。
func NewClusterManager(clusterStore store.ClusterStore, config *configs.Config) (*ClusterManager, error) {
	manager := &ClusterManager{
		clients:     make(map[string]*Client),
		clientInfo:  make(map[string]store.Cluster),
		store:       clusterStore,
		statusCache: make(map[string]ClusterInfoResponse),
	}
	log.Println("正在初始化集群管理器...")

	// 1. 从数据库加载动态管理的集群。
	if clusterStore != nil {
		// 1. 从数据库加载动态管理的集群
		dbClusters, err := clusterStore.GetAllClusters()
		if err != nil {
			log.Printf("警告: 从数据库加载集群失败: %v", err)
		} else {
			for _, cluster := range dbClusters {
				manager.addClient(cluster.Name, cluster.KubeconfigData, "database", cluster.Environment)
				manager.clientInfo[cluster.Name] = cluster
			}
		}
	} else {
		log.Println("提示: 数据库存储未初始化，跳过从数据库加载集群。")
	}

	// 2. [兼容逻辑] 从 config.yaml 文件加载静态配置的集群。
	if len(config.Clusters) > 0 {
		for _, clusterInfo := range config.Clusters {
			if _, exists := manager.clients[clusterInfo.Name]; exists {
				continue // 跳过已从数据库加载的同名集群
			}
			// 对于文件类型，我们直接传递路径，由 NewClient 内部处理
			manager.addClient(clusterInfo.Name, nil, "file", "default", clusterInfo.ConfigPath)
		}
	}

	// 3. 启动一个后台 goroutine，异步刷新所有集群的状态。
	go manager.startStatusUpdater()

	// 4. 根据配置文件设置活动集群。
	if config.Server.ActiveCluster != "" {
		if err := manager.SetActiveCluster(config.Server.ActiveCluster); err != nil {
			log.Printf("警告: 无法将配置文件中的 '%s' 设置为活动集群: %v", config.Server.ActiveCluster, err)
		}
	} else if len(manager.clients) > 0 {
		for name := range manager.clients {
			if err := manager.SetActiveCluster(name); err == nil {
				log.Printf("提示: 未指定活动集群，已默认设置为 '%s'。", name)
				break
			}
		}
	}

	log.Printf("集群管理器初始化完成，共加载 %d 个客户端。活动集群: '%s'", len(manager.clients), manager.GetActiveClusterName())
	return manager, nil
}

// addClient 是一个内部辅助函数，用于添加客户端并初始化状态缓存。
// 可变参数 configPath 仅用于文件类型的集群。
func (cm *ClusterManager) addClient(name string, kubeconfigData []byte, source, environment string, configPath ...string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	var client *Client
	var err error

	if source == "database" {
		client, err = NewClientFromContent(kubeconfigData)
	} else if source == "file" && len(configPath) > 0 {
		client, err = NewClient(configPath[0])
	} else {
		err = fmt.Errorf("无效的 addClient 调用")
	}

	if err != nil {
		log.Printf("警告: 为集群 '%s' 创建客户端失败: %v", name, err)
		cm.statusCache[name] = ClusterInfoResponse{
			Name:        name,
			Status:      fmt.Sprintf("初始化失败: %v", err),
			Source:      source,
			Environment: environment,
		}
		return
	}
	cm.clients[name] = client
	cm.statusCache[name] = ClusterInfoResponse{
		Name:        name,
		Server:      client.Config.Host,
		Status:      "检查中...",
		Source:      source,
		Environment: environment,
	}
}

// startStatusUpdater 在后台定期刷新集群状态。
func (cm *ClusterManager) startStatusUpdater() {
	time.Sleep(5 * time.Second) // 延迟首次执行，等待服务完全启动
	log.Println("正在执行首次集群状态检查...")
	cm.RefreshAllClusterStatus()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		cm.RefreshAllClusterStatus()
	}
}

// RefreshAllClusterStatus 并发地检查所有集群的健康状况。
func (cm *ClusterManager) RefreshAllClusterStatus() {
	cm.lock.RLock()
	clientsToUpdate := make(map[string]*Client)
	for name, client := range cm.clients {
		clientsToUpdate[name] = client
	}
	cm.lock.RUnlock()

	var wg sync.WaitGroup
	for name, client := range clientsToUpdate {
		wg.Add(1)
		go func(name string, client *Client) {
			defer wg.Done()
			var status, version string
			serverVersion, err := client.Clientset.Discovery().ServerVersion()
			if err != nil {
				status = fmt.Sprintf("不可用: %v", err)
				version = "N/A"
			} else {
				status = "可用"
				version = serverVersion.GitVersion
			}

			cm.lock.Lock()
			cachedInfo := cm.statusCache[name]
			cachedInfo.Status = status
			cachedInfo.Version = version
			cm.statusCache[name] = cachedInfo
			cm.lock.Unlock()

			// 如果是数据库来源的集群，则更新数据库中的版本信息
			if cachedInfo.Source == "database" && cm.store != nil {
				cluster, err := cm.store.GetClusterByName(name)
				if err == nil && cluster.Version != version {
					cluster.Version = version
					cm.store.UpdateCluster(cluster)
				}
			}
		}(name, client)
	}
	wg.Wait()
	log.Println("集群状态刷新完成。")
}

// ListClusterInfo 返回所有受管集群的信息列表。
func (cm *ClusterManager) ListClusterInfo() []ClusterInfoResponse {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	clusterList := make([]ClusterInfoResponse, 0, len(cm.statusCache))
	for _, info := range cm.statusCache {
		clusterList = append(clusterList, info)
	}
	return clusterList
}

// AddDBCluster 将新集群信息存入数据库并激活。
func (cm *ClusterManager) AddDBCluster(cluster *store.Cluster) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if cm.store == nil {
		return fmt.Errorf("数据库未初始化，无法添加集群")
	}

	if err := cm.store.CreateCluster(cluster); err != nil {
		return fmt.Errorf("保存集群到数据库失败: %w", err)
	}

	// 使用内部方法添加客户端
	cm.addClient(cluster.Name, cluster.KubeconfigData, "database", cluster.Environment)
	// 立即触发一次状态刷新
	go cm.RefreshAllClusterStatus()
	return nil
}

// RemoveDBCluster 从数据库和内存中移除一个集群。
func (cm *ClusterManager) RemoveDBCluster(name string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	if _, exists := cm.clients[name]; !exists {
		return fmt.Errorf("集群 '%s' 未找到", name)
	}

	if info, ok := cm.statusCache[name]; ok && info.Source == "file" {
		return fmt.Errorf("集群 '%s' 是一个基于文件的集群，无法通过API删除", name)
	}

	if err := cm.store.DeleteClusterByName(name); err != nil {
		return fmt.Errorf("从数据库删除集群失败: %w", err)
	}

	delete(cm.clients, name)
	delete(cm.statusCache, name)
	delete(cm.clientInfo, name)

	if cm.activeClientName == name {
		cm.activeClient = nil
		cm.activeClientName = ""
	}
	return nil
}

// SetActiveCluster 切换当前活动的集群。
func (cm *ClusterManager) SetActiveCluster(name string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	client, exists := cm.clients[name]
	if !exists {
		return fmt.Errorf("集群 '%s' 未找到或未初始化", name)
	}

	cm.activeClient = client
	cm.activeClientName = name
	return nil
}

// GetActiveClient 返回当前活动集群的客户端实例。
func (cm *ClusterManager) GetActiveClient() (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	if cm.activeClient == nil {
		return nil, fmt.Errorf("当前没有配置或可用的活动集群")
	}
	return cm.activeClient, nil
}

// GetActiveClusterName 返回当前活动集群的名称。
func (cm *ClusterManager) GetActiveClusterName() string {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	return cm.activeClientName
}

// GetClient 根据指定的名称从内存中获取一个集群的客户端实例。
func (cm *ClusterManager) GetClient(name string) (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	client, exists := cm.clients[name]
	if !exists {
		return nil, fmt.Errorf("名为 '%s' 的客户端未在内存中找到", name)
	}
	return client, nil
}
