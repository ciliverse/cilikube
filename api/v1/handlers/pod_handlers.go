package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	// Keep for potential future use (like WebSocket ping)
	"github.com/ciliverse/cilikube/api/v1/models"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/utils" // Assuming utils package exists
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type PodHandler struct {
	service *service.PodService
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for dev, **restrict in production**
		return true
	},
}

func NewPodHandler(svc *service.PodService) *PodHandler {
	return &PodHandler{service: svc}
}

// ListNamespaces ... (保持不变)
func (h *PodHandler) ListNamespaces(c *gin.Context) {
	namespaces, err := h.service.ListNamespaces()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取命名空间失败: "+err.Error())
		return
	}
	// Use respondSuccess for consistency
	respondSuccess(c, http.StatusOK, namespaces)
}

// GetPod ... (保持不变)
func (h *PodHandler) GetPod(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}
	if !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的Pod名称格式")
		return
	}

	pod, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取Pod失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, models.ToPodResponse(pod))
}

// CreatePod 创建Pod (支持 JSON 或 YAML)
func (h *PodHandler) CreatePod(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	contentType := c.ContentType()
	var createdPod *corev1.Pod
	var err error

	if strings.Contains(contentType, "yaml") || strings.Contains(contentType, "x-yaml") {
		yamlBody, errRead := io.ReadAll(c.Request.Body)
		if errRead != nil {
			respondError(c, http.StatusBadRequest, "读取 YAML 请求体失败: "+errRead.Error())
			return
		}
		if len(yamlBody) == 0 {
			respondError(c, http.StatusBadRequest, "请求体不能为空 (YAML)")
			return
		}
		createdPod, err = h.service.CreateFromYAML(namespace, yamlBody)

	} else if strings.Contains(contentType, "json") { // Explicitly check for JSON
		var req models.CreatePodRequest
		if errBind := c.ShouldBindJSON(&req); errBind != nil {
			respondError(c, http.StatusBadRequest, "无效的 JSON 请求格式: "+errBind.Error())
			return
		}
		// Validate name from JSON body if present
		if !utils.ValidateResourceName(req.Name) {
			respondError(c, http.StatusBadRequest, "无效的 Pod 名称格式 (来自 JSON body)")
			return
		}

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        req.Name,
				Namespace:   namespace, // Use namespace from path, ignore req.Namespace if any
				Labels:      req.Labels,
				Annotations: req.Annotations,
			},
			Spec: req.Spec,
		}
		// Use the original service.Create method for JSON objects
		createdPod, err = h.service.Create(namespace, pod)
	} else {
		respondError(c, http.StatusUnsupportedMediaType, "不支持的 Content-Type，请使用 application/json 或 application/yaml")
		return
	}

	// --- Handle Response ---
	if err != nil {
		if e, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, e.Error())
			return
		}
		if errors.IsAlreadyExists(err) {
			respondError(c, http.StatusConflict, "Pod已存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "创建Pod失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, models.ToPodResponse(createdPod))
}

// UpdatePod 更新Pod (支持 JSON 或 YAML)
func (h *PodHandler) UpdatePod(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}

	contentType := c.ContentType()
	var result *corev1.Pod
	var err error

	if strings.Contains(contentType, "yaml") || strings.Contains(contentType, "x-yaml") {
		yamlBody, errRead := io.ReadAll(c.Request.Body)
		if errRead != nil {
			respondError(c, http.StatusBadRequest, "读取 YAML 请求体失败: "+errRead.Error())
			return
		}
		if len(yamlBody) == 0 {
			respondError(c, http.StatusBadRequest, "请求体不能为空 (YAML)")
			return
		}
		result, err = h.service.UpdateFromYAML(namespace, name, yamlBody)

	} else if strings.Contains(contentType, "json") { // Explicitly check for JSON
		// --- Handle JSON Input ---
		// Get the existing Pod first to apply changes correctly
		existingPod, errGet := h.service.Get(namespace, name)
		if errGet != nil {
			if errors.IsNotFound(errGet) {
				respondError(c, http.StatusNotFound, "Pod不存在，无法更新")
				return
			}
			respondError(c, http.StatusInternalServerError, "获取Pod失败: "+errGet.Error())
			return
		}

		// Bind the JSON request which contains only the fields to update
		var req models.UpdatePodRequest // Assumes this model only contains fields allowed to change
		if errBind := c.ShouldBindJSON(&req); errBind != nil {
			respondError(c, http.StatusBadRequest, "无效的 JSON 请求格式: "+errBind.Error())
			return
		}

		// Create a deep copy and apply changes from the request
		updatedPod := existingPod.DeepCopy()
		updatedPod.Labels = req.Labels           // Replace labels
		updatedPod.Annotations = req.Annotations // Replace annotations
		updatedPod.Spec = req.Spec               // Replace the entire spec

		// *** Call the correct Update method in the service ***
		result, err = h.service.Update(namespace, updatedPod) // Use the method taking a Pod object

	} else {
		respondError(c, http.StatusUnsupportedMediaType, "不支持的 Content-Type，请使用 application/json 或 application/yaml")
		return
	}

	// --- Handle Response ---
	if err != nil {
		if e, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, e.Error())
			return
		}
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod不存在 (可能在更新期间被删除)")
			return
		}
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "Pod已被修改，请重试 (ResourceVersion conflict)")
			return
		}
		respondError(c, http.StatusInternalServerError, "更新Pod失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, models.ToPodResponse(result))
}

