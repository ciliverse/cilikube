package models

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateNetworkPolicyRequest struct {
	Name        string                         `json:"name" binding:"required"`
	Namespace   string                         `json:"namespace" binding:"required"`
	Labels      map[string]string              `json:"labels,omitempty"`
	Annotations map[string]string              `json:"annotations,omitempty"`
	Spec        networkingv1.NetworkPolicySpec `json:"spec" binding:"required"`
}

type UpdateNetworkPolicyRequest struct {
	Labels      map[string]string              `json:"labels,omitempty"`
	Annotations map[string]string              `json:"annotations,omitempty"`
	Spec        networkingv1.NetworkPolicySpec `json:"spec" binding:"required"`
}

// 响应结构
type NetworkPolicyResponse struct {
	Name        string                         `json:"name"`
	Namespace   string                         `json:"namespace"`
	Labels      map[string]string              `json:"labels,omitempty"`
	Annotations map[string]string              `json:"annotations,omitempty"`
	Spec        networkingv1.NetworkPolicySpec `json:"spec"`
	CreatedAt   metav1.Time                    `json:"createdAt"`
}

type NetworkPolicyListResponse struct {
	Items []NetworkPolicyResponse `json:"items"`
	Total int                     `json:"total"`
}

func ToNetworkPolicyResponse(networkPolicy *networkingv1.NetworkPolicy) NetworkPolicyResponse {
	return NetworkPolicyResponse{
		Name:        networkPolicy.Name,
		Namespace:   networkPolicy.Namespace,
		Labels:      networkPolicy.Labels,
		Annotations: networkPolicy.Annotations,
		Spec:        networkPolicy.Spec,

		CreatedAt: networkPolicy.CreationTimestamp,
	}
}
