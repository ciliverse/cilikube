package service

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// --- NodeClient (Cluster-scoped) ---
type NodeClient struct{}

func (c *NodeClient) Get(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.GetOptions) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Get(ctx, name, opts)
}
func (c *NodeClient) List(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().Nodes().List(ctx, opts)
}
func (c *NodeClient) Create(ctx context.Context, clientset kubernetes.Interface, _ string, obj *corev1.Node, opts metav1.CreateOptions) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Create(ctx, obj, opts)
}
func (c *NodeClient) Update(ctx context.Context, clientset kubernetes.Interface, _ string, obj *corev1.Node, opts metav1.UpdateOptions) (*corev1.Node, error) {
	return clientset.CoreV1().Nodes().Update(ctx, obj, opts)
}
func (c *NodeClient) Delete(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().Nodes().Delete(ctx, name, opts)
}
func (c *NodeClient) Watch(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().Nodes().Watch(ctx, opts)
}

// --- PodClient (Namespaced) ---
type PodClient struct{}

func (c *PodClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*corev1.Pod, error) {
	return clientset.CoreV1().Pods(namespace).Get(ctx, name, opts)
}
func (c *PodClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().Pods(namespace).List(ctx, opts)
}
func (c *PodClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	return clientset.CoreV1().Pods(namespace).Create(ctx, obj, opts)
}
func (c *PodClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	return clientset.CoreV1().Pods(namespace).Update(ctx, obj, opts)
}
func (c *PodClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}
func (c *PodClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().Pods(namespace).Watch(ctx, opts)
}

// --- DeploymentClient (Namespaced) ---
type DeploymentClient struct{}

func (c *DeploymentClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*appsv1.Deployment, error) {
	return clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}
func (c *DeploymentClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.AppsV1().Deployments(namespace).List(ctx, opts)
}
func (c *DeploymentClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *appsv1.Deployment, opts metav1.CreateOptions) (*appsv1.Deployment, error) {
	return clientset.AppsV1().Deployments(namespace).Create(ctx, obj, opts)
}
func (c *DeploymentClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *appsv1.Deployment, opts metav1.UpdateOptions) (*appsv1.Deployment, error) {
	return clientset.AppsV1().Deployments(namespace).Update(ctx, obj, opts)
}
func (c *DeploymentClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}
func (c *DeploymentClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.AppsV1().Deployments(namespace).Watch(ctx, opts)
}

// --- ServiceClient (Namespaced) ---
type ServiceClient struct{}

func (c *ServiceClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*corev1.Service, error) {
	return clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
}
func (c *ServiceClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().Services(namespace).List(ctx, opts)
}
func (c *ServiceClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.Service, opts metav1.CreateOptions) (*corev1.Service, error) {
	return clientset.CoreV1().Services(namespace).Create(ctx, obj, opts)
}
func (c *ServiceClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.Service, opts metav1.UpdateOptions) (*corev1.Service, error) {
	return clientset.CoreV1().Services(namespace).Update(ctx, obj, opts)
}
func (c *ServiceClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().Services(namespace).Delete(ctx, name, opts)
}
func (c *ServiceClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().Services(namespace).Watch(ctx, opts)
}

// --- DaemonSetClient (Namespaced) ---
type DaemonSetClient struct{}

func (c *DaemonSetClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*appsv1.DaemonSet, error) {
	return clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, opts)
}
func (c *DaemonSetClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.AppsV1().DaemonSets(namespace).List(ctx, opts)
}
func (c *DaemonSetClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *appsv1.DaemonSet, opts metav1.CreateOptions) (*appsv1.DaemonSet, error) {
	return clientset.AppsV1().DaemonSets(namespace).Create(ctx, obj, opts)
}
func (c *DaemonSetClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *appsv1.DaemonSet, opts metav1.UpdateOptions) (*appsv1.DaemonSet, error) {
	return clientset.AppsV1().DaemonSets(namespace).Update(ctx, obj, opts)
}
func (c *DaemonSetClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.AppsV1().DaemonSets(namespace).Delete(ctx, name, opts)
}
func (c *DaemonSetClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.AppsV1().DaemonSets(namespace).Watch(ctx, opts)
}

// --- IngressClient (Namespaced) ---
type IngressClient struct{}

