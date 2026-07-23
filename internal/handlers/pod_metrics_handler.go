package handlers

import (
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

// PodMetricsHandler serves pod CPU/MEM snapshots from metrics-server.
type PodMetricsHandler struct {
	service        *service.PodMetricsService
	clusterManager *k8s.ClusterManager
}

func NewPodMetricsHandler(svc *service.PodMetricsService, k8sManager *k8s.ClusterManager) *PodMetricsHandler {
	return &PodMetricsHandler{service: svc, clusterManager: k8sManager}
}

// ListPodMetrics GET /api/v1/pods/metrics?namespace=&clusterId=
func (h *PodMetricsHandler) ListPodMetrics(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	namespace := c.Query("namespace")
	data := h.service.ListPodMetrics(k8sClient.Config, namespace)
	// Always 200 — Available flag tells the UI whether metrics-server works.
	utils.ApiSuccess(c, data, "pod metrics retrieved")
}
