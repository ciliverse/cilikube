package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type PVService struct {
	client kubernetes.Interface
}

func NewPVService(client kubernetes.Interface) *PVService {
	return &PVService{client: client}
}

// 获取单个PV
func (s *PVService) Get(name string) (*corev1.PersistentVolume, error) {
	return s.client.CoreV1().PersistentVolumes().Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建PV
func (s *PVService) Create(pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	return s.client.CoreV1().PersistentVolumes().Create(
		context.TODO(),
		pv,
		metav1.CreateOptions{},
	)
}

// 更新PV
func (s *PVService) Update(pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	return s.client.CoreV1().PersistentVolumes().Update(
		context.TODO(),
		pv,
		metav1.UpdateOptions{},
	)
}

// 删除PV
func (s *PVService) Delete(name string) error {
	return s.client.CoreV1().PersistentVolumes().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *PVService) List(selector string, limit int64) (*corev1.PersistentVolumeList, error) {
	return s.client.CoreV1().PersistentVolumes().List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *PVService) Watch(selector string) (watch.Interface, error) {
	return s.client.CoreV1().PersistentVolumes().Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
