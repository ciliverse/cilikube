// package handlers

// import (
// 	"io"
// 	"net/http"

// 	"github.com/ciliverse/cilikube/api/v1/models"
// 	"github.com/ciliverse/cilikube/internal/service"
// 	"github.com/ciliverse/cilikube/pkg/utils"
// 	"github.com/gin-gonic/gin"
// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/apimachinery/pkg/api/errors"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// // PVCHandler ...
// type PVCHandler struct {
// 	service *service.PVCService
// }

// // NewPVCHandler ...
// func NewPVCHandler(svc *service.PVCService) *PVCHandler {
// 	return &PVCHandler{service: svc}
// }

// // ListPVCs ...
// func (h *PVCHandler) ListPVCs(c *gin.Context) {
// 	namespace := c.Param("namespace")
// 	// 1. 参数校验
// 	if !utils.ValidateNamespace(namespace) {
// 		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
// 		return
// 	}

// 	// 2. 调用服务层获取PVC列表
// 	pvcs, err := h.service.List(namespace, c.Query("selector"), 0)
// 	if err != nil {
// 		respondError(c, http.StatusInternalServerError, "获取PVC列表失败: "+err.Error())
// 		return
// 	}

// 	// 3. 返回结果
// 	respondSuccess(c, http.StatusOK, pvcs)
// }

// // CreatePVC ...
// func (h *PVCHandler) CreatePVC(c *gin.Context) {
// 	namespace := c.Param("namespace")
// 	var req models.CreatePVCRequest

// 	// 1. 参数校验
// 	if !utils.ValidateNamespace(namespace) {
// 		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		respondError(c, http.StatusBadRequest, "无效的PVC格式: "+err.Error())
// 		return
// 	}

// 	// 2. 调用服务层创建PVC
// 	pvc := &corev1.PersistentVolumeClaim{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:        req.Name,
// 			Namespace:   req.Namespace,
// 			Labels:      req.Labels,
// 			Annotations: req.Annotations,
// 		},
// 		Spec: req.Spec,
// 	}

// 	createdPVC, err := h.service.Create(namespace, pvc)
// 	if err != nil {
// 		respondError(c, http.StatusInternalServerError, "创建PVC失败: "+err.Error())
// 		return
// 	}

// 	// 3. 返回结果
// 	respondSuccess(c, http.StatusOK, models.ToPVCResponse(createdPVC))
// }

// // GetPVC ...
// func (h *PVCHandler) GetPVC(c *gin.Context) {
// 	namespace := c.Param("namespace")
// 	name := c.Param("name")
// 	// 1. 参数校验
// 	if !utils.ValidateNamespace(namespace) {
// 		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
// 		return
// 	}

// 	if !utils.ValidateResourceName(name) {
// 		respondError(c, http.StatusBadRequest, "无效的PVC名称格式")
// 		return
// 	}

// 	// 2. 调用服务层获取PVC详情
// 	pvc, err := h.service.Get(namespace, name)
// 	if err != nil {
// 		if errors.IsNotFound(err) {
// 			respondError(c, http.StatusNotFound, "PVC不存在")
// 			return
// 		}
// 		respondError(c, http.StatusInternalServerError, "获取PVC失败: "+err.Error())
// 		return
// 	}

// 	// 3. 返回结果
// 	respondSuccess(c, http.StatusOK, models.ToPVCResponse(pvc))
// }

// // UpdatePVC ...
// func (h *PVCHandler) UpdatePVC(c *gin.Context) {
// 	namespace := c.Param("namespace")
// 	name := c.Param("name")
// 	var req models.UpdatePVCRequest

// 	// 1. 参数校验
// 	if !utils.ValidateNamespace(namespace) {
// 		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
// 		return
// 	}

// 	if !utils.ValidateResourceName(name) {
// 		respondError(c, http.StatusBadRequest, "无效的PVC名称格式")
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		respondError(c, http.StatusBadRequest, "无效的PVC格式: "+err.Error())
// 		return
// 	}

// 	// 2. 调用服务层更新PVC
// 	pvc := &corev1.PersistentVolumeClaim{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:        name,
// 			Namespace:   namespace,
// 			Labels:      req.Labels,
// 			Annotations: req.Annotations,
// 		},
// 		Spec: req.Spec,
// 	}

