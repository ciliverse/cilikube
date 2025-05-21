// internal/service/node_service.go
package service

import (
	"context" // Kubernetes API 调用通常需要 context
	"fmt"
	"log"

	"github.com/ciliverse/cilikube/pkg/k8s" // 引入 ClientManager 包

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes" // 不再直接需要 Clientset 作为参数
)

// NodeService 结构体封装了与节点相关的业务逻辑
type NodeService struct {
	// 此服务不再直接持有 k8s 客户端。
	// 它将在每个需要与 k8s API 交互的方法中，
	// 通过 k8s.GetClientManager() 获取 ClientManager，
	// 然后通过 ClientManager.GetClient(clusterName) 获取特定集群的客户端。
}

// NewNodeService 创建一个新的 NodeService 实例。
// 构造函数不再接收 Kubernetes 客户端。
func NewNodeService() *NodeService {
	return &NodeService{}
}

// ListNodes 从指定的集群检索节点列表。
// clusterName 参数用于指定目标集群。
func (s *NodeService) ListNodes(clusterName string) (*corev1.NodeList, error) {
	log.Printf("服务层: 尝试列出集群 '%s' 中的节点。\n", clusterName)
	// 1. 获取 ClientManager 实例
	manager, err := k8s.GetClientManager()
	if err != nil {
		log.Printf("错误: 获取 ClientManager 失败: %v\n", err)
		return nil, fmt.Errorf("获取客户端管理器失败: %w", err)
	}

	// 2. 根据集群名称获取对应的 K8s 客户端
	// 如果 clusterName 为空字符串，GetClient 内部会尝试获取默认集群客户端。
	client, err := manager.GetClient(clusterName)
	if err != nil {
		log.Printf("错误: 获取集群 '%s' 的客户端失败: %v\n", clusterName, err)
		return nil, fmt.Errorf("获取集群 '%s' 的客户端失败: %w", clusterName, err)
	}

	// 3. 使用获取到的客户端执行操作
	nodes, err := client.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("错误: 调用 Kubernetes API 列出集群 '%s' 中的节点失败: %v\n", clusterName, err)
		return nil, fmt.Errorf("列出集群 '%s' 中的节点失败: %w", clusterName, err)
	}

	log.Printf("服务层: 成功列出集群 '%s' 中的 %d 个节点。\n", clusterName, len(nodes.Items))
	return nodes, nil
}

// GetNodeDetails 从指定的集群检索特定节点的详细信息。
func (s *NodeService) GetNodeDetails(clusterName string, nodeName string) (*corev1.Node, error) {
	log.Printf("服务层: 尝试获取集群 '%s' 中节点 '%s' 的详细信息。\n", clusterName, nodeName)
	if nodeName == "" {
		return nil, fmt.Errorf("节点名称不能为空")
	}

	manager, err := k8s.GetClientManager()
	if err != nil {
		log.Printf("错误: 获取 ClientManager 失败: %v\n", err)
		return nil, fmt.Errorf("获取客户端管理器失败: %w", err)
	}

	client, err := manager.GetClient(clusterName)
	if err != nil {
		log.Printf("错误: 获取集群 '%s' 的客户端失败: %v\n", clusterName, err)
		return nil, fmt.Errorf("获取集群 '%s' 的客户端失败: %w", clusterName, err)
	}

	node, err := client.Clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("错误: 调用 Kubernetes API 获取集群 '%s' 中节点 '%s' 的详细信息失败: %v\n", clusterName, nodeName, err)
		return nil, fmt.Errorf("获取集群 '%s' 中节点 '%s' 的详细信息失败: %w", clusterName, nodeName, err)
	}

	log.Printf("服务层: 成功获取集群 '%s' 中节点 '%s' 的详细信息。\n", clusterName, nodeName)
	return node, nil
}
