package handlers

import (
	"fmt"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"io"
	"net/http"
	"strconv"
	"strings"

	// Keep for potential future use (like WebSocket ping)
	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils" // Assuming utils package exists
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type PodHandler struct {
	service        *service.PodService
	clusterManager *k8s.ClusterManager
}

func NewPodHandler(svc *service.PodService, cm *k8s.ClusterManager) *PodHandler {
	return &PodHandler{service: svc, clusterManager: cm}
}

// ListNamespaces ... (保持不变)
func (h *PodHandler) ListNamespaces(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespaces, err := h.service.ListNamespaces(k8sClient.Clientset)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取命名空间失败: "+err.Error())
		return
	}
	// Use respondSuccess for consistency
	respondSuccess(c, http.StatusOK, namespaces)
}

// GetPod ... (保持不变)
func (h *PodHandler) GetPod(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

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

	pod, err := h.service.Get(k8sClient.Clientset, namespace, name)
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

// CreatePod 创建Pod (支持 JSON 或 YAML)
func (h *PodHandler) CreatePod(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	contentType := c.ContentType()
	var createdPod *corev1.Pod
	var err error

	if strings.Contains(contentType, "yaml") || strings.Contains(contentType, "x-yaml") {
		yamlBody, errRead := io.ReadAll(c.Request.Body)
		if errRead != nil {
			respondError(c, http.StatusBadRequest, "读取 YAML 请求体失败: "+errRead.Error())
			return
		}
		if len(yamlBody) == 0 {
			respondError(c, http.StatusBadRequest, "请求体不能为空 (YAML)")
			return
		}
		createdPod, err = h.service.CreateFromYAML(k8sClient.Clientset, namespace, yamlBody)

	} else if strings.Contains(contentType, "json") { // Explicitly check for JSON
		var req models.CreatePodRequest
		if errBind := c.ShouldBindJSON(&req); errBind != nil {
			respondError(c, http.StatusBadRequest, "无效的 JSON 请求格式: "+errBind.Error())
			return
		}
		// Validate name from JSON body if present
		if !utils.ValidateResourceName(req.Name) {
			respondError(c, http.StatusBadRequest, "无效的 Pod 名称格式 (来自 JSON body)")
			return
		}

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        req.Name,
				Namespace:   namespace, // Use namespace from path, ignore req.Namespace if any
				Labels:      req.Labels,
				Annotations: req.Annotations,
			},
			Spec: req.Spec,
		}
		// Use the original service.Create method for JSON objects
		createdPod, err = h.service.Create(k8sClient.Clientset, namespace, pod)
	} else {
		respondError(c, http.StatusUnsupportedMediaType, "不支持的 Content-Type，请使用 application/json 或 application/yaml")
		return
	}

	// --- Handle Response ---
	if err != nil {
		if e, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, e.Error())
			return
		}
		if errors.IsAlreadyExists(err) {
			respondError(c, http.StatusConflict, "Pod已存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "创建Pod失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, models.ToPodResponse(createdPod))
}

// UpdatePod 更新Pod (支持 JSON 或 YAML)
func (h *PodHandler) UpdatePod(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}

	contentType := c.ContentType()
	var result *corev1.Pod
	var err error

	if strings.Contains(contentType, "yaml") || strings.Contains(contentType, "x-yaml") {
		yamlBody, errRead := io.ReadAll(c.Request.Body)
		if errRead != nil {
			respondError(c, http.StatusBadRequest, "读取 YAML 请求体失败: "+errRead.Error())
			return
		}
		if len(yamlBody) == 0 {
			respondError(c, http.StatusBadRequest, "请求体不能为空 (YAML)")
			return
		}
		result, err = h.service.UpdateFromYAML(k8sClient.Clientset, namespace, name, yamlBody)

	} else if strings.Contains(contentType, "json") { // Explicitly check for JSON
		// --- Handle JSON Input ---
		// Get the existing Pod first to apply changes correctly
		existingPod, errGet := h.service.Get(k8sClient.Clientset, namespace, name)
		if errGet != nil {
			if errors.IsNotFound(errGet) {
				respondError(c, http.StatusNotFound, "Pod不存在，无法更新")
				return
			}
			respondError(c, http.StatusInternalServerError, "获取Pod失败: "+errGet.Error())
			return
		}

		// Bind the JSON request which contains only the fields to update
		var req models.UpdatePodRequest // Assumes this model only contains fields allowed to change
		if errBind := c.ShouldBindJSON(&req); errBind != nil {
			respondError(c, http.StatusBadRequest, "无效的 JSON 请求格式: "+errBind.Error())
			return
		}

		// Create a deep copy and apply changes from the request
		updatedPod := existingPod.DeepCopy()
		updatedPod.Labels = req.Labels           // Replace labels
		updatedPod.Annotations = req.Annotations // Replace annotations
		updatedPod.Spec = req.Spec               // Replace the entire spec

		// *** Call the correct Update method in the service ***
		result, err = h.service.Update(k8sClient.Clientset, namespace, updatedPod) // Use the method taking a Pod object

	} else {
		respondError(c, http.StatusUnsupportedMediaType, "不支持的 Content-Type，请使用 application/json 或 application/yaml")
		return
	}

	// --- Handle Response ---
	if err != nil {
		if e, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, e.Error())
			return
		}
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod不存在 (可能在更新期间被删除)")
			return
		}
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "Pod已被修改，请重试 (ResourceVersion conflict)")
			return
		}
		respondError(c, http.StatusInternalServerError, "更新Pod失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, models.ToPodResponse(result))
}