func (c *IngressClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*networkingv1.Ingress, error) {
	return clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, opts)
}
func (c *IngressClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.NetworkingV1().Ingresses(namespace).List(ctx, opts)
}
func (c *IngressClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *networkingv1.Ingress, opts metav1.CreateOptions) (*networkingv1.Ingress, error) {
	return clientset.NetworkingV1().Ingresses(namespace).Create(ctx, obj, opts)
}
func (c *IngressClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *networkingv1.Ingress, opts metav1.UpdateOptions) (*networkingv1.Ingress, error) {
	return clientset.NetworkingV1().Ingresses(namespace).Update(ctx, obj, opts)
}
func (c *IngressClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.NetworkingV1().Ingresses(namespace).Delete(ctx, name, opts)
}
func (c *IngressClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.NetworkingV1().Ingresses(namespace).Watch(ctx, opts)
}

// --- ConfigMapClient (Namespaced) ---
type ConfigMapClient struct{}

func (c *ConfigMapClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*corev1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, opts)
}
func (c *ConfigMapClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().ConfigMaps(namespace).List(ctx, opts)
}
func (c *ConfigMapClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.ConfigMap, opts metav1.CreateOptions) (*corev1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps(namespace).Create(ctx, obj, opts)
}
func (c *ConfigMapClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.ConfigMap, opts metav1.UpdateOptions) (*corev1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps(namespace).Update(ctx, obj, opts)
}
func (c *ConfigMapClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
}
func (c *ConfigMapClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().ConfigMaps(namespace).Watch(ctx, opts)
}

// --- SecretClient (Namespaced) ---
type SecretClient struct{}

func (c *SecretClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*corev1.Secret, error) {
	return clientset.CoreV1().Secrets(namespace).Get(ctx, name, opts)
}
func (c *SecretClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().Secrets(namespace).List(ctx, opts)
}
func (c *SecretClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.Secret, opts metav1.CreateOptions) (*corev1.Secret, error) {
	return clientset.CoreV1().Secrets(namespace).Create(ctx, obj, opts)
}
func (c *SecretClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.Secret, opts metav1.UpdateOptions) (*corev1.Secret, error) {
	return clientset.CoreV1().Secrets(namespace).Update(ctx, obj, opts)
}
func (c *SecretClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().Secrets(namespace).Delete(ctx, name, opts)
}
func (c *SecretClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().Secrets(namespace).Watch(ctx, opts)
}

// --- PVCClient (Namespaced) ---
type PVCClient struct{}

func (c *PVCClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*corev1.PersistentVolumeClaim, error) {
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, opts)
}
func (c *PVCClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, opts)
}
func (c *PVCClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.PersistentVolumeClaim, opts metav1.CreateOptions) (*corev1.PersistentVolumeClaim, error) {
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, obj, opts)
}
func (c *PVCClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.PersistentVolumeClaim, opts metav1.UpdateOptions) (*corev1.PersistentVolumeClaim, error) {
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, obj, opts)
}
func (c *PVCClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, opts)
}
func (c *PVCClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Watch(ctx, opts)
}

// --- PVClient (Cluster-scoped) ---
type PVClient struct{}

func (c *PVClient) Get(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.GetOptions) (*corev1.PersistentVolume, error) {
	return clientset.CoreV1().PersistentVolumes().Get(ctx, name, opts)
}
func (c *PVClient) List(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().PersistentVolumes().List(ctx, opts)
}
func (c *PVClient) Create(ctx context.Context, clientset kubernetes.Interface, _ string, obj *corev1.PersistentVolume, opts metav1.CreateOptions) (*corev1.PersistentVolume, error) {
	return clientset.CoreV1().PersistentVolumes().Create(ctx, obj, opts)
}
func (c *PVClient) Update(ctx context.Context, clientset kubernetes.Interface, _ string, obj *corev1.PersistentVolume, opts metav1.UpdateOptions) (*corev1.PersistentVolume, error) {
	return clientset.CoreV1().PersistentVolumes().Update(ctx, obj, opts)
}
func (c *PVClient) Delete(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().PersistentVolumes().Delete(ctx, name, opts)
}
func (c *PVClient) Watch(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().PersistentVolumes().Watch(ctx, opts)
}

// --- StatefulSetClient (Namespaced) ---
type StatefulSetClient struct{}

