package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/ciliverse/cilikube/pkg/utils"
	"github.com/gin-gonic/gin"
)

// InformerHandler exposes informer cache status and cached listings.
type InformerHandler struct {
	clusterManager *k8s.ClusterManager
}

func NewInformerHandler(cm *k8s.ClusterManager) *InformerHandler {
	return &InformerHandler{clusterManager: cm}
}

func (h *InformerHandler) GetStatus(c *gin.Context) {
	clusterID := c.Query("clusterId")
	status := h.clusterManager.GetInformerStatus(clusterID)
	utils.ApiSuccess(c, status, "informer status retrieved successfully")
}

func (h *InformerHandler) ListPods(c *gin.Context) {
	clusterID := c.Query("clusterId")
	namespace := c.Query("namespace")
	pods, err := h.clusterManager.ListPodsFromCache(clusterID, namespace)
	if err != nil {
		utils.ApiError(c, http.StatusServiceUnavailable, "failed to list pods from cache", err.Error())
		return
	}
	utils.ApiSuccess(c, gin.H{"items": pods, "total": len(pods), "source": "informer"}, "pods retrieved from informer cache")
}

func (h *InformerHandler) ListNodes(c *gin.Context) {
	clusterID := c.Query("clusterId")
	nodes, err := h.clusterManager.ListNodesFromCache(clusterID)
	if err != nil {
		utils.ApiError(c, http.StatusServiceUnavailable, "failed to list nodes from cache", err.Error())
		return
	}
	utils.ApiSuccess(c, gin.H{"items": nodes, "total": len(nodes), "source": "informer"}, "nodes retrieved from informer cache")
}

func (h *InformerHandler) ListNamespaces(c *gin.Context) {
	clusterID := c.Query("clusterId")
	namespaces, err := h.clusterManager.ListNamespacesFromCache(clusterID)
	if err != nil {
		utils.ApiError(c, http.StatusServiceUnavailable, "failed to list namespaces from cache", err.Error())
		return
	}
	utils.ApiSuccess(c, gin.H{"items": namespaces, "total": len(namespaces), "source": "informer"}, "namespaces retrieved from informer cache")
}

func (h *InformerHandler) ListServices(c *gin.Context) {
	clusterID := c.Query("clusterId")
	namespace := c.Query("namespace")
	services, err := h.clusterManager.ListServicesFromCache(clusterID, namespace)
	if err != nil {
		utils.ApiError(c, http.StatusServiceUnavailable, "failed to list services from cache", err.Error())
		return
	}
	utils.ApiSuccess(c, gin.H{"items": services, "total": len(services), "source": "informer"}, "services retrieved from informer cache")
}

func (h *InformerHandler) ListDeployments(c *gin.Context) {
	clusterID := c.Query("clusterId")
	namespace := c.Query("namespace")
	deployments, err := h.clusterManager.ListDeploymentsFromCache(clusterID, namespace)
	if err != nil {
		utils.ApiError(c, http.StatusServiceUnavailable, "failed to list deployments from cache", err.Error())
		return
	}
	utils.ApiSuccess(c, gin.H{"items": deployments, "total": len(deployments), "source": "informer"}, "deployments retrieved from informer cache")
}
