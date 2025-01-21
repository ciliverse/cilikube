package service

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
)

type Service corev1.Service

func ListServices() ([]Service, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	services, err := client.ListServices()
	if err != nil {
		return nil, err
	}
	var result []Service
	for _, svc := range services {
		result = append(result, Service(svc))
	}
	return result, nil
}

func CreateService(svc *Service) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.CreateService((*corev1.Service)(svc))
}

func GetService(name string) (*Service, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	svc, err := client.GetService(name)
	if err != nil {
		return nil, err
	}
	return (*Service)(svc), nil
}

func UpdateService(name string, svc *Service) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.UpdateService(name, (*corev1.Service)(svc))
}

func DeleteService(name string) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.DeleteService(name)
}
