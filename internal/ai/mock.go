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
	lang := normalizeLang(req.Language)
	zh := isZh(lang)
	var refs []ResourceRef
	var parts []string

	call := func(name string, args map[string]interface{}) {
		send(SSEEvent{Event: "tool_call", Data: map[string]interface{}{"name": name, "arguments": args}})
		res, err := executeTool(ctx, client, name, args, lang)
		if err != nil {
			send(SSEEvent{Event: "tool_result", Data: map[string]interface{}{
				"name": name, "ok": false, "content": err.Error(),
			}})
			if zh {
				parts = append(parts, fmt.Sprintf("工具 %s 失败：%s", name, err.Error()))
			} else {
				parts = append(parts, fmt.Sprintf("Tool %s failed: %s", name, err.Error()))
			}
			return
		}
		send(SSEEvent{Event: "tool_result", Data: map[string]interface{}{
			"name": name, "ok": true, "content": truncate(res.Text, 2000),
		}})
		mergeRefs(&refs, res.Resources)
		parts = append(parts, res.Text)
	}

	// nsMode: "" = UI/request namespace; "*" = all namespaces; other = force that ns.
	listKind := func(kind string, limit int, nsMode string, extra map[string]interface{}) {
		args := map[string]interface{}{"kind": kind, "limit": limit}
		for k, v := range extra {
			args[k] = v
		}
		switch {
		case nsMode == "*":
			// omit namespace → cluster-wide
		case nsMode != "":
			args["namespace"] = nsMode
		case ns != "" && kind != "nodes" && kind != "namespaces":
			args["namespace"] = ns
		}
		call("list_resources", args)
	}

	// Combo skills (explicit skill_id or prompt keywords, zh + en)
	comboInspect := skill == "skill_inspect_combo" ||
		strings.Contains(raw, "智能点检") || strings.Contains(q, "点检") ||
		strings.Contains(q, "smart inspect")
	comboTriage := skill == "skill_triage_combo" ||
		strings.Contains(raw, "快速调查") ||
		(strings.Contains(q, "调查") && strings.Contains(q, "日志")) ||
		strings.Contains(q, "quick triage")
	comboWorkload := skill == "skill_workload_snap" ||
		strings.Contains(raw, "工作负载快照") ||
		strings.Contains(q, "workload snapshot") ||
		(strings.Contains(q, "deploy") && strings.Contains(q, "service")) ||
		(strings.Contains(raw, "Deployments") && strings.Contains(raw, "Services"))

	wantOverview := comboInspect || strings.Contains(q, "概览") || strings.Contains(q, "overview") ||
		strings.Contains(q, "集群") || strings.Contains(q, "怎么样") || strings.Contains(q, "状态") ||
		strings.Contains(q, "脉诊") || strings.Contains(q, "健康") ||
		strings.Contains(q, "cluster pulse") || strings.Contains(q, "how is the cluster") ||
		q == "" || strings.Contains(q, "hello") || strings.Contains(q, "你好")
	wantPods := comboInspect || comboTriage || strings.Contains(q, "pod") || strings.Contains(q, "容器") ||
		strings.Contains(q, "failed") || strings.Contains(q, "失败") || strings.Contains(q, "pending") ||
		strings.Contains(q, "crash") || strings.Contains(q, "崩坏") || strings.Contains(q, "imagepull") ||
		strings.Contains(q, "重启") || strings.Contains(q, "restart") ||
		strings.Contains(q, "fault scan") || skill == "skill_fault_scan" || skill == "skill_crashloop" || skill == "skill_restarts"
	wantLogs := comboTriage || strings.Contains(q, "日志") || strings.Contains(q, "log") ||
		skill == "skill_log_sample"
	wantDeploy := comboWorkload || strings.Contains(q, "deploy") || strings.Contains(q, "部署") ||
		skill == "skill_deploy_list"
	wantSvc := comboWorkload || strings.Contains(q, "service") || strings.Contains(q, "服务") ||
		skill == "skill_svc_map"
	wantNode := strings.Contains(q, "node") || strings.Contains(q, "节点") ||
		skill == "skill_node_pressure"
	wantEvents := comboInspect || strings.Contains(q, "event") || strings.Contains(q, "事件") ||
		strings.Contains(q, "warning") || strings.Contains(q, "error") ||
		skill == "skill_events"
	wantNS := strings.Contains(q, "namespace") || strings.Contains(q, "命名空间") ||
		skill == "skill_namespaces"

	clusterWidePods := comboInspect || comboTriage ||
		strings.Contains(q, "failed") || strings.Contains(raw, "Failed") ||
		strings.Contains(q, "pending") || strings.Contains(raw, "Pending") ||
		strings.Contains(q, "异常") || strings.Contains(q, "故障") ||
		strings.Contains(q, "fault") || strings.Contains(q, "unhealthy")
	clusterWideEvents := comboInspect ||
		strings.Contains(q, "warning") || strings.Contains(raw, "Warning") ||
		strings.Contains(q, "error") || strings.Contains(raw, "Error")

	if wantOverview || (!wantPods && !wantDeploy && !wantSvc && !wantNode && !wantLogs && !wantEvents && !wantNS) {
		call("get_cluster_overview", map[string]interface{}{})
	}
	if wantPods || wantLogs {
		nsMode := ""
		extra := map[string]interface{}(nil)
		if clusterWidePods {
			nsMode = "*"
		}
		if comboInspect {
			extra = map[string]interface{}{"phases": "Pending,Failed,Unknown"}
		}
		listKind("pods", 20, nsMode, extra)
	}
	if wantDeploy {
		nsMode := ""
		if comboWorkload || strings.Contains(q, "default") {
			nsMode = "default"
		}
		listKind("deployments", 20, nsMode, nil)
	}
	if wantSvc {
		nsMode := ""
		if comboWorkload || strings.Contains(q, "default") {
			nsMode = "default"
		}
		listKind("services", 20, nsMode, nil)
	}
	if wantNode {
		listKind("nodes", 20, "*", nil)
	}
	if wantEvents {
		nsMode := ""
		extra := map[string]interface{}(nil)
		if clusterWideEvents {
			nsMode = "*"
		}
		if comboInspect {
			extra = map[string]interface{}{"event_types": "Warning,Error"}
		}
		listKind("events", 30, nsMode, extra)
	}
	if wantNS {
		listKind("namespaces", 40, "*", nil)
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

	msg := buildMockAnswer(lang, q, raw, skill, parts, refs, comboInspect, comboTriage, comboWorkload)
	send(SSEEvent{Event: "message", Data: map[string]string{"content": msg}})
	send(SSEEvent{Event: "done", Data: map[string]string{}})
}

