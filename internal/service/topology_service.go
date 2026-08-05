package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const topologyProbeTimeout = 12 * time.Second

// TopologyNode is a graph vertex for the Topology page.
type TopologyNode struct {
	ID        string            `json:"id"`
	Kind      string            `json:"kind"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Group     string            `json:"group"`
	Status    string            `json:"status"` // ok | warn | danger | unknown
	Subtitle  string            `json:"subtitle,omitempty"`
	Href      string            `json:"href"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// TopologyEdge connects two nodes (ownership, selector, or ingress backend).
type TopologyEdge struct {
	ID     string  `json:"id"`
	Source string  `json:"source"`
	Target string  `json:"target"`
	Kind   string  `json:"kind"` // owner | selector | ingress | hpa
	Weight float64 `json:"weight,omitempty"`
}

// TopologyKindCount is a filter chip count.
type TopologyKindCount struct {
	Kind  string `json:"kind"`
	Count int    `json:"count"`
}

// TopologyGraphResponse is returned by GET /api/v1/topology.
type TopologyGraphResponse struct {
	Namespace string              `json:"namespace"`
	GroupBy   string              `json:"groupBy"`
	Nodes     []TopologyNode      `json:"nodes"`
	Edges     []TopologyEdge      `json:"edges"`
	Counts    []TopologyKindCount `json:"counts"`
	Truncated bool                `json:"truncated,omitempty"`
}

// TopologyTrafficEdge is a weighted edge for Traffic mode.
type TopologyTrafficEdge struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	RPS    float64 `json:"rps"`
	Mode   string  `json:"mode"` // prometheus | synthetic
}

// TopologyTrafficResponse is returned by GET /api/v1/topology/traffic.
type TopologyTrafficResponse struct {
	Namespace string                `json:"namespace"`
	Mode      string                `json:"mode"`
	Edges     []TopologyTrafficEdge `json:"edges"`
}

type TopologyService struct {
	prom *PrometheusService
}

func NewTopologyService(prom *PrometheusService) *TopologyService {
	return &TopologyService{prom: prom}
}

type TopologyBuildOpts struct {
	Namespace string
	GroupBy   string // app | namespace
	Kinds     map[string]bool
	MaxPods   int
}

// TopologyOptsForHandler builds options from query params.
func TopologyOptsForHandler(ns, groupBy string, kinds map[string]bool) TopologyBuildOpts {
	return TopologyBuildOpts{Namespace: ns, GroupBy: groupBy, Kinds: kinds}
}

