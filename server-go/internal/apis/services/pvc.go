package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PVCService 提供对 Kubernetes PersistentVolumeClaim 资源的操作
type PVCService struct {
	client *client.K8sClient
}

// NewPVCService 创建一个新的 PVCService 实例
func NewPVCService(k8sClient *client.K8sClient) *PVCService {
	return &PVCService{
		client: k8sClient,
	}
}

// CreatePVC 创建一个新的 PersistentVolumeClaim
func (s *PVCService) CreatePVC(namespace string, pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return s.client.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
}

// GetPVC 获取指定命名空间和名称的 PersistentVolumeClaim
func (s *PVCService) GetPVC(namespace, name string) (*v1.PersistentVolumeClaim, error) {
	return s.client.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// UpdatePVC 更新指定的 PersistentVolumeClaim
func (s *PVCService) UpdatePVC(namespace string, pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return s.client.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvc, metav1.UpdateOptions{})
}

// DeletePVC 删除指定命名空间和名称的 PersistentVolumeClaim
func (s *PVCService) DeletePVC(namespace, name string) error {
	return s.client.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
