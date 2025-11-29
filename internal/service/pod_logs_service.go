package service

import (
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodLogsService handles Pod logs related operations
type PodLogsService struct{}

// NewPodLogsService creates Pod logs service
func NewPodLogsService() *PodLogsService {
	return &PodLogsService{}
}

// Get retrieves Pod information
func (s *PodLogsService) Get(clientset kubernetes.Interface, namespace, name string) (*v1.Pod, error) {
	return clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

// GetPodLogs retrieves Pod log stream
func (s *PodLogsService) GetPodLogs(clientset kubernetes.Interface, namespace, name string, opts *v1.PodLogOptions) (io.ReadCloser, error) {
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
	stream, err := req.Stream(context.Background())
	if err != nil {
		return nil, err
	}
	// Automatically detect and convert GBK -> UTF-8
	return ConvertIfGBK(stream), nil
}

// GetLogs retrieves Pod logs
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
