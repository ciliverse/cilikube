package k8s

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetK8sClientFromContext 是连接 API 请求和多集群客户端的“桥梁”。
// 它从 Gin 上下文中提取 ":cluster_name" URL 参数，
// 然后使用 ClusterManager 获取该集群的客户端。
// 如果失败，它会直接向客户端写入错误响应并返回 false。
func GetK8sClientFromContext(c *gin.Context, cm *ClusterManager) (*Client, bool) {
	clusterName := c.Param("cluster_name")
	if clusterName == "" {
		// 在新的路由设计下，所有需要客户端的请求都应包含 :cluster_name。
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "缺少集群名称",
			"message": "请求的URL中必须包含目标集群的名称，例如 /api/v1/clusters/my-cluster/pods",
		})
		return nil, false
	}

	// 从 manager 获取指定名称的客户端
	client, err := cm.GetClient(clusterName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   fmt.Sprintf("集群 '%s' 未找到或不可用", clusterName),
			"message": err.Error(),
		})
		return nil, false
	}

	return client, true
}
