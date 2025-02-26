package models

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 请求结构
type CreateDaemonSetRequest struct {
	Name        string               `json:"name" binding:"required"`
	Namespace   string               `json:"namespace" binding:"required"`
	Labels      map[string]string    `json:"labels,omitempty"`
	Annotations map[string]string    `json:"annotations,omitempty"`
	Spec        appsv1.DaemonSetSpec `json:"spec" binding:"required"`
}

type UpdateDaemonSetRequest struct {
	Labels      map[string]string    `json:"labels,omitempty"`
	Annotations map[string]string    `json:"annotations,omitempty"`
	Spec        appsv1.DaemonSetSpec `json:"spec" binding:"required"`
}

// 响应结构
type DaemonSetResponse struct {
	Name                   string            `json:"name"`
	Namespace              string            `json:"namespace"`
	Labels                 map[string]string `json:"labels,omitempty"`
	Annotations            map[string]string `json:"annotations,omitempty"`
	Status                 string            `json:"status"`
	CurrentNumberScheduled int32             `json:"currentNumberScheduled"`
	NumberMisscheduled     int32             `json:"numberMisscheduled"`
	DesiredNumberScheduled int32             `json:"desiredNumberScheduled"`
	NumberReady            int32             `json:"numberReady"`
	CreatedAt              metav1.Time       `json:"createdAt"`
}

type DaemonSetListResponse struct {
	Items []DaemonSetResponse `json:"items"`
	Total int                 `json:"total"`
}

func ToDaemonSetResponse(daemonset *appsv1.DaemonSet) DaemonSetResponse {
	return DaemonSetResponse{
		Name:                   daemonset.Name,
		Namespace:              daemonset.Namespace,
		Labels:                 daemonset.Labels,
		Annotations:            daemonset.Annotations,
		CurrentNumberScheduled: daemonset.Status.CurrentNumberScheduled,
		NumberMisscheduled:     daemonset.Status.NumberMisscheduled,
		DesiredNumberScheduled: daemonset.Status.DesiredNumberScheduled,
		NumberReady:            daemonset.Status.NumberReady,
		CreatedAt:              daemonset.CreationTimestamp,
	}
}
