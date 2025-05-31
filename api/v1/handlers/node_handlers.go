package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors" // Import for errors.IsNotFound
)

type NodeHandler struct {
	service        *service.NodeService
	clusterManager *k8s.ClusterManager
}

func NewNodeHandler(nodeService *service.NodeService, cm *k8s.ClusterManager) *NodeHandler {
	return &NodeHandler{
		service:        nodeService,
		clusterManager: cm,
	}
}

// Standard API response helper
func apiSuccess(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = "success"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    data,
		"message": message,
	})
}

func apiError(c *gin.Context, statusCode int, message string, details ...string) {
	detailStr := ""
	if len(details) > 0 {
		detailStr = details[0]
	}
	log.Printf("API Error: Status %d, Message: %s, Details: %s, Path: %s", statusCode, message, detailStr, c.Request.URL.Path)
	c.JSON(statusCode, gin.H{
		"code":    statusCode,
		"data":    nil,
		"message": message,
		"details": detailStr,
	})
}

func (h *NodeHandler) getK8sClient(c *gin.Context) (*k8s.Client, bool) {
	clusterName := c.Param("cluster_name")
	if clusterName == "" {
		apiError(c, http.StatusBadRequest, "路径中缺少 'cluster_name' 参数")
		return nil, false
	}

	k8sClient, err := h.clusterManager.GetClient(clusterName)
	if err != nil {
		apiError(c, http.StatusNotFound, fmt.Sprintf("名为 '%s' 的集群未找到或不可用", clusterName), err.Error())
		return nil, false
	}
	if k8sClient.Clientset == nil {
		apiError(c, http.StatusInternalServerError, fmt.Sprintf("集群 '%s' 的客户端内部 Clientset 为空", clusterName))
		return nil, false
	}
	return k8sClient, true
}

func (h *NodeHandler) ListNodes(c *gin.Context) {
	k8sClient, ok := h.getK8sClient(c)
	if !ok {
		return
	}

	labelSelector := c.Query("labelSelector")
	limitStr := c.DefaultQuery("limit", "0")
	continueToken := c.Query("continue")

	limit, convErr := strconv.ParseInt(limitStr, 10, 64)
	if convErr != nil {
		apiError(c, http.StatusBadRequest, "无效的 'limit' 参数", convErr.Error())
		return
	}

	nodeList, serviceErr := h.service.List(k8sClient.Clientset, labelSelector, limit, continueToken)
	if serviceErr != nil {
		apiError(c, http.StatusInternalServerError, "列出 Node 失败", serviceErr.Error())
		return
	}
	apiSuccess(c, nodeList, "节点列表获取成功")
}

func (h *NodeHandler) GetNode(c *gin.Context) {
	k8sClient, ok := h.getK8sClient(c)
	if !ok {
		return
	}
	nodeName := c.Param("name")
	if nodeName == "" {
		apiError(c, http.StatusBadRequest, "路径中缺少 Node 名称")
		return
	}
	node, serviceErr := h.service.Get(k8sClient.Clientset, nodeName)
	if serviceErr != nil {
		if errors.IsNotFound(serviceErr) {
			apiError(c, http.StatusNotFound, fmt.Sprintf("Node '%s' 未找到", nodeName), serviceErr.Error())
		} else {
			apiError(c, http.StatusInternalServerError, fmt.Sprintf("获取 Node '%s' 失败", nodeName), serviceErr.Error())
		}
		return
	}
	apiSuccess(c, node, "Node 获取成功")
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
	k8sClient, ok := h.getK8sClient(c)
	if !ok {
		return
	}
	var node corev1.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		apiError(c, http.StatusBadRequest, "无效的请求体", err.Error())
		return
	}
	createdNode, serviceErr := h.service.Create(k8sClient.Clientset, &node)
	if serviceErr != nil {
		apiError(c, http.StatusInternalServerError, "创建 Node 失败", serviceErr.Error())
		return
	}
	// For create, use http.StatusCreated
	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"data":    createdNode,
		"message": "Node 创建成功",
	})
}

