package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterStatefulSetRoutes 注册StatefulSet相关路由
func RegisterStatefulSetRoutes(router *gin.RouterGroup, handler *handlers.StatefulSetHandler) {
	// 基础资源操作
	statefulSetGroup := router.Group("/namespaces/:namespace/statefulsets")
	{
		statefulSetGroup.GET("", handler.ListStatefulSets)
		statefulSetGroup.POST("", handler.CreateStatefulSet)
		statefulSetGroup.GET("/:name", handler.GetStatefulSet)
		statefulSetGroup.PUT("/:name", handler.UpdateStatefulSet)
		statefulSetGroup.DELETE("/:name", handler.DeleteStatefulSet)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces/:namespace/statefulsets")
	{
		watchGroup.GET("", handler.WatchStatefulSets)
	}
}
