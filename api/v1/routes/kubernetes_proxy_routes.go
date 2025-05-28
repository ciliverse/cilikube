package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/ciliverse/cilikube/api/v1/handlers"
)

func KubernetesProxyRoutes(router *gin.RouterGroup, handler *handlers.ProxyHandler) {
	proxyGroup := router.Group("/proxy")
	{
		proxyGroup.Any("/*act", handler.Proxy)
	}
}
