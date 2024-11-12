package services

import (
	"context"

	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StorageClassService struct {
	client *client.K8sClient
}

func NewStorageClassService(client *client.K8sClient) *StorageClassService {
	return &StorageClassService{
		client: client,
	}
}

func (s *StorageClassService) CreateStorageClass(sc *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	// 调用 client.K8sClient() 获取 Kubernetes 客户端实例
	return s.client.K8sClient().StorageV1().StorageClasses().Create(context.TODO(), sc, metav1.CreateOptions{})
}

func (s *StorageClassService) GetStorageClass(name string) (*storagev1.StorageClass, error) {
	// 调用 client.K8sClient() 获取 Kubernetes 客户端实例
	return s.client.K8sClient().StorageV1().StorageClasses().Get(context.TODO(), name, metav1.GetOptions{})
}

func (s *StorageClassService) UpdateStorageClass(sc *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	// 调用 client.K8sClient() 获取 Kubernetes 客户端实例
	return s.client.K8sClient().StorageV1().StorageClasses().Update(context.TODO(), sc, metav1.UpdateOptions{})
}

func (s *StorageClassService) DeleteStorageClass(name string) error {
	// 调用 client.K8sClient() 获取 Kubernetes 客户端实例
	return s.client.K8sClient().StorageV1().StorageClasses().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
