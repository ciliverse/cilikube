package k8s

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/models"
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
	informers      *informerState
}

func NewClusterManager(clusterStore store.ClusterStore, config *configs.Config) (*ClusterManager, error) {
	manager := &ClusterManager{
		clients:     make(map[string]*Client),
		clientInfo:  make(map[string]store.Cluster),
		nameToID:    make(map[string]string),
		store:       clusterStore,
		statusCache: make(map[string]ClusterInfoResponse),
		informers:   newInformerState(),
	}
	log.Println("initializing cluster manager...")

	if IsShowcase() {
		// Public exhibit: in-memory fake cluster only — never dial a real apiserver / k3s.
		log.Println("CILIKUBE_SHOWCASE=1 — registering simulated demo cluster (no real kubeconfig)")
		manager.registerShowcaseCluster()
	} else if clusterStore != nil {
		dbClusters, err := clusterStore.GetAllClusters()
		if err != nil {
			log.Printf("warning: failed to load clusters from database: %v", err)
		} else {
			for _, cluster := range dbClusters {
				manager.addClient(cluster.ID, cluster.Name, cluster.KubeconfigData, "database", cluster.Environment, "")
				manager.clientInfo[cluster.ID] = cluster
				manager.nameToID[cluster.Name] = cluster.ID
			}
		}
	} else {
		log.Println("Note: Database storage not initialized, skipping loading clusters from database.")
	}

	// File clusters in config.yaml are deprecated. Prefer DB registration via UI
	// paste or "import from local kubeconfig". Keeping the section only logs a hint.
	if len(config.Clusters) > 0 {
		log.Printf("Note: config.yaml lists %d cluster(s) under clusters: — ignored. Add clusters in the UI or import local kubeconfig contexts.", len(config.Clusters))
	}

	go manager.startStatusUpdater()

	if IsShowcase() {
		// already set active to demo
	} else if config.Server.ActiveClusterID != "" {
		if err := manager.SetActiveClusterByID(config.Server.ActiveClusterID); err != nil {
			log.Printf("Warning: Unable to set active cluster ID '%s' from config file as active: %v", config.Server.ActiveClusterID, err)
		}
	} else if len(manager.clients) > 0 {
		for id := range manager.clients {
			if err := manager.SetActiveClusterByID(id); err == nil {
				break
			}
		}
	}

	log.Printf("Cluster manager initialization completed, loaded %d clients. Active cluster ID: '%s'", len(manager.clients), manager.GetActiveClusterID())
	return manager, nil
}

func (cm *ClusterManager) addClient(id, name string, kubeconfigData []byte, source, environment string, configPath string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.addClientLocked(id, name, kubeconfigData, source, environment, configPath)
}

// addClientLocked registers a client; caller must hold cm.lock.
func (cm *ClusterManager) addClientLocked(id, name string, kubeconfigData []byte, source, environment string, configPath string) {
	var client *Client
	var err error
	if source == "database" {
		client, err = NewClientFromContent(kubeconfigData)
	} else if source == "file" {
		client, err = NewClient(configPath)
	} else {
		err = fmt.Errorf("invalid addClient call for ID %s", id)
	}

	if err != nil {
		log.Printf("Warning: Failed to create client for cluster '%s' (ID: %s): %v", name, id, err)
		cm.statusCache[id] = ClusterInfoResponse{ID: id, Name: name, Status: fmt.Sprintf("Initialization failed: %v", err), Source: source}
		return
	}
	cm.clients[id] = client
	cm.statusCache[id] = ClusterInfoResponse{
		ID:          id,
		Name:        name,
		Server:      client.Config.Host,
		Status:      "Checking...",
		Source:      source,
		Environment: environment,
	}
}

