import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getPodMetrics, type PodMetricsItem } from '@/api/cluster'
import { useCluster } from '@/store/cluster'

export function podMetricsKey(namespace: string, name: string) {
  return `${namespace}/${name}`
}

/** Fetch metrics-server pod snapshots; returns map keyed by namespace/name. */
export function usePodMetricsMap(namespace?: string) {
  const { clusterId } = useCluster()
  const q = useQuery({
    queryKey: ['pod-metrics', clusterId, namespace || '__all__'],
    queryFn: () => getPodMetrics(namespace || undefined),
    enabled: Boolean(clusterId),
    refetchInterval: 15_000,
    staleTime: 10_000,
  })

  const map = useMemo(() => {
    const m = new Map<string, PodMetricsItem>()
    for (const p of q.data?.pods || []) {
      m.set(podMetricsKey(p.namespace, p.name), p)
    }
    return m
  }, [q.data])

  return {
    map,
    available: Boolean(q.data?.available),
    message: q.data?.message || '',
    isLoading: q.isLoading,
    refetch: q.refetch,
  }
}
