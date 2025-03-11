package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateNodeRequest struct {
	Name        string            `json:"name" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        corev1.NodeSpec   `json:"spec" binding:"required"`
}

type UpdateNodeRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        corev1.NodeSpec   `json:"spec" binding:"required"`
}

// 响应结构
type NodeResponse struct {
	Name        string            `json:"name"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        corev1.NodeSpec   `json:"spec"`
	Status      corev1.NodeStatus `json:"status"`
	CreatedAt   metav1.Time       `json:"createdAt"`
}

type NodeListResponse struct {
	Items []NodeResponse `json:"items"`
	Total int            `json:"total"`
}

func ToNodeResponse(node *corev1.Node) NodeResponse {
	return NodeResponse{
		Name:        node.Name,
		Labels:      node.Labels,
		Annotations: node.Annotations,
		Spec:        node.Spec,
		Status:      node.Status,
		CreatedAt:   node.CreationTimestamp,
	}
}
