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

// PVHandler ...
type PVHandler struct {
	service *service.PVService
}

// NewPVHandler ...
func NewPVHandler(svc *service.PVService) *PVHandler {
	return &PVHandler{service: svc}
}

// ListPVs ...
func (h *PVHandler) ListPVs(c *gin.Context) {
	// 1. 调用服务层获取PV列表
	pvs, err := h.service.List(c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取PV列表失败: "+err.Error())
		return
	}

	// 2. 返回结果
	respondSuccess(c, http.StatusOK, pvs)
}

// CreatePV ...
func (h *PVHandler) CreatePV(c *gin.Context) {
	var req models.CreatePVRequest

	// 1. 参数校验
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的PV格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建PV
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdPV, err := h.service.Create(pv)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建PV失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToPVResponse(createdPV))
}

// GetPV ...
func (h *PVHandler) GetPV(c *gin.Context) {
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的PV名称格式")
		return
	}

	// 2. 调用服务层获取PV详情
	pv, err := h.service.Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "PV不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取PV失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToPVResponse(pv))
}

// UpdatePV ...
func (h *PVHandler) UpdatePV(c *gin.Context) {
	name := c.Param("name")
	var req models.UpdatePVRequest

	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的PV名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的PV格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新PV
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	updatedPV, err := h.service.Update(pv)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新PV失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToPVResponse(updatedPV))
}

// DeletePV ...
func (h *PVHandler) DeletePV(c *gin.Context) {
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的PV名称格式")
		return
	}

	// 2. 调用服务层删除PV
	if err := h.service.Delete(name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "PV不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除PV失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchPVs ...
func (h *PVHandler) WatchPVs(c *gin.Context) {
	// 1. 调用服务层Watch PVs
	watcher, err := h.service.Watch(c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch PVs失败: "+err.Error())
		return
	}

	// 2. 返回结果
	c.Stream(func(w io.Writer) bool {
		event, ok := <-watcher.ResultChan()
		if !ok {
			return false
		}
		c.SSEvent("message", event)
		return true
	})
}
