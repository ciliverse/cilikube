// package models

// import (
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// // 请求结构
// type CreatePVCRequest struct {
// 	Name        string                           `json:"name" binding:"required"`
// 	Namespace   string                           `json:"namespace" binding:"required"`
// 	Labels      map[string]string                `json:"labels,omitempty"`
// 	Annotations map[string]string                `json:"annotations,omitempty"`
// 	Spec        corev1.PersistentVolumeClaimSpec `json:"spec" binding:"required"`
// }

// type UpdatePVCRequest struct {
// 	Labels      map[string]string                `json:"labels,omitempty"`
// 	Annotations map[string]string                `json:"annotations,omitempty"`
// 	Spec        corev1.PersistentVolumeClaimSpec `json:"spec" binding:"required"`
// }

// // 响应结构
// type PVCResponse struct {
// 	Name        string                             `json:"name"`
// 	Namespace   string                             `json:"namespace"`
// 	Labels      map[string]string                  `json:"labels,omitempty"`
// 	Annotations map[string]string                  `json:"annotations,omitempty"`
// 	Spec        corev1.PersistentVolumeClaimSpec   `json:"spec"`
// 	Status      corev1.PersistentVolumeClaimStatus `json:"status"`
// 	CreatedAt   metav1.Time                        `json:"createdAt"`
// }

// type PVCListResponse struct {
// 	Items []PVCResponse `json:"items"`
// 	Total int           `json:"total"`
// }

//	func ToPVCResponse(pvc *corev1.PersistentVolumeClaim) PVCResponse {
//		return PVCResponse{
//			Name:        pvc.Name,
//			Namespace:   pvc.Namespace,
//			Labels:      pvc.Labels,
//			Annotations: pvc.Annotations,
//			Spec:        pvc.Spec,
//			Status:      pvc.Status,
//			CreatedAt:   pvc.CreationTimestamp,
//		}
//	}
package models

import (
	corev1 "k8s.io/api/core/v1"
)

// PVCResponse represents the data structure for a PVC sent to the frontend.
type PVCResponse struct {
	Name             string            `json:"name"`
	Namespace        string            `json:"namespace"`
	UID              string            `json:"uid"`
	Status           string            `json:"status"`           // Phase: Pending, Bound, Lost
	VolumeName       string            `json:"volumeName"`       // Bound PV name
	StorageClassName *string           `json:"storageClass"`     // Pointer because it can be omitted
	AccessModes      []string          `json:"accessModes"`      // Simpler string array
	RequestedStorage string            `json:"requestedStorage"` // User requested size
	ActualCapacity   string            `json:"actualCapacity"`   // Actual size from bound PV (if available)
	VolumeMode       string            `json:"volumeMode"`       // Filesystem or Block
	CreatedAt        string            `json:"createdAt"`
	Labels           map[string]string `json:"labels,omitempty"`
	Annotations      map[string]string `json:"annotations,omitempty"`
	ResourceVersion  string            `json:"resourceVersion"`
}

// PVCListResponse is the response structure for listing PVCs.
type PVCListResponse struct {
	Items []PVCResponse `json:"items"`
	Total int           `json:"total"`
}

// ToPVCResponse converts a corev1.PersistentVolumeClaim to our PVCResponse model.
func ToPVCResponse(pvc *corev1.PersistentVolumeClaim) PVCResponse {
	requestedStorageStr := ""
	if storageRequest, ok := pvc.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
		requestedStorageStr = storageRequest.String()
	}

	actualCapacityStr := ""
	if pvc.Status.Phase == corev1.ClaimBound { // Only show actual capacity if bound
		if capacity, ok := pvc.Status.Capacity[corev1.ResourceStorage]; ok {
			actualCapacityStr = capacity.String()
		}
	}

	accessModesStr := make([]string, len(pvc.Spec.AccessModes))
	for i, mode := range pvc.Spec.AccessModes {
		accessModesStr[i] = string(mode)
	}

	volumeModeStr := string(corev1.PersistentVolumeFilesystem) // Default
	if pvc.Spec.VolumeMode != nil {
		volumeModeStr = string(*pvc.Spec.VolumeMode)
	}

	return PVCResponse{
		Name:             pvc.Name,
		Namespace:        pvc.Namespace,
		UID:              string(pvc.UID),
		Status:           string(pvc.Status.Phase),
		VolumeName:       pvc.Spec.VolumeName, // Name of the PV it's bound to
		StorageClassName: pvc.Spec.StorageClassName,
		AccessModes:      accessModesStr,
		RequestedStorage: requestedStorageStr,
		ActualCapacity:   actualCapacityStr,
		VolumeMode:       volumeModeStr,
		CreatedAt:        pvc.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		Labels:           pvc.Labels,
		Annotations:      pvc.Annotations,
		ResourceVersion:  pvc.ResourceVersion,
	}
}
