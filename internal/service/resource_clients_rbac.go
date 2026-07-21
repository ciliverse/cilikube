package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// --- StorageClassClient (Cluster-scoped) ---
type StorageClassClient struct{}

func (c *StorageClassClient) Get(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.GetOptions) (*storagev1.StorageClass, error) {
	return clientset.StorageV1().StorageClasses().Get(ctx, name, opts)
}
func (c *StorageClassClient) List(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.StorageV1().StorageClasses().List(ctx, opts)
}
func (c *StorageClassClient) Create(ctx context.Context, clientset kubernetes.Interface, _ string, obj *storagev1.StorageClass, opts metav1.CreateOptions) (*storagev1.StorageClass, error) {
	return clientset.StorageV1().StorageClasses().Create(ctx, obj, opts)
}
func (c *StorageClassClient) Update(ctx context.Context, clientset kubernetes.Interface, _ string, obj *storagev1.StorageClass, opts metav1.UpdateOptions) (*storagev1.StorageClass, error) {
	return clientset.StorageV1().StorageClasses().Update(ctx, obj, opts)
}
func (c *StorageClassClient) Delete(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.DeleteOptions) error {
	return clientset.StorageV1().StorageClasses().Delete(ctx, name, opts)
}
func (c *StorageClassClient) Watch(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.StorageV1().StorageClasses().Watch(ctx, opts)
}

// --- ServiceAccountClient (Namespaced) ---
type ServiceAccountClient struct{}

func (c *ServiceAccountClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*corev1.ServiceAccount, error) {
	return clientset.CoreV1().ServiceAccounts(namespace).Get(ctx, name, opts)
}
func (c *ServiceAccountClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().ServiceAccounts(namespace).List(ctx, opts)
}
func (c *ServiceAccountClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.ServiceAccount, opts metav1.CreateOptions) (*corev1.ServiceAccount, error) {
	return clientset.CoreV1().ServiceAccounts(namespace).Create(ctx, obj, opts)
}
func (c *ServiceAccountClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.ServiceAccount, opts metav1.UpdateOptions) (*corev1.ServiceAccount, error) {
	return clientset.CoreV1().ServiceAccounts(namespace).Update(ctx, obj, opts)
}
func (c *ServiceAccountClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, opts)
}
func (c *ServiceAccountClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().ServiceAccounts(namespace).Watch(ctx, opts)
}

// --- RoleClient (Namespaced) ---
type RoleClient struct{}

func (c *RoleClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*rbacv1.Role, error) {
	return clientset.RbacV1().Roles(namespace).Get(ctx, name, opts)
}
func (c *RoleClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.RbacV1().Roles(namespace).List(ctx, opts)
}
func (c *RoleClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *rbacv1.Role, opts metav1.CreateOptions) (*rbacv1.Role, error) {
	return clientset.RbacV1().Roles(namespace).Create(ctx, obj, opts)
}
func (c *RoleClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *rbacv1.Role, opts metav1.UpdateOptions) (*rbacv1.Role, error) {
	return clientset.RbacV1().Roles(namespace).Update(ctx, obj, opts)
}
func (c *RoleClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.RbacV1().Roles(namespace).Delete(ctx, name, opts)
}
func (c *RoleClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.RbacV1().Roles(namespace).Watch(ctx, opts)
}

// --- RoleBindingClient (Namespaced) ---
type RoleBindingClient struct{}

func (c *RoleBindingClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*rbacv1.RoleBinding, error) {
	return clientset.RbacV1().RoleBindings(namespace).Get(ctx, name, opts)
}
func (c *RoleBindingClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.RbacV1().RoleBindings(namespace).List(ctx, opts)
}
func (c *RoleBindingClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *rbacv1.RoleBinding, opts metav1.CreateOptions) (*rbacv1.RoleBinding, error) {
	return clientset.RbacV1().RoleBindings(namespace).Create(ctx, obj, opts)
}
func (c *RoleBindingClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *rbacv1.RoleBinding, opts metav1.UpdateOptions) (*rbacv1.RoleBinding, error) {
	return clientset.RbacV1().RoleBindings(namespace).Update(ctx, obj, opts)
}
func (c *RoleBindingClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.RbacV1().RoleBindings(namespace).Delete(ctx, name, opts)
}
func (c *RoleBindingClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.RbacV1().RoleBindings(namespace).Watch(ctx, opts)
}

// --- ClusterRoleClient (Cluster-scoped) ---
type ClusterRoleClient struct{}

func (c *ClusterRoleClient) Get(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.GetOptions) (*rbacv1.ClusterRole, error) {
	return clientset.RbacV1().ClusterRoles().Get(ctx, name, opts)
}
func (c *ClusterRoleClient) List(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.RbacV1().ClusterRoles().List(ctx, opts)
}
func (c *ClusterRoleClient) Create(ctx context.Context, clientset kubernetes.Interface, _ string, obj *rbacv1.ClusterRole, opts metav1.CreateOptions) (*rbacv1.ClusterRole, error) {
	return clientset.RbacV1().ClusterRoles().Create(ctx, obj, opts)
}
func (c *ClusterRoleClient) Update(ctx context.Context, clientset kubernetes.Interface, _ string, obj *rbacv1.ClusterRole, opts metav1.UpdateOptions) (*rbacv1.ClusterRole, error) {
	return clientset.RbacV1().ClusterRoles().Update(ctx, obj, opts)
}
func (c *ClusterRoleClient) Delete(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.DeleteOptions) error {
	return clientset.RbacV1().ClusterRoles().Delete(ctx, name, opts)
}
func (c *ClusterRoleClient) Watch(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.RbacV1().ClusterRoles().Watch(ctx, opts)
}

// --- ClusterRoleBindingClient (Cluster-scoped) ---
type ClusterRoleBindingClient struct{}

func (c *ClusterRoleBindingClient) Get(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.GetOptions) (*rbacv1.ClusterRoleBinding, error) {
	return clientset.RbacV1().ClusterRoleBindings().Get(ctx, name, opts)
}
func (c *ClusterRoleBindingClient) List(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.RbacV1().ClusterRoleBindings().List(ctx, opts)
}
func (c *ClusterRoleBindingClient) Create(ctx context.Context, clientset kubernetes.Interface, _ string, obj *rbacv1.ClusterRoleBinding, opts metav1.CreateOptions) (*rbacv1.ClusterRoleBinding, error) {
	return clientset.RbacV1().ClusterRoleBindings().Create(ctx, obj, opts)
}
func (c *ClusterRoleBindingClient) Update(ctx context.Context, clientset kubernetes.Interface, _ string, obj *rbacv1.ClusterRoleBinding, opts metav1.UpdateOptions) (*rbacv1.ClusterRoleBinding, error) {
	return clientset.RbacV1().ClusterRoleBindings().Update(ctx, obj, opts)
}
func (c *ClusterRoleBindingClient) Delete(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.DeleteOptions) error {
	return clientset.RbacV1().ClusterRoleBindings().Delete(ctx, name, opts)
}
func (c *ClusterRoleBindingClient) Watch(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.RbacV1().ClusterRoleBindings().Watch(ctx, opts)
}
