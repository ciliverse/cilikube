package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// PodExecHandler handles pod execution requests
type PodExecHandler struct {
	service        *service.PodExecService
	clusterManager *k8s.ClusterManager
	upgrader       websocket.Upgrader
}

// NewPodExecHandler creates a new PodExecHandler
func NewPodExecHandler(svc *service.PodExecService, cm *k8s.ClusterManager) *PodExecHandler {
	return &PodExecHandler{
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

// ExecPod handles WebSocket requests for pod execution
func (h *PodExecHandler) ExecPod(c *gin.Context) {
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

	namespace := c.Param("namespace")
	podName := c.Param("name")
	container := c.Query("container")
	command := c.QueryArray("command")

	shell := c.Query("shell")
	if len(command) == 0 {
		if shell != "" {
			command = []string{shell}
		} else {
			command = []string{"/bin/sh"}
		}
	}

	wsStreamHandler := &WebSocketStreamHandler{
		conn:        ws,
		stdinChan:   make(chan []byte, 100),
		stdoutChan:  make(chan []byte, 100),
		closeChan:   make(chan struct{}),
		stdinClosed: false,
	}
	defer wsStreamHandler.Close()

	go wsStreamHandler.readMessages()
	go wsStreamHandler.writeMessages()

	options := &service.ExecOptions{
		Command:   command,
		Container: container,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}

	err = h.service.Exec(k8sClient.Clientset, namespace, podName, options, wsStreamHandler, wsStreamHandler)
	if err != nil {
		errmsg := []byte(fmt.Sprintf("\r\n--- Command Execution Failed ---\r\nError: %v\r\n", err))
		wsStreamHandler.WriteMessage(websocket.TextMessage, errmsg)
		log.Printf("Exec error: %v", err)
		return
	}

	log.Println("Exec finished without error.")
}

// WebSocketStreamHandler implements io.Reader and io.Writer for WebSocket data
type WebSocketStreamHandler struct {
	conn        *websocket.Conn
	stdinChan   chan []byte
	stdoutChan  chan []byte
	closeChan   chan struct{}
	mu          sync.Mutex
	stdinClosed bool
	buffer      []byte
}

// readMessages reads messages from WebSocket and sends to stdinChan
func (h *WebSocketStreamHandler) readMessages() {
	defer h.closeStdin()
	for {
		_, message, err := h.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}
		if message != nil {
			h.stdinChan <- message
		}
	}
}

// closeStdin closes stdinChan
func (h *WebSocketStreamHandler) closeStdin() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.stdinClosed {
		close(h.stdinChan)
		h.stdinClosed = true
	}
}

// writeMessages reads from stdoutChan and writes to WebSocket
func (h *WebSocketStreamHandler) writeMessages() {
	for {
		select {
		case data, ok := <-h.stdoutChan:
			if !ok {
				return
			}
			if err := h.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		case <-h.closeChan:
			return
		}
	}
}

// Read reads data from stdinChan for container execution
func (h *WebSocketStreamHandler) Read(p []byte) (n int, err error) {
	if len(h.buffer) > 0 {
		n = copy(p, h.buffer)
		h.buffer = h.buffer[n:]
		return n, nil
	}
	data, ok := <-h.stdinChan
	if !ok {
		return 0, io.EOF
	}
	n = copy(p, data)
	h.buffer = append(h.buffer, data[n:]...)
	return n, nil
}

// Write writes container output to stdoutChan
func (h *WebSocketStreamHandler) Write(p []byte) (n int, err error) {
	h.stdoutChan <- p
	return len(p), nil
}

// WriteMessage sends a WebSocket message
func (h *WebSocketStreamHandler) WriteMessage(messageType int, data []byte) error {
	return h.conn.WriteMessage(messageType, data)
}

// Close closes the WebSocket connection
func (h *WebSocketStreamHandler) Close() error {
	close(h.closeChan)
	return h.conn.Close()
}

// buildCommand builds the command array
func buildCommand(commandStr, argsStr string) []string {
	command := []string{commandStr}
	if argsStr != "" {
		args := strings.Split(argsStr, " ")
		command = append(command, args...)
	}
	return command
}