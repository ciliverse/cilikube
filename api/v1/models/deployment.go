package models

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateDeploymentRequest struct {
	Name        string                `json:"name" binding:"required"`
	Namespace   string                `json:"namespace" binding:"required"`
	Labels      map[string]string     `json:"labels,omitempty"`
	Annotations map[string]string     `json:"annotations,omitempty"`
	Spec        appsv1.DeploymentSpec `json:"spec" binding:"required"`
}

type UpdateDeploymentRequest struct {
	Labels      map[string]string     `json:"labels,omitempty"`
	Annotations map[string]string     `json:"annotations,omitempty"`
	Spec        appsv1.DeploymentSpec `json:"spec" binding:"required"`
}

// 响应结构
type DeploymentResponse struct {
	Name                string            `json:"name"`
	Namespace           string            `json:"namespace"`
	Labels              map[string]string `json:"labels,omitempty"`
	Annotations         map[string]string `json:"annotations,omitempty"`
	Status              string            `json:"status"`
	Replicas            int32             `json:"replicas"`
	ReadyReplicas       int32             `json:"readyReplicas"`
	AvailableReplicas   int32             `json:"availableReplicas"`
	UnavailableReplicas int32             `json:"unavailableReplicas"`
	CreatedAt           metav1.Time       `json:"createdAt"`
}

type DeploymentListResponse struct {
	Items []DeploymentResponse `json:"items"`
	Total int                  `json:"total"`
}

type ScaleDeploymentRequest struct {
	Replicas int32 `json:"replicas" binding:"required"`
}

func ToDeploymentResponse(deployment *appsv1.Deployment) DeploymentResponse {
	return DeploymentResponse{
		Name:                deployment.Name,
		Namespace:           deployment.Namespace,
		Labels:              deployment.Labels,
		Annotations:         deployment.Annotations,
		Replicas:            *deployment.Spec.Replicas,
		ReadyReplicas:       deployment.Status.ReadyReplicas,
		AvailableReplicas:   deployment.Status.AvailableReplicas,
		UnavailableReplicas: deployment.Status.UnavailableReplicas,
		CreatedAt:           deployment.CreationTimestamp,
	}
}
