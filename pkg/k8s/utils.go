package k8s

import (
	"net/http"

	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GetK8sClientFromContext(c *gin.Context, clusterManager *ClusterManager) (*Client, bool) {
	clusterName := c.Param("cluster_name")
	if clusterName == "" {
		utils.ApiError(c, http.StatusBadRequest, "路径中缺少 'cluster_name' 参数")
		return nil, false
	}
	k8sClient, err := clusterManager.GetClient(clusterName)
	if err != nil {
		utils.ApiError(c, http.StatusNotFound, "集群未找到或不可用: "+clusterName)
		return nil, false
	}
	if k8sClient.Clientset == nil {
		utils.ApiError(c, http.StatusInternalServerError, "集群的客户端内部 Clientset 为空: "+clusterName)
		return nil, false
	}
	return k8sClient, true
}
