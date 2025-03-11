package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterIngressRoutes 注册Ingress相关路由
func RegisterIngressRoutes(router *gin.RouterGroup, handler *handlers.IngressHandler) {
	// 基础资源操作
	ingressGroup := router.Group("/namespaces/:namespace/ingresses")
	{
		ingressGroup.GET("", handler.ListIngresses)
		ingressGroup.POST("", handler.CreateIngress)
		ingressGroup.GET("/:name", handler.GetIngress)
		ingressGroup.PUT("/:name", handler.UpdateIngress)
		ingressGroup.DELETE("/:name", handler.DeleteIngress)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces/:namespace/ingresses")
	{
		watchGroup.GET("", handler.WatchIngresses)
	}
}
