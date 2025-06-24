package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

// NodeMetricsHandler 处理节点指标相关的请求
type NodeMetricsHandler struct {
	service        *service.NodeMetricsService
	clusterManager *k8s.ClusterManager
}

// NewNodeMetricsHandler 创建一个新的 NodeMetricsHandler 实例
func NewNodeMetricsHandler(svc *service.NodeMetricsService, k8sManager *k8s.ClusterManager) *NodeMetricsHandler {
	return &NodeMetricsHandler{
		service:        svc,
		clusterManager: k8sManager,
	}
}

// GetNodeMetrics 是获取单个节点实时指标的 HTTP 处理函数
func (h *NodeMetricsHandler) GetNodeMetrics(c *gin.Context) {
	// 1. 从查询参数中获取 clusterId，并得到对应集群的 k8s 客户端
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return // 错误已在 GetClientFromQuery 中处理
	}

	// 2. 从路径参数中获取节点名称
	nodeName := c.Param("name")
	if nodeName == "" {
		utils.ApiError(c, http.StatusBadRequest, "节点名称不能为空", "")
		return
	}

	// 3. 调用 service 层获取指标，注意需要传入 k8sClient.Config
	metrics, err := h.service.GetNodeMetrics(k8sClient.Config, nodeName)
	if err != nil {
		// 这里对错误进行判断，如果是因为 metrics-server 未安装导致的，给出友好提示
		if clientErr, ok := err.(interface{ IsNotFound() bool }); ok && clientErr.IsNotFound() {
			utils.ApiError(c, http.StatusNotFound, "获取指标失败", "请确认 Metrics-Server 是否已在目标集群中正确安装并运行。")
			return
		}
		utils.ApiError(c, http.StatusInternalServerError, "获取节点指标失败", err.Error())
		return
	}

	// 4. 成功返回数据
	utils.ApiSuccess(c, metrics, "成功获取节点指标")
}