// DeletePod ... (保持不变, 使用 204)
func (h *PodHandler) DeletePod(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}

	err := h.service.Delete(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			// Idempotent: Return success even if not found
			c.Status(http.StatusNoContent)
			return
		}
		respondError(c, http.StatusInternalServerError, "删除Pod失败: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent) // Success
}

// ListPods ... (保持不变)
func (h *PodHandler) ListPods(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}

	labelSelector := c.Query("labelSelector")
	limitStr := c.DefaultQuery("limit", "500") // Sensible default limit
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 500 // Fallback
	}

	pods, err := h.service.List(namespace, labelSelector, limit)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取Pod列表失败: "+err.Error())
		return
	}

	response := models.PodListResponse{
		Items: make([]models.PodResponse, 0, len(pods.Items)),
		// Total reflects items *in this batch*. K8s list doesn't give total count easily.
		Total: len(pods.Items),
	}
	for _, pod := range pods.Items {
		response.Items = append(response.Items, models.ToPodResponse(&pod))
	}

	respondSuccess(c, http.StatusOK, response)
}

// WatchPods 监听Pod变化 (**修正**)
// @Summary 实时监听Pod变化 (SSE)
// @Tags Pods
// @Param namespace path string true "Namespace"
// @Param labelSelector query string false "标签选择器"
// @Success 200 {object} models.WatchEvent "Server-Sent Event stream"
// @Router /api/v1/namespaces/{namespace}/watch/pods [get] // Adjusted route based on routes.go
func (h *PodHandler) WatchPods(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	if !utils.ValidateNamespace(namespace) {
		respondError(c, http.StatusBadRequest, "无效的命名空间格式")
		return
	}
	labelSelector := c.Query("labelSelector")

	watcher, err := h.service.Watch(namespace, labelSelector)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "开始监听Pod失败: "+err.Error())
		return
	}
	defer watcher.Stop()

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Adjust in production

	// Use c.Stream to handle the streaming goroutine
	// ** REMOVED chanStream := c.Stream(...) and <-chanStream **
	c.Stream(func(w io.Writer) bool {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				// Channel closed by Watcher (e.g., timeout, error)
				fmt.Println("Watcher channel closed")
				c.SSEvent("close", gin.H{"message": "Watcher channel closed"}) // Send a close event
				return false                                                   // Stop streaming
			}
			// Send event data
			c.SSEvent("message", toWatchEvent(event))
			// c.Writer.Flush() // Gin's SSEvent might handle flushing
			return true // Keep connection open and continue streaming

		case <-c.Request.Context().Done():
			// Client closed the connection
			fmt.Println("Client disconnected from watch stream")
			return false // Stop streaming
		}
	})

	// The handlers function returns here, but the goroutine inside c.Stream continues
	fmt.Println("WatchPods handlers finished setup, streaming started.")
}

