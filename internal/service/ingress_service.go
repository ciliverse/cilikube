package service

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type IngressService struct {
	client kubernetes.Interface
}

func NewIngressService(client kubernetes.Interface) *IngressService {
	return &IngressService{client: client}
}

// 获取单个Ingress
func (s *IngressService) Get(namespace, name string) (*networkingv1.Ingress, error) {
	return s.client.NetworkingV1().Ingresses(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Ingress
func (s *IngressService) Create(namespace string, ingress *networkingv1.Ingress) (*networkingv1.Ingress, error) {

	if ingress.Namespace != "" && ingress.Namespace != namespace {
		return nil, NewValidationError("ingress namespace conflicts with path parameter")
	}

	return s.client.NetworkingV1().Ingresses(namespace).Create(
		context.TODO(),
		ingress,
		metav1.CreateOptions{},
	)
}

// 更新Ingress
func (s *IngressService) Update(namespace string, ingress *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	return s.client.NetworkingV1().Ingresses(namespace).Update(
		context.TODO(),
		ingress,
		metav1.UpdateOptions{},
	)
}

// 删除Ingress
func (s *IngressService) Delete(namespace, name string) error {
	return s.client.NetworkingV1().Ingresses(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *IngressService) List(namespace, selector string, limit int64) (*networkingv1.IngressList, error) {
	return s.client.NetworkingV1().Ingresses(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *IngressService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.NetworkingV1().Ingresses(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
