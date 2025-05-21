package models

import (
	"fmt"
	"time" // Import time

	corev1 "k8s.io/api/core/v1"
)

// --- Request Structures (Primarily for JSON API, YAML handled separately) ---

// CreatePodRequest (Used if Content-Type is application/json)
type CreatePodRequest struct {
	Name string `json:"name" binding:"required"`
	// Namespace is typically taken from the URL path, not the body for POST
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        corev1.PodSpec    `json:"spec" binding:"required"`
}

// UpdatePodRequest (Used if Content-Type is application/json)
// Note: This replaces the entire spec, labels, annotations. Use with caution or prefer PATCH/YAML.
type UpdatePodRequest struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        corev1.PodSpec    `json:"spec" binding:"required"`
}

// --- Response Structures ---

// ContainerNameResponse 提供容器状态信息给前端,目前只展示日志接口部分使用的字段
type ContainerNameResponse struct {
	Name string `json:"name"`
}

type PodSpecResponse struct {
	Containers     []ContainerNameResponse `json:"containers"`
	InitContainers []ContainerNameResponse `json:"initContainers"`
}

// PodResponse represents the data sent back to the client for a single Pod.
type PodResponse struct {
	UID         string            `json:"uid"` // Added UID
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Status      string            `json:"status"`            // Phase: Pending, Running, Succeeded, Failed, Unknown
	Reason      string            `json:"reason,omitempty"`  // Added Reason (e.g., Evicted)
	Message     string            `json:"message,omitempty"` // Added Message (more details on status)
	IP          string            `json:"ip,omitempty"`      // Pod IP
	Node        string            `json:"node,omitempty"`    // Node name where the pod is scheduled
	CreatedAt   string            `json:"createdAt"`         // Formatted timestamp string
	Spec        *PodSpecResponse  `json:"spec,omitempty"`    // 根据前端解析需求传参

	// Add Container Statuses if needed by frontend
	//ContainerStatuses []ContainerStatusResponse `json:"spec,omitempty"`
}

// PodListResponse represents the paginated list of Pods.
type PodListResponse struct {
	Items []PodResponse `json:"items"`
	Total int           `json:"total"` // Total items matching filter (in this batch, see handler note)
}

// ToPodResponse converts a Kubernetes Pod object to our API response format.
func ToPodResponse(pod *corev1.Pod) PodResponse {
	// Format timestamp for better readability in frontend
	createdAtFormatted := "N/A"
	if !pod.CreationTimestamp.IsZero() {
		createdAtFormatted = pod.CreationTimestamp.Format(time.RFC3339) // Use standard format
	}

	// Determine a more detailed status if possible
	status := string(pod.Status.Phase)
	reason := pod.Status.Reason
	message := pod.Status.Message

	// Try to get a better status from conditions or container statuses if phase is Pending/Failed/Unknown
	if pod.Status.Phase == corev1.PodPending || pod.Status.Phase == corev1.PodFailed || pod.Status.Phase == corev1.PodUnknown {
		// Check conditions
		for _, cond := range pod.Status.Conditions {
			if cond.Type == corev1.PodScheduled && cond.Status == corev1.ConditionFalse {
				if reason == "" {
					reason = cond.Reason
				}
				if message == "" {
					message = cond.Message
				}
				break
			}
			if cond.Type == corev1.PodReasonUnschedulable {
				if reason == "" {
					reason = cond.Reason
				}
				if message == "" {
					message = cond.Message
				}
				break
			}
		}
		// Check container statuses for waiting reasons like ImagePullBackOff, ContainerCreating, CrashLoopBackOff
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Waiting != nil {
				if reason == "" {
					reason = cs.State.Waiting.Reason
				}
				if message == "" {
					message = cs.State.Waiting.Message
				}
				// You might want to prioritize certain reasons like CrashLoopBackOff
				if reason == "CrashLoopBackOff" || reason == "ImagePullBackOff" || reason == "ErrImagePull" {
					status = cs.State.Waiting.Reason // Override phase with a more specific waiting reason
				}
				break // Take the first waiting reason
			}
			if cs.State.Terminated != nil {
				if reason == "" {
					reason = cs.State.Terminated.Reason
				}
				if message == "" {
					message = cs.State.Terminated.Message
				}
				if status != string(corev1.PodFailed) && cs.State.Terminated.ExitCode != 0 {
					status = "Error" // Or use Terminated Reason
					if reason == "" {
						reason = "Terminated"
					}
				}
				break // Take the first terminated reason
			}
		}
		// Check init container statuses too
		for _, cs := range pod.Status.InitContainerStatuses {
			if cs.State.Waiting != nil {
				if reason == "" {
					reason = fmt.Sprintf("Init: %s", cs.State.Waiting.Reason)
				} // Prefix with Init:
				if message == "" {
					message = cs.State.Waiting.Message
				}
				if cs.State.Waiting.Reason != "PodInitializing" { // PodInitializing is normal
					status = fmt.Sprintf("Init: %s", cs.State.Waiting.Reason)
				}
				break
			}
			if cs.State.Terminated != nil {
				if reason == "" {
					reason = fmt.Sprintf("Init: %s", cs.State.Terminated.Reason)
				}
				if message == "" {
					message = cs.State.Terminated.Message
				}
				if cs.State.Terminated.ExitCode != 0 {
					status = fmt.Sprintf("Init: %s", cs.State.Terminated.Reason) // Or "Init:Error"
					if reason == "" {
						reason = "Init:Terminated"
					}
				}
				break
			}
		}
	}

	// Ensure status has a value
	if status == "" {
		status = "Unknown"
	}
	// 正式容器和初始化容器信息
	var containers []ContainerNameResponse
	for _, c := range pod.Spec.Containers {
		containers = append(containers, ContainerNameResponse{
			Name: c.Name,
		})
	}
	var initContainers []ContainerNameResponse
	for _, c := range pod.Spec.InitContainers {
		initContainers = append(initContainers, ContainerNameResponse{
			Name: c.Name,
		})
	}

	return PodResponse{
		UID:         string(pod.UID), // Include UID
		Name:        pod.Name,
		Namespace:   pod.Namespace,
		Labels:      pod.Labels,
		Annotations: pod.Annotations,
		Status:      status,  // Use the potentially refined status
		Reason:      reason,  // Include reason
		Message:     message, // Include message
		IP:          pod.Status.PodIP,
		Node:        pod.Spec.NodeName,
		CreatedAt:   createdAtFormatted, // Use formatted string
		Spec: &PodSpecResponse{
			Containers:     containers,
			InitContainers: initContainers,
		},
	}
}
