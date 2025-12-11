package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

// RegisterRoleManagementRoutes registers role management routes for administrators
func RegisterRoleManagementRoutes(router *gin.RouterGroup, roleService *service.RoleService) {
	roleHandler := handlers.NewRoleManagementHandler(roleService)

	// Apply JWT middleware and admin permission check to all role management routes
	roleRoutes := router.Group("/roles")
	roleRoutes.Use(auth.JWTAuthMiddleware())
	// TODO: Add admin permission middleware here
	{
		// Role CRUD operations
		roleRoutes.GET("", roleHandler.ListRoles)
		roleRoutes.POST("", roleHandler.CreateRole)
		roleRoutes.GET("/:id", roleHandler.GetRole)
		roleRoutes.PUT("/:id", roleHandler.UpdateRole)
		roleRoutes.DELETE("/:id", roleHandler.DeleteRole)

		// Role permission operations
		roleRoutes.GET("/:id/permissions", roleHandler.GetRolePermissions)
		roleRoutes.PUT("/:id/permissions", roleHandler.UpdateRolePermissions)

		// Role assignment operations
		roleRoutes.POST("/assign", roleHandler.AssignRoleToUser)
		roleRoutes.POST("/remove", roleHandler.RemoveRoleFromUser)

		// Role-user relationship queries
		roleRoutes.GET("/users/:userId", roleHandler.GetUserRoles)
		roleRoutes.GET("/:id/users", roleHandler.GetRoleUsers)
	}

	// Permission management routes
	permissionRoutes := router.Group("/permissions")
	permissionRoutes.Use(auth.JWTAuthMiddleware())
	// TODO: Add admin permission middleware here
	{
		permissionRoutes.GET("", roleHandler.GetAvailablePermissions)
	}
}
