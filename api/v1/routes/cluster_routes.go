package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterClusterRoutes(router *gin.RouterGroup, handler *handlers.ClusterHandler) {
	// 这个路由组现在只负责集群本身的元数据管理
	clusterRoutes := router.Group("/clusters")
	{
		clusterRoutes.GET("", handler.ListClusters)
		clusterRoutes.POST("", handler.CreateCluster)
		clusterRoutes.GET("/:id", handler.GetCluster)
		clusterRoutes.PUT("/:id", handler.UpdateCluster)
		clusterRoutes.DELETE("/:id", handler.DeleteCluster)

		// 激活集群的API
		activeRoutes := clusterRoutes.Group("/active")
		{
			// activeRoutes.GET("", handler.GetActiveCluster)
			activeRoutes.POST("", handler.SetActiveCluster)
		}
	}
}
