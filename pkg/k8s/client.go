package k8s

import (
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ClusterInfo struct {
	ServerVersion  string                    `json:"serverVersion"`
	APIResources   []*metav1.APIResourceList `json:"apiResources,omitempty"`
	NodeCount      int                       `json:"nodeCount"`
	NamespaceCount int                       `json:"namespaceCount"`
	Status         string                    `json:"status"`
	LastUpdated    metav1.Time               `json:"lastUpdated"`
}

type Client struct {
	Clientset kubernetes.Interface

	DynamicClient dynamic.Interface

	DiscoveryClient discovery.DiscoveryInterface

	Config *rest.Config

	clusterInfo *ClusterInfo
}

func NewClient(kubeconfig string) (*Client, error) {
	config, err := buildConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build Kubernetes config: %w", err)
	}

	return newClientFromConfig(config)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {

	if kubeconfig == "in-cluster" {
		return rest.InClusterConfig()
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	kubeconfigPath := resolveKubeconfigPath(kubeconfig)
	if kubeconfigPath != "" {
		loadingRules.ExplicitPath = kubeconfigPath
	}

	configOverrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	return clientConfig.ClientConfig()
}

func resolveKubeconfigPath(kubeconfig string) string {

	if kubeconfig != "" && kubeconfig != "default" {
		return kubeconfig
	}

	home := homedir.HomeDir()
	if home == "" {
		return ""
	}

	defaultKubeconfig := filepath.Join(home, ".kube", "config")
	return defaultKubeconfig
}

func newClientFromConfig(config *rest.Config) (*Client, error) {
	// Create configuration copy to avoid modifying original configuration
	clientConfig := *config

	if clientConfig.QPS == 0 {
		clientConfig.QPS = 50.0
	}
	if clientConfig.Burst == 0 {
		clientConfig.Burst = 100
	}

	// Try to create client using original configuration
	clientset, err := kubernetes.NewForConfig(&clientConfig)
	if err != nil {
		// If failed, try to skip TLS verification
		fmt.Printf("warning: failed to create clientset with original config, trying insecure mode: %v\n", err)

		insecureConfig := &rest.Config{
			Host:    clientConfig.Host,
			APIPath: clientConfig.APIPath,
			QPS:     clientConfig.QPS,
			Burst:   clientConfig.Burst,
			// Skip TLS verification
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
			// Preserve authentication information
			Username:    clientConfig.Username,
			Password:    clientConfig.Password,
			BearerToken: clientConfig.BearerToken,
			Timeout:     clientConfig.Timeout,
		}

		clientset, err = kubernetes.NewForConfig(insecureConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create Kubernetes clientset even with insecure mode: %w", err)
		}
		clientConfig = *insecureConfig
	}

	dynamicClient, err := dynamic.NewForConfig(&clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(&clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %w", err)
	}

	client := &Client{
		Clientset:       clientset,
		DynamicClient:   dynamicClient,
		DiscoveryClient: discoveryClient,
		Config:          &clientConfig,
	}

	if err := client.initClusterInfo(); err != nil {
		fmt.Printf("warning: failed to initialize cluster info: %v\n", err)
	}

	return client, nil
}

func NewClientFromContent(kubeconfigData []byte) (*Client, error) {
	if len(kubeconfigData) == 0 {
		return nil, fmt.Errorf("kubeconfig content cannot be empty")
	}

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigData)
	if err != nil {
		return nil, fmt.Errorf("failed to create client config from bytes: %w", err)
	}

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get REST config from client config: %w", err)
	}

	return newClientFromConfig(restConfig)
}

func (c *Client) initClusterInfo() error {

	serverVersion, err := c.DiscoveryClient.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to get server version: %w", err)
	}

	apiResources, err := c.DiscoveryClient.ServerPreferredResources()
	if err != nil {

		fmt.Printf("warning: failed to get API resource list: %v\n", err)
	}

	c.clusterInfo = &ClusterInfo{
		ServerVersion: serverVersion.GitVersion,
		APIResources:  apiResources,
		Status:        "connected",
		LastUpdated:   metav1.Now(),
	}

	return nil
}

func (c *Client) CheckConnection() error {
	if c == nil || c.Clientset == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	_, err := c.DiscoveryClient.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to check Kubernetes API Server connection: %w", err)
	}
	return nil
}

func (c *Client) GetServerVersion() (string, error) {
	if c.clusterInfo != nil && c.clusterInfo.ServerVersion != "" {
		return c.clusterInfo.ServerVersion, nil
	}

	version, err := c.DiscoveryClient.ServerVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get server version: %w", err)
	}

	return version.GitVersion, nil
}

func (c *Client) GetAPIResources() ([]*metav1.APIResourceList, error) {
	if c.clusterInfo != nil && c.clusterInfo.APIResources != nil {
		return c.clusterInfo.APIResources, nil
	}

	resources, err := c.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, fmt.Errorf("failed to get API resource list: %w", err)
	}

	return resources, nil
}

func (c *Client) RefreshClusterInfo() error {
	return c.initClusterInfo()
}

// GetClusterInfo gets cluster information
func (c *Client) GetClusterInfo() *ClusterInfo {
	return c.clusterInfo
}

func (c *Client) GetCustomResource(gvr schema.GroupVersionResource, namespace, name string) (dynamic.ResourceInterface, error) {
	if c.DynamicClient == nil {
		return nil, fmt.Errorf("dynamic client not initialized")
	}

	if namespace != "" {
		return c.DynamicClient.Resource(gvr).Namespace(namespace), nil
	}
	return c.DynamicClient.Resource(gvr), nil
}

func (c *Client) ListCustomResources(gvr schema.GroupVersionResource, namespace string) (dynamic.ResourceInterface, error) {
	if c.DynamicClient == nil {
		return nil, fmt.Errorf("dynamic client not initialized")
	}

	if namespace != "" {
		return c.DynamicClient.Resource(gvr).Namespace(namespace), nil
	}
	return c.DynamicClient.Resource(gvr), nil
}

func (c *Client) HasCustomResourceDefinition(group, version, kind string) (bool, error) {
	if c.clusterInfo == nil || c.clusterInfo.APIResources == nil {
		if err := c.RefreshClusterInfo(); err != nil {
			return false, fmt.Errorf("failed to refresh cluster info: %w", err)
		}
	}

	for _, resourceList := range c.clusterInfo.APIResources {
		if resourceList.GroupVersion == group+"/"+version {
			for _, resource := range resourceList.APIResources {
				if resource.Kind == kind {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (c *Client) GetSupportedAPIVersions() ([]string, error) {
	groups, err := c.DiscoveryClient.ServerGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get API groups: %w", err)
	}

	var versions []string
	for _, group := range groups.Groups {
		for _, version := range group.Versions {
			versions = append(versions, version.GroupVersion)
		}
	}

	return versions, nil
}

func (c *Client) IsResourceNamespaced(gvr schema.GroupVersionResource) (bool, error) {
	if c.clusterInfo == nil || c.clusterInfo.APIResources == nil {
		if err := c.RefreshClusterInfo(); err != nil {
			return false, fmt.Errorf("failed to refresh cluster info: %w", err)
		}
	}

	groupVersion := gvr.Group + "/" + gvr.Version
	if gvr.Group == "" {
		groupVersion = gvr.Version
	}

	for _, resourceList := range c.clusterInfo.APIResources {
		if resourceList.GroupVersion == groupVersion {
			for _, resource := range resourceList.APIResources {
				if resource.Name == gvr.Resource {
					return resource.Namespaced, nil
				}
			}
		}
	}

	return false, fmt.Errorf("resource %s not found", gvr.String())
}