func (s *TopologyService) BuildGraph(ctx context.Context, cs kubernetes.Interface, opts TopologyBuildOpts) (*TopologyGraphResponse, error) {
	if opts.GroupBy != "namespace" {
		opts.GroupBy = "app"
	}
	if opts.MaxPods <= 0 {
		opts.MaxPods = 120
	}
	ns := strings.TrimSpace(opts.Namespace)
	if ns == "" || ns == "all" {
		return nil, fmt.Errorf("namespace is required (pick a namespace; all-namespaces topology is not supported yet)")
	}

	ctx, cancel := context.WithTimeout(ctx, topologyProbeTimeout)
	defer cancel()

	want := opts.Kinds
	if len(want) == 0 {
		want = map[string]bool{
			"ingress": true, "service": true, "deployment": true, "statefulset": true,
			"daemonset": true, "pod": true, "job": true, "cronjob": true, "hpa": true,
			"configmap": true,
		}
	}

	out := &TopologyGraphResponse{
		Namespace: ns,
		GroupBy:   opts.GroupBy,
		Nodes:     []TopologyNode{},
		Edges:     []TopologyEdge{},
		Counts:    []TopologyKindCount{},
	}

	nodes := map[string]TopologyNode{}
	edges := map[string]TopologyEdge{}
	addNode := func(n TopologyNode) {
		if n.ID == "" {
			return
		}
		nodes[n.ID] = n
	}
	addEdge := func(e TopologyEdge) {
		if e.Source == "" || e.Target == "" || e.Source == e.Target {
			return
		}
		if e.ID == "" {
			e.ID = e.Source + "->" + e.Target + ":" + e.Kind
		}
		edges[e.ID] = e
	}
	groupOf := func(lbl map[string]string, namespace string) string {
		if opts.GroupBy == "namespace" {
			return namespace
		}
		for _, k := range []string{"app.kubernetes.io/name", "app", "app.kubernetes.io/instance", "k8s-app"} {
			if v := strings.TrimSpace(lbl[k]); v != "" {
				return v
			}
		}
		return "_ungrouped"
	}

	var pods []corev1.Pod
	if want["pod"] || want["deployment"] || want["statefulset"] || want["daemonset"] || want["job"] || want["service"] {
		list, err := cs.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{Limit: int64(opts.MaxPods)})
		if err != nil {
			return nil, err
		}
		pods = list.Items
		if list.RemainingItemCount != nil && *list.RemainingItemCount > 0 {
			out.Truncated = true
		}
	}

	podByName := map[string]*corev1.Pod{}
	for i := range pods {
		p := &pods[i]
		podByName[p.Name] = p
		if !want["pod"] {
			continue
		}
		st, sub := podStatus(p)
		addNode(TopologyNode{
			ID:        nodeID("pod", ns, p.Name),
			Kind:      "pod",
			Name:      p.Name,
			Namespace: ns,
			Group:     groupOf(p.Labels, ns),
			Status:    st,
			Subtitle:  sub,
			Href:      fmt.Sprintf("/pods/%s/%s", ns, p.Name),
			Labels:    p.Labels,
		})
		for _, own := range p.OwnerReferences {
			if own.Controller != nil && !*own.Controller {
				continue
			}
			ownerKind := strings.ToLower(own.Kind)
			// Pod → ReplicaSet → Deployment collapsed later; keep RS edge via synthetic
			if ownerKind == "replicaset" || ownerKind == "statefulset" || ownerKind == "daemonset" || ownerKind == "job" {
				src := nodeID(ownerKind, ns, own.Name)
				addEdge(TopologyEdge{Source: src, Target: nodeID("pod", ns, p.Name), Kind: "owner"})
			}
		}
	}

	if want["deployment"] {
		list, err := cs.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
		if err == nil {
			for i := range list.Items {
				d := &list.Items[i]
				st, sub := deployStatus(d)
				id := nodeID("deployment", ns, d.Name)
				podLabels := d.Spec.Template.Labels
				addNode(TopologyNode{
					ID: id, Kind: "deployment", Name: d.Name, Namespace: ns,
					Group: groupOf(d.Labels, ns), Status: st, Subtitle: sub,
					Href: fmt.Sprintf("/deployments/%s/%s", ns, d.Name), Labels: podLabels,
				})
				sel, _ := metav1.LabelSelectorAsSelector(d.Spec.Selector)
				for _, p := range pods {
					if sel.Matches(labels.Set(p.Labels)) {
						addEdge(TopologyEdge{Source: id, Target: nodeID("pod", ns, p.Name), Kind: "owner"})
					}
				}
			}
		}
	}

	if want["statefulset"] {
		list, err := cs.AppsV1().StatefulSets(ns).List(ctx, metav1.ListOptions{})
		if err == nil {
			for i := range list.Items {
				d := &list.Items[i]
				st, sub := stsStatus(d)
				id := nodeID("statefulset", ns, d.Name)
				podLabels := d.Spec.Template.Labels
				addNode(TopologyNode{
					ID: id, Kind: "statefulset", Name: d.Name, Namespace: ns,
					Group: groupOf(d.Labels, ns), Status: st, Subtitle: sub,
					Href: fmt.Sprintf("/statefulsets/%s/%s", ns, d.Name), Labels: podLabels,
				})
				sel, _ := metav1.LabelSelectorAsSelector(d.Spec.Selector)
				for _, p := range pods {
					if sel.Matches(labels.Set(p.Labels)) {
						addEdge(TopologyEdge{Source: id, Target: nodeID("pod", ns, p.Name), Kind: "owner"})
					}
				}
			}
		}
	}

	if want["daemonset"] {
		list, err := cs.AppsV1().DaemonSets(ns).List(ctx, metav1.ListOptions{})
		if err == nil {
			for i := range list.Items {
				d := &list.Items[i]
				st, sub := dsStatus(d)
				id := nodeID("daemonset", ns, d.Name)
				podLabels := d.Spec.Template.Labels
				addNode(TopologyNode{
					ID: id, Kind: "daemonset", Name: d.Name, Namespace: ns,
					Group: groupOf(d.Labels, ns), Status: st, Subtitle: sub,
					Href: fmt.Sprintf("/daemonsets/%s/%s", ns, d.Name), Labels: podLabels,
				})
				sel, _ := metav1.LabelSelectorAsSelector(d.Spec.Selector)
				for _, p := range pods {
					if sel.Matches(labels.Set(p.Labels)) {
						addEdge(TopologyEdge{Source: id, Target: nodeID("pod", ns, p.Name), Kind: "owner"})
					}
				}
			}
		}
	}

	if want["job"] {
		list, err := cs.BatchV1().Jobs(ns).List(ctx, metav1.ListOptions{})
		if err == nil {
			for i := range list.Items {
				j := &list.Items[i]
				st, sub := jobStatus(j)
				id := nodeID("job", ns, j.Name)
				addNode(TopologyNode{
					ID: id, Kind: "job", Name: j.Name, Namespace: ns,
					Group: groupOf(j.Labels, ns), Status: st, Subtitle: sub,
					Href: fmt.Sprintf("/jobs/%s/%s", ns, j.Name), Labels: j.Labels,
				})
			}
		}
	}

	if want["cronjob"] {
		list, err := cs.BatchV1().CronJobs(ns).List(ctx, metav1.ListOptions{})
		if err == nil {
			for i := range list.Items {
				c := &list.Items[i]
				id := nodeID("cronjob", ns, c.Name)
				addNode(TopologyNode{
					ID: id, Kind: "cronjob", Name: c.Name, Namespace: ns,
					Group: groupOf(c.Labels, ns), Status: "ok", Subtitle: c.Spec.Schedule,
					Href: fmt.Sprintf("/cronjobs/%s/%s", ns, c.Name), Labels: c.Labels,
				})
			}
		}
	}

	if want["service"] {
		list, err := cs.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
		if err == nil {
			for i := range list.Items {
				svc := &list.Items[i]
				id := nodeID("service", ns, svc.Name)
				sub := string(svc.Spec.Type)
				if svc.Spec.ClusterIP != "" && svc.Spec.ClusterIP != "None" {
					sub = svc.Spec.ClusterIP
				}
				addNode(TopologyNode{
					ID: id, Kind: "service", Name: svc.Name, Namespace: ns,
					Group: groupOf(svc.Labels, ns), Status: "ok", Subtitle: sub,
					Href: fmt.Sprintf("/services/%s/%s", ns, svc.Name), Labels: svc.Labels,
				})
				sel := labels.Set(svc.Spec.Selector).AsSelector()
				if len(svc.Spec.Selector) == 0 {
					continue
				}
				// Prefer linking to workloads that match; else pods
				linkedWorkload := false
				for _, n := range nodes {
					if n.Kind != "deployment" && n.Kind != "statefulset" && n.Kind != "daemonset" {
						continue
					}
					if sel.Matches(labels.Set(n.Labels)) {
						addEdge(TopologyEdge{Source: id, Target: n.ID, Kind: "selector"})
						linkedWorkload = true
					}
				}
				if !linkedWorkload {
					for _, p := range pods {
						if sel.Matches(labels.Set(p.Labels)) {
							addEdge(TopologyEdge{Source: id, Target: nodeID("pod", ns, p.Name), Kind: "selector"})
						}
					}
				}
			}
		}
	}

	if want["ingress"] {
		list, err := cs.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{})
		if err == nil {
			for i := range list.Items {
				ing := &list.Items[i]
				id := nodeID("ingress", ns, ing.Name)
				host := ""
				if len(ing.Spec.Rules) > 0 {
					host = ing.Spec.Rules[0].Host
				}
				addNode(TopologyNode{
					ID: id, Kind: "ingress", Name: ing.Name, Namespace: ns,
					Group: groupOf(ing.Labels, ns), Status: "ok", Subtitle: host,
					Href: fmt.Sprintf("/ingresses/%s/%s", ns, ing.Name), Labels: ing.Labels,
				})
				for _, rule := range ing.Spec.Rules {
					if rule.HTTP == nil {
						continue
					}
					for _, path := range rule.HTTP.Paths {
						if path.Backend.Service == nil {
							continue
						}
						svcName := path.Backend.Service.Name
						addEdge(TopologyEdge{
							Source: id, Target: nodeID("service", ns, svcName), Kind: "ingress",
						})
					}
				}
			}
		}
	}

	if want["hpa"] {
		list, err := cs.AutoscalingV2().HorizontalPodAutoscalers(ns).List(ctx, metav1.ListOptions{})
		if err == nil {
			for i := range list.Items {
				h := &list.Items[i]
				id := nodeID("hpa", ns, h.Name)
				sub := fmt.Sprintf("%d→%d", h.Status.CurrentReplicas, h.Status.DesiredReplicas)
				addNode(TopologyNode{
					ID: id, Kind: "hpa", Name: h.Name, Namespace: ns,
					Group: groupOf(h.Labels, ns), Status: "ok", Subtitle: sub,
					Href: fmt.Sprintf("/horizontalpodautoscalers/%s/%s", ns, h.Name), Labels: h.Labels,
				})
				ref := h.Spec.ScaleTargetRef
				kind := strings.ToLower(ref.Kind)
				addEdge(TopologyEdge{
					Source: id, Target: nodeID(kind, ns, ref.Name), Kind: "hpa",
				})
			}
		}
	}

	if want["configmap"] {
		list, err := cs.CoreV1().ConfigMaps(ns).List(ctx, metav1.ListOptions{Limit: 80})
		if err == nil {
			for i := range list.Items {
				cm := &list.Items[i]
				if strings.HasPrefix(cm.Name, "kube-root-ca") {
					continue
				}
				id := nodeID("configmap", ns, cm.Name)
				addNode(TopologyNode{
					ID: id, Kind: "configmap", Name: cm.Name, Namespace: ns,
					Group: groupOf(cm.Labels, ns), Status: "unknown", Subtitle: fmt.Sprintf("%d keys", len(cm.Data)),
					Href: fmt.Sprintf("/configmaps/%s/%s", ns, cm.Name), Labels: cm.Labels,
				})
			}
		}
	}

	// Drop edges whose endpoints were filtered out; drop orphan RS-only edges
	_ = podByName
	finalEdges := make([]TopologyEdge, 0, len(edges))
	for _, e := range edges {
		if _, ok := nodes[e.Source]; !ok {
			continue
		}
		if _, ok := nodes[e.Target]; !ok {
			continue
		}
		finalEdges = append(finalEdges, e)
	}

	// Counts
	countMap := map[string]int{}
	finalNodes := make([]TopologyNode, 0, len(nodes))
	for _, n := range nodes {
		// Skip dangling replicaset ids that were never materialized as nodes
		if strings.HasPrefix(n.Kind, "replicaset") {
			continue
		}
		finalNodes = append(finalNodes, n)
		countMap[n.Kind]++
	}
	sort.Slice(finalNodes, func(i, j int) bool {
		if finalNodes[i].Group != finalNodes[j].Group {
			return finalNodes[i].Group < finalNodes[j].Group
		}
		if finalNodes[i].Kind != finalNodes[j].Kind {
			return finalNodes[i].Kind < finalNodes[j].Kind
		}
		return finalNodes[i].Name < finalNodes[j].Name
	})
	kindsOrder := []string{"ingress", "service", "deployment", "statefulset", "daemonset", "pod", "job", "cronjob", "hpa", "configmap"}
	counts := make([]TopologyKindCount, 0, len(kindsOrder))
	for _, k := range kindsOrder {
		if c := countMap[k]; c > 0 {
			counts = append(counts, TopologyKindCount{Kind: k, Count: c})
		}
	}

	out.Nodes = finalNodes
	out.Edges = finalEdges
	out.Counts = counts
	return out, nil
}

