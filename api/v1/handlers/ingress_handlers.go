package handlers

import (
	"io"
	"net/http"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressHandler ...
type IngressHandler struct {
	service *service.IngressService
}

// NewIngressHandler ...
func NewIngressHandler(svc *service.IngressService) *IngressHandler {
	return &IngressHandler{service: svc}
}

// ListIngresses ...
func (h *IngressHandler) ListIngresses(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层获取Ingress列表
	ingresses, err := h.service.List(namespace, c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Ingress列表失败: "+err.Error())
		return
	}

	// 修复空列表报错问题
	if ingresses.Items == nil {
		ingresses.Items = make([]networkingv1.Ingress, 0)
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, ingresses)
}

// CreateIngress ...
func (h *IngressHandler) CreateIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	var req models.CreateIngressRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Ingress格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建Ingress
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdIngress, err := h.service.Create(namespace, ingress)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建Ingress失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToIngressResponse(createdIngress))
}

// GetIngress ...
func (h *IngressHandler) GetIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Ingress名称格式")
		return
	}

	// 2. 调用服务层获取Ingress详情
	ingress, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Ingress不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Ingress失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToIngressResponse(ingress))
}

// UpdateIngress ...
func (h *IngressHandler) UpdateIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	var req models.UpdateIngressRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Ingress名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Ingress格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新Ingress
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	updatedIngress, err := h.service.Update(namespace, ingress)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新Ingress失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToIngressResponse(updatedIngress))
}

// DeleteIngress ...
func (h *IngressHandler) DeleteIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Ingress名称格式")
		return
	}

	// 2. 调用服务层删除Ingress
	if err := h.service.Delete(namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Ingress不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Ingress失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchIngresses ...
func (h *IngressHandler) WatchIngresses(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层Watch Ingresses
	watcher, err := h.service.Watch(namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch Ingresses失败: "+err.Error())
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
