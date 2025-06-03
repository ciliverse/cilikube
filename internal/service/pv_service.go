package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PVService 结构体不再持有 client 字段
type PVService struct {
	// 不需要 client kubernetes.Interface 字段了
}

// NewPVService 构造函数不再接收 kubernetes.Interface 参数
func NewPVService() *PVService {
	return &PVService{}
}

// Get retrieves a single PersistentVolume by name.
func (s *PVService) Get(clientSet kubernetes.Interface, name string) (*corev1.PersistentVolume, error) {
	return clientSet.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
}

// List retrieves a list of PersistentVolumes.
// Supports label selector filtering and limit.
// Note: Pagination for cluster-scoped resources requires careful handling with 'continue' tokens
//
//	if dealing with very large numbers. For simplicity, limit is used here.
func (s *PVService) List(clientSet kubernetes.Interface, labelSelector string, limit int64) (*corev1.PersistentVolumeList, error) {
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}
	if limit > 0 {
		listOptions.Limit = limit
	}

	return clientSet.CoreV1().PersistentVolumes().List(context.TODO(), listOptions)
}

// Create creates a new PersistentVolume.
func (s *PVService) Create(clientSet kubernetes.Interface, pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	// Basic validation (optional, more can be added)
	if pv.Name == "" {
		return nil, NewValidationError("PersistentVolume name cannot be empty")
	}
	// Ensure namespace is not set for cluster-scoped resource
	pv.Namespace = ""

	return clientSet.CoreV1().PersistentVolumes().Create(context.TODO(), pv, metav1.CreateOptions{})
}

// Update updates an existing PersistentVolume.
// Note: Many PV fields are immutable after creation. Updates usually involve labels, annotations,
//
//	or potentially capacity/reclaim policy depending on the provisioner and status.
func (s *PVService) Update(clientSet kubernetes.Interface, pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	if pv.Name == "" {
		return nil, NewValidationError("PersistentVolume name cannot be empty for update")
	}
	// Ensure namespace is not set
	pv.Namespace = ""

	// Fetch existing to ensure resource version for optimistic concurrency (optional but good practice)
	// existing, err := s.Get(pv.Name)
	// if err != nil {
	//     return nil, err // Handle not found etc.
	// }
	// pv.ResourceVersion = existing.ResourceVersion // Set for update

	return clientSet.CoreV1().PersistentVolumes().Update(context.TODO(), pv, metav1.UpdateOptions{})
}

// Delete deletes a PersistentVolume by name.
func (s *PVService) Delete(clientSet kubernetes.Interface, name string) error {
	return clientSet.CoreV1().PersistentVolumes().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// --- Error Handling (reuse or define locally if not shared) ---
// type ValidationError struct { Message string }
// func (e *ValidationError) Error() string { return e.Message }
// func NewValidationError(msg string) error { return &ValidationError{Message: msg} }
