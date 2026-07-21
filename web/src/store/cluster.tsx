import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { useQuery } from '@tanstack/react-query'
import { listClusters, type ClusterItem } from '@/api/cluster'
import { getClusterId, setClusterId as persistClusterId } from '@/lib/api'
import { useAuth } from './auth'

type ClusterContextValue = {
  clusters: ClusterItem[]
  clusterId: string
  setClusterId: (id: string) => void
  loading: boolean
  activeCluster?: ClusterItem
}

const ClusterContext = createContext<ClusterContextValue | null>(null)

export function ClusterProvider({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth()
  const [clusterId, setClusterIdState] = useState(() => getClusterId())

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

  const setClusterId = useCallback((id: string) => {
    setClusterIdState(id)
    persistClusterId(id)
  }, [])

  const activeCluster = data.find((c) => c.id === clusterId || c.name === clusterId)

  const value = useMemo(
    () => ({
      clusters: data,
      clusterId,
      setClusterId,
      loading: isLoading,
      activeCluster,
    }),
    [data, clusterId, setClusterId, isLoading, activeCluster],
  )

  return <ClusterContext.Provider value={value}>{children}</ClusterContext.Provider>
}

export function useCluster() {
  const ctx = useContext(ClusterContext)
  if (!ctx) throw new Error('useCluster must be used within ClusterProvider')
  return ctx
}
