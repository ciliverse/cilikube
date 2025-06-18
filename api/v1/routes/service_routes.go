package routes

import (
	"github.com/ciliverse/cilikube/api/v1/handlers"
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

// RegisterServiceRoutes 注册 Service 相关路由
func RegisterServiceRoutes(router *gin.RouterGroup, service service.ResourceService[*corev1.Service], clusterManager *k8s.ClusterManager) {
	handler := handlers.NewServiceHandler(service, clusterManager)
	RegisterResourceRoutes(router, handler.ResourceHandler, clusterManager, "services")
}
