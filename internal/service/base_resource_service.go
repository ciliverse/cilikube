package service

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// ResourceClient 资源客户端接口
// 为保持一致性，所有方法都接收 namespace 参数。对于非命名空间资源，实现时可以忽略该参数。
type ResourceClient[T runtime.Object] interface {
	Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (T, error)
	List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error)
	Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj T, opts metav1.CreateOptions) (T, error)
	Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj T, opts metav1.UpdateOptions) (T, error)
	Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error
	Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error)
}

// ResourceService 资源服务接口
type ResourceService[T runtime.Object] interface {
	List(clientset kubernetes.Interface, namespace, selector string, limit int64, continueToken string) (runtime.Object, error)
	Get(clientset kubernetes.Interface, namespace, name string) (T, error)
	Create(clientset kubernetes.Interface, namespace string, obj T) (T, error)
	Update(clientset kubernetes.Interface, namespace string, obj T) (T, error)
	Delete(clientset kubernetes.Interface, namespace, name string) error
	Watch(clientset kubernetes.Interface, namespace, selector string, resourceVersion string, timeoutSeconds int64) (watch.Interface, error)
}

// BaseResourceService 基础资源服务实现
type BaseResourceService[T runtime.Object] struct {
	client ResourceClient[T]
}

// NewBaseResourceService 创建基础资源服务
func NewBaseResourceService[T runtime.Object](client ResourceClient[T]) *BaseResourceService[T] {
	return &BaseResourceService[T]{
		client: client,
	}
}

// Get 获取单个资源
func (s *BaseResourceService[T]) Get(clientset kubernetes.Interface, namespace, name string) (T, error) {
	ctx := context.Background()
	return s.client.Get(ctx, clientset, namespace, name, metav1.GetOptions{})
}

// List 获取资源列表
func (s *BaseResourceService[T]) List(clientset kubernetes.Interface, namespace, selector string, limit int64, continueToken string) (runtime.Object, error) {
	ctx := context.Background()
	opts := metav1.ListOptions{
		LabelSelector: selector,
		Limit:         limit,
		Continue:      continueToken,
	}
	return s.client.List(ctx, clientset, namespace, opts)
}

// Create 创建资源
func (s *BaseResourceService[T]) Create(clientset kubernetes.Interface, namespace string, obj T) (T, error) {
	ctx := context.Background()
	return s.client.Create(ctx, clientset, namespace, obj, metav1.CreateOptions{})
}

// Update 更新资源
func (s *BaseResourceService[T]) Update(clientset kubernetes.Interface, namespace string, obj T) (T, error) {
	ctx := context.Background()
	return s.client.Update(ctx, clientset, namespace, obj, metav1.UpdateOptions{})
}

// Delete 删除资源
func (s *BaseResourceService[T]) Delete(clientset kubernetes.Interface, namespace, name string) error {
	ctx := context.Background()
	return s.client.Delete(ctx, clientset, namespace, name, metav1.DeleteOptions{})
}

// Watch 监听资源变化
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
