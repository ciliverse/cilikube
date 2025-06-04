package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// NetworkPolicyService 结构体不再持有 client 字段
type NetworkPolicyService struct {
	// 不需要 client kubernetes.Interface 字段了
}

func NewNetworkPolicyService() *NetworkPolicyService {
	return &NetworkPolicyService{}
}

// 获取单个NetworkPolicy
func (s *NetworkPolicyService) Get(clientSet kubernetes.Interface, namespace, name string) (*networkingv1.NetworkPolicy, error) {
	return clientSet.NetworkingV1().NetworkPolicies(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建NetworkPolicy
func (s *NetworkPolicyService) Create(clientSet kubernetes.Interface, namespace string, networkPolicy *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {

	if networkPolicy.Namespace != "" && networkPolicy.Namespace != namespace {
		return nil, NewValidationError("networkPolicy namespace conflicts with path parameter")
	}

	return clientSet.NetworkingV1().NetworkPolicies(namespace).Create(
		context.TODO(),
		networkPolicy,
		metav1.CreateOptions{},
	)
}

// 更新NetworkPolicy
func (s *NetworkPolicyService) Update(clientSet kubernetes.Interface, namespace string, networkPolicy *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {
	return clientSet.NetworkingV1().NetworkPolicies(namespace).Update(
		context.TODO(),
		networkPolicy,
		metav1.UpdateOptions{},
	)
}

// 删除NetworkPolicy
func (s *NetworkPolicyService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.NetworkingV1().NetworkPolicies(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *NetworkPolicyService) List(clientSet kubernetes.Interface, namespace, selector string, limit int64) (*networkingv1.NetworkPolicyList, error) {
	return clientSet.NetworkingV1().NetworkPolicies(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *NetworkPolicyService) Watch(clientSet kubernetes.Interface, namespace, selector string) (watch.Interface, error) {
	return clientSet.NetworkingV1().NetworkPolicies(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
