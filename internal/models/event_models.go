package models

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// ClusterEvent represents a Kubernetes cluster event
type ClusterEvent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Namespace  string    `json:"namespace"`
	Reason     string    `json:"reason"`
	Message    string    `json:"message"`
	Type       string    `json:"type"`       // Normal, Warning
	Source     string    `json:"source"`     // Component that reported this event
	Object     string    `json:"object"`     // Involved object (Pod, Service, etc.)
	ObjectKind string    `json:"objectKind"` // Kind of the involved object
	Count      int32     `json:"count"`      // Number of times this event has occurred
	FirstTime  time.Time `json:"firstTime"`  // First time this event was observed
	LastTime   time.Time `json:"lastTime"`   // Last time this event was observed
	CreatedAt  time.Time `json:"createdAt"`  // Event creation time
}

// EventListRequest represents the request parameters for listing events
type EventListRequest struct {
	Namespace string `form:"namespace" json:"namespace"` // Filter by namespace, empty means all namespaces
	Type      string `form:"type" json:"type"`           // Filter by event type (Normal, Warning)
	Limit     int    `form:"limit" json:"limit"`         // Limit number of events returned
	Since     string `form:"since" json:"since"`         // Filter events since this time (RFC3339 format)
}

// EventListResponse represents the response for event listing
type EventListResponse struct {
	Events []ClusterEvent `json:"events"`
	Total  int            `json:"total"`
}

// ConvertK8sEventToClusterEvent converts Kubernetes Event to ClusterEvent
func ConvertK8sEventToClusterEvent(event *corev1.Event) ClusterEvent {
	eventType := "Normal"
	if event.Type == corev1.EventTypeWarning {
		eventType = "Warning"
	}

	objectName := ""
	objectKind := ""
	if event.InvolvedObject.Name != "" {
		objectName = event.InvolvedObject.Name
		objectKind = event.InvolvedObject.Kind
	}

	source := ""
	if event.Source.Component != "" {
		source = event.Source.Component
	} else if event.ReportingController != "" {
		source = event.ReportingController
	}

	return ClusterEvent{
		ID:         string(event.UID),
		Name:       event.Name,
		Namespace:  event.Namespace,
		Reason:     event.Reason,
		Message:    event.Message,
		Type:       eventType,
		Source:     source,
		Object:     objectName,
		ObjectKind: objectKind,
		Count:      event.Count,
		FirstTime:  event.FirstTimestamp.Time,
		LastTime:   event.LastTimestamp.Time,
		CreatedAt:  event.CreationTimestamp.Time,
	}
}
