package services

import (
	"context"
	"log"

	"github.com/cilliandevops/kopsadmin/server-go/internal/apis/models/k8s"
	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterService struct {
	client *client.K8sClient
}

// NewClusterService creates a new ClusterService instance
func NewClusterService(k8sClient *client.K8sClient) *ClusterService {
	return &ClusterService{
		client: k8sClient,
	}
}

// getNodeInfo extracts the node details
func (s *ClusterService) getNodeInfo(node *v1.Node) k8s.NodeInfo {
	return k8s.NodeInfo{
		Name:        node.Name,
		Labels:      node.Labels,
		Annotations: node.Annotations,
		CreatedAt:   node.CreationTimestamp,
	}
}

// GetClusterInfo retrieves general information about the Kubernetes cluster
func (s *ClusterService) GetClusterInfo(ctx context.Context) (*k8s.ClusterInfo, error) {
	// 使用 s.client.K8sClient() 获取底层的 Clientset
	version, err := s.client.K8sClient().Discovery().ServerVersion()
	if err != nil {
		log.Printf("Failed to get server version: %v", err)
		return nil, err
	}

	return &k8s.ClusterInfo{
		Version: version.String(),
	}, nil
}

// ListNodes lists all nodes in the cluster
func (s *ClusterService) ListNodes(ctx context.Context) ([]k8s.NodeInfo, error) {
	// 使用 s.client.K8sClient() 获取底层的 Clientset
	nodes, err := s.client.K8sClient().CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list nodes: %v", err)
		return nil, err
	}

	var nodeInfos []k8s.NodeInfo
	for _, node := range nodes.Items {
		nodeInfos = append(nodeInfos, s.getNodeInfo(&node))
	}

	return nodeInfos, nil
}

// GetNode retrieves detailed information of a specific node
func (s *ClusterService) GetNode(ctx context.Context, nodeName string) (*k8s.NodeInfo, error) {
	// 使用 s.client.K8sClient() 获取底层的 Clientset
	node, err := s.client.K8sClient().CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get node %s: %v", nodeName, err)
		return nil, err
	}

	nodeInfo := s.getNodeInfo(node)
	return &nodeInfo, nil
}

// DeleteNode deletes a node from the cluster
func (s *ClusterService) DeleteNode(ctx context.Context, nodeName string) error {
	// 使用 s.client.K8sClient() 获取底层的 Clientset
	err := s.client.K8sClient().CoreV1().Nodes().Delete(ctx, nodeName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Failed to delete node %s: %v", nodeName, err)
	}
	return err
}
