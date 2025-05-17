package models

// ClusterInfo represents metadata about a configured Kubernetes cluster.
type ClusterInfo struct {
	Name           string `json:"name"`
	KubeconfigPath string `json:"kubeconfigPath"` // Path where the kubeconfig file is stored on the server
	Description    string `json:"description,omitempty"`
	IsActive       bool   `json:"isActive,omitempty"` // Transient field, set at runtime
}

// AddClusterRequest is the request body for adding a new cluster.
type AddClusterRequest struct {
	Name              string `json:"name" binding:"required"`
	KubeconfigContent string `json:"kubeconfigContent" binding:"required"`
	Description       string `json:"description"`
}
