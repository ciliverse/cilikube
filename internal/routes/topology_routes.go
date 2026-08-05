package routes

import (
	"github.com/ciliverse/cilikube/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterTopologyRoutes(router *gin.RouterGroup, handler *handlers.TopologyHandler) {
	g := router.Group("/topology")
	{
		g.GET("", handler.GetGraph)
		g.GET("/traffic", handler.GetTraffic)
	}
}
