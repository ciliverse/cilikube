package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PVService 提供对 Kubernetes PersistentVolume 资源的操作
type PVService struct {
	client *client.K8sClient
}

// NewPVService 创建一个新的 PVService 实例
func NewPVService(k8sClient *client.K8sClient) *PVService {
	return &PVService{
		client: k8sClient,
	}
}

// CreatePV 创建一个新的 PersistentVolume
func (s *PVService) CreatePV(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	return s.client.K8sClient().CoreV1().PersistentVolumes().Create(context.TODO(), pv, metav1.CreateOptions{})
}

// GetPV 获取指定名称的 PersistentVolume
func (s *PVService) GetPV(name string) (*v1.PersistentVolume, error) {
	return s.client.K8sClient().CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
}

// UpdatePV 更新指定的 PersistentVolume
func (s *PVService) UpdatePV(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	return s.client.K8sClient().CoreV1().PersistentVolumes().Update(context.TODO(), pv, metav1.UpdateOptions{})
}

// DeletePV 删除指定名称的 PersistentVolume
func (s *PVService) DeletePV(name string) error {
	return s.client.K8sClient().CoreV1().PersistentVolumes().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
