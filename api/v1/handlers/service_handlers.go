package handlers

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
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

// ServiceHandler ...
type ServiceHandler struct {
	service        *service.ServiceService
	clusterManager *k8s.ClusterManager
}

// NewServiceHandler ...
func NewServiceHandler(svc *service.ServiceService, cm *k8s.ClusterManager) *ServiceHandler {
	return &ServiceHandler{service: svc, clusterManager: cm}
}

// ListServices ...
func (h *ServiceHandler) ListServices(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层获取Service列表
	services, err := h.service.List(k8sClient.Clientset, namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Service列表失败: "+err.Error())
		return
	}

	// 修复空列表报错问题
	if services.Items == nil {
		services.Items = make([]corev1.Service, 0)
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, services)
}

// CreateService ...
func (h *ServiceHandler) CreateService(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	var req models.CreateServiceRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Service格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdService, err := h.service.Create(k8sClient.Clientset, namespace, service)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建Service失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToServiceResponse(createdService))
}

// GetService ...
func (h *ServiceHandler) GetService(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Service名称格式")
		return
	}

	// 2. 调用服务层获取Service详情
	service, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Service不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Service失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToServiceResponse(service))
}

// UpdateService ...
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	var req models.UpdateServiceRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Service名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Service格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	updatedService, err := h.service.Update(k8sClient.Clientset, namespace, service)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新Service失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToServiceResponse(updatedService))
}

// DeleteService ...
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Service名称格式")
		return
	}

	// 2. 调用服务层删除Service
	if err := h.service.Delete(k8sClient.Clientset, namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Service不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Service失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchServices ...
func (h *ServiceHandler) WatchServices(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层Watch Services
	watcher, err := h.service.Watch(k8sClient.Clientset, namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch Services失败: "+err.Error())
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
