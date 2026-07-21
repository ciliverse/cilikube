package k8s

import (
	"fmt"
	"log"
	"sync"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

const defaultInformerResync = 30 * time.Second

// InformerStatus describes the cache sync state for one cluster.
type InformerStatus struct {
	ClusterID string          `json:"cluster_id"`
	Running   bool            `json:"running"`
	Synced    bool            `json:"synced"`
	Resources map[string]bool `json:"resources"`
	StartedAt *time.Time      `json:"started_at,omitempty"`
	LastError string          `json:"last_error,omitempty"`
}

type clusterInformers struct {
	factory      informers.SharedInformerFactory
	stopCh       chan struct{}
	startedAt    time.Time
	podLister    corelisters.PodLister
	nodeLister   corelisters.NodeLister
	nsLister     corelisters.NamespaceLister
	svcLister    corelisters.ServiceLister
	deployLister appslisters.DeploymentLister
	synced       map[string]cache.InformerSynced
	lastError    string
}

// informerState holds per-cluster informer factories for the manager.
type informerState struct {
	mu        sync.RWMutex
	byCluster map[string]*clusterInformers
}

func newInformerState() *informerState {
	return &informerState{byCluster: make(map[string]*clusterInformers)}
}

func (cm *ClusterManager) startInformersForCluster(clusterID string, client *Client) {
	if client == nil || client.Clientset == nil {
		return
	}

	cm.informers.mu.Lock()
	// Stop previous factory for this cluster if any
	if prev, ok := cm.informers.byCluster[clusterID]; ok {
		close(prev.stopCh)
		delete(cm.informers.byCluster, clusterID)
	}

	factory := informers.NewSharedInformerFactory(client.Clientset, defaultInformerResync)
	podInformer := factory.Core().V1().Pods().Informer()
	nodeInformer := factory.Core().V1().Nodes().Informer()
	nsInformer := factory.Core().V1().Namespaces().Informer()
	svcInformer := factory.Core().V1().Services().Informer()
	deployInformer := factory.Apps().V1().Deployments().Informer()

	entry := &clusterInformers{
		factory:      factory,
		stopCh:       make(chan struct{}),
		startedAt:    time.Now(),
		podLister:    factory.Core().V1().Pods().Lister(),
		nodeLister:   factory.Core().V1().Nodes().Lister(),
		nsLister:     factory.Core().V1().Namespaces().Lister(),
		svcLister:    factory.Core().V1().Services().Lister(),
		deployLister: factory.Apps().V1().Deployments().Lister(),
		synced: map[string]cache.InformerSynced{
			"pods":        podInformer.HasSynced,
			"nodes":       nodeInformer.HasSynced,
			"namespaces":  nsInformer.HasSynced,
			"services":    svcInformer.HasSynced,
			"deployments": deployInformer.HasSynced,
		},
	}
	cm.informers.byCluster[clusterID] = entry
	cm.informers.mu.Unlock()

	factory.Start(entry.stopCh)

	go func() {
		if !cache.WaitForCacheSync(entry.stopCh, podInformer.HasSynced, nodeInformer.HasSynced, nsInformer.HasSynced, svcInformer.HasSynced, deployInformer.HasSynced) {
			cm.informers.mu.Lock()
			if current, ok := cm.informers.byCluster[clusterID]; ok && current == entry {
				current.lastError = "cache sync timed out or stopped"
			}
			cm.informers.mu.Unlock()
			log.Printf("informer cache sync incomplete for cluster %s", clusterID)
			return
		}
		log.Printf("informer caches synced for cluster %s", clusterID)
	}()
}

func (cm *ClusterManager) stopInformersForCluster(clusterID string) {
	cm.informers.mu.Lock()
	defer cm.informers.mu.Unlock()
	if entry, ok := cm.informers.byCluster[clusterID]; ok {
		close(entry.stopCh)
		delete(cm.informers.byCluster, clusterID)
	}
}

// GetInformerStatus returns informer sync status for a cluster (active cluster if empty).
func (cm *ClusterManager) GetInformerStatus(clusterID string) InformerStatus {
	if clusterID == "" {
		clusterID = cm.GetActiveClusterID()
	}

	status := InformerStatus{
		ClusterID: clusterID,
		Resources: map[string]bool{},
	}

	cm.informers.mu.RLock()
	defer cm.informers.mu.RUnlock()
	entry, ok := cm.informers.byCluster[clusterID]
	if !ok {
		return status
	}

	status.Running = true
	started := entry.startedAt
	status.StartedAt = &started
	status.LastError = entry.lastError
	allSynced := true
	for name, syncedFn := range entry.synced {
		synced := syncedFn()
		status.Resources[name] = synced
		if !synced {
			allSynced = false
		}
	}
	status.Synced = allSynced && entry.lastError == ""
	return status
}

func (cm *ClusterManager) getClusterInformers(clusterID string) (*clusterInformers, error) {
	if clusterID == "" {
		clusterID = cm.GetActiveClusterID()
	}
	cm.informers.mu.RLock()
	defer cm.informers.mu.RUnlock()
	entry, ok := cm.informers.byCluster[clusterID]
	if !ok {
		return nil, fmt.Errorf("informers not running for cluster %s", clusterID)
	}
	return entry, nil
}

// ListPodsFromCache lists pods from the informer cache.
func (cm *ClusterManager) ListPodsFromCache(clusterID, namespace string) ([]*corev1.Pod, error) {
	entry, err := cm.getClusterInformers(clusterID)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return entry.podLister.List(labels.Everything())
	}
	return entry.podLister.Pods(namespace).List(labels.Everything())
}

// ListNodesFromCache lists nodes from the informer cache.
func (cm *ClusterManager) ListNodesFromCache(clusterID string) ([]*corev1.Node, error) {
	entry, err := cm.getClusterInformers(clusterID)
	if err != nil {
		return nil, err
	}
	return entry.nodeLister.List(labels.Everything())
}

// ListNamespacesFromCache lists namespaces from the informer cache.
func (cm *ClusterManager) ListNamespacesFromCache(clusterID string) ([]*corev1.Namespace, error) {
	entry, err := cm.getClusterInformers(clusterID)
	if err != nil {
		return nil, err
	}
	return entry.nsLister.List(labels.Everything())
}

// ListServicesFromCache lists services from the informer cache.
func (cm *ClusterManager) ListServicesFromCache(clusterID, namespace string) ([]*corev1.Service, error) {
	entry, err := cm.getClusterInformers(clusterID)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return entry.svcLister.List(labels.Everything())
	}
	return entry.svcLister.Services(namespace).List(labels.Everything())
}

// ListDeploymentsFromCache lists deployments from the informer cache.
func (cm *ClusterManager) ListDeploymentsFromCache(clusterID, namespace string) ([]*appsv1.Deployment, error) {
	entry, err := cm.getClusterInformers(clusterID)
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return entry.deployLister.List(labels.Everything())
	}
	return entry.deployLister.Deployments(namespace).List(labels.Everything())
}
