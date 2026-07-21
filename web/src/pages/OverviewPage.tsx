import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Activity, Boxes, Server, Siren } from 'lucide-react'
import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { getNodeMetrics, getSummary, listEvents, listNamespaces, listNodes } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, PageHeader } from '@/components/ui'
import { formatPercent } from '@/lib/utils'

const container = {
  hidden: { opacity: 0 },
  show: { opacity: 1, transition: { staggerChildren: 0.08 } },
}
const item = {
  hidden: { opacity: 0, y: 12 },
  show: { opacity: 1, y: 0 },
}

export function OverviewPage() {
  const { clusterId, activeCluster } = useCluster()

  const nodesQ = useQuery({
    queryKey: ['nodes', clusterId],
    queryFn: listNodes,
    enabled: Boolean(clusterId),
  })
  const nsQ = useQuery({
    queryKey: ['namespaces', clusterId],
    queryFn: listNamespaces,
    enabled: Boolean(clusterId),
  })
  const eventsQ = useQuery({
    queryKey: ['events-recent', clusterId],
    queryFn: () => listEvents({ limit: 8 }),
    enabled: Boolean(clusterId),
  })
  const metricsQ = useQuery({
    queryKey: ['node-metrics', clusterId],
    queryFn: getNodeMetrics,
    enabled: Boolean(clusterId),
    refetchInterval: 20_000,
  })
  useQuery({
    queryKey: ['summary', clusterId],
    queryFn: getSummary,
    enabled: Boolean(clusterId),
  })

  const nodes = nodesQ.data || []
  const metrics = metricsQ.data?.nodes || []
  const avgCpu =
    metrics.reduce((sum, n) => sum + parseFloat(String(n.cpuPercent || '0')), 0) /
    (metrics.length || 1)
  const avgMem =
    metrics.reduce((sum, n) => sum + parseFloat(String(n.memoryPercent || '0')), 0) /
    (metrics.length || 1)

  const chartData = metrics.slice(0, 8).map((n: any) => ({
    name: String(n.nodeName || '').replace(/^(.{8}).*/, '$1…'),
    cpu: parseFloat(String(n.cpuPercent || '0')) || 0,
    mem: parseFloat(String(n.memoryPercent || '0')) || 0,
  }))

  const cards = [
    {
      label: 'Nodes',
      value: nodes.length,
      icon: Server,
      tone: 'accent' as const,
    },
    {
      label: 'Namespaces',
      value: nsQ.data?.length ?? '-',
      icon: Boxes,
      tone: 'ok' as const,
    },
    {
      label: 'Avg CPU',
      value: metrics.length ? formatPercent(avgCpu) : '-',
      icon: Activity,
      tone: 'warn' as const,
    },
    {
      label: 'Warnings',
      value: (eventsQ.data?.events || []).filter((e: any) => e.type === 'Warning').length,
      icon: Siren,
      tone: 'danger' as const,
    },
  ]

  return (
    <div>
      <PageHeader
        title="Overview"
        subtitle={`Live posture for ${activeCluster?.name || clusterId || 'selected cluster'}`}
        action={<Badge tone="accent">{activeCluster?.name || 'no cluster'}</Badge>}
      />

      <motion.div
        variants={container}
        initial="hidden"
        animate="show"
        className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4"
      >
        {cards.map((card) => (
          <motion.div key={card.label} variants={item}>
            <Card className="p-5">
              <div className="flex items-start justify-between">
                <div>
                  <div className="text-xs font-semibold uppercase tracking-wider text-ink-soft">
                    {card.label}
                  </div>
                  <div className="mt-2 font-display text-3xl font-extrabold tracking-tight">
                    {card.value}
                  </div>
                </div>
                <div className="grid h-10 w-10 place-items-center rounded-xl bg-accent-soft text-accent">
                  <card.icon className="h-5 w-5" />
                </div>
              </div>
            </Card>
          </motion.div>
        ))}
      </motion.div>

      <div className="mt-6 grid gap-4 xl:grid-cols-[1.4fr_1fr]">
        <Card className="p-5">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="font-display text-xl font-bold">Node pressure</h2>
            <Badge tone="neutral">CPU / Memory</Badge>
          </div>
          <div className="h-64">
            {chartData.length ? (
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={chartData}>
                  <defs>
                    <linearGradient id="cpu" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#14b8a6" stopOpacity={0.45} />
                      <stop offset="95%" stopColor="#14b8a6" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <XAxis dataKey="name" tick={{ fontSize: 12 }} />
                  <YAxis tick={{ fontSize: 12 }} />
                  <Tooltip />
                  <Area
                    type="monotone"
                    dataKey="cpu"
                    stroke="#0f766e"
                    fill="url(#cpu)"
                    strokeWidth={2}
                  />
                  <Area
                    type="monotone"
                    dataKey="mem"
                    stroke="#ea580c"
                    fillOpacity={0.08}
                    strokeWidth={2}
                  />
                </AreaChart>
              </ResponsiveContainer>
            ) : (
              <div className="grid h-full place-items-center text-sm text-ink-soft">
                {metricsQ.isLoading ? 'Loading metrics…' : 'No metrics-server data yet'}
              </div>
            )}
          </div>
        </Card>

        <Card className="p-5">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="font-display text-xl font-bold">Recent events</h2>
            <Badge tone="neutral">{eventsQ.data?.total ?? 0}</Badge>
          </div>
          <div className="space-y-3">
            {(eventsQ.data?.events || []).slice(0, 6).map((event: any) => (
              <div
                key={event.id || `${event.name}-${event.lastTime}`}
                className="rounded-xl border border-line/80 bg-mist/50 px-3 py-2.5"
              >
                <div className="flex items-center justify-between gap-2">
                  <div className="truncate text-sm font-semibold">{event.reason || event.name}</div>
                  <Badge tone={event.type === 'Warning' ? 'warn' : 'ok'}>
                    {event.type || 'Normal'}
                  </Badge>
                </div>
                <div className="mt-1 line-clamp-2 text-xs text-ink-soft">
                  {event.message || '-'}
                </div>
              </div>
            ))}
            {!eventsQ.isLoading && !(eventsQ.data?.events || []).length ? (
              <div className="text-sm text-ink-soft">No recent events</div>
            ) : null}
          </div>
        </Card>
      </div>

      <Card className="mt-4 overflow-hidden">
        <div className="border-b border-line px-5 py-4">
          <h2 className="font-display text-xl font-bold">Cluster memory</h2>
          <p className="text-sm text-ink-soft">
            Avg memory {metrics.length ? formatPercent(avgMem) : '-'} across reported nodes
          </p>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full text-left text-sm">
            <thead className="bg-mist/80 text-xs uppercase tracking-wide text-ink-soft">
              <tr>
                <th className="px-5 py-3 font-semibold">Node</th>
                <th className="px-5 py-3 font-semibold">CPU</th>
                <th className="px-5 py-3 font-semibold">Memory</th>
                <th className="px-5 py-3 font-semibold">Requests</th>
              </tr>
            </thead>
            <tbody>
              {metrics.map((n: any) => (
                <tr key={n.nodeName} className="border-t border-line/70">
                  <td className="px-5 py-3 font-semibold">{n.nodeName}</td>
                  <td className="px-5 py-3">{n.cpuPercent || '-'}</td>
                  <td className="px-5 py-3">{n.memoryPercent || '-'}</td>
                  <td className="px-5 py-3 text-ink-soft">
                    CPU {n.cpuRequestsPercent || '-'} · Mem {n.memoryRequestsPercent || '-'}
                  </td>
                </tr>
              ))}
              {!metrics.length ? (
                <tr>
                  <td className="px-5 py-8 text-ink-soft" colSpan={4}>
                    Connect metrics-server to populate this table.
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
