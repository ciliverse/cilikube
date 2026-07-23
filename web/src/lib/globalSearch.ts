import { metaName, metaNamespace } from '@/api/resources'

export type SearchHit = {
  kind: string
  name: string
  namespace?: string
  to: string
}

export function filterResourceHits(
  needle: string,
  resources: {
    pods?: any[]
    deployments?: any[]
    statefulsets?: any[]
    services?: any[]
    nodes?: any[]
  },
  limit = 100,
): SearchHit[] {
  const q = needle.trim().toLowerCase()
  if (!q) return []

  const out: SearchHit[] = []

  const pushNamespaced = (items: any[] | undefined, kind: string, path: string) => {
    for (const item of items || []) {
      const name = metaName(item)
      if (!name.toLowerCase().includes(q)) continue
      const ns = metaNamespace(item)
      out.push({ kind, name, namespace: ns, to: `/${path}/${ns}/${name}` })
      if (out.length >= limit) return
    }
  }

  pushNamespaced(resources.pods, 'Pod', 'pods')
  if (out.length >= limit) return out
  pushNamespaced(resources.deployments, 'Deployment', 'deployments')
  if (out.length >= limit) return out
  pushNamespaced(resources.statefulsets, 'StatefulSet', 'statefulsets')
  if (out.length >= limit) return out
  pushNamespaced(resources.services, 'Service', 'services')
  if (out.length >= limit) return out

  for (const n of resources.nodes || []) {
    const name = metaName(n)
    if (!name.toLowerCase().includes(q)) continue
    out.push({ kind: 'Node', name, to: `/nodes/${name}` })
    if (out.length >= limit) break
  }

  return out
}
