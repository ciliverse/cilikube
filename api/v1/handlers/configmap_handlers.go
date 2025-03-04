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

// ConfigMapHandler ...
type ConfigMapHandler struct {
	service *service.ConfigMapService
}

// NewConfigMapHandler ...
func NewConfigMapHandler(svc *service.ConfigMapService) *ConfigMapHandler {
	return &ConfigMapHandler{service: svc}
}

// ListConfigMaps ...
func (h *ConfigMapHandler) ListConfigMaps(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层获取ConfigMap列表
	configMaps, err := h.service.List(namespace, c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取ConfigMap列表失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, configMaps)
}

// CreateConfigMap ...
func (h *ConfigMapHandler) CreateConfigMap(c *gin.Context) {
	namespace := c.Param("namespace")
	var req models.CreateConfigMapRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的ConfigMap格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Data: req.Data,
	}

	createdConfigMap, err := h.service.Create(namespace, configMap)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建ConfigMap失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToConfigMapResponse(createdConfigMap))
}

// GetConfigMap ...
func (h *ConfigMapHandler) GetConfigMap(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的ConfigMap名称格式")
		return
	}

	// 2. 调用服务层获取ConfigMap详情
	configMap, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "ConfigMap不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取ConfigMap失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToConfigMapResponse(configMap))
}

// UpdateConfigMap ...
func (h *ConfigMapHandler) UpdateConfigMap(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	var req models.UpdateConfigMapRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的ConfigMap名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的ConfigMap格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Data: req.Data,
	}

	updatedConfigMap, err := h.service.Update(namespace, configMap)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新ConfigMap失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToConfigMapResponse(updatedConfigMap))
}

// DeleteConfigMap ...
func (h *ConfigMapHandler) DeleteConfigMap(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的ConfigMap名称格式")
		return
	}

	// 2. 调用服务层删除ConfigMap
	if err := h.service.Delete(namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "ConfigMap不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除ConfigMap失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchConfigMaps ...
func (h *ConfigMapHandler) WatchConfigMaps(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层Watch ConfigMaps
	watcher, err := h.service.Watch(namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch ConfigMaps失败: "+err.Error())
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
