import { apiDelete, apiGet, apiPost, apiPut, getClusterId, getToken } from '@/lib/api'

export async function listNamespacedResource(namespace: string, resource: string) {
  // Empty namespace → cluster-wide list (GET /api/v1/<resource>)
  const path = namespace
    ? `/api/v1/namespaces/${encodeURIComponent(namespace)}/${resource}`
    : `/api/v1/${resource}`
  const data = await apiGet<{ items: any[] }>(path)
  return data.items || []
}

export async function listClusterResource(resource: string) {
  const data = await apiGet<{ items: any[] }>(`/api/v1/${resource}`)
  return data.items || []
}

export async function getNamespacedResource(namespace: string, resource: string, name: string) {
  return apiGet<any>(`/api/v1/namespaces/${namespace}/${resource}/${name}`)
}

export async function getClusterScopedResource(resource: string, name: string) {
  return apiGet<any>(`/api/v1/${resource}/${name}`)
}

export async function deleteNamespacedResource(namespace: string, resource: string, name: string) {
  return apiDelete<null>(`/api/v1/namespaces/${namespace}/${resource}/${name}`)
}

export async function deleteClusterScopedResource(resource: string, name: string) {
  return apiDelete<null>(`/api/v1/${resource}/${name}`)
}

export async function updateNamespacedResource(
  namespace: string,
  resource: string,
  name: string,
  body: unknown,
) {
  return apiPut<any>(`/api/v1/namespaces/${namespace}/${resource}/${name}`, body)
}

export async function updateClusterScopedResource(resource: string, name: string, body: unknown) {
  return apiPut<any>(`/api/v1/${resource}/${name}`, body)
}

export async function createNamespacedResource(
  namespace: string,
  resource: string,
  body: unknown,
) {
  return apiPost<any>(`/api/v1/namespaces/${namespace}/${resource}`, body)
}

export async function createClusterScopedResource(resource: string, body: unknown) {
  return apiPost<any>(`/api/v1/${resource}`, body)
}

export function metaName(item: any) {
  return item?.metadata?.name || '-'
}

export function metaNamespace(item: any) {
  return item?.metadata?.namespace || '-'
}

export function metaCreated(item: any) {
  const m = item?.metadata || {}
  return (
    m.creationTimestamp ||
    m.creation_timestamp ||
    m.CreationTimestamp ||
    item?.creationTimestamp ||
    ''
  )
}

export function podContainers(pod: any): Array<{ name: string; image: string }> {
  const specs = [
    ...(pod?.spec?.initContainers || []),
    ...(pod?.spec?.containers || []),
  ]
  return specs.map((c: any) => ({ name: c.name, image: c.image || '' }))
}

/** Build authenticated WebSocket URL for pod logs / exec / attach / portforward. */
export function buildPodWsUrl(
  kind: 'logs' | 'exec' | 'attach' | 'portforward',
  namespace: string,
  podName: string,
  params: Record<string, string | number | boolean | undefined> = {},
) {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const qs = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === '') continue
    // Backend accepts ports as local:remote (repeatable query).
    if (key === 'ports') {
      String(value)
        .split(',')
        .map((p) => p.trim())
        .filter(Boolean)
        .forEach((p) => qs.append('ports', p))
      continue
    }
    qs.set(key, String(value))
  }
  const token = getToken()
  const clusterId = getClusterId()
  if (token) qs.set('token', token)
  if (clusterId) qs.set('clusterId', clusterId)
  return `${protocol}//${window.location.host}/api/v1/namespaces/${encodeURIComponent(namespace)}/pods/${encodeURIComponent(podName)}/${kind}?${qs.toString()}`
}
