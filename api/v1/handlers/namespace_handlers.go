package handlers

import (
	"io"
	"net/http"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceHandler ...
type NamespaceHandler struct {
	service        *service.NamespaceService
	clusterManager *k8s.ClusterManager
}

// NewNamespaceHandler ...
func NewNamespaceHandler(svc *service.NamespaceService, cm *k8s.ClusterManager) *NamespaceHandler {
	return &NamespaceHandler{service: svc, clusterManager: cm}
}

// ListNamespaces ...
func (h *NamespaceHandler) ListNamespaces(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}
	namespaces, err := h.service.List(k8sClient.Clientset, c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Namespace列表失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, namespaces)
}

// CreateNamespace ...
func (h *NamespaceHandler) CreateNamespace(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}
	var req models.CreateNamespaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Namespace格式: "+err.Error())
		return
	}
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
	}
	createdNamespace, err := h.service.Create(k8sClient.Clientset, namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建Namespace失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, models.ToNamespaceResponse(createdNamespace))
}

// GetNamespace ...
func (h *NamespaceHandler) GetNamespace(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}
	name := c.Param("name")
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Namespace名称格式")
		return
	}
	namespace, err := h.service.Get(k8sClient.Clientset, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Namespace不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Namespace失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, models.ToNamespaceResponse(namespace))
}

// UpdateNamespace ...
func (h *NamespaceHandler) UpdateNamespace(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}
	name := c.Param("name")
	var req models.UpdateNamespaceRequest
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Namespace名称格式")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Namespace格式: "+err.Error())
		return
	}
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
	}
	updatedNamespace, err := h.service.Update(k8sClient.Clientset, namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新Namespace失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, models.ToNamespaceResponse(updatedNamespace))
}

// DeleteNamespace ...
func (h *NamespaceHandler) DeleteNamespace(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}
	name := c.Param("name")
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Namespace名称格式")
		return
	}
	if err := h.service.Delete(k8sClient.Clientset, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Namespace不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Namespace失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchNamespaces ...
func (h *NamespaceHandler) WatchNamespaces(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}
	watcher, err := h.service.Watch(k8sClient.Clientset, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch Namespaces失败: "+err.Error())
		return
	}
	c.Stream(func(w io.Writer) bool {
		event, ok := <-watcher.ResultChan()
		if !ok {
			return false
		}
		c.SSEvent("message", event)
		return true
	})
}
