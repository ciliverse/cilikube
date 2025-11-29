package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers" // Import your handlers
	"github.com/gin-gonic/gin"
)

// RegisterInstallerRoutes registers routes related to the Minikube installer.
func RegisterInstallerRoutes(router *gin.RouterGroup, installerHandler *handlers.InstallerHandler) {
	// Health check endpoint
	router.GET("/healthz", installerHandler.HealthCheck)

	installerRoutes := router.Group("/system") // Group under /system or choose another name
	{
		installerRoutes.GET("/install-minikube", installerHandler.StreamMinikubeInstallation)
	}
}
