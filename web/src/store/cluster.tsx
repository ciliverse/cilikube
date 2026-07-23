import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { listClusters, type ClusterItem } from '@/api/cluster'
import { apiPost, getClusterId, setClusterId as persistClusterId } from '@/lib/api'
import { useAuth } from './auth'

type ClusterContextValue = {
  clusters: ClusterItem[]
  clusterId: string
  setClusterId: (id: string) => void
  loading: boolean
  switching: boolean
  activeCluster?: ClusterItem
}

const ClusterContext = createContext<ClusterContextValue | null>(null)

export function ClusterProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth()
  const queryClient = useQueryClient()
  const [clusterId, setClusterIdState] = useState(() => getClusterId())
  const [switching, setSwitching] = useState(false)

  const { data = [], isLoading } = useQuery({
    queryKey: ['clusters'],
    queryFn: listClusters,
    enabled: isAuthenticated,
    staleTime: 30_000,
  })

  useEffect(() => {
    if (!data.length) return
    if (!clusterId || !data.some((c) => c.id === clusterId || c.name === clusterId)) {
      const next = data.find((c) => c.is_active)?.id || data[0].id || data[0].name
      setClusterIdState(next)
      persistClusterId(next)
    }
  }, [data, clusterId])

  const setClusterId = useCallback(
    (id: string) => {
      if (!id || id === clusterId) return
      setClusterIdState(id)
      persistClusterId(id)
      setSwitching(true)

      void (async () => {
        try {
          // Keep server "active" cluster in sync (informers / fallbacks)
          await apiPost('/api/v1/clusters/active', { id })
        } catch {
          /* still switch client-side; resource APIs carry ?clusterId= */
        } finally {
          // Drop cached K8s views so lists refetch for the new cluster
          await queryClient.invalidateQueries({
            predicate: (q) => {
              const key = q.queryKey[0]
              return key !== 'clusters' && key !== 'settings-system' && key !== 'settings-oauth'
            },
          })
          void queryClient.invalidateQueries({ queryKey: ['clusters'] })
          setSwitching(false)
        }
      })()
    },
    [clusterId, queryClient],
  )

  const activeCluster = data.find((c) => c.id === clusterId || c.name === clusterId)

  const value = useMemo(
    () => ({
      clusters: data,
      clusterId,
      setClusterId,
      loading: isLoading,
      switching,
      activeCluster,
    }),
    [data, clusterId, setClusterId, isLoading, switching, activeCluster],
  )

  return <ClusterContext.Provider value={value}>{children}</ClusterContext.Provider>
}

export function useCluster() {
  const ctx = useContext(ClusterContext)
  if (!ctx) throw new Error('useCluster must be used within ClusterProvider')
  return ctx
}
