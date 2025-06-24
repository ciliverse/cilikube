package service

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// NodeMetrics 表示单个节点的实时资源使用情况
type NodeMetrics struct {
	CPUUsage    string `json:"cpuUsage"`
	MemoryUsage string `json:"memoryUsage"`
}

// NodeMetricsService 提供了获取节点指标的服务
type NodeMetricsService struct{}

// NewNodeMetricsService 创建一个新的 NodeMetricsService 实例
func NewNodeMetricsService() *NodeMetricsService {
	return &NodeMetricsService{}
}

// GetNodeMetrics 通过 metrics-server 获取单个节点的实时指标。
// 它需要传入目标集群的 rest.Config 来创建专门的 metrics 客户端。
func (s *NodeMetricsService) GetNodeMetrics(config *rest.Config, nodeName string) (*NodeMetrics, error) {
	// 基于传入的集群配置，创建 Metrics API 的客户端
	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 metrics 客户端失败: %w", err)
	}

	// 调用 Metrics API 获取指定节点的指标
	nodeMetrics, err := metricsClientset.MetricsV1beta1().NodeMetricses().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("从 metrics API 获取节点 '%s' 的指标失败: %w", nodeName, err)
	}

	// 格式化返回的数据
	metrics := &NodeMetrics{
		CPUUsage:    nodeMetrics.Usage.Cpu().String(),
		MemoryUsage: nodeMetrics.Usage.Memory().String(),
	}

	return metrics, nil
}
