package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterEventsRoutes(router *gin.RouterGroup, handler *handlers.EventsHandler) {

	eventsGroup := router.Group("/namespaces/:namespace/events")
	{
		eventsGroup.GET("", handler.ListEventsHandler)
		eventsGroup.GET("/:name", handler.GetEventsHandler)
	}
}
