package handlers

import (
	"net/http"
	"strings"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TopologyHandler struct {
	service        *service.TopologyService
	clusterManager *k8s.ClusterManager
}

func NewTopologyHandler(svc *service.TopologyService, cm *k8s.ClusterManager) *TopologyHandler {
	return &TopologyHandler{service: svc, clusterManager: cm}
}

// GetGraph godoc
// @Summary Cluster resource topology graph
// @Tags Topology
// @Router /api/v1/topology [get]
func (h *TopologyHandler) GetGraph(c *gin.Context) {
	client, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	ns := strings.TrimSpace(c.Query("namespace"))
	groupBy := strings.TrimSpace(c.Query("groupBy"))
	kindsRaw := strings.TrimSpace(c.Query("kinds"))
	kinds := map[string]bool{}
	if kindsRaw != "" {
		for _, k := range strings.Split(kindsRaw, ",") {
			k = strings.ToLower(strings.TrimSpace(k))
			if k != "" {
				kinds[k] = true
			}
		}
	}

	graph, err := h.service.BuildGraph(c.Request.Context(), client.Clientset, service.TopologyOptsForHandler(ns, groupBy, kinds))
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, err.Error(), "")
		return
	}
	utils.ApiSuccess(c, graph, "ok")
}

// GetTraffic godoc
// @Summary Topology traffic edges
// @Tags Topology
// @Router /api/v1/topology/traffic [get]
func (h *TopologyHandler) GetTraffic(c *gin.Context) {
	client, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	ns := strings.TrimSpace(c.Query("namespace"))
	traffic, err := h.service.BuildTraffic(c.Request.Context(), client.Clientset, ns)
	if err != nil {
		utils.ApiError(c, http.StatusBadRequest, err.Error(), "")
		return
	}
	utils.ApiSuccess(c, traffic, "ok")
}
