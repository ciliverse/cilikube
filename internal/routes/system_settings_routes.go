package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/gin-gonic/gin"
)

// RegisterSystemSettingsRoutes registers system settings routes for administrators
func RegisterSystemSettingsRoutes(router *gin.RouterGroup) {
	settingsHandler := handlers.NewSystemSettingsHandler()

	// Apply JWT middleware and admin permission check to all system settings routes
	settingsRoutes := router.Group("/settings")
	settingsRoutes.Use(auth.JWTAuthMiddleware())
	// TODO: Add admin permission middleware here
	{
		// System information
		settingsRoutes.GET("/system", settingsHandler.GetSystemInfo)

		// OAuth settings
		settingsRoutes.GET("/oauth", settingsHandler.GetOAuthSettings)
		settingsRoutes.PUT("/oauth", settingsHandler.UpdateOAuthSettings)

		// Security settings
		settingsRoutes.GET("/security", settingsHandler.GetSecuritySettings)
		settingsRoutes.PUT("/security", settingsHandler.UpdateSecuritySettings)

		// System preferences
		settingsRoutes.GET("/preferences", settingsHandler.GetSystemPreferences)
		settingsRoutes.PUT("/preferences", settingsHandler.UpdateSystemPreferences)
	}
}
