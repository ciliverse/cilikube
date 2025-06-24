package k8s

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
)

type ClusterInfoResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Server      string `json:"server"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	Source      string `json:"source"`
	Environment string `json:"environment"`
}

type ClusterManager struct {
	clients        map[string]*Client
	clientInfo     map[string]store.Cluster
	nameToID       map[string]string
	store          store.ClusterStore
	statusCache    map[string]ClusterInfoResponse
	lock           sync.RWMutex
	activeClientID string
	activeClient   *Client
}

func NewClusterManager(clusterStore store.ClusterStore, config *configs.Config) (*ClusterManager, error) {
	manager := &ClusterManager{
		clients:     make(map[string]*Client),
		clientInfo:  make(map[string]store.Cluster),
		nameToID:    make(map[string]string),
		store:       clusterStore,
		statusCache: make(map[string]ClusterInfoResponse),
	}
	log.Println("正在初始化集群管理器...")

	if clusterStore != nil {
		dbClusters, err := clusterStore.GetAllClusters()
		if err != nil {
			log.Printf("警告: 从数据库加载集群失败: %v", err)
		} else {
			for _, cluster := range dbClusters {
				manager.addClient(cluster.ID, cluster.Name, cluster.KubeconfigData, "database", cluster.Environment, "")
				manager.clientInfo[cluster.ID] = cluster
				manager.nameToID[cluster.Name] = cluster.ID
			}
		}
	} else {
		log.Println("提示: 数据库存储未初始化，跳过从数据库加载集群。")
	}

	if len(config.Clusters) > 0 {
		for _, clusterInfo := range config.Clusters {
			fileID := clusterInfo.Name
			if _, exists := manager.clients[fileID]; exists {
				continue
			}
			if _, nameExists := manager.nameToID[clusterInfo.Name]; nameExists {
				log.Printf("警告: 文件集群 '%s' 与已加载的集群名称冲突，跳过。", clusterInfo.Name)
				continue
			}
			manager.addClient(fileID, clusterInfo.Name, nil, "file", "default", clusterInfo.ConfigPath)
			manager.clientInfo[fileID] = store.Cluster{ID: fileID, Name: clusterInfo.Name, Provider: "file"}
			manager.nameToID[clusterInfo.Name] = fileID
		}
	}

	go manager.startStatusUpdater()

	if config.Server.ActiveClusterID != "" {
		if err := manager.SetActiveClusterByID(config.Server.ActiveClusterID); err != nil {
			log.Printf("警告: 无法将配置文件中的活动集群 ID '%s' 设置为活动: %v", config.Server.ActiveClusterID, err)
		}
	} else if len(manager.clients) > 0 {
		for id := range manager.clients {
			if err := manager.SetActiveClusterByID(id); err == nil {
				break
			}
		}
	}

	log.Printf("集群管理器初始化完成，共加载 %d 个客户端。活动集群ID: '%s'", len(manager.clients), manager.GetActiveClusterID())
	return manager, nil
}

func (cm *ClusterManager) addClient(id, name string, kubeconfigData []byte, source, environment string, configPath string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	var client *Client
	var err error
	if source == "database" {
		client, err = NewClientFromContent(kubeconfigData)
	} else if source == "file" {
		client, err = NewClient(configPath)
	} else {
		err = fmt.Errorf("无效的 addClient 调用 for ID %s", id)
	}

	if err != nil {
		log.Printf("警告: 为集群 '%s' (ID: %s) 创建客户端失败: %v", name, id, err)
		cm.statusCache[id] = ClusterInfoResponse{ID: id, Name: name, Status: fmt.Sprintf("初始化失败: %v", err), Source: source}
		return
	}
	cm.clients[id] = client
	cm.statusCache[id] = ClusterInfoResponse{
		ID:          id,
		Name:        name,
		Server:      client.Config.Host,
		Status:      "检查中...",
		Source:      source,
		Environment: environment,
	}
}

func (cm *ClusterManager) startStatusUpdater() {
	time.Sleep(5 * time.Second)
	log.Println("正在执行首次集群状态检查...")
	cm.RefreshAllClusterStatus()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		cm.RefreshAllClusterStatus()
	}
}

func (cm *ClusterManager) RefreshAllClusterStatus() {
	cm.lock.RLock()
	clientsToUpdate := make(map[string]*Client)
	for id, client := range cm.clients {
		clientsToUpdate[id] = client
	}
	cm.lock.RUnlock()

	var wg sync.WaitGroup
	for id, client := range clientsToUpdate {
		wg.Add(1)
		go func(id string, client *Client) {
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
			cm.statusCache[id] = cachedInfo
			cm.lock.Unlock()
			if cachedInfo.Source == "database" && cm.store != nil {
				dbCluster, err := cm.store.GetClusterByID(id)
				if err == nil && dbCluster.Version != version {
					dbCluster.Version = version
					_ = cm.store.UpdateCluster(dbCluster)
				}
			}
		}(id, client)
	}
	wg.Wait()
}

func (cm *ClusterManager) ListClusterInfo() []ClusterInfoResponse {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	var list []ClusterInfoResponse
	for _, info := range cm.statusCache {
		list = append(list, info)
	}
	return list
}

func (cm *ClusterManager) AddDBCluster(cluster *store.Cluster) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if cm.store == nil {
		return fmt.Errorf("数据库未初始化，无法添加集群")
	}
	if _, nameExists := cm.nameToID[cluster.Name]; nameExists {
		return fmt.Errorf("集群名称 '%s' 已存在", cluster.Name)
	}
	if err := cm.store.CreateCluster(cluster); err != nil {
		return fmt.Errorf("保存集群到数据库失败: %w", err)
	}
	cm.addClient(cluster.ID, cluster.Name, cluster.KubeconfigData, "database", cluster.Environment, "")
	cm.clientInfo[cluster.ID] = *cluster
	cm.nameToID[cluster.Name] = cluster.ID
	go cm.RefreshAllClusterStatus()
	return nil
}

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
		for newActiveID := range cm.clients {
			_ = cm.SetActiveClusterByID(newActiveID)
			break
		}
	}
	return nil
}

func (cm *ClusterManager) SetActiveClusterByID(id string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	client, exists := cm.clients[id]
	if !exists {
		return fmt.Errorf("集群ID '%s' 未找到或未初始化", id)
	}
	cm.activeClient = client
	cm.activeClientID = id
	log.Printf("活动集群已设置为 ID: %s (名称: %s)", id, cm.clientInfo[id].Name)
	return nil
}

func (cm *ClusterManager) GetActiveClient() (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	if cm.activeClient == nil {
		return nil, fmt.Errorf("当前没有配置或可用的活动集群")
	}
	return cm.activeClient, nil
}

func (cm *ClusterManager) GetActiveClusterID() string {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	return cm.activeClientID
}

func (cm *ClusterManager) GetClientByID(id string) (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	client, exists := cm.clients[id]
	if !exists {
		return nil, fmt.Errorf("ID为 '%s' 的客户端未在内存中找到", id)
	}
	return client, nil
}

func (cm *ClusterManager) GetClusterDetailFromDB(id string) (*store.Cluster, error) {
	if cm.store == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return cm.store.GetClusterByID(id)
}

func (cm *ClusterManager) UpdateDBCluster(id string, req models.UpdateClusterRequest) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if cm.store == nil {
		return fmt.Errorf("数据库未初始化，无法更新集群")
	}
	cluster, err := cm.store.GetClusterByID(id)
	if err != nil {
		return fmt.Errorf("集群 ID '%s' 未找到: %w", id, err)
	}
	oldName := cluster.Name
	kubeconfigUpdated := false
	if req.Name != "" {
		cluster.Name = req.Name
	}
	// ... 其他字段更新 ...
	if req.KubeconfigData != "" {
		kubeconfigBytes, err := base64.StdEncoding.DecodeString(req.KubeconfigData)
		if err != nil {
			return fmt.Errorf("kubeconfig 数据不是有效的 Base64 编码: %w", err)
		}
		cluster.KubeconfigData = kubeconfigBytes
		kubeconfigUpdated = true
	}
	if err := cm.store.UpdateCluster(cluster); err != nil {
		return fmt.Errorf("更新数据库失败: %w", err)
	}
	cm.clientInfo[id] = *cluster
	if oldName != cluster.Name {
		delete(cm.nameToID, oldName)
		cm.nameToID[cluster.Name] = id
	}
	if kubeconfigUpdated {
		delete(cm.clients, id)
		delete(cm.statusCache, id)
		cm.addClient(id, cluster.Name, cluster.KubeconfigData, "database", cluster.Environment, "")
		go cm.RefreshAllClusterStatus()
	}
	return nil
}

// [核心修正] 添加这个缺失的方法
// GetStatusFromCache 从内存缓存中获取集群的状态信息。
func (cm *ClusterManager) GetStatusFromCache(id string) (ClusterInfoResponse, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	info, ok := cm.statusCache[id]
	return info, ok
}
