package handlers

import (
	"fmt"
	"net/http"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	service *service.ClusterService
}

func NewClusterHandler(svc *service.ClusterService) *ClusterHandler {
	return &ClusterHandler{service: svc}
}

// ListClusters gets cluster list
func (h *ClusterHandler) ListClusters(c *gin.Context) {
	clusters := h.service.ListClusters()
	utils.ApiSuccess(c, clusters, "successfully retrieved cluster list")
}

// GetCluster gets single cluster details
func (h *ClusterHandler) GetCluster(c *gin.Context) {
	clusterID := c.Param("id")
	cluster, err := h.service.GetClusterByID(clusterID)
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, "failed to get cluster", err.Error())
		return
	}
	utils.ApiSuccess(c, cluster, "successfully retrieved cluster details")
}

// CreateCluster creates a new cluster
func (h *ClusterHandler) CreateCluster(c *gin.Context) {
	var req models.CreateClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "request parameter error", err.Error())
		return
	}
	if err := h.service.CreateCluster(req); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to create cluster", err.Error())
		return
	}
	utils.ApiSuccess(c, nil, "cluster created successfully")
}

// UpdateCluster updates an existing cluster
func (h *ClusterHandler) UpdateCluster(c *gin.Context) {
	clusterID := c.Param("id")
	var req models.UpdateClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "request parameter error", err.Error())
		return
	}
	if err := h.service.UpdateCluster(clusterID, req); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to update cluster", err.Error())
		return
	}
	utils.ApiSuccess(c, nil, "cluster updated successfully")
}

// DeleteCluster deletes a cluster
func (h *ClusterHandler) DeleteCluster(c *gin.Context) {
	clusterID := c.Param("id")
	if err := h.service.DeleteClusterByID(clusterID); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to delete cluster", err.Error())
		return
	}
	utils.ApiSuccess(c, nil, "cluster deleted successfully")
}

// SetActiveCluster sets the current active cluster
func (h *ClusterHandler) SetActiveCluster(c *gin.Context) {
	var req struct {
		ID   string `json:"id"`
		Name string `json:"name"` // Maintain backward compatibility
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "request parameter error", err.Error())
		return
	}

	var targetID string
	if req.ID != "" {
		// Prioritize using ID
		targetID = req.ID
	} else if req.Name != "" {
		// Backward compatibility: find cluster ID by name
		clusters := h.service.ListClusters()
		for _, cluster := range clusters {
			if cluster.Name == req.Name {
				targetID = cluster.ID
				break
			}
		}
		if targetID == "" {
			utils.ApiError(c, http.StatusNotFound, "cluster does not exist", fmt.Sprintf("cluster named '%s' not found", req.Name))
			return
		}
	} else {
		utils.ApiError(c, http.StatusBadRequest, "request parameter error", "must provide id or name parameter")
		return
	}

	if err := h.service.SetActiveCluster(targetID); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to switch active cluster", err.Error())
		return
	}

	// Return detailed cluster information
	activeCluster, err := h.service.GetClusterByID(targetID)
	if err != nil {
		utils.ApiSuccess(c, gin.H{"activeClusterID": targetID}, "active cluster switched successfully")
	} else {
		utils.ApiSuccess(c, gin.H{
			"activeClusterID":   targetID,
			"activeClusterName": activeCluster.Name,
			"cluster":           activeCluster,
		}, "active cluster switched successfully")
	}
}

// GetActiveCluster gets the current active cluster
func (h *ClusterHandler) GetActiveCluster(c *gin.Context) {
	activeClusterID := h.service.GetActiveClusterID()
	if activeClusterID == "" {
		utils.ApiError(c, http.StatusNotFound, "no active cluster currently", "please add and activate a cluster first")
		return
	}

	// Return the active cluster name directly, if cluster details are needed, they can be obtained through other APIs
	utils.ApiSuccess(c, activeClusterID, "successfully retrieved active cluster")
}
