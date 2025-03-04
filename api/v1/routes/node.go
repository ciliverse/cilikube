package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterNodeRoutes 注册Node相关路由
func RegisterNodeRoutes(router *gin.RouterGroup, handler *handlers.NodeHandler) {
	// 基础资源操作
	nodeGroup := router.Group("/nodes")
	{
		nodeGroup.GET("", handler.ListNodes)
		nodeGroup.POST("", handler.CreateNode)
		nodeGroup.GET("/:name", handler.GetNode)
		nodeGroup.PUT("/:name", handler.UpdateNode)
		nodeGroup.DELETE("/:name", handler.DeleteNode)
	}

	// Watch端点
	watchGroup := router.Group("/watch/nodes")
	{
		watchGroup.GET("", handler.WatchNodes)
	}
}
