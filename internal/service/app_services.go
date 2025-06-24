package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// AppServices 作为所有应用服务的集合，在这里统一定义
type AppServices struct {
	// 集群与安装服务
	ClusterService   *ClusterService
	InstallerService InstallerService

	// [新增] 节点指标服务
	NodeMetricsService *NodeMetricsService

	// Kubernetes 资源服务
	NodeService        ResourceService[*corev1.Node]
	NamespaceService   ResourceService[*corev1.Namespace]
	PVService          ResourceService[*corev1.PersistentVolume]
	PodService         ResourceService[*corev1.Pod]
	DeploymentService  ResourceService[*appsv1.Deployment]
	ServiceService     ResourceService[*corev1.Service]
	DaemonSetService   ResourceService[*appsv1.DaemonSet]
	IngressService     ResourceService[*networkingv1.Ingress]
	ConfigMapService   ResourceService[*corev1.ConfigMap]
	SecretService      ResourceService[*corev1.Secret]
	PVCService         ResourceService[*corev1.PersistentVolumeClaim]
	StatefulSetService ResourceService[*appsv1.StatefulSet]

	// Pod 日志与终端服务
	PodLogsService *PodLogsService
	PodExecService *PodExecService
}
