import { apiGet } from '@/lib/api'

export async function listNamespacedResource(namespace: string, resource: string) {
  const data = await apiGet<{ items: any[] }>(`/api/v1/namespaces/${namespace}/${resource}`)
  return data.items || []
}

export async function listClusterResource(resource: string) {
  const data = await apiGet<{ items: any[] }>(`/api/v1/${resource}`)
  return data.items || []
}

export function metaName(item: any) {
  return item?.metadata?.name || '-'
}

export function metaNamespace(item: any) {
  return item?.metadata?.namespace || '-'
}

export function metaCreated(item: any) {
  return item?.metadata?.creationTimestamp || ''
}
