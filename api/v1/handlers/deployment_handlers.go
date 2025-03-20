package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
)

// DeploymentHandler ...
type DeploymentHandler struct {
	service *service.DeploymentService
}

// NewDeploymentHandler ...
func NewDeploymentHandler(svc *service.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{service: svc}
}

// ListDeployments ...
func (h *DeploymentHandler) ListDeployments(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层获取Deployment列表
	deployments, err := h.service.List(namespace)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Deployment列表失败: "+err.Error())
		return
	}
	// 3. 返回结果
	respondSuccess(c, http.StatusOK, deployments)

}

// CreateDeployment ...
func (h *DeploymentHandler) CreateDeployment(c *gin.Context) {}

// GetDeployment ...
func (h *DeploymentHandler) GetDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Deployment名称格式")
		return
	}

	// 2. 调用服务层获取Deployment详情
	deployment, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Deployment不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Deployment失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToDeploymentResponse(deployment))

}

// UpdateDeployment ...
func (h *DeploymentHandler) UpdateDeployment(c *gin.Context) {}

// DeleteDeployment ...
func (h *DeploymentHandler) DeleteDeployment(c *gin.Context) {}

// WatchDeployments ...
func (h *DeploymentHandler) WatchDeployments(c *gin.Context) {}
