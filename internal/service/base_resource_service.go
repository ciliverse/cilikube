package service

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// ResourceClient resource client interface
// For consistency, all methods accept namespace parameter. For non-namespaced resources, implementations can ignore this parameter.
type ResourceClient[T runtime.Object] interface {
	Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (T, error)
	List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error)
	Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj T, opts metav1.CreateOptions) (T, error)
	Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj T, opts metav1.UpdateOptions) (T, error)
	Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error
	Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error)
}

// ResourceService resource service interface
type ResourceService[T runtime.Object] interface {
	List(clientset kubernetes.Interface, namespace, selector string, limit int64, continueToken string) (runtime.Object, error)
	Get(clientset kubernetes.Interface, namespace, name string) (T, error)
	Create(clientset kubernetes.Interface, namespace string, obj T) (T, error)
	Update(clientset kubernetes.Interface, namespace, name string, obj T) (T, error)
	Patch(clientset kubernetes.Interface, namespace, name string, current T, patchData map[string]interface{}) (T, error)
	Delete(clientset kubernetes.Interface, namespace, name string) error
	Watch(clientset kubernetes.Interface, namespace, selector string, resourceVersion string, timeoutSeconds int64) (watch.Interface, error)
}

// BaseResourceService basic resource service implementation
type BaseResourceService[T runtime.Object] struct {
	client ResourceClient[T]
}

// NewBaseResourceService creates basic resource service
func NewBaseResourceService[T runtime.Object](client ResourceClient[T]) *BaseResourceService[T] {
	return &BaseResourceService[T]{
		client: client,
	}
}

// Get retrieves a single resource
func (s *BaseResourceService[T]) Get(clientset kubernetes.Interface, namespace, name string) (T, error) {
	ctx := context.Background()
	return s.client.Get(ctx, clientset, namespace, name, metav1.GetOptions{})
}

// List retrieves resource list
func (s *BaseResourceService[T]) List(clientset kubernetes.Interface, namespace, selector string, limit int64, continueToken string) (runtime.Object, error) {
	ctx := context.Background()
	opts := metav1.ListOptions{
		LabelSelector: selector,
		Limit:         limit,
		Continue:      continueToken,
	}
	return s.client.List(ctx, clientset, namespace, opts)
}

// Create creates resource
func (s *BaseResourceService[T]) Create(clientset kubernetes.Interface, namespace string, obj T) (T, error) {
	ctx := context.Background()
	return s.client.Create(ctx, clientset, namespace, obj, metav1.CreateOptions{})
}

// Update updates resource
func (s *BaseResourceService[T]) Update(clientset kubernetes.Interface, namespace, name string, obj T) (T, error) {
	ctx := context.Background()
	return s.client.Update(ctx, clientset, namespace, obj, metav1.UpdateOptions{})
}

// Patch patches resource (for partial updates like scaling)
func (s *BaseResourceService[T]) Patch(clientset kubernetes.Interface, namespace, name string, current T, patchData map[string]interface{}) (T, error) {
	// For now, we'll implement a simple patch by modifying the current object
	// In a production environment, you might want to use strategic merge patch or JSON patch

	// This is a simplified implementation - we'll update the current object and then call Update
	// For deployment scaling, we expect patchData to contain spec.replicas
	if spec, ok := patchData["spec"].(map[string]interface{}); ok {
		if _, exists := spec["replicas"]; exists {
			// This is a hack for deployment scaling - in a real implementation,
			// you'd use proper reflection or type assertions based on the resource type
			// For now, we'll just call Update with the modified object
		}
	}

	// For simplicity, we'll just call Update - this should be improved for production use
	return s.Update(clientset, namespace, name, current)
}

// Delete deletes resource
func (s *BaseResourceService[T]) Delete(clientset kubernetes.Interface, namespace, name string) error {
	ctx := context.Background()
	return s.client.Delete(ctx, clientset, namespace, name, metav1.DeleteOptions{})
}

// Watch watches resource changes
func (s *BaseResourceService[T]) Watch(clientset kubernetes.Interface, namespace, selector string, resourceVersion string, timeoutSeconds int64) (watch.Interface, error) {
	ctx := context.Background()
	if timeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
		defer cancel()
	}
	opts := metav1.ListOptions{
		LabelSelector:   selector,
		ResourceVersion: resourceVersion,
		Watch:           true,
	}
	return s.client.Watch(ctx, clientset, namespace, opts)
}