func (cm *ClusterManager) startStatusUpdater() {
	time.Sleep(5 * time.Second)
	log.Println("Performing initial cluster status check...")
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
			if IsShowcaseConfig(client.Config) {
				status = "Available"
				version = ShowcaseVersion
				if client.clusterInfo != nil && client.clusterInfo.ServerVersion != "" {
					version = client.clusterInfo.ServerVersion
				}
			} else if serverVersion, err := client.Clientset.Discovery().ServerVersion(); err != nil {
				status = fmt.Sprintf("Unavailable: %v", err)
				version = "N/A"
			} else {
				status = "Available"
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
	if cm.store == nil {
		cm.lock.Unlock()
		return fmt.Errorf("cluster store not initialized, cannot add cluster")
	}
	if _, nameExists := cm.nameToID[cluster.Name]; nameExists {
		cm.lock.Unlock()
		return fmt.Errorf("cluster name '%s' already exists", cluster.Name)
	}
	if err := cm.store.CreateCluster(cluster); err != nil {
		cm.lock.Unlock()
		return fmt.Errorf("failed to save cluster: %w", err)
	}
	// Use locked helper — we already hold cm.lock (addClient would deadlock).
	cm.addClientLocked(cluster.ID, cluster.Name, cluster.KubeconfigData, "database", cluster.Environment, "")
	cm.clientInfo[cluster.ID] = *cluster
	cm.nameToID[cluster.Name] = cluster.ID
	cm.lock.Unlock()
	// Refresh outside the write lock so create API is not blocked by status probes.
	go cm.RefreshAllClusterStatus()
	return nil
}

func (cm *ClusterManager) RemoveDBClusterByID(id string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	clientInfo, clientInfoExists := cm.clientInfo[id]
	if !clientInfoExists {
		return fmt.Errorf("cluster ID '%s' not found", id)
	}
	if info, ok := cm.statusCache[id]; ok && info.Source == "file" {
		return fmt.Errorf("cluster '%s' (ID: %s) is a file-based cluster, cannot be deleted via API", clientInfo.Name, id)
	}
	if cm.store == nil {
		return fmt.Errorf("cluster store not initialized, cannot delete cluster '%s' (ID: %s)", clientInfo.Name, id)
	}
	if err := cm.store.DeleteClusterByID(id); err != nil {
		return fmt.Errorf("failed to delete cluster '%s' (ID: %s): %w", clientInfo.Name, id, err)
	}
	delete(cm.clients, id)
	delete(cm.statusCache, id)
	delete(cm.clientInfo, id)
	delete(cm.nameToID, clientInfo.Name)
	cm.stopInformersForCluster(id)
	if cm.activeClientID == id {
		cm.activeClient = nil
		cm.activeClientID = ""
		for newActiveID, newClient := range cm.clients {
			cm.activeClient = newClient
			cm.activeClientID = newActiveID
			log.Printf("Active cluster set to ID: %s (name: %s)", newActiveID, cm.clientInfo[newActiveID].Name)
			go cm.startInformersForCluster(newActiveID, newClient)
			break
		}
	}
	return nil
}

func (cm *ClusterManager) SetActiveClusterByID(id string) error {
	cm.lock.Lock()
	client, exists := cm.clients[id]
	if !exists {
		cm.lock.Unlock()
		return fmt.Errorf("cluster ID '%s' not found or not initialized", id)
	}
	prevID := cm.activeClientID
	cm.activeClient = client
	cm.activeClientID = id
	name := cm.clientInfo[id].Name
	cm.lock.Unlock()

	log.Printf("Active cluster set to ID: %s (name: %s)", id, name)
	if prevID != id {
		if prevID != "" {
			cm.stopInformersForCluster(prevID)
		}
		cm.startInformersForCluster(id, client)
	} else {
		// Ensure informers are running even if active cluster unchanged (e.g. startup)
		cm.informers.mu.RLock()
		_, running := cm.informers.byCluster[id]
		cm.informers.mu.RUnlock()
		if !running {
			cm.startInformersForCluster(id, client)
		}
	}
	return nil
}

func (cm *ClusterManager) GetActiveClient() (*Client, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	if cm.activeClient == nil {
		return nil, fmt.Errorf("no active cluster currently configured or available")
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
		return nil, fmt.Errorf("client with ID '%s' not found in memory", id)
	}
	return client, nil
}

func (cm *ClusterManager) GetClusterDetailFromDB(id string) (*store.Cluster, error) {
	if cm.store == nil {
		return nil, fmt.Errorf("cluster store not initialized")
	}
	return cm.store.GetClusterByID(id)
}

func (cm *ClusterManager) UpdateDBCluster(id string, req models.UpdateClusterRequest) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if cm.store == nil {
		return fmt.Errorf("cluster store not initialized, cannot update cluster")
	}
	cluster, err := cm.store.GetClusterByID(id)
	if err != nil {
		return fmt.Errorf("cluster ID '%s' not found: %w", id, err)
	}
	oldName := cluster.Name
	kubeconfigUpdated := false
	if req.Name != "" {
		cluster.Name = req.Name
	}
	if req.Provider != "" {
		cluster.Provider = req.Provider
	}
	if req.Description != "" {
		cluster.Description = req.Description
	}
	if req.Environment != nil {
		cluster.Environment = *req.Environment
	}
	if req.Region != "" {
		cluster.Region = req.Region
	}
	if req.Status != "" {
		cluster.Status = req.Status
	}
	if req.Labels != nil {
		cluster.Labels = store.Labels(req.Labels)
	}
	if req.KubeconfigData != "" {
		kubeconfigBytes, err := base64.StdEncoding.DecodeString(req.KubeconfigData)
		if err != nil {
			return fmt.Errorf("kubeconfig data is not valid Base64 encoding: %w", err)
		}
		cluster.KubeconfigData = kubeconfigBytes
		kubeconfigUpdated = true
	}
	if err := cm.store.UpdateCluster(cluster); err != nil {
		return fmt.Errorf("failed to update cluster: %w", err)
	}
	cm.clientInfo[id] = *cluster
	if oldName != cluster.Name {
		delete(cm.nameToID, oldName)
		cm.nameToID[cluster.Name] = id
	}
	if kubeconfigUpdated {
		delete(cm.clients, id)
		delete(cm.statusCache, id)
		cm.addClientLocked(id, cluster.Name, cluster.KubeconfigData, "database", cluster.Environment, "")
		go cm.RefreshAllClusterStatus()
	} else if cached, ok := cm.statusCache[id]; ok {
		cached.Name = cluster.Name
		cached.Environment = cluster.Environment
		cm.statusCache[id] = cached
	}
	return nil
}

// GetStatusFromCache gets cluster status information from memory cache
func (cm *ClusterManager) GetStatusFromCache(id string) (ClusterInfoResponse, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	info, ok := cm.statusCache[id]
	return info, ok
}
