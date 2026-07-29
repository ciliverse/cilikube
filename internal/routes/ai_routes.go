package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/ai"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

func RegisterAIRoutes(router *gin.RouterGroup, cfg *configs.Config, manager *k8s.ClusterManager) {
	h := ai.NewHandler(cfg, manager)
	g := router.Group("/ai")
	{
		g.GET("/status", h.Status)
		g.POST("/chat", h.Chat)
	}
}
