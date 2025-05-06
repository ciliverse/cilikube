package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRbacRoutes(router *gin.RouterGroup, handler *handlers.RbacHandler) {
	rbacRouter := router.Group("/namespaces/:namespace/rbac")
	// 在 RBAC 中细分成角色和绑定 的路由组
	// Roles
	rolesRouter := rbacRouter.Group("/roles")
	{
		rolesRouter.GET("", handler.ListRoles)
		rolesRouter.GET("/:name", handler.GetRole)
	}

	// RolesBinding
	roleBindingsRouter := rbacRouter.Group("/roleBindings")
	{
		roleBindingsRouter.GET("", handler.ListRoleBindings)
		roleBindingsRouter.GET("/:name", handler.GetRoleBindings)
	}

	// ClusterRoles
	clusterRolesRouter := rbacRouter.Group("/clusterRoles")
	{
		clusterRolesRouter.GET("", handler.ListClusterRoles)
		clusterRolesRouter.GET("/:name", handler.GetClusterRoles)
	}

	// ClusterRolesBinding
	clusterRoleBindingsRouter := rbacRouter.Group("/clusterRoleBindings")
	{
		clusterRoleBindingsRouter.GET("", handler.ListClusterRoleBindings)
		clusterRoleBindingsRouter.GET("/:name", handler.GetClusterRoleBindings)
	}

	// ServiceAccount
	serviceAccountsRouter := rbacRouter.Group("/serviceAccounts")
	{
		serviceAccountsRouter.GET("", handler.ListServiceAccounts)
		serviceAccountsRouter.GET("/:name", handler.GetServiceAccounts)
	}
}
