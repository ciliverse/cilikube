package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

// [Core Change] Function name changed to RegisterAuthRoutes and receives *gin.RouterGroup
// This way it stays consistent with other Register... functions
func RegisterAuthRoutes(authGroup *gin.RouterGroup) {
	authHandler := handlers.NewAuthHandler()

	// Routes are registered directly on the passed authGroup, no longer creating our own

	// Public routes (no authentication required)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/register", authHandler.Register)

	// Routes requiring authentication
	authenticated := authGroup.Group("")
	authenticated.Use(auth.JWTAuthMiddleware())
	{
		authenticated.GET("/profile", authHandler.GetProfile)
		authenticated.PUT("/profile", authHandler.UpdateProfile)
		authenticated.POST("/change-password", authHandler.ChangePassword)
		authenticated.POST("/logout", authHandler.Logout)
	}

	// Admin-only routes
	admin := authGroup.Group("/admin") // Grouping admin routes under /admin for clarity
	admin.Use(auth.JWTAuthMiddleware(), auth.AdminRequiredMiddleware())
	{
		admin.GET("/users", authHandler.GetUserList)
		admin.PUT("/users/:id/status", authHandler.UpdateUserStatus)
		admin.DELETE("/users/:id", authHandler.DeleteUser)
	}
}
