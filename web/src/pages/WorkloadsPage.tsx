import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import dayjs from 'dayjs'
import { listNamespaces, listPods } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, PageHeader } from '@/components/ui'

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
        title="Workloads"
        subtitle="Pods across the selected namespace"
        action={
          <select
            value={namespace}
            onChange={(e) => setNamespace(e.target.value)}
            className="rounded-xl border border-line bg-white px-3 py-2 text-sm outline-none focus:border-accent"
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
        {[
          { label: 'Total', value: pods.length },
          {
            label: 'Running',
            value: pods.filter((p: any) => p.status?.phase === 'Running').length,
          },
          {
            label: 'Pending / Failed',
            value: pods.filter((p: any) =>
              ['Pending', 'Failed'].includes(p.status?.phase),
            ).length,
          },
        ].map((stat) => (
          <Card key={stat.label} className="p-4">
            <div className="text-xs font-semibold uppercase tracking-wider text-ink-soft">
              {stat.label}
            </div>
            <div className="mt-1 font-display text-2xl font-extrabold">{stat.value}</div>
          </Card>
        ))}
      </div>

      <Card className="overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full text-left text-sm">
            <thead className="bg-mist/80 text-xs uppercase tracking-wide text-ink-soft">
              <tr>
                <th className="px-5 py-3">Name</th>
                <th className="px-5 py-3">Namespace</th>
                <th className="px-5 py-3">Status</th>
                <th className="px-5 py-3">Restarts</th>
                <th className="px-5 py-3">Created</th>
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
                    className="border-t border-line/70"
                  >
                    <td className="px-5 py-3 font-semibold">{pod.metadata?.name}</td>
                    <td className="px-5 py-3">{pod.metadata?.namespace}</td>
                    <td className="px-5 py-3">
                      <Badge tone={podTone(pod.status?.phase)}>
                        {pod.status?.phase || 'Unknown'}
                      </Badge>
                    </td>
                    <td className="px-5 py-3">{restarts}</td>
                    <td className="px-5 py-3 text-ink-soft">
                      {pod.metadata?.creationTimestamp
                        ? dayjs(pod.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm')
                        : '-'}
                    </td>
                  </motion.tr>
                )
              })}
              {!podsQ.isLoading && !pods.length ? (
                <tr>
                  <td colSpan={5} className="px-5 py-10 text-ink-soft">
                    No pods in namespace {namespace}.
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
