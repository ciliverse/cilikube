import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { listClusterResource, listNamespacedResource } from '@/api/resources'
import { filterResourceHits, type SearchHit } from '@/lib/globalSearch'
import { useCluster } from '@/store/cluster'
import { useNamespace } from '@/store/namespace'

/** Client-side multi-kind search. Empty namespace = all namespaces. */
export function useGlobalSearchHits(query: string, enabled = true) {
  const { clusterId } = useCluster()
  const { namespace } = useNamespace()
  const ready = Boolean(clusterId) && enabled

  const podsQ = useQuery({
    queryKey: ['search-pods', clusterId, namespace],
    enabled: ready,
    queryFn: () => listNamespacedResource(namespace, 'pods'),
  })
  const depsQ = useQuery({
    queryKey: ['search-deps', clusterId, namespace],
    enabled: ready,
    queryFn: () => listNamespacedResource(namespace, 'deployments'),
  })
  const svcsQ = useQuery({
    queryKey: ['search-svcs', clusterId, namespace],
    enabled: ready,
    queryFn: () => listNamespacedResource(namespace, 'services'),
  })
  const stsQ = useQuery({
    queryKey: ['search-sts', clusterId, namespace],
    enabled: ready,
    queryFn: () => listNamespacedResource(namespace, 'statefulsets'),
  })
  const nodesQ = useQuery({
    queryKey: ['search-nodes', clusterId],
    enabled: ready,
    queryFn: () => listClusterResource('nodes'),
  })

  const loading =
    podsQ.isLoading || depsQ.isLoading || svcsQ.isLoading || stsQ.isLoading || nodesQ.isLoading

  const hits: SearchHit[] = useMemo(
    () =>
      filterResourceHits(query, {
        pods: podsQ.data,
        deployments: depsQ.data,
        statefulsets: stsQ.data,
        services: svcsQ.data,
        nodes: nodesQ.data,
      }),
    [query, podsQ.data, depsQ.data, stsQ.data, svcsQ.data, nodesQ.data],
  )

  return { hits, loading, namespace }
}
