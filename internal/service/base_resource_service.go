package service

import (
	"context"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
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
	Watch(ctx context.Context, clientset kubernetes.Interface, namespace, selector string, resourceVersion string, timeoutSeconds int64) (watch.Interface, error)
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
	obj := any(current)

	if spec, ok := patchData["spec"].(map[string]interface{}); ok {
		if rawReplicas, exists := spec["replicas"]; exists {
			replicas, convOK := toInt32(rawReplicas)
			if !convOK {
				var zero T
				return zero, fmt.Errorf("invalid spec.replicas value: %v", rawReplicas)
			}
			switch v := obj.(type) {
			case *appsv1.Deployment:
				v.Spec.Replicas = &replicas
			case *appsv1.StatefulSet:
				v.Spec.Replicas = &replicas
			case *appsv1.ReplicaSet:
				v.Spec.Replicas = &replicas
			default:
				var zero T
				return zero, fmt.Errorf("scaling not supported for resource type %T", current)
			}
		}
	}

	return s.Update(clientset, namespace, name, current)
}

func toInt32(v interface{}) (int32, bool) {
	switch n := v.(type) {
	case float64:
		return int32(n), true
	case float32:
		return int32(n), true
	case int:
		return int32(n), true
	case int32:
		return n, true
	case int64:
		return int32(n), true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return int32(i), true
	default:
		return 0, false
	}
}

// Delete deletes resource
func (s *BaseResourceService[T]) Delete(clientset kubernetes.Interface, namespace, name string) error {
	ctx := context.Background()
	return s.client.Delete(ctx, clientset, namespace, name, metav1.DeleteOptions{})
}

// Watch watches resource changes
func (s *BaseResourceService[T]) Watch(ctx context.Context, clientset kubernetes.Interface, namespace, selector string, resourceVersion string, timeoutSeconds int64) (watch.Interface, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	opts := metav1.ListOptions{
		LabelSelector:   selector,
		ResourceVersion: resourceVersion,
		Watch:           true,
	}
	if timeoutSeconds > 0 {
		opts.TimeoutSeconds = &timeoutSeconds
	}
	return s.client.Watch(ctx, clientset, namespace, opts)
}
