package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

const (
	maxScanPods    = 20
	maxScanDeploys = 10
	maxScanNodes   = 10
	maxScanEvents  = 8
	maxMetricLines = 40
	maxPromQLLen   = 800
	maxPromResult  = 6000
)

func toolScanHealth(ctx context.Context, client *k8s.Client, args map[string]interface{}, zh bool) (toolResult, error) {
	if client == nil || client.Clientset == nil {
		if zh {
			return toolResult{}, fmt.Errorf("集群客户端不可用")
		}
		return toolResult{}, fmt.Errorf("cluster client unavailable")
	}
	ns := strArg(args, "namespace")

	pods, err := client.Clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return toolResult{}, err
	}
	deploys, err := client.Clientset.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return toolResult{}, err
	}
	nodes, err := client.Clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return toolResult{}, err
	}
	events, err := client.Clientset.CoreV1().Events(ns).List(ctx, metav1.ListOptions{Limit: 500})
	if err != nil {
		return toolResult{}, err
	}

	var lines []string
	var refs []ResourceRef
	if zh {
		if ns == "" {
			lines = append(lines, "健康扫描 · 全集群（只读）")
		} else {
			lines = append(lines, "健康扫描 · 命名空间="+ns+"（节点仍看全集群，只读）")
		}
	} else {
		if ns == "" {
			lines = append(lines, "health scan · cluster-wide (read-only)")
		} else {
			lines = append(lines, "health scan · namespace="+ns+" (nodes stay cluster-wide, read-only)")
		}
	}

	podHits := 0
	for i := range pods.Items {
		p := &pods.Items[i]
		reason, detail, bad := podIssue(p)
		if !bad {
			continue
		}
		if podHits >= maxScanPods {
			lines = append(lines, scanMore(zh, "pods"))
			break
		}
		podHits++
		msg := strings.TrimSpace(detail)
		if len(msg) > 180 {
			msg = msg[:180] + "…"
		}
		if msg != "" {
			lines = append(lines, fmt.Sprintf("pod %s/%s %s — %s", p.Namespace, p.Name, reason, msg))
		} else {
			lines = append(lines, fmt.Sprintf("pod %s/%s %s", p.Namespace, p.Name, reason))
		}
		refs = append(refs, podRef(p.Namespace, p.Name, ""))
	}
	if podHits == 0 {
		if zh {
			lines = append(lines, "pod：没有 CrashLoop / ImagePull / Pending / Failed")
		} else {
			lines = append(lines, "pods: no CrashLoop, ImagePull, Pending, or Failed")
		}
	}

	deployHits := 0
	for i := range deploys.Items {
		d := &deploys.Items[i]
		if !deploymentUnhealthy(d) {
			continue
		}
		if deployHits >= maxScanDeploys {
			lines = append(lines, scanMore(zh, "deployments"))
			break
		}
		deployHits++
		want := int32(0)
		if d.Spec.Replicas != nil {
			want = *d.Spec.Replicas
		}
		lines = append(lines, fmt.Sprintf("deployment %s/%s ready=%d/%d unavailable=%d", d.Namespace, d.Name, d.Status.ReadyReplicas, want, d.Status.UnavailableReplicas))
		refs = append(refs, resourceRef("deployments", d.Namespace, d.Name, true, ""))
	}
	if deployHits == 0 {
		if zh {
			lines = append(lines, "deployment：副本都就绪")
		} else {
			lines = append(lines, "deployments: replicas ready")
		}
	}

	nodeHits := 0
	for i := range nodes.Items {
		n := &nodes.Items[i]
		ready, pressures := nodeHealth(n)
		if ready && len(pressures) == 0 {
			continue
		}
		if nodeHits >= maxScanNodes {
			lines = append(lines, scanMore(zh, "nodes"))
			break
		}
		nodeHits++
		state := "NotReady"
		if ready {
			state = "Ready"
		}
		extra := ""
		if len(pressures) > 0 {
			extra = " " + strings.Join(pressures, ",")
		}
		lines = append(lines, fmt.Sprintf("node %s %s%s", n.Name, state, extra))
		refs = append(refs, resourceRef("nodes", "", n.Name, false, ""))
	}
	if nodeHits == 0 {
		if zh {
			lines = append(lines, "node：全部 Ready，没有压力状态")
		} else {
			lines = append(lines, "nodes: all Ready, no pressure")
		}
	}

	warnings := warningEvents(events.Items)
	if len(warnings) == 0 {
		if zh {
			lines = append(lines, "events：最近一页里没有 Warning")
		} else {
			lines = append(lines, "events: no Warning in the sampled page")
		}
	} else {
		for i, e := range warnings {
			if i >= maxScanEvents {
				lines = append(lines, scanMore(zh, "events"))
				break
			}
			msg := strings.TrimSpace(e.Message)
			if len(msg) > 180 {
				msg = msg[:180] + "…"
			}
			lines = append(lines, fmt.Sprintf("event [%s] %s/%s %s: %s", e.Type, e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, msg))
			if strings.EqualFold(e.InvolvedObject.Kind, "Pod") && e.InvolvedObject.Name != "" {
				refs = append(refs, podRef(e.InvolvedObject.Namespace, e.InvolvedObject.Name, ""))
			}
		}
	}

	return toolResult{Text: strings.Join(lines, "\n"), Resources: refs}, nil
}

