package k8s

import (
	"context"
	"net/http"

	"github.com/cilliandevops/kopsadmin/server-go/internal/apis/services"
	"github.com/gin-gonic/gin"
	networkingv1 "k8s.io/api/networking/v1"
)

type IngressHandler struct {
	ingressService *services.IngressService
}

// NewIngressHandler creates a new IngressHandler
func NewIngressHandler(ingressService *services.IngressService) *IngressHandler {
	return &IngressHandler{ingressService: ingressService}
}

// handleError is a helper function to handle error responses
func (h *IngressHandler) handleError(c *gin.Context, err error, message string) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": message,
		})
	}
}

// ListIngresses handles the listing of ingresses
func (h *IngressHandler) ListIngresses(c *gin.Context) {
	namespace := c.Param("namespace")
	ingresses, err := h.ingressService.ListIngresses(context.Background(), namespace) // Pass context.Background() here
	if err != nil {
		h.handleError(c, err, "Failed to list ingresses")
		return
	}
	c.JSON(http.StatusOK, ingresses)
}

// GetIngress handles getting a specific ingress
func (h *IngressHandler) GetIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	ingress, err := h.ingressService.GetIngress(context.Background(), namespace, name) // Pass context.Background() here
	if err != nil {
		h.handleError(c, err, "Failed to get ingress")
		return
	}
	c.JSON(http.StatusOK, ingress)
}

// CreateIngress handles the creation of a new ingress
func (h *IngressHandler) CreateIngress(c *gin.Context) {
	var ingress networkingv1.Ingress
	if err := c.ShouldBindJSON(&ingress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	namespace := c.Param("namespace")
	newIngress, err := h.ingressService.CreateIngress(context.Background(), namespace, &ingress) // Pass context.Background() here
	if err != nil {
		h.handleError(c, err, "Failed to create ingress")
		return
	}
	c.JSON(http.StatusOK, newIngress)
}

// DeleteIngress handles the deletion of an ingress
// DeleteIngress handles the deletion of an ingress
func (h *IngressHandler) DeleteIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if err := h.ingressService.DeleteIngress(context.Background(), namespace, name); err != nil { // Fixed: combined declaration and condition with ';'
		h.handleError(c, err, "Failed to delete ingress")
		return
	}
	c.Status(http.StatusNoContent)
}
