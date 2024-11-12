package client

import (
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sClient provides methods for interacting with the Kubernetes API.
type K8sClient struct {
	k8sClient *kubernetes.Clientset
}

// NewK8sClient creates a new K8sClient instance.
func NewK8sClient(kubeconfig string) (*K8sClient, error) {
	// Build the configuration from the provided kubeconfig path
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Error building Kubernetes config: %v", err)
		return nil, fmt.Errorf("failed to build Kubernetes config: %w", err)
	}

	// Create the clientset
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error creating Kubernetes clientset: %v", err)
		return nil, fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}

	log.Println("Successfully created Kubernetes client")
	return &K8sClient{k8sClient: k8sClient}, nil
}

// K8sClient returns the underlying Kubernetes Clientset.
func (c *K8sClient) K8sClient() *kubernetes.Clientset {
	return c.k8sClient
}
