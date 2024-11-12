package services

import (
	"context"
	"log"

	"github.com/cilliandevops/kopsadmin/server-go/internal/apis/models/k8s"
	"github.com/cilliandevops/kopsadmin/server-go/internal/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodService struct {
	client *client.K8sClient
}

func NewPodService(c *client.K8sClient) *PodService {
	return &PodService{client: c}
}

func (ps *PodService) ListPods(namespace string) ([]k8s.Pod, error) {
	podList, err := ps.client.K8sClient().CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error listing pods in namespace %s: %v", namespace, err)
		return nil, err
	}

	var pods []k8s.Pod
	for _, pod := range podList.Items {
		pods = append(pods, k8s.Pod{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Labels:       pod.Labels,
			Annotations:  pod.Annotations,
			Status:       pod.Status,
			CreationTime: pod.CreationTimestamp,
		})
	}

	return pods, nil
}

func (ps *PodService) GetPod(namespace, podName string) (*k8s.Pod, error) {
	pod, err := ps.client.K8sClient().CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting pod %s in namespace %s: %v", podName, namespace, err)
		return nil, err
	}

	return &k8s.Pod{
		Name:         pod.Name,
		Namespace:    pod.Namespace,
		Labels:       pod.Labels,
		Annotations:  pod.Annotations,
		Status:       pod.Status,
		CreationTime: pod.CreationTimestamp,
	}, nil
}

func (ps *PodService) CreatePod(namespace string, pod *corev1.Pod) (*k8s.Pod, error) {
	createdPod, err := ps.client.K8sClient().CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating pod in namespace %s: %v", namespace, err)
		return nil, err
	}

	return &k8s.Pod{
		Name:         createdPod.Name,
		Namespace:    createdPod.Namespace,
		Labels:       createdPod.Labels,
		Annotations:  createdPod.Annotations,
		Status:       createdPod.Status,
		CreationTime: createdPod.CreationTimestamp,
	}, nil
}

func (ps *PodService) DeletePod(namespace, podName string) error {
	err := ps.client.K8sClient().CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Error deleting pod %s in namespace %s: %v", podName, namespace, err)
	}
	return err
}
