package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupCRDRoutes sets up CRD routes
func SetupCRDRoutes(router *gin.RouterGroup, crdHandler *handlers.CRDHandler) {
	crdGroup := router.Group("/crds")
	{
		// CRD management
		crdGroup.GET("", crdHandler.ListCRDs)                // Get CRD list
		crdGroup.GET("/definition/:name", crdHandler.GetCRD) // Get CRD details

		// Custom resource management
		resourceGroup := crdGroup.Group("/resources/:group/:version/:plural")
		{
			resourceGroup.GET("", crdHandler.ListCustomResources)           // Get custom resource list
			resourceGroup.POST("", crdHandler.CreateCustomResource)         // Create custom resource
			resourceGroup.GET("/:name", crdHandler.GetCustomResource)       // Get custom resource details
			resourceGroup.PUT("/:name", crdHandler.UpdateCustomResource)    // Update custom resource
			resourceGroup.DELETE("/:name", crdHandler.DeleteCustomResource) // Delete custom resource
		}
	}
}