func buildMockAnswer(
	lang, q, raw, skill string,
	parts []string,
	refs []ResourceRef,
	comboInspect, comboTriage, comboWorkload bool,
) string {
	zh := isZh(lang)
	var b strings.Builder
	if len(parts) == 0 {
		if zh {
			b.WriteString("没有查到数据。可以点 Skills：智能点检、故障扫描、日志取样，或直接问集群状态。")
		} else {
			b.WriteString("No data returned. Try Skills: Smart inspect, Fault scan, Log sample — or ask about cluster health.")
		}
		return b.String()
	}

	switch {
	case comboInspect || skill == "skill_inspect_combo":
		if zh {
			b.WriteString("【智能点检】集群调查员已跑完概览 / 异常 Pod / 事件取样：\n\n")
		} else {
			b.WriteString("**Smart inspect** — overview / unhealthy pods / event sample:\n\n")
		}
	case comboTriage || skill == "skill_triage_combo":
		if zh {
			b.WriteString("【快速调查】已定位可疑 Pod 并抽样日志：\n\n")
		} else {
			b.WriteString("**Quick triage** — suspect pods and sample logs:\n\n")
		}
	case comboWorkload || skill == "skill_workload_snap":
		if zh {
			b.WriteString("【工作负载快照】Deployments / Services：\n\n")
		} else {
			b.WriteString("**Workload snapshot** — Deployments / Services:\n\n")
		}
	default:
		if zh {
			b.WriteString("根据工具查询结果：\n\n")
		} else {
			b.WriteString("Based on tool results:\n\n")
		}
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
		if zh {
			b.WriteString("下方卡片可直接打开详情 / 日志 / 终端；需要改配置时请进控制台操作。")
		} else {
			b.WriteString("Use the cards below for details / logs / terminal. Make changes in the console.")
		}
	} else if strings.Contains(q, "日志") || strings.Contains(raw, "日志") || strings.Contains(q, "log") {
		if zh {
			b.WriteString("没有找到可打开的 Pod，试试先点「故障扫描」或问「有哪些 Pod」。")
		} else {
			b.WriteString("No openable Pod found — try Fault scan, or ask which pods exist.")
		}
	}
	return strings.TrimSpace(b.String())
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
