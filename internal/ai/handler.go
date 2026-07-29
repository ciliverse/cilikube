package ai

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
)

type Handler struct {
	cfg     *configs.Config
	manager *k8s.ClusterManager
}

func NewHandler(cfg *configs.Config, manager *k8s.ClusterManager) *Handler {
	return &Handler{cfg: cfg, manager: manager}
}

func (h *Handler) Status(c *gin.Context) {
	cfg := h.cfg
	if cfg == nil {
		cfg = configs.GlobalConfig
	}
	provider := "mock"
	if cfg != nil && cfg.AI.Provider != "" {
		provider = cfg.AI.Provider
	}
	utils.ApiSuccess(c, gin.H{
		"enabled":  cfg != nil && cfg.AI.Enabled,
		"ready":    Ready(cfg),
		"provider": provider,
		"model":    cfg.AI.Model,
	}, "ok")
}

func (h *Handler) Chat(c *gin.Context) {
	cfg := h.cfg
	if cfg == nil {
		cfg = configs.GlobalConfig
	}
	if !Ready(cfg) {
		utils.ApiError(c, http.StatusBadRequest, "AI is not enabled or not configured")
		return
	}
	client, ok := k8s.GetClientFromQuery(c, h.manager)
	if !ok {
		return
	}
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}
	if len(req.Messages) > 40 {
		utils.ApiError(c, http.StatusBadRequest, "Too many messages")
		return
	}
	for i := range req.Messages {
		if len(req.Messages[i].Content) > 8000 {
			req.Messages[i].Content = req.Messages[i].Content[:8000]
		}
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		utils.ApiError(c, http.StatusInternalServerError, "Streaming unsupported")
		return
	}

	send := func(ev SSEEvent) {
		_, _ = fmt.Fprint(c.Writer, MarshalSSE(ev))
		flusher.Flush()
	}

	ProcessChat(c.Request.Context(), cfg, client, &req, send)
}
