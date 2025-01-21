package service

import (
	"github.com/ciliverse/cilikube/pkg/k8s"
	appsv1 "k8s.io/api/apps/v1"
)

type Deployment appsv1.Deployment

func ListDeployments() ([]Deployment, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	deployments, err := client.ListDeployments()
	if err != nil {
		return nil, err
	}
	result := make([]Deployment, len(deployments))
	for i, d := range deployments {
		result[i] = Deployment(d)
	}
	return result, nil
}

func CreateDeployment(deployment *Deployment) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.CreateDeployment((*appsv1.Deployment)(deployment))
}

func GetDeployment(name string) (*Deployment, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	deployment, err := client.GetDeployment(name)
	if err != nil {
		return nil, err
	}
	return (*Deployment)(deployment), nil
}

func UpdateDeployment(name string, deployment *Deployment) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.UpdateDeployment(name, (*appsv1.Deployment)(deployment))
}

func DeleteDeployment(name string) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return client.DeleteDeployment(name)
}
