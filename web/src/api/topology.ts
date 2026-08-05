import { apiGet } from '@/lib/api'

export type TopologyNode = {
  id: string
  kind: string
  name: string
  namespace: string
  group: string
  status: 'ok' | 'warn' | 'danger' | 'unknown' | string
  subtitle?: string
  href: string
  labels?: Record<string, string>
}

export type TopologyEdge = {
  id: string
  source: string
  target: string
  kind: string
  weight?: number
}

export type TopologyKindCount = { kind: string; count: number }

export type TopologyGraph = {
  namespace: string
  groupBy: string
  nodes: TopologyNode[]
  edges: TopologyEdge[]
  counts: TopologyKindCount[]
  truncated?: boolean
}

export type TopologyTrafficEdge = {
  source: string
  target: string
  rps: number
  mode: string
}

export type TopologyTraffic = {
  namespace: string
  mode: string
  edges: TopologyTrafficEdge[]
}

export function getTopologyGraph(params: {
  namespace: string
  groupBy?: 'app' | 'namespace'
  kinds?: string[]
}) {
  return apiGet<TopologyGraph>('/api/v1/topology', {
    namespace: params.namespace,
    groupBy: params.groupBy || 'app',
    kinds: params.kinds?.length ? params.kinds.join(',') : undefined,
  })
}

export function getTopologyTraffic(namespace: string) {
  return apiGet<TopologyTraffic>('/api/v1/topology/traffic', { namespace })
}
