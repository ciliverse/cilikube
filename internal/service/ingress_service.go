package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// IngressService 结构体不再持有 client 字段
type IngressService struct {
	// 不需要 client kubernetes.Interface 字段了
}

func NewIngressService() *IngressService {
	return &IngressService{}
}

// 获取单个Ingress
func (s *IngressService) Get(clientSet kubernetes.Interface, namespace, name string) (*networkingv1.Ingress, error) {
	return clientSet.NetworkingV1().Ingresses(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Ingress
func (s *IngressService) Create(clientSet kubernetes.Interface, namespace string, ingress *networkingv1.Ingress) (*networkingv1.Ingress, error) {

	if ingress.Namespace != "" && ingress.Namespace != namespace {
		return nil, NewValidationError("ingress namespace conflicts with path parameter")
	}

	return clientSet.NetworkingV1().Ingresses(namespace).Create(
		context.TODO(),
		ingress,
		metav1.CreateOptions{},
	)
}

// 更新Ingress
func (s *IngressService) Update(clientSet kubernetes.Interface, namespace string, ingress *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	return clientSet.NetworkingV1().Ingresses(namespace).Update(
		context.TODO(),
		ingress,
		metav1.UpdateOptions{},
	)
}

// 删除Ingress
func (s *IngressService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.NetworkingV1().Ingresses(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *IngressService) List(clientSet kubernetes.Interface, namespace, selector string, limit int64) (*networkingv1.IngressList, error) {
	return clientSet.NetworkingV1().Ingresses(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *IngressService) Watch(clientSet kubernetes.Interface, namespace, selector string) (watch.Interface, error) {
	return clientSet.NetworkingV1().Ingresses(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
