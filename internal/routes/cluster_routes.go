package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterClusterRoutes(router *gin.RouterGroup, handler *handlers.ClusterHandler) {
	// This route group is now only responsible for cluster metadata management
	clusterRoutes := router.Group("/clusters")
	{
		clusterRoutes.GET("", handler.ListClusters)
		clusterRoutes.POST("", handler.CreateCluster)

		// Static paths must be registered before /:id
		clusterRoutes.GET("/local-contexts", handler.ListLocalKubeContexts)
		clusterRoutes.POST("/import-local", handler.ImportLocalClusters)

		activeRoutes := clusterRoutes.Group("/active")
		{
			activeRoutes.GET("", handler.GetActiveCluster)
			activeRoutes.POST("", handler.SetActiveCluster)
		}

		clusterRoutes.GET("/:id", handler.GetCluster)
		clusterRoutes.PUT("/:id", handler.UpdateCluster)
		clusterRoutes.DELETE("/:id", handler.DeleteCluster)
	}
}
