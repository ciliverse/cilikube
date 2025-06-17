package service

// import (
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	corev1 "k8s.io/api/core/v1"
// 	rbacv1 "k8s.io/api/rbac/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"
// )

// // 使用 fake 客户端的测试示例
// func TestRbacService_ListRoles(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 Role
// 	testRole := &rbacv1.Role{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "test-role",
// 			Namespace: "default",
// 		},
// 		Rules: []rbacv1.PolicyRule{
// 			{
// 				Verbs:     []string{"get", "list"},
// 				APIGroups: []string{""},
// 				Resources: []string{"pods"},
// 			},
// 		},
// 	}

// 	// 将测试 Role 添加到假客户端
// 	_, err := fakeClient.RbacV1().Roles("default").Create(context.TODO(), testRole, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 ListRoles
// 	roles, err := service.ListRoles("default")
// 	assert.NoError(t, err)
// 	assert.Len(t, roles, 1)
// 	assert.Equal(t, "test-role", roles[0].Name)
// 	assert.Equal(t, "default", roles[0].Namespace)
// }

// // 测试 GetRole 方法
// func TestRbacService_GetRole(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 Role
// 	testRole := &rbacv1.Role{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "test-role",
// 			Namespace: "default",
// 		},
// 		Rules: []rbacv1.PolicyRule{
// 			{
// 				Verbs:     []string{"get", "list"},
// 				APIGroups: []string{""},
// 				Resources: []string{"pods"},
// 			},
// 		},
// 	}

// 	// 将测试 Role 添加到假客户端
// 	_, err := fakeClient.RbacV1().Roles("default").Create(context.TODO(), testRole, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 GetRole
// 	role, err := service.GetRole("default", "test-role")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, role)
// 	assert.Equal(t, "test-role", role.Name)
// }

// // 测试 ListRoleBindings 方法
// func TestRbacService_ListRoleBindings(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 RoleBinding
// 	testRoleBinding := &rbacv1.RoleBinding{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "test-rolebinding",
// 			Namespace: "default",
// 		},
// 		RoleRef: rbacv1.RoleRef{
// 			APIGroup: "rbac.authorization.k8s.io",
// 			Kind:     "Role",
// 			Name:     "test-role",
// 		},
// 		Subjects: []rbacv1.Subject{
// 			{
// 				Kind:      "ServiceAccount",
// 				Name:      "test-sa",
// 				Namespace: "default",
// 			},
// 		},
// 	}

// 	// 将测试 RoleBinding 添加到假客户端
// 	_, err := fakeClient.RbacV1().RoleBindings("default").Create(context.TODO(), testRoleBinding, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 ListRoleBindings
// 	roleBindings, err := service.ListRoleBindings("default")
// 	assert.NoError(t, err)
// 	assert.Len(t, roleBindings, 1)
// 	assert.Equal(t, "test-rolebinding", roleBindings[0].Name)
// 	assert.Equal(t, "default", roleBindings[0].Namespace)
// 	assert.Equal(t, "test-role", roleBindings[0].RoleRef.Name)
// }

// // 测试 GetRoleBinding 方法
// func TestRbacService_GetRoleBinding(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 RoleBinding
// 	testRoleBinding := &rbacv1.RoleBinding{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "test-rolebinding",
// 			Namespace: "default",
// 		},
// 		RoleRef: rbacv1.RoleRef{
// 			APIGroup: "rbac.authorization.k8s.io",
// 			Kind:     "Role",
// 			Name:     "test-role",
// 		},
// 		Subjects: []rbacv1.Subject{
// 			{
// 				Kind:      "ServiceAccount",
// 				Name:      "test-sa",
// 				Namespace: "default",
// 			},
// 		},
// 	}

// 	// 将测试 RoleBinding 添加到假客户端
// 	_, err := fakeClient.RbacV1().RoleBindings("default").Create(context.TODO(), testRoleBinding, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 GetRoleBinding
// 	roleBinding, err := service.GetRoleBinding("default", "test-rolebinding")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, roleBinding)
// 	assert.Equal(t, "test-rolebinding", roleBinding.Name)
// 	assert.Equal(t, "test-role", roleBinding.RoleRef.Name)
// }

// // 测试 ListClusterRoles 方法
// func TestRbacService_ListClusterRoles(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 ClusterRole
// 	testClusterRole := &rbacv1.ClusterRole{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "test-clusterrole",
// 		},
// 		Rules: []rbacv1.PolicyRule{
// 			{
// 				Verbs:     []string{"get", "list"},
// 				APIGroups: []string{""},
// 				Resources: []string{"pods"},
// 			},
// 		},
// 	}

// 	// 将测试 ClusterRole 添加到假客户端
// 	_, err := fakeClient.RbacV1().ClusterRoles().Create(context.TODO(), testClusterRole, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 ListClusterRoles
// 	clusterRoles, err := service.ListClusterRoles()
// 	assert.NoError(t, err)
// 	assert.Len(t, clusterRoles, 1)
// 	assert.Equal(t, "test-clusterrole", clusterRoles[0].Name)
// }

