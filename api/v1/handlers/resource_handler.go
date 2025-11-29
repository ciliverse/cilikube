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

// ResourceHandler generic handler
type ResourceHandler[T runtime.Object] struct {
	service        service.ResourceService[T]
	clusterManager *k8s.ClusterManager
	resourceType   string
}

// NewResourceHandler creates generic handler
func NewResourceHandler[T runtime.Object](svc service.ResourceService[T], k8sManager *k8s.ClusterManager, resourceType string) *ResourceHandler[T] {
	return &ResourceHandler[T]{
		service:        svc,
		clusterManager: k8sManager,
		resourceType:   resourceType,
	}
}

// List handles list requests
func (h *ResourceHandler[T]) List(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return // Error already handled in GetClientFromQuery
	}

	// For namespaced resources, get from path; for cluster resources, this parameter is empty
	namespace := c.Param("namespace")
	selector := c.Query("labelSelector")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 64)
	continueToken := c.Query("continue")

	items, err := h.service.List(k8sClient.Clientset, namespace, selector, limit, continueToken)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to get resource list", err.Error())
		return
	}

	utils.ApiSuccess(c, items, "successfully retrieved resource list")
}

// Get handles single resource retrieval requests
func (h *ResourceHandler[T]) Get(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	namespace := c.Param("namespace")
	name := c.Param("name")

	item, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to get resource", err.Error())
		return
	}
	utils.ApiSuccess(c, item, "successfully retrieved resource")
}

// Create handles resource creation requests
func (h *ResourceHandler[T]) Create(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	namespace := c.Param("namespace")

	var obj T
	// Kubernetes Create API requires a complete object, so we bind from request body
	if err := c.ShouldBindJSON(&obj); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "invalid request body format", err.Error())
		return
	}

	created, err := h.service.Create(k8sClient.Clientset, namespace, obj)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to create resource", err.Error())
		return
	}
	utils.ApiSuccess(c, created, "resource created successfully")
}

// Update handles resource update requests
func (h *ResourceHandler[T]) Update(c *gin.Context) {
	utils.ApiError(c, http.StatusNotImplemented, "Update not yet implemented", "")
}

// Delete handles resource deletion requests
func (h *ResourceHandler[T]) Delete(c *gin.Context) {
	utils.ApiError(c, http.StatusNotImplemented, "Delete not yet implemented", "")
}

// Watch handles resource watch requests
func (h *ResourceHandler[T]) Watch(c *gin.Context) {
	utils.ApiError(c, http.StatusNotImplemented, "Watch not yet implemented", "")
}
