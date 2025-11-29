package service

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// NodeMetrics represents real-time resource usage of a single node
type NodeMetrics struct {
	CPUUsage    string `json:"cpuUsage"`
	MemoryUsage string `json:"memoryUsage"`
}

// NodeMetricsService provides service for getting node metrics
type NodeMetricsService struct{}

// NewNodeMetricsService creates a new NodeMetricsService instance
func NewNodeMetricsService() *NodeMetricsService {
	return &NodeMetricsService{}
}

// GetNodeMetrics gets real-time metrics of a single node through metrics-server.
// It requires passing the target cluster's rest.Config to create a dedicated metrics client.
func (s *NodeMetricsService) GetNodeMetrics(config *rest.Config, nodeName string) (*NodeMetrics, error) {
	// Create Metrics API client based on the passed cluster configuration
	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	// Call Metrics API to get metrics for the specified node
	nodeMetrics, err := metricsClientset.MetricsV1beta1().NodeMetricses().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics for node '%s' from metrics API: %w", nodeName, err)
	}

	// Format the returned data
	metrics := &NodeMetrics{
		CPUUsage:    nodeMetrics.Usage.Cpu().String(),
		MemoryUsage: nodeMetrics.Usage.Memory().String(),
	}

	return metrics, nil
}
