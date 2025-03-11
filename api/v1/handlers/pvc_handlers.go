package handlers

import (
	"io"
	"net/http"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PVCHandler ...
type PVCHandler struct {
	service *service.PVCService
}

// NewPVCHandler ...
func NewPVCHandler(svc *service.PVCService) *PVCHandler {
	return &PVCHandler{service: svc}
}

// ListPVCs ...
func (h *PVCHandler) ListPVCs(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层获取PVC列表
	pvcs, err := h.service.List(namespace, c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取PVC列表失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, pvcs)
}

// CreatePVC ...
func (h *PVCHandler) CreatePVC(c *gin.Context) {
	namespace := c.Param("namespace")
	var req models.CreatePVCRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的PVC格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建PVC
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdPVC, err := h.service.Create(namespace, pvc)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建PVC失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToPVCResponse(createdPVC))
}

// GetPVC ...
func (h *PVCHandler) GetPVC(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的PVC名称格式")
		return
	}

	// 2. 调用服务层获取PVC详情
	pvc, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "PVC不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取PVC失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToPVCResponse(pvc))
}

// UpdatePVC ...
func (h *PVCHandler) UpdatePVC(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	var req models.UpdatePVCRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的PVC名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的PVC格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新PVC
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	updatedPVC, err := h.service.Update(namespace, pvc)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新PVC失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToPVCResponse(updatedPVC))
}

// DeletePVC ...
func (h *PVCHandler) DeletePVC(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的PVC名称格式")
		return
	}

	// 2. 调用服务层删除PVC
	if err := h.service.Delete(namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "PVC不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除PVC失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchPVCs ...
func (h *PVCHandler) WatchPVCs(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层Watch PVCs
	watcher, err := h.service.Watch(namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch PVCs失败: "+err.Error())
		return
	}

	// 3. 返回结果
	c.Stream(func(w io.Writer) bool {
		event, ok := <-watcher.ResultChan()
		if !ok {
			return false
		}
		c.SSEvent("message", event)
		return true
	})
}