func scanMore(zh bool, kind string) string {
	if zh {
		return "… " + kind + " 还有更多，已截断"
	}
	return "... more " + kind + " truncated"
}

func podIssue(p *corev1.Pod) (reason, detail string, bad bool) {
	if p.Status.Phase == corev1.PodSucceeded {
		return "", "", false
	}
	statuses := append([]corev1.ContainerStatus{}, p.Status.InitContainerStatuses...)
	statuses = append(statuses, p.Status.ContainerStatuses...)
	for _, c := range statuses {
		if c.State.Waiting != nil {
			r := c.State.Waiting.Reason
			if r != "" && r != "ContainerCreating" && r != "PodInitializing" {
				return r, c.State.Waiting.Message, true
			}
		}
		if c.State.Terminated != nil {
			r := c.State.Terminated.Reason
			if r != "" && r != "Completed" {
				return r, c.State.Terminated.Message, true
			}
		}
	}
	if p.Status.Phase == corev1.PodFailed || p.Status.Phase == corev1.PodUnknown {
		return string(p.Status.Phase), p.Status.Message, true
	}
	if p.Status.Phase == corev1.PodPending {
		for _, c := range p.Status.Conditions {
			if c.Type == corev1.PodScheduled && c.Status == corev1.ConditionFalse {
				msg := c.Message
				if msg == "" {
					msg = c.Reason
				}
				return "Unschedulable", msg, true
			}
		}
		return "Pending", p.Status.Message, true
	}
	return "", "", false
}

func deploymentUnhealthy(d *appsv1.Deployment) bool {
	if d.Status.UnavailableReplicas > 0 {
		return true
	}
	if d.Spec.Replicas != nil && d.Status.ReadyReplicas < *d.Spec.Replicas {
		return true
	}
	return false
}

func nodeHealth(n *corev1.Node) (ready bool, pressures []string) {
	seenReady := false
	for _, c := range n.Status.Conditions {
		switch c.Type {
		case corev1.NodeReady:
			seenReady = true
			ready = c.Status == corev1.ConditionTrue
		case corev1.NodeMemoryPressure, corev1.NodeDiskPressure, corev1.NodePIDPressure, corev1.NodeNetworkUnavailable:
			if c.Status == corev1.ConditionTrue {
				pressures = append(pressures, string(c.Type))
			}
		}
	}
	if !seenReady {
		ready = false
	}
	return ready, pressures
}

func warningEvents(items []corev1.Event) []corev1.Event {
	var out []corev1.Event
	for _, e := range items {
		if e.Type == corev1.EventTypeWarning || strings.EqualFold(e.Type, "Error") {
			out = append(out, e)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return eventTime(out[i]).After(eventTime(out[j]))
	})
	return out
}

func eventTime(e corev1.Event) time.Time {
	if !e.LastTimestamp.IsZero() {
		return e.LastTimestamp.Time
	}
	if !e.EventTime.IsZero() {
		return e.EventTime.Time
	}
	return e.FirstTimestamp.Time
}

func toolGetMetrics(client *k8s.Client, args map[string]interface{}, zh bool) (toolResult, error) {
	if client == nil || client.Config == nil {
		if zh {
			return toolResult{Text: "没有集群配置，读不到 metrics-server。"}, nil
		}
		return toolResult{Text: "No cluster config, so metrics-server cannot be read."}, nil
	}
	scope := strings.ToLower(strArg(args, "scope"))
	ns := strArg(args, "namespace")
	name := strArg(args, "name")
	switch scope {
	case "node", "nodes":
		return nodeMetricsText(client, name, zh)
	case "pod", "pods":
		return podMetricsText(client, ns, name, zh)
	default:
		if zh {
			return toolResult{}, fmt.Errorf("scope 必须是 pod 或 node")
		}
		return toolResult{}, fmt.Errorf("scope must be pod or node")
	}
}

