package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/ai"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

func RegisterAIRoutes(router *gin.RouterGroup, cfg *configs.Config, manager *k8s.ClusterManager, prom *service.PrometheusService) {
	h := ai.NewHandler(cfg, manager, prom)
	g := router.Group("/ai")
	{
		g.GET("/status", h.Status)
		g.POST("/chat", h.Chat)
		g.GET("/mcp", h.MCP)
		g.POST("/mcp", h.MCP)
	}
}
