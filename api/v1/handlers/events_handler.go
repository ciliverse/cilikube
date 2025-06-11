package handlers

import (
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type EventsHandler struct {
	service        *service.EventsService
	clusterManager *k8s.ClusterManager
}

func NewEventsHandler(svc *service.EventsService, cm *k8s.ClusterManager) *EventsHandler {
	return &EventsHandler{
		service:        svc,
		clusterManager: cm,
	}
}

func (h *EventsHandler) ListEventsHandler(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间")
		return
	}
	events := h.service.List(k8sClient.Clientset, namespace)
	respondSuccess(c, http.StatusOK, events)
}

func (h *EventsHandler) GetEventsHandler(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或事件名称格式")
		return
	}
	if name == "" {
		respondError(c, http.StatusBadRequest, "事件名称不能为空")
		return
	}
	event := h.service.Get(k8sClient.Clientset, namespace, name)
	respondSuccess(c, http.StatusOK, event)
}
