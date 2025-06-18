package handlers

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
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
func NewResourceHandler[T runtime.Object](service service.ResourceService[T], k8sManager *k8s.ClusterManager, resourceType string) *ResourceHandler[T] {
	return &ResourceHandler[T]{
		service:        service,
		clusterManager: k8sManager,
		resourceType:   resourceType,
	}
}

// List 处理列表请求
func (h *ResourceHandler[T]) List(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	selector := c.Query("labelSelector")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 64)
	continueToken := c.Query("continue")

	// 设置命名空间到上下文
	ctx := context.WithValue(c.Request.Context(), "namespace", namespace)
	c.Request = c.Request.WithContext(ctx)

	items, err := h.service.List(k8sClient.Clientset, selector, limit, continueToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// Get 处理获取单个资源请求
func (h *ResourceHandler[T]) Get(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	name := c.Param("name")
	item, err := h.service.Get(k8sClient.Clientset, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Create 处理创建资源请求
func (h *ResourceHandler[T]) Create(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	var obj T
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.Create(k8sClient.Clientset, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// Update 处理更新资源请求
func (h *ResourceHandler[T]) Update(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	var obj T
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.service.Update(k8sClient.Clientset, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Delete 处理删除资源请求
func (h *ResourceHandler[T]) Delete(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	name := c.Param("name")
	if err := h.service.Delete(k8sClient.Clientset, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Watch 处理监听资源请求
func (h *ResourceHandler[T]) Watch(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	selector := c.Query("labelSelector")
	resourceVersion := c.Query("resourceVersion")
	timeoutSeconds, _ := strconv.ParseInt(c.DefaultQuery("timeoutSeconds", "0"), 10, 64)

	watcher, err := h.service.Watch(k8sClient.Clientset, selector, resourceVersion, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
