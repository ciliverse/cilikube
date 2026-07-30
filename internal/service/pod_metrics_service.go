package service

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// PodMetrics is a pod-level usage snapshot (k9s-style: usage + % of request/limit).
type PodMetrics struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	CPU         string `json:"cpu"`
	Memory      string `json:"memory"`
	CPUMilli    int64  `json:"cpuMilli"`
	MemoryBytes int64  `json:"memoryBytes"`

	CPURequest         string `json:"cpuRequest,omitempty"`
	CPULimit           string `json:"cpuLimit,omitempty"`
	MemoryRequest      string `json:"memoryRequest,omitempty"`
	MemoryLimit        string `json:"memoryLimit,omitempty"`
	CPURequestPercent  string `json:"cpuRequestPercent"`  // e.g. "45%" or "-"
	CPULimitPercent    string `json:"cpuLimitPercent"`
	MemoryRequestPercent string `json:"memoryRequestPercent"`
	MemoryLimitPercent   string `json:"memoryLimitPercent"`

	CPURequestRatio    float64 `json:"cpuRequestRatio"`    // 0..n for UI bars/sort
	CPULimitRatio      float64 `json:"cpuLimitRatio"`
	MemoryRequestRatio float64 `json:"memoryRequestRatio"`
	MemoryLimitRatio   float64 `json:"memoryLimitRatio"`

	Timestamp string `json:"timestamp,omitempty"`
}

// PodsMetricsResponse lists pod metrics; Available=false when metrics-server is missing.
type PodsMetricsResponse struct {
	Pods      []PodMetrics `json:"pods"`
	Total     int          `json:"total"`
	Available bool         `json:"available"`
	Message   string       `json:"message,omitempty"`
}

// PodMetricsService reads PodMetrics from metrics.k8s.io + pod specs for request/limit %.
type PodMetricsService struct{}

func NewPodMetricsService() *PodMetricsService {
	return &PodMetricsService{}
}

// ListPodMetrics returns usage for pods. Empty namespace lists all namespaces.
func (s *PodMetricsService) ListPodMetrics(config *rest.Config, namespace string) *PodsMetricsResponse {
	out := &PodsMetricsResponse{Pods: []PodMetrics{}, Available: false}
	if config == nil {
		out.Message = "no cluster config"
		return out
	}
	if k8s.IsShowcaseConfig(config) {
		return showcasePodMetrics(config, namespace)
	}

	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		out.Message = fmt.Sprintf("failed to create metrics client: %v", err)
		return out
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		out.Message = fmt.Sprintf("failed to create kubernetes client: %v", err)
		return out
	}

	ctx := context.Background()
	pm, err := metricsClientset.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		out.Message = friendlyMetricsErr(err)
		return out
	}

	// Pod specs for requests/limits (same scope as metrics list)
	podList, err := k8sClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		out.Message = fmt.Sprintf("failed to list pods for request/limit: %v", err)
		return out
	}
	podByKey := make(map[string]*corev1.Pod, len(podList.Items))
	for i := range podList.Items {
		p := &podList.Items[i]
		podByKey[p.Namespace+"/"+p.Name] = p
	}

	items := make([]PodMetrics, 0, len(pm.Items))
	for i := range pm.Items {
		m := &pm.Items[i]
		item := sumPodContainerMetrics(m)
		if pod := podByKey[m.Namespace+"/"+m.Name]; pod != nil {
			enrichWithRequestsLimits(&item, pod)
		} else {
			item.CPURequestPercent = "-"
			item.CPULimitPercent = "-"
			item.MemoryRequestPercent = "-"
			item.MemoryLimitPercent = "-"
		}
		items = append(items, item)
	}
	out.Available = true
	out.Pods = items
	out.Total = len(items)
	return out
}

func sumPodContainerMetrics(m *metricsv1beta1.PodMetrics) PodMetrics {
	var cpuMilli, memBytes int64
	for _, c := range m.Containers {
		if q := c.Usage.Cpu(); q != nil {
			cpuMilli += q.MilliValue()
		}
		if q := c.Usage.Memory(); q != nil {
			memBytes += q.Value()
		}
	}
	return PodMetrics{
		Namespace:   m.Namespace,
		Name:        m.Name,
		CPU:         formatCPU(cpuMilli),
		Memory:      formatMemory(memBytes),
		CPUMilli:    cpuMilli,
		MemoryBytes: memBytes,
		Timestamp:   m.Timestamp.Time.Format("2006-01-02T15:04:05Z"),
	}
}

func enrichWithRequestsLimits(item *PodMetrics, pod *corev1.Pod) {
	var cpuReq, cpuLim, memReq, memLim resource.Quantity
	for _, c := range pod.Spec.Containers {
		if q := c.Resources.Requests.Cpu(); q != nil {
			cpuReq.Add(*q)
		}
		if q := c.Resources.Limits.Cpu(); q != nil {
			cpuLim.Add(*q)
		}
		if q := c.Resources.Requests.Memory(); q != nil {
			memReq.Add(*q)
		}
		if q := c.Resources.Limits.Memory(); q != nil {
			memLim.Add(*q)
		}
	}

	if cpuReq.MilliValue() > 0 {
		item.CPURequest = formatCPU(cpuReq.MilliValue())
		item.CPURequestRatio = float64(item.CPUMilli) / float64(cpuReq.MilliValue())
		item.CPURequestPercent = formatRatioPercent(item.CPURequestRatio)
	} else {
		item.CPURequestPercent = "-"
	}
	if cpuLim.MilliValue() > 0 {
		item.CPULimit = formatCPU(cpuLim.MilliValue())
		item.CPULimitRatio = float64(item.CPUMilli) / float64(cpuLim.MilliValue())
		item.CPULimitPercent = formatRatioPercent(item.CPULimitRatio)
	} else {
		item.CPULimitPercent = "-"
	}
	if memReq.Value() > 0 {
		item.MemoryRequest = formatMemory(memReq.Value())
		item.MemoryRequestRatio = float64(item.MemoryBytes) / float64(memReq.Value())
		item.MemoryRequestPercent = formatRatioPercent(item.MemoryRequestRatio)
	} else {
		item.MemoryRequestPercent = "-"
	}
	if memLim.Value() > 0 {
		item.MemoryLimit = formatMemory(memLim.Value())
		item.MemoryLimitRatio = float64(item.MemoryBytes) / float64(memLim.Value())
		item.MemoryLimitPercent = formatRatioPercent(item.MemoryLimitRatio)
	} else {
		item.MemoryLimitPercent = "-"
	}
}

func formatRatioPercent(ratio float64) string {
	if ratio < 0 {
		ratio = 0
	}
	pct := int(ratio*100 + 0.5)
	return fmt.Sprintf("%d%%", pct)
}

func friendlyMetricsErr(err error) string {
	msg := err.Error()
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "metrics.k8s.io") ||
		strings.Contains(lower, "could not find the requested resource") ||
		strings.Contains(lower, "not found") {
		return "metrics-server not available in this cluster (metrics.k8s.io missing). Install metrics-server to show CPU/MEM."
	}
	return msg
}
