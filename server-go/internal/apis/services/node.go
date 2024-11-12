package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/apis/models/k8s"
	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeService struct {
	client *client.K8sClient
}

// NewNodeService creates a new NodeService
func NewNodeService(client *client.K8sClient) *NodeService {
	return &NodeService{
		client: client,
	}
}

// ListNodes returns a list of all Nodes with detailed information
func (s *NodeService) ListNodes() ([]k8s.NodeModel, error) {
	nodes, err := s.client.K8sClient().CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeModels []k8s.NodeModel
	for _, node := range nodes.Items {
		nodeModels = append(nodeModels, *k8s.NewNodeModel(&node))
	}
	return nodeModels, nil
}

// GetNode returns the details of a specific Node
func (s *NodeService) GetNode(name string) (*k8s.NodeModel, error) {
	node, err := s.client.K8sClient().CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return k8s.NewNodeModel(node), nil
}

// CreateNode creates a new Node
func (s *NodeService) CreateNode(node *v1.Node) (*k8s.NodeModel, error) {
	createdNode, err := s.client.K8sClient().CoreV1().Nodes().Create(context.Background(), node, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return k8s.NewNodeModel(createdNode), nil
}

// DeleteNode deletes a Node by name
func (s *NodeService) DeleteNode(name string) error {
	return s.client.K8sClient().CoreV1().Nodes().Delete(context.Background(), name, metav1.DeleteOptions{})
}
