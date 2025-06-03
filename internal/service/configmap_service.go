package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapService 结构体不再持有 client 字段
type ConfigMapService struct {
	// 不需要 client kubernetes.Interface 字段了
}

func NewConfigMapService() *ConfigMapService {
	return &ConfigMapService{}
}

// Get retrieves a single ConfigMap by namespace and name.
func (s *ConfigMapService) Get(clientSet kubernetes.Interface, namespace, name string) (*corev1.ConfigMap, error) {
	return clientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// List retrieves ConfigMaps within a specific namespace.
func (s *ConfigMapService) List(clientSet kubernetes.Interface, namespace, labelSelector string, limit int64) (*corev1.ConfigMapList, error) {
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}
	if limit > 0 {
		listOptions.Limit = limit
	}
	return clientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), listOptions)
}

// Create creates a new ConfigMap in the specified namespace.
func (s *ConfigMapService) Create(clientSet kubernetes.Interface, namespace string, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if cm.Namespace != "" && cm.Namespace != namespace {
		return nil, NewValidationError("ConfigMap namespace conflicts with path parameter")
	}
	if cm.Namespace == "" {
		cm.Namespace = namespace
	}
	if cm.Name == "" {
		return nil, NewValidationError("ConfigMap name cannot be empty")
	}

	return clientSet.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
}

// Update updates an existing ConfigMap.
func (s *ConfigMapService) Update(clientSet kubernetes.Interface, namespace string, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if cm.Namespace != "" && cm.Namespace != namespace {
		return nil, NewValidationError("ConfigMap namespace conflicts with path parameter")
	}
	if cm.Namespace == "" {
		cm.Namespace = namespace
	}
	if cm.Name == "" {
		return nil, NewValidationError("ConfigMap name cannot be empty for update")
	}

	// Fetch existing for ResourceVersion recommended
	// existingCM, err := s.Get(namespace, cm.Name)
	// if err != nil { return nil, err }
	// cm.ResourceVersion = existingCM.ResourceVersion

	return clientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
}

// Delete deletes a ConfigMap by namespace and name.
func (s *ConfigMapService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// --- Re-use or define ValidationError ---
// type ValidationError struct { Message string }
// func (e *ValidationError) Error() string { return e.Message }
// func NewValidationError(msg string) error { return &ValidationError{Message: msg} }
