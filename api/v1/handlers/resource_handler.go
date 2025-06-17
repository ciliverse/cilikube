package handlers

import (
	"net/http"
	"strconv"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

// ResourceHandler 通用 handler
// T 为 K8s 资源类型（如 *corev1.Node）
type ResourceHandler[T any] struct {
	Service      service.ResourceService[T]
	GetK8sClient func(c *gin.Context) (kubernetes.Interface, bool)
}

func NewResourceHandler[T any](service service.ResourceService[T], k8sManager *k8s.ClusterManager) *ResourceHandler[T] {
	return &ResourceHandler[T]{
		Service: service,
		GetK8sClient: func(c *gin.Context) (kubernetes.Interface, bool) {
			client, ok := k8s.GetK8sClientFromContext(c, k8sManager)
			if !ok {
				return nil, false
			}
			return client.Clientset, true
		},
	}
}

func (h *ResourceHandler[T]) List(c *gin.Context) {
	client, ok := h.GetK8sClient(c)
	if !ok {
		return
	}
	selector := c.Query("labelSelector")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 64)
	continueToken := c.Query("continue")
	items, err := h.Service.List(client, selector, limit, continueToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ResourceHandler[T]) Get(c *gin.Context) {
	client, ok := h.GetK8sClient(c)
	if !ok {
		return
	}
	name := c.Param("name")
	item, err := h.Service.Get(client, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ResourceHandler[T]) Create(c *gin.Context) {
	client, ok := h.GetK8sClient(c)
	if !ok {
		return
	}
	var obj T
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := h.Service.Create(client, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *ResourceHandler[T]) Update(c *gin.Context) {
	client, ok := h.GetK8sClient(c)
	if !ok {
		return
	}
	var obj T
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.Service.Update(client, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *ResourceHandler[T]) Delete(c *gin.Context) {
	client, ok := h.GetK8sClient(c)
	if !ok {
		return
	}
	name := c.Param("name")
	if err := h.Service.Delete(client, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
