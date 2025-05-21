package k8s

import (
	"fmt"
	"sync"
)

// MultiClient 管理多个集群的 clientset 和 config
type MultiClient struct {
	mu      sync.RWMutex
	clients map[string]*Client // key: 集群标识
}

// NewMultiClient 初始化 MultiClient
func NewMultiClient() *MultiClient {
	return &MultiClient{
		clients: make(map[string]*Client),
	}
}

// AddCluster 添加或更新一个集群的客户端配置
// name: 集群名称, kubeconfigPath: kubeconfig 文件路径，空则使用 in-cluster
func (m *MultiClient) AddCluster(name, kubeconfigPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 创建单集群 Client
	client, err := NewClient(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("[MultiClient] 添加集群 %s 失败: %w", name, err)
	}
	m.clients[name] = client
	return nil
}

// RemoveCluster 删除指定集群
func (m *MultiClient) RemoveCluster(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, name)
}

// GetClient 获取指定集群 client
func (m *MultiClient) GetClient(name string) (*Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.clients[name]
	if !ok {
		return nil, fmt.Errorf("集群 %s 未注册", name)
	}
	return client, nil
}

// ListClusters 返回所有已注册的集群名称
func (m *MultiClient) ListClusters() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}
