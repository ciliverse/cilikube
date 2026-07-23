import { useMemo } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Activity, Boxes, Container, Server, Siren } from 'lucide-react'
import { getNodeMetrics, getSummary, listEvents } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { useAuth } from '@/store/auth'
import { Badge, Card, EmptyState, PageHeader, StatCard } from '@/components/ui'
import { HudTable, HudTableScroll } from '@/components/HudTableScroll'
import { formatPercent } from '@/lib/utils'
import { parseNodeMetricsSummary } from '@/lib/nodeMetricsSummary'

const container = {
  hidden: { opacity: 0 },
  show: { opacity: 1, transition: { staggerChildren: 0.06 } },
}
const item = {
  hidden: { opacity: 0, y: 10 },
  show: { opacity: 1, y: 0 },
}

type ResourceTile = {
  key: string
  label: string
  field: string
  to: string
  resource?: string
}

const RESOURCE_TILES: ResourceTile[] = [
  { key: 'nodes', label: 'Nodes', field: 'nodes', to: '/nodes', resource: 'nodes' },
  { key: 'namespaces', label: 'Namespaces', field: 'namespaces', to: '/namespaces', resource: 'namespaces' },
  { key: 'pods', label: 'Pods', field: 'pods', to: '/pods', resource: 'pods' },
  { key: 'deployments', label: 'Deployments', field: 'deployments', to: '/deployments', resource: 'deployments' },
  { key: 'statefulSets', label: 'StatefulSets', field: 'statefulSets', to: '/statefulsets', resource: 'statefulsets' },
  { key: 'daemonSets', label: 'DaemonSets', field: 'daemonSets', to: '/daemonsets', resource: 'daemonsets' },
  { key: 'services', label: 'Services', field: 'services', to: '/services', resource: 'services' },
  { key: 'ingresses', label: 'Ingresses', field: 'ingresses', to: '/ingresses', resource: 'ingresses' },
  { key: 'pvs', label: 'PVs', field: 'persistentVolumes', to: '/persistentvolumes', resource: 'persistentvolumes' },
  { key: 'pvcs', label: 'PVCs', field: 'pvcs', to: '/persistentvolumeclaims', resource: 'persistentvolumeclaims' },
  { key: 'configMaps', label: 'ConfigMaps', field: 'configMaps', to: '/configmaps', resource: 'configmaps' },
  { key: 'secrets', label: 'Secrets', field: 'secrets', to: '/secrets', resource: 'secrets' },
]

function UsageBar({ label, percent, detail }: { label: string; percent: number; detail: string }) {
  const pct = Math.max(0, Math.min(100, percent))
  const tone = pct >= 80 ? 'bg-danger' : pct >= 50 ? 'bg-orange' : 'bg-cyan'
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between gap-2 text-xs">
        <span className="text-text-dim">{label}</span>
        <span className="font-mono text-text">
          {formatPercent(pct)} · {detail}
        </span>
      </div>
      <div className="h-1.5 overflow-hidden rounded-full bg-mist">
        <div className={`h-full rounded-full ${tone}`} style={{ width: `${pct}%` }} />
      </div>
    </div>
  )
}

