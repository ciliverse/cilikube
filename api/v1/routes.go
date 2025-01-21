package v1

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/pods", handlers.ListPods)
		api.POST("/pods", handlers.CreatePod)
		api.GET("/pods/:name", handlers.GetPod)
		api.PUT("/pods/:name", handlers.UpdatePod)
		api.DELETE("/pods/:name", handlers.DeletePod)

		api.GET("/deployments", handlers.ListDeployments)
		api.POST("/deployments", handlers.CreateDeployment)
		api.GET("/deployments/:name", handlers.GetDeployment)
		api.PUT("/deployments/:name", handlers.UpdateDeployment)
		api.DELETE("/deployments/:name", handlers.DeleteDeployment)

		api.GET("/services", handlers.ListServices)
		api.POST("/services", handlers.CreateService)
		api.GET("/services/:name", handlers.GetService)
		api.PUT("/services/:name", handlers.UpdateService)
		api.DELETE("/services/:name", handlers.DeleteService)
	}
}
