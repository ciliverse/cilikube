package handlers

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	"net/http"
	"strings"

	"github.com/ciliverse/cilikube/api/v1/models" // Adjust path
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type SecretHandler struct {
	service        *service.SecretService
	clusterManager *k8s.ClusterManager
}

func NewSecretHandler(svc *service.SecretService, cm *k8s.ClusterManager) *SecretHandler {
	return &SecretHandler{service: svc, clusterManager: cm}
}

// ListSecrets godoc
// @Summary List Secrets
// @Description Get a list of Secrets in a specific namespace (data values are omitted).
// @Tags Secrets
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param labelSelector query string false "Label selector for filtering"
// @Param limit query int false "Maximum number of items to return" default(100)
// @Success 200 {object} models.SecretListResponse "List of Secrets (metadata only)"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Namespace"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/secrets [get]
func (h *SecretHandler) ListSecrets(c *gin.Context) {
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

	secretList, err := h.service.List(k8sClient.Clientset, namespace, labelSelector, int64(limit))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Secret列表失败: "+err.Error())
		return
	}

	response := models.SecretListResponse{
		Items: make([]models.SecretResponse, 0, len(secretList.Items)),
		Total: len(secretList.Items), // Adjust if needed
	}
	for _, secret := range secretList.Items {
		response.Items = append(response.Items, models.ToSecretResponse(&secret))
	}

	respondSuccess(c, http.StatusOK, response)
}

// GetSecret godoc
// @Summary Get a Secret
// @Description Get details of a specific Secret by namespace and name, including data keys (values base64 encoded).
// @Tags Secrets
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "Secret Name"
// @Success 200 {object} models.SecretDetailResponse "Secret details"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Name/Namespace"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/secrets/{name} [get]
func (h *SecretHandler) GetSecret(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或Secret名称格式")
		return
	}

	secret, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Secret不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Secret失败: "+err.Error())
		return
	}
	// Use detail response model including Data (base64 encoded) and StringData
	respondSuccess(c, http.StatusOK, models.ToSecretDetailResponse(secret))
}

// CreateSecret godoc
// @Summary Create a Secret
// @Description Create a new Secret using YAML/JSON definition. Data should be base64 encoded, StringData plain text.
// @Tags Secrets
// @Accept json,yaml
// @Produce json
// @Param namespace path string true "Namespace"
// @Param secret body corev1.Secret true "Secret definition (ensure kind: Secret, apiVersion: v1)"
// @Success 201 {object} models.SecretResponse "Secret created (metadata only)"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Input"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Already Exists"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/secrets [post]
func (h *SecretHandler) CreateSecret(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	var secret corev1.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求体格式: "+err.Error())
		return
	}

	if secret.Kind != "Secret" || (secret.APIVersion != "v1" && secret.APIVersion != "") {
		respondError(c, http.StatusBadRequest, "无效的Kind或APIVersion，应为 Secret/v1")
		return
	}
	if secret.APIVersion == "" {
		secret.APIVersion = "v1"
	}

	// Note: If frontend sends base64 encoded strings in 'Data', they need to be decoded before passing to K8s client?
	// The Go client expects map[string][]byte for secret.Data. If ShouldBindJSON correctly handles base64 -> []byte, it's fine.
	// If not, you might need manual decoding here based on how the frontend sends it.
	// However, K8s usually handles encoding StringData into Data automatically. Prefer using StringData for text.

	createdSecret, err := h.service.Create(k8sClient.Clientset, namespace, &secret)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			respondError(c, http.StatusConflict, "Secret已存在")
			return
		}
		if _, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, "创建Secret验证失败: "+err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "创建Secret失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusCreated, models.ToSecretResponse(createdSecret)) // Return basic info
}

// UpdateSecret godoc
// @Summary Update a Secret
// @Description Update an existing Secret using YAML/JSON definition.
// @Tags Secrets
// @Accept json,yaml
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "Secret Name"
// @Param secret body corev1.Secret true "Updated Secret definition"
// @Success 200 {object} models.SecretResponse "Secret updated (metadata only)"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Input or Name/Namespace Mismatch"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Resource modified"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/secrets/{name} [put]
func (h *SecretHandler) UpdateSecret(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或Secret名称格式")
		return
	}

	var secret corev1.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		respondError(c, http.StatusBadRequest, "无效的请求体格式: "+err.Error())
		return
	}
	if secret.Name != name || secret.Namespace != namespace {
		respondError(c, http.StatusBadRequest, "路径参数与请求体中的名称/命名空间不匹配")
		return
	}
	if secret.Kind != "Secret" || (secret.APIVersion != "v1" && secret.APIVersion != "") {
		respondError(c, http.StatusBadRequest, "无效的Kind或APIVersion，应为 Secret/v1")
		return
	}
	if secret.APIVersion == "" {
		secret.APIVersion = "v1"
	}

	updatedSecret, err := h.service.Update(k8sClient.Clientset, namespace, &secret)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Secret不存在")
			return
		}
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "资源已被修改，请获取最新版本后重试")
			return
		}
		if _, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, "更新Secret验证失败: "+err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "更新Secret失败: "+err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, models.ToSecretResponse(updatedSecret)) // Return basic info
}

// DeleteSecret godoc
// @Summary Delete a Secret
// @Description Delete a specific Secret by namespace and name.
// @Tags Secrets
// @Accept json
// @Produce json
// @Param namespace path string true "Namespace"
// @Param name path string true "Secret Name"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} handlers.ErrorResponse "Bad Request - Invalid Name/Namespace"
// @Failure 404 {object} handlers.ErrorResponse "Not Found"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error"
// @Router /api/v1/namespaces/{namespace}/secrets/{name} [delete]
func (h *SecretHandler) DeleteSecret(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或Secret名称格式")
		return
	}

	err := h.service.Delete(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			c.Status(http.StatusNoContent)
			return
		} // Idempotent
		respondError(c, http.StatusInternalServerError, "删除Secret失败: "+err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Re-use or define respond helpers ---
/*
type ErrorResponse struct { Code int `json:"code"`; Message string `json:"message"`}
func respondSuccess(c *gin.Context, code int, data interface{}) { ... }
func respondError(c *gin.Context, code int, message string) { ... }
*/
