package service

import (
	"fmt"
	"math"
	"time"

	"k8s.io/client-go/rest"
)

// showcaseNodeMetrics returns synthetic, gently oscillating node metrics for the public demo.
func showcaseNodeMetrics(config *rest.Config) (*NodesMetricsResponse, error) {
	_ = config
	t := float64(time.Now().Unix())
	nodes := []string{"demo-master-1", "demo-worker-1", "demo-worker-2"}
	capsCPU := []float64{4, 8, 8}
	capsMemGi := []float64{8, 16, 16}

	out := make([]NodeMetrics, 0, len(nodes))
	for i, name := range nodes {
		cpuPct := 18 + 12*math.Sin(t/7+float64(i)) + 6*math.Sin(t/3+float64(i)*1.7)
		if cpuPct < 5 {
			cpuPct = 5
		}
		if cpuPct > 92 {
			cpuPct = 92
		}
		memPct := 42 + 10*math.Cos(t/11+float64(i)) + 4*math.Sin(t/5)
		if memPct < 20 {
			memPct = 20
		}
		if memPct > 88 {
			memPct = 88
		}
		cpuCores := capsCPU[i] * cpuPct / 100
		memMi := capsMemGi[i] * 1024 * memPct / 100
		out = append(out, NodeMetrics{
			NodeName:              name,
			CPUCores:              fmt.Sprintf("%.0fm", cpuCores*1000),
			CPUPercent:            fmt.Sprintf("%.0f%%", cpuPct),
			MemoryBytes:           fmt.Sprintf("%.0fMi", memMi),
			MemoryPercent:         fmt.Sprintf("%.0f%%", memPct),
			CPUCapacity:           fmt.Sprintf("%.0f", capsCPU[i]),
			MemoryCapacity:        fmt.Sprintf("%.0fGi", capsMemGi[i]),
			CPURequests:           fmt.Sprintf("%.0fm", capsCPU[i]*1000*0.35),
			CPURequestsPercent:    "35%",
			MemoryRequests:        fmt.Sprintf("%.0fMi", capsMemGi[i]*1024*0.4),
			MemoryRequestsPercent: "40%",
			CPULimits:             fmt.Sprintf("%.0fm", capsCPU[i]*1000*0.8),
			CPULimitsPercent:      "80%",
			MemoryLimits:          fmt.Sprintf("%.0fMi", capsMemGi[i]*1024*0.85),
			MemoryLimitsPercent:   "85%",
			Timestamp:             time.Now().UTC().Format(time.RFC3339),
		})
	}
	return &NodesMetricsResponse{Nodes: out, Total: len(out)}, nil
}

func showcaseSingleNodeMetrics(config *rest.Config, nodeName string) (*NodeMetrics, error) {
	all, err := showcaseNodeMetrics(config)
	if err != nil {
		return nil, err
	}
	for i := range all.Nodes {
		if all.Nodes[i].NodeName == nodeName {
			return &all.Nodes[i], nil
		}
	}
	return nil, fmt.Errorf("showcase node %q not found", nodeName)
}

func showcasePodMetrics(namespace string) *PodsMetricsResponse {
	t := float64(time.Now().Unix())
	type seed struct {
		ns, name string
		cpuBase  float64
		memMi    float64
	}
	seeds := []seed{
		{"default", "web-frontend-7d9f8b-abc12", 45, 96},
		{"default", "web-frontend-7d9f8b-def34", 52, 110},
		{"default", "api-gateway-6c4d5-jkl78", 80, 128},
		{"production", "orders-api-8f2a1-aa111", 120, 180},
		{"production", "orders-api-8f2a1-bb222", 95, 160},
		{"production", "postgres-0", 70, 512},
		{"production", "postgres-1", 65, 480},
		{"monitoring", "prometheus-0abc1", 200, 640},
		{"kube-system", "coredns-xyz01", 15, 48},
		{"cilibase", "cilikube-demo-1a2b3", 40, 90},
	}
	items := make([]PodMetrics, 0, len(seeds))
	for i, s := range seeds {
		if namespace != "" && s.ns != namespace {
			continue
		}
		cpu := s.cpuBase + 20*math.Sin(t/5+float64(i))
		if cpu < 5 {
			cpu = 5
		}
		mem := s.memMi + 30*math.Cos(t/8+float64(i))
		items = append(items, PodMetrics{
			Namespace:            s.ns,
			Name:                 s.name,
			CPU:                  fmt.Sprintf("%.0fm", cpu),
			Memory:               fmt.Sprintf("%.0fMi", mem),
			CPUMilli:             int64(cpu),
			MemoryBytes:          int64(mem * 1024 * 1024),
			CPURequest:           "50m",
			CPULimit:             "500m",
			MemoryRequest:        "64Mi",
			MemoryLimit:          "256Mi",
			CPURequestPercent:    fmt.Sprintf("%.0f%%", cpu/50*100),
			CPULimitPercent:      fmt.Sprintf("%.0f%%", cpu/500*100),
			MemoryRequestPercent: fmt.Sprintf("%.0f%%", mem/64*100),
			MemoryLimitPercent:   fmt.Sprintf("%.0f%%", mem/256*100),
			CPURequestRatio:      cpu / 50,
			CPULimitRatio:        cpu / 500,
			MemoryRequestRatio:   mem / 64,
			MemoryLimitRatio:     mem / 256,
			Timestamp:            time.Now().UTC().Format(time.RFC3339),
		})
	}
	return &PodsMetricsResponse{
		Pods:      items,
		Total:     len(items),
		Available: true,
		Message:   "showcase simulated metrics",
	}
}
