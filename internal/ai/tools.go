package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/ciliverse/cilikube/pkg/k8s"
)

type toolDef struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}

func toolDefinitions() []toolDef {
	return []toolDef{
		{
			Name:        "get_cluster_overview",
			Description: "Summarize nodes, namespaces, and pod phase counts for the active cluster.",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "list_resources",
			Description: "List Kubernetes resources by kind. Supported: pods, deployments, services, nodes, namespaces, events.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"kind":      map[string]interface{}{"type": "string"},
					"namespace": map[string]interface{}{"type": "string"},
					"limit":     map[string]interface{}{"type": "integer"},
				},
				"required": []string{"kind"},
			},
		},
		{
			Name:        "get_resource",
			Description: "Get a single resource as YAML. Secrets/ConfigMaps are redacted.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"kind":      map[string]interface{}{"type": "string"},
					"namespace": map[string]interface{}{"type": "string"},
					"name":      map[string]interface{}{"type": "string"},
				},
				"required": []string{"kind", "name"},
			},
		},
		{
			Name:        "get_pod_logs",
			Description: "Fetch recent logs from a pod (tail lines).",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{"type": "string"},
					"name":      map[string]interface{}{"type": "string"},
					"container": map[string]interface{}{"type": "string"},
					"tail":      map[string]interface{}{"type": "integer"},
				},
				"required": []string{"namespace", "name"},
			},
		},
	}
}

type toolResult struct {
	Text      string
	Resources []ResourceRef
}

func executeTool(ctx context.Context, client *k8s.Client, name string, args map[string]interface{}, lang string) (toolResult, error) {
	zh := isZh(lang)
	switch name {
	case "get_cluster_overview":
		return toolClusterOverview(ctx, client, zh)
	case "list_resources":
		return toolListResources(ctx, client, args, zh)
	case "get_resource":
		return toolGetResource(ctx, client, args, zh)
	case "get_pod_logs":
		return toolPodLogs(ctx, client, args, zh)
	default:
		if zh {
			return toolResult{}, fmt.Errorf("未知工具 %q", name)
		}
		return toolResult{}, fmt.Errorf("unknown tool %q", name)
	}
}

func toolClusterOverview(ctx context.Context, client *k8s.Client, zh bool) (toolResult, error) {
	nodes, err := client.Clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return toolResult{}, err
	}
	nss, err := client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return toolResult{}, err
	}
	pods, err := client.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return toolResult{}, err
	}
	phase := map[string]int{}
	var bad []ResourceRef
	for _, p := range pods.Items {
		ph := string(p.Status.Phase)
		phase[ph]++
		if ph == "Failed" || ph == "Pending" || ph == "Unknown" {
			if len(bad) < 12 {
				bad = append(bad, podRef(p.Namespace, p.Name, ""))
			}
		}
	}
	readyNodes := 0
	for _, n := range nodes.Items {
		for _, c := range n.Status.Conditions {
			if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
				readyNodes++
				break
			}
		}
	}
	var text string
	if zh {
		text = fmt.Sprintf(
			"节点=%d 就绪=%d 命名空间=%d Pod=%d 相位=%v",
			len(nodes.Items), readyNodes, len(nss.Items), len(pods.Items), phase,
		)
	} else {
		text = fmt.Sprintf(
			"nodes=%d ready=%d namespaces=%d pods=%d phases=%v",
			len(nodes.Items), readyNodes, len(nss.Items), len(pods.Items), phase,
		)
	}
	return toolResult{Text: text, Resources: bad}, nil
}

