package handlers

import (
	"net/http"
	"strconv"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceHandler 通用处理器
type ResourceHandler[T runtime.Object] struct {
	service        service.ResourceService[T]
	clusterManager *k8s.ClusterManager
	resourceType   string
}

// NewResourceHandler 创建通用处理器
func NewResourceHandler[T runtime.Object](svc service.ResourceService[T], k8sManager *k8s.ClusterManager, resourceType string) *ResourceHandler[T] {
	return &ResourceHandler[T]{
		service:        svc,
		clusterManager: k8sManager,
		resourceType:   resourceType,
	}
}

// List 处理列表请求
func (h *ResourceHandler[T]) List(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return // 错误已在 GetClientFromQuery 中处理
	}

	// 对于命名空间资源，从路径获取；对于集群资源，此参数为空
	namespace := c.Param("namespace")
	selector := c.Query("labelSelector")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 64)
	continueToken := c.Query("continue")

	items, err := h.service.List(k8sClient.Clientset, namespace, selector, limit, continueToken)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "获取资源列表失败", err.Error())
		return
	}

	utils.ApiSuccess(c, items, "成功获取资源列表")
}

// Get 处理获取单个资源请求
func (h *ResourceHandler[T]) Get(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	namespace := c.Param("namespace")
	name := c.Param("name")

	item, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "获取资源失败", err.Error())
		return
	}
	utils.ApiSuccess(c, item, "成功获取资源")
}

// Create 处理创建资源请求
func (h *ResourceHandler[T]) Create(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	namespace := c.Param("namespace")

	var obj T
	// Kubernetes 的 Create API 需要一个完整的对象，所以我们从请求体中绑定
	if err := c.ShouldBindJSON(&obj); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "请求体格式无效", err.Error())
		return
	}

	created, err := h.service.Create(k8sClient.Clientset, namespace, obj)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "创建资源失败", err.Error())
		return
	}
	utils.ApiSuccess(c, created, "资源创建成功")
}

// Update 处理更新资源请求
func (h *ResourceHandler[T]) Update(c *gin.Context) {
	utils.ApiError(c, http.StatusNotImplemented, "Update尚未实现", "")
}

// Delete 处理删除资源请求
func (h *ResourceHandler[T]) Delete(c *gin.Context) {
	utils.ApiError(c, http.StatusNotImplemented, "Delete尚未实现", "")
}

// Watch 处理监听资源请求
func (h *ResourceHandler[T]) Watch(c *gin.Context) {
	utils.ApiError(c, http.StatusNotImplemented, "Watch尚未实现", "")
}
