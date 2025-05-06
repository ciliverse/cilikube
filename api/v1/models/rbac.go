package models

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"time"
)

type RoleResponse struct {
	Name        string              `json:"name"`
	Namespace   string              `json:"namespace"`
	UID         string              `json:"uid"`
	Rules       []rbacv1.PolicyRule `json:"rules"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Annotations map[string]string   `json:"annotations,omitempty"`
	CreatedAt   string              `json:"createdAt"`
}

type RoleBindingResponse struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	UID         string            `json:"uid"`
	RoleRef     rbacv1.RoleRef    `json:"roleRef"`
	Subjects    []rbacv1.Subject  `json:"subjects"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   string            `json:"createdAt"`
}

type ClusterRoleResponse struct {
	Name        string              `json:"name"`
	UID         string              `json:"uid"`
	Rules       []rbacv1.PolicyRule `json:"rules"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Annotations map[string]string   `json:"annotations,omitempty"`
	CreatedAt   string              `json:"createdAt"`
}

type ClusterRoleBindingsResponse struct {
	Name        string            `json:"name"`
	UID         string            `json:"uid"`
	RoleRef     rbacv1.RoleRef    `json:"roleRef"`
	Subjects    []rbacv1.Subject  `json:"subjects"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   string            `json:"createdAt"`
}

type ServiceAccountsResponse struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	UID         string            `json:"uid"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   string            `json:"createdAt"`
}

func ToRoleResponse(role *rbacv1.Role) *RoleResponse {
	return &RoleResponse{
		Name:        role.Name,
		Namespace:   role.Namespace,
		UID:         string(role.UID),
		Rules:       role.Rules,
		Labels:      role.Labels,
		Annotations: role.Annotations,
		CreatedAt:   role.CreationTimestamp.Format(time.RFC3339),
	}
}

func ToRoleBindingResponse(roleBinding *rbacv1.RoleBinding) *RoleBindingResponse {
	return &RoleBindingResponse{
		Name:        roleBinding.Name,
		Namespace:   roleBinding.Namespace,
		UID:         string(roleBinding.UID),
		RoleRef:     roleBinding.RoleRef,
		Subjects:    roleBinding.Subjects,
		Labels:      roleBinding.Labels,
		Annotations: roleBinding.Annotations,
		CreatedAt:   roleBinding.CreationTimestamp.Format(time.RFC3339),
	}
}
func ToClusterRoleResponse(clusterRole *rbacv1.ClusterRole) *ClusterRoleResponse {
	return &ClusterRoleResponse{
		Name:        clusterRole.Name,
		UID:         string(clusterRole.UID),
		Rules:       clusterRole.Rules,
		Labels:      clusterRole.Labels,
		Annotations: clusterRole.Annotations,
		CreatedAt:   clusterRole.CreationTimestamp.Format(time.RFC3339),
	}
}

func ToClusterRoleBindingsResponse(clusterRoleBinding *rbacv1.ClusterRoleBinding) *ClusterRoleBindingsResponse {
	return &ClusterRoleBindingsResponse{
		Name:        clusterRoleBinding.Name,
		UID:         string(clusterRoleBinding.UID),
		RoleRef:     clusterRoleBinding.RoleRef,
		Subjects:    clusterRoleBinding.Subjects,
		Labels:      clusterRoleBinding.Labels,
		Annotations: clusterRoleBinding.Annotations,
		CreatedAt:   clusterRoleBinding.CreationTimestamp.Format(time.RFC3339),
	}
}

func ToServiceAccountsResponse(serviceAccount *corev1.ServiceAccount) *ServiceAccountsResponse {
	return &ServiceAccountsResponse{
		Name:        serviceAccount.Name,
		Namespace:   serviceAccount.Namespace,
		UID:         string(serviceAccount.UID),
		Labels:      serviceAccount.Labels,
		Annotations: serviceAccount.Annotations,
		CreatedAt:   serviceAccount.CreationTimestamp.Format(time.RFC3339),
	}
}

type CreateRoleRequest struct {
	Name        string              `json:"name" binding:"required"`
	Namespace   string              `json:"namespace" binding:"required"`
	Rules       []rbacv1.PolicyRule `json:"rules" binding:"required"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Annotations map[string]string   `json:"annotations,omitempty"`
}

type UpdateRoleRequest struct {
	Rules       []rbacv1.PolicyRule `json:"rules" binding:"required"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Annotations map[string]string   `json:"annotations,omitempty"`
}
type CreateRoleBindingRequest struct {
	Name        string            `json:"name" binding:"required"`
	Namespace   string            `json:"namespace" binding:"required"`
	RoleRef     rbacv1.RoleRef    `json:"roleRef" binding:"required"`
	Subjects    []rbacv1.Subject  `json:"subjects" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type UpdateRoleBindingRequest struct {
	RoleRef     rbacv1.RoleRef    `json:"roleRef" binding:"required"`
	Subjects    []rbacv1.Subject  `json:"subjects" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type CreateClusterRoleRequest struct {
	Name        string              `json:"name" binding:"required"`
	Rules       []rbacv1.PolicyRule `json:"rules" binding:"required"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Annotations map[string]string   `json:"annotations,omitempty"`
}

type UpdateClusterRoleRequest struct {
	Rules       []rbacv1.PolicyRule `json:"rules" binding:"required"`
	Labels      map[string]string   `json:"labels,omitempty"`
	Annotations map[string]string   `json:"annotations,omitempty"`
}

type CreateClusterRoleBindingRequest struct {
	Name        string            `json:"name" binding:"required"`
	RoleRef     rbacv1.RoleRef    `json:"roleRef" binding:"required"`
	Subjects    []rbacv1.Subject  `json:"subjects" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type UpdateClusterRoleBindingRequest struct {
	RoleRef     rbacv1.RoleRef    `json:"roleRef" binding:"required"`
	Subjects    []rbacv1.Subject  `json:"subjects" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type CreateServiceAccountRequest struct {
	Name        string            `json:"name" binding:"required"`
	Namespace   string            `json:"namespace" binding:"required"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type UpdateServiceAccountRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
