package service

import (
	"encoding/base64"
	"fmt"

	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

// ClusterService 提供了围绕集群管理的业务逻辑。
// 它是 API handlers 和底层 k8s.ClusterManager 之间的桥梁。
type ClusterService struct {
	k8sManager *k8s.ClusterManager
}

// NewClusterService 创建一个新的 ClusterService 实例。
func NewClusterService(k8sManager *k8s.ClusterManager) *ClusterService {
	return &ClusterService{
		k8sManager: k8sManager,
	}
}

// CreateClusterRequest 定义了创建集群API的请求体结构。
// 这是一个数据传输对象 (DTO)，用于 API 层的输入。
type CreateClusterRequest struct {
	Name           string `json:"name" binding:"required"`
	KubeconfigData string `json:"kubeconfigData" binding:"required"` // Base64 编码的 kubeconfig 字符串
	Provider       string `json:"provider"`
	Description    string `json:"description"`
}

// ListClusters 返回所有受管集群的信息列表。
func (s *ClusterService) ListClusters() []k8s.ClusterInfoResponse {
	return s.k8sManager.ListClusterInfo()
}

// CreateCluster 处理创建新集群的逻辑。
func (s *ClusterService) CreateCluster(req CreateClusterRequest) error {
	// 解码 Base64 格式的 kubeconfig 数据。
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(req.KubeconfigData)
	if err != nil {
		return fmt.Errorf("kubeconfig 数据不是有效的 Base64 编码: %w", err)
	}

	// 将请求数据映射到数据库模型。
	cluster := &store.Cluster{
		Name:           req.Name,
		KubeconfigData: kubeconfigBytes,
		Provider:       req.Provider,
		Description:    req.Description,
	}

	return s.k8sManager.AddDBCluster(cluster)
}

// DeleteCluster 处理删除集群的逻辑（按名称）。
// Deprecated: 优先使用 DeleteClusterByID。
func (s *ClusterService) DeleteCluster(name string) error {
	return s.k8sManager.RemoveDBClusterByName(name) // 假设 k8sManager 有 RemoveDBClusterByName
}

// DeleteClusterByID 处理删除集群的逻辑（按ID）。
func (s *ClusterService) DeleteClusterByID(id string) error {
	return s.k8sManager.RemoveDBClusterByID(id) // 假设 k8sManager 有 RemoveDBClusterByID
}

// SetActiveCluster 处理切换活动集群的逻辑（按名称）。
// Deprecated: 优先使用 SetActiveClusterByID。
func (s *ClusterService) SetActiveCluster(name string) error {
	return s.k8sManager.SetActiveClusterByName(name) // 假设 k8sManager 有 SetActiveClusterByName
}

// SetActiveClusterByID 处理切换活动集群的逻辑（按ID）。
func (s *ClusterService) SetActiveClusterByID(id string) error {
	return s.k8sManager.SetActiveClusterByID(id) // 假设 k8sManager 有 SetActiveClusterByID
}

// GetActiveCluster 获取当前活动集群的名称。
func (s *ClusterService) GetActiveCluster() string {
	// TODO: 考虑是否也需要 GetActiveClusterID() 或让此方法返回更丰富的集群信息对象
	return s.k8sManager.GetActiveClusterName()
}

// GetActiveClusterID 获取当前活动集群的ID。
func (s *ClusterService) GetActiveClusterID() string {
	return s.k8sManager.GetActiveClusterID() // 假设 k8sManager 有 GetActiveClusterID
}
