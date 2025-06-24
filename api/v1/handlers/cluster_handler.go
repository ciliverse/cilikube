package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/api/v1/models"
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

// ListClusters 获取集群列表
func (h *ClusterHandler) ListClusters(c *gin.Context) {
	clusters := h.service.ListClusters()
	utils.ApiSuccess(c, clusters, "成功获取集群列表")
}

// GetCluster 获取单个集群详情
func (h *ClusterHandler) GetCluster(c *gin.Context) {
	clusterID := c.Param("id")
	cluster, err := h.service.GetClusterByID(clusterID)
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, "获取集群失败", err.Error())
		return
	}
	utils.ApiSuccess(c, cluster, "成功获取集群详情")
}

// CreateCluster 创建一个新集群
func (h *ClusterHandler) CreateCluster(c *gin.Context) {
	var req models.CreateClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}
	if err := h.service.CreateCluster(req); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "创建集群失败", err.Error())
		return
	}
	utils.ApiSuccess(c, nil, "集群创建成功")
}

// UpdateCluster 更新一个现有集群
func (h *ClusterHandler) UpdateCluster(c *gin.Context) {
	clusterID := c.Param("id")
	var req models.UpdateClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}
	if err := h.service.UpdateCluster(clusterID, req); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "更新集群失败", err.Error())
		return
	}
	utils.ApiSuccess(c, nil, "集群更新成功")
}

// DeleteCluster 删除一个集群
func (h *ClusterHandler) DeleteCluster(c *gin.Context) {
	clusterID := c.Param("id")
	if err := h.service.DeleteClusterByID(clusterID); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "删除集群失败", err.Error())
		return
	}
	utils.ApiSuccess(c, nil, "集群删除成功")
}

// SetActiveCluster 设定当前活动集群
func (h *ClusterHandler) SetActiveCluster(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}
	if err := h.service.SetActiveCluster(req.ID); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "切换活动集群失败", err.Error())
		return
	}
	utils.ApiSuccess(c, gin.H{"activeClusterID": req.ID}, "活动集群切换成功")
}
