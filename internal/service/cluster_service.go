package service

import (
	"encoding/base64"
	"fmt"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

// ClusterService 提供了围绕集群管理的业务逻辑。
type ClusterService struct {
	k8sManager *k8s.ClusterManager
}

// NewClusterService 创建一个新的 ClusterService 实例。
func NewClusterService(k8sManager *k8s.ClusterManager) *ClusterService {
	return &ClusterService{
		k8sManager: k8sManager,
	}
}

// ListClusters 返回所有受管集群的概要信息列表。
func (s *ClusterService) ListClusters() []models.ClusterListResponse {
	// k8sManager 返回的信息结构已经很适合列表页，我们直接转换一下
	managerInfo := s.k8sManager.ListClusterInfo()
	response := make([]models.ClusterListResponse, len(managerInfo))
	for i, info := range managerInfo {
		response[i] = models.ClusterListResponse{
			ID:          info.ID, // 确保 k8s.ClusterInfoResponse 中有 ID 字段
			Name:        info.Name,
			Server:      info.Server,
			Version:     info.Version,
			Status:      info.Status,
			Source:      info.Source,
			Environment: info.Environment,
		}
	}
	return response
}

// GetClusterByID 获取单个集群的详细信息。
func (s *ClusterService) GetClusterByID(id string) (*models.ClusterResponse, error) {
	cluster, err := s.k8sManager.GetClusterDetailFromDB(id)
	if err != nil {
		// 如果数据库中没有，可能是文件类型的集群，我们从缓存中组装一个简版
		if info, ok := s.k8sManager.GetStatusFromCache(id); ok {
			return &models.ClusterResponse{
				ID:          info.ID,
				Name:        info.Name,
				Version:     info.Version,
				Status:      info.Status,
				Environment: info.Environment,
				Source:      info.Source,
			}, nil
		}
		return nil, fmt.Errorf("集群 ID '%s' 未找到: %w", id, err)
	}

	return &models.ClusterResponse{
		ID:          cluster.ID,
		Name:        cluster.Name,
		Provider:    cluster.Provider,
		Description: cluster.Description,
		Environment: cluster.Environment,
		Region:      cluster.Region,
		Version:     cluster.Version,
		Status:      cluster.Status,
		Labels:      cluster.Labels,
		CreatedAt:   cluster.CreatedAt,
		UpdatedAt:   cluster.UpdatedAt,
	}, nil
}

// CreateCluster 处理创建新集群的逻辑。
func (s *ClusterService) CreateCluster(req models.CreateClusterRequest) error {
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(req.KubeconfigData)
	if err != nil {
		return fmt.Errorf("kubeconfig 数据不是有效的 Base64 编码: %w", err)
	}
	cluster := &store.Cluster{
		Name:           req.Name,
		KubeconfigData: kubeconfigBytes,
		Provider:       req.Provider,
		Description:    req.Description,
		Environment:    req.Environment,
		Region:         req.Region,
	}
	return s.k8sManager.AddDBCluster(cluster)
}

// UpdateCluster 更新集群信息。
func (s *ClusterService) UpdateCluster(id string, req models.UpdateClusterRequest) error {
	return s.k8sManager.UpdateDBCluster(id, req)
}

// DeleteClusterByID 处理删除集群的逻辑。
func (s *ClusterService) DeleteClusterByID(id string) error {
	return s.k8sManager.RemoveDBClusterByID(id)
}

// SetActiveCluster 处理切换活动集群的逻辑。
func (s *ClusterService) SetActiveCluster(id string) error {
	return s.k8sManager.SetActiveClusterByID(id)
}
