package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressService provides methods for managing Kubernetes Ingresses
type IngressService struct {
	client *client.K8sClient
}

// NewIngressService creates a new instance of IngressService
func NewIngressService(client *client.K8sClient) *IngressService {
	return &IngressService{client: client}
}

// ListIngresses lists all ingresses in a given namespace
func (s *IngressService) ListIngresses(ctx context.Context, namespace string) ([]networkingv1.Ingress, error) {
	// Use the K8sClient's clientset to list ingresses
	ingresses, err := s.client.K8sClient().NetworkingV1().Ingresses(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ingresses.Items, nil
}

// GetIngress retrieves a specific ingress by name in a given namespace
func (s *IngressService) GetIngress(ctx context.Context, namespace, name string) (*networkingv1.Ingress, error) {
	// Use the K8sClient's clientset to get a specific ingress
	return s.client.K8sClient().NetworkingV1().Ingresses(namespace).Get(ctx, name, v1.GetOptions{})
}

// CreateIngress creates a new ingress in a given namespace
func (s *IngressService) CreateIngress(ctx context.Context, namespace string, ingress *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	// Use the K8sClient's clientset to create an ingress
	return s.client.K8sClient().NetworkingV1().Ingresses(namespace).Create(ctx, ingress, v1.CreateOptions{})
}

// DeleteIngress deletes an ingress by name in a given namespace
func (s *IngressService) DeleteIngress(ctx context.Context, namespace, name string) error {
	// Use the K8sClient's clientset to delete an ingress
	return s.client.K8sClient().NetworkingV1().Ingresses(namespace).Delete(ctx, name, v1.DeleteOptions{})
}
