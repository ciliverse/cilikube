import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { listNamespaces, listPods } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { AgeCell, CreatedCell } from '@/components/AgeCell'
import { HudTable, HudTablePanel, ListPageFrame } from '@/components/HudTableScroll'
import { Badge, EmptyState, HudSelect, PageHeader, StatCard } from '@/components/ui'
import { metaCreated } from '@/api/resources'
import { podMetricsKey, usePodMetricsMap } from '@/hooks/usePodMetricsMap'
import { PercentCell } from '@/components/PodMetricCells'
import { shouldSkipEnterAnim } from '@/lib/motionPrefs'

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
  const metrics = usePodMetricsMap(namespace)

  const pods = podsQ.data || []

  return (
    <ListPageFrame>
      <PageHeader
        title="WORKLOADS"
        subtitle={
          metrics.available
            ? 'Pods · CPU/MEM + % request/limit'
            : 'Pods · metrics-server unavailable'
        }
        action={
          <HudSelect
            aria-label="Namespace"
            className="w-auto min-w-[160px]"
            value={namespace}
            onChange={setNamespace}
            searchableWhen={0}
            options={(nsQ.data || []).map((ns) => ({ value: ns, label: ns }))}
          />
        }
      />

      <div className="mb-3 grid shrink-0 gap-3 sm:grid-cols-3">
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

      <HudTablePanel wide pinFirst>
          <HudTable wide pinFirst>
            <thead>
              <tr>
                <th>Name</th>
                <th>Namespace</th>
                <th>Status</th>
                <th title="Current CPU usage">CPU</th>
                <th title="CPU % of Request">%CPU/R</th>
                <th title="CPU % of Limit">%CPU/L</th>
                <th title="Current memory usage">MEM</th>
                <th title="Memory % of Request">%MEM/R</th>
                <th title="Memory % of Limit">%MEM/L</th>
                <th>Restarts</th>
                <th>Age</th>
                <th>Created</th>
              </tr>
            </thead>
            <tbody>
              {pods.map((pod: any, index: number) => {
                const restarts = (pod.status?.containerStatuses || []).reduce(
                  (sum: number, c: any) => sum + (c.restartCount || 0),
                  0,
                )
                const name = pod.metadata?.name
                const ns = pod.metadata?.namespace
                const m = metrics.map.get(podMetricsKey(ns || '', name || ''))
                return (
                  <motion.tr
                    key={pod.metadata?.uid || name}
                    initial={shouldSkipEnterAnim() ? false : { opacity: 0, y: 8 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: shouldSkipEnterAnim() ? 0 : Math.min(index, 12) * 0.02 }}
                  >
                    <td className="font-semibold text-cyan">
                      {name && ns ? (
                        <Link className="hover:underline" to={`/pods/${ns}/${name}`}>
                          {name}
                        </Link>
                      ) : (
                        name || '-'
                      )}
                    </td>
                    <td>{ns}</td>
                    <td>
                      <Badge tone={podTone(pod.status?.phase)}>
                        {pod.status?.phase || 'Unknown'}
                      </Badge>
                    </td>
                    <td className="font-mono text-[12px] tabular-nums">{m?.cpu || '-'}</td>
                    <td>
                      <PercentCell
                        percent={m?.cpuRequestPercent}
                        ratio={m?.cpuRequestRatio}
                        hint={m?.cpuRequest}
                      />
                    </td>
                    <td>
                      <PercentCell
                        percent={m?.cpuLimitPercent}
                        ratio={m?.cpuLimitRatio}
                        hint={m?.cpuLimit}
                      />
                    </td>
                    <td className="font-mono text-[12px] tabular-nums">{m?.memory || '-'}</td>
                    <td>
                      <PercentCell
                        percent={m?.memoryRequestPercent}
                        ratio={m?.memoryRequestRatio}
                        hint={m?.memoryRequest}
                      />
                    </td>
                    <td>
                      <PercentCell
                        percent={m?.memoryLimitPercent}
                        ratio={m?.memoryLimitRatio}
                        hint={m?.memoryLimit}
                      />
                    </td>
                    <td>{restarts}</td>
                    <td>
                      <AgeCell value={metaCreated(pod)} />
                    </td>
                    <td>
                      <CreatedCell value={metaCreated(pod)} />
                    </td>
                  </motion.tr>
                )
              })}
              {!podsQ.isLoading && !pods.length ? (
                <tr>
                  <td colSpan={12}>
                    <EmptyState>No pods in namespace {namespace}.</EmptyState>
                  </td>
                </tr>
              ) : null}
            </tbody>
          </HudTable>
      </HudTablePanel>
    </ListPageFrame>
  )
}
