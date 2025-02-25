package service

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type DeploymentService struct {
	client kubernetes.Interface
}

func NewDeploymentService(client kubernetes.Interface) *DeploymentService {
	return &DeploymentService{client: client}
}

// 获取单个Deployment
func (s *DeploymentService) Get(namespace, name string) (*appsv1.Deployment, error) {
	return s.client.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Deployment
func (s *DeploymentService) Create(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {

	if deployment.Namespace != "" && deployment.Namespace != namespace {
		return nil, NewValidationError("deployment namespace conflicts with path parameter")
	}

	return s.client.AppsV1().Deployments(namespace).Create(
		context.TODO(),
		deployment,
		metav1.CreateOptions{},
	)
}

// 更新Deployment（包含冲突检测）
func (s *DeploymentService) Update(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return s.client.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// 删除Deployment
func (s *DeploymentService) Delete(namespace, name string) error {
	return s.client.AppsV1().Deployments(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// ListDeployments 列出所有Deployment
func (s *DeploymentService) List(namespace string) (*appsv1.DeploymentList, error) {
	return s.client.AppsV1().Deployments(namespace).List(
		context.TODO(),
		metav1.ListOptions{},
	)
}

// ListDeploymentsByLabels 根据标签过滤列出Deployment
func (s *DeploymentService) ListByLabels(namespace, selector string) (*appsv1.DeploymentList, error) {
	return s.client.AppsV1().Deployments(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
		},
	)
}

// WatchDeployments 实现Watch机制
func (s *DeploymentService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.AppsV1().Deployments(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
		},
	)
}

// PatchDeployment 实现Patch机制
// func (s *DeploymentService) Patch(namespace, name string, patchData []byte) (*appsv1.Deployment, error) {
// 	return s.client.AppsV1().Deployments(namespace).Patch(
// 		context.TODO(),
// 		name,
// 		types.StrategicMergePatchType,
// 		patchData,
// 		metav1.PatchOptions{},
// 	)
// }

// ReplaceDeployment 实现Replace机制
func (s *DeploymentService) Replace(namespace, name string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return s.client.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// RollbackDeployment 实现Deployment回滚

// ScaleDeployment 实现Deployment扩缩容
func (s *DeploymentService) Scale(namespace, name string, replicas int32) (*appsv1.Deployment, error) {
	deployment, err := s.Get(namespace, name)
	if err != nil {
		return nil, err
	}

	deployment.Spec.Replicas = &replicas
	return s.client.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// PauseDeployment 实现Deployment暂停
func (s *DeploymentService) Pause(namespace, name string) (*appsv1.Deployment, error) {
	deployment, err := s.Get(namespace, name)
	if err != nil {
		return nil, err
	}

	deployment.Spec.Paused = true
	return s.client.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// ResumeDeployment 实现Deployment恢复
func (s *DeploymentService) Resume(namespace, name string) (*appsv1.Deployment, error) {
	deployment, err := s.Get(namespace, name)
	if err != nil {
		return nil, err
	}

	deployment.Spec.Paused = false
	return s.client.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}
