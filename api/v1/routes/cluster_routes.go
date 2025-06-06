package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterClusterRoutes(router *gin.RouterGroup, handler *handlers.ClusterHandler) {
	clusterRoutes := router.Group("/clusters")
	{
		clusterRoutes.GET("", handler.ListClusters)
		clusterRoutes.POST("", handler.CreateCluster)

		// [修复] 将 /:name 修改为 /:cluster_name 以解决路由冲突
		clusterRoutes.DELETE("/:cluster_name", handler.DeleteCluster)

		activeRoutes := clusterRoutes.Group("/active")
		{
			activeRoutes.GET("", handler.GetActiveCluster)
			activeRoutes.POST("", handler.SetActiveCluster)
		}
	}
}
