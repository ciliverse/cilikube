package service

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// baseClient base resource client
type baseClient[T runtime.Object] struct{}

// Get retrieves a resource
func (c *baseClient[T]) Get(ctx context.Context, clientset kubernetes.Interface, name string, opts metav1.GetOptions) (T, error) {
	var zero T
	return zero, nil
}

// List retrieves a list of resources
func (c *baseClient[T]) List(ctx context.Context, clientset kubernetes.Interface, opts metav1.ListOptions) ([]T, error) {
	return nil, nil
}

// Create creates a resource
func (c *baseClient[T]) Create(ctx context.Context, clientset kubernetes.Interface, obj T, opts metav1.CreateOptions) (T, error) {
	var zero T
	return zero, nil
}

// Update updates a resource
func (c *baseClient[T]) Update(ctx context.Context, clientset kubernetes.Interface, obj T, opts metav1.UpdateOptions) (T, error) {
	var zero T
	return zero, nil
}

// Delete deletes a resource
func (c *baseClient[T]) Delete(ctx context.Context, clientset kubernetes.Interface, name string, opts metav1.DeleteOptions) error {
	return nil
}

// Watch watches for resource changes
func (c *baseClient[T]) Watch(ctx context.Context, clientset kubernetes.Interface, opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
