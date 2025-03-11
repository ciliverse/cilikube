package models

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateIngressRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Namespace   string                   `json:"namespace" binding:"required"`
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
	Spec        networkingv1.IngressSpec `json:"spec" binding:"required"`
}

type UpdateIngressRequest struct {
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
	Spec        networkingv1.IngressSpec `json:"spec" binding:"required"`
}

// 响应结构
type IngressResponse struct {
	Name        string                     `json:"name"`
	Namespace   string                     `json:"namespace"`
	Labels      map[string]string          `json:"labels,omitempty"`
	Annotations map[string]string          `json:"annotations,omitempty"`
	Spec        networkingv1.IngressSpec   `json:"spec"`
	Status      networkingv1.IngressStatus `json:"status"`
	CreatedAt   metav1.Time                `json:"createdAt"`
}

type IngressListResponse struct {
	Items []IngressResponse `json:"items"`
	Total int               `json:"total"`
}

func ToIngressResponse(ingress *networkingv1.Ingress) IngressResponse {
	return IngressResponse{
		Name:        ingress.Name,
		Namespace:   ingress.Namespace,
		Labels:      ingress.Labels,
		Annotations: ingress.Annotations,
		Spec:        ingress.Spec,
		Status:      ingress.Status,
		CreatedAt:   ingress.CreationTimestamp,
	}
}
