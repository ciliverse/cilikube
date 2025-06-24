package service

import (
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// ResourceServiceFactory 资源服务工厂
type ResourceServiceFactory struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

// NewResourceServiceFactory 创建资源服务工厂
func NewResourceServiceFactory() *ResourceServiceFactory {
	return &ResourceServiceFactory{
		services: make(map[string]interface{}),
	}
}

// RegisterService 注册资源服务
func (f *ResourceServiceFactory) RegisterService(name string, service interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.services[name] = service
}

// GetService 获取资源服务
func (f *ResourceServiceFactory) GetService(name string) interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.services[name]
}

// InitializeDefaultServices 初始化所有默认服务
func (f *ResourceServiceFactory) InitializeDefaultServices() {
	f.RegisterService("nodes", NewBaseResourceService[*corev1.Node](new(NodeClient)))
	f.RegisterService("pods", NewBaseResourceService[*corev1.Pod](new(PodClient)))
	f.RegisterService("deployments", NewBaseResourceService[*appsv1.Deployment](new(DeploymentClient)))
	f.RegisterService("services", NewBaseResourceService[*corev1.Service](new(ServiceClient)))
	f.RegisterService("daemonsets", NewBaseResourceService[*appsv1.DaemonSet](new(DaemonSetClient)))
	f.RegisterService("ingresses", NewBaseResourceService[*networkingv1.Ingress](new(IngressClient)))
	f.RegisterService("configmaps", NewBaseResourceService[*corev1.ConfigMap](new(ConfigMapClient)))
	f.RegisterService("secrets", NewBaseResourceService[*corev1.Secret](new(SecretClient)))
	f.RegisterService("persistentvolumeclaims", NewBaseResourceService[*corev1.PersistentVolumeClaim](new(PVCClient)))
	f.RegisterService("persistentvolumes", NewBaseResourceService[*corev1.PersistentVolume](new(PVClient)))
	f.RegisterService("statefulsets", NewBaseResourceService[*appsv1.StatefulSet](new(StatefulSetClient)))
	f.RegisterService("namespaces", NewBaseResourceService[*corev1.Namespace](new(NamespaceClient)))
}
