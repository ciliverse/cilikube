package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// PodPortForwardHandler handles WebSocket port-forward sessions.
type PodPortForwardHandler struct {
	service        *service.PodPortForwardService
	clusterManager *k8s.ClusterManager
	upgrader       websocket.Upgrader
}

func NewPodPortForwardHandler(svc *service.PodPortForwardService, cm *k8s.ClusterManager) *PodPortForwardHandler {
	return &PodPortForwardHandler{
		service:        svc,
		clusterManager: cm,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type portForwardStatusMsg struct {
	Type    string   `json:"type"`
	Ports   []string `json:"ports,omitempty"`
	Message string   `json:"message,omitempty"`
}

// PortForward upgrades to WebSocket, starts a server-local forward, and keeps the
// session alive until the client disconnects. Status JSON frames are sent on the socket.
func (h *PodPortForwardHandler) PortForward(c *gin.Context) {
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("portforward upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		_ = ws.WriteJSON(portForwardStatusMsg{Type: "error", Message: "Failed to get Kubernetes client"})
		return
	}

	namespace := c.Param("namespace")
	podName := c.Param("name")
	rawPorts := c.QueryArray("ports")
	if len(rawPorts) == 0 {
		if p := c.Query("ports"); p != "" {
			rawPorts = []string{p}
		}
	}
	ports, err := service.ParsePortPairs(rawPorts)
	if err != nil {
		_ = ws.WriteJSON(portForwardStatusMsg{Type: "error", Message: err.Error()})
		return
	}

	stopCh := make(chan struct{})
	readyCh := make(chan struct{})
	var closeOnce sync.Once
	stop := func() { closeOnce.Do(func() { close(stopCh) }) }
	defer stop()

	// Client disconnect / ping reader
	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				stop()
				return
			}
		}
	}()

	statusWriter := &wsStatusWriter{conn: ws}
	errWriter := &wsStatusWriter{conn: ws, asError: true}

	errCh := make(chan error, 1)
	go func() {
		errCh <- h.service.Forward(
			k8sClient.Config,
			k8sClient.Clientset,
			namespace,
			podName,
			ports,
			stopCh,
			readyCh,
			statusWriter,
			errWriter,
		)
	}()

	select {
	case <-readyCh:
		_ = ws.WriteJSON(portForwardStatusMsg{
			Type:    "ready",
			Ports:   ports,
			Message: fmt.Sprintf(
				"Ready on CiliKube API host 127.0.0.1 for %s (not your browser localhost)",
				strings.Join(ports, ", "),
			),
		})
	case err := <-errCh:
		_ = ws.WriteJSON(portForwardStatusMsg{Type: "error", Message: err.Error()})
		return
	case <-time.After(30 * time.Second):
		stop()
		_ = ws.WriteJSON(portForwardStatusMsg{Type: "error", Message: "port-forward ready timeout"})
		return
	}

	select {
	case err := <-errCh:
		if err != nil {
			_ = ws.WriteJSON(portForwardStatusMsg{Type: "error", Message: err.Error()})
		} else {
			_ = ws.WriteJSON(portForwardStatusMsg{Type: "closed", Message: "port-forward stopped"})
		}
	case <-stopCh:
		_ = ws.WriteJSON(portForwardStatusMsg{Type: "closed", Message: "client disconnected"})
	}
}

type wsStatusWriter struct {
	conn    *websocket.Conn
	asError bool
	mu      sync.Mutex
}

func (w *wsStatusWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	msgType := "log"
	if w.asError {
		msgType = "stderr"
	}
	b, _ := json.Marshal(portForwardStatusMsg{Type: msgType, Message: strings.TrimSpace(string(p))})
	if err := w.conn.WriteMessage(websocket.TextMessage, b); err != nil {
		return 0, err
	}
	return len(p), nil
}
