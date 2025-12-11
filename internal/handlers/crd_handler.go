package handlers

import (
	"net/http"
	"strconv"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

// CRDHandler handles CRD operations
type CRDHandler struct {
	crdService     service.CRDService
	clusterManager *k8s.ClusterManager
}

// NewCRDHandler creates a new CRD handler
func NewCRDHandler(crdService service.CRDService, clusterManager *k8s.ClusterManager) *CRDHandler {
	return &CRDHandler{
		crdService:     crdService,
		clusterManager: clusterManager,
	}
}

// ListCRDs retrieves the list of CRDs
func (h *CRDHandler) ListCRDs(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	crds, err := h.crdService.ListCRDs(k8sClient)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to get CRD list", err.Error())
		return
	}

	utils.ApiSuccess(c, crds, "successfully retrieved CRD list")
}

// GetCRD retrieves CRD details
func (h *CRDHandler) GetCRD(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	name := c.Param("name")
	if name == "" {
		utils.ApiError(c, http.StatusBadRequest, "CRD name is required", "")
		return
	}

	crd, err := h.crdService.GetCRD(k8sClient, name)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to get CRD", err.Error())
		return
	}

	utils.ApiSuccess(c, crd, "successfully retrieved CRD")
}

// ListCustomResources retrieves the list of custom resources
func (h *CRDHandler) ListCustomResources(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	group := c.Param("group")
	version := c.Param("version")
	plural := c.Param("plural")
	namespace := c.Query("namespace")

	if group == "" || version == "" || plural == "" {
		utils.ApiError(c, http.StatusBadRequest, "group, version and plural are required", "")
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "0"), 10, 64)
	continueToken := c.Query("continue")

	resources, err := h.crdService.ListCustomResources(k8sClient, group, version, plural, namespace, limit, continueToken)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to get custom resource list", err.Error())
		return
	}

	utils.ApiSuccess(c, resources, "successfully retrieved custom resource list")
}

// GetCustomResource retrieves custom resource details
func (h *CRDHandler) GetCustomResource(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	group := c.Param("group")
	version := c.Param("version")
	plural := c.Param("plural")
	namespace := c.Query("namespace")
	name := c.Param("name")

	if group == "" || version == "" || plural == "" || name == "" {
		utils.ApiError(c, http.StatusBadRequest, "group, version, plural and name are required", "")
		return
	}

	resource, err := h.crdService.GetCustomResource(k8sClient, group, version, plural, namespace, name)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to get custom resource", err.Error())
		return
	}

	utils.ApiSuccess(c, resource, "successfully retrieved custom resource")
}

// CreateCustomResource creates a custom resource
func (h *CRDHandler) CreateCustomResource(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	group := c.Param("group")
	version := c.Param("version")
	plural := c.Param("plural")
	namespace := c.Query("namespace")

	if group == "" || version == "" || plural == "" {
		utils.ApiError(c, http.StatusBadRequest, "group, version and plural are required", "")
		return
	}

	var req models.CustomResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	resource, err := h.crdService.CreateCustomResource(k8sClient, group, version, plural, namespace, &req)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to create custom resource", err.Error())
		return
	}

	utils.ApiSuccess(c, resource, "custom resource created successfully")
}

// UpdateCustomResource updates a custom resource
func (h *CRDHandler) UpdateCustomResource(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	group := c.Param("group")
	version := c.Param("version")
	plural := c.Param("plural")
	namespace := c.Query("namespace")
	name := c.Param("name")

	if group == "" || version == "" || plural == "" || name == "" {
		utils.ApiError(c, http.StatusBadRequest, "group, version, plural and name are required", "")
		return
	}

	var req models.CustomResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ApiError(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	resource, err := h.crdService.UpdateCustomResource(k8sClient, group, version, plural, namespace, name, &req)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to update custom resource", err.Error())
		return
	}

	utils.ApiSuccess(c, resource, "custom resource updated successfully")
}

// DeleteCustomResource deletes a custom resource
func (h *CRDHandler) DeleteCustomResource(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	group := c.Param("group")
	version := c.Param("version")
	plural := c.Param("plural")
	namespace := c.Query("namespace")
	name := c.Param("name")

	if group == "" || version == "" || plural == "" || name == "" {
		utils.ApiError(c, http.StatusBadRequest, "group, version, plural and name are required", "")
		return
	}

	err := h.crdService.DeleteCustomResource(k8sClient, group, version, plural, namespace, name)
	if err != nil {
		utils.ApiError(c, http.StatusInternalServerError, "failed to delete custom resource", err.Error())
		return
	}

	utils.ApiSuccess(c, nil, "custom resource deleted successfully")
}
