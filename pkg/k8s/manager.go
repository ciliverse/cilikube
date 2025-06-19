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
	clients          map[string]*Client             // 存储所有激活的客户端实例，键为集群 ID。
	clientInfo       map[string]store.Cluster       // 存储从数据库加载的集群元数据，键为集群 ID。
	nameToID         map[string]string              // 集群名称到ID的映射。
	store            store.ClusterStore             // 数据库存储接口，用于持久化集群信息。
	statusCache      map[string]ClusterInfoResponse // 缓存各集群的健康状态和信息，键为集群 ID。
	lock             sync.RWMutex                   // 保护对 manager 内部共享资源的并发访问。
	activeClientID   string                         // 当前活动集群的 ID。
	activeClient     *Client                        // 当前活动集群的客户端实例。
}

// NewClusterManager 创建并初始化一个新的 ClusterManager 实例。
func NewClusterManager(clusterStore store.ClusterStore, config *configs.Config) (*ClusterManager, error) {
	manager := &ClusterManager{
		clients:     make(map[string]*Client),
		clientInfo:  make(map[string]store.Cluster),
		nameToID:    make(map[string]string),
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
				manager.addClient(cluster.ID, cluster.Name, cluster.KubeconfigData, "database", cluster.Environment)
				manager.clientInfo[cluster.ID] = cluster
				manager.nameToID[cluster.Name] = cluster.ID
			}
		}
	} else {
		log.Println("提示: 数据库存储未初始化，跳过从数据库加载集群。")
	}

	// 2. [兼容逻辑] 从 config.yaml 文件加载静态配置的集群。
	if len(config.Clusters) > 0 {
		for _, clusterInfo := range config.Clusters {
            // 文件类型集群的 ID 可以是其名称，或者生成一个UUID
            // 这里为了简化，我们暂时使用名称作为ID，但理想情况下应该有独立的ID机制
            fileID := clusterInfo.Name 
			if _, exists := manager.clients[fileID]; exists {
				continue // 跳过已从数据库加载的同名集群
			}
            if _, nameExists := manager.nameToID[clusterInfo.Name]; nameExists {
                log.Printf("警告: 文件集群 '%s' 与已从数据库加载的集群名称冲突，跳过。", clusterInfo.Name)
                continue
            }
			// 对于文件类型，我们直接传递路径，由 NewClient 内部处理
			manager.addClient(fileID, clusterInfo.Name, nil, "file", "default", clusterInfo.ConfigPath)
            manager.nameToID[clusterInfo.Name] = fileID
		}
	}

	// 3. 启动一个后台 goroutine，异步刷新所有集群的状态。
	go manager.startStatusUpdater()

	// 4. 根据配置文件设置活动集群。
	if config.Server.ActiveClusterID != "" {
		if _, exists := manager.clients[config.Server.ActiveClusterID]; exists {
			if err := manager.SetActiveClusterByID(config.Server.ActiveClusterID); err != nil {
				log.Printf("警告: 无法将配置文件中的活动集群 ID '%s' 设置为活动集群: %v", config.Server.ActiveClusterID, err)
			}
		} else {
		    log.Printf("警告: 配置文件中指定的活动集群 ID '%s' 未在已加载的集群中找到。", config.Server.ActiveClusterID)
		}
	} else if len(manager.clients) > 0 {
		// 尝试设置第一个加载的集群为活动集群
		for id := range manager.clients {
			if err := manager.SetActiveClusterByID(id); err == nil {
				log.Printf("提示: 未指定活动集群，已默认设置为 ID '%s' (名称: '%s')。", id, manager.clientInfo[id].Name)
				break
			}
		}
	}

	log.Printf("集群管理器初始化完成，共加载 %d 个客户端。活动集群ID: '%s'", len(manager.clients), manager.GetActiveClusterID())
	return manager, nil
}

