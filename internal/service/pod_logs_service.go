package service

import (
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodLogsService 处理 Pod 日志相关操作
type PodLogsService struct{}

// NewPodLogsService 创建 Pod 日志服务
func NewPodLogsService() *PodLogsService {
	return &PodLogsService{}
}

// Get 获取 Pod 信息
func (s *PodLogsService) Get(clientset kubernetes.Interface, namespace, name string) (*v1.Pod, error) {
	return clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

// GetPodLogs 获取 Pod 日志流
func (s *PodLogsService) GetPodLogs(clientset kubernetes.Interface, namespace, name string, opts *v1.PodLogOptions) (io.ReadCloser, error) {
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
	return req.Stream(context.Background())
}

// GetLogs 获取 Pod 日志
func (s *PodLogsService) GetLogs(clientset kubernetes.Interface, namespace, podName, container string, follow, previous bool, tailLines int64, writer io.Writer) error {
	opts := &v1.PodLogOptions{
		Container: container,
		Follow:    follow,
		Previous:  previous,
		TailLines: &tailLines,
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	reader, err := req.Stream(context.Background())
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(writer, reader)
	return err
}
