package service

import (
	"context"
	"github.com/ciliverse/cilikube/api/v1/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RbacService 结构体不再持有 client 字段
type RbacService struct {
	// 不需要 client kubernetes.Interface 字段了
}

func NewRbacService() *RbacService {
	return &RbacService{}
}

// Roles
func (s *RbacService) ListRoles(clientSet kubernetes.Interface, namespace string) ([]*models.RoleResponse, error) {
	roleList, err := clientSet.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
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
func (s *RbacService) GetRole(clientSet kubernetes.Interface, namespace string, name string) (*models.RoleResponse, error) {
	role, err := clientSet.RbacV1().Roles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToRoleResponse(role), nil
}

// RoleBindings
func (s *RbacService) ListRoleBindings(clientSet kubernetes.Interface, namespace string) ([]*models.RoleBindingResponse, error) {
	roleBindingList, err := clientSet.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var roleBindings []*models.RoleBindingResponse
	for _, roleBinding := range roleBindingList.Items {
		roleBindings = append(roleBindings, models.ToRoleBindingResponse(&roleBinding))
	}
	return roleBindings, nil
}

func (s *RbacService) GetRoleBinding(clientSet kubernetes.Interface, namespace string, name string) (*models.RoleBindingResponse, error) {
	roleBinding, err := clientSet.RbacV1().RoleBindings(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToRoleBindingResponse(roleBinding), nil
}

// ClusterRoles
func (s *RbacService) ListClusterRoles(clientSet kubernetes.Interface) ([]*models.ClusterRoleResponse, error) {
	clusterRoleList, err := clientSet.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var clusterRoles []*models.ClusterRoleResponse
	for _, clusterRole := range clusterRoleList.Items {
		clusterRoles = append(clusterRoles, models.ToClusterRoleResponse(&clusterRole))
	}
	return clusterRoles, nil
}

func (s *RbacService) GetClusterRole(clientSet kubernetes.Interface, name string) (*models.ClusterRoleResponse, error) {
	clusterRole, err := clientSet.RbacV1().ClusterRoles().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToClusterRoleResponse(clusterRole), nil
}

// ClusterRoleBindings
func (s *RbacService) ListClusterRoleBindings(clientSet kubernetes.Interface) ([]*models.ClusterRoleBindingsResponse, error) {
	clusterRoleBindingsList, err := clientSet.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var clusterRoleBindings []*models.ClusterRoleBindingsResponse
	for _, clusterRoleBinding := range clusterRoleBindingsList.Items {
		clusterRoleBindings = append(clusterRoleBindings, models.ToClusterRoleBindingsResponse(&clusterRoleBinding))
	}
	return clusterRoleBindings, nil
}

func (s *RbacService) GetClusterRoleBinding(clientSet kubernetes.Interface, name string) (*models.ClusterRoleBindingsResponse, error) {
	clusterRoleBinding, err := clientSet.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToClusterRoleBindingsResponse(clusterRoleBinding), nil
}

// ServiceAccount
func (s *RbacService) ListServiceAccounts(clientSet kubernetes.Interface, namespace string) ([]*models.ServiceAccountsResponse, error) {
	serviceAccountList, err := clientSet.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var serviceAccounts []*models.ServiceAccountsResponse
	for _, serviceAccount := range serviceAccountList.Items {
		serviceAccounts = append(serviceAccounts, models.ToServiceAccountsResponse(&serviceAccount))
	}
	return serviceAccounts, nil
}

func (s *RbacService) GetServiceAccounts(clientSet kubernetes.Interface, namespace string, name string) (*models.ServiceAccountsResponse, error) {
	serviceAccount, err := clientSet.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return models.ToServiceAccountsResponse(serviceAccount), nil
}
