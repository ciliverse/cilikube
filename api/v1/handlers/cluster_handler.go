package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

// ClusterHandler 封装了所有与集群管理相关的 HTTP 处理函数。
type ClusterHandler struct {
	service *service.ClusterService
}

// NewClusterHandler 创建一个新的 ClusterHandler 实例。
func NewClusterHandler(svc *service.ClusterService) *ClusterHandler {
	return &ClusterHandler{service: svc}
}

// ListClusters 处理获取集群列表的 API 请求。
func (h *ClusterHandler) ListClusters(c *gin.Context) {
	clusters := h.service.ListClusters()
	c.JSON(http.StatusOK, gin.H{"data": clusters})
}

// CreateCluster 处理创建新集群的 API 请求。
func (h *ClusterHandler) CreateCluster(c *gin.Context) {
	var req service.CreateClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateCluster(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "集群创建成功"})
}

// DeleteCluster 处理删除集群的 API 请求。
func (h *ClusterHandler) DeleteCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	if clusterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "集群ID不能为空"})
		return
	}
	if err := h.service.DeleteClusterByID(clusterID); err != nil { // 假设服务层有 DeleteClusterByID 方法
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "集群删除成功"})
}

// SetActiveCluster 处理切换活动集群的 API 请求。
func (h *ClusterHandler) SetActiveCluster(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"` // 修改为 ID
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SetActiveClusterByID(req.ID); err != nil { // 假设服务层有 SetActiveClusterByID 方法
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "活动集群切换成功"})
}

// GetActiveCluster 处理获取当前活动集群名称的 API 请求。
func (h *ClusterHandler) GetActiveCluster(c *gin.Context) {
	activeCluster := h.service.GetActiveCluster()
	c.JSON(http.StatusOK, gin.H{"data": activeCluster})
}
