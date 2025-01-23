package service

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod v1.Pod

func ListPods() ([]Pod, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	pods, err := client.ListPods()
	if err != nil {
		return nil, err
	}
	var result []Pod
	for _, pod := range pods {
		result = append(result, Pod(pod))
	}
	return result, nil
}

func CreatePod(pod *Pod) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.CreatePod((*v1.Pod)(pod))
}

func GetPod(name string) (*Pod, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	pod, err := client.GetPod(name)
	if err != nil {
		return nil, err
	}
	return (*Pod)(pod), nil
}

func UpdatePod(name string, pod *Pod) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.UpdatePod((*v1.Pod)(pod))
}

func DeletePod(name string) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	deleteOptions := metav1.DeleteOptions{}
	return client.DeletePod(name, &deleteOptions)
}
