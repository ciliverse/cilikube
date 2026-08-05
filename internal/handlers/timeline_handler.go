package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	service        *service.TimelineService
	clusterManager *k8s.ClusterManager
}

func NewTimelineHandler(svc *service.TimelineService, cm *k8s.ClusterManager) *TimelineHandler {
	return &TimelineHandler{service: svc, clusterManager: cm}
}

// GetTimeline godoc
// @Summary Resource timeline (status segments + events)
// @Tags Timeline
// @Router /api/v1/timeline [get]
func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	client, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	clusterID := strings.TrimSpace(c.Query("clusterId"))
	if clusterID == "" && h.clusterManager != nil {
		clusterID = h.clusterManager.GetActiveClusterID()
	}

	to := time.Now().UTC()
	from := to.Add(-15 * time.Minute)
	if v := strings.TrimSpace(c.Query("to")); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			to = t.UTC()
		}
	}
	if v := strings.TrimSpace(c.Query("from")); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			from = t.UTC()
		}
	}
	if window := strings.TrimSpace(c.Query("window")); window != "" {
		to = time.Now().UTC()
		switch window {
		case "1h":
			from = to.Add(-time.Hour)
		case "6h":
			from = to.Add(-6 * time.Hour)
		default:
			from = to.Add(-15 * time.Minute)
		}
	}

	kinds := map[string]bool{}
	if raw := strings.TrimSpace(c.Query("kinds")); raw != "" {
		for _, k := range strings.Split(raw, ",") {
			k = strings.ToLower(strings.TrimSpace(k))
			if k != "" {
				kinds[k] = true
			}
		}
	}

	resp, err := h.service.Build(c.Request.Context(), client.Clientset, service.TimelineQuery{
		ClusterID: clusterID,
		Namespace: strings.TrimSpace(c.Query("namespace")),
		From:      from,
		To:        to,
		GroupBy:   strings.TrimSpace(c.Query("groupBy")),
		Q:         strings.TrimSpace(c.Query("q")),
		Kinds:     kinds,
	})
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, err.Error(), "")
		return
	}
	utils.ApiSuccess(c, resp, "ok")
}

// GetMeta godoc
// @Summary Timeline sampler metadata
// @Tags Timeline
// @Router /api/v1/timeline/meta [get]
func (h *TimelineHandler) GetMeta(c *gin.Context) {
	clusterID := strings.TrimSpace(c.Query("clusterId"))
	if clusterID == "" && h.clusterManager != nil {
		clusterID = h.clusterManager.GetActiveClusterID()
	}
	utils.ApiSuccess(c, h.service.Meta(clusterID), "ok")
}
