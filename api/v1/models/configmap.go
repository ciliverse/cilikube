package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateConfigMapRequest struct {
	Name        string            `json:"name" binding:"required"`
	Namespace   string            `json:"namespace" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Data        map[string]string `json:"data" binding:"required"`
}

type UpdateConfigMapRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Data        map[string]string `json:"data" binding:"required"`
}

// 响应结构
type ConfigMapResponse struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Data        map[string]string `json:"data"`
	CreatedAt   metav1.Time       `json:"createdAt"`
}

type ConfigMapListResponse struct {
	Items []ConfigMapResponse `json:"items"`
	Total int                 `json:"total"`
}

func ToConfigMapResponse(configMap *corev1.ConfigMap) ConfigMapResponse {
	return ConfigMapResponse{
		Name:        configMap.Name,
		Namespace:   configMap.Namespace,
		Labels:      configMap.Labels,
		Annotations: configMap.Annotations,
		Data:        configMap.Data,
		CreatedAt:   configMap.CreationTimestamp,
	}
}
