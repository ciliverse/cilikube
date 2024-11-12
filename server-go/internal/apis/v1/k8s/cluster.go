package k8s

import (
	"net/http"

	"github.com/cilliandevops/kopsadmin/server-go/internal/apis/services"
	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	clusterService *services.ClusterService
}

func NewClusterHandler(clusterService *services.ClusterService) *ClusterHandler {
	return &ClusterHandler{
		clusterService: clusterService,
	}
}

// GetClusterInfo handles getting cluster information
func (h *ClusterHandler) GetClusterInfo(c *gin.Context) {
	// Get the context from the Gin request
	ctx := c.Request.Context()

	// Call the service method with the context
	info, err := h.clusterService.GetClusterInfo(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// ListNodes handles listing all nodes in the cluster
func (h *ClusterHandler) ListNodes(c *gin.Context) {
	// Get the context from the Gin request
	ctx := c.Request.Context()

	// Call the service method with the context
	nodes, err := h.clusterService.ListNodes(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

// GetNode handles getting a specific node's information
func (h *ClusterHandler) GetNode(c *gin.Context) {
	// Get the context from the Gin request
	ctx := c.Request.Context()

	// Extract node name from URL parameter
	nodeName := c.Param("name")

	// Call the service method with the context
	node, err := h.clusterService.GetNode(ctx, nodeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

// DeleteNode handles deleting a specific node from the cluster
func (h *ClusterHandler) DeleteNode(c *gin.Context) {
	// Get the context from the Gin request
	ctx := c.Request.Context()

	// Extract node name from URL parameter
	nodeName := c.Param("name")

	// Call the service method with the context
	err := h.clusterService.DeleteNode(ctx, nodeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Node deleted"})
}
