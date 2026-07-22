import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import dayjs from 'dayjs'
import { listNamespaces, listPods } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, EmptyState, PageHeader, StatCard } from '@/components/ui'

function podTone(phase?: string) {
  const p = (phase || '').toLowerCase()
  if (p === 'running') return 'ok' as const
  if (p === 'pending') return 'warn' as const
  if (p === 'failed' || p === 'unknown') return 'danger' as const
  return 'neutral' as const
}

export function WorkloadsPage() {
  const { clusterId } = useCluster()
  const [namespace, setNamespace] = useState('default')

  const nsQ = useQuery({
    queryKey: ['namespaces', clusterId],
    queryFn: listNamespaces,
    enabled: Boolean(clusterId),
  })

  useEffect(() => {
    if (!nsQ.data?.length) return
    if (!nsQ.data.includes(namespace)) {
      setNamespace(nsQ.data.includes('default') ? 'default' : nsQ.data[0])
    }
  }, [nsQ.data, namespace])

  const podsQ = useQuery({
    queryKey: ['pods', clusterId, namespace],
    queryFn: () => listPods(namespace),
    enabled: Boolean(clusterId && namespace),
  })

  const pods = podsQ.data || []

  return (
    <div>
      <PageHeader
        title="WORKLOADS"
        subtitle="Pods across the selected namespace"
        action={
          <select
            value={namespace}
            onChange={(e) => setNamespace(e.target.value)}
            className="hud-select w-auto min-w-[160px]"
          >
            {(nsQ.data || []).map((ns) => (
              <option key={ns} value={ns}>
                {ns}
              </option>
            ))}
          </select>
        }
      />

      <div className="mb-4 grid gap-3 sm:grid-cols-3">
        <StatCard label="Total" value={pods.length} />
        <StatCard
          label="Running"
          value={pods.filter((p: any) => p.status?.phase === 'Running').length}
        />
        <StatCard
          label="Pending / Failed"
          value={pods.filter((p: any) => ['Pending', 'Failed'].includes(p.status?.phase)).length}
        />
      </div>

      <Card className="overflow-hidden">
        <div className="overflow-x-auto">
          <table className="hud-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Namespace</th>
                <th>Status</th>
                <th>Restarts</th>
                <th>Created</th>
              </tr>
            </thead>
            <tbody>
              {pods.map((pod: any, index: number) => {
                const restarts = (pod.status?.containerStatuses || []).reduce(
                  (sum: number, c: any) => sum + (c.restartCount || 0),
                  0,
                )
                return (
                  <motion.tr
                    key={pod.metadata?.uid || pod.metadata?.name}
                    initial={{ opacity: 0, y: 8 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: Math.min(index, 12) * 0.02 }}
                  >
                    <td className="font-semibold text-cyan">{pod.metadata?.name}</td>
                    <td>{pod.metadata?.namespace}</td>
                    <td>
                      <Badge tone={podTone(pod.status?.phase)}>
                        {pod.status?.phase || 'Unknown'}
                      </Badge>
                    </td>
                    <td>{restarts}</td>
                    <td className="text-text-dim">
                      {pod.metadata?.creationTimestamp
                        ? dayjs(pod.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm')
                        : '-'}
                    </td>
                  </motion.tr>
                )
              })}
              {!podsQ.isLoading && !pods.length ? (
                <tr>
                  <td colSpan={5}>
                    <EmptyState>No pods in namespace {namespace}.</EmptyState>
                  </td>
                </tr>
              ) : null}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  )
}
