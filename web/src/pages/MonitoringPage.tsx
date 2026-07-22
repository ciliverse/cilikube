import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import {
  getMonitoringDashboard,
  getNodeMetrics,
  getPrometheusStatus,
} from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, EmptyState, PageHeader, StatCard } from '@/components/ui'
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
        title="MONITORING"
        subtitle="Security posture, Prometheus health, and node telemetry"
        action={
          <Badge tone={prom?.healthy ? 'ok' : prom?.enabled ? 'danger' : 'neutral'}>
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
            <StatCard label={card.label} value={card.value} />
          </motion.div>
        ))}
      </div>

      <div className="mt-4 grid gap-4 xl:grid-cols-2">
        <Card className="p-5">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">PROMETHEUS</h2>
          <p className="mt-1 text-sm text-text-dim">
            Endpoint: {prom?.url || 'not configured'}
          </p>
          {prom?.error ? (
            <div className="mt-3 rounded border border-warn/40 bg-warn/10 px-3 py-2 text-sm text-warn">
              {prom.error}
            </div>
          ) : (
            <div className="mt-3 text-sm text-text-dim">
              Configure `prometheus.enabled` and `prometheus.url` in backend config to enable
              PromQL queries.
            </div>
          )}
          <div className="mt-4 rounded border border-line bg-mist px-3 py-2 font-mono text-xs text-text-dim">
            GET /api/v1/prometheus/query?query=up
          </div>
        </Card>

        <Card className="p-5">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">SECURITY HEALTH</h2>
          <div className="mt-3 space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-text-dim">Status</span>
              <Badge tone="accent">{dashQ.data?.health?.status || summary?.status || '-'}</Badge>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-text-dim">Violations</span>
              <span className="font-semibold">{summary?.security_violations ?? '-'}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-text-dim">Alerts</span>
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
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">NODE TELEMETRY</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="hud-table">
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
              {nodes.map((n: any) => (
                <tr key={n.nodeName}>
                  <td className="font-semibold text-cyan">{n.nodeName}</td>
                  <td>{formatPercent(n.cpuPercent)}</td>
                  <td>{formatPercent(n.memoryPercent)}</td>
                  <td>{n.cpuRequestsPercent || '-'}</td>
                  <td>{n.memoryRequestsPercent || '-'}</td>
                </tr>
              ))}
              {!nodes.length ? (
                <tr>
                  <td colSpan={5}>
                    <EmptyState>No node metrics available.</EmptyState>
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
