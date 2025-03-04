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

// NamespaceHandler ...
type NamespaceHandler struct {
	service *service.NamespaceService
}

// NewNamespaceHandler ...
func NewNamespaceHandler(svc *service.NamespaceService) *NamespaceHandler {
	return &NamespaceHandler{service: svc}
}

// ListNamespaces ...
func (h *NamespaceHandler) ListNamespaces(c *gin.Context) {
	// 1. 调用服务层获取Namespace列表
	namespaces, err := h.service.List(c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Namespace列表失败: "+err.Error())
		return
	}

	// 2. 返回结果
	respondSuccess(c, http.StatusOK, namespaces)
}

// CreateNamespace ...
func (h *NamespaceHandler) CreateNamespace(c *gin.Context) {
	var req models.CreateNamespaceRequest

	// 1. 参数校验
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Namespace格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建Namespace
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
	}

	createdNamespace, err := h.service.Create(namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建Namespace失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNamespaceResponse(createdNamespace))
}

// GetNamespace ...
func (h *NamespaceHandler) GetNamespace(c *gin.Context) {
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Namespace名称格式")
		return
	}

	// 2. 调用服务层获取Namespace详情
	namespace, err := h.service.Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Namespace不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Namespace失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNamespaceResponse(namespace))
}

// UpdateNamespace ...
func (h *NamespaceHandler) UpdateNamespace(c *gin.Context) {
	name := c.Param("name")
	var req models.UpdateNamespaceRequest

	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Namespace名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Namespace格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新Namespace
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
	}

	updatedNamespace, err := h.service.Update(namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新Namespace失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNamespaceResponse(updatedNamespace))
}

// DeleteNamespace ...
func (h *NamespaceHandler) DeleteNamespace(c *gin.Context) {
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Namespace名称格式")
		return
	}

	// 2. 调用服务层删除Namespace
	if err := h.service.Delete(name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Namespace不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Namespace失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchNamespaces ...
func (h *NamespaceHandler) WatchNamespaces(c *gin.Context) {
	// 1. 调用服务层Watch Namespaces
	watcher, err := h.service.Watch(c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch Namespaces失败: "+err.Error())
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
