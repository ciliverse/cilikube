package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterDeploymentRoutes 注册Deployment相关路由
func RegisterDeploymentRoutes(router *gin.RouterGroup, handler *handlers.DeploymentHandler) {
	// 基础资源操作
	deploymentGroup := router.Group("/namespaces/:namespace/deployments")
	{
		deploymentGroup.GET("", handler.ListDeployments)
		deploymentGroup.POST("", handler.CreateDeployment)
		deploymentGroup.GET("/:name", handler.GetDeployment)
		deploymentGroup.PUT("/:name", handler.UpdateDeployment)
		deploymentGroup.DELETE("/:name", handler.DeleteDeployment)
		deploymentGroup.PUT("/:name/scale", handler.ScaleDeployment)
		deploymentGroup.GET("/:name/pods", handler.GetDeploymentPods)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces/:namespace/deployments")
	{
		watchGroup.GET("", handler.WatchDeployments)
	}
}