// addClient 是一个内部辅助函数，用于添加客户端并初始化状态缓存。
// configPath 仅用于文件类型的集群。
func (cm *ClusterManager) addClient(id, name string, kubeconfigData []byte, source, environment string, configPath ...string) {
	// 注意：此方法不再直接加锁，调用方负责锁管理，因为它可能在已持有锁的上下文中被调用
	var client *Client
	var err error

	if source == "database" {
		client, err = NewClientFromContent(kubeconfigData)
	} else if source == "file" && len(configPath) > 0 {
		client, err = NewClient(configPath[0])
	} else {
		err = fmt.Errorf("无效的 addClient 调用 for ID %s, Name %s", id, name)
		log.Printf("%v", err)
		// 即使创建失败，也记录一个条目，以便UI可以显示错误状态
		cm.statusCache[id] = ClusterInfoResponse{
			Name:        name, // 仍然使用 Name 进行显示
			Status:      fmt.Sprintf("初始化失败: %v", err),
			Source:      source,
			Environment: environment,
		}
		return
	}

	if err != nil {
		log.Printf("警告: 为集群 '%s' (ID: %s) 创建客户端失败: %v", name, id, err)
		cm.statusCache[id] = ClusterInfoResponse{
			Name:        name,
			Status:      fmt.Sprintf("初始化失败: %v", err),
			Source:      source,
			Environment: environment,
		}
		return
	}
	cm.clients[id] = client
	cm.statusCache[id] = ClusterInfoResponse{
		Name:        name, // UI层面仍然显示Name
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
	clientsToUpdate := make(map[string]*Client) // Key is ID
	clientInfoSnapshot := make(map[string]store.Cluster) // Key is ID
	for id, client := range cm.clients {
		clientsToUpdate[id] = client
		clientInfoSnapshot[id] = cm.clientInfo[id] // Capture name for logging/DB ops
	}
	cm.lock.RUnlock()

	var wg sync.WaitGroup
	for id, client := range clientsToUpdate {
		wg.Add(1)
		go func(id string, client *Client, clusterInfo store.Cluster) {
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
			cachedInfo := cm.statusCache[id]
			cachedInfo.Status = status
			cachedInfo.Version = version
			// Name is already set in statusCache during addClient
			cm.statusCache[id] = cachedInfo
			cm.lock.Unlock()

			// 如果是数据库来源的集群，则更新数据库中的版本信息
			if cachedInfo.Source == "database" && cm.store != nil {
				// 使用 GetClusterByID 进行查找和更新
				dbCluster, err := cm.store.GetClusterByID(id)
				if err == nil && dbCluster.Version != version {
					dbCluster.Version = version
					cm.store.UpdateCluster(dbCluster)
				}
			}
		}(id, client, clientInfoSnapshot[id])
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

	// 检查名称是否已存在
	if existingID, nameExists := cm.nameToID[cluster.Name]; nameExists {
		// 如果名称存在但ID不同，则不允许（名称应唯一）
		// 如果ID也相同，则可能是重复添加，也应处理
		return fmt.Errorf("集群名称 '%s' 已存在，对应ID: %s", cluster.Name, existingID)
	}

	if err := cm.store.CreateCluster(cluster); err != nil {
		return fmt.Errorf("保存集群到数据库失败: %w", err)
	}

	// 使用内部方法添加客户端，此时 cluster.ID 应该已经被 store.CreateCluster 填充
	cm.addClient(cluster.ID, cluster.Name, cluster.KubeconfigData, "database", cluster.Environment)
	cm.clientInfo[cluster.ID] = *cluster
	cm.nameToID[cluster.Name] = cluster.ID

	// 立即触发一次状态刷新
	go cm.RefreshAllClusterStatus()
	return nil
}

// RemoveDBClusterByName 从数据库和内存中移除一个集群（按名称）。
// Deprecated: 优先使用 RemoveDBClusterByID。
func (cm *ClusterManager) RemoveDBClusterByName(name string) error {
	cm.lock.Lock()
	id, exists := cm.nameToID[name]
	cm.lock.Unlock() // 在调用 RemoveDBClusterByID 之前解锁，因为它会再次加锁
	if !exists {
		return fmt.Errorf("集群名称 '%s' 未找到", name)
	}
	return cm.RemoveDBClusterByID(id)
}

// RemoveDBClusterByID 从数据库和内存中移除一个集群（按ID）。
func (cm *ClusterManager) RemoveDBClusterByID(id string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	clientInfo, clientInfoExists := cm.clientInfo[id]
	if !clientInfoExists {
		return fmt.Errorf("集群ID '%s' 未找到", id)
	}

	if info, ok := cm.statusCache[id]; ok && info.Source == "file" {
		return fmt.Errorf("集群 '%s' (ID: %s) 是一个基于文件的集群，无法通过API删除", clientInfo.Name, id)
	}

	if cm.store == nil {
		return fmt.Errorf("数据库未初始化，无法删除集群 '%s' (ID: %s)", clientInfo.Name, id)
	}

	if err := cm.store.DeleteClusterByID(id); err != nil {
		return fmt.Errorf("从数据库删除集群 '%s' (ID: %s) 失败: %w", clientInfo.Name, id, err)
	}

	delete(cm.clients, id)
	delete(cm.statusCache, id)
	delete(cm.clientInfo, id)
	delete(cm.nameToID, clientInfo.Name)

	if cm.activeClientID == id {
		cm.activeClient = nil
		cm.activeClientID = ""
		// 尝试选择另一个集群作为活动集群
		for newActiveID := range cm.clients {
			cm.activeClientID = newActiveID
			cm.activeClient = cm.clients[newActiveID]
			log.Printf("活动集群已切换到 ID: %s (名称: %s)", newActiveID, cm.clientInfo[newActiveID].Name)
			break
		}
	}
	return nil
}

// SetActiveClusterByName 切换当前活动的集群（按名称）。
// Deprecated: 优先使用 SetActiveClusterByID。
func (cm *ClusterManager) SetActiveClusterByName(name string) error {
	cm.lock.RLock()
	id, exists := cm.nameToID[name]
	cm.lock.RUnlock()
	if !exists {
		return fmt.Errorf("集群名称 '%s' 未找到", name)
	}
	return cm.SetActiveClusterByID(id)
}

// SetActiveClusterByID 切换当前活动的集群（按ID）。
func (cm *ClusterManager) SetActiveClusterByID(id string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	client, exists := cm.clients[id]
	if !exists {
		// 尝试从 clientInfo 中获取名称以提供更友好的错误消息
		clusterName := "<未知>"
		if info, ok := cm.clientInfo[id]; ok {
			clusterName = info.Name
		}
		return fmt.Errorf("集群ID '%s' (名称: %s) 未找到或未初始化", id, clusterName)
	}

	cm.activeClient = client
	cm.activeClientID = id
	log.Printf("活动集群已设置为 ID: %s (名称: %s)", id, cm.clientInfo[id].Name)
	return nil
}

// GetActiveClient 返回当前活动集群的客户端实例。
func (cm *ClusterManager) GetActiveClient() (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	if cm.activeClient == nil || cm.activeClientID == "" {
		return nil, fmt.Errorf("当前没有配置或可用的活动集群")
	}
	return cm.activeClient, nil
}

// GetActiveClusterName 返回当前活动集群的名称。
func (cm *ClusterManager) GetActiveClusterName() string {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	if cm.activeClientID == "" {
		return ""
	}
	if info, ok := cm.clientInfo[cm.activeClientID]; ok {
		return info.Name
	}
	return ""
}

// GetActiveClusterID 返回当前活动集群的ID。
func (cm *ClusterManager) GetActiveClusterID() string {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	return cm.activeClientID
}

// GetClientByName 根据指定的名称从内存中获取一个集群的客户端实例。
// Deprecated: 优先使用 GetClientByID。
func (cm *ClusterManager) GetClientByName(name string) (*Client, error) {
	cm.lock.RLock()
	id, exists := cm.nameToID[name]
	cm.lock.RUnlock()
	if !exists {
		return nil, fmt.Errorf("集群名称 '%s' 未找到对应的ID", name)
	}
	return cm.GetClientByID(id)
}

// GetClientByID 根据指定的ID从内存中获取一个集群的客户端实例。
func (cm *ClusterManager) GetClientByID(id string) (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	client, exists := cm.clients[id]
	if !exists {
		// 尝试从 clientInfo 中获取名称以提供更友好的错误消息
		clusterName := "<未知>"
		if info, ok := cm.clientInfo[id]; ok {
			clusterName = info.Name
		}
		return nil, fmt.Errorf("ID为 '%s' (名称: %s) 的客户端未在内存中找到", id, clusterName)
	}
	return client, nil
}
