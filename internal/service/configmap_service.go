package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapService struct {
	client kubernetes.Interface
}

func NewConfigMapService(client kubernetes.Interface) *ConfigMapService {
	return &ConfigMapService{client: client}
}

// 获取单个ConfigMap
func (s *ConfigMapService) Get(namespace, name string) (*corev1.ConfigMap, error) {
	return s.client.CoreV1().ConfigMaps(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建ConfigMap
func (s *ConfigMapService) Create(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {

	if configMap.Namespace != "" && configMap.Namespace != namespace {
		return nil, NewValidationError("configMap namespace conflicts with path parameter")
	}

	return s.client.CoreV1().ConfigMaps(namespace).Create(
		context.TODO(),
		configMap,
		metav1.CreateOptions{},
	)
}

// 更新ConfigMap
func (s *ConfigMapService) Update(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return s.client.CoreV1().ConfigMaps(namespace).Update(
		context.TODO(),
		configMap,
		metav1.UpdateOptions{},
	)
}

// 删除ConfigMap
func (s *ConfigMapService) Delete(namespace, name string) error {
	return s.client.CoreV1().ConfigMaps(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *ConfigMapService) List(namespace, selector string, limit int64) (*corev1.ConfigMapList, error) {
	return s.client.CoreV1().ConfigMaps(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *ConfigMapService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.CoreV1().ConfigMaps(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
