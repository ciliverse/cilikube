package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/apis/models/k8s"
	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulSetService 提供对 Kubernetes StatefulSet 资源的操作
type StatefulSetService struct {
	client *client.K8sClient
}

// NewStatefulSetService 创建一个新的 StatefulSetService 实例
func NewStatefulSetService(client *client.K8sClient) *StatefulSetService {
	return &StatefulSetService{
		client: client,
	}
}

// ListStatefulSets 列出指定命名空间中的所有 StatefulSets
func (s *StatefulSetService) ListStatefulSets(namespace string) ([]k8s.StatefulSetModel, error) {
	statefulSets, err := s.client.K8sClient().AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var statefulSetModels []k8s.StatefulSetModel
	for _, statefulSet := range statefulSets.Items {
		statefulSetModels = append(statefulSetModels, *k8s.NewStatefulSetModel(&statefulSet))
	}
	return statefulSetModels, nil
}

// GetStatefulSet 获取指定命名空间中指定 StatefulSet 的详细信息
func (s *StatefulSetService) GetStatefulSet(namespace, name string) (*k8s.StatefulSetModel, error) {
	statefulSet, err := s.client.K8sClient().AppsV1().StatefulSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return k8s.NewStatefulSetModel(statefulSet), nil
}

// CreateStatefulSet 创建一个新的 StatefulSet
func (s *StatefulSetService) CreateStatefulSet(namespace string, statefulSet *appsv1.StatefulSet) (*k8s.StatefulSetModel, error) {
	createdStatefulSet, err := s.client.K8sClient().AppsV1().StatefulSets(namespace).Create(context.Background(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return k8s.NewStatefulSetModel(createdStatefulSet), nil
}

// UpdateStatefulSet 更新指定命名空间中的 StatefulSet
func (s *StatefulSetService) UpdateStatefulSet(namespace string, statefulSet *appsv1.StatefulSet) (*k8s.StatefulSetModel, error) {
	updatedStatefulSet, err := s.client.K8sClient().AppsV1().StatefulSets(namespace).Update(context.Background(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return k8s.NewStatefulSetModel(updatedStatefulSet), nil
}

// DeleteStatefulSet 删除指定命名空间中的指定 StatefulSet
func (s *StatefulSetService) DeleteStatefulSet(namespace, name string) error {
	return s.client.K8sClient().AppsV1().StatefulSets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}
