package k8s

import (
	"fmt"
	"net/http"

	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

// GetClientFromQuery gets clusterId from URL query parameters and returns the corresponding k8s client.
// This is the "gatekeeper" for all resource operation handler functions.
func GetClientFromQuery(c *gin.Context, cm *ClusterManager) (*Client, bool) {
	clusterID := c.Query("clusterId")
	if clusterID == "" {
		// If no clusterId is provided, try to use the currently active cluster as fallback
		activeID := cm.GetActiveClusterID()
		if activeID == "" {
			utils.ApiError(c, http.StatusBadRequest, "missing 'clusterId' query parameter and no active default cluster", "e.g., /api/v1/nodes?clusterId=cls-xxxxx")
			return nil, false
		}
		clusterID = activeID
	}

	client, err := cm.GetClientByID(clusterID)
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, fmt.Sprintf("cluster ID '%s' not found or unavailable", clusterID), err.Error())
		return nil, false
	}

	return client, true
}
