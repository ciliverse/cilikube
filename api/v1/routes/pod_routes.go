package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPodRoutes(router *gin.RouterGroup, handler *handlers.PodHandler) {
	// Namespace listing (keep outside the namespaced group)
	router.GET("/namespaces", handler.ListNamespaces)

	// Operations within a specific namespace
	namespaceGroup := router.Group("/namespaces/:namespace")
	{
		// Pod CRUD and Listing
		podGroup := namespaceGroup.Group("/pods")
		{
			podGroup.GET("", handler.ListPods)   // List Pods
			podGroup.POST("", handler.CreatePod) // Create Pod (JSON or YAML)

			// Pod specific operations
			podNameGroup := podGroup.Group("/:name")
			{
				podNameGroup.GET("", handler.GetPod)       // Get Pod details
				podNameGroup.PUT("", handler.UpdatePod)    // Update Pod (JSON or YAML) - Prefer YAML or PATCH
				podNameGroup.DELETE("", handler.DeletePod) // Delete Pod

				// --- New Endpoints ---
				podNameGroup.GET("/logs", handler.GetPodLogs)    // Get Pod Logs
				podNameGroup.GET("/exec", handler.ExecIntoPod)   // Execute command in Pod (WebSocket)
				podNameGroup.GET("/yaml", handler.GetPodYAML)    // Get Pod as YAML
				podNameGroup.PUT("/yaml", handler.UpdatePodYAML) // Update Pod from YAML
			}
		}

		// Watch endpoints within a namespace
		watchGroup := namespaceGroup.Group("/watch/pods")
		{
			watchGroup.GET("", handler.WatchPods) // Watch Pod changes (SSE)
		}
	}

	// Note: Watch endpoint is now also under /namespaces/:namespace/watch/pods
	// The old /watch/namespaces/:namespace/pods route can be removed or kept for compatibility
	// Let's keep the namespaced structure consistent.

	// Example: Remove the old top-level watch route if desired
	// router.DELETE("/watch/namespaces/:namespace/pods") // Or handle redirection if needed
}