func (h *NodeHandler) UpdateNode(c *gin.Context) {
	k8sClient, ok := h.getK8sClient(c)
	if !ok {
		return
	}
	nodeName := c.Param("name")
	if nodeName == "" {
		apiError(c, http.StatusBadRequest, "路径中缺少 Node 名称")
		return
	}
	var node corev1.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		apiError(c, http.StatusBadRequest, "无效的请求体", err.Error())
		return
	}
	if node.Name != "" && node.Name != nodeName { // Ensure consistency
		apiError(c, http.StatusBadRequest, "请求体中的 Node 名称与路径参数不匹配")
		return
	}
	node.Name = nodeName // Set name from path param to be sure

	updatedNode, serviceErr := h.service.Update(k8sClient.Clientset, &node)
	if serviceErr != nil {
		apiError(c, http.StatusInternalServerError, fmt.Sprintf("更新 Node '%s' 失败", nodeName), serviceErr.Error())
		return
	}
	apiSuccess(c, updatedNode, fmt.Sprintf("Node '%s' 更新成功", nodeName))
}

func (h *NodeHandler) DeleteNode(c *gin.Context) {
	k8sClient, ok := h.getK8sClient(c)
	if !ok {
		return
	}
	nodeName := c.Param("name")
	if nodeName == "" {
		apiError(c, http.StatusBadRequest, "路径中缺少 Node 名称")
		return
	}
	serviceErr := h.service.Delete(k8sClient.Clientset, nodeName)
	if serviceErr != nil {
		apiError(c, http.StatusInternalServerError, fmt.Sprintf("删除 Node '%s' 失败", nodeName), serviceErr.Error())
		return
	}
	apiSuccess(c, gin.H{"name": nodeName}, fmt.Sprintf("Node '%s' 删除请求已接受", nodeName))
}

func (h *NodeHandler) WatchNodes(c *gin.Context) {
	k8sClient, ok := h.getK8sClient(c)
	if !ok {
		return
	}

	labelSelector := c.Query("labelSelector")
	resourceVersion := c.Query("resourceVersion")
	timeoutSecondsStr := c.DefaultQuery("timeoutSeconds", "0") // 0 for server default/indefinite

	timeoutSeconds, _ := strconv.ParseInt(timeoutSecondsStr, 10, 64)

	watcher, serviceErr := h.service.Watch(k8sClient.Clientset, labelSelector, resourceVersion, timeoutSeconds)
	if serviceErr != nil {
		apiError(c, http.StatusInternalServerError, "启动 Node watch 失败", serviceErr.Error())
		return
	}
	defer watcher.Stop()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, flushOK := c.Writer.(http.Flusher)
	if !flushOK {
		apiError(c, http.StatusInternalServerError, "Streaming 不被支持 (HTTP Flusher 不可用)")
		return
	}

	log.Printf("在集群 '%s' 上开始推送 Node watch 事件...", c.Param("cluster_name"))
	ctx := c.Request.Context()

	for {
		select {
		case <-ctx.Done():
			log.Printf("客户端断开连接，停止推送集群 '%s' 的 Node watch 事件。", c.Param("cluster_name"))
			return
		case event, open := <-watcher.ResultChan():
			if !open {
				log.Printf("集群 '%s' 的 Node Watcher channel 已关闭。", c.Param("cluster_name"))
				fmt.Fprintf(c.Writer, "event: close\ndata: Watcher closed by server\n\n")
				flusher.Flush()
				return
			}

			// 构建符合 { type: "ADDED/MODIFIED/DELETED", object: K8sObject } 的事件结构
			responseEvent := gin.H{
				"type":   event.Type,
				"object": event.Object,
			}
			jsonData, marshalErr := json.Marshal(responseEvent)
			if marshalErr != nil {
				log.Printf("错误: 序列化 watch 事件失败 (集群 %s): %v", c.Param("cluster_name"), marshalErr)
				fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", fmt.Sprintf("Error marshalling event: %v", marshalErr))
				flusher.Flush()
				continue
			}
			fmt.Fprintf(c.Writer, "event: watch_update\ndata: %s\n\n", string(jsonData)) // 使用自定义事件名或 event.Type
			flusher.Flush()
		case <-time.After(30 * time.Second): // Keep-alive for some proxies/browsers
			fmt.Fprintf(c.Writer, ": keep-alive\n\n")
			flusher.Flush()
		}
	}
}
