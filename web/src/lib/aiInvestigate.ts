/** Deep-link helpers: resource page → /ai investigate with context. */

import { isZhLang } from '@/lib/aiSkills'

/** True when the active cluster context has settled on the requested id/name. */
export function clusterContextReady(
  targetId: string | undefined,
  clusterId: string,
  active?: { id?: string; name?: string } | null,
): boolean {
  if (!targetId) return true
  if (clusterId === targetId) return true
  if (active?.id === targetId || active?.name === targetId) return true
  return false
}

export type AiInvestigateTarget = {
  kind: string
  name: string
  namespace?: string
  /** Route plural key, e.g. pods / deployments (optional, for console hints). */
  resource?: string
}

const KIND_FROM_RESOURCE: Record<string, string> = {
  pods: 'Pod',
  deployments: 'Deployment',
  statefulsets: 'StatefulSet',
  daemonsets: 'DaemonSet',
  services: 'Service',
  configmaps: 'ConfigMap',
  secrets: 'Secret',
  jobs: 'Job',
  cronjobs: 'CronJob',
  ingresses: 'Ingress',
  networkpolicies: 'NetworkPolicy',
  gatewayclasses: 'GatewayClass',
  gateways: 'Gateway',
  httproutes: 'HTTPRoute',
  serviceaccounts: 'ServiceAccount',
  persistentvolumeclaims: 'PersistentVolumeClaim',
  persistentvolumes: 'PersistentVolume',
  storageclasses: 'StorageClass',
  roles: 'Role',
  rolebindings: 'RoleBinding',
  clusterroles: 'ClusterRole',
  clusterrolebindings: 'ClusterRoleBinding',
  nodes: 'Node',
  events: 'Event',
  namespaces: 'Namespace',
}

export function kindFromResource(resource: string): string {
  return KIND_FROM_RESOURCE[resource] || resource
}

export function resourceRefLabel(t: AiInvestigateTarget): string {
  const ns = t.namespace?.trim()
  return ns ? `${t.kind} ${ns}/${t.name}` : `${t.kind} ${t.name}`
}

/** Build the user prompt sent when opening investigate from a resource page. */
export function buildInvestigatePrompt(t: AiInvestigateTarget, lang?: string): string {
  const label = resourceRefLabel(t)
  const kind = t.kind
  const lower = kind.toLowerCase()
  const zh = isZhLang(lang)

  if (lower === 'pod') {
    return zh
      ? `请调查 ${label}：查看该 Pod 当前状态、相关 Warning/Error 事件，必要时抽样最近日志，用中文给出简洁结论，并指出下一步该进控制台看哪里（详情 / 日志）。只做只读查证。`
      : `Investigate ${label}: check Pod status, related Warning/Error events, sample recent logs if needed. Give a concise English conclusion and where to go next in the console (details / logs). Read-only only.`
  }
  if (lower === 'node') {
    return zh
      ? `请调查 ${label}：查看节点状态与资源压力相关线索，列出该节点上异常 Pod（若有），用中文给出简洁结论与控制台下一步。只做只读查证。`
      : `Investigate ${label}: check node status and resource pressure, list unhealthy pods on this node if any. Give a concise English conclusion and console next steps. Read-only only.`
  }
  if (['deployment', 'statefulset', 'daemonset', 'job', 'cronjob'].includes(lower)) {
    return zh
      ? `请调查 ${label}：查看该工作负载状态与相关 Pod / 事件，定位是否有 Failed、Pending、CrashLoop 等问题，用中文给出简洁结论与控制台下一步。只做只读查证。`
      : `Investigate ${label}: check workload status and related pods/events for Failed, Pending, CrashLoop, etc. Give a concise English conclusion and console next steps. Read-only only.`
  }
  return zh
    ? `请调查 ${label}：拉取该资源详情与相关事件（若适用），说明当前是否异常、可能原因，并用中文给出简洁结论与建议进控制台查看的入口。只做只读查证。`
    : `Investigate ${label}: fetch resource details and related events if applicable. Say whether it looks unhealthy, likely causes, and a concise English conclusion with console entry points. Read-only only.`
}

export function buildInvestigateHref(t: AiInvestigateTarget, opts?: { auto?: boolean }): string {
  const sp = new URLSearchParams()
  sp.set('investigate', '1')
  sp.set('kind', t.kind)
  sp.set('name', t.name)
  if (t.namespace?.trim()) sp.set('namespace', t.namespace.trim())
  if (t.resource?.trim()) sp.set('resource', t.resource.trim())
  if (opts?.auto !== false) sp.set('auto', '1')
  return `/ai?${sp.toString()}`
}

export function parseInvestigateSearch(sp: URLSearchParams): AiInvestigateTarget | null {
  if (sp.get('investigate') !== '1') return null
  const name = (sp.get('name') || '').trim()
  if (!name) return null
  const resource = (sp.get('resource') || '').trim()
  const kind = (sp.get('kind') || '').trim() || (resource ? kindFromResource(resource) : '')
  if (!kind) return null
  const namespace = (sp.get('namespace') || '').trim()
  return {
    kind,
    name,
    namespace: namespace || undefined,
    resource: resource || undefined,
  }
}

/** Fleet card → /ai cluster-level inspect (smart pulse). */
export type AiFleetInspectTarget = {
  clusterId: string
  clusterName: string
}

