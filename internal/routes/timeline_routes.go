package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterTimelineRoutes(router *gin.RouterGroup, handler *handlers.TimelineHandler) {
	router.GET("/timeline", handler.GetTimeline)
	router.GET("/timeline/meta", handler.GetMeta)
}
