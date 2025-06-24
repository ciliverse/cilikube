package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

// PodLogsHandler 处理 Pod 日志相关请求
type PodLogsHandler struct {
	service        *service.PodLogsService
	clusterManager *k8s.ClusterManager
}

// NewPodLogsHandler 创建 Pod 日志处理器
func NewPodLogsHandler(service *service.PodLogsService, clusterManager *k8s.ClusterManager) *PodLogsHandler {
	return &PodLogsHandler{
		service:        service,
		clusterManager: clusterManager,
	}
}

// GetPodLogs 获取 Pod 日志
func (h *PodLogsHandler) GetPodLogs(c *gin.Context) {
	k8sClient, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Param("namespace"))
	name := strings.TrimSpace(c.Param("name"))
	container := c.Query("container")
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

	// 检查 Pod 和容器是否存在
	pod, err := h.service.Get(k8sClient.Clientset, namespace, name)
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

	// 配置日志选项
	logOptions := buildLogOptions(container, timestamps, tailLinesStr)

	// 获取日志流
	logStream, err := h.service.GetPodLogs(k8sClient.Clientset, namespace, name, logOptions)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "获取日志失败: "+err.Error())
		return
	}
	defer func() {
		if closeErr := logStream.Close(); closeErr != nil {
			fmt.Printf("关闭日志流出错: %v\n", closeErr)
		}
	}()

	// 设置 SSE 响应头
	setSSEHeaders(c)

	// 检查是否支持 Flush
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		respondError(c, http.StatusInternalServerError, "当前响应不支持 SSE")
		return
	}

	// 初始化 Scanner 并设置缓冲区
	scanner := initScanner(logStream)

	// 异步处理日志流
	logChan := make(chan string)
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	go func() {
		defer close(logChan)
		defer close(errChan)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case logChan <- scanner.Text():
			}
		}
		if err := scanner.Err(); err != nil {
			errChan <- err
		}
	}()

	// 处理日志输出
	for {
		select {
		case line, ok := <-logChan:
			if !ok {
				// 日志流结束，主动推送 event: end
				fmt.Fprintf(c.Writer, "event: end\ndata: [END]\n\n")
				flusher.Flush()
				return
			}
			if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", line); err != nil {
				fmt.Printf("写入 SSE 数据出错: %v\n", err)
				return
			}
			flusher.Flush()
		case err := <-errChan:
			fmt.Printf("读取日志出错: %v\n", err)
			return
		case <-ctx.Done():
			fmt.Println("客户端断开连接")
			return
		}
	}
}

// buildLogOptions 构建日志选项
func buildLogOptions(container string, timestamps bool, tailLinesStr string) *corev1.PodLogOptions {
	defaultTailLines := int64(100)
	var tailLines *int64

	if tailLinesStr != "" {
		val, err := strconv.ParseInt(tailLinesStr, 10, 64)
		if err == nil && val > 0 {
			tailLines = &val
		}
	} else {
		tailLines = &defaultTailLines
	}

	return &corev1.PodLogOptions{
		Container:  container,
		Timestamps: timestamps,
		TailLines:  tailLines,
	}
}

// setSSEHeaders 设置 SSE 响应头
func setSSEHeaders(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Access-Control-Allow-Origin", "*") // 支持跨域
	c.Header("Access-Control-Allow-Credentials", "true")
}

// initScanner 初始化日志扫描器
func initScanner(logStream io.ReadCloser) *bufio.Scanner {
	scanner := bufio.NewScanner(logStream)
	scanner.Buffer(make([]byte, 4096), 1024*1024) // 设置更大的缓冲区
	return scanner
}
