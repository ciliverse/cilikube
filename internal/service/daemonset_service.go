package service

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type DaemonSetService struct {
	client kubernetes.Interface
}

func NewDaemonSetService(client kubernetes.Interface) *DaemonSetService {
	return &DaemonSetService{client: client}
}

// 获取单个DaemonSet
func (s *DaemonSetService) Get(namespace, name string) (*appsv1.DaemonSet, error) {
	return s.client.AppsV1().DaemonSets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建DaemonSet
func (s *DaemonSetService) Create(namespace string, daemonset *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	if daemonset.Namespace != "" && daemonset.Namespace != namespace {
		return nil, NewValidationError("daemonset namespace conflicts with path parameter")
	}

	return s.client.AppsV1().DaemonSets(namespace).Create(
		context.TODO(),
		daemonset,
		metav1.CreateOptions{},
	)
}

// 更新DaemonSet（包含冲突检测）
func (s *DaemonSetService) Update(namespace string, daemonset *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return s.client.AppsV1().DaemonSets(namespace).Update(
		context.TODO(),
		daemonset,
		metav1.UpdateOptions{},
	)
}

// 删除DaemonSet
func (s *DaemonSetService) Delete(namespace, name string) error {
	return s.client.AppsV1().DaemonSets(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *DaemonSetService) List(namespace, selector string) (*appsv1.DaemonSetList, error) {
	if namespace == "" {
		namespace = corev1.NamespaceAll
	}
	return s.client.AppsV1().DaemonSets(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
		},
	)
}

// Watch机制实现
func (s *DaemonSetService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.AppsV1().DaemonSets(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
