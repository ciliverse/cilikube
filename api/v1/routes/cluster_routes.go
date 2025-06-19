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

		// 使用集群 ID 进行删除操作
		clusterRoutes.DELETE("/:cluster_id", handler.DeleteCluster)

		activeRoutes := clusterRoutes.Group("/active")
		{
			activeRoutes.GET("", handler.GetActiveCluster)
			activeRoutes.POST("", handler.SetActiveCluster)
		}
	}
}
