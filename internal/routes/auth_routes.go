package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication and OAuth routes
func RegisterAuthRoutes(authGroup *gin.RouterGroup, authService *service.AuthService, oauthService *service.OAuthService) {
	authHandler := handlers.NewAuthHandler(authService)
	oauthHandler := handlers.NewOAuthHandler(oauthService)

	// Routes are registered directly on the passed authGroup, no longer creating our own

	// Public routes (no authentication required)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/register", authHandler.Register)

	// OAuth routes (public)
	oauth := authGroup.Group("/oauth")
	{
		oauth.GET("/:provider/auth", oauthHandler.GetAuthURL)
		oauth.POST("/callback", oauthHandler.HandleCallback)
	}

	// Routes requiring authentication
	authenticated := authGroup.Group("")
	authenticated.Use(auth.JWTAuthMiddleware())
	{
		authenticated.GET("/profile", authHandler.GetProfile)
		authenticated.GET("/profile/detailed", authHandler.GetDetailedProfile)
		authenticated.PUT("/profile", authHandler.UpdateProfile)
		authenticated.POST("/change-password", authHandler.ChangePassword)
		authenticated.POST("/refresh", authHandler.RefreshToken)
		authenticated.POST("/logout", authHandler.Logout)

		// OAuth account management (authenticated)
		authenticated.POST("/oauth/link", oauthHandler.LinkAccount)
		authenticated.POST("/oauth/unlink", oauthHandler.UnlinkAccount)
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
