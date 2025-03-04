package handlers

import (
	"io"
	"net/http"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeHandler ...
type NodeHandler struct {
	service *service.NodeService
}

// NewNodeHandler ...
func NewNodeHandler(svc *service.NodeService) *NodeHandler {
	return &NodeHandler{service: svc}
}

// ListNodes ...
func (h *NodeHandler) ListNodes(c *gin.Context) {
	// 1. 调用服务层获取Node列表
	nodes, err := h.service.List(c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Node列表失败: "+err.Error())
		return
	}

	// 2. 返回结果
	respondSuccess(c, http.StatusOK, nodes)
}

// CreateNode ...
func (h *NodeHandler) CreateNode(c *gin.Context) {
	var req models.CreateNodeRequest

	// 1. 参数校验
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Node格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建Node
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdNode, err := h.service.Create(node)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建Node失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNodeResponse(createdNode))
}

// GetNode ...
func (h *NodeHandler) GetNode(c *gin.Context) {
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Node名称格式")
		return
	}

	// 2. 调用服务层获取Node详情
	node, err := h.service.Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Node不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Node失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNodeResponse(node))
}

// UpdateNode ...
func (h *NodeHandler) UpdateNode(c *gin.Context) {
	name := c.Param("name")
	var req models.UpdateNodeRequest

	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Node名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Node格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新Node
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	updatedNode, err := h.service.Update(node)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新Node失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNodeResponse(updatedNode))
}

// DeleteNode ...
func (h *NodeHandler) DeleteNode(c *gin.Context) {
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Node名称格式")
		return
	}

	// 2. 调用服务层删除Node
	if err := h.service.Delete(name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Node不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Node失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchNodes ...
func (h *NodeHandler) WatchNodes(c *gin.Context) {
	// 1. 调用服务层Watch Nodes
	watcher, err := h.service.Watch(c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch Nodes失败: "+err.Error())
		return
	}

	// 2. 返回结果
	c.Stream(func(w io.Writer) bool {
		event, ok := <-watcher.ResultChan()
		if !ok {
			return false
		}
		c.SSEvent("message", event)
		return true
	})
}
