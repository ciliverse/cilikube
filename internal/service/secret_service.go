package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type SecretService struct {
	client kubernetes.Interface
}

func NewSecretService(client kubernetes.Interface) *SecretService {
	return &SecretService{client: client}
}

// 获取单个Secret
func (s *SecretService) Get(namespace, name string) (*corev1.Secret, error) {
	return s.client.CoreV1().Secrets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Secret
func (s *SecretService) Create(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {

	if secret.Namespace != "" && secret.Namespace != namespace {
		return nil, NewValidationError("secret namespace conflicts with path parameter")
	}

	return s.client.CoreV1().Secrets(namespace).Create(
		context.TODO(),
		secret,
		metav1.CreateOptions{},
	)
}

// 更新Secret
func (s *SecretService) Update(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	return s.client.CoreV1().Secrets(namespace).Update(
		context.TODO(),
		secret,
		metav1.UpdateOptions{},
	)
}

// 删除Secret
func (s *SecretService) Delete(namespace, name string) error {
	return s.client.CoreV1().Secrets(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *SecretService) List(namespace, selector string, limit int64) (*corev1.SecretList, error) {
	return s.client.CoreV1().Secrets(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *SecretService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.CoreV1().Secrets(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
