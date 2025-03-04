package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreatePVCRequest struct {
	Name        string                           `json:"name" binding:"required"`
	Namespace   string                           `json:"namespace" binding:"required"`
	Labels      map[string]string                `json:"labels,omitempty"`
	Annotations map[string]string                `json:"annotations,omitempty"`
	Spec        corev1.PersistentVolumeClaimSpec `json:"spec" binding:"required"`
}

type UpdatePVCRequest struct {
	Labels      map[string]string                `json:"labels,omitempty"`
	Annotations map[string]string                `json:"annotations,omitempty"`
	Spec        corev1.PersistentVolumeClaimSpec `json:"spec" binding:"required"`
}

// 响应结构
type PVCResponse struct {
	Name        string                             `json:"name"`
	Namespace   string                             `json:"namespace"`
	Labels      map[string]string                  `json:"labels,omitempty"`
	Annotations map[string]string                  `json:"annotations,omitempty"`
	Spec        corev1.PersistentVolumeClaimSpec   `json:"spec"`
	Status      corev1.PersistentVolumeClaimStatus `json:"status"`
	CreatedAt   metav1.Time                        `json:"createdAt"`
}

type PVCListResponse struct {
	Items []PVCResponse `json:"items"`
	Total int           `json:"total"`
}

func ToPVCResponse(pvc *corev1.PersistentVolumeClaim) PVCResponse {
	return PVCResponse{
		Name:        pvc.Name,
		Namespace:   pvc.Namespace,
		Labels:      pvc.Labels,
		Annotations: pvc.Annotations,
		Spec:        pvc.Spec,
		Status:      pvc.Status,
		CreatedAt:   pvc.CreationTimestamp,
	}
}
