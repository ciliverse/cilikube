package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type ServiceService struct {
	client kubernetes.Interface
}

func NewServiceService(client kubernetes.Interface) *ServiceService {
	return &ServiceService{client: client}
}

// 列表查询（支持分页和标签过滤）
func (s *ServiceService) List(namespace string) (*corev1.ServiceList, error) {
	if namespace == "" {
		namespace = corev1.NamespaceAll
	}

	return s.client.CoreV1().Services(namespace).List(
		context.TODO(),
		metav1.ListOptions{},
	)
}

// 获取单个Service
func (s *ServiceService) Get(namespace, name string) (*corev1.Service, error) {
	return s.client.CoreV1().Services(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Service
func (s *ServiceService) Create(namespace string, service *corev1.Service) (*corev1.Service, error) {

	if service.Namespace != "" && service.Namespace != namespace {
		return nil, NewValidationError("service namespace conflicts with path parameter")
	}

	return s.client.CoreV1().Services(namespace).Create(
		context.TODO(),
		service,
		metav1.CreateOptions{},
	)
}

// 更新Service
func (s *ServiceService) Update(namespace string, service *corev1.Service) (*corev1.Service, error) {
	return s.client.CoreV1().Services(namespace).Update(
		context.TODO(),
		service,
		metav1.UpdateOptions{},
	)
}

// 删除Service
func (s *ServiceService) Delete(namespace, name string) error {
	return s.client.CoreV1().Services(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// Watch机制实现
func (s *ServiceService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.CoreV1().Services(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
