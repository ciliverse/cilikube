package handlers

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	"net/http"
	"strings"

	"github.com/ciliverse/cilikube/api/v1/models" // Adjust import path
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type ConfigMapHandler struct {
	service        *service.ConfigMapService
	clusterManager *k8s.ClusterManager
}

func NewConfigMapHandler(svc *service.ConfigMapService, cm *k8s.ClusterManager) *ConfigMapHandler {
	return &ConfigMapHandler{service: svc, clusterManager: cm}
}

// ListConfigMaps godoc
// @Summary List ConfigMaps
// @Description Get a list of all ConfigMaps in a specific namespace
// @Tags ConfigMaps
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param labelSelector query string false "Label selector for filtering"
// @Param limit query int false "Maximum number of items to return" default(100)
// @Success 200 {object} models.ConfigMapListResponse "List of ConfigMaps"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Namespace"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/configmaps [get]
func (h *ConfigMapHandler) ListConfigMaps(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	labelSelector := c.Query("labelSelector")
	limit := utils.ParseInt(c.DefaultQuery("limit", "100"), 100)

	cmList, err := h.service.List(k8sClient.Clientset, namespace, labelSelector, int64(limit))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取ConfigMap列表失败: "+err.Error())
		return
	}

	response := models.ConfigMapListResponse{
		Items: make([]models.ConfigMapResponse, 0, len(cmList.Items)),
		Total: len(cmList.Items), // Adjust if backend provides actual total
	}
	for _, cm := range cmList.Items {
		response.Items = append(response.Items, models.ToConfigMapResponse(&cm))
	}

	respondSuccess(c, http.StatusOK, response)
}

// GetConfigMap godoc
// @Summary Get a ConfigMap
// @Description Get details of a specific ConfigMap by namespace and name, including data
// @Tags ConfigMaps
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "ConfigMap Name"
// @Success 200 {object} models.ConfigMapDetailResponse "ConfigMap details" // Use detail response
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Name/Namespace"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/configmaps/{name} [get]
func (h *ConfigMapHandler) GetConfigMap(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或ConfigMap名称格式")
		return
	}

	cm, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "ConfigMap不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取ConfigMap失败: "+err.Error())
		return
	}
	// Use the detail response model including Data
	respondSuccess(c, http.StatusOK, models.ToConfigMapDetailResponse(cm))
}

// CreateConfigMap godoc
// @Summary Create a ConfigMap
// @Description Create a new ConfigMap using YAML/JSON definition
// @Tags ConfigMaps
// @Accept json,yaml
// @Produce json
// @Param namespace path string true "Namespace"
// @Param configmap body corev1.ConfigMap true "ConfigMap definition (ensure kind: ConfigMap, apiVersion: v1)"
// @Success 201 {object} models.ConfigMapResponse "ConfigMap created" // Return basic info
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Input"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Already Exists"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/configmaps [post]
func (h *ConfigMapHandler) CreateConfigMap(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	var cm corev1.ConfigMap
	if err := c.ShouldBindJSON(&cm); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求体格式: "+err.Error())
		return
	}

	if cm.Kind != "ConfigMap" || (cm.APIVersion != "v1" && cm.APIVersion != "") {
		respondError(c, http.StatusBadRequest, "无效的Kind或APIVersion，应为 ConfigMap/v1")
		return
	}
	if cm.APIVersion == "" {
		cm.APIVersion = "v1"
	}

	createdCM, err := h.service.Create(k8sClient.Clientset, namespace, &cm)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			respondError(c, http.StatusConflict, "ConfigMap已存在")
			return
		}
		if _, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, "创建ConfigMap验证失败: "+err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "创建ConfigMap失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusCreated, models.ToConfigMapResponse(createdCM)) // Return basic info
}

// UpdateConfigMap godoc
// @Summary Update a ConfigMap
// @Description Update an existing ConfigMap using YAML/JSON definition
// @Tags ConfigMaps
// @Accept json,yaml
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "ConfigMap Name"
// @Param configmap body corev1.ConfigMap true "Updated ConfigMap definition"
// @Success 200 {object} models.ConfigMapResponse "ConfigMap updated" // Return basic info
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Input or Name Mismatch"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Resource modified"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/configmaps/{name} [put]
func (h *ConfigMapHandler) UpdateConfigMap(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或ConfigMap名称格式")
		return
	}

	var cm corev1.ConfigMap
	if err := c.ShouldBindJSON(&cm); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求体格式: "+err.Error())
		return
	}
	if cm.Name != name || cm.Namespace != namespace {
		respondError(c, http.StatusBadRequest, "路径参数与请求体中的名称/命名空间不匹配")
		return
	}
	if cm.Kind != "ConfigMap" || (cm.APIVersion != "v1" && cm.APIVersion != "") {
		respondError(c, http.StatusBadRequest, "无效的Kind或APIVersion，应为 ConfigMap/v1")
		return
	}
	if cm.APIVersion == "" {
		cm.APIVersion = "v1"
	}

	updatedCM, err := h.service.Update(k8sClient.Clientset, namespace, &cm)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "ConfigMap不存在")
			return
		}
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "资源已被修改，请获取最新版本后重试")
			return
		}
		if _, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, "更新ConfigMap验证失败: "+err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "更新ConfigMap失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, models.ToConfigMapResponse(updatedCM)) // Return basic info
}

// DeleteConfigMap godoc
// @Summary Delete a ConfigMap
// @Description Delete a specific ConfigMap by namespace and name
// @Tags ConfigMaps
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "ConfigMap Name"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Name/Namespace"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/configmaps/{name} [delete]
func (h *ConfigMapHandler) DeleteConfigMap(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或ConfigMap名称格式")
		return
	}

	err := h.service.Delete(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			c.Status(http.StatusNoContent)
			return
		} // Idempotent
		respondError(c, http.StatusInternalServerError, "删除ConfigMap失败: "+err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Re-use or define respond helpers ---
/*
type ErrorResponse struct { Code int `json:"code"`; Message string `json:"message"`}
func respondSuccess(c *gin.Context, code int, data interface{}) { ... }
func respondError(c *gin.Context, code int, message string) { ... }
*/
