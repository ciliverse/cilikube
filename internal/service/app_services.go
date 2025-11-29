package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// AppServices serves as a collection of all application services, defined here uniformly
type AppServices struct {
	// Cluster and installer services
	ClusterService   *ClusterService
	InstallerService InstallerService

	// [Added] Node metrics service
	NodeMetricsService *NodeMetricsService

	// [Added] Summary service
	SummaryService *SummaryService

	// Kubernetes resource services
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

	// Pod logs and terminal services
	PodLogsService *PodLogsService
	PodExecService *PodExecService
}
