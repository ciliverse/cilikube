package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretService 提供对 Kubernetes Secret 资源的操作
type SecretService struct {
	client *client.K8sClient
}

// NewSecretService 创建一个新的 SecretService 实例
func NewSecretService(k8sClient *client.K8sClient) *SecretService {
	return &SecretService{
		client: k8sClient,
	}
}

// CreateSecret 创建一个新的 Secret
func (s *SecretService) CreateSecret(namespace string, secret *v1.Secret) (*v1.Secret, error) {
	return s.client.K8sClient().CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
}

// GetSecret 获取指定命名空间和名称的 Secret
func (s *SecretService) GetSecret(namespace, name string) (*v1.Secret, error) {
	return s.client.K8sClient().CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// UpdateSecret 更新指定的 Secret
func (s *SecretService) UpdateSecret(namespace string, secret *v1.Secret) (*v1.Secret, error) {
	return s.client.K8sClient().CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
}

// DeleteSecret 删除指定命名空间和名称的 Secret
func (s *SecretService) DeleteSecret(namespace, name string) error {
	return s.client.K8sClient().CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
