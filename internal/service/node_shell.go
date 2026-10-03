package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	nodeShellNamespace = "kube-system"
	nodeShellImage     = "busybox:1.36"
)

// BuildNodeShellPod is a privileged host-PID pod pinned to one node.
// chroot /host lands the shell in that node's root filesystem.
func BuildNodeShellPod(nodeName, image string) (*corev1.Pod, error) {
	nodeName = strings.TrimSpace(nodeName)
	if nodeName == "" {
		return nil, fmt.Errorf("node name is required")
	}
	image = strings.TrimSpace(image)
	if image == "" {
		image = nodeShellImage
	}
	if strings.ContainsAny(image, " \t\r\n") {
		return nil, fmt.Errorf("invalid shell image")
	}
	privileged := true
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "ck-node-shell-",
			Namespace:    nodeShellNamespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "cilikube",
				"app.kubernetes.io/name":       "node-shell",
			},
		},
		Spec: corev1.PodSpec{
			NodeName:      nodeName,
			HostPID:       true,
			HostNetwork:   true,
			RestartPolicy: corev1.RestartPolicyNever,
			Tolerations:   []corev1.Toleration{{Operator: corev1.TolerationOpExists}},
			Containers: []corev1.Container{{
				Name:            "shell",
				Image:           image,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Stdin:           true,
				TTY:             true,
				Command:         []string{"sleep", "3600"},
				SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
				VolumeMounts: []corev1.VolumeMount{{
					Name:      "host-root",
					MountPath: "/host",
				}},
			}},
			Volumes: []corev1.Volume{{
				Name: "host-root",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/"},
				},
			}},
		},
	}, nil
}

// StartNodeShell creates the helper pod and waits until it is Running.
func StartNodeShell(ctx context.Context, client kubernetes.Interface, nodeName, image string) (*corev1.Pod, error) {
	pod, err := BuildNodeShellPod(nodeName, image)
	if err != nil {
		return nil, err
	}
	created, err := client.CoreV1().Pods(nodeShellNamespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	waitCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	err = wait.PollUntilContextTimeout(waitCtx, time.Second, 2*time.Minute, true, func(ctx context.Context) (bool, error) {
		current, getErr := client.CoreV1().Pods(nodeShellNamespace).Get(ctx, created.Name, metav1.GetOptions{})
		if getErr != nil {
			return false, getErr
		}
		switch current.Status.Phase {
		case corev1.PodRunning:
			return true, nil
		case corev1.PodFailed, corev1.PodSucceeded:
			return false, fmt.Errorf("node shell pod %s is %s", current.Name, current.Status.Phase)
		default:
			return false, nil
		}
	})
	if err != nil {
		_ = DeleteNodeShell(context.Background(), client, created.Name)
		return nil, fmt.Errorf("node shell pod %s did not start: %w", created.Name, err)
	}
	return created, nil
}

// DeleteNodeShell removes the helper pod. Missing pods are not an error.
func DeleteNodeShell(ctx context.Context, client kubernetes.Interface, name string) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	err := client.CoreV1().Pods(nodeShellNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}
