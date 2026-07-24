package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/pkg/k8s"
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
	if k8s.IsShowcase() {
		return io.NopCloser(strings.NewReader(showcasePodLogText(namespace, name))), nil
	}
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
	stream, err := req.Stream(context.Background())
	if err != nil {
		return nil, err
	}
	// Automatically detect and convert GBK -> UTF-8
	return ConvertIfGBK(stream), nil
}

func showcasePodLogText(namespace, name string) string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return fmt.Sprintf(`%s INFO  starting container (showcase)
%s INFO  listening on :8080
%s INFO  connected to demo postgres
%s INFO  ready to serve traffic
%s WARN  showcase mode — logs are simulated; no real pod on this host
%s INFO  request id=demo-1 path=/health status=200
%s INFO  request id=demo-2 path=/api/v1/ping status=200
`, ts, ts, ts, ts, ts, ts, ts) + fmt.Sprintf("# namespace=%s pod=%s\n", namespace, name)
}

// GetLogs retrieves Pod logs
func (s *PodLogsService) GetLogs(clientset kubernetes.Interface, namespace, podName, container string, follow, previous bool, tailLines int64, writer io.Writer) error {
	if k8s.IsShowcase() {
		_, err := io.WriteString(writer, showcasePodLogText(namespace, podName))
		return err
	}
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
