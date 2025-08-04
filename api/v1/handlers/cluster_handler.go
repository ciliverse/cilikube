package handlers

import (
	"fmt"
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
		ID   string `json:"id"`
		Name string `json:"name"` // 保持向后兼容
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}
	
	var targetID string
	if req.ID != "" {
		targetID = req.ID
	} else if req.Name != "" {
		// 向后兼容：根据名称查找集群ID
		clusters := h.service.ListClusters()
		for _, cluster := range clusters {
			if cluster.Name == req.Name {
				targetID = cluster.ID
				break
			}
		}
		if targetID == "" {
			utils.ApiError(c, http.StatusNotFound, "集群不存在", fmt.Sprintf("未找到名为 '%s' 的集群", req.Name))
			return
		}
	} else {
		utils.ApiError(c, http.StatusBadRequest, "请求参数错误", "必须提供 id 或 name 参数")
		return
	}
	
	if err := h.service.SetActiveCluster(targetID); err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "切换活动集群失败", err.Error())
		return
	}
	utils.ApiSuccess(c, gin.H{"activeClusterID": targetID}, "活动集群切换成功")
}

// GetActiveCluster 获取当前活动集群
func (h *ClusterHandler) GetActiveCluster(c *gin.Context) {
	activeClusterID := h.service.GetActiveClusterID()
	if activeClusterID == "" {
		utils.ApiError(c, http.StatusNotFound, "当前没有活动集群", "请先添加并激活一个集群")
		return
	}
	
	// 直接返回活动集群的名称，如果需要集群详情，可以通过其他API获取
	utils.ApiSuccess(c, activeClusterID, "成功获取活动集群")
}
