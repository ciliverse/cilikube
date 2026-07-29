package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/ciliverse/cilikube/pkg/k8s"
)

// processMock runs a lightweight heuristic agent so local UI works without an API key.
func processMock(ctx context.Context, client *k8s.Client, req *ChatRequest, send sendFn) {
	q := strings.ToLower(lastUserText(req.Messages))
	ns := req.Namespace
	var refs []ResourceRef
	var parts []string

	call := func(name string, args map[string]interface{}) {
		send(SSEEvent{Event: "tool_call", Data: map[string]interface{}{"name": name, "arguments": args}})
		res, err := executeTool(ctx, client, name, args)
		if err != nil {
			send(SSEEvent{Event: "tool_result", Data: map[string]interface{}{
				"name": name, "ok": false, "content": err.Error(),
			}})
			parts = append(parts, fmt.Sprintf("工具 %s 失败：%s", name, err.Error()))
			return
		}
		send(SSEEvent{Event: "tool_result", Data: map[string]interface{}{
			"name": name, "ok": true, "content": truncate(res.Text, 2000),
		}})
		mergeRefs(&refs, res.Resources)
		parts = append(parts, res.Text)
	}

	wantOverview := strings.Contains(q, "概览") || strings.Contains(q, "overview") ||
		strings.Contains(q, "集群") || strings.Contains(q, "怎么样") || strings.Contains(q, "状态") ||
		q == "" || strings.Contains(q, "hello") || strings.Contains(q, "你好")
	wantPods := strings.Contains(q, "pod") || strings.Contains(q, "容器") ||
		strings.Contains(q, "failed") || strings.Contains(q, "失败") || strings.Contains(q, "pending")
	wantLogs := strings.Contains(q, "日志") || strings.Contains(q, "log")
	wantDeploy := strings.Contains(q, "deploy") || strings.Contains(q, "部署")
	wantSvc := strings.Contains(q, "service") || strings.Contains(q, "服务")
	wantNode := strings.Contains(q, "node") || strings.Contains(q, "节点")

	if wantOverview || (!wantPods && !wantDeploy && !wantSvc && !wantNode && !wantLogs) {
		call("get_cluster_overview", map[string]interface{}{})
	}
	if wantPods || wantLogs {
		args := map[string]interface{}{"kind": "pods", "limit": 20}
		if ns != "" {
			args["namespace"] = ns
		}
		call("list_resources", args)
	}
	if wantDeploy {
		args := map[string]interface{}{"kind": "deployments", "limit": 20}
		if ns != "" {
			args["namespace"] = ns
		}
		call("list_resources", args)
	}
	if wantSvc {
		args := map[string]interface{}{"kind": "services", "limit": 20}
		if ns != "" {
			args["namespace"] = ns
		}
		call("list_resources", args)
	}
	if wantNode {
		call("list_resources", map[string]interface{}{"kind": "nodes", "limit": 20})
	}

	// If user asked for logs and we have a pod ref, fetch first pod logs
	if wantLogs {
		for _, r := range refs {
			if r.Kind == "pods" && r.Namespace != "" && r.Name != "" {
				call("get_pod_logs", map[string]interface{}{
					"namespace": r.Namespace, "name": r.Name, "tail": 60,
				})
				break
			}
		}
	}

	if len(refs) > 0 {
		send(SSEEvent{Event: "resources", Data: map[string]interface{}{"items": refs}})
	}

	msg := buildMockAnswer(q, parts, refs)
	send(SSEEvent{Event: "message", Data: map[string]string{"content": msg}})
	send(SSEEvent{Event: "done", Data: map[string]string{}})
}

func buildMockAnswer(q string, parts []string, refs []ResourceRef) string {
	var b strings.Builder
	if len(parts) == 0 {
		b.WriteString("没有查到数据。可以问：集群概览、有哪些 Pod、某个 Failed Pod 的日志。")
		return b.String()
	}
	b.WriteString("根据工具查询结果：\n\n")
	for i, p := range parts {
		if i >= 3 {
			break
		}
		snippet := truncate(p, 800)
		b.WriteString(snippet)
		b.WriteString("\n\n")
	}
	if len(refs) > 0 {
		b.WriteString("下方卡片可直接打开详情 / 日志 / 终端。")
	} else if strings.Contains(q, "日志") {
		b.WriteString("没有找到可打开的 Pod，试试先问「有哪些 Pod」。")
	}
	return strings.TrimSpace(b.String())
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
