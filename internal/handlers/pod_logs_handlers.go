package handlers

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

// PodLogsHandler struct
type PodLogsHandler struct {
	service        *service.PodLogsService
	clusterManager *k8s.ClusterManager
	upgrader       websocket.Upgrader
}

// NewPodLogsHandler creates a new PodLogsHandler
func NewPodLogsHandler(service *service.PodLogsService, clusterManager *k8s.ClusterManager) *PodLogsHandler {
	return &PodLogsHandler{
		service:        service,
		clusterManager: clusterManager,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// GetPodLogs handles WebSocket requests for pod logs
func (h *PodLogsHandler) GetPodLogs(c *gin.Context) {
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	defer ws.Close()

	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		ws.WriteMessage(websocket.TextMessage, []byte("Failed to get Kubernetes client"))
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	container := c.Query("container")
	timestamps := c.Query("timestamps") == "true"
	tailLinesStr := c.Query("tailLines")

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		ws.WriteMessage(websocket.TextMessage, []byte("Invalid namespace or pod name"))
		return
	}
	if container == "" {
		ws.WriteMessage(websocket.TextMessage, []byte("Container name is required"))
		return
	}

	pod, err := h.service.Get(k8sClient.Clientset, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			ws.WriteMessage(websocket.TextMessage, []byte("Pod not found"))
			return
		}
		ws.WriteMessage(websocket.TextMessage, []byte("Failed to get pod info: "+err.Error()))
		return
	}

	containerFound := false
	for _, cont := range append(pod.Spec.Containers, pod.Spec.InitContainers...) {
		if cont.Name == container {
			containerFound = true
			break
		}
	}
	if !containerFound {
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Container '%s' not found in pod '%s'", container, name)))
		return
	}

	follow := c.Query("follow") == "true"

	logOptions := buildLogOptions(container, timestamps, tailLinesStr, follow)
	logStream, err := h.service.GetPodLogs(k8sClient.Clientset, namespace, name, logOptions)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Failed to get log stream: "+err.Error()))
		return
	}
	defer logStream.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				cancel()
				break
			}
		}
	}()

	scanner := bufio.NewScanner(logStream)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			err := ws.WriteMessage(websocket.TextMessage, scanner.Bytes())
			if err != nil {
				return
			}
		}
	}
}

func buildLogOptions(container string, timestamps bool, tailLinesStr string, follow bool) *corev1.PodLogOptions {
	var tailLines int64 = 1000
	if val, err := strconv.ParseInt(tailLinesStr, 10, 64); err == nil {
		tailLines = val
	}

	return &corev1.PodLogOptions{
		Container:  container,
		Follow:     follow,
		Timestamps: timestamps,
		TailLines:  &tailLines,
	}
}