// GetPodLogs ... (保持不变)
func (h *PodHandler) GetPodLogs(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	container := c.Query("container")
	follow := c.Query("follow") == "true"
	timestamps := c.Query("timestamps") == "true"
	tailLinesStr := c.Query("tailLines")

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}
	if container == "" {
		respondError(c, http.StatusBadRequest, "必须提供 'container' 查询参数")
		return
	}

	// Optional: Check container exists
	pod, err := h.service.Get(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod 不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取 Pod 信息失败: "+err.Error())
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
		respondError(c, http.StatusNotFound, fmt.Sprintf("容器 '%s' 在 Pod '%s' 中未找到", container, name))
		return
	}

	logOptions := &corev1.PodLogOptions{
		Container:  container,
		Follow:     follow,
		Timestamps: timestamps,
	}

	if tailLinesStr != "" {
		tailLines, err := strconv.ParseInt(tailLinesStr, 10, 64)
		if err != nil || tailLines <= 0 {
			respondError(c, http.StatusBadRequest, "无效的 'tailLines' 参数")
			return
		}
		logOptions.TailLines = &tailLines
	}

	logStream, err := h.service.GetPodLogs(namespace, name, logOptions)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取日志失败: "+err.Error())
		return
	}
	defer logStream.Close()

	c.Header("Content-Type", "text/plain; charset=utf-8")

	if follow {
		// Use c.Stream for proper handling of follow=true
		c.Stream(func(w io.Writer) bool {
			// Copy in a loop or use a buffer? io.Copy might block.
			// Let's try a simple copy first.
			_, errCopy := io.Copy(w, logStream)
			if errCopy != nil {
				fmt.Printf("Log stream copy error (follow=true): %v\n", errCopy)
				return false // Stop streaming on error
			}
			// If io.Copy returns without error (e.g. stream closed by server), stop.
			// Need to also check for client disconnect.
			select {
			case <-c.Request.Context().Done():
				fmt.Println("Client disconnected during log follow")
				return false
			default:
				// If copy returned EOF, we should stop.
				// If copy returned no error, but context isn't done,
				// it implies the stream finished? This path needs care.
				// For simplicity, let's assume io.Copy blocks until error/EOF/cancel.
				// If it returned without error, assume EOF.
				fmt.Println("Log stream finished (follow=true, EOF?)")
				return false
			}
		})
	} else {
		// For non-follow, just copy everything once.
		_, errCopy := io.Copy(c.Writer, logStream)
		if errCopy != nil {
			// Cannot send error response here as headers/body might be partially sent
			fmt.Printf("Log stream copy error (follow=false): %v\n", errCopy)
			// c.AbortWithError(http.StatusInternalServerError, errCopy) // May not work
		}
	}
}

// ExecIntoPod ... (保持不变)
func (h *PodHandler) ExecIntoPod(c *gin.Context) {
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

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		// Cannot use respondError here reliably
		return
	}
	defer ws.Close()

	command := []string{commandStr}
	if argsStr != "" {
		// Basic splitting, consider more robust parsing if needed
		args := strings.Split(argsStr, ",")
		command = append(command, args...)
	}

	wsStreamHandler := NewWebSocketStreamHandler(ws, enableStdin, enableStdout, enableStderr)
	defer wsStreamHandler.Close() // Ensure pipes are closed

	execOptions := service.ExecOptions{
		Namespace:     namespace,
		PodName:       name,
		ContainerName: container,
		Command:       command,
		Stdin:         wsStreamHandler,
		Stdout:        wsStreamHandler,
		Stderr:        wsStreamHandler, // Combine stdout/stderr for simplicity here
		Tty:           enableTty,
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	var execErr error
	execDone := make(chan struct{})
	go func() {
		defer close(execDone)
		fmt.Printf("Executing command: %v in %s/%s/%s\n", command, namespace, name, container)
		execErr = h.service.ExecIntoPod(ctx, execOptions)
		if execErr != nil {
			// Attempt to send error back via WebSocket
			errMsg := fmt.Sprintf("\r\n--- Command Execution Failed ---\r\nError: %v\r\n", execErr)
			wsStreamHandler.WriteMessage(websocket.TextMessage, []byte(errMsg))
			fmt.Printf("ExecIntoPod error: %v\n", execErr)
			// Close the WebSocket connection from server-side on error?
			// ws.Close() // Defer already handles closing
			cancel() // Cancel context to potentially stop read/write loops
		} else {
			fmt.Println("ExecIntoPod finished without error.")
			// Optionally send a success message or just rely on stream closing
			// wsStreamHandler.WriteMessage(websocket.TextMessage, []byte("\r\n--- Command Finished ---\r\n"))
		}
		// Close the pipes from the exec side to signal EOF to ws loops
		wsStreamHandler.ClosePipes()

	}()

	// Wait for execution goroutine to finish OR client to disconnect
	select {
	case <-execDone:
		fmt.Println("Exec goroutine completed.")
		// Maybe wait a tiny bit for final output to be written to WS?
		// time.Sleep(100 * time.Millisecond)
	case <-ctx.Done(): // Triggered by client disconnect or cancel() call
		fmt.Println("Exec context done (client disconnected or error).")
	}

	fmt.Println("Exec handlers exiting.")
	// ws.Close() is handled by defer
}

