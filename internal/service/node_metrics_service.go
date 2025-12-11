package service

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// NodeMetrics represents real-time resource usage of a single node
type NodeMetrics struct {
	NodeName              string `json:"nodeName"`              // Node name
	CPUCores              string `json:"cpuCores"`              // CPU usage (e.g.: "574m")
	CPUPercent            string `json:"cpuPercent"`            // CPU usage percentage (e.g.: "7%")
	MemoryBytes           string `json:"memoryBytes"`           // Memory usage (e.g.: "8820Mi")
	MemoryPercent         string `json:"memoryPercent"`         // Memory usage percentage (e.g.: "60%")
	CPUCapacity           string `json:"cpuCapacity"`           // Total CPU capacity
	MemoryCapacity        string `json:"memoryCapacity"`        // Total memory capacity
	CPURequests           string `json:"cpuRequests"`           // Total CPU requests
	CPURequestsPercent    string `json:"cpuRequestsPercent"`    // CPU requests percentage
	MemoryRequests        string `json:"memoryRequests"`        // Total memory requests
	MemoryRequestsPercent string `json:"memoryRequestsPercent"` // Memory requests percentage
	CPULimits             string `json:"cpuLimits"`             // Total CPU limits
	CPULimitsPercent      string `json:"cpuLimitsPercent"`      // CPU limits percentage
	MemoryLimits          string `json:"memoryLimits"`          // Total memory limits
	MemoryLimitsPercent   string `json:"memoryLimitsPercent"`   // Memory limits percentage
	Timestamp             string `json:"timestamp"`             // Monitoring data timestamp
}

// NodesMetricsResponse represents the response for all nodes metrics
type NodesMetricsResponse struct {
	Nodes []NodeMetrics `json:"nodes"`
	Total int           `json:"total"`
}

// NodeMetricsService provides service for getting node metrics
type NodeMetricsService struct{}

// NewNodeMetricsService creates a new NodeMetricsService instance
func NewNodeMetricsService() *NodeMetricsService {
	return &NodeMetricsService{}
}

// GetNodeMetrics gets real-time metrics of a single node through metrics-server.
// It requires passing the target cluster's rest.Config to create a dedicated metrics client.
func (s *NodeMetricsService) GetNodeMetrics(config *rest.Config, nodeName string) (*NodeMetrics, error) {
	// Create Metrics API client and Kubernetes client
	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Get node metrics from metrics-server
	nodeMetrics, err := metricsClientset.MetricsV1beta1().NodeMetricses().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics for node '%s' from metrics API: %w", nodeName, err)
	}

	// Get node info to calculate capacity and percentages
	node, err := k8sClient.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node info for '%s': %w", nodeName, err)
	}

	// Get pods on this node to calculate requests and limits
	pods, err := k8sClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods for node '%s': %w", nodeName, err)
	}

	// Calculate resource usage, capacity, requests and limits
	cpuUsage := nodeMetrics.Usage.Cpu()
	memoryUsage := nodeMetrics.Usage.Memory()
	cpuCapacity := node.Status.Capacity.Cpu()
	memoryCapacity := node.Status.Capacity.Memory()

	// Calculate requests and limits from all pods on this node
	var cpuRequests, memoryRequests, cpuLimits, memoryLimits resource.Quantity
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		for _, container := range pod.Spec.Containers {
			if req := container.Resources.Requests.Cpu(); req != nil {
				cpuRequests.Add(*req)
			}
			if req := container.Resources.Requests.Memory(); req != nil {
				memoryRequests.Add(*req)
			}
			if limit := container.Resources.Limits.Cpu(); limit != nil {
				cpuLimits.Add(*limit)
			}
			if limit := container.Resources.Limits.Memory(); limit != nil {
				memoryLimits.Add(*limit)
			}
		}
	}

	// Calculate percentages
	cpuPercent := calculatePercentage(cpuUsage.MilliValue(), cpuCapacity.MilliValue())
	memoryPercent := calculatePercentage(memoryUsage.Value(), memoryCapacity.Value())
	cpuRequestsPercent := calculatePercentage(cpuRequests.MilliValue(), cpuCapacity.MilliValue())
	memoryRequestsPercent := calculatePercentage(memoryRequests.Value(), memoryCapacity.Value())
	cpuLimitsPercent := calculatePercentage(cpuLimits.MilliValue(), cpuCapacity.MilliValue())
	memoryLimitsPercent := calculatePercentage(memoryLimits.Value(), memoryCapacity.Value())

	// Format the returned data
	metrics := &NodeMetrics{
		NodeName:              nodeMetrics.Name,
		CPUCores:              formatCPU(cpuUsage.MilliValue()),
		CPUPercent:            cpuPercent + "%",
		MemoryBytes:           formatMemory(memoryUsage.Value()),
		MemoryPercent:         memoryPercent + "%",
		CPUCapacity:           formatCPU(cpuCapacity.MilliValue()),
		MemoryCapacity:        formatMemory(memoryCapacity.Value()),
		CPURequests:           formatCPU(cpuRequests.MilliValue()),
		CPURequestsPercent:    cpuRequestsPercent + "%",
		MemoryRequests:        formatMemory(memoryRequests.Value()),
		MemoryRequestsPercent: memoryRequestsPercent + "%",
		CPULimits:             formatCPU(cpuLimits.MilliValue()),
		CPULimitsPercent:      cpuLimitsPercent + "%",
		MemoryLimits:          formatMemory(memoryLimits.Value()),
		MemoryLimitsPercent:   memoryLimitsPercent + "%",
		Timestamp:             nodeMetrics.Timestamp.Time.Format("2006-01-02T15:04:05Z"),
	}

	return metrics, nil
}

