package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type NodeService struct {
	client kubernetes.Interface
}

func NewNodeService(client kubernetes.Interface) *NodeService {
	return &NodeService{client: client}
}

// 获取单个Node
func (s *NodeService) Get(name string) (*corev1.Node, error) {
	return s.client.CoreV1().Nodes().Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Node
func (s *NodeService) Create(node *corev1.Node) (*corev1.Node, error) {
	return s.client.CoreV1().Nodes().Create(
		context.TODO(),
		node,
		metav1.CreateOptions{},
	)
}

// 更新Node
func (s *NodeService) Update(node *corev1.Node) (*corev1.Node, error) {
	return s.client.CoreV1().Nodes().Update(
		context.TODO(),
		node,
		metav1.UpdateOptions{},
	)
}

// 删除Node
func (s *NodeService) Delete(name string) error {
	return s.client.CoreV1().Nodes().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *NodeService) List(selector string, limit int64) (*corev1.NodeList, error) {
	return s.client.CoreV1().Nodes().List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *NodeService) Watch(selector string) (watch.Interface, error) {
	return s.client.CoreV1().Nodes().Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
