package service

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type NetworkPolicyService struct {
	client kubernetes.Interface
}

func NewNetworkPolicyService(client kubernetes.Interface) *NetworkPolicyService {
	return &NetworkPolicyService{client: client}
}

// 获取单个NetworkPolicy
func (s *NetworkPolicyService) Get(namespace, name string) (*networkingv1.NetworkPolicy, error) {
	return s.client.NetworkingV1().NetworkPolicies(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建NetworkPolicy
func (s *NetworkPolicyService) Create(namespace string, networkPolicy *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {

	if networkPolicy.Namespace != "" && networkPolicy.Namespace != namespace {
		return nil, NewValidationError("networkPolicy namespace conflicts with path parameter")
	}

	return s.client.NetworkingV1().NetworkPolicies(namespace).Create(
		context.TODO(),
		networkPolicy,
		metav1.CreateOptions{},
	)
}

// 更新NetworkPolicy
func (s *NetworkPolicyService) Update(namespace string, networkPolicy *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {
	return s.client.NetworkingV1().NetworkPolicies(namespace).Update(
		context.TODO(),
		networkPolicy,
		metav1.UpdateOptions{},
	)
}

// 删除NetworkPolicy
func (s *NetworkPolicyService) Delete(namespace, name string) error {
	return s.client.NetworkingV1().NetworkPolicies(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *NetworkPolicyService) List(namespace, selector string, limit int64) (*networkingv1.NetworkPolicyList, error) {
	return s.client.NetworkingV1().NetworkPolicies(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *NetworkPolicyService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.NetworkingV1().NetworkPolicies(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