func (c *StatefulSetClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*appsv1.StatefulSet, error) {
	return clientset.AppsV1().StatefulSets(namespace).Get(ctx, name, opts)
}
func (c *StatefulSetClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.AppsV1().StatefulSets(namespace).List(ctx, opts)
}
func (c *StatefulSetClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *appsv1.StatefulSet, opts metav1.CreateOptions) (*appsv1.StatefulSet, error) {
	return clientset.AppsV1().StatefulSets(namespace).Create(ctx, obj, opts)
}
func (c *StatefulSetClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *appsv1.StatefulSet, opts metav1.UpdateOptions) (*appsv1.StatefulSet, error) {
	return clientset.AppsV1().StatefulSets(namespace).Update(ctx, obj, opts)
}
func (c *StatefulSetClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.AppsV1().StatefulSets(namespace).Delete(ctx, name, opts)
}
func (c *StatefulSetClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.AppsV1().StatefulSets(namespace).Watch(ctx, opts)
}

// --- NamespaceClient (Cluster-scoped) ---
type NamespaceClient struct{}

func (c *NamespaceClient) Get(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.GetOptions) (*corev1.Namespace, error) {
	return clientset.CoreV1().Namespaces().Get(ctx, name, opts)
}
func (c *NamespaceClient) List(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().Namespaces().List(ctx, opts)
}
func (c *NamespaceClient) Create(ctx context.Context, clientset kubernetes.Interface, _ string, obj *corev1.Namespace, opts metav1.CreateOptions) (*corev1.Namespace, error) {
	return clientset.CoreV1().Namespaces().Create(ctx, obj, opts)
}
func (c *NamespaceClient) Update(ctx context.Context, clientset kubernetes.Interface, _ string, obj *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error) {
	return clientset.CoreV1().Namespaces().Update(ctx, obj, opts)
}
func (c *NamespaceClient) Delete(ctx context.Context, clientset kubernetes.Interface, _ string, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().Namespaces().Delete(ctx, name, opts)
}
func (c *NamespaceClient) Watch(ctx context.Context, clientset kubernetes.Interface, _ string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().Namespaces().Watch(ctx, opts)
}

// --- JobClient (Namespaced) ---
type JobClient struct{}

func (c *JobClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*batchv1.Job, error) {
	return clientset.BatchV1().Jobs(namespace).Get(ctx, name, opts)
}
func (c *JobClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.BatchV1().Jobs(namespace).List(ctx, opts)
}
func (c *JobClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *batchv1.Job, opts metav1.CreateOptions) (*batchv1.Job, error) {
	return clientset.BatchV1().Jobs(namespace).Create(ctx, obj, opts)
}
func (c *JobClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *batchv1.Job, opts metav1.UpdateOptions) (*batchv1.Job, error) {
	return clientset.BatchV1().Jobs(namespace).Update(ctx, obj, opts)
}
func (c *JobClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.BatchV1().Jobs(namespace).Delete(ctx, name, opts)
}
func (c *JobClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.BatchV1().Jobs(namespace).Watch(ctx, opts)
}

// --- CronJobClient (Namespaced) ---
type CronJobClient struct{}

func (c *CronJobClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*batchv1.CronJob, error) {
	return clientset.BatchV1().CronJobs(namespace).Get(ctx, name, opts)
}
func (c *CronJobClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.BatchV1().CronJobs(namespace).List(ctx, opts)
}
func (c *CronJobClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *batchv1.CronJob, opts metav1.CreateOptions) (*batchv1.CronJob, error) {
	return clientset.BatchV1().CronJobs(namespace).Create(ctx, obj, opts)
}
func (c *CronJobClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *batchv1.CronJob, opts metav1.UpdateOptions) (*batchv1.CronJob, error) {
	return clientset.BatchV1().CronJobs(namespace).Update(ctx, obj, opts)
}
func (c *CronJobClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.BatchV1().CronJobs(namespace).Delete(ctx, name, opts)
}
func (c *CronJobClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.BatchV1().CronJobs(namespace).Watch(ctx, opts)
}

// --- NetworkPolicyClient (Namespaced) ---
type NetworkPolicyClient struct{}

func (c *NetworkPolicyClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*networkingv1.NetworkPolicy, error) {
	return clientset.NetworkingV1().NetworkPolicies(namespace).Get(ctx, name, opts)
}
func (c *NetworkPolicyClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.NetworkingV1().NetworkPolicies(namespace).List(ctx, opts)
}
func (c *NetworkPolicyClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *networkingv1.NetworkPolicy, opts metav1.CreateOptions) (*networkingv1.NetworkPolicy, error) {
	return clientset.NetworkingV1().NetworkPolicies(namespace).Create(ctx, obj, opts)
}
func (c *NetworkPolicyClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *networkingv1.NetworkPolicy, opts metav1.UpdateOptions) (*networkingv1.NetworkPolicy, error) {
	return clientset.NetworkingV1().NetworkPolicies(namespace).Update(ctx, obj, opts)
}
func (c *NetworkPolicyClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.NetworkingV1().NetworkPolicies(namespace).Delete(ctx, name, opts)
}
func (c *NetworkPolicyClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.NetworkingV1().NetworkPolicies(namespace).Watch(ctx, opts)
}

