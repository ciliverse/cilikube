package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretService 结构体不再持有 client 字段
type SecretService struct {
	// 不需要 client kubernetes.Interface 字段了
}

func NewSecretService() *SecretService {
	return &SecretService{}
}

// Get retrieves a single Secret by namespace and name.
func (s *SecretService) Get(clientSet kubernetes.Interface, namespace, name string) (*corev1.Secret, error) {
	return clientSet.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// List retrieves Secrets within a specific namespace.
func (s *SecretService) List(clientSet kubernetes.Interface, namespace, labelSelector string, limit int64) (*corev1.SecretList, error) {
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}
	if limit > 0 {
		listOptions.Limit = limit
	}
	return clientSet.CoreV1().Secrets(namespace).List(context.TODO(), listOptions)
}

// Create creates a new Secret in the specified namespace.
func (s *SecretService) Create(clientSet kubernetes.Interface, namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	if secret.Namespace != "" && secret.Namespace != namespace {
		return nil, NewValidationError("Secret namespace conflicts")
	}
	if secret.Namespace == "" {
		secret.Namespace = namespace
	}
	if secret.Name == "" {
		return nil, NewValidationError("Secret name cannot be empty")
	}
	// K8s automatically base64 encodes StringData into Data if Data[key] doesn't exist.
	// No need for manual encoding here if receiving corev1.Secret object.
	return clientSet.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
}

// Update updates an existing Secret.
func (s *SecretService) Update(clientSet kubernetes.Interface, namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	if secret.Namespace != "" && secret.Namespace != namespace {
		return nil, NewValidationError("Secret namespace conflicts")
	}
	if secret.Namespace == "" {
		secret.Namespace = namespace
	}
	if secret.Name == "" {
		return nil, NewValidationError("Secret name required for update")
	}
	// Fetch existing for ResourceVersion recommended
	return clientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
}

// Delete deletes a Secret by namespace and name.
func (s *SecretService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// --- Re-use or define ValidationError ---
// type ValidationError struct { Message string } ...
