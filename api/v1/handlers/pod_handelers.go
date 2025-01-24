package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type PodHandler struct {
	service *service.PodService
}

func NewPodHandler(svc *service.PodService) *PodHandler {
	return &PodHandler{service: svc}
}

// GetPod 获取Pod详情
// @Summary 获取Pod详情
// @Tags Pods
// @Param namespace path string true "Namespace"
// @Param name path string true "Pod名称"
// @Success 200 {object} models.PodResponse
// @Router /api/v1/namespaces/{namespace}/pods/{name} [get]
func (h *PodHandler) GetPod(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Pod名称格式")
		return
	}

	pod, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Pod失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, models.ToPodResponse(pod))
}

// CreatePod 创建Pod
// @Summary 创建新Pod
// @Tags Pods
// @Accept json
// @Param namespace path string true "Namespace"
// @Param pod body models.CreatePodRequest true "Pod配置"
// @Success 201 {object} models.PodResponse
// @Router /api/v1/namespaces/{namespace}/pods [post]
func (h *PodHandler) CreatePod(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))

	var req models.CreatePodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求格式: "+err.Error())
		return
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdPod, err := h.service.Create(namespace, pod)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			respondError(c, http.StatusConflict, "Pod已存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "创建Pod失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, models.ToPodResponse(createdPod))
}

// UpdatePod 更新Pod
// @Summary 更新现有Pod
// @Tags Pods
// @Accept json
// @Param namespace path string true "Namespace"
// @Param name path string true "Pod名称"
// @Param pod body models.UpdatePodRequest true "更新配置"
// @Success 200 {object} models.PodResponse
// @Router /api/v1/namespaces/{namespace}/pods/{name} [put]
func (h *PodHandler) UpdatePod(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	var req models.UpdatePodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求格式: "+err.Error())
		return
	}

	existingPod, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Pod失败: "+err.Error())
		return
	}

	// 保留不可变字段
	updatedPod := existingPod.DeepCopy()
	updatedPod.Labels = req.Labels
	updatedPod.Annotations = req.Annotations
	updatedPod.Spec = req.Spec

	result, err := h.service.Update(namespace, updatedPod)
	if err != nil {
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "Pod已被修改，请重试")
			return
		}
		respondError(c, http.StatusInternalServerError, "更新Pod失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, models.ToPodResponse(result))
}

// DeletePod 删除Pod
// @Summary 删除指定Pod
// @Tags Pods
// @Param namespace path string true "Namespace"
// @Param name path string true "Pod名称"
// @Success 204
// @Router /api/v1/namespaces/{namespace}/pods/{name} [delete]
func (h *PodHandler) DeletePod(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if err := h.service.Delete(namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Pod失败: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// ListPods 列出Pod
// @Summary 列出命名空间下的所有Pod
// @Tags Pods
// @Param namespace path string true "Namespace"
// @Param labelSelector query string false "标签选择器"
// @Param limit query int false "返回数量限制"
// @Success 200 {object} models.PodListResponse
// @Router /api/v1/namespaces/{namespace}/pods [get]
func (h *PodHandler) ListPods(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	labelSelector := c.Query("labelSelector")
	limit := utils.ParseInt(c.DefaultQuery("limit", "100"), 100)

	pods, err := h.service.List(namespace, labelSelector, int64(limit))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
		return
	}

	response := models.PodListResponse{
		Items: make([]models.PodResponse, 0, len(pods.Items)),
		Total: len(pods.Items),
	}

	for _, pod := range pods.Items {
		response.Items = append(response.Items, models.ToPodResponse(&pod))
	}

	respondSuccess(c, http.StatusOK, response)
}

// WatchPods 监听Pod变化
// @Summary 实时监听Pod变化
// @Tags Pods
// @Param namespace path string true "Namespace"
// @Param labelSelector query string false "标签选择器"
// @Success 200 {object} models.WatchEvent
// @Router /api/v1/watch/namespaces/{namespace}/pods [get]
func (h *PodHandler) WatchPods(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	labelSelector := c.Query("labelSelector")

	watcher, err := h.service.Watch(namespace, labelSelector)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "监听Pod失败: "+err.Error())
		return
	}
	defer watcher.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return false
			}
			c.SSEvent("message", toWatchEvent(event))
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

func toWatchEvent(event watch.Event) interface{} {
	pod, ok := event.Object.(*corev1.Pod)
	if !ok {
		return map[string]interface{}{
			"type":  event.Type,
			"error": "类型转换失败",
		}
	}

	return map[string]interface{}{
		"type":   event.Type,
		"object": models.ToPodResponse(pod),
	}
}

// 统一响应处理
func respondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"code":    code,
		"data":    data,
		"message": "success",
	})
}

func respondError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
