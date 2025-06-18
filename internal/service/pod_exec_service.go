package service

import (
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecOptions 执行选项
type ExecOptions struct {
	Command   []string
	Container string
	Stdin     bool
	Stdout    bool
	Stderr    bool
	TTY       bool
}

// PodExecService 处理 Pod 执行相关操作
type PodExecService struct {
	config *rest.Config
}

// NewPodExecService 创建 Pod 执行服务
func NewPodExecService(config *rest.Config) *PodExecService {
	return &PodExecService{
		config: config,
	}
}

// Exec 在 Pod 中执行命令
func (s *PodExecService) Exec(clientset kubernetes.Interface, namespace, podName string, options *ExecOptions, stdout io.Writer, stdin io.Reader) error {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: options.Container,
		Command:   options.Command,
		Stdin:     options.Stdin,
		Stdout:    options.Stdout,
		Stderr:    options.Stderr,
		TTY:       options.TTY,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(s.config, "POST", req.URL())
	if err != nil {
		return err
	}

	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stdout,
		Tty:    options.TTY,
	})
}