// GetPodYAML ... (保持不变)
func (h *PodHandler) GetPodYAML(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}

	yamlBytes, err := h.service.GetPodYAML(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod 不存在")
			return
		}
		respondError(c, http.StatusInternalServerError, "获取 Pod YAML 失败: "+err.Error())
		return
	}

	c.Header("Content-Type", "application/yaml")
	c.Data(http.StatusOK, "application/yaml", yamlBytes)
}

// UpdatePodYAML ... (保持不变)
func (h *PodHandler) UpdatePodYAML(c *gin.Context) {
	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))

	if !utils.ValidateNamespace(namespace) || !utils.ValidateResourceName(name) {
		respondError(c, http.StatusBadRequest, "无效的命名空间或 Pod 名称格式")
		return
	}

	contentType := c.ContentType()
	if !strings.Contains(contentType, "yaml") && !strings.Contains(contentType, "x-yaml") {
		respondError(c, http.StatusUnsupportedMediaType, "请求体必须是 application/yaml 类型")
		return
	}

	yamlBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		respondError(c, http.StatusBadRequest, "读取 YAML 请求体失败: "+err.Error())
		return
	}
	if len(yamlBody) == 0 {
		respondError(c, http.StatusBadRequest, "请求体不能为空")
		return
	}

	updatedPod, err := h.service.UpdateFromYAML(namespace, name, yamlBody)
	if err != nil {
		if e, ok := err.(*service.ValidationError); ok {
			respondError(c, http.StatusBadRequest, e.Error())
			return
		}
		if errors.IsNotFound(err) {
			respondError(c, http.StatusNotFound, "Pod 不存在")
			return
		}
		if errors.IsConflict(err) {
			respondError(c, http.StatusConflict, "Pod已被修改，请重试 (ResourceVersion conflict)")
			return
		}
		respondError(c, http.StatusInternalServerError, "更新 Pod YAML 失败: "+err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, models.ToPodResponse(updatedPod))
}

// --- Helper Functions ---

// toWatchEvent ... (保持不变)
func toWatchEvent(event watch.Event) interface{} {
	pod, ok := event.Object.(*corev1.Pod)
	resp := gin.H{
		"type": string(event.Type),
	}
	if ok {
		resp["object"] = models.ToPodResponse(pod)
	} else {
		if status, okStatus := event.Object.(*metav1.Status); okStatus {
			resp["error"] = fmt.Sprintf("K8s API Error: %s (Code: %d)", status.Message, status.Code)
			resp["status"] = status
		} else {
			resp["error"] = "事件对象类型不是 Pod 或 Status"
			resp["rawObject"] = fmt.Sprintf("%T", event.Object) // Show type
		}
	}
	return resp
}

// respondSuccess ... (保持不变)
func respondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"code":    code,
		"data":    data,
		"message": "success",
	})
}

