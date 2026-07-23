import { apiGet, apiPut } from '@/lib/api'

export type ClusterItem = {
  id: string
  name: string
  status?: string
  version?: string
  is_active?: boolean
  description?: string
  [key: string]: unknown
}

export type UpdateClusterPayload = {
  name?: string
  description?: string
  provider?: string
  environment?: string
  region?: string
  status?: string
  kubeconfigData?: string
}

export async function listClusters() {
  const data = await apiGet<ClusterItem[] | { items: ClusterItem[]; clusters?: ClusterItem[] }>(
    '/api/v1/clusters',
  )
  if (Array.isArray(data)) return data
  return data.items || data.clusters || []
}

export async function updateCluster(id: string, body: UpdateClusterPayload) {
  return apiPut<null>(`/api/v1/clusters/${encodeURIComponent(id)}`, body)
}

export async function listNamespaces() {
  const data = await apiGet<{ items: Array<{ metadata: { name: string } }> }>(
    '/api/v1/namespaces',
  )
  return (data.items || []).map((item) => item.metadata.name)
}

export async function listNodes() {
  const data = await apiGet<{ items: any[] }>('/api/v1/nodes')
  return data.items || []
}

export async function listPods(namespace: string) {
  const data = await apiGet<{ items: any[] }>(`/api/v1/namespaces/${namespace}/pods`)
  return data.items || []
}

export async function listEvents(params?: { namespace?: string; limit?: number }) {
  return apiGet<{ events: any[]; total: number }>('/api/v1/events', params)
}

export async function getNodeMetrics() {
  return apiGet<{ nodes: any[]; total: number }>('/api/v1/nodes/metrics')
}

export type PodMetricsItem = {
  namespace: string
  name: string
  cpu: string
  memory: string
  cpuMilli?: number
  memoryBytes?: number
  cpuRequest?: string
  cpuLimit?: string
  memoryRequest?: string
  memoryLimit?: string
  cpuRequestPercent?: string
  cpuLimitPercent?: string
  memoryRequestPercent?: string
  memoryLimitPercent?: string
  cpuRequestRatio?: number
  cpuLimitRatio?: number
  memoryRequestRatio?: number
  memoryLimitRatio?: number
  timestamp?: string
}

export async function getPodMetrics(namespace?: string) {
  return apiGet<{
    pods: PodMetricsItem[]
    total: number
    available: boolean
    message?: string
  }>('/api/v1/pods/metrics', namespace ? { namespace } : undefined)
}

export async function getMonitoringDashboard() {
  return apiGet<any>('/api/v1/monitoring/dashboard')
}

export async function getPrometheusStatus() {
  return apiGet<{ enabled: boolean; url: string; healthy: boolean; error?: string }>(
    '/api/v1/prometheus/status',
  )
}

export async function getSummary() {
  try {
    return await apiGet<any>('/api/v1/summary/resources')
  } catch {
    return null
  }
}

export async function getObjectEvents(kind: string, name: string, namespace?: string) {
  return apiGet<{ events: any[]; total: number }>(`/api/v1/events/object/${kind}/${name}`, {
    namespace,
  })
}

export async function prometheusQuery(query: string) {
  return apiGet<any>('/api/v1/prometheus/query', { query })
}

export async function prometheusQueryRange(params: {
  query: string
  start: string | Date
  end: string | Date
  step?: string
}) {
  const start = typeof params.start === 'string' ? params.start : params.start.toISOString()
  const end = typeof params.end === 'string' ? params.end : params.end.toISOString()
  return apiGet<any>('/api/v1/prometheus/query_range', {
    query: params.query,
    start,
    end,
    step: params.step || '60s',
  })
}
