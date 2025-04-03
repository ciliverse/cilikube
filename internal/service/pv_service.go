package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PVService struct {
	client kubernetes.Interface
}

func NewPVService(client kubernetes.Interface) *PVService {
	return &PVService{client: client}
}

// Get retrieves a single PersistentVolume by name.
func (s *PVService) Get(name string) (*corev1.PersistentVolume, error) {
	return s.client.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
}

// List retrieves a list of PersistentVolumes.
// Supports label selector filtering and limit.
// Note: Pagination for cluster-scoped resources requires careful handling with 'continue' tokens
//
//	if dealing with very large numbers. For simplicity, limit is used here.
func (s *PVService) List(labelSelector string, limit int64) (*corev1.PersistentVolumeList, error) {
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}
	if limit > 0 {
		listOptions.Limit = limit
	}

	return s.client.CoreV1().PersistentVolumes().List(context.TODO(), listOptions)
}

// Create creates a new PersistentVolume.
func (s *PVService) Create(pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	// Basic validation (optional, more can be added)
	if pv.Name == "" {
		return nil, NewValidationError("PersistentVolume name cannot be empty")
	}
	// Ensure namespace is not set for cluster-scoped resource
	pv.Namespace = ""

	return s.client.CoreV1().PersistentVolumes().Create(context.TODO(), pv, metav1.CreateOptions{})
}

// Update updates an existing PersistentVolume.
// Note: Many PV fields are immutable after creation. Updates usually involve labels, annotations,
//
//	or potentially capacity/reclaim policy depending on the provisioner and status.
func (s *PVService) Update(pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
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

	return s.client.CoreV1().PersistentVolumes().Update(context.TODO(), pv, metav1.UpdateOptions{})
}

// Delete deletes a PersistentVolume by name.
func (s *PVService) Delete(name string) error {
	return s.client.CoreV1().PersistentVolumes().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// --- Error Handling (reuse or define locally if not shared) ---
// type ValidationError struct { Message string }
// func (e *ValidationError) Error() string { return e.Message }
// func NewValidationError(msg string) error { return &ValidationError{Message: msg} }
