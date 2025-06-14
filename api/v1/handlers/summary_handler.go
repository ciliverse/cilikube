package handlers

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	"net/http"

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

// Existing GetResourceSummary handlers...
func (h *SummaryHandler) GetResourceSummary(c *gin.Context) { /* ... as before ... */
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	summary, _ := h.service.GetResourceSummary(k8sClient.Clientset)
	respondSuccess(c, http.StatusOK, summary)
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
		respondError(c, http.StatusInternalServerError, "获取后端依赖失败: "+err.Error())
		return
	}
	// Use a different response structure if needed, but returning the slice directly is fine
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    dependencies, // Return the slice directly
		"message": "success",
	})
}

// --- Re-use or define respond helpers ---
/* ... */
