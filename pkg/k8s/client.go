package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func NewClient() (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{clientset: clientset}, nil
}

func (c *Client) ListPods() ([]v1.Pod, error) {
	pods, err := c.clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func (c *Client) CreatePod(pod *v1.Pod) error {
	_, err := c.clientset.CoreV1().Pods(pod.Namespace).Create(pod)
	return err
}

func (c *Client) GetPod(name string) (*v1.Pod, error) {
	pod, err := c.clientset.CoreV1().Pods("").Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (c *Client) UpdatePod(name string, pod *v1.Pod) error {
	_, err := c.clientset.CoreV1().Pods(pod.Namespace).Update(pod)
	return err
}

func (c *Client) DeletePod(name string) error {
	return c.clientset.CoreV1().Pods("").Delete(name, &metav1.DeleteOptions{})
}

func (c *Client) ListDeployments() ([]appsv1.Deployment, error) {
	deployments, err := c.clientset.AppsV1().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return deployments.Items, nil
}

func (c *Client) CreateDeployment(deployment *appsv1.Deployment) error {
	_, err := c.clientset.AppsV1().Deployments(deployment.Namespace).Create(deployment)
	return err
}

func (c *Client) GetDeployment(name string) (*appsv1.Deployment, error) {
	deployment, err := c.clientset.AppsV1().Deployments("").Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (c *Client) UpdateDeployment(name string, deployment *appsv1.Deployment) error {
	_, err := c.clientset.AppsV1().Deployments(deployment.Namespace).Update(deployment)
	return err
}

func (c *Client) DeleteDeployment(name string) error {
	return c.clientset.AppsV1().Deployments("").Delete(name, &metav1.DeleteOptions{})
}

func (c *Client) ListServices() ([]v1.Service, error) {
	services, err := c.clientset.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return services.Items, nil
}

func (c *Client) CreateService(svc *v1.Service) error {
	_, err := c.clientset.CoreV1().Services(svc.Namespace).Create(svc)
	return err
}

func (c *Client) GetService(name string) (*v1.Service, error) {
	svc, err := c.clientset.CoreV1().Services("").Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func (c *Client) UpdateService(name string, svc *v1.Service) error {
	_, err := c.clientset.CoreV1().Services(svc.Namespace).Update(svc)
	return err
}

func (c *Client) DeleteService(name string) error {
	return c.clientset.CoreV1().Services("").Delete(name, &metav1.DeleteOptions{})
}
