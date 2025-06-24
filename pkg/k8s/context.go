// package k8s

// GetK8sClientFromContext 是连接 API 请求和多集群客户端的“桥梁”。
// 它从 Gin 上下文中提取 ":cluster_name" URL 参数，
// 然后使用 ClusterManager 获取该集群的客户端。
// 如果失败，它会直接向客户端写入错误响应并返回 false。
// func GetK8sClientFromContext(c *gin.Context, cm *ClusterManager) (*Client, bool) {
// 	clusterName := c.Param("cluster_name")
// 	if clusterName == "" {
// 		// 在新的路由设计下，所有需要客户端的请求都应包含 :cluster_name。
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "缺少集群名称",
// 			"message": "请求的URL中必须包含目标集群的名称，例如 /api/v1/clusters/my-cluster/pods",
// 		})
// 		return nil, false
// 	}

// 	// 从 manager 获取指定名称的客户端
// 	client, err := cm.GetClientByName(clusterName)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error":   fmt.Sprintf("集群 '%s' 未找到或不可用", clusterName),
// 			"message": err.Error(),
// 		})
// 		return nil, false
// 	}

//		return client, true
//	}
package k8s

import (
	"fmt"
	"net/http"

	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

// GetClientFromQuery 从 URL 查询参数中获取 clusterId 并返回对应的 k8s 客户端。
// 这是所有资源操作处理函数的“守门员”。
func GetClientFromQuery(c *gin.Context, cm *ClusterManager) (*Client, bool) {
	clusterID := c.Query("clusterId")
	if clusterID == "" {
		// 如果没有提供 clusterId，可以尝试使用当前激活的集群作为后备
		activeID := cm.GetActiveClusterID()
		if activeID == "" {
			utils.ApiError(c, http.StatusBadRequest, "缺少 'clusterId' 查询参数, 且没有活动的默认集群", "e.g., /api/v1/nodes?clusterId=cls-xxxxx")
			return nil, false
		}
		clusterID = activeID
	}

	client, err := cm.GetClientByID(clusterID)
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, fmt.Sprintf("集群ID '%s' 未找到或不可用", clusterID), err.Error())
		return nil, false
	}

	return client, true
}
