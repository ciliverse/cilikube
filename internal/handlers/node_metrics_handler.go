package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

// NodeMetricsHandler handles node metrics related requests
type NodeMetricsHandler struct {
	service        *service.NodeMetricsService
	clusterManager *k8s.ClusterManager
}

// NewNodeMetricsHandler creates a new NodeMetricsHandler instance
func NewNodeMetricsHandler(svc *service.NodeMetricsService, k8sManager *k8s.ClusterManager) *NodeMetricsHandler {
	return &NodeMetricsHandler{
		service:        svc,
		clusterManager: k8sManager,
	}
}

// GetNodeMetrics is the HTTP handler function for getting real-time metrics of a single node
func (h *NodeMetricsHandler) GetNodeMetrics(c *gin.Context) {
	// 1. Get clusterId from query parameters and get the corresponding cluster's k8s client
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return // Error already handled in GetClientFromQuery
	}

	// 2. Get node name from path parameters
	nodeName := c.Param("name")
	if nodeName == "" {
		utils.ApiError(c, http.StatusBadRequest, "node name cannot be empty", "")
		return
	}

	// 3. Call service layer to get metrics, note that k8sClient.Config needs to be passed
	metrics, err := h.service.GetNodeMetrics(k8sClient.Config, nodeName)
	if err != nil {
		// Judge the error here, if it's caused by metrics-server not being installed, give a friendly prompt
		if clientErr, ok := err.(interface{ IsNotFound() bool }); ok && clientErr.IsNotFound() {
			utils.ApiError(c, http.StatusNotFound, "failed to get metrics", "Please confirm that Metrics-Server is properly installed and running in the target cluster.")
			return
		}
		utils.ApiError(c, http.StatusInternalServerError, "failed to get node metrics", err.Error())
		return
	}

	// 4. Successfully return data
	utils.ApiSuccess(c, metrics, "successfully retrieved node metrics")
}
