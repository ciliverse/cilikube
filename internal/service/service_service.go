package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// ServiceService 结构体不再持有 client 字段
type ServiceService struct {
	// 不需要 client kubernetes.Interface 字段了
}

func NewServiceService() *ServiceService {
	return &ServiceService{}
}

// 列表查询（支持分页和标签过滤）
func (s *ServiceService) List(clientSet kubernetes.Interface, namespace string) (*corev1.ServiceList, error) {
	if namespace == "" {
		namespace = corev1.NamespaceAll
	}

	return clientSet.CoreV1().Services(namespace).List(
		context.TODO(),
		metav1.ListOptions{},
	)
}

// 获取单个Service
func (s *ServiceService) Get(clientSet kubernetes.Interface, namespace, name string) (*corev1.Service, error) {
	return clientSet.CoreV1().Services(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Service
func (s *ServiceService) Create(clientSet kubernetes.Interface, namespace string, service *corev1.Service) (*corev1.Service, error) {

	if service.Namespace != "" && service.Namespace != namespace {
		return nil, NewValidationError("service namespace conflicts with path parameter")
	}

	return clientSet.CoreV1().Services(namespace).Create(
		context.TODO(),
		service,
		metav1.CreateOptions{},
	)
}

// 更新Service
func (s *ServiceService) Update(clientSet kubernetes.Interface, namespace string, service *corev1.Service) (*corev1.Service, error) {
	return clientSet.CoreV1().Services(namespace).Update(
		context.TODO(),
		service,
		metav1.UpdateOptions{},
	)
}

// 删除Service
func (s *ServiceService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.CoreV1().Services(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// Watch机制实现
func (s *ServiceService) Watch(clientSet kubernetes.Interface, namespace, selector string) (watch.Interface, error) {
	return clientSet.CoreV1().Services(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
