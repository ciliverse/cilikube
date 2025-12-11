package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(router *gin.RouterGroup, handler *handlers.EventHandler) {
	eventRoutes := router.Group("/events")
	{
		// List all events with optional filters
		eventRoutes.GET("", handler.ListEvents)

		// Get recent events (for dashboard)
		eventRoutes.GET("/recent", handler.GetRecentEvents)

		// Get events related to a specific object
		eventRoutes.GET("/object/:kind/:name", handler.GetEventsByObject)
	}
}