func (s *TopologyService) BuildTraffic(ctx context.Context, cs kubernetes.Interface, namespace string) (*TopologyTrafficResponse, error) {
	graph, err := s.BuildGraph(ctx, cs, TopologyBuildOpts{Namespace: namespace, GroupBy: "app"})
	if err != nil {
		return nil, err
	}

	// Default: stable synthetic RPS (hash of edge endpoints) so Traffic mode always has a signal.
	// When Prometheus is enabled, try rate(http_requests_total[1m]) keyed by service name.
	mode := "synthetic"
	byService := map[string]float64{}

	if s.prom != nil {
		if st, serr := s.prom.GetStatus(ctx); serr == nil && st != nil {
			promMode, _ := st["mode"].(string)
			enabled, _ := st["enabled"].(bool)
			if enabled || promMode == "showcase" {
				queries := []string{
					`sum by (service) (rate(http_requests_total{namespace="` + namespace + `"}[1m]))`,
					`sum by (service) (rate(http_requests_total[1m]))`,
					`sum by (destination_service_name) (rate(istio_requests_total{reporter="destination",destination_service_namespace="` + namespace + `"}[1m]))`,
				}
				for _, q := range queries {
					res, qerr := s.prom.Query(ctx, q, nil)
					if qerr != nil || res == nil {
						continue
					}
					parsed := parsePromServiceRates(res.Data)
					if len(parsed) == 0 {
						continue
					}
					byService = parsed
					if promMode == "showcase" {
						mode = "showcase"
					} else {
						mode = "prometheus"
					}
					break
				}
				if mode == "synthetic" && (promMode == "showcase" || strings.Contains(fmt.Sprint(st["url"]), "showcase")) {
					mode = "showcase"
				}
			}
		}
	}

	out := &TopologyTrafficResponse{Namespace: namespace, Mode: mode, Edges: []TopologyTrafficEdge{}}
	for _, e := range graph.Edges {
		if e.Kind != "ingress" && e.Kind != "selector" {
			continue
		}
		rps := rateForEdge(e, byService)
		edgeMode := mode
		if rps <= 0 {
			rps = syntheticRPS(e.Source + e.Target)
			if mode == "showcase" {
				rps *= 1.8
			}
			if mode == "prometheus" {
				edgeMode = "synthetic"
			}
		}
		out.Edges = append(out.Edges, TopologyTrafficEdge{
			Source: e.Source, Target: e.Target, RPS: rps, Mode: edgeMode,
		})
	}
	// If Prom was intended but every edge fell back, report synthetic at top level.
	if mode == "prometheus" {
		allSynthetic := true
		for _, e := range out.Edges {
			if e.Mode == "prometheus" {
				allSynthetic = false
				break
			}
		}
		if allSynthetic {
			out.Mode = "synthetic"
		}
	}
	return out, nil
}