// 	updatedPVC, err := h.service.Update(namespace, pvc)
// 	if err != nil {
// 		respondError(c, http.StatusInternalServerError, "更新PVC失败: "+err.Error())
// 		return
// 	}

// 	// 3. 返回结果
// 	respondSuccess(c, http.StatusOK, models.ToPVCResponse(updatedPVC))
// }

// // DeletePVC ...
// func (h *PVCHandler) DeletePVC(c *gin.Context) {
// 	namespace := c.Param("namespace")
// 	name := c.Param("name")

// 	// 1. 参数校验
// 	if !utils.ValidateNamespace(namespace) {
// 		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
// 		return
// 	}

// 	if !utils.ValidateResourceName(name) {
// 		respondError(c, http.StatusBadRequest, "无效的PVC名称格式")
// 		return
// 	}

// 	// 2. 调用服务层删除PVC
// 	if err := h.service.Delete(namespace, name); err != nil {
// 		if errors.IsNotFound(err) {
// 			respondError(c, http.StatusNotFound, "PVC不存在")
// 			return
// 		}
// 		respondError(c, http.StatusInternalServerError, "删除PVC失败: "+err.Error())
// 		return
// 	}

// 	// 3. 返回结果
// 	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
// }

// // WatchPVCs ...
// func (h *PVCHandler) WatchPVCs(c *gin.Context) {
// 	namespace := c.Param("namespace")
// 	// 1. 参数校验
// 	if !utils.ValidateNamespace(namespace) {
// 		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
// 		return
// 	}

// 	// 2. 调用服务层Watch PVCs
// 	watcher, err := h.service.Watch(namespace, c.Query("selector"))
// 	if err != nil {
// 		respondError(c, http.StatusInternalServerError, "Watch PVCs失败: "+err.Error())
// 		return
// 	}

//		// 3. 返回结果
//		c.Stream(func(w io.Writer) bool {
//			event, ok := <-watcher.ResultChan()
//			if !ok {
//				return false
//			}
//			c.SSEvent("message", event)
//			return true
//		})
//	}
package handlers

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	"net/http"
	"strings"

	"github.com/ciliverse/cilikube/api/v1/models" // Adjust import path
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils" // Assuming utils package exists
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type PVCHandler struct {
	service        *service.PVCService
	clusterManager *k8s.ClusterManager
}

func NewPVCHandler(svc *service.PVCService, cm *k8s.ClusterManager) *PVCHandler {
	return &PVCHandler{service: svc, clusterManager: cm}
}

// ListPVCs godoc
// @Summary List Persistent Volume Claims
// @Description Get a list of PVCs within a specific namespace
// @Tags PersistentVolumeClaims
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param labelSelector query string false "Label selector for filtering"
// @Param limit query int false "Maximum number of items to return" default(100)
// @Success 200 {object} models.PVCListResponse "List of Persistent Volume Claims"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Namespace"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/persistentvolumeclaims [get]
func (h *PVCHandler) ListPVCs(c *gin.Context) {
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

	pvcList, err := h.service.List(k8sClient.Clientset, namespace, labelSelector, int64(limit))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取PVC列表失败: "+err.Error())
		return
	}

	response := models.PVCListResponse{
		Items: make([]models.PVCResponse, 0, len(pvcList.Items)),
		Total: len(pvcList.Items), // See notes in PV handlers about total vs returned
	}
	for _, pvc := range pvcList.Items {
		response.Items = append(response.Items, models.ToPVCResponse(&pvc))
	}

	respondSuccess(c, http.StatusOK, response)
}

// GetPVC godoc
// @Summary Get a Persistent Volume Claim
// @Description Get details of a specific PVC by namespace and name
// @Tags PersistentVolumeClaims
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "PVC Name"
// @Success 200 {object} models.PVCResponse "Persistent Volume Claim details"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Name/Namespace"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name} [get]
func (h *PVCHandler) GetPVC(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或PVC名称格式")
		return
	}

	pvc, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "PVC不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取PVC失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, models.ToPVCResponse(pvc))
}

