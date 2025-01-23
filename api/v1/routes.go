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
	}
}
