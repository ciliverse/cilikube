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
import { listNamespaces } from '@/api/cluster'
import { useCluster } from './cluster'

type NamespaceContextValue = {
  namespaces: string[]
  namespace: string
  setNamespace: (ns: string) => void
  loading: boolean
}

const NamespaceContext = createContext<NamespaceContextValue | null>(null)
const STORAGE_KEY = 'cilikube_namespace'

export function NamespaceProvider({ children }: { children: ReactNode }) {
  const { clusterId } = useCluster()
  const [namespace, setNamespaceState] = useState(
    () => localStorage.getItem(STORAGE_KEY) || 'default',
  )

  const { data = [], isLoading } = useQuery({
    queryKey: ['namespaces', clusterId],
    queryFn: listNamespaces,
    enabled: Boolean(clusterId),
  })

  useEffect(() => {
    if (!data.length) return
    if (!data.includes(namespace)) {
      const next = data.includes('default') ? 'default' : data[0]
      setNamespaceState(next)
      localStorage.setItem(STORAGE_KEY, next)
    }
  }, [data, namespace])

  const setNamespace = useCallback((ns: string) => {
    setNamespaceState(ns)
    localStorage.setItem(STORAGE_KEY, ns)
  }, [])

  const value = useMemo(
    () => ({
      namespaces: data,
      namespace,
      setNamespace,
      loading: isLoading,
    }),
    [data, namespace, setNamespace, isLoading],
  )

  return <NamespaceContext.Provider value={value}>{children}</NamespaceContext.Provider>
}

export function useNamespace() {
  const ctx = useContext(NamespaceContext)
  if (!ctx) throw new Error('useNamespace must be used within NamespaceProvider')
  return ctx
}
