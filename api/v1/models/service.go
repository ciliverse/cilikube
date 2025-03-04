package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateServiceRequest struct {
	Name        string             `json:"name" binding:"required"`
	Namespace   string             `json:"namespace" binding:"required"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Spec        corev1.ServiceSpec `json:"spec" binding:"required"`
}

type UpdateServiceRequest struct {
	Labels      map[string]string  `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Spec        corev1.ServiceSpec `json:"spec" binding:"required"`
}

// 响应结构
type ServiceResponse struct {
	Name        string               `json:"name"`
	Namespace   string               `json:"namespace"`
	Labels      map[string]string    `json:"labels,omitempty"`
	Annotations map[string]string    `json:"annotations,omitempty"`
	Spec        corev1.ServiceSpec   `json:"spec"`
	Status      corev1.ServiceStatus `json:"status"`
	CreatedAt   metav1.Time          `json:"createdAt"`
}

type ServiceListResponse struct {
	Items []ServiceResponse `json:"items"`
	Total int               `json:"total"`
}

func ToServiceResponse(service *corev1.Service) ServiceResponse {
	return ServiceResponse{
		Name:        service.Name,
		Namespace:   service.Namespace,
		Labels:      service.Labels,
		Annotations: service.Annotations,
		Spec:        service.Spec,
		Status:      service.Status,
		CreatedAt:   service.CreationTimestamp,
	}
}
