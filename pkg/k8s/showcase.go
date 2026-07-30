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

	ShowcaseProdClusterID   = "showcase-prod"
	ShowcaseProdClusterName = "prod-east"
	ShowcaseProdAPIHost     = "https://prod-east.showcase.cilikube.local"
	ShowcaseProdVersion     = "v1.33.4-showcase"

	ShowcaseStagingClusterID   = "showcase-staging"
	ShowcaseStagingClusterName = "staging-lab"
	ShowcaseStagingAPIHost     = "https://staging-lab.showcase.cilikube.local"
	ShowcaseStagingVersion     = "v1.35.0-showcase"
)

// IsShowcase reports whether the process is running the public read-only exhibit.
// Local/dev defaults to false — set CILIKUBE_SHOWCASE=1 only on the public host.
func IsShowcase() bool {
	v := strings.TrimSpace(os.Getenv("CILIKUBE_SHOWCASE"))
	return v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(os.Getenv("CILIKUBE_MODE"), "showcase")
}

// IsShowcaseConfig detects an in-memory demo rest.Config (any fleet showcase cluster).
func IsShowcaseConfig(cfg *rest.Config) bool {
	if cfg == nil || cfg.Host == "" {
		return false
	}
	return strings.Contains(cfg.Host, "showcase.cilikube.local")
}

// ShowcaseProfile is one simulated cluster in the public multi-cluster fleet.
type ShowcaseProfile struct {
	ID          string
	Name        string
	APIHost     string
	Version     string
	Environment string
	Description string
	Region      string
	Seed        func() []runtime.Object
}

// ShowcaseProfiles returns the public demo fleet (primary demo + prod + staging).
func ShowcaseProfiles() []ShowcaseProfile {
	return []ShowcaseProfile{
		{
			ID:          ShowcaseClusterID,
			Name:        ShowcaseClusterName,
			APIHost:     ShowcaseAPIHost,
			Version:     ShowcaseVersion,
			Environment: "demo",
			Description: "Public showcase cluster (simulated — full resource inventory)",
			Region:      "exhibit",
			Seed:        ShowcaseSeedObjects,
		},
		{
			ID:          ShowcaseProdClusterID,
			Name:        ShowcaseProdClusterName,
			APIHost:     ShowcaseProdAPIHost,
			Version:     ShowcaseProdVersion,
			Environment: "production",
			Description: "Simulated production-east fleet member (includes unhealthy pods / warnings)",
			Region:      "demo-east",
			Seed:        ShowcaseProdSeedObjects,
		},
		{
			ID:          ShowcaseStagingClusterID,
			Name:        ShowcaseStagingClusterName,
			APIHost:     ShowcaseStagingAPIHost,
			Version:     ShowcaseStagingVersion,
			Environment: "staging",
			Description: "Simulated staging lab (smaller, mostly healthy)",
			Region:      "demo-west",
			Seed:        ShowcaseStagingSeedObjects,
		},
	}
}

// NewShowcaseClient builds a client-go fake stack with seeded objects (primary demo host).
func NewShowcaseClient(objs ...runtime.Object) *Client {
	return NewShowcaseClientAt(ShowcaseAPIHost, ShowcaseVersion, objs...)
}

// NewShowcaseClientAt builds a fake client bound to a specific showcase API host/version.
func NewShowcaseClientAt(apiHost, version string, objs ...runtime.Object) *Client {
	if apiHost == "" {
		apiHost = ShowcaseAPIHost
	}
	if version == "" {
		version = ShowcaseVersion
	}

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
		Host:  apiHost,
		QPS:   50,
		Burst: 100,
	}

	return &Client{
		Clientset:       cs,
		DynamicClient:   dyn,
		DiscoveryClient: cs.Discovery(),
		Config:          cfg,
		clusterInfo: &ClusterInfo{
			ServerVersion: version,
			Status:        "connected",
		},
	}
}

// registerShowcaseCluster installs the multi-cluster simulated fleet.
func (cm *ClusterManager) registerShowcaseCluster() {
	now := time.Now()
	profiles := ShowcaseProfiles()
	cm.lock.Lock()
	for _, p := range profiles {
		client := NewShowcaseClientAt(p.APIHost, p.Version, p.Seed()...)
		info := store.Cluster{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Provider:    "showcase",
			Environment: p.Environment,
			Region:      p.Region,
			Version:     p.Version,
			Status:      "Active",
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		cm.clients[p.ID] = client
		cm.clientInfo[p.ID] = info
		cm.nameToID[p.Name] = p.ID
		cm.statusCache[p.ID] = ClusterInfoResponse{
			ID:          p.ID,
			Name:        p.Name,
			Server:      p.APIHost,
			Version:     p.Version,
			Status:      "Available",
			Source:      "showcase",
			Environment: p.Environment,
		}
	}
	cm.lock.Unlock()

	_ = cm.SetActiveClusterByID(ShowcaseClusterID)
}
