package k8s

import (
	"os"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/internal/store"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayfake "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned/fake"
)

const (
	ShowcaseClusterID   = "showcase-demo"
	ShowcaseClusterName = "demo"
	ShowcaseAPIHost     = "https://showcase.cilikube.local"
	ShowcaseVersion     = "v1.36.2-showcase"
)

// IsShowcase reports whether the process is running the public read-only exhibit.
// Local/dev defaults to false — set CILIKUBE_SHOWCASE=1 only on the public host.
func IsShowcase() bool {
	v := strings.TrimSpace(os.Getenv("CILIKUBE_SHOWCASE"))
	return v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(os.Getenv("CILIKUBE_MODE"), "showcase")
}

// IsShowcaseConfig detects the in-memory demo rest.Config.
func IsShowcaseConfig(cfg *rest.Config) bool {
	return cfg != nil && cfg.Host == ShowcaseAPIHost
}

// NewShowcaseClient builds a client-go fake stack with seeded objects.
func NewShowcaseClient(objs ...runtime.Object) *Client {
	var coreObjs, gwObjs []runtime.Object
	for _, o := range objs {
		switch o.(type) {
		case *gatewayv1.GatewayClass, *gatewayv1.Gateway, *gatewayv1.HTTPRoute:
			gwObjs = append(gwObjs, o)
		default:
			coreObjs = append(coreObjs, o)
		}
	}
	gwObjs = append(gwObjs, ShowcaseGatewaySeedObjects()...)

	cs := fake.NewSimpleClientset(coreObjs...)
	gwcs := gatewayfake.NewSimpleClientset(gwObjs...)
	RegisterGatewayClient(cs, gwcs)

	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = batchv1.AddToScheme(scheme)
	_ = networkingv1.AddToScheme(scheme)
	_ = rbacv1.AddToScheme(scheme)
	_ = storagev1.AddToScheme(scheme)
	_ = policyv1.AddToScheme(scheme)
	_ = autoscalingv2.AddToScheme(scheme)
	_ = gatewayv1.Install(scheme)
	allDyn := append(append([]runtime.Object{}, coreObjs...), gwObjs...)
	dyn := dynamicfake.NewSimpleDynamicClient(scheme, allDyn...)

	cfg := &rest.Config{
		Host:  ShowcaseAPIHost,
		QPS:   50,
		Burst: 100,
	}

	return &Client{
		Clientset:       cs,
		DynamicClient:   dyn,
		DiscoveryClient: cs.Discovery(),
		Config:          cfg,
		clusterInfo: &ClusterInfo{
			ServerVersion: ShowcaseVersion,
			Status:        "connected",
		},
	}
}

// registerShowcaseCluster installs the demo cluster as the only active client.
func (cm *ClusterManager) registerShowcaseCluster() {
	client := NewShowcaseClient(ShowcaseSeedObjects()...)
	now := time.Now()
	info := store.Cluster{
		ID:          ShowcaseClusterID,
		Name:        ShowcaseClusterName,
		Description: "Public showcase cluster (simulated — no real kube-apiserver)",
		Provider:    "showcase",
		Environment: "demo",
		Region:      "exhibit",
		Version:     ShowcaseVersion,
		Status:      "Active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	cm.lock.Lock()
	cm.clients[ShowcaseClusterID] = client
	cm.clientInfo[ShowcaseClusterID] = info
	cm.nameToID[ShowcaseClusterName] = ShowcaseClusterID
	cm.statusCache[ShowcaseClusterID] = ClusterInfoResponse{
		ID:          ShowcaseClusterID,
		Name:        ShowcaseClusterName,
		Server:      ShowcaseAPIHost,
		Version:     ShowcaseVersion,
		Status:      "Available",
		Source:      "showcase",
		Environment: "demo",
	}
	cm.lock.Unlock()

	_ = cm.SetActiveClusterByID(ShowcaseClusterID)
}
