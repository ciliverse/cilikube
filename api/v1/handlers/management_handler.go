// api/v1/handlers/management_handler.go
package handlers

import (
	"log"
	"net/http"

	"github.com/ciliverse/cilikube/api/v1/response" // 你的响应封装
	"github.com/ciliverse/cilikube/pkg/k8s"         // 直接使用 ClientManager

	"github.com/gin-gonic/gin"
)

// ManagementHandler 负责处理与平台管理相关的 HTTP 请求
type ManagementHandler struct {
	// 此 Handler 可能不需要依赖特定的 service，而是直接与 ClientManager 交互
}

// NewManagementHandler 创建 ManagementHandler 实例
func NewManagementHandler() *ManagementHandler {
	return &ManagementHandler{}
}

// ListAvailableClusters godoc
// @Summary      列出所有已配置且客户端可用的集群名称
// @Description  获取平台当前能够成功连接或已配置的所有 Kubernetes 集群的名称列表。此列表可用于前端的集群选择器。
// @Tags         Management
// @Accept       json
// @Produce      json
// @Success      200  {object} response.SuccessResponse{data=[]string} "成功响应，包含已排序的集群名称数组"
// @Failure      500  {object} response.ErrorResponse "服务器内部错误 (例如 ClientManager 获取失败)"
// @Router       /api/v1/management/clusters [get]
func (h *ManagementHandler) ListAvailableClusters(c *gin.Context) {
	log.Println("处理器层: 收到列出可用集群列表的请求。")
	manager, err := k8s.GetClientManager()
	if err != nil {
		log.Printf("错误: ListAvailableClusters 获取 ClientManager 失败: %v\n", err)
		response.SendError(c, http.StatusInternalServerError, "获取客户端管理器失败", err.Error())
		return
	}

	clusterNames := manager.ListClusterNames() // ListClusterNames 内部已实现排序
	if clusterNames == nil {                   // ListClusterNames 应该返回空切片而不是 nil
		clusterNames = []string{} // 确保前端总是收到一个数组
	}

	log.Printf("处理器层: 成功获取可用集群列表: %v\n", clusterNames)
	response.SendSuccess(c, "成功获取可用集群列表", clusterNames)
}
