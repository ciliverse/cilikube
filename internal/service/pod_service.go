package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type PodService struct {
	client kubernetes.Interface
}

func NewPodService(client kubernetes.Interface) *PodService {
	return &PodService{client: client}
}

func (s *PodService) Get(namespace, name string) (*corev1.Pod, error) {
	return s.client.CoreV1().Pods(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

func (s *PodService) Create(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	return s.client.CoreV1().Pods(namespace).Create(
		context.TODO(),
		pod,
		metav1.CreateOptions{},
	)
}

func (s *PodService) Update(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	return s.client.CoreV1().Pods(namespace).Update(
		context.TODO(),
		pod,
		metav1.UpdateOptions{},
	)
}

func (s *PodService) Delete(namespace, name string) error {
	return s.client.CoreV1().Pods(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

func (s *PodService) List(namespace, selector string, limit int64) (*corev1.PodList, error) {
	return s.client.CoreV1().Pods(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

func (s *PodService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.CoreV1().Pods(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}

func int64ptr(i int64) *int64 { return &i }
