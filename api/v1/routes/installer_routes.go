package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers" // Import your handlers
	"github.com/gin-gonic/gin"
)

// RegisterInstallerRoutes registers routes related to the Minikube installer.
// RegisterInstallerRoutes 注册与 Minikube 安装程序相关的路由。
func RegisterInstallerRoutes(router *gin.RouterGroup, installerHandler *handlers.InstallerHandler) {
	// Health check endpoint
	router.GET("/healthz", installerHandler.HealthCheck)
	
	installerRoutes := router.Group("/system") // Group under /system or choose another name // 分组到 /system 下或选择其他名称
	{
		installerRoutes.GET("/install-minikube", installerHandler.StreamMinikubeInstallation)
	}
}