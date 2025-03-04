package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateSecretRequest struct {
	Name        string            `json:"name" binding:"required"`
	Namespace   string            `json:"namespace" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Data        map[string][]byte `json:"data" binding:"required"`
	Type        corev1.SecretType `json:"type" binding:"required"`
}

type UpdateSecretRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Data        map[string][]byte `json:"data" binding:"required"`
	Type        corev1.SecretType `json:"type" binding:"required"`
}

// 响应结构
type SecretResponse struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Data        map[string][]byte `json:"data"`
	Type        corev1.SecretType `json:"type"`
	CreatedAt   metav1.Time       `json:"createdAt"`
}

type SecretListResponse struct {
	Items []SecretResponse `json:"items"`
	Total int              `json:"total"`
}

func ToSecretResponse(secret *corev1.Secret) SecretResponse {
	return SecretResponse{
		Name:        secret.Name,
		Namespace:   secret.Namespace,
		Labels:      secret.Labels,
		Annotations: secret.Annotations,
		Data:        secret.Data,
		Type:        secret.Type,
		CreatedAt:   secret.CreationTimestamp,
	}
}
