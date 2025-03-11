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

// SecretHandler ...
type SecretHandler struct {
	service *service.SecretService
}

// NewSecretHandler ...
func NewSecretHandler(svc *service.SecretService) *SecretHandler {
	return &SecretHandler{service: svc}
}

// ListSecrets ...
func (h *SecretHandler) ListSecrets(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层获取Secret列表
	secrets, err := h.service.List(namespace, c.Query("selector"), 0)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Secret列表失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, secrets)
}

// CreateSecret ...
func (h *SecretHandler) CreateSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	var req models.CreateSecretRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Secret格式: "+err.Error())
		return
	}

	// 2. 调用服务层创建Secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Data: req.Data,
		Type: req.Type,
	}

	createdSecret, err := h.service.Create(namespace, secret)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "创建Secret失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToSecretResponse(createdSecret))
}

// GetSecret ...
func (h *SecretHandler) GetSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Secret名称格式")
		return
	}

	// 2. 调用服务层获取Secret详情
	secret, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Secret不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Secret失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToSecretResponse(secret))
}

// UpdateSecret ...
func (h *SecretHandler) UpdateSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	var req models.UpdateSecretRequest

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Secret名称格式")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "无效的Secret格式: "+err.Error())
		return
	}

	// 2. 调用服务层更新Secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Data: req.Data,
		Type: req.Type,
	}

	updatedSecret, err := h.service.Update(namespace, secret)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "更新Secret失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, models.ToSecretResponse(updatedSecret))
}

// DeleteSecret ...
func (h *SecretHandler) DeleteSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Secret名称格式")
		return
	}

	// 2. 调用服务层删除Secret
	if err := h.service.Delete(namespace, name); err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Secret不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Secret失败: "+err.Error())
		return
	}

	// 3. 返回结果
	respondSuccess(c, http.StatusOK, gin.H{"message": "删除成功"})
}

// WatchSecrets ...
func (h *SecretHandler) WatchSecrets(c *gin.Context) {
	namespace := c.Param("namespace")
	// 1. 参数校验
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	// 2. 调用服务层Watch Secrets
	watcher, err := h.service.Watch(namespace, c.Query("selector"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Watch Secrets失败: "+err.Error())
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
