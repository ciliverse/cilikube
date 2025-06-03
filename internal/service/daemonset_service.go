package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// DaemonSetService 结构体不再持有 client 字段
type DaemonSetService struct {
	// 不需要 client kubernetes.Interface 字段了
}

// NewDaemonSetService 构造函数不再接收 kubernetes.Interface 参数
func NewDaemonSetService() *DaemonSetService {
	return &DaemonSetService{}
}

// 获取单个DaemonSet
func (s *DaemonSetService) Get(clientSet kubernetes.Interface, namespace, name string) (*appsv1.DaemonSet, error) {
	return clientSet.AppsV1().DaemonSets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建DaemonSet
func (s *DaemonSetService) Create(clientSet kubernetes.Interface, namespace string, daemonset *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	if daemonset.Namespace != "" && daemonset.Namespace != namespace {
		return nil, NewValidationError("daemonset namespace conflicts with path parameter")
	}

	return clientSet.AppsV1().DaemonSets(namespace).Create(
		context.TODO(),
		daemonset,
		metav1.CreateOptions{},
	)
}

// 更新DaemonSet（包含冲突检测）
func (s *DaemonSetService) Update(clientSet kubernetes.Interface, namespace string, daemonset *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return clientSet.AppsV1().DaemonSets(namespace).Update(
		context.TODO(),
		daemonset,
		metav1.UpdateOptions{},
	)
}

// 删除DaemonSet
func (s *DaemonSetService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.AppsV1().DaemonSets(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *DaemonSetService) List(clientSet kubernetes.Interface, namespace, selector string) (*appsv1.DaemonSetList, error) {
	if namespace == "" {
		namespace = corev1.NamespaceAll
	}
	return clientSet.AppsV1().DaemonSets(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
		},
	)
}

// Watch机制实现
func (s *DaemonSetService) Watch(clientSet kubernetes.Interface, namespace, selector string) (watch.Interface, error) {
	return clientSet.AppsV1().DaemonSets(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
