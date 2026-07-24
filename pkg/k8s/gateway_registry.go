package k8s

import (
	"fmt"
	"sync"

	"k8s.io/client-go/kubernetes"
	gatewayclientset "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

var gatewayClients sync.Map // kubernetes.Interface → gatewayclientset.Interface

// RegisterGatewayClient associates a Gateway API clientset with a core clientset.
func RegisterGatewayClient(core kubernetes.Interface, gw gatewayclientset.Interface) {
	if core == nil || gw == nil {
		return
	}
	gatewayClients.Store(core, gw)
}

// GatewayClientFor returns the Gateway API client registered for the given core clientset.
func GatewayClientFor(core kubernetes.Interface) (gatewayclientset.Interface, error) {
	if core == nil {
		return nil, fmt.Errorf("kubernetes clientset is nil")
	}
	if v, ok := gatewayClients.Load(core); ok {
		if gw, ok := v.(gatewayclientset.Interface); ok && gw != nil {
			return gw, nil
		}
	}
	return nil, fmt.Errorf("Gateway API client not available for this cluster")
}