export function buildFleetInspectPrompt(clusterName: string, lang?: string): string {
  const zh = isZhLang(lang)
  const label = clusterName.trim() || (zh ? '当前集群' : 'the current cluster')
  return zh
    ? `请对集群「${label}」做一次智能点检：先给集群概览，再列出 Failed 或 Pending 的 Pod，并汇总最近 Warning / Error 事件。用中文给出简洁结论，并指出下一步该进控制台看哪里。只做只读查证。`
    : `Run a smart inspect on cluster "${label}": give a cluster overview, list Failed or Pending pods, and summarize recent Warning / Error events. Concise English conclusion and console next steps. Read-only only.`
}

export function buildFleetInspectHref(
  t: AiFleetInspectTarget,
  opts?: { auto?: boolean },
): string {
  const sp = new URLSearchParams()
  sp.set('inspect', '1')
  sp.set('cluster', t.clusterId)
  if (t.clusterName.trim()) sp.set('name', t.clusterName.trim())
  if (opts?.auto !== false) sp.set('auto', '1')
  return `/ai?${sp.toString()}`
}

export function parseFleetInspectSearch(sp: URLSearchParams): AiFleetInspectTarget | null {
  if (sp.get('inspect') !== '1') return null
  const clusterId = (sp.get('cluster') || '').trim()
  if (!clusterId) return null
  const clusterName = (sp.get('name') || '').trim() || clusterId
  return { clusterId, clusterName }
}

/** Focused drill-down from fleet metric chips. */
export type AiFleetFocus = 'unhealthy' | 'warnings'

export type AiFleetFocusTarget = AiFleetInspectTarget & {
  focus: AiFleetFocus
}

export function buildFleetFocusPrompt(
  clusterName: string,
  focus: AiFleetFocus,
  lang?: string,
): string {
  const zh = isZhLang(lang)
  const label = clusterName.trim() || (zh ? '当前集群' : 'the current cluster')
  if (focus === 'unhealthy') {
    return zh
      ? `请聚焦集群「${label}」的异常 Pod：列出 Failed / Pending / CrashLoop 等不健康 Pod，挑 1～2 个代表性对象说明原因线索，并指出控制台下一步（详情 / 日志）。只做只读查证。`
      : `Focus on unhealthy pods in cluster "${label}": list Failed / Pending / CrashLoop pods, pick 1–2 representatives with cause clues, and point to console next steps (details / logs). Read-only only.`
  }
  return zh
    ? `请聚焦集群「${label}」的近期 Warning / Error 事件：按严重度或频次摘要，关联到可疑工作负载（若有），用中文给出简洁结论与控制台入口。只做只读查证。`
    : `Focus on recent Warning / Error events in cluster "${label}": summarize by severity or frequency, link to suspect workloads if any. Concise English conclusion and console entry points. Read-only only.`
}

export function buildFleetFocusHref(t: AiFleetFocusTarget, opts?: { auto?: boolean }): string {
  const sp = new URLSearchParams()
  sp.set('inspect', '1')
  sp.set('cluster', t.clusterId)
  if (t.clusterName.trim()) sp.set('name', t.clusterName.trim())
  sp.set('focus', t.focus)
  if (opts?.auto !== false) sp.set('auto', '1')
  return `/ai?${sp.toString()}`
}

export function parseFleetFocus(sp: URLSearchParams): AiFleetFocus | null {
  const f = (sp.get('focus') || '').trim()
  if (f === 'unhealthy' || f === 'warnings') return f
  return null
}

/** Fleet page → serial multi-cluster inspect + summary. */
export function buildFleetTourHref(opts?: { auto?: boolean }): string {
  const sp = new URLSearchParams()
  sp.set('tour', '1')
  if (opts?.auto !== false) sp.set('auto', '1')
  return `/ai?${sp.toString()}`
}

export function parseFleetTourSearch(sp: URLSearchParams): { auto: boolean } | null {
  if (sp.get('tour') !== '1') return null
  return { auto: sp.get('auto') !== '0' }
}

export function fleetIssueScore(c: {
  unhealthy_pods?: number
  warning_events?: number
  not_ready_nodes?: number
}): number {
  return (c.unhealthy_pods || 0) * 10 + (c.not_ready_nodes || 0) * 5 + (c.warning_events || 0)
}

export function buildFleetTourStepPrompt(
  clusterName: string,
  index: number,
  total: number,
  lang?: string,
): string {
  const zh = isZhLang(lang)
  const prefix = zh
    ? `[舰队巡检 ${index}/${total} · ${clusterName}] `
    : `[Fleet inspect ${index}/${total} · ${clusterName}] `
  const suffix = zh
    ? ` 本轮只查当前已切换的集群，结论标题请写明集群名。`
    : ` Only inspect the currently selected cluster this round; put the cluster name in the conclusion title.`
  return prefix + buildFleetInspectPrompt(clusterName, lang) + suffix
}

export function buildFleetTourSummaryPrompt(clusterNames: string[], lang?: string): string {
  const zh = isZhLang(lang)
  const list = clusterNames.map((n, i) => `${i + 1}. ${n}`).join('\n')
  return zh
    ? `以上已依次点检这些可达集群：\n${list}\n\n请做「舰队健康」汇总：按风险高低排序，标出优先处理的集群与原因，用中文给出简洁结论（不必再调用工具，基于上文即可）。`
    : `The reachable clusters above were inspected in order:\n${list}\n\nWrite a "fleet health" summary: rank by risk, call out which clusters to handle first and why. Concise English conclusion (no more tool calls — use the context above).`
}
