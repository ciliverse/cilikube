package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/ciliverse/cilikube/pkg/k8s"
)

// processMock runs a lightweight heuristic agent so local UI works without an API key.
func processMock(ctx context.Context, client *k8s.Client, req *ChatRequest, send sendFn) {
	raw := lastUserText(req.Messages)
	q := strings.ToLower(raw)
	ns := req.Namespace
	skill := strings.ToLower(strings.TrimSpace(req.SkillID))
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

	listKind := func(kind string, limit int, forceNS string) {
		args := map[string]interface{}{"kind": kind, "limit": limit}
		useNS := forceNS
		if useNS == "" {
			useNS = ns
		}
		if useNS != "" && kind != "nodes" && kind != "namespaces" {
			args["namespace"] = useNS
		}
		call("list_resources", args)
	}

	// Combo skills (explicit skill_id or prompt keywords)
	comboInspect := skill == "skill_inspect_combo" ||
		strings.Contains(raw, "智能点检") || strings.Contains(q, "点检")
	comboTriage := skill == "skill_triage_combo" ||
		strings.Contains(raw, "快速调查") || (strings.Contains(q, "调查") && strings.Contains(q, "日志"))
	comboWorkload := skill == "skill_workload_snap" ||
		strings.Contains(raw, "工作负载快照") ||
		(strings.Contains(q, "deploy") && strings.Contains(q, "service")) ||
		(strings.Contains(raw, "Deployments") && strings.Contains(raw, "Services"))

	wantOverview := comboInspect || strings.Contains(q, "概览") || strings.Contains(q, "overview") ||
		strings.Contains(q, "集群") || strings.Contains(q, "怎么样") || strings.Contains(q, "状态") ||
		strings.Contains(q, "脉诊") || strings.Contains(q, "健康") ||
		q == "" || strings.Contains(q, "hello") || strings.Contains(q, "你好")
	wantPods := comboInspect || comboTriage || strings.Contains(q, "pod") || strings.Contains(q, "容器") ||
		strings.Contains(q, "failed") || strings.Contains(q, "失败") || strings.Contains(q, "pending") ||
		strings.Contains(q, "crash") || strings.Contains(q, "崩坏") || strings.Contains(q, "imagepull") ||
		strings.Contains(q, "重启")
	wantLogs := comboTriage || strings.Contains(q, "日志") || strings.Contains(q, "log")
	wantDeploy := comboWorkload || strings.Contains(q, "deploy") || strings.Contains(q, "部署")
	wantSvc := comboWorkload || strings.Contains(q, "service") || strings.Contains(q, "服务")
	wantNode := strings.Contains(q, "node") || strings.Contains(q, "节点")
	wantEvents := comboInspect || strings.Contains(q, "event") || strings.Contains(q, "事件") ||
		strings.Contains(q, "warning") || strings.Contains(q, "error")
	wantNS := strings.Contains(q, "namespace") || strings.Contains(q, "命名空间")

	if wantOverview || (!wantPods && !wantDeploy && !wantSvc && !wantNode && !wantLogs && !wantEvents && !wantNS) {
		call("get_cluster_overview", map[string]interface{}{})
	}
	if wantPods || wantLogs {
		listKind("pods", 20, "")
	}
	if wantDeploy {
		force := ""
		if comboWorkload || strings.Contains(q, "default") {
			force = "default"
		}
		listKind("deployments", 20, force)
	}
	if wantSvc {
		force := ""
		if comboWorkload || strings.Contains(q, "default") {
			force = "default"
		}
		listKind("services", 20, force)
	}
	if wantNode {
		listKind("nodes", 20, "")
	}
	if wantEvents {
		listKind("events", 30, "")
	}
	if wantNS {
		listKind("namespaces", 40, "")
	}

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

	msg := buildMockAnswer(q, raw, skill, parts, refs, comboInspect, comboTriage, comboWorkload)
	send(SSEEvent{Event: "message", Data: map[string]string{"content": msg}})
	send(SSEEvent{Event: "done", Data: map[string]string{}})
}

func buildMockAnswer(
	q, raw, skill string,
	parts []string,
	refs []ResourceRef,
	comboInspect, comboTriage, comboWorkload bool,
) string {
	var b strings.Builder
	if len(parts) == 0 {
		b.WriteString("没有查到数据。可以点 Skills：智能点检、故障扫描、日志取样，或直接问集群状态。")
		return b.String()
	}

	switch {
	case comboInspect || skill == "skill_inspect_combo":
		b.WriteString("【智能点检】集群调查员已跑完概览 / 异常 Pod / 事件取样：\n\n")
	case comboTriage || skill == "skill_triage_combo":
		b.WriteString("【快速调查】已定位可疑 Pod 并抽样日志：\n\n")
	case comboWorkload || skill == "skill_workload_snap":
		b.WriteString("【工作负载快照】Deployments / Services：\n\n")
	default:
		b.WriteString("根据工具查询结果：\n\n")
	}

	maxParts := 4
	if comboInspect {
		maxParts = 5
	}
	for i, p := range parts {
		if i >= maxParts {
			break
		}
		b.WriteString(truncate(p, 800))
		b.WriteString("\n\n")
	}
	if len(refs) > 0 {
		b.WriteString("下方卡片可直接打开详情 / 日志 / 终端；需要改配置时请进控制台操作。")
	} else if strings.Contains(q, "日志") || strings.Contains(raw, "日志") {
		b.WriteString("没有找到可打开的 Pod，试试先点「故障扫描」或问「有哪些 Pod」。")
	}
	return strings.TrimSpace(b.String())
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
