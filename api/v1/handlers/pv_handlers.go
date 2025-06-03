package handlers

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	"net/http"
	"strings"

	// "github.com/ciliverse/cilikube/api/v1/models" // Assuming models are defined here
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils" // Assuming utils package exists
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Might need if using req models
)

// --- Define Response Models (Example - Ideally in models package) ---

type PVResponse struct {
	Name            string                               `json:"name"`
	UID             string                               `json:"uid"`
	Capacity        string                               `json:"capacity"` // Keep as string "5Gi", "100Mi" etc.
	AccessModes     []corev1.PersistentVolumeAccessMode  `json:"accessModes"`
	ReclaimPolicy   corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	Status          string                               `json:"status"` // Phase: Pending, Available, Bound, Released, Failed
	Claim           string                               `json:"claim"`  // Format as "namespace/pvcName" if bound
	StorageClass    string                               `json:"storageClass"`
	VolumeMode      string                               `json:"volumeMode"`       // Filesystem or Block
	Reason          string                               `json:"reason,omitempty"` // Reason for Failed/Pending status
	CreatedAt       string                               `json:"createdAt"`
	Labels          map[string]string                    `json:"labels,omitempty"`
	Annotations     map[string]string                    `json:"annotations,omitempty"`
	ResourceVersion string                               `json:"resourceVersion"` // Needed for updates
}

type PVListResponse struct {
	Items []PVResponse `json:"items"`
	Total int          `json:"total"` // Total matching items
}

// Mapping function (Ideally in models package)
func ToPVResponse(pv *corev1.PersistentVolume) PVResponse {
	capacityStr := ""
	if storage, ok := pv.Spec.Capacity[corev1.ResourceStorage]; ok {
		capacityStr = storage.String()
	}

	claimStr := ""
	if pv.Spec.ClaimRef != nil {
		claimStr = pv.Spec.ClaimRef.Namespace + "/" + pv.Spec.ClaimRef.Name
	}

	volumeModeStr := string(corev1.PersistentVolumeFilesystem) // Default
	if pv.Spec.VolumeMode != nil {
		volumeModeStr = string(*pv.Spec.VolumeMode)
	}

	return PVResponse{
		Name:            pv.Name,
		UID:             string(pv.UID),
		Capacity:        capacityStr,
		AccessModes:     pv.Spec.AccessModes,
		ReclaimPolicy:   pv.Spec.PersistentVolumeReclaimPolicy,
		Status:          string(pv.Status.Phase),
		Claim:           claimStr,
		StorageClass:    pv.Spec.StorageClassName,
		VolumeMode:      volumeModeStr,
		Reason:          pv.Status.Reason,
		CreatedAt:       pv.CreationTimestamp.Format("2006-01-02T15:04:05Z"), // Use standard ISO 8601 for consistency
		Labels:          pv.Labels,
		Annotations:     pv.Annotations,
		ResourceVersion: pv.ResourceVersion,
	}
}

// --- Handler ---

type PVHandler struct {
	service        *service.PVService
	clusterManager *k8s.ClusterManager
}

func NewPVHandler(svc *service.PVService, cm *k8s.ClusterManager) *PVHandler {
	return &PVHandler{service: svc, clusterManager: cm}
}

// ListPVs godoc
// @Summary List Persistent Volumes
// @Description Get a list of all Persistent Volumes in the cluster
// @Tags PersistentVolumes
// @Accept json
// @Produce json
// @Param labelSelector query string false "Label selector for filtering"
// @Param limit query int false "Maximum number of items to return" default(500)
// @Success 200 {object} PVListResponse "List of Persistent Volumes"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/persistentvolumes [get]
func (h *PVHandler) ListPVs(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	labelSelector := c.Query("labelSelector")
	// Limit query param, default to a reasonable number, e.g., 500
	// Frontend pagination will handle displaying pageSize items from this list.
	limit := utils.ParseInt(c.DefaultQuery("limit", "500"), 500)

	pvList, err := h.service.List(k8sClient.Clientset, labelSelector, int64(limit))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取PV列表失败: "+err.Error())
		return
	}

	response := PVListResponse{
		Items: make([]PVResponse, 0, len(pvList.Items)),
		// Note: The 'total' here reflects the number *returned* by the List call (due to limit),
		//       not necessarily the *absolute* total in the cluster unless limit wasn't hit.
		//       For accurate total with server-side pagination, list without limit first (expensive)
		//       or rely on frontend count if limit is high enough for most cases.
		Total: len(pvList.Items),
	}
	for _, pv := range pvList.Items {
		response.Items = append(response.Items, ToPVResponse(&pv))
	}

	respondSuccess(c, http.StatusOK, response)
}

