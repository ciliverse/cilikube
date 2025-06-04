package handlers

import (
	"context"
	"fmt"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketStreamHandler 实现 io.Reader 和 io.Writer 接口，处理 WebSocket 数据
type WebSocketStreamHandler struct {
	conn        *websocket.Conn
	stdinChan   chan []byte
	stdoutChan  chan []byte
	closeChan   chan struct{}
	mu          sync.Mutex
	stdinClosed bool
}

// NewWebSocketStreamHandler 创建 WebSocketStreamHandler 实例
func NewWebSocketStreamHandler(conn *websocket.Conn, enableStdin, enableStdout bool) *WebSocketStreamHandler {
	handler := &WebSocketStreamHandler{
		conn:        conn,
		stdinChan:   make(chan []byte, 100),
		stdoutChan:  make(chan []byte, 100),
		closeChan:   make(chan struct{}),
		stdinClosed: false,
	}

	if enableStdin {
		go handler.readMessages()
	}

	if enableStdout {
		go handler.writeMessages()
	}

	return handler
}

// ExecIntoPod 处理 WebSocket 连接，执行容器命令
func (h *PodHandler) ExecIntoPod(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	container := c.Query("container")
	commandStr := c.Query("command")
	argsStr := c.Query("args")

	enableStdin := c.DefaultQuery("stdin", "true") == "true"
	enableStdout := c.DefaultQuery("stdout", "true") == "true"
	enableStderr := c.DefaultQuery("stderr", "true") == "true"
	enableTty := c.Query("tty") == "true"

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) || container == "" || commandStr == "" {
		respondError(c, http.StatusBadRequest, "无效的命名空间/Pod名称/容器/命令参数")
		return
	}

	// 构建命令
	command := buildCommand(commandStr, argsStr)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {

		}
	}(ws)

	wsStreamHandler := NewWebSocketStreamHandler(ws, enableStdin, enableStdout || enableStderr)
	defer func(wsStreamHandler *WebSocketStreamHandler) {
		err := wsStreamHandler.Close()
		if err != nil {

		}
	}(wsStreamHandler)

	execOptions := service.ExecOptions{
		Namespace:     namespace,
		PodName:       name,
		ContainerName: container,
		Command:       command,
		Stdin:         wsStreamHandler,
		Stdout:        wsStreamHandler,
		Stderr:        wsStreamHandler,
		Tty:           enableTty,
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	var execErr error
	execDone := make(chan struct{})
	go func() {
		defer close(execDone)
		log.Printf("Executing command: %v in %s/%s/%s", command, namespace, name, container)
		execErr = h.service.ExecIntoPod(k8sClient.Clientset, k8sClient.Config, ctx, execOptions)
		if execErr != nil {
			errMsg := []byte(fmt.Sprintf("\r\n--- Command Execution Failed ---\r\nError: %v\r\n", execErr))
			if err := wsStreamHandler.WriteMessage(websocket.TextMessage, errMsg); err != nil {
				log.Printf("Failed to send error message: %v", err)
			}
			log.Printf("ExecIntoPod error: %v", execErr)
			cancel()
		} else {
			log.Println("ExecIntoPod finished without error.")
		}
		wsStreamHandler.ClosePipes()
	}()

	select {
	case <-execDone:
		log.Println("Exec goroutine completed.")
	case <-ctx.Done():
		log.Println("Exec context done (client disconnected or error).")
	}

	log.Println("Exec handlers exiting.")
}

// readMessages 从 WebSocket 读取前端输入，发送到 stdinChan
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

// closeStdin 关闭 stdinChan
func (h *WebSocketStreamHandler) closeStdin() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.stdinClosed {
		close(h.stdinChan)
		h.stdinClosed = true
	}
}

// ClosePipes 关闭 stdinChan
func (h *WebSocketStreamHandler) ClosePipes() {
	h.closeStdin()
}

// writeMessages 从 stdoutChan 读取容器输出，发送到 WebSocket
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

// --- Helper Functions ---
// Read 从 stdinChan 读取数据，供容器执行命令使用
func (h *WebSocketStreamHandler) Read(p []byte) (n int, err error) {
	data, ok := <-h.stdinChan
	if !ok {
		return 0, io.EOF
	}
	if data == nil {
		return 0, nil
	}
	n = copy(p, data)
	return
}

// Write 将容器输出写入 stdoutChan，发送到前端
func (h *WebSocketStreamHandler) Write(p []byte) (n int, err error) {
	h.stdoutChan <- append([]byte(nil), p...)
	return len(p), nil
}

// WriteMessage 向 WebSocket 发送消息
func (h *WebSocketStreamHandler) WriteMessage(messageType int, data []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.conn.WriteMessage(messageType, data)
}

// Close 关闭 WebSocket 连接和相关通道
func (h *WebSocketStreamHandler) Close() error {
	close(h.closeChan)
	close(h.stdoutChan)
	return h.conn.Close()
}

// buildCommand 构建执行命令
func buildCommand(commandStr, argsStr string) []string {
	var command []string
	if commandStr == "sh" {
		command = []string{"sh"}
	} else {
		command = []string{commandStr}
		if argsStr != "" {
			args := strings.Split(argsStr, ",")
			command = append(command, args...)
		}
	}
	return command
}
