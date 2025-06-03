package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// StatefulSetService 结构体不再持有 client 字段
type StatefulSetService struct {
	// 不需要 client kubernetes.Interface 字段了
}

// NewStatefulSetService 构造函数不再接收 kubernetes.Interface 参数
func NewStatefulSetService() *StatefulSetService {
	return &StatefulSetService{}
}

// 获取单个StatefulSet
func (s *StatefulSetService) Get(clientSet kubernetes.Interface, namespace, name string) (*appsv1.StatefulSet, error) {
	return clientSet.AppsV1().StatefulSets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建StatefulSet
func (s *StatefulSetService) Create(clientSet kubernetes.Interface, namespace string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {

	if statefulSet.Namespace != "" && statefulSet.Namespace != namespace {
		return nil, NewValidationError("statefulSet namespace conflicts with path parameter")
	}

	return clientSet.AppsV1().StatefulSets(namespace).Create(
		context.TODO(),
		statefulSet,
		metav1.CreateOptions{},
	)
}

// 更新StatefulSet
func (s *StatefulSetService) Update(clientSet kubernetes.Interface, namespace string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return clientSet.AppsV1().StatefulSets(namespace).Update(
		context.TODO(),
		statefulSet,
		metav1.UpdateOptions{},
	)
}

// 删除StatefulSet
func (s *StatefulSetService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.AppsV1().StatefulSets(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *StatefulSetService) List(clientSet kubernetes.Interface, namespace, selector string, limit int64) (*appsv1.StatefulSetList, error) {
	return clientSet.AppsV1().StatefulSets(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *StatefulSetService) Watch(clientSet kubernetes.Interface, namespace, selector string) (watch.Interface, error) {
	return clientSet.AppsV1().StatefulSets(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
