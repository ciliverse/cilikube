package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterNamespaceRoutes 注册Namespace相关路由
func RegisterNamespaceRoutes(router *gin.RouterGroup, handler *handlers.NamespaceHandler) {
	// 基础资源操作
	namespaceGroup := router.Group("/namespace")
	{
		namespaceGroup.GET("", handler.ListNamespaces)
		namespaceGroup.POST("", handler.CreateNamespace)
		namespaceGroup.GET("/:name", handler.GetNamespace)
		namespaceGroup.PUT("/:name", handler.UpdateNamespace)
		namespaceGroup.DELETE("/:name", handler.DeleteNamespace)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces")
	{
		watchGroup.GET("", handler.WatchNamespaces)
	}
}
