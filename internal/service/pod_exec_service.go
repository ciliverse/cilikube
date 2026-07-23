package service

import (
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecOptions execution options
type ExecOptions struct {
	Command   []string
	Container string
	Stdin     bool
	Stdout    bool
	Stderr    bool
	TTY       bool
}

// PodExecService handles Pod execution related operations
type PodExecService struct{}

// NewPodExecService creates Pod execution service
func NewPodExecService() *PodExecService {
	return &PodExecService{}
}

// Exec executes command in Pod using the request cluster's rest.Config
func (s *PodExecService) Exec(config *rest.Config, clientset kubernetes.Interface, namespace, podName string, options *ExecOptions, stdout io.Writer, stdin io.Reader) error {
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

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
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

// Attach attaches to a running container process using the request cluster's rest.Config
func (s *PodExecService) Attach(config *rest.Config, clientset kubernetes.Interface, namespace, podName string, options *ExecOptions, stdout io.Writer, stdin io.Reader) error {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("attach")

	req.VersionedParams(&corev1.PodAttachOptions{
		Container: options.Container,
		Stdin:     options.Stdin,
		Stdout:    options.Stdout,
		Stderr:    options.Stderr,
		TTY:       options.TTY,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	return executor.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stdout,
		Tty:    options.TTY,
	})
}