func nodeMetricsText(client *k8s.Client, name string, zh bool) (toolResult, error) {
	svc := service.NewNodeMetricsService()
	if name != "" {
		item, err := svc.GetNodeMetrics(client.Config, name)
		if err != nil {
			return toolResult{}, err
		}
		return toolResult{
			Text:      formatNodeMetric(*item),
			Resources: []ResourceRef{resourceRef("nodes", "", item.NodeName, false, "")},
		}, nil
	}
	all, err := svc.GetAllNodesMetrics(client.Config)
	if err != nil {
		if zh {
			return toolResult{Text: "metrics-server 不可用：" + err.Error()}, nil
		}
		return toolResult{Text: "metrics-server unavailable: " + err.Error()}, nil
	}
	if all == nil || len(all.Nodes) == 0 {
		if zh {
			return toolResult{Text: "没有节点指标。"}, nil
		}
		return toolResult{Text: "No node metrics."}, nil
	}
	var lines []string
	var refs []ResourceRef
	for i, n := range all.Nodes {
		if i >= maxMetricLines {
			lines = append(lines, scanMore(zh, "nodes"))
			break
		}
		lines = append(lines, formatNodeMetric(n))
		if len(refs) < 12 {
			refs = append(refs, resourceRef("nodes", "", n.NodeName, false, ""))
		}
	}
	return toolResult{Text: strings.Join(lines, "\n"), Resources: refs}, nil
}

func formatNodeMetric(n service.NodeMetrics) string {
	return fmt.Sprintf("node %s cpu=%s (%s) mem=%s (%s)", n.NodeName, n.CPUCores, n.CPUPercent, n.MemoryBytes, n.MemoryPercent)
}

func podMetricsText(client *k8s.Client, ns, name string, zh bool) (toolResult, error) {
	resp := service.NewPodMetricsService().ListPodMetrics(client.Config, ns)
	if resp == nil || !resp.Available {
		msg := ""
		if resp != nil {
			msg = resp.Message
		}
		if msg == "" {
			if zh {
				msg = "metrics-server 不可用"
			} else {
				msg = "metrics-server unavailable"
			}
		}
		return toolResult{Text: msg}, nil
	}
	var lines []string
	var refs []ResourceRef
	for _, p := range resp.Pods {
		if name != "" && p.Name != name {
			continue
		}
		if len(lines) >= maxMetricLines {
			lines = append(lines, scanMore(zh, "pods"))
			break
		}
		lines = append(lines, fmt.Sprintf(
			"pod %s/%s cpu=%s req=%s lim=%s mem=%s req=%s lim=%s",
			p.Namespace, p.Name, p.CPU, emptyDash(p.CPURequestPercent), emptyDash(p.CPULimitPercent),
			p.Memory, emptyDash(p.MemoryRequestPercent), emptyDash(p.MemoryLimitPercent),
		))
		if len(refs) < 12 {
			refs = append(refs, podRef(p.Namespace, p.Name, ""))
		}
	}
	if len(lines) == 0 {
		if zh {
			return toolResult{Text: "没有匹配的 Pod 指标。"}, nil
		}
		return toolResult{Text: "No matching pod metrics."}, nil
	}
	return toolResult{Text: strings.Join(lines, "\n"), Resources: refs}, nil
}

func emptyDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func toolQueryPrometheus(ctx context.Context, prom *service.PrometheusService, args map[string]interface{}, zh bool) (toolResult, error) {
	query, err := cleanPromQL(strArg(args, "query"))
	if err != nil {
		if zh && err.Error() == "query is required" {
			return toolResult{}, fmt.Errorf("需要 query")
		}
		if zh && strings.Contains(err.Error(), "800") {
			return toolResult{}, fmt.Errorf("query 超过 800 个字符")
		}
		if zh && strings.Contains(err.Error(), "single line") {
			return toolResult{}, fmt.Errorf("query 必须是单行")
		}
		return toolResult{}, err
	}
	if prom == nil || !prom.Enabled() {
		if zh {
			return toolResult{Text: "未配置 Prometheus，已跳过 PromQL。可改用 get_metrics 看 metrics-server。"}, nil
		}
		return toolResult{Text: "Prometheus is not configured; skipped PromQL. Use get_metrics for metrics-server instead."}, nil
	}
	res, err := prom.Query(ctx, query, nil)
	if err != nil {
		return toolResult{}, err
	}
	if res == nil {
		if zh {
			return toolResult{Text: "Prometheus 返回空结果。"}, nil
		}
		return toolResult{Text: "Prometheus returned an empty result."}, nil
	}
	if res.Status == "error" || res.Error != "" {
		return toolResult{}, fmt.Errorf("%s", res.Error)
	}
	text := string(res.Data)
	if !json.Valid(res.Data) {
		b, mErr := json.Marshal(res.Data)
		if mErr == nil {
			text = string(b)
		}
	}
	if len(text) > maxPromResult {
		text = text[:maxPromResult] + "\n... truncated ..."
	}
	return toolResult{Text: text}, nil
}

func cleanPromQL(q string) (string, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return "", fmt.Errorf("query is required")
	}
	if len(q) > maxPromQLLen {
		return "", fmt.Errorf("query exceeds 800 characters")
	}
	if strings.ContainsAny(q, "\r\n") {
		return "", fmt.Errorf("query must be a single line")
	}
	return q, nil
}
