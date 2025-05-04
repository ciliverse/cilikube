package handlers

import (
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type EventsHandler struct {
	service *service.EventsService
}

func NewEventsHandler(svc *service.EventsService) *EventsHandler {
	return &EventsHandler{
		service: svc,
	}
}

func (h *EventsHandler) ListEventsHandler(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间")
		return
	}
	events := h.service.List(namespace)
	respondSuccess(c, http.StatusOK, events)
}

func (h *EventsHandler) GetEventsHandler(c *gin.Context) {
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
	event := h.service.Get(namespace, name)
	respondSuccess(c, http.StatusOK, event)
}