func toolListResources(ctx context.Context, client *k8s.Client, args map[string]interface{}, zh bool) (toolResult, error) {
	kind := strings.ToLower(strArg(args, "kind"))
	ns := strArg(args, "namespace")
	limit := intArg(args, "limit", 30)
	if limit <= 0 || limit > 100 {
		limit = 30
	}

	var lines []string
	var refs []ResourceRef

	switch kind {
	case "pod", "pods":
		list, err := client.Clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return toolResult{}, err
		}
		phaseAllow := parseCSVSet(strArg(args, "phases"))
		n := 0
		for _, p := range list.Items {
			ph := string(p.Status.Phase)
			if len(phaseAllow) > 0 && !phaseAllow[ph] {
				continue
			}
			if n >= limit {
				break
			}
			lines = append(lines, fmt.Sprintf("%s/%s phase=%s", p.Namespace, p.Name, p.Status.Phase))
			refs = append(refs, podRef(p.Namespace, p.Name, ""))
			n++
		}
	case "deployment", "deployments":
		list, err := client.Clientset.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return toolResult{}, err
		}
		for i, d := range list.Items {
			if i >= limit {
				break
			}
			lines = append(lines, fmt.Sprintf("%s/%s ready=%d/%d", d.Namespace, d.Name, d.Status.ReadyReplicas, d.Status.Replicas))
			refs = append(refs, resourceRef("deployments", d.Namespace, d.Name, true, ""))
		}
	case "service", "services":
		list, err := client.Clientset.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return toolResult{}, err
		}
		for i, s := range list.Items {
			if i >= limit {
				break
			}
			lines = append(lines, fmt.Sprintf("%s/%s type=%s", s.Namespace, s.Name, s.Spec.Type))
			refs = append(refs, resourceRef("services", s.Namespace, s.Name, true, ""))
		}
	case "node", "nodes":
		list, err := client.Clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			return toolResult{}, err
		}
		for i, n := range list.Items {
			if i >= limit {
				break
			}
			lines = append(lines, n.Name)
			refs = append(refs, resourceRef("nodes", "", n.Name, false, ""))
		}
	case "namespace", "namespaces":
		list, err := client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return toolResult{}, err
		}
		for i, n := range list.Items {
			if i >= limit {
				break
			}
			lines = append(lines, n.Name)
			refs = append(refs, resourceRef("namespaces", "", n.Name, false, ""))
		}
	case "event", "events":
		list, err := client.Clientset.CoreV1().Events(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return toolResult{}, err
		}
		typeAllow := parseCSVSet(strArg(args, "event_types"))
		n := 0
		for _, e := range list.Items {
			if len(typeAllow) > 0 && !typeAllow[e.Type] {
				continue
			}
			if n >= limit {
				break
			}
			lines = append(lines, fmt.Sprintf("[%s] %s/%s %s: %s", e.Type, e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message))
			n++
		}
	default:
		if zh {
			return toolResult{}, fmt.Errorf("不支持的 kind %q", kind)
		}
		return toolResult{}, fmt.Errorf("unsupported kind %q", kind)
	}

	if len(lines) == 0 {
		scope := "cluster-wide"
		scopeZh := "全集群"
		if ns != "" {
			scope = "namespace=" + ns
			scopeZh = "命名空间=" + ns
		}
		filter := ""
		filterZh := ""
		if p := strArg(args, "phases"); p != "" {
			filter = "; phases=" + p
			filterZh = "；相位=" + p
		}
		if t := strArg(args, "event_types"); t != "" {
			filter += "; event_types=" + t
			filterZh += "；事件类型=" + t
		}
		if zh {
			return toolResult{Text: fmt.Sprintf("无结果（%s%s）", scopeZh, filterZh)}, nil
		}
		return toolResult{Text: fmt.Sprintf("no items (%s%s)", scope, filter)}, nil
	}
	return toolResult{Text: strings.Join(lines, "\n"), Resources: refs}, nil
}

