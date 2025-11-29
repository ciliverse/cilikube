package service

import (
	"encoding/base64"
	"fmt"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

// ClusterService provides business logic around cluster management.
type ClusterService struct {
	k8sManager *k8s.ClusterManager
}

// NewClusterService creates a new ClusterService instance.
func NewClusterService(k8sManager *k8s.ClusterManager) *ClusterService {
	return &ClusterService{
		k8sManager: k8sManager,
	}
}

// ListClusters returns a list of summary information for all managed clusters.
func (s *ClusterService) ListClusters() []models.ClusterListResponse {
	// The information structure returned by k8sManager is already suitable for the list page, we just convert it
	managerInfo := s.k8sManager.ListClusterInfo()
	response := make([]models.ClusterListResponse, len(managerInfo))
	for i, info := range managerInfo {
		response[i] = models.ClusterListResponse{
			ID:          info.ID, // Ensure k8s.ClusterInfoResponse has ID field
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

// GetClusterByID gets detailed information for a single cluster.
func (s *ClusterService) GetClusterByID(id string) (*models.ClusterResponse, error) {
	cluster, err := s.k8sManager.GetClusterDetailFromDB(id)
	if err != nil {
		// If not in database, it might be a file-type cluster, we assemble a simple version from cache
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
		return nil, fmt.Errorf("cluster ID '%s' not found: %w", id, err)
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

// CreateCluster handles the logic for creating a new cluster.
func (s *ClusterService) CreateCluster(req models.CreateClusterRequest) error {
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(req.KubeconfigData)
	if err != nil {
		return fmt.Errorf("kubeconfig data is not valid Base64 encoding: %w", err)
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

// UpdateCluster updates cluster information.
func (s *ClusterService) UpdateCluster(id string, req models.UpdateClusterRequest) error {
	return s.k8sManager.UpdateDBCluster(id, req)
}

// DeleteClusterByID handles the logic for deleting a cluster.
func (s *ClusterService) DeleteClusterByID(id string) error {
	return s.k8sManager.RemoveDBClusterByID(id)
}

// SetActiveCluster handles the logic for switching the active cluster.
func (s *ClusterService) SetActiveCluster(id string) error {
	return s.k8sManager.SetActiveClusterByID(id)
}

// GetActiveClusterID gets the current active cluster ID
func (s *ClusterService) GetActiveClusterID() string {
	return s.k8sManager.GetActiveClusterID()
}
