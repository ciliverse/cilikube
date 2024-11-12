package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/apis/models/k8s"
	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceService provides methods for managing Kubernetes Namespaces
type NamespaceService struct {
	client *client.K8sClient
}

// NewNamespaceService creates a new NamespaceService
func NewNamespaceService(client *client.K8sClient) *NamespaceService {
	return &NamespaceService{
		client: client,
	}
}

// ListNamespaces retrieves a list of Namespaces
func (s *NamespaceService) ListNamespaces() ([]*k8s.Namespace, error) {
	nsList, err := s.client.K8sClient().CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaces []*k8s.Namespace
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, k8s.ToNamespaceModel(&ns))
	}
	return namespaces, nil
}

// GetNamespace retrieves a single Namespace by name
func (s *NamespaceService) GetNamespace(name string) (*k8s.Namespace, error) {
	ns, err := s.client.K8sClient().CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return k8s.ToNamespaceModel(ns), nil
}

// CreateNamespace creates a new Namespace
func (s *NamespaceService) CreateNamespace(name string, labels map[string]string) (*k8s.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}

	createdNs, err := s.client.K8sClient().CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return k8s.ToNamespaceModel(createdNs), nil
}

// DeleteNamespace deletes an existing Namespace by name
// DeleteNamespace deletes an existing Namespace by name
func (s *NamespaceService) DeleteNamespace(name string) error {
	return s.client.K8sClient().CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
