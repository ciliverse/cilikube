package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterNetworkPolicyRoutes 注册NetworkPolicy相关路由
func RegisterNetworkPolicyRoutes(router *gin.RouterGroup, handler *handlers.NetworkPolicyHandler) {
	// 基础资源操作
	networkPolicyGroup := router.Group("/namespaces/:namespace/networkpolicies")
	{
		networkPolicyGroup.GET("", handler.ListNetworkPolicies)
		networkPolicyGroup.POST("", handler.CreateNetworkPolicy)
		networkPolicyGroup.GET("/:name", handler.GetNetworkPolicy)
		networkPolicyGroup.PUT("/:name", handler.UpdateNetworkPolicy)
		networkPolicyGroup.DELETE("/:name", handler.DeleteNetworkPolicy)
	}

	// Watch端点
	watchGroup := router.Group("/watch/namespaces/:namespace/networkpolicies")
	{
		watchGroup.GET("", handler.WatchNetworkPolicies)
	}
}
