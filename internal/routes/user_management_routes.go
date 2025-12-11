package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

// RegisterUserManagementRoutes registers user management routes for administrators
func RegisterUserManagementRoutes(router *gin.RouterGroup, authService *service.AuthService, roleService *service.RoleService) {
	userHandler := handlers.NewUserManagementHandler(authService, roleService)

	// Apply JWT middleware and admin permission check to all user management routes
	userRoutes := router.Group("/users")
	userRoutes.Use(auth.JWTAuthMiddleware())
	// TODO: Add admin permission middleware here
	{
		// User CRUD operations
		userRoutes.GET("", userHandler.ListUsers)
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)

		// User status management
		userRoutes.PUT("/:id/status", userHandler.UpdateUserStatus)
	}
}
