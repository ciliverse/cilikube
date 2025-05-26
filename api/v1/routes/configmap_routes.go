package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterConfigMapRoutes 注册ConfigMap相关路由
func RegisterConfigMapRoutes(router *gin.RouterGroup, handler *handlers.ConfigMapHandler) {
	// 基础资源操作
	configMapGroup := router.Group("/namespaces/:namespace/configmaps")
	{
		configMapGroup.GET("", handler.ListConfigMaps)
		configMapGroup.POST("", handler.CreateConfigMap)
		configMapGroup.GET("/:name", handler.GetConfigMap)
		configMapGroup.PUT("/:name", handler.UpdateConfigMap)
		configMapGroup.DELETE("/:name", handler.DeleteConfigMap)
	}

	// // Watch端点
	// watchGroup := router.Group("/watch/namespaces/:namespace/configmaps")
	// {
	// 	watchGroup.GET("", handlers.WatchConfigMaps)
	// }
}
