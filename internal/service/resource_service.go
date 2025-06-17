package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// ResourceService 通用接口
// T 为 K8s 资源类型（如 *corev1.Node）
type ResourceService[T any] interface {
	Get(client kubernetes.Interface, name string) (T, error)
	List(client kubernetes.Interface, selector string, limit int64, continueToken string) ([]T, error)
	Create(client kubernetes.Interface, obj T) (T, error)
	Update(client kubernetes.Interface, obj T) (T, error)
	Delete(client kubernetes.Interface, name string) error
	Watch(client kubernetes.Interface, selector, resourceVersion string, timeoutSeconds int64) (watch.Interface, error)
}

// NodeResourceService 实现 ResourceService[*corev1.Node]
// 其它资源可用类似方式实现
type NodeResourceService struct{}

func (s *NodeResourceService) Get(client kubernetes.Interface, name string) (*corev1.Node, error) {
	return client.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
}
func (s *NodeResourceService) List(client kubernetes.Interface, selector string, limit int64, continueToken string) ([]*corev1.Node, error) {
	opts := metav1.ListOptions{LabelSelector: selector, Limit: limit, Continue: continueToken}
	list, err := client.CoreV1().Nodes().List(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	nodes := make([]*corev1.Node, len(list.Items))
	for i := range list.Items {
		nodes[i] = &list.Items[i]
	}
	return nodes, nil
}
func (s *NodeResourceService) Create(client kubernetes.Interface, obj *corev1.Node) (*corev1.Node, error) {
	return client.CoreV1().Nodes().Create(context.TODO(), obj, metav1.CreateOptions{})
}
func (s *NodeResourceService) Update(client kubernetes.Interface, obj *corev1.Node) (*corev1.Node, error) {
	return client.CoreV1().Nodes().Update(context.TODO(), obj, metav1.UpdateOptions{})
}
func (s *NodeResourceService) Delete(client kubernetes.Interface, name string) error {
	return client.CoreV1().Nodes().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
func (s *NodeResourceService) Watch(client kubernetes.Interface, selector, resourceVersion string, timeoutSeconds int64) (watch.Interface, error) {
	opts := metav1.ListOptions{LabelSelector: selector, ResourceVersion: resourceVersion, Watch: true}
	if timeoutSeconds > 0 {
		opts.TimeoutSeconds = &timeoutSeconds
	}
	return client.CoreV1().Nodes().Watch(context.TODO(), opts)
}
