package handlers

import (
	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
)

// ServiceHandler 处理 Service 相关的请求
type ServiceHandler struct {
	*ResourceHandler[*corev1.Service]
}

// NewServiceHandler 创建 Service 处理器
func NewServiceHandler(service service.ResourceService[*corev1.Service], clusterManager *k8s.ClusterManager) *ServiceHandler {
	return &ServiceHandler{
		ResourceHandler: NewResourceHandler(service, clusterManager, "services"),
	}
}
