package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes" // 仍然需要导入以获取 kubernetes.Interface 类型
)

// NodeService 结构体不再持有 client 字段。
type NodeService struct {
	// 不需要 client kubernetes.Interface 字段了
}

// NewNodeService 构造函数不再接收 kubernetes.Interface 参数。
func NewNodeService() *NodeService {
	return &NodeService{}
}

// Get 获取单个Node。
// 第一个参数 'clientset kubernetes.Interface' 是动态传入的特定集群的客户端。
func (s *NodeService) Get(clientset kubernetes.Interface, name string) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// Create 创建Node。
// 第一个参数 'clientset kubernetes.Interface' 是动态传入的特定集群的客户端。
func (s *NodeService) Create(clientset kubernetes.Interface, node *corev1.Node) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Create(
		context.TODO(),
		node,
		metav1.CreateOptions{},
	)
}

// Update 更新Node。
// 第一个参数 'clientset kubernetes.Interface' 是动态传入的特定集群的客户端。
func (s *NodeService) Update(clientset kubernetes.Interface, node *corev1.Node) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Update(
		context.TODO(),
		node,
		metav1.UpdateOptions{},
	)
}

// Delete 删除Node。
// 第一个参数 'clientset kubernetes.Interface' 是动态传入的特定集群的客户端。
func (s *NodeService) Delete(clientset kubernetes.Interface, name string) error {
	return clientset.CoreV1().Nodes().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// List 列表查询（支持分页和标签过滤）。
// 第一个参数 'clientset kubernetes.Interface' 是动态传入的特定集群的客户端。
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

// Watch Watch机制实现。
// 第一个参数 'clientset kubernetes.Interface' 是动态传入的特定集群的客户端。
func (s *NodeService) Watch(clientset kubernetes.Interface, selector string, resourceVersion string, timeoutSeconds int64) (watch.Interface, error) {
	var timeoutSecondsPtr *int64
	if timeoutSeconds > 0 {
		timeoutSecondsPtr = &timeoutSeconds
	} else {
		// 默认超时时间，例如1800秒 (30分钟)
		// 或者根据你的需求设定，如果 timeoutSeconds 为0或负数，则可能表示不设置或使用服务器默认
		// 在你的原始代码中是 int64ptr(1800)，这里我们保持一致，如果 timeoutSeconds <=0 则使用默认
		defaultTimeout := int64(1800) // 和你原代码一致
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

/*
// 如果你没有全局的 int64ptr 辅助函数，可以在这里或公共包中定义一个：
func int64ptr(i int64) *int64 {
    return &i
}
// 你的代码中已经使用了 int64ptr(1800)，所以你可能已经有了这个辅助函数。
// 注意：对于 ListOptions.Limit，如果 limit 为0，通常表示不限制（或由服务器决定）。
// 对于 WatchOptions.TimeoutSeconds，如果为 nil，表示没有客户端超时（依赖服务器端超时）。
// 你之前的代码是写死的 `TimeoutSeconds: int64ptr(1800)`，
// 在新的方法签名中，我将其改为了一个参数 `timeoutSeconds int64`，以便更灵活地控制。
// 并在方法内部处理转换成 `*int64`。
// 我还为 List 方法添加了 continueToken 参数以支持分页。
// 为 Watch 方法添加了 resourceVersion 和 timeoutSeconds 参数。
*/
