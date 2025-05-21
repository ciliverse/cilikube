// api/v1/handlers/node_handler.go
package handlers

import (
	"fmt"
	"log"
	"net/http"

	// 确保路径正确
	"github.com/ciliverse/cilikube/api/v1/response" // 假设你有一个通用的响应封装
	"github.com/ciliverse/cilikube/internal/service"

	"github.com/gin-gonic/gin"
	// corev1 "k8s.io/api/core/v1" // 用于 godoc 中的返回类型示例
)

// NodeHandler 负责处理与节点相关的 HTTP 请求
type NodeHandler struct {
	nodeService *service.NodeService // 依赖注入 NodeService
}

// NewNodeHandler 创建 NodeHandler 实例
func NewNodeHandler(nodeSvc *service.NodeService) *NodeHandler {
	if nodeSvc == nil {
		// 在实际应用中，依赖注入失败通常应该导致程序启动失败
		panic("NodeService 依赖不能为空，NodeHandler 初始化失败。")
	}
	return &NodeHandler{nodeService: nodeSvc}
}

// ListNodes godoc
// @Summary      列出指定集群中的所有节点
// @Description  从指定的 Kubernetes 集群获取节点列表。集群名称从路径参数中提取。
// @Tags         Nodes
// @Accept       json
// @Produce      json
// @Param        clusterName path string true "目标集群的名称 (例如 'prod-cluster' 或 'dev-local')"
// @Success      200  {object} response.SuccessResponse{data=corev1.NodeList} "成功响应，包含节点列表"
// @Failure      400  {object} response.ErrorResponse "请求错误 (例如路径参数缺失)"
// @Failure      404  {object} response.ErrorResponse "指定的集群未找到或其客户端不可用"
// @Failure      500  {object} response.ErrorResponse "服务器内部错误 (例如 K8s API 调用失败或 ClientManager 获取失败)"
// @Router       /api/v1/clusters/{clusterName}/nodes [get]
func (h *NodeHandler) ListNodes(c *gin.Context) {
	clusterName := c.Param("clusterName") // 从路径参数中获取集群名称
	log.Printf("处理器层: 收到列出集群 '%s' 节点列表的请求。\n", clusterName)

	if clusterName == "" {
		// 理论上 Gin 的路由匹配会确保 clusterName 存在，但作为防御性检查
		log.Println("错误: ListNodes 请求中的 clusterName 路径参数为空。")
		response.SendError(c, http.StatusBadRequest, "路径参数 'clusterName' 不能为空", nil)
		return
	}

	nodes, err := h.nodeService.ListNodes(clusterName)
	if err != nil {
		log.Printf("错误: 服务层在列出集群 '%s' 节点时发生错误: %v\n", clusterName, err)
		// TODO: 更细致的错误处理。例如，如果错误表明集群客户端不可用，可能返回 404 或 503。
		// 目前统一返回 500，但错误信息会传递给客户端。
		response.SendError(c, http.StatusInternalServerError, fmt.Sprintf("获取集群 '%s' 节点列表失败", clusterName), err.Error())
		return
	}

	log.Printf("处理器层: 成功处理列出集群 '%s' 节点列表的请求。\n", clusterName)
	response.SendSuccess(c, fmt.Sprintf("成功获取集群 '%s' 的节点列表", clusterName), nodes)
}

// GetNodeDetails godoc
// @Summary      获取指定集群中特定节点的详细信息
// @Description  根据集群名称和节点名称获取节点的详细规格和状态。
// @Tags         Nodes
// @Accept       json
// @Produce      json
// @Param        clusterName path string true "目标集群的名称"
// @Param        nodeName    path string true "目标节点的名称"
// @Success      200  {object} response.SuccessResponse{data=corev1.Node} "成功响应，包含节点详细信息"
// @Failure      400  {object} response.ErrorResponse "请求参数错误 (例如节点名称为空)"
// @Failure      404  {object} response.ErrorResponse "指定的集群或该集群中的节点未找到"
// @Failure      500  {object} response.ErrorResponse "服务器内部错误"
// @Router       /api/v1/clusters/{clusterName}/nodes/{nodeName} [get]
func (h *NodeHandler) GetNodeDetails(c *gin.Context) {
	clusterName := c.Param("clusterName")
	nodeName := c.Param("nodeName") // 从路径参数获取节点名称
	log.Printf("处理器层: 收到获取集群 '%s' 节点 '%s' 详细信息的请求。\n", clusterName, nodeName)

	if clusterName == "" {
		log.Println("错误: GetNodeDetails 请求中的 clusterName 路径参数为空。")
		response.SendError(c, http.StatusBadRequest, "路径参数 'clusterName' 不能为空", nil)
		return
	}
	if nodeName == "" {
		log.Println("错误: GetNodeDetails 请求中的 nodeName 路径参数为空。")
		response.SendError(c, http.StatusBadRequest, "路径参数 'nodeName' 不能为空", nil)
		return
	}

	node, err := h.nodeService.GetNodeDetails(clusterName, nodeName)
	if err != nil {
		log.Printf("错误: 服务层在获取集群 '%s' 节点 '%s' 详细信息时发生错误: %v\n", clusterName, nodeName, err)
		// TODO: 更细致的错误处理 (例如区分 k8serrors.IsNotFound(err))
		response.SendError(c, http.StatusInternalServerError, fmt.Sprintf("获取集群 '%s' 节点 '%s' 详细信息失败", clusterName, nodeName), err.Error())
		return
	}

	log.Printf("处理器层: 成功处理获取集群 '%s' 节点 '%s' 详细信息的请求。\n", clusterName, nodeName)
	response.SendSuccess(c, fmt.Sprintf("成功获取集群 '%s' 节点 '%s' 的详细信息", clusterName, nodeName), node)
}
