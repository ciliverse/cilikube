package service

import (
	"context"
	"github.com/ciliverse/cilikube/api/v1/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RbacService struct {
	client kubernetes.Interface
}

func NewRbacService(client kubernetes.Interface) *RbacService { return &RbacService{client: client} }

// Roles
func (s *RbacService) ListRoles(namespace string) ([]*models.RoleResponse, error) {
	roleList, err := s.client.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var roles []*models.RoleResponse
	for _, role := range roleList.Items {
		roles = append(roles, models.ToRoleResponse(&role))
	}
	return roles, nil
}

// GetRole retrieves a single Role by namespace and name.
func (s *RbacService) GetRole(namespace string, name string) (*models.RoleResponse, error) {
	role, err := s.client.RbacV1().Roles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToRoleResponse(role), nil
}

// RoleBindings
func (s *RbacService) ListRoleBindings(namespace string) ([]*models.RoleBindingResponse, error) {
	roleBindingList, err := s.client.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var roleBindings []*models.RoleBindingResponse
	for _, roleBinding := range roleBindingList.Items {
		roleBindings = append(roleBindings, models.ToRoleBindingResponse(&roleBinding))
	}
	return roleBindings, nil
}

func (s *RbacService) GetRoleBinding(namespace string, name string) (*models.RoleBindingResponse, error) {
	roleBinding, err := s.client.RbacV1().RoleBindings(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToRoleBindingResponse(roleBinding), nil
}

// ClusterRoles
func (s *RbacService) ListClusterRoles() ([]*models.ClusterRoleResponse, error) {
	clusterRoleList, err := s.client.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var clusterRoles []*models.ClusterRoleResponse
	for _, clusterRole := range clusterRoleList.Items {
		clusterRoles = append(clusterRoles, models.ToClusterRoleResponse(&clusterRole))
	}
	return clusterRoles, nil
}

func (s *RbacService) GetClusterRole(name string) (*models.ClusterRoleResponse, error) {
	clusterRole, err := s.client.RbacV1().ClusterRoles().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToClusterRoleResponse(clusterRole), nil
}

// ClusterRoleBindings
func (s *RbacService) ListClusterRoleBindings() ([]*models.ClusterRoleBindingsResponse, error) {
	clusterRoleBindingsList, err := s.client.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var clusterRoleBindings []*models.ClusterRoleBindingsResponse
	for _, clusterRoleBinding := range clusterRoleBindingsList.Items {
		clusterRoleBindings = append(clusterRoleBindings, models.ToClusterRoleBindingsResponse(&clusterRoleBinding))
	}
	return clusterRoleBindings, nil
}

func (s *RbacService) GetClusterRoleBinding(name string) (*models.ClusterRoleBindingsResponse, error) {
	clusterRoleBinding, err := s.client.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToClusterRoleBindingsResponse(clusterRoleBinding), nil
}

// ServiceAccount
func (s *RbacService) ListServiceAccounts(namespace string) ([]*models.ServiceAccountsResponse, error) {
	serviceAccountList, err := s.client.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var serviceAccounts []*models.ServiceAccountsResponse
	for _, serviceAccount := range serviceAccountList.Items {
		serviceAccounts = append(serviceAccounts, models.ToServiceAccountsResponse(&serviceAccount))
	}
	return serviceAccounts, nil
}

func (s *RbacService) GetServiceAccounts(namespace string, name string) (*models.ServiceAccountsResponse, error) {
	serviceAccount, err := s.client.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToServiceAccountsResponse(serviceAccount), nil
}