// respondError ... (保持不变)
func respondError(c *gin.Context, code int, message string) {
	fmt.Printf("API Error: Status=%d, Message=%s\n", code, message)
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

// --- WebSocket Stream Handler Helper (Improved Closing) ---
type WebSocketStreamHandler struct {
	ws          *websocket.Conn
	stdinPipeR  *io.PipeReader
	stdinPipeW  *io.PipeWriter
	stdoutPipeR *io.PipeReader // Renamed for clarity
	stdoutPipeW *io.PipeWriter
	readDone    chan struct{}
	writeDone   chan struct{}
}

func NewWebSocketStreamHandler(ws *websocket.Conn, stdin, stdout, stderr bool) *WebSocketStreamHandler {
	// Stdin pipe: ws -> pipe -> remotecommand
	prStdin, pwStdin := io.Pipe()
	// Stdout pipe: remotecommand -> pipe -> ws
	prStdout, pwStdout := io.Pipe()

	h := &WebSocketStreamHandler{
		ws:          ws,
		stdinPipeR:  prStdin,
		stdinPipeW:  pwStdin,
		stdoutPipeR: prStdout, // Read from this pipe in writeLoop
		stdoutPipeW: pwStdout, // Write to this pipe from remotecommand
		readDone:    make(chan struct{}),
		writeDone:   make(chan struct{}),
	}

	go h.readLoop(stdin)
	go h.writeLoop() // writeLoop now reads from h.stdoutPipeR

	return h
}

func (h *WebSocketStreamHandler) Read(p []byte) (n int, err error) {
	n, err = h.stdinPipeR.Read(p)
	// fmt.Printf("stdinPipeR Read: n=%d, err=%v\n", n, err)
	return
}

func (h *WebSocketStreamHandler) Write(p []byte) (n int, err error) {
	// This is called by remotecommand (stdout/stderr)
	n, err = h.stdoutPipeW.Write(p)
	// fmt.Printf("stdoutPipeW Write: n=%d, err=%v\n", n, err)
	return
}

func (h *WebSocketStreamHandler) readLoop(stdinEnabled bool) {
	defer close(h.readDone)
	defer h.stdinPipeW.CloseWithError(io.EOF) // Signal EOF to reader when done
	defer fmt.Println("WS Read loop exited")

	if !stdinEnabled {
		fmt.Println("WS Read loop: stdin disabled.")
		return
	}

	fmt.Println("WS Read loop started")
	for {
		msgType, payload, err := h.ws.ReadMessage()
		if err != nil {
			// Propagate error to the pipe reader
			h.stdinPipeW.CloseWithError(fmt.Errorf("WebSocket read error: %w", err))
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("WebSocket closed normally by client.")
			} else {
				fmt.Printf("WebSocket read error: %v\n", err)
			}
			return
		}

		if msgType == websocket.TextMessage || msgType == websocket.BinaryMessage {
			_, err = h.stdinPipeW.Write(payload)
			if err != nil {
				fmt.Printf("Error writing to stdin pipe: %v\n", err)
				// Close WebSocket? Or just return? Let's return.
				return
			}
		} else if msgType == websocket.CloseMessage {
			fmt.Println("WS Read loop: Received close message")
			h.stdinPipeW.CloseWithError(io.EOF) // Normal close
			return
		}
	}
}

func (h *WebSocketStreamHandler) writeLoop() {
	defer close(h.writeDone)
	// No need to close h.stdoutPipeR, reader closes itself on EOF/error
	defer fmt.Println("WS Write loop exited")
	fmt.Println("WS Write loop started")

	buf := make([]byte, 4096) // Slightly larger buffer
	for {
		// Read from the pipe where remotecommand writes stdout/stderr
		n, err := h.stdoutPipeR.Read(buf)
		if n > 0 {
			// Determine message type (Binary for TTY usually, Text otherwise)
			// Let's default to Binary for now, maybe adjust based on TTY later.
			errWrite := h.ws.WriteMessage(websocket.BinaryMessage, buf[:n])
			if errWrite != nil {
				fmt.Printf("WebSocket write error: %v\n", errWrite)
				// If we can't write to WS, we should stop reading the pipe.
				// Closing the pipe reader will signal the writer (remotecommand)
				h.stdoutPipeR.Close() // Close the reader side
				return
			}
		}
		if err != nil {
			if err == io.EOF {
				fmt.Println("Stdout/Stderr pipe closed (EOF). Signaling WS close.")
				// Send a WebSocket close message before exiting
				h.ws.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Stream finished"))
			} else {
				fmt.Printf("Error reading from stdout/stderr pipe: %v\n", err)
				h.ws.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Pipe read error"))
			}
			return // Exit loop on EOF or error
		}
	}
}

// WriteMessage sends an explicit message (e.g., error text) to the WebSocket client
func (h *WebSocketStreamHandler) WriteMessage(messageType int, data []byte) error {
	// Consider adding locking if multiple goroutines could call this,
	// but in this setup, only the main handlers calls it for errors.
	return h.ws.WriteMessage(messageType, data)
}

// ClosePipes closes the writing ends of the pipes, typically called by the exec goroutine when done.
func (h *WebSocketStreamHandler) ClosePipes() {
	fmt.Println("Closing WS stream handlers pipes")
	// Closing the writer signals EOF to the reader in the other goroutine
	h.stdinPipeW.CloseWithError(io.EOF)
	h.stdoutPipeW.CloseWithError(io.EOF)
}

// Close closes the WebSocket connection and waits for loops to finish.
// This should likely be called by the main handlers function (e.g., in defer).
func (h *WebSocketStreamHandler) Close() {
	fmt.Println("Closing WebSocket stream handlers")
	// Closing the WebSocket connection should cause ReadMessage/WriteMessage
	// in the loops to return errors, thus stopping the loops.
	h.ws.Close()
	// Wait for loops to finish cleaning up (optional, but good practice)
	<-h.readDone
	<-h.writeDone
	fmt.Println("WebSocket stream handlers fully closed")
}