// CreatePVC godoc
// @Summary Create a Persistent Volume Claim
// @Description Create a new PVC using YAML/JSON definition
// @Tags PersistentVolumeClaims
// @Accept json,yaml
// @Produce json
// @Param namespace path string true "Namespace"
// @Param pvc body corev1.PersistentVolumeClaim true "PVC definition (ensure kind: PersistentVolumeClaim, apiVersion: v1)"
// @Success 201 {object} models.PVCResponse "PVC created"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Input"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Already Exists"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/persistentvolumeclaims [post]
func (h *PVCHandler) CreatePVC(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	var pvc corev1.PersistentVolumeClaim
	if err := c.ShouldBindJSON(&pvc); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求体格式: "+err.Error())
		return
	}

	// Validate Kind and APIVersion
	if pvc.Kind != "PersistentVolumeClaim" || (pvc.APIVersion != "v1" && pvc.APIVersion != "") {
		respondError(c, http.StatusBadRequest, "无效的Kind或APIVersion，应为 PersistentVolumeClaim/v1")
		return
	}
	if pvc.APIVersion == "" {
		pvc.APIVersion = "v1"
	}

	// Let service handle namespace assignment/validation based on path param
	createdPVC, err := h.service.Create(k8sClient.Clientset, namespace, &pvc)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			respondError(c, http.StatusConflict, "PVC已存在")
			return
		}
		if vErr, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, "创建PVC验证失败: "+vErr.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "创建PVC失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusCreated, models.ToPVCResponse(createdPVC))
}

// UpdatePVC godoc
// @Summary Update a Persistent Volume Claim
// @Description Update an existing PVC (primarily labels/annotations)
// @Tags PersistentVolumeClaims
// @Accept json,yaml
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "PVC Name"
// @Param pvc body corev1.PersistentVolumeClaim true "Updated PVC definition (Spec changes likely rejected by API server)"
// @Success 200 {object} models.PVCResponse "PVC updated"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Input or Name/Namespace Mismatch"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Resource modified"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name} [put]
func (h *PVCHandler) UpdatePVC(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或PVC名称格式")
		return
	}

	var pvc corev1.PersistentVolumeClaim
	if err := c.ShouldBindJSON(&pvc); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求体格式: "+err.Error())
		return
	}

	// Validate name/namespace consistency
	if pvc.Name != name || (pvc.Namespace != "" && pvc.Namespace != namespace) {
		respondError(c, http.StatusBadRequest, "路径参数与请求体中的名称/命名空间不匹配")
		return
	}
	// Validate Kind and APIVersion
	if pvc.Kind != "PersistentVolumeClaim" || (pvc.APIVersion != "v1" && pvc.APIVersion != "") {
		respondError(c, http.StatusBadRequest, "无效的Kind或APIVersion，应为 PersistentVolumeClaim/v1")
		return
	}
	if pvc.APIVersion == "" {
		pvc.APIVersion = "v1"
	}

	// Service Update handles the actual call, API server enforces immutability
	updatedPVC, err := h.service.Update(k8sClient.Clientset, namespace, &pvc)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "PVC不存在")
			return
		}
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "资源已被修改，请获取最新版本后重试")
			return
		}
		if vErr, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, "更新PVC验证失败: "+vErr.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "更新PVC失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, models.ToPVCResponse(updatedPVC))
}

// DeletePVC godoc
// @Summary Delete a Persistent Volume Claim
// @Description Delete a specific PVC by namespace and name
// @Tags PersistentVolumeClaims
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "PVC Name"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Name/Namespace"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/persistentvolumeclaims/{name} [delete]
func (h *PVCHandler) DeletePVC(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或PVC名称格式")
		return
	}

	err := h.service.Delete(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			// respondError(c, http.StatusNotFound, "PVC不存在")
			c.Status(http.StatusNoContent) // Idempotent delete
			return
		}
		respondError(c, http.StatusInternalServerError, "删除PVC失败: "+err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Re-use or define respond helpers ---
/* Defined in pv_handler or a shared place
type ErrorResponse struct { Code int `json:"code"`; Message string `json:"message"`}
func respondSuccess(c *gin.Context, code int, data interface{}) { ... }
func respondError(c *gin.Context, code int, message string) { ... }
*/
