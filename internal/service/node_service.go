package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes" // 仍然需要导入以获取 kubernetes.Interface 类型
)

// NodeService 结构体不再持有 client 字段
type NodeService struct {
	// 不需要 client kubernetes.Interface 字段了
}

// NewNodeService 构造函数不再接收 kubernetes.Interface 参数
func NewNodeService() *NodeService {
	return &NodeService{}
}

// Get 获取单个Node

func (s *NodeService) Get(clientset kubernetes.Interface, name string) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// Create 创建Node

func (s *NodeService) Create(clientset kubernetes.Interface, node *corev1.Node) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Create(
		context.TODO(),
		node,
		metav1.CreateOptions{},
	)
}

// Update 更新Node

func (s *NodeService) Update(clientset kubernetes.Interface, node *corev1.Node) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Update(
		context.TODO(),
		node,
		metav1.UpdateOptions{},
	)
}

// Delete 删除Node

func (s *NodeService) Delete(clientset kubernetes.Interface, name string) error {
	return clientset.CoreV1().Nodes().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// List 列表查询（支持分页和标签过滤）

func (s *NodeService) List(clientset kubernetes.Interface, selector string, limit int64, continueToken string) (*corev1.NodeList, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: selector,
		Limit:         limit,
		Continue:      continueToken,
	}
	return clientset.CoreV1().Nodes().List(
		context.TODO(),
		listOptions,
	)
}

// Watch Watch机制实现

func (s *NodeService) Watch(clientset kubernetes.Interface, selector string, resourceVersion string, timeoutSeconds int64) (watch.Interface, error) {
	var timeoutSecondsPtr *int64
	if timeoutSeconds > 0 {
		timeoutSecondsPtr = &timeoutSeconds
	} else {

		defaultTimeout := int64(1800)
		if timeoutSeconds > 0 {
			timeoutSecondsPtr = &timeoutSeconds
		} else if timeoutSeconds == 0 { // 如果明确传入0，可能表示使用默认
			timeoutSecondsPtr = &defaultTimeout
		}
		// 如果timeoutSeconds < 0, 则timeoutSecondsPtr为nil, 表示没有客户端超时，依赖服务端
	}

	return clientset.CoreV1().Nodes().Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:   selector,
			ResourceVersion: resourceVersion, // 允许从特定资源版本开始watch
			Watch:           true,
			TimeoutSeconds:  timeoutSecondsPtr,
		},
	)
}
