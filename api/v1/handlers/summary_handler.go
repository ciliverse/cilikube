package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/pkg/k8s"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

// Existing SummaryHandler struct...
type SummaryHandler struct {
	service        *service.SummaryService
	clusterManager *k8s.ClusterManager
}

func NewSummaryHandler(svc *service.SummaryService, cm *k8s.ClusterManager) *SummaryHandler {
	return &SummaryHandler{service: svc, clusterManager: cm}
}

// GetResourceSummary 获取资源汇总信息
func (h *SummaryHandler) GetResourceSummary(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	summary, errs := h.service.GetResourceSummary(k8sClient.Clientset)
	if len(errs) > 0 {
		// Log errors but still return the summary with available data
		for resource, err := range errs {
			c.Header("X-Resource-Error-"+resource, err.Error())
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    summary,
		"message": "成功获取资源汇总",
	})
}

// --- New Handler for Backend Dependencies ---

// GetBackendDependencies godoc
// @Summary Get Backend Dependencies
// @Description Retrieves the list of direct Go module dependencies and their versions from go.mod.
// @Tags Summary
// @Accept json
// @Produce json
// @Success 200 {array} service.BackendDependency "List of backend dependencies"
// @Failure 500 {object} handlers.ErrorResponse "Internal Server Error - Failed to read/parse go.mod"
// @Router /api/v1/summary/backend-dependencies [get]
func (h *SummaryHandler) GetBackendDependencies(c *gin.Context) {
	dependencies, err := h.service.GetBackendDependencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    nil,
			"message": "获取后端依赖失败: " + err.Error(),
		})
		return
	}
	// Use a different response structure if needed, but returning the slice directly is fine
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    dependencies, // Return the slice directly
		"message": "success",
	})
}
