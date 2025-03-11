package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterServiceRoutes 注册Service相关路由
func RegisterServiceRoutes(router *gin.RouterGroup, handler *handlers.ServiceHandler) {
	// 基础资源操作
	serviceGroup := router.Group("/namespaces/:namespace/services")
	{
		serviceGroup.GET("", handler.ListServices)
		serviceGroup.POST("", handler.CreateService)
		serviceGroup.GET("/:name", handler.GetService)
		serviceGroup.PUT("/:name", handler.UpdateService)
		serviceGroup.DELETE("/:name", handler.DeleteService)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces/:namespace/services")
	{
		watchGroup.GET("", handler.WatchServices)
	}
}
