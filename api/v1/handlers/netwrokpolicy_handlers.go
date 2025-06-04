package handlers

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
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

// NetworkPolicyHandler ...
type NetworkPolicyHandler struct {
	service        *service.NetworkPolicyService
	clusterManager *k8s.ClusterManager
}

// NewNetworkPolicyHandler ...
func NewNetworkPolicyHandler(svc *service.NetworkPolicyService, cm *k8s.ClusterManager) *NetworkPolicyHandler {
	return &NetworkPolicyHandler{service: svc, clusterManager: cm}
}

// ListNetworkPolicies ...
func (h *NetworkPolicyHandler) ListNetworkPolicies(c *gin.Context) {
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

	// 2. 调用服务层获取NetworkPolicy列表
	networkPolicies, err := h.service.List(k8sClient.Clientset, namespace, c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取NetworkPolicy列表失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, networkPolicies)
}

// CreateNetworkPolicy ...
func (h *NetworkPolicyHandler) CreateNetworkPolicy(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	var req models.CreateNetworkPolicyRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的NetworkPolicy格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建NetworkPolicy
	networkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	createdNetworkPolicy, err := h.service.Create(k8sClient.Clientset, namespace, networkPolicy)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建NetworkPolicy失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNetworkPolicyResponse(createdNetworkPolicy))
}

// GetNetworkPolicy ...
func (h *NetworkPolicyHandler) GetNetworkPolicy(c *gin.Context) {
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
		respondError(c, http.StatusBadRequest, "无效的NetworkPolicy名称格式")
		return
	}

	// 2. 调用服务层获取NetworkPolicy详情
	networkPolicy, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "NetworkPolicy不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取NetworkPolicy失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNetworkPolicyResponse(networkPolicy))
}

// UpdateNetworkPolicy ...
func (h *NetworkPolicyHandler) UpdateNetworkPolicy(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	var req models.UpdateNetworkPolicyRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的NetworkPolicy名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的NetworkPolicy格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新NetworkPolicy
	networkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Spec: req.Spec,
	}

	updatedNetworkPolicy, err := h.service.Update(k8sClient.Clientset, namespace, networkPolicy)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新NetworkPolicy失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToNetworkPolicyResponse(updatedNetworkPolicy))
}

// DeleteNetworkPolicy ...
func (h *NetworkPolicyHandler) DeleteNetworkPolicy(c *gin.Context) {
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
		respondError(c, http.StatusBadRequest, "无效的NetworkPolicy名称格式")
		return
	}

	// 2. 调用服务层删除NetworkPolicy
	if err := h.service.Delete(k8sClient.Clientset, namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "NetworkPolicy不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除NetworkPolicy失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchNetworkPolicies ...
func (h *NetworkPolicyHandler) WatchNetworkPolicies(c *gin.Context) {
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

	// 2. 调用服务层Watch NetworkPolicies
	watcher, err := h.service.Watch(k8sClient.Clientset, namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch NetworkPolicies失败: "+err.Error())
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
