package k8s

import (
	"fmt"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ClientManager manages multiple Kubernetes client instances and the active client.
type ClientManager struct {
	mu           sync.RWMutex
	clients      map[string]*Client // Map of cluster name to Client
	activeClient *Client
	activeName   string
}

// NewClientManager creates a new ClientManager.
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*Client),
	}
}

// AddOrReplaceClient adds a new Kubernetes client or replaces an existing one for the given cluster name and kubeconfig path.
// If this is the first client being added, it becomes the active client.
func (cm *ClientManager) AddOrReplaceClient(clusterName, kubeconfigPath string) (*Client, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	k8sClient, err := NewClient(kubeconfigPath) // NewClient is your existing func
	if err != nil {
		return nil, fmt.Errorf("failed to create client for cluster '%s' with path '%s': %w", clusterName, kubeconfigPath, err)
	}

	cm.clients[clusterName] = k8sClient
	fmt.Printf("Client for cluster '%s' added/updated.\n", clusterName)

	// If no active client or replacing the active one, set this as active
	if cm.activeClient == nil || cm.activeName == clusterName {
		cm.activeClient = k8sClient
		cm.activeName = clusterName
		fmt.Printf("Cluster '%s' is now the active cluster.\n", clusterName)
	}
	return k8sClient, nil
}

// SetActiveClient sets the active Kubernetes client.
func (cm *ClientManager) SetActiveClient(clusterName string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client, exists := cm.clients[clusterName]
	if !exists {
		// Attempt to load it if it wasn't pre-loaded but exists in config
		// This part depends on how you manage loading vs. on-demand creation
		// For now, we assume it must be added first via AddOrReplaceClient
		return fmt.Errorf("client for cluster '%s' not found. Add it first", clusterName)
	}

	cm.activeClient = client
	cm.activeName = clusterName
	fmt.Printf("Switched active cluster to: %s\n", clusterName)
	return nil
}

// GetActiveClient returns the active Kubernetes client (Clientset and Config).
func (cm *ClientManager) GetActiveClient() (*Client, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if cm.activeClient == nil {
		return nil, fmt.Errorf("no active Kubernetes client configured")
	}
	return cm.activeClient, nil
}

// GetActiveClientset returns the Kubernetes.Interface for the active cluster.
func (cm *ClientManager) GetActiveClientset() (kubernetes.Interface, error) {
	client, err := cm.GetActiveClient()
	if err != nil {
		return nil, err
	}
	return client.Clientset, nil
}

// GetActiveConfig returns the *rest.Config for the active cluster.
func (cm *ClientManager) GetActiveConfig() (*rest.Config, error) {
	client, err := cm.GetActiveClient()
	if err != nil {
		return nil, err
	}
	return client.Config, nil
}

// GetClientByName returns a specific client by name, does not change active client
func (cm *ClientManager) GetClientByName(clusterName string) (*Client, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	client, exists := cm.clients[clusterName]
	if !exists {
		return nil, fmt.Errorf("client for cluster '%s' not found", clusterName)
	}
	return client, nil
}

// RemoveClient removes a client from the manager.
// If the removed client was active, the active client becomes nil.
// Consider logic to switch to a default or another client if needed.
func (cm *ClientManager) RemoveClient(clusterName string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.activeName == clusterName {
		cm.activeClient = nil
		cm.activeName = ""
		fmt.Printf("Active cluster '%s' removed. No active cluster set.\n", clusterName)
	}
	delete(cm.clients, clusterName)
	fmt.Printf("Client for cluster '%s' removed.\n", clusterName)
}

func (cm *ClientManager) GetActiveClusterName() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.activeName
}
