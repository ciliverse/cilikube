package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateNamespaceRequest struct {
	Name        string            `json:"name" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type UpdateNamespaceRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// 响应结构
type NamespaceResponse struct {
	Name        string                 `json:"name"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	Status      corev1.NamespaceStatus `json:"status"`
	CreatedAt   metav1.Time            `json:"createdAt"`
}

type NamespaceListResponse struct {
	Items []NamespaceResponse `json:"items"`
	Total int                 `json:"total"`
}

func ToNamespaceResponse(namespace *corev1.Namespace) NamespaceResponse {
	return NamespaceResponse{
		Name:        namespace.Name,
		Labels:      namespace.Labels,
		Annotations: namespace.Annotations,
		Status:      namespace.Status,
		CreatedAt:   namespace.CreationTimestamp,
	}
}
