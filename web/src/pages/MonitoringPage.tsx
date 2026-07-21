import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import {
  getMonitoringDashboard,
  getNodeMetrics,
  getPrometheusStatus,
} from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, PageHeader } from '@/components/ui'
import { formatPercent } from '@/lib/utils'

export function MonitoringPage() {
  const { clusterId } = useCluster()

  const dashQ = useQuery({
    queryKey: ['monitoring-dashboard'],
    queryFn: getMonitoringDashboard,
    refetchInterval: 20_000,
  })
  const promQ = useQuery({
    queryKey: ['prometheus-status'],
    queryFn: getPrometheusStatus,
    refetchInterval: 30_000,
  })
  const metricsQ = useQuery({
    queryKey: ['node-metrics', clusterId],
    queryFn: getNodeMetrics,
    enabled: Boolean(clusterId),
    refetchInterval: 20_000,
  })

  const summary = dashQ.data?.summary
  const prom = promQ.data
  const nodes = metricsQ.data?.nodes || []

  return (
    <div>
      <PageHeader
        title="Monitoring"
        subtitle="Security posture, Prometheus health, and node telemetry"
        action={
          <Badge
            tone={
              prom?.healthy ? 'ok' : prom?.enabled ? 'danger' : 'neutral'
            }
          >
            Prometheus {prom?.enabled ? (prom.healthy ? 'healthy' : 'down') : 'off'}
          </Badge>
        }
      />

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {[
          { label: 'Active sessions', value: summary?.active_sessions ?? '-' },
          { label: 'Active users', value: summary?.active_users ?? '-' },
          { label: 'Active threats', value: summary?.active_threats ?? '-' },
          { label: 'Failed login rate', value: summary?.failed_logins_rate ?? '-' },
        ].map((card, index) => (
          <motion.div
            key={card.label}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.05 }}
          >
            <Card className="p-5">
              <div className="text-xs font-semibold uppercase tracking-wider text-ink-soft">
                {card.label}
              </div>
              <div className="mt-2 font-display text-3xl font-extrabold">{card.value}</div>
            </Card>
          </motion.div>
        ))}
      </div>

      <div className="mt-4 grid gap-4 xl:grid-cols-2">
        <Card className="p-5">
          <h2 className="font-display text-xl font-bold">Prometheus</h2>
          <p className="mt-1 text-sm text-ink-soft">
            Endpoint: {prom?.url || 'not configured'}
          </p>
          {prom?.error ? (
            <div className="mt-3 rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-warn">
              {prom.error}
            </div>
          ) : (
            <div className="mt-3 text-sm text-ink-soft">
              Configure `prometheus.enabled` and `prometheus.url` in backend config to enable
              PromQL queries.
            </div>
          )}
          <div className="mt-4 rounded-xl bg-mist px-3 py-2 font-mono text-xs text-ink-soft">
            GET /api/v1/prometheus/query?query=up
          </div>
        </Card>

        <Card className="p-5">
          <h2 className="font-display text-xl font-bold">Security health</h2>
          <div className="mt-3 space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-ink-soft">Status</span>
              <Badge tone="accent">{dashQ.data?.health?.status || summary?.status || '-'}</Badge>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-ink-soft">Violations</span>
              <span className="font-semibold">{summary?.security_violations ?? '-'}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-ink-soft">Alerts</span>
              <span className="font-semibold">
                {dashQ.data?.alerts
                  ? dashQ.data.alerts.critical +
                    dashQ.data.alerts.warning +
                    dashQ.data.alerts.info
                  : '-'}
              </span>
            </div>
          </div>
        </Card>
      </div>

      <Card className="mt-4 overflow-hidden">
        <div className="border-b border-line px-5 py-4">
          <h2 className="font-display text-xl font-bold">Node telemetry</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full text-left text-sm">
            <thead className="bg-mist/80 text-xs uppercase tracking-wide text-ink-soft">
              <tr>
                <th className="px-5 py-3">Node</th>
                <th className="px-5 py-3">CPU</th>
                <th className="px-5 py-3">Memory</th>
                <th className="px-5 py-3">CPU req</th>
                <th className="px-5 py-3">Mem req</th>
              </tr>
            </thead>
            <tbody>
              {nodes.map((n: any) => (
                <tr key={n.nodeName} className="border-t border-line/70">
                  <td className="px-5 py-3 font-semibold">{n.nodeName}</td>
                  <td className="px-5 py-3">{formatPercent(n.cpuPercent)}</td>
                  <td className="px-5 py-3">{formatPercent(n.memoryPercent)}</td>
                  <td className="px-5 py-3">{n.cpuRequestsPercent || '-'}</td>
                  <td className="px-5 py-3">{n.memoryRequestsPercent || '-'}</td>
                </tr>
              ))}
              {!nodes.length ? (
                <tr>
                  <td colSpan={5} className="px-5 py-8 text-ink-soft">
                    No node metrics available.
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
