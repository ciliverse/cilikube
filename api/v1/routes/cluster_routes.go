package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterClusterRoutes(router *gin.RouterGroup, handler *handlers.ClusterHandler) {
	// This route group is now only responsible for cluster metadata management
	clusterRoutes := router.Group("/clusters")
	{
		clusterRoutes.GET("", handler.ListClusters)
		clusterRoutes.POST("", handler.CreateCluster)
		clusterRoutes.GET("/:id", handler.GetCluster)
		clusterRoutes.PUT("/:id", handler.UpdateCluster)
		clusterRoutes.DELETE("/:id", handler.DeleteCluster)

		// Active cluster API
		activeRoutes := clusterRoutes.Group("/active")
		{
			activeRoutes.GET("", handler.GetActiveCluster)
			activeRoutes.POST("", handler.SetActiveCluster)
		}
	}
}
