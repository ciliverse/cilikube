package models

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateStatefulSetRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Namespace   string                 `json:"namespace" binding:"required"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	Spec        appsv1.StatefulSetSpec `json:"spec" binding:"required"`
}

type UpdateStatefulSetRequest struct {
	Labels      map[string]string      `json:"labels,omitempty"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	Spec        appsv1.StatefulSetSpec `json:"spec" binding:"required"`
}

// 响应结构
type StatefulSetResponse struct {
	Name        string                   `json:"name"`
	Namespace   string                   `json:"namespace"`
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
	Spec        appsv1.StatefulSetSpec   `json:"spec"`
	Status      appsv1.StatefulSetStatus `json:"status"`
	CreatedAt   metav1.Time              `json:"createdAt"`
}

type StatefulSetListResponse struct {
	Items []StatefulSetResponse `json:"items"`
	Total int                   `json:"total"`
}

func ToStatefulSetResponse(statefulSet *appsv1.StatefulSet) StatefulSetResponse {
	return StatefulSetResponse{
		Name:        statefulSet.Name,
		Namespace:   statefulSet.Namespace,
		Labels:      statefulSet.Labels,
		Annotations: statefulSet.Annotations,
		Spec:        statefulSet.Spec,
		Status:      statefulSet.Status,
		CreatedAt:   statefulSet.CreationTimestamp,
	}
}
