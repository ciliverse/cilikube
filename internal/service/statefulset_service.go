package service

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type StatefulSetService struct {
	client kubernetes.Interface
}

func NewStatefulSetService(client kubernetes.Interface) *StatefulSetService {
	return &StatefulSetService{client: client}
}

// 获取单个StatefulSet
func (s *StatefulSetService) Get(namespace, name string) (*appsv1.StatefulSet, error) {
	return s.client.AppsV1().StatefulSets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建StatefulSet
func (s *StatefulSetService) Create(namespace string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {

	if statefulSet.Namespace != "" && statefulSet.Namespace != namespace {
		return nil, NewValidationError("statefulSet namespace conflicts with path parameter")
	}

	return s.client.AppsV1().StatefulSets(namespace).Create(
		context.TODO(),
		statefulSet,
		metav1.CreateOptions{},
	)
}

// 更新StatefulSet
func (s *StatefulSetService) Update(namespace string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return s.client.AppsV1().StatefulSets(namespace).Update(
		context.TODO(),
		statefulSet,
		metav1.UpdateOptions{},
	)
}

// 删除StatefulSet
func (s *StatefulSetService) Delete(namespace, name string) error {
	return s.client.AppsV1().StatefulSets(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *StatefulSetService) List(namespace, selector string, limit int64) (*appsv1.StatefulSetList, error) {
	return s.client.AppsV1().StatefulSets(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *StatefulSetService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.AppsV1().StatefulSets(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