func parseCSVSet(raw string) map[string]bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	out := map[string]bool{}
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out[p] = true
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func toolGetResource(ctx context.Context, client *k8s.Client, args map[string]interface{}, zh bool) (toolResult, error) {
	kind := strings.ToLower(strArg(args, "kind"))
	ns := strArg(args, "namespace")
	name := strArg(args, "name")
	if name == "" {
		if zh {
			return toolResult{}, fmt.Errorf("需要 name")
		}
		return toolResult{}, fmt.Errorf("name required")
	}

	var obj interface{}
	var refs []ResourceRef
	var err error

	switch kind {
	case "pod", "pods":
		obj, err = client.Clientset.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
		refs = []ResourceRef{podRef(ns, name, "logs"), podRef(ns, name, "exec")}
	case "deployment", "deployments":
		obj, err = client.Clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
		refs = []ResourceRef{resourceRef("deployments", ns, name, true, "")}
	case "service", "services":
		obj, err = client.Clientset.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
		refs = []ResourceRef{resourceRef("services", ns, name, true, "")}
	case "node", "nodes":
		obj, err = client.Clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
		refs = []ResourceRef{resourceRef("nodes", "", name, false, "")}
	case "configmap", "configmaps":
		cm, e := client.Clientset.CoreV1().ConfigMaps(ns).Get(ctx, name, metav1.GetOptions{})
		err = e
		if cm != nil {
			redacted := cm.DeepCopy()
			for k := range redacted.Data {
				redacted.Data[k] = "***"
			}
			for k := range redacted.BinaryData {
				redacted.BinaryData[k] = []byte("***")
			}
			obj = redacted
			refs = []ResourceRef{resourceRef("configmaps", ns, name, true, "")}
		}
	case "secret", "secrets":
		sec, e := client.Clientset.CoreV1().Secrets(ns).Get(ctx, name, metav1.GetOptions{})
		err = e
		if sec != nil {
			redacted := sec.DeepCopy()
			for k := range redacted.Data {
				redacted.Data[k] = []byte("***")
			}
			obj = redacted
			refs = []ResourceRef{resourceRef("secrets", ns, name, true, "")}
		}
	case "namespace", "namespaces":
		obj, err = client.Clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
		refs = []ResourceRef{resourceRef("namespaces", "", name, false, "")}
	default:
		if zh {
			return toolResult{}, fmt.Errorf("不支持的 kind %q", kind)
		}
		return toolResult{}, fmt.Errorf("unsupported kind %q", kind)
	}
	if err != nil {
		return toolResult{}, err
	}
	b, err := yaml.Marshal(obj)
	if err != nil {
		return toolResult{}, err
	}
	text := string(b)
	if len(text) > 12000 {
		if zh {
			text = text[:12000] + "\n…（已截断）"
		} else {
			text = text[:12000] + "\n... truncated ..."
		}
	}
	return toolResult{Text: text, Resources: refs}, nil
}

func toolPodLogs(ctx context.Context, client *k8s.Client, args map[string]interface{}, zh bool) (toolResult, error) {
	ns := strArg(args, "namespace")
	name := strArg(args, "name")
	container := strArg(args, "container")
	tail := int64(intArg(args, "tail", 80))
	if tail <= 0 || tail > 500 {
		tail = 80
	}
	opts := &corev1.PodLogOptions{TailLines: &tail}
	if container != "" {
		opts.Container = container
	}
	stream, err := client.Clientset.CoreV1().Pods(ns).GetLogs(name, opts).Stream(ctx)
	if err != nil {
		return toolResult{}, err
	}
	defer stream.Close()
	b, err := io.ReadAll(io.LimitReader(stream, 64*1024))
	if err != nil {
		return toolResult{}, err
	}
	text := string(b)
	if text == "" {
		if zh {
			text = "（日志为空）"
		} else {
			text = "(empty logs)"
		}
	}
	return toolResult{
		Text: text,
		Resources: []ResourceRef{
			podRef(ns, name, "logs"),
			podRef(ns, name, "exec"),
		},
	}, nil
}

func podRef(ns, name, console string) ResourceRef {
	href := fmt.Sprintf("/pods/%s/%s", ns, name)
	label := fmt.Sprintf("Pod %s/%s", ns, name)
	if console != "" {
		href += "?console=" + console
		label += " · " + console
	}
	return ResourceRef{
		Kind: "pods", Namespace: ns, Name: name, Href: href, Console: console, Label: label,
	}
}

func resourceRef(resource, ns, name string, namespaced bool, console string) ResourceRef {
	var href string
	if namespaced {
		href = fmt.Sprintf("/%s/%s/%s", resource, ns, name)
	} else {
		href = fmt.Sprintf("/%s/%s", resource, name)
	}
	if console != "" {
		href += "?console=" + console
	}
	label := name
	if ns != "" {
		label = ns + "/" + name
	}
	return ResourceRef{Kind: resource, Namespace: ns, Name: name, Href: href, Console: console, Label: label}
}

func strArg(args map[string]interface{}, key string) string {
	if args == nil {
		return ""
	}
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

func intArg(args map[string]interface{}, key string, def int) int {
	if args == nil {
		return def
	}
	v, ok := args[key]
	if !ok || v == nil {
		return def
	}
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case json.Number:
		i, _ := t.Int64()
		return int(i)
	default:
		var n int
		_, _ = fmt.Sscanf(fmt.Sprint(t), "%d", &n)
		if n == 0 {
			return def
		}
		return n
	}
}
