package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreatePodRequest struct {
	Name        string            `json:"name" binding:"required"`
	Namespace   string            `json:"namespace" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        corev1.PodSpec    `json:"spec" binding:"required"`
}

type UpdatePodRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        corev1.PodSpec    `json:"spec" binding:"required"`
}

// 响应结构
type PodResponse struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Status      string            `json:"status"`
	IP          string            `json:"ip,omitempty"`
	Node        string            `json:"node,omitempty"`
	CreatedAt   metav1.Time       `json:"createdAt"`
}

type PodListResponse struct {
	Items []PodResponse `json:"items"`
	Total int           `json:"total"`
}

func ToPodResponse(pod *corev1.Pod) PodResponse {
	return PodResponse{
		Name:        pod.Name,
		Namespace:   pod.Namespace,
		Labels:      pod.Labels,
		Annotations: pod.Annotations,
		Status:      string(pod.Status.Phase),
		IP:          pod.Status.PodIP,
		Node:        pod.Spec.NodeName,
		CreatedAt:   pod.CreationTimestamp,
	}
}
