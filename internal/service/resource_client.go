package service

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// baseClient 基础资源客户端
type baseClient[T runtime.Object] struct{}

// Get 获取资源
func (c *baseClient[T]) Get(ctx context.Context, clientset kubernetes.Interface, name string, opts metav1.GetOptions) (T, error) {
	var zero T
	return zero, nil
}

// List 获取资源列表
func (c *baseClient[T]) List(ctx context.Context, clientset kubernetes.Interface, opts metav1.ListOptions) ([]T, error) {
	return nil, nil
}

// Create 创建资源
func (c *baseClient[T]) Create(ctx context.Context, clientset kubernetes.Interface, obj T, opts metav1.CreateOptions) (T, error) {
	var zero T
	return zero, nil
}

// Update 更新资源
func (c *baseClient[T]) Update(ctx context.Context, clientset kubernetes.Interface, obj T, opts metav1.UpdateOptions) (T, error) {
	var zero T
	return zero, nil
}

// Delete 删除资源
func (c *baseClient[T]) Delete(ctx context.Context, clientset kubernetes.Interface, name string, opts metav1.DeleteOptions) error {
	return nil
}

// Watch 监听资源变化
func (c *baseClient[T]) Watch(ctx context.Context, clientset kubernetes.Interface, opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
