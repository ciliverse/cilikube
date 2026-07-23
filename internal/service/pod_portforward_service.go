package service

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// PortForwardOptions configures a port-forward session.
type PortForwardOptions struct {
	// Ports entries like "8080:80" (local:remote) or "80" (same local/remote).
	Ports []string
}

// ParsePortPairs normalizes port query values into local:remote pairs.
func ParsePortPairs(raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("ports query is required (e.g. ports=8080:80)")
	}
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, ",") {
			for _, part := range strings.Split(p, ",") {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				normalized, err := normalizePortPair(part)
				if err != nil {
					return nil, err
				}
				out = append(out, normalized)
			}
			continue
		}
		normalized, err := normalizePortPair(p)
		if err != nil {
			return nil, err
		}
		out = append(out, normalized)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid ports provided")
	}
	return out, nil
}

func normalizePortPair(p string) (string, error) {
	if !strings.Contains(p, ":") {
		if _, err := strconv.Atoi(p); err != nil {
			return "", fmt.Errorf("invalid port %q", p)
		}
		return p + ":" + p, nil
	}
	parts := strings.Split(p, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid port pair %q (expected local:remote)", p)
	}
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return "", fmt.Errorf("invalid local port in %q", p)
	}
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return "", fmt.Errorf("invalid remote port in %q", p)
	}
	return p, nil
}

// PodPortForwardService starts server-local port forwards for a pod.
type PodPortForwardService struct{}

func NewPodPortForwardService() *PodPortForwardService {
	return &PodPortForwardService{}
}

// Forward listens on the API host (127.0.0.1) and forwards to the pod until stopCh is closed.
// readyCh is closed when forwarding is ready. out/errOut receive status lines.
func (s *PodPortForwardService) Forward(
	config *rest.Config,
	clientset kubernetes.Interface,
	namespace, podName string,
	ports []string,
	stopCh <-chan struct{},
	readyCh chan struct{},
	out, errOut io.Writer,
) error {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())
	fw, err := portforward.NewOnAddresses(dialer, []string{"127.0.0.1"}, ports, stopCh, readyCh, out, errOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

// PortForwardSession tracks an active forward for UI status messages.
type PortForwardSession struct {
	mu      sync.Mutex
	Ports   []string
	Ready   bool
	Message string
}

func (s *PortForwardSession) SetReady(ports []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Ready = true
	s.Ports = ports
	s.Message = "ready"
}

func (s *PortForwardSession) SetError(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Message = msg
}
