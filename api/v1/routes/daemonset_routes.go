package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterDaemonSetRoutes 注册DaemonSet相关路由
func RegisterDaemonSetRoutes(router *gin.RouterGroup, handler *handlers.DaemonSetHandler) {
	// 基础资源操作
	daemonSetGroup := router.Group("/namespaces/:namespace/daemonsets")
	{
		daemonSetGroup.GET("", handler.ListDaemonSets)
		daemonSetGroup.POST("", handler.CreateDaemonSet)
		daemonSetGroup.GET("/:name", handler.GetDaemonSet)
		daemonSetGroup.PUT("/:name", handler.UpdateDaemonSet)
		daemonSetGroup.DELETE("/:name", handler.DeleteDaemonSet)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces/:namespace/daemonsets")
	{
		watchGroup.GET("", handler.WatchDaemonSets)
	}
}