// GetPV godoc
// @Summary Get a Persistent Volume
// @Description Get details of a specific Persistent Volume by name
// @Tags PersistentVolumes
// @Accept json
// @Produce json
// @Param name path string true "Persistent Volume Name"
// @Success 200 {object} PVResponse "Persistent Volume details"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Name"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/persistentvolumes/{name} [get]
func (h *PVHandler) GetPV(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateResourceName(name) { // Assuming you have this validation
		respondError(c, http.StatusBadRequest, "无效的PV名称格式")
		return
	}

	pv, err := h.service.Get(k8sClient.Clientset, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "PV不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取PV失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, ToPVResponse(pv))
}

// CreatePV godoc
// @Summary Create a Persistent Volume
// @Description Create a new Persistent Volume using YAML/JSON definition
// @Tags PersistentVolumes
// @Accept json,yaml
// @Produce json
// @Param pv body corev1.PersistentVolume true "Persistent Volume definition (ensure kind: PersistentVolume, apiVersion: v1)"
// @Success 201 {object} PVResponse "Persistent Volume created"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Input"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Already Exists"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/persistentvolumes [post]
func (h *PVHandler) CreatePV(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	var pv corev1.PersistentVolume
	// Bind JSON or YAML depending on Content-Type? Gin might handle JSON by default.
	// For YAML, you might need custom binding or check Content-Type.
	if err := c.ShouldBindJSON(&pv); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求体格式: "+err.Error())
		return
	}

	// Ensure kind and apiVersion are correct (optional but good practice)
	if pv.Kind != "PersistentVolume" || (pv.APIVersion != "v1" && pv.APIVersion != "") {
		respondError(c, http.StatusBadRequest, "无效的Kind或APIVersion，应为 PersistentVolume/v1")
		return
	}
	if pv.APIVersion == "" {
		pv.APIVersion = "v1"
	} // Default if missing

	createdPV, err := h.service.Create(k8sClient.Clientset, &pv)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			respondError(c, http.StatusConflict, "PV已存在")
			return
		}
		// Handle validation errors from service
		if _, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, "创建PV验证失败: "+err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "创建PV失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusCreated, ToPVResponse(createdPV))
}

// UpdatePV godoc
// @Summary Update a Persistent Volume
// @Description Update an existing Persistent Volume using YAML/JSON definition
// @Tags PersistentVolumes
// @Accept json,yaml
// @Produce json
// @Param name path string true "Persistent Volume Name"
// @Param pv body corev1.PersistentVolume true "Updated Persistent Volume definition"
// @Success 200 {object} PVResponse "Persistent Volume updated"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Input or Name Mismatch"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Resource modified"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/persistentvolumes/{name} [put]
func (h *PVHandler) UpdatePV(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的PV名称格式")
		return
	}

	var pv corev1.PersistentVolume
	if err := c.ShouldBindJSON(&pv); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求体格式: "+err.Error())
		return
	}

	// Ensure name in path matches name in body
	if pv.Name != name {
		respondError(c, http.StatusBadRequest, "路径中的名称与请求体中的名称不匹配")
		return
	}
	// Ensure kind/apiVersion if necessary
	if pv.Kind != "PersistentVolume" || (pv.APIVersion != "v1" && pv.APIVersion != "") {
		respondError(c, http.StatusBadRequest, "无效的Kind或APIVersion，应为 PersistentVolume/v1")
		return
	}
	if pv.APIVersion == "" {
		pv.APIVersion = "v1"
	}

	updatedPV, err := h.service.Update(k8sClient.Clientset, &pv) // Service needs to handle potential conflicts
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "PV不存在")
			return
		}
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "资源已被修改，请获取最新版本后重试")
			return
		}
		if _, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, "更新PV验证失败: "+err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "更新PV失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, ToPVResponse(updatedPV))
}

// DeletePV godoc
// @Summary Delete a Persistent Volume
// @Description Delete a specific Persistent Volume by name
// @Tags PersistentVolumes
// @Accept json
// @Produce json
// @Param name path string true "Persistent Volume Name"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Name"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/persistentvolumes/{name} [delete]
func (h *PVHandler) DeletePV(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的PV名称格式")
		return
	}

	err := h.service.Delete(k8sClient.Clientset, name)
	if err != nil {
		if errors.IsNotFound(err) {
			// Consider returning 204 even if not found, idempotent delete
			// respondError(c, http.StatusNotFound, "PV不存在")
			c.Status(http.StatusNoContent) // Or 200 OK
			return
		}
		respondError(c, http.StatusInternalServerError, "删除PV失败: "+err.Error())
		return
	}
	c.Status(http.StatusNoContent) // Standard success for DELETE
}

// --- Re-use or define respond helpers ---
/*
type ErrorResponse struct { Code int `json:"code"`; Message string `json:"message"`}
func respondSuccess(c *gin.Context, code int, data interface{}) { ... }
func respondError(c *gin.Context, code int, message string) { ... }
*/