func rateForEdge(e TopologyEdge, byService map[string]float64) float64 {
	if len(byService) == 0 {
		return 0
	}
	for _, id := range []string{e.Source, e.Target} {
		parts := strings.Split(id, "/")
		if len(parts) < 3 {
			continue
		}
		kind, name := parts[0], parts[len(parts)-1]
		if kind != "service" && kind != "ingress" {
			continue
		}
		if v, ok := byService[name]; ok && v > 0 {
			return v
		}
		// destination_service_name sometimes includes ".ns.svc.cluster.local"
		for k, v := range byService {
			if v <= 0 {
				continue
			}
			if k == name || strings.HasPrefix(k, name+".") || strings.Contains(k, "/"+name) {
				return v
			}
		}
	}
	return 0
}

// parsePromServiceRates extracts service→rps from a Prometheus instant-query data blob.
func parsePromServiceRates(raw json.RawMessage) map[string]float64 {
	out := map[string]float64{}
	if len(raw) == 0 {
		return out
	}
	var data struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return out
	}
	labelKeys := []string{"service", "destination_service_name", "kubernetes_service", "svc", "job"}
	for _, row := range data.Result {
		name := ""
		for _, k := range labelKeys {
			if v := strings.TrimSpace(row.Metric[k]); v != "" {
				name = v
				break
			}
		}
		if name == "" {
			continue
		}
		if len(row.Value) < 2 {
			continue
		}
		vs, _ := row.Value[1].(string)
		f, err := strconv.ParseFloat(vs, 64)
		if err != nil || f < 0 {
			continue
		}
		out[name] += math.Round(f*10) / 10
	}
	return out
}