export function OverviewPage() {
  const { clusterId, activeCluster } = useCluster()
  const { checkPermission } = useAuth()

  const summaryQ = useQuery({
    queryKey: ['summary', clusterId],
    queryFn: getSummary,
    enabled: Boolean(clusterId),
  })
  const eventsQ = useQuery({
    queryKey: ['events-recent', clusterId],
    queryFn: () => listEvents({ limit: 50 }),
    enabled: Boolean(clusterId),
  })
  const metricsQ = useQuery({
    queryKey: ['node-metrics', clusterId],
    queryFn: getNodeMetrics,
    enabled: Boolean(clusterId),
    refetchInterval: 20_000,
  })

  const metrics = metricsQ.data?.nodes || []
  const events = eventsQ.data?.events || []
  const warningEvents = events.filter((e: any) => e.type === 'Warning')
  const summary = summaryQ.data

  const rollup = useMemo(() => parseNodeMetricsSummary(metrics), [metrics])

  const cards = [
    {
      label: 'Nodes',
      value: summary?.nodes ?? rollup.totalNodes ?? '-',
      icon: <Server className="h-5 w-5" />,
    },
    {
      label: 'Namespaces',
      value: summary?.namespaces ?? '-',
      icon: <Boxes className="h-5 w-5" />,
    },
    {
      label: 'Pods',
      value: summary?.pods ?? '-',
      icon: <Container className="h-5 w-5" />,
    },
    {
      label: 'Avg CPU',
      value: rollup.totalNodes ? formatPercent(rollup.avgCpuUsagePercent) : '-',
      icon: <Activity className="h-5 w-5" />,
    },
    {
      label: 'Recent warnings',
      value: warningEvents.length,
      icon: <Siren className="h-5 w-5" />,
    },
  ]

  return (
    <div className="flex w-full min-w-0 flex-col gap-4">
      <PageHeader
        title="OVERVIEW"
        subtitle={`Live posture for ${activeCluster?.name || clusterId || 'selected cluster'}`}
        action={<Badge tone="accent">{activeCluster?.name || 'no cluster'}</Badge>}
      />

      <motion.div
        variants={container}
        initial="hidden"
        animate="show"
        className="grid w-full grid-cols-2 gap-3 sm:grid-cols-3 xl:grid-cols-5"
      >
        {cards.map((card) => (
          <motion.div key={card.label} variants={item}>
            <StatCard label={card.label} value={card.value} icon={card.icon} />
          </motion.div>
        ))}
      </motion.div>

      <Card className="min-w-0 p-5">
        <div className="mb-4 flex flex-wrap items-center justify-between gap-2">
          <div>
            <h2 className="font-display text-lg font-bold tracking-[0.12em]">CLUSTER RESOURCES</h2>
            <p className="mt-1 text-xs text-text-dim">
              Cluster-wide counts from /api/v1/summary/resources · click to open list
            </p>
          </div>
          <Badge tone="neutral">
            {summaryQ.isLoading ? 'loading…' : `${RESOURCE_TILES.length} types`}
          </Badge>
        </div>
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6">
          {RESOURCE_TILES.map((tile) => {
            const count = summary?.[tile.field]
            const allowed = !tile.resource || checkPermission(tile.resource, 'read')
            const value = count == null ? '—' : count
            const inner = (
              <>
                <div className="hud-label">{tile.label}</div>
                <div className="mt-1 font-display text-2xl font-bold tracking-wide text-cyan">
                  {value}
                </div>
              </>
            )
            if (!allowed) {
              return (
                <div
                  key={tile.key}
                  className="rounded border border-line/60 bg-mist/40 px-3 py-3 opacity-50"
                  title="No permission"
                >
                  {inner}
                </div>
              )
            }
            return (
              <Link
                key={tile.key}
                to={tile.to}
                className="rounded border border-line bg-mist px-3 py-3 transition hover:border-cyan/40 hover:bg-cyan/10"
              >
                {inner}
              </Link>
            )
          })}
        </div>
        {!summaryQ.isLoading && !summary ? (
          <div className="mt-3">
            <EmptyState>Summary unavailable for this cluster.</EmptyState>
          </div>
        ) : null}
      </Card>

      <Card className="min-w-0 p-5">
        <div className="mb-4 flex flex-wrap items-center justify-between gap-2">
          <div>
            <h2 className="font-display text-lg font-bold tracking-[0.12em]">NODE SNAPSHOT</h2>
            <p className="mt-1 text-xs text-text-dim">
              Aggregated from metrics-server · point-in-time (not Prometheus history)
            </p>
          </div>
          <Badge tone="neutral">
            {rollup.totalNodes} nodes
          </Badge>
        </div>
        {rollup.totalNodes ? (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <div className="rounded border border-line bg-mist px-3 py-3">
              <div className="hud-label">Avg CPU</div>
              <div className="mt-1 font-display text-xl font-bold text-cyan">
                {formatPercent(rollup.avgCpuUsagePercent)}
              </div>
            </div>
            <div className="rounded border border-line bg-mist px-3 py-3">
              <div className="hud-label">Avg Memory</div>
              <div className="mt-1 font-display text-xl font-bold text-orange">
                {formatPercent(rollup.avgMemoryUsagePercent)}
              </div>
            </div>
            <div className="rounded border border-line bg-mist px-3 py-3">
              <div className="hud-label">Capacity</div>
              <div className="mt-1 font-mono text-sm text-text">
                {rollup.totalCpuCores} cores · {rollup.totalMemoryGB} Gi
              </div>
            </div>
            <div className="space-y-3 rounded border border-line bg-mist px-3 py-3 md:col-span-2 lg:col-span-1">
              <UsageBar
                label="CPU requests"
                percent={rollup.totalCpuRequestsPercent}
                detail={`${rollup.totalCpuRequests} cores`}
              />
              <UsageBar
                label="Mem requests"
                percent={rollup.totalMemoryRequestsPercent}
                detail={`${rollup.totalMemoryRequests} Gi`}
              />
            </div>
          </div>
        ) : (
          <EmptyState>
            {metricsQ.isLoading ? 'Loading metrics…' : 'Connect metrics-server to populate node snapshot.'}
          </EmptyState>
        )}
      </Card>

      <div className="grid w-full gap-4 lg:grid-cols-[minmax(0,1.4fr)_minmax(0,1fr)]">
        <Card className="min-w-0 overflow-hidden">
          <div className="border-b border-line px-5 py-4">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <h2 className="font-display text-lg font-bold tracking-[0.12em]">
                NODE RESOURCE USAGE
              </h2>
              <Badge tone="neutral">metrics-server · snapshot</Badge>
            </div>
            <p className="mt-1 text-sm text-text-dim">
              Per-node CPU / memory / request saturation
            </p>
          </div>
          <HudTableScroll>
          <HudTable>
              <thead>
                <tr>
                  <th>Node</th>
                  <th>CPU</th>
                  <th>Memory</th>
                  <th>CPU req</th>
                  <th>Mem req</th>
                </tr>
              </thead>
              <tbody>
                {metrics.map((n: any) => (
                  <tr key={n.nodeName}>
                    <td>
                      <Link
                        className="font-semibold text-cyan hover:underline"
                        to={`/nodes/${encodeURIComponent(n.nodeName)}`}
                      >
                        {n.nodeName}
                      </Link>
                    </td>
                    <td>{n.cpuPercent || '-'}</td>
                    <td>{n.memoryPercent || '-'}</td>
                    <td>{n.cpuRequestsPercent || '-'}</td>
                    <td>{n.memoryRequestsPercent || '-'}</td>
                  </tr>
                ))}
                {!metricsQ.isLoading && !metrics.length ? (
                  <tr>
                    <td colSpan={5}>
                      <EmptyState>Connect metrics-server to populate this table.</EmptyState>
                    </td>
                  </tr>
                ) : null}
                {metricsQ.isLoading && !metrics.length ? (
                  <tr>
                    <td colSpan={5}>
                      <EmptyState>Loading metrics…</EmptyState>
                    </td>
                  </tr>
                ) : null}
              </tbody>
            </HudTable>
          </HudTableScroll>
        </Card>

        <Card className="min-w-0 p-5">
          <div className="mb-4 flex items-center justify-between gap-2">
            <h2 className="font-display text-lg font-bold tracking-[0.12em]">RECENT EVENTS</h2>
            <div className="flex items-center gap-2">
              <Badge tone="warn">{warningEvents.length} warnings</Badge>
              <Link to="/events" className="text-xs font-semibold text-cyan hover:underline">
                View all
              </Link>
            </div>
          </div>
          <div className="space-y-3">
            {events.slice(0, 8).map((event: any) => (
              <div
                key={event.id || `${event.name}-${event.lastTime}`}
                className="rounded border border-line bg-mist px-3 py-2.5"
              >
                <div className="flex items-center justify-between gap-2">
                  <div className="truncate text-sm font-semibold">{event.reason || event.name}</div>
                  <Badge tone={event.type === 'Warning' ? 'warn' : 'ok'}>
                    {event.type || 'Normal'}
                  </Badge>
                </div>
                <div className="mt-1 line-clamp-2 text-xs text-text-dim">
                  {event.message || '-'}
                </div>
              </div>
            ))}
            {!eventsQ.isLoading && !events.length ? (
              <div className="text-sm text-text-dim">No recent events</div>
            ) : null}
          </div>
        </Card>
      </div>
    </div>
  )
}
