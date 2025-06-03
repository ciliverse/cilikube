package handlers

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	"io"
	"net/http"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSetHandler ...
type DaemonSetHandler struct {
	service        *service.DaemonSetService
	clusterManager *k8s.ClusterManager
}

// NewDaemonSetHandler ...
func NewDaemonSetHandler(svc *service.DaemonSetService, cm *k8s.ClusterManager) *DaemonSetHandler {
	return &DaemonSetHandler{service: svc, clusterManager: cm}
}

// ListDaemonSets ...
func (h *DaemonSetHandler) ListDaemonSets(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层获取DaemonSet列表
	daemonsets, err := h.service.List(k8sClient.Clientset, namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取DaemonSet列表失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, daemonsets)
}

// CreateDaemonSet ...
func (h *DaemonSetHandler) CreateDaemonSet(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	var req models.CreateDaemonSetRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的DaemonSet格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建DaemonSet
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdDaemonset, err := h.service.Create(k8sClient.Clientset, namespace, daemonset)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建DaemonSet失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToDaemonSetResponse(createdDaemonset))
}

// GetDaemonSet ...
func (h *DaemonSetHandler) GetDaemonSet(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的DaemonSet名称格式")
		return
	}

	// 2. 调用服务层获取DaemonSet详情
	daemonset, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "DaemonSet不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取DaemonSet失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToDaemonSetResponse(daemonset))
}

// UpdateDaemonSet ...
func (h *DaemonSetHandler) UpdateDaemonSet(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	var req models.UpdateDaemonSetRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的DaemonSet名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的DaemonSet格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新DaemonSet
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	updatedDaemonset, err := h.service.Update(k8sClient.Clientset, namespace, daemonset)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新DaemonSet失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToDaemonSetResponse(updatedDaemonset))
}

// DeleteDaemonSet ...
func (h *DaemonSetHandler) DeleteDaemonSet(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的DaemonSet名称格式")
		return
	}

	// 2. 调用服务层删除DaemonSet
	if err := h.service.Delete(k8sClient.Clientset, namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "DaemonSet不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除DaemonSet失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchDaemonSets ...
func (h *DaemonSetHandler) WatchDaemonSets(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层Watch DaemonSets
	watcher, err := h.service.Watch(k8sClient.Clientset, namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch DaemonSets失败: "+err.Error())
		return
	}

	// 3. 返回结果
	c.Stream(func(w io.Writer) bool {
		event, ok := <-watcher.ResultChan()
		if !ok {
			return false
		}
		c.SSEvent("message", event)
		return true
	})
}