// --- HPAClient (Namespaced) ---
type HPAClient struct{}

func (c *HPAClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	return clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(ctx, name, opts)
}
func (c *HPAClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).List(ctx, opts)
}
func (c *HPAClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *autoscalingv2.HorizontalPodAutoscaler, opts metav1.CreateOptions) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	return clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Create(ctx, obj, opts)
}
func (c *HPAClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *autoscalingv2.HorizontalPodAutoscaler, opts metav1.UpdateOptions) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	return clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Update(ctx, obj, opts)
}
func (c *HPAClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Delete(ctx, name, opts)
}
func (c *HPAClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Watch(ctx, opts)
}

// --- PDBClient (Namespaced) ---
type PDBClient struct{}

func (c *PDBClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*policyv1.PodDisruptionBudget, error) {
	return clientset.PolicyV1().PodDisruptionBudgets(namespace).Get(ctx, name, opts)
}
func (c *PDBClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.PolicyV1().PodDisruptionBudgets(namespace).List(ctx, opts)
}
func (c *PDBClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *policyv1.PodDisruptionBudget, opts metav1.CreateOptions) (*policyv1.PodDisruptionBudget, error) {
	return clientset.PolicyV1().PodDisruptionBudgets(namespace).Create(ctx, obj, opts)
}
func (c *PDBClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *policyv1.PodDisruptionBudget, opts metav1.UpdateOptions) (*policyv1.PodDisruptionBudget, error) {
	return clientset.PolicyV1().PodDisruptionBudgets(namespace).Update(ctx, obj, opts)
}
func (c *PDBClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.PolicyV1().PodDisruptionBudgets(namespace).Delete(ctx, name, opts)
}
func (c *PDBClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.PolicyV1().PodDisruptionBudgets(namespace).Watch(ctx, opts)
}

// --- ResourceQuotaClient (Namespaced) ---
type ResourceQuotaClient struct{}

func (c *ResourceQuotaClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*corev1.ResourceQuota, error) {
	return clientset.CoreV1().ResourceQuotas(namespace).Get(ctx, name, opts)
}
func (c *ResourceQuotaClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().ResourceQuotas(namespace).List(ctx, opts)
}
func (c *ResourceQuotaClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.ResourceQuota, opts metav1.CreateOptions) (*corev1.ResourceQuota, error) {
	return clientset.CoreV1().ResourceQuotas(namespace).Create(ctx, obj, opts)
}
func (c *ResourceQuotaClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.ResourceQuota, opts metav1.UpdateOptions) (*corev1.ResourceQuota, error) {
	return clientset.CoreV1().ResourceQuotas(namespace).Update(ctx, obj, opts)
}
func (c *ResourceQuotaClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().ResourceQuotas(namespace).Delete(ctx, name, opts)
}
func (c *ResourceQuotaClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().ResourceQuotas(namespace).Watch(ctx, opts)
}

// --- LimitRangeClient (Namespaced) ---
type LimitRangeClient struct{}

func (c *LimitRangeClient) Get(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.GetOptions) (*corev1.LimitRange, error) {
	return clientset.CoreV1().LimitRanges(namespace).Get(ctx, name, opts)
}
func (c *LimitRangeClient) List(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	return clientset.CoreV1().LimitRanges(namespace).List(ctx, opts)
}
func (c *LimitRangeClient) Create(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.LimitRange, opts metav1.CreateOptions) (*corev1.LimitRange, error) {
	return clientset.CoreV1().LimitRanges(namespace).Create(ctx, obj, opts)
}
func (c *LimitRangeClient) Update(ctx context.Context, clientset kubernetes.Interface, namespace string, obj *corev1.LimitRange, opts metav1.UpdateOptions) (*corev1.LimitRange, error) {
	return clientset.CoreV1().LimitRanges(namespace).Update(ctx, obj, opts)
}
func (c *LimitRangeClient) Delete(ctx context.Context, clientset kubernetes.Interface, namespace, name string, opts metav1.DeleteOptions) error {
	return clientset.CoreV1().LimitRanges(namespace).Delete(ctx, name, opts)
}
func (c *LimitRangeClient) Watch(ctx context.Context, clientset kubernetes.Interface, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return clientset.CoreV1().LimitRanges(namespace).Watch(ctx, opts)
}
