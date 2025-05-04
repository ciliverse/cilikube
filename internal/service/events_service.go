package service

import (
	"context"
	"github.com/ciliverse/cilikube/api/v1/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

type EventsService struct {
	client kubernetes.Interface
}

func NewEventsService(client kubernetes.Interface) *EventsService {
	return &EventsService{
		client: client,
	}
}

func (s *EventsService) List(namespace string) *models.EventList {
	results := &models.EventList{
		Items: []models.Event{},
		Total: 0,
	}
	events, err := s.client.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("获取事件列表失败：%s", err)
		return results
	}
	for _, event := range events.Items {
		results.Items = append(results.Items, models.K8sEventToEvent(&event))
	}
	results.Total = len(results.Items)
	return results
}

func (s *EventsService) Get(namespace, name string) models.Event {
	event, err := s.client.CoreV1().Events(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Printf("获取事件失败：%s", err)
		return models.Event{}
	}
	return models.K8sEventToEvent(event)
}