// // 测试 GetClusterRole 方法
// func TestRbacService_GetClusterRole(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 ClusterRole
// 	testClusterRole := &rbacv1.ClusterRole{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "test-clusterrole",
// 		},
// 		Rules: []rbacv1.PolicyRule{
// 			{
// 				Verbs:     []string{"get", "list"},
// 				APIGroups: []string{""},
// 				Resources: []string{"pods"},
// 			},
// 		},
// 	}

// 	// 将测试 ClusterRole 添加到假客户端
// 	_, err := fakeClient.RbacV1().ClusterRoles().Create(context.TODO(), testClusterRole, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 GetClusterRole
// 	clusterRole, err := service.GetClusterRole("test-clusterrole")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, clusterRole)
// 	assert.Equal(t, "test-clusterrole", clusterRole.Name)
// }

// // 测试 ListClusterRoleBindings 方法
// func TestRbacService_ListClusterRoleBindings(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 ClusterRoleBinding
// 	testClusterRoleBinding := &rbacv1.ClusterRoleBinding{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "test-clusterrolebinding",
// 		},
// 		RoleRef: rbacv1.RoleRef{
// 			APIGroup: "rbac.authorization.k8s.io",
// 			Kind:     "ClusterRole",
// 			Name:     "test-clusterrole",
// 		},
// 		Subjects: []rbacv1.Subject{
// 			{
// 				Kind:      "ServiceAccount",
// 				Name:      "test-sa",
// 				Namespace: "default",
// 			},
// 		},
// 	}

// 	// 将测试 ClusterRoleBinding 添加到假客户端
// 	_, err := fakeClient.RbacV1().ClusterRoleBindings().Create(context.TODO(), testClusterRoleBinding, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 ListClusterRoleBindings
// 	clusterRoleBindings, err := service.ListClusterRoleBindings()
// 	assert.NoError(t, err)
// 	assert.Len(t, clusterRoleBindings, 1)
// 	assert.Equal(t, "test-clusterrolebinding", clusterRoleBindings[0].Name)
// 	assert.Equal(t, "test-clusterrole", clusterRoleBindings[0].RoleRef.Name)
// }

// // 测试 GetClusterRoleBinding 方法
// func TestRbacService_GetClusterRoleBinding(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 ClusterRoleBinding
// 	testClusterRoleBinding := &rbacv1.ClusterRoleBinding{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: "test-clusterrolebinding",
// 		},
// 		RoleRef: rbacv1.RoleRef{
// 			APIGroup: "rbac.authorization.k8s.io",
// 			Kind:     "ClusterRole",
// 			Name:     "test-clusterrole",
// 		},
// 		Subjects: []rbacv1.Subject{
// 			{
// 				Kind:      "ServiceAccount",
// 				Name:      "test-sa",
// 				Namespace: "default",
// 			},
// 		},
// 	}

// 	// 将测试 ClusterRoleBinding 添加到假客户端
// 	_, err := fakeClient.RbacV1().ClusterRoleBindings().Create(context.TODO(), testClusterRoleBinding, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 GetClusterRoleBinding
// 	clusterRoleBinding, err := service.GetClusterRoleBinding("test-clusterrolebinding")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, clusterRoleBinding)
// 	assert.Equal(t, "test-clusterrolebinding", clusterRoleBinding.Name)
// 	assert.Equal(t, "test-clusterrole", clusterRoleBinding.RoleRef.Name)
// }

// // 测试 ListServiceAccounts 方法
// func TestRbacService_ListServiceAccounts(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 ServiceAccount
// 	testServiceAccount := &corev1.ServiceAccount{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "test-sa",
// 			Namespace: "default",
// 		},
// 	}

// 	// 将测试 ServiceAccount 添加到假客户端
// 	_, err := fakeClient.CoreV1().ServiceAccounts("default").Create(context.TODO(), testServiceAccount, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 ListServiceAccounts
// 	serviceAccounts, err := service.ListServiceAccounts("default")
// 	assert.NoError(t, err)
// 	assert.Len(t, serviceAccounts, 1)
// 	assert.Equal(t, "test-sa", serviceAccounts[0].Name)
// 	assert.Equal(t, "default", serviceAccounts[0].Namespace)
// }

// // 测试 GetServiceAccounts 方法
// func TestRbacService_GetServiceAccounts(t *testing.T) {
// 	// 创建一个假的 Kubernetes 客户端
// 	fakeClient := fake.NewSimpleClientset()

// 	// 创建一个测试用的 ServiceAccount
// 	testServiceAccount := &corev1.ServiceAccount{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "test-sa",
// 			Namespace: "default",
// 		},
// 	}

// 	// 将测试 ServiceAccount 添加到假客户端
// 	_, err := fakeClient.CoreV1().ServiceAccounts("default").Create(context.TODO(), testServiceAccount, metav1.CreateOptions{})
// 	assert.NoError(t, err)

// 	// 创建服务
// 	service := NewRbacService(fakeClient)

// 	// 测试 GetServiceAccounts
// 	serviceAccount, err := service.GetServiceAccounts("default", "test-sa")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, serviceAccount)
// 	assert.Equal(t, "test-sa", serviceAccount.Name)
// 	assert.Equal(t, "default", serviceAccount.Namespace)
// }