// DeletePod ... (保持不变, 使用 204)
func (h *PodHandler) DeletePod(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}

	err := h.service.Delete(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			// Idempotent: Return success even if not found
			c.Status(http.StatusNoContent)
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Pod失败: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent) // Success
}

// ListPods ... (保持不变)
func (h *PodHandler) ListPods(c *gin.Context) {
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
	limitStr := c.DefaultQuery("limit", "500") // Sensible default limit
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 500 // Fallback
	}

	pods, err := h.service.List(k8sClient.Clientset, namespace, labelSelector, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
		return
	}

	response := models.PodListResponse{
		Items: make([]models.PodResponse, 0, len(pods.Items)),
		// Total reflects items *in this batch*. K8s list doesn't give total count easily.
		Total: len(pods.Items),
	}
	for _, pod := range pods.Items {
		response.Items = append(response.Items, models.ToPodResponse(&pod))
	}

	respondSuccess(c, http.StatusOK, response)
}

// WatchPods 监听Pod变化 (**修正**)
// @Summary 实时监听Pod变化 (SSE)
// @Tags Pods
// @Param namespace path string true "Namespace"
// @Param labelSelector query string false "标签选择器"
// @Success 200 {object} models.WatchEvent "Server-Sent Event stream"
// @Router /api/v1/namespaces/{namespace}/watch/pods [get] // Adjusted route based on routes.go
func (h *PodHandler) WatchPods(c *gin.Context) {
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

	watcher, err := h.service.Watch(k8sClient.Clientset, namespace, labelSelector)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "开始监听Pod失败: "+err.Error())
		return
	}
	defer watcher.Stop()

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Adjust in production

	// Use c.Stream to handle the streaming goroutine
	// ** REMOVED chanStream := c.Stream(...) and <-chanStream **
	c.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				// Channel closed by Watcher (e.g., timeout, error)
				fmt.Println("Watcher channel closed")
				c.SSEvent("close", gin.H{"message": "Watcher channel closed"}) // Send a close event
				return false                                                   // Stop streaming
			}
			// Send event data
			c.SSEvent("message", toWatchEvent(event))
			// c.Writer.Flush() // Gin's SSEvent might handle flushing
			return true // Keep connection open and continue streaming

		case <-c.Request.Context().Done():
			// Client closed the connection
			fmt.Println("Client disconnected from watch stream")
			return false // Stop streaming
		}
	})

	// The handlers function returns here, but the goroutine inside c.Stream continues
	fmt.Println("WatchPods handlers finished setup, streaming started.")
}

// GetPodYAML ... (保持不变)
func (h *PodHandler) GetPodYAML(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}

	yamlBytes, err := h.service.GetPodYAML(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod 不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取 Pod YAML 失败: "+err.Error())
		return
	}

	c.Header("Content-Type", "application/yaml")
	respondSuccess(c, http.StatusOK, string(yamlBytes))
}

// UpdatePodYAML ... (保持不变)
func (h *PodHandler) UpdatePodYAML(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}

	contentType := c.ContentType()
	if !strings.Contains(contentType, "yaml") && !strings.Contains(contentType, "x-yaml") {
		respondError(c, http.StatusUnsupportedMediaType, "请求体必须是 application/yaml 类型")
		return
	}

	yamlBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		respondError(c, http.StatusBadRequest, "读取 YAML 请求体失败: "+err.Error())
		return
	}
	if len(yamlBody) == 0 {
		respondError(c, http.StatusBadRequest, "请求体不能为空")
		return
	}

	updatedPod, err := h.service.UpdateFromYAML(k8sClient.Clientset, namespace, name, yamlBody)
	if err != nil {
		if e, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, e.Error())
			return
		}
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod 不存在")
			return
		}
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "Pod已被修改，请重试 (ResourceVersion conflict)")
			return
		}
		respondError(c, http.StatusInternalServerError, "更新 Pod YAML 失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, models.ToPodResponse(updatedPod))
}

// --- Helper Functions ---

// toWatchEvent ... (保持不变)
func toWatchEvent(event watch.Event) interface{} {
	pod, ok := event.Object.(*corev1.Pod)
	resp := gin.H{
		"type": string(event.Type),
	}
	if ok {
		resp["object"] = models.ToPodResponse(pod)
	} else {
		if status, okStatus := event.Object.(*metav1.Status); okStatus {
			resp["error"] = fmt.Sprintf("K8s API Error: %s (Code: %d)", status.Message, status.Code)
			resp["status"] = status
		} else {
			resp["error"] = "事件对象类型不是 Pod 或 Status"
			resp["rawObject"] = fmt.Sprintf("%T", event.Object) // Show type
		}
	}
	return resp
}

// respondSuccess ... (保持不变)
func respondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"code":    code,
		"data":    data,
		"message": "success",
	})
}

// respondError ... (保持不变)
func respondError(c *gin.Context, code int, message string) {
	fmt.Printf("API Error: Status=%d, Message=%s\n", code, message)
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
