package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreatePVRequest struct {
	Name        string                      `json:"name" binding:"required"`
	Labels      map[string]string           `json:"labels,omitempty"`
	Annotations map[string]string           `json:"annotations,omitempty"`
	Spec        corev1.PersistentVolumeSpec `json:"spec" binding:"required"`
}

type UpdatePVRequest struct {
	Labels      map[string]string           `json:"labels,omitempty"`
	Annotations map[string]string           `json:"annotations,omitempty"`
	Spec        corev1.PersistentVolumeSpec `json:"spec" binding:"required"`
}

// 响应结构
type PVResponse struct {
	Name        string                        `json:"name"`
	Labels      map[string]string             `json:"labels,omitempty"`
	Annotations map[string]string             `json:"annotations,omitempty"`
	Spec        corev1.PersistentVolumeSpec   `json:"spec"`
	Status      corev1.PersistentVolumeStatus `json:"status"`
	CreatedAt   metav1.Time                   `json:"createdAt"`
}

type PVListResponse struct {
	Items []PVResponse `json:"items"`
	Total int          `json:"total"`
}

func ToPVResponse(pv *corev1.PersistentVolume) PVResponse {
	return PVResponse{
		Name:        pv.Name,
		Labels:      pv.Labels,
		Annotations: pv.Annotations,
		Spec:        pv.Spec,
		Status:      pv.Status,
		CreatedAt:   pv.CreationTimestamp,
	}
}
