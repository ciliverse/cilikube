package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EventService provides business logic for cluster events
type EventService struct {
	k8sManager *k8s.ClusterManager
}

// NewEventService creates a new EventService instance
func NewEventService(k8sManager *k8s.ClusterManager) *EventService {
	return &EventService{
		k8sManager: k8sManager,
	}
}

// ListEvents retrieves cluster events based on the provided filters
func (s *EventService) ListEvents(req models.EventListRequest) (*models.EventListResponse, error) {
	client, err := s.k8sManager.GetActiveClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get active cluster client: %w", err)
	}

	// Set default limit if not specified
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 200 {
		req.Limit = 200 // Maximum limit to prevent performance issues
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var events []corev1.Event

	if req.Namespace != "" {
		// Get events from specific namespace
		eventList, err := client.Clientset.CoreV1().Events(req.Namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list events in namespace %s: %w", req.Namespace, err)
		}
		events = eventList.Items
	} else {
		// Get events from all namespaces
		eventList, err := client.Clientset.CoreV1().Events("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list events from all namespaces: %w", err)
		}
		events = eventList.Items
	}

	// Filter events based on request parameters
	filteredEvents := s.filterEvents(events, req)

	// Sort events by last timestamp (newest first)
	sort.Slice(filteredEvents, func(i, j int) bool {
		return filteredEvents[i].LastTimestamp.Time.After(filteredEvents[j].LastTimestamp.Time)
	})

	// Apply limit
	total := len(filteredEvents)
	if len(filteredEvents) > req.Limit {
		filteredEvents = filteredEvents[:req.Limit]
	}

	// Convert to response format
	clusterEvents := make([]models.ClusterEvent, len(filteredEvents))
	for i, event := range filteredEvents {
		clusterEvents[i] = models.ConvertK8sEventToClusterEvent(&event)
	}

	return &models.EventListResponse{
		Events: clusterEvents,
		Total:  total,
	}, nil
}

// GetRecentEvents retrieves the most recent cluster events (for dashboard)
func (s *EventService) GetRecentEvents(limit int) ([]models.ClusterEvent, error) {
	if limit <= 0 {
		limit = 10
	}

	req := models.EventListRequest{
		Limit: limit,
	}

	response, err := s.ListEvents(req)
	if err != nil {
		return nil, err
	}

	return response.Events, nil
}

// filterEvents applies filters to the event list
func (s *EventService) filterEvents(events []corev1.Event, req models.EventListRequest) []corev1.Event {
	var filtered []corev1.Event

	for _, event := range events {
		// Filter by event type
		if req.Type != "" {
			eventType := "Normal"
			if event.Type == corev1.EventTypeWarning {
				eventType = "Warning"
			}
			if eventType != req.Type {
				continue
			}
		}

		// Filter by time (since parameter)
		if req.Since != "" {
			sinceTime, err := time.Parse(time.RFC3339, req.Since)
			if err == nil && event.LastTimestamp.Time.Before(sinceTime) {
				continue
			}
		}

		filtered = append(filtered, event)
	}

	return filtered
}

// GetEventsByObject retrieves events related to a specific Kubernetes object
func (s *EventService) GetEventsByObject(namespace, kind, name string) ([]models.ClusterEvent, error) {
	client, err := s.k8sManager.GetActiveClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get active cluster client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get events from the specific namespace or all namespaces
	searchNamespace := namespace
	if searchNamespace == "" {
		searchNamespace = metav1.NamespaceAll
	}

	eventList, err := client.Clientset.CoreV1().Events(searchNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	var relatedEvents []corev1.Event
	for _, event := range eventList.Items {
		if event.InvolvedObject.Kind == kind && event.InvolvedObject.Name == name {
			if namespace == "" || event.InvolvedObject.Namespace == namespace {
				relatedEvents = append(relatedEvents, event)
			}
		}
	}

	// Sort by last timestamp (newest first)
	sort.Slice(relatedEvents, func(i, j int) bool {
		return relatedEvents[i].LastTimestamp.Time.After(relatedEvents[j].LastTimestamp.Time)
	})

	// Convert to response format
	clusterEvents := make([]models.ClusterEvent, len(relatedEvents))
	for i, event := range relatedEvents {
		clusterEvents[i] = models.ConvertK8sEventToClusterEvent(&event)
	}

	return clusterEvents, nil
}
