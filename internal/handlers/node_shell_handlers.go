package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// NodeShellHandler starts a short-lived privileged pod on a node and attaches a shell.
type NodeShellHandler struct {
	exec           *service.PodExecService
	clusterManager *k8s.ClusterManager
	audit          *service.AuditService
	upgrader       websocket.Upgrader
}

func NewNodeShellHandler(exec *service.PodExecService, cm *k8s.ClusterManager, audit *service.AuditService) *NodeShellHandler {
	return &NodeShellHandler{
		exec:           exec,
		clusterManager: cm,
		audit:          audit,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (h *NodeShellHandler) Shell(c *gin.Context) {
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("node shell upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	if !adminRole(c) {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("admin role required\r\n"))
		return
	}
	client, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("failed to get Kubernetes client\r\n"))
		return
	}
	nodeName := c.Param("name")
	pod, err := service.StartNodeShell(c.Request.Context(), client.Clientset, nodeName, c.Query("image"))
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("node shell: %v\r\n", err)))
		return
	}
	logShell(h.audit, c, "node", nodeName)
	defer func() {
		if delErr := service.DeleteNodeShell(context.Background(), client.Clientset, pod.Name); delErr != nil {
			log.Printf("delete node shell pod %s: %v", pod.Name, delErr)
		}
	}()

	stream := &WebSocketStreamHandler{
		conn:        ws,
		stdinChan:   make(chan []byte, 100),
		stdoutChan:  make(chan []byte, 100),
		closeChan:   make(chan struct{}),
		stdinClosed: false,
	}
	go stream.readMessages()
	go stream.writeMessages()
	defer stream.Close()

	err = h.exec.Exec(client.Config, client.Clientset, pod.Namespace, pod.Name, &service.ExecOptions{
		Command:   []string{"chroot", "/host"},
		Container: "shell",
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, stream, stream)
	if err != nil {
		_, _ = stream.Write([]byte(fmt.Sprintf("\r\n--- node shell failed ---\r\n%v\r\n", err)))
	}
}
