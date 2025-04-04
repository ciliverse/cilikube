package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterSecretRoutes 注册Secret相关路由
func RegisterSecretRoutes(router *gin.RouterGroup, handler *handlers.SecretHandler) {
	// 基础资源操作
	secretGroup := router.Group("/namespaces/:namespace/secrets")
	{
		secretGroup.GET("", handler.ListSecrets)
		secretGroup.POST("", handler.CreateSecret)
		secretGroup.GET("/:name", handler.GetSecret)
		secretGroup.PUT("/:name", handler.UpdateSecret)
		secretGroup.DELETE("/:name", handler.DeleteSecret)
	}

	// // Watch端点
	// watchGroup := router.Group("/watch/namespaces/:namespace/secrets")
	// {
	// 	watchGroup.GET("", handler.WatchSecrets)
	// }
}
