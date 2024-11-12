package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapService 提供对 Kubernetes ConfigMap 资源的操作
type ConfigMapService struct {
	client *client.K8sClient
}

// NewConfigMapService 创建一个新的 ConfigMapService 实例
func NewConfigMapService(k8sClient *client.K8sClient) *ConfigMapService {
	return &ConfigMapService{
		client: k8sClient,
	}
}

// CreateConfigMap 创建一个新的 ConfigMap
func (s *ConfigMapService) CreateConfigMap(namespace string, configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	return s.client.K8sClient().CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
}

func (s *ConfigMapService) ListConfigMaps(namespace string) (*v1.ConfigMapList, error) {
	return s.client.K8sClient().CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
}

// UpdateConfigMap 更新指定的 ConfigMap
func (s *ConfigMapService) UpdateConfigMap(namespace string, configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	return s.client.K8sClient().CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
}

// DeleteConfigMap 删除指定命名空间中的指定 ConfigMap
func (s *ConfigMapService) DeleteConfigMap(namespace, name string) error {
	return s.client.K8sClient().CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
