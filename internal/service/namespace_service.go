package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService struct {
	// 不再持有 client 字段
}

// 构造函数不再接收 kubernetes.Interface 参数
func NewNamespaceService() *NamespaceService {
	return &NamespaceService{}
}

// 获取单个Namespace
func (s *NamespaceService) Get(client kubernetes.Interface, name string) (*corev1.Namespace, error) {
	return client.CoreV1().Namespaces().Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Namespace
func (s *NamespaceService) Create(client kubernetes.Interface, namespace *corev1.Namespace) (*corev1.Namespace, error) {
	return client.CoreV1().Namespaces().Create(
		context.TODO(),
		namespace,
		metav1.CreateOptions{},
	)
}

// 更新Namespace
func (s *NamespaceService) Update(client kubernetes.Interface, namespace *corev1.Namespace) (*corev1.Namespace, error) {
	return client.CoreV1().Namespaces().Update(
		context.TODO(),
		namespace,
		metav1.UpdateOptions{},
	)
}

// 删除Namespace
func (s *NamespaceService) Delete(client kubernetes.Interface, name string) error {
	return client.CoreV1().Namespaces().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *NamespaceService) List(client kubernetes.Interface, selector string, limit int64) (*corev1.NamespaceList, error) {
	return client.CoreV1().Namespaces().List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *NamespaceService) Watch(client kubernetes.Interface, selector string) (watch.Interface, error) {
	return client.CoreV1().Namespaces().Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