func syntheticRPS(key string) float64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	v := float64(h.Sum32()%900)/10.0 + 1.5
	return math.Round(v*10) / 10
}

func nodeID(kind, ns, name string) string {
	return strings.ToLower(kind) + "/" + ns + "/" + name
}

func podStatus(p *corev1.Pod) (string, string) {
	ready, total := 0, 0
	waiting := ""
	for _, cs := range p.Status.ContainerStatuses {
		total++
		if cs.Ready {
			ready++
		}
		if cs.State.Waiting != nil && cs.State.Waiting.Reason != "" {
			waiting = cs.State.Waiting.Reason
		}
	}
	sub := fmt.Sprintf("%d/%d %s", ready, total, p.Status.Phase)
	if waiting != "" {
		sub = fmt.Sprintf("%d/%d %s", ready, total, waiting)
	}
	switch p.Status.Phase {
	case corev1.PodFailed, corev1.PodUnknown:
		return "danger", sub
	case corev1.PodPending:
		return "warn", sub
	case corev1.PodSucceeded:
		return "ok", sub
	}
	if waiting != "" && (strings.Contains(waiting, "BackOff") || strings.Contains(waiting, "Err") || strings.Contains(waiting, "Image")) {
		return "danger", sub
	}
	if total > 0 && ready < total {
		return "warn", sub
	}
	return "ok", sub
}

