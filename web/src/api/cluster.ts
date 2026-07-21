import { apiGet } from '@/lib/api'

export type ClusterItem = {
  id: string
  name: string
  status?: string
  version?: string
  is_active?: boolean
  [key: string]: unknown
}

export async function listClusters() {
  const data = await apiGet<ClusterItem[] | { items: ClusterItem[]; clusters?: ClusterItem[] }>(
    '/api/v1/clusters',
  )
  if (Array.isArray(data)) return data
  return data.items || data.clusters || []
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
    return await apiGet<any>('/api/v1/summary')
  } catch {
    return null
  }
}
