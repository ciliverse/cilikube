package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type PVCService struct {
	client kubernetes.Interface
}

func NewPVCService(client kubernetes.Interface) *PVCService {
	return &PVCService{client: client}
}

// 获取单个PVC
func (s *PVCService) Get(namespace, name string) (*corev1.PersistentVolumeClaim, error) {
	return s.client.CoreV1().PersistentVolumeClaims(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建PVC
func (s *PVCService) Create(namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {

	if pvc.Namespace != "" && pvc.Namespace != namespace {
		return nil, NewValidationError("pvc namespace conflicts with path parameter")
	}

	return s.client.CoreV1().PersistentVolumeClaims(namespace).Create(
		context.TODO(),
		pvc,
		metav1.CreateOptions{},
	)
}

// 更新PVC
func (s *PVCService) Update(namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	return s.client.CoreV1().PersistentVolumeClaims(namespace).Update(
		context.TODO(),
		pvc,
		metav1.UpdateOptions{},
	)
}

// 删除PVC
func (s *PVCService) Delete(namespace, name string) error {
	return s.client.CoreV1().PersistentVolumeClaims(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *PVCService) List(namespace, selector string, limit int64) (*corev1.PersistentVolumeClaimList, error) {
	return s.client.CoreV1().PersistentVolumeClaims(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *PVCService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.CoreV1().PersistentVolumeClaims(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}