// GetAllNodesMetrics gets real-time metrics of all nodes through metrics-server.
// It requires passing the target cluster's rest.Config to create a dedicated metrics client.
func (s *NodeMetricsService) GetAllNodesMetrics(config *rest.Config) (*NodesMetricsResponse, error) {
	// Create Metrics API client and Kubernetes client
	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Get metrics for all nodes
	nodeMetricsList, err := metricsClientset.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics for all nodes from metrics API: %w", err)
	}

	// Get all nodes info
	nodesList, err := k8sClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes info: %w", err)
	}

	// Create a map for quick node lookup
	nodesMap := make(map[string]corev1.Node)
	for _, node := range nodesList.Items {
		nodesMap[node.Name] = node
	}

	// Get all pods to calculate requests and limits per node
	allPods, err := k8sClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all pods: %w", err)
	}

	// Group pods by node
	podsByNode := make(map[string][]corev1.Pod)
	for _, pod := range allPods.Items {
		if pod.Spec.NodeName != "" {
			podsByNode[pod.Spec.NodeName] = append(podsByNode[pod.Spec.NodeName], pod)
		}
	}

	// Process each node metrics
	nodes := make([]NodeMetrics, 0, len(nodeMetricsList.Items))
	for _, nodeMetrics := range nodeMetricsList.Items {
		node, exists := nodesMap[nodeMetrics.Name]
		if !exists {
			continue
		}

		// Get resource usage and capacity
		cpuUsage := nodeMetrics.Usage.Cpu()
		memoryUsage := nodeMetrics.Usage.Memory()
		cpuCapacity := node.Status.Capacity.Cpu()
		memoryCapacity := node.Status.Capacity.Memory()

		// Calculate requests and limits from pods on this node
		var cpuRequests, memoryRequests, cpuLimits, memoryLimits resource.Quantity
		if pods, exists := podsByNode[nodeMetrics.Name]; exists {
			for _, pod := range pods {
				if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
					continue
				}
				for _, container := range pod.Spec.Containers {
					if req := container.Resources.Requests.Cpu(); req != nil {
						cpuRequests.Add(*req)
					}
					if req := container.Resources.Requests.Memory(); req != nil {
						memoryRequests.Add(*req)
					}
					if limit := container.Resources.Limits.Cpu(); limit != nil {
						cpuLimits.Add(*limit)
					}
					if limit := container.Resources.Limits.Memory(); limit != nil {
						memoryLimits.Add(*limit)
					}
				}
			}
		}

		// Calculate percentages
		cpuPercent := calculatePercentage(cpuUsage.MilliValue(), cpuCapacity.MilliValue())
		memoryPercent := calculatePercentage(memoryUsage.Value(), memoryCapacity.Value())
		cpuRequestsPercent := calculatePercentage(cpuRequests.MilliValue(), cpuCapacity.MilliValue())
		memoryRequestsPercent := calculatePercentage(memoryRequests.Value(), memoryCapacity.Value())
		cpuLimitsPercent := calculatePercentage(cpuLimits.MilliValue(), cpuCapacity.MilliValue())
		memoryLimitsPercent := calculatePercentage(memoryLimits.Value(), memoryCapacity.Value())

		metrics := NodeMetrics{
			NodeName:              nodeMetrics.Name,
			CPUCores:              formatCPU(cpuUsage.MilliValue()),
			CPUPercent:            cpuPercent + "%",
			MemoryBytes:           formatMemory(memoryUsage.Value()),
			MemoryPercent:         memoryPercent + "%",
			CPUCapacity:           formatCPU(cpuCapacity.MilliValue()),
			MemoryCapacity:        formatMemory(memoryCapacity.Value()),
			CPURequests:           formatCPU(cpuRequests.MilliValue()),
			CPURequestsPercent:    cpuRequestsPercent + "%",
			MemoryRequests:        formatMemory(memoryRequests.Value()),
			MemoryRequestsPercent: memoryRequestsPercent + "%",
			CPULimits:             formatCPU(cpuLimits.MilliValue()),
			CPULimitsPercent:      cpuLimitsPercent + "%",
			MemoryLimits:          formatMemory(memoryLimits.Value()),
			MemoryLimitsPercent:   memoryLimitsPercent + "%",
			Timestamp:             nodeMetrics.Timestamp.Time.Format("2006-01-02T15:04:05Z"),
		}
		nodes = append(nodes, metrics)
	}

	response := &NodesMetricsResponse{
		Nodes: nodes,
		Total: len(nodes),
	}

	return response, nil
}

// calculatePercentage calculates the percentage of used resources
func calculatePercentage(used, total int64) string {
	if total == 0 {
		return "0"
	}
	percentage := (used * 100) / total
	return fmt.Sprintf("%d", percentage)
}

// formatCPU formats CPU value from milliCores to human readable format
func formatCPU(milliCores int64) string {
	if milliCores == 0 {
		return "0m"
	}
	if milliCores < 1000 {
		return fmt.Sprintf("%dm", milliCores)
	}
	cores := float64(milliCores) / 1000.0
	return fmt.Sprintf("%.1f", cores)
}

// formatMemory formats memory value from bytes to human readable format
func formatMemory(bytes int64) string {
	if bytes == 0 {
		return "0"
	}

	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.1fTi", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.1fGi", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.0fMi", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.0fKi", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d", bytes)
	}
}
