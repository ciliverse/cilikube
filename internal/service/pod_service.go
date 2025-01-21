package service

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	v1 "k8s.io/api/core/v1"
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
	for _, p := range pods {
		result = append(result, Pod(p))
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
	return client.UpdatePod(name, (*v1.Pod)(pod))
}

func DeletePod(name string) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.DeletePod(name)
}
