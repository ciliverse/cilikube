package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService struct {
	client kubernetes.Interface
}

func NewNamespaceService(client kubernetes.Interface) *NamespaceService {
	return &NamespaceService{client: client}
}

// 获取单个Namespace
func (s *NamespaceService) Get(name string) (*corev1.Namespace, error) {
	return s.client.CoreV1().Namespaces().Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Namespace
func (s *NamespaceService) Create(namespace *corev1.Namespace) (*corev1.Namespace, error) {
	return s.client.CoreV1().Namespaces().Create(
		context.TODO(),
		namespace,
		metav1.CreateOptions{},
	)
}

// 更新Namespace
func (s *NamespaceService) Update(namespace *corev1.Namespace) (*corev1.Namespace, error) {
	return s.client.CoreV1().Namespaces().Update(
		context.TODO(),
		namespace,
		metav1.UpdateOptions{},
	)
}

// 删除Namespace
func (s *NamespaceService) Delete(name string) error {
	return s.client.CoreV1().Namespaces().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *NamespaceService) List(selector string, limit int64) (*corev1.NamespaceList, error) {
	return s.client.CoreV1().Namespaces().List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *NamespaceService) Watch(selector string) (watch.Interface, error) {
	return s.client.CoreV1().Namespaces().Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
