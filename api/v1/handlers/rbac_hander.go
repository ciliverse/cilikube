package handlers

import (
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RbacHandler struct {
	service        *service.RbacService
	clusterManager *k8s.ClusterManager
}

func NewRbacHandler(svc *service.RbacService, cm *k8s.ClusterManager) *RbacHandler {
	return &RbacHandler{service: svc, clusterManager: cm}
}

// Roles
func (h *RbacHandler) ListRoles(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}
	roles, err := h.service.ListRoles(k8sClient.Clientset, namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Role列表失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, roles)
}

func (h *RbacHandler) GetRole(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的资源名称格式")
		return
	}
	role, err := h.service.GetRole(k8sClient.Clientset, namespace, name)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Role失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, role)
}

// ListRoleBinding
func (h *RbacHandler) ListRoleBindings(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}
	roleBindings, err := h.service.ListRoleBindings(k8sClient.Clientset, namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取RoleBinding列表失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, roleBindings)
}

func (h *RbacHandler) GetRoleBindings(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的资源名称格式")
		return
	}
	roleBinding, err := h.service.GetRoleBinding(k8sClient.Clientset, namespace, name)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取RoleBinding失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, roleBinding)
}

// ClusterRoles
func (h *RbacHandler) ListClusterRoles(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	clusterRoles, err := h.service.ListClusterRoles(k8sClient.Clientset)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取ClusterRole列表失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, clusterRoles)
}

func (h *RbacHandler) GetClusterRoles(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的资源名称格式")
		return
	}
	clusterRole, err := h.service.GetClusterRole(k8sClient.Clientset, name)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取ClusterRole失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, clusterRole)
}

// ClusterRoleBindings
func (h *RbacHandler) ListClusterRoleBindings(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	clusterRoleBindings, err := h.service.ListClusterRoleBindings(k8sClient.Clientset)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取ClusterRoleBinding列表失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, clusterRoleBindings)
}

func (h *RbacHandler) GetClusterRoleBindings(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的资源名称格式")
		return
	}
	clusterRoleBinding, err := h.service.GetClusterRoleBinding(k8sClient.Clientset, name)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取ClusterRoleBinding失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, clusterRoleBinding)
}

// ServiceAccounts
func (h *RbacHandler) ListServiceAccounts(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}
	serviceAccounts, err := h.service.ListServiceAccounts(k8sClient.Clientset, namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取ServiceAccount列表失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, serviceAccounts)
}

func (h *RbacHandler) GetServiceAccounts(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的资源名称格式")
		return
	}
	serviceAccount, err := h.service.GetServiceAccounts(k8sClient.Clientset, namespace, name)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取ServiceAccount失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, serviceAccount)
}
