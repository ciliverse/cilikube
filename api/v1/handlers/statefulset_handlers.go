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

// StatefulSetHandler ...
type StatefulSetHandler struct {
	service        *service.StatefulSetService
	clusterManager *k8s.ClusterManager
}

// NewStatefulSetHandler ...
func NewStatefulSetHandler(svc *service.StatefulSetService, cm *k8s.ClusterManager) *StatefulSetHandler {
	return &StatefulSetHandler{service: svc, clusterManager: cm}
}

// ListStatefulSets ...
func (h *StatefulSetHandler) ListStatefulSets(c *gin.Context) {
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

	// 2. 调用服务层获取StatefulSet列表
	statefulSets, err := h.service.List(k8sClient.Clientset, namespace, c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取StatefulSet列表失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, statefulSets)
}

// CreateStatefulSet ...
func (h *StatefulSetHandler) CreateStatefulSet(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	var req models.CreateStatefulSetRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的StatefulSet格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建StatefulSet
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdStatefulSet, err := h.service.Create(k8sClient.Clientset, namespace, statefulSet)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建StatefulSet失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToStatefulSetResponse(createdStatefulSet))
}

// GetStatefulSet ...
func (h *StatefulSetHandler) GetStatefulSet(c *gin.Context) {
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
		respondError(c, http.StatusBadRequest, "无效的StatefulSet名称格式")
		return
	}

	// 2. 调用服务层获取StatefulSet详情
	statefulSet, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "StatefulSet不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取StatefulSet失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToStatefulSetResponse(statefulSet))
}

// UpdateStatefulSet ...
func (h *StatefulSetHandler) UpdateStatefulSet(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	var req models.UpdateStatefulSetRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的StatefulSet名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的StatefulSet格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新StatefulSet
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	updatedStatefulSet, err := h.service.Update(k8sClient.Clientset, namespace, statefulSet)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新StatefulSet失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToStatefulSetResponse(updatedStatefulSet))
}

// DeleteStatefulSet ...
func (h *StatefulSetHandler) DeleteStatefulSet(c *gin.Context) {
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
		respondError(c, http.StatusBadRequest, "无效的StatefulSet名称格式")
		return
	}

	// 2. 调用服务层删除StatefulSet
	if err := h.service.Delete(k8sClient.Clientset, namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "StatefulSet不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除StatefulSet失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchStatefulSets ...
func (h *StatefulSetHandler) WatchStatefulSets(c *gin.Context) {
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

	// 2. 调用服务层Watch StatefulSets
	watcher, err := h.service.Watch(k8sClient.Clientset, namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch StatefulSets失败: "+err.Error())
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