func deployStatus(d *appsv1.Deployment) (string, string) {
	want := int32(1)
	if d.Spec.Replicas != nil {
		want = *d.Spec.Replicas
	}
	sub := fmt.Sprintf("%d/%d ready", d.Status.ReadyReplicas, want)
	if d.Status.ReadyReplicas < want {
		if d.Status.UnavailableReplicas > 0 {
			return "danger", sub
		}
		return "warn", sub
	}
	return "ok", sub
}

func stsStatus(d *appsv1.StatefulSet) (string, string) {
	want := int32(1)
	if d.Spec.Replicas != nil {
		want = *d.Spec.Replicas
	}
	sub := fmt.Sprintf("%d/%d ready", d.Status.ReadyReplicas, want)
	if d.Status.ReadyReplicas < want {
		return "warn", sub
	}
	return "ok", sub
}

func dsStatus(d *appsv1.DaemonSet) (string, string) {
	sub := fmt.Sprintf("%d/%d ready", d.Status.NumberReady, d.Status.DesiredNumberScheduled)
	if d.Status.NumberReady < d.Status.DesiredNumberScheduled {
		return "warn", sub
	}
	return "ok", sub
}

func jobStatus(j *batchv1.Job) (string, string) {
	sub := fmt.Sprintf("%d succeeded", j.Status.Succeeded)
	for _, c := range j.Status.Conditions {
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
			return "danger", "Failed"
		}
	}
	if j.Status.Active > 0 {
		return "warn", fmt.Sprintf("%d active", j.Status.Active)
	}
	return "ok", sub
}
