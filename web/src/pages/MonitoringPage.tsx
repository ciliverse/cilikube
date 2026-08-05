import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { motion } from 'framer-motion'
import {
  getMonitoringDashboard,
  getNodeMetrics,
  getPrometheusStatus,
} from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, EmptyState, PageHeader, StatCard } from '@/components/ui'
import { HudTable, HudTableScroll } from '@/components/HudTableScroll'
import { PromTimeChart } from '@/components/PromTimeChart'
import { shouldSkipEnterAnim } from '@/lib/motionPrefs'
import { formatPercent } from '@/lib/utils'

export function MonitoringPage() {
  const { t } = useTranslation()
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
    <div className="flex w-full min-w-0 flex-col gap-4">
      <PageHeader
        title={t('monitoring.title')}
        subtitle={t('monitoring.subtitle')}
        action={
          <Badge tone={prom?.healthy ? 'ok' : prom?.enabled ? 'danger' : 'neutral'}>
            Prometheus{' '}
            {prom?.mode === 'showcase'
              ? 'showcase'
              : prom?.enabled
                ? prom.healthy
                  ? 'healthy'
                  : 'down'
                : 'off'}
          </Badge>
        }
      />

      <div className="grid w-full grid-cols-2 gap-4 lg:grid-cols-4">
        {[
          { label: 'Active sessions', value: summary?.active_sessions ?? '-' },
          { label: 'Active users', value: summary?.active_users ?? '-' },
          { label: 'Active threats', value: summary?.active_threats ?? '-' },
          { label: 'Failed login rate', value: summary?.failed_logins_rate ?? '-' },
        ].map((card, index) => (
          <motion.div
            key={card.label}
            className="min-w-0"
            initial={shouldSkipEnterAnim() ? false : { opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: shouldSkipEnterAnim() ? 0 : index * 0.05 }}
          >
            <StatCard label={card.label} value={card.value} />
          </motion.div>
        ))}
      </div>

      <div className="grid w-full gap-4 lg:grid-cols-2 lg:items-stretch">
        <Card className="flex min-w-0 flex-col p-5">
          <div className="mb-2 flex items-center justify-between gap-2">
            <h2 className="font-display text-lg font-bold tracking-[0.12em]">CLUSTER CPU (1h)</h2>
            <Badge tone="neutral">cores</Badge>
          </div>
          <p className="mb-3 min-h-[1.25rem] truncate text-xs text-text-dim" title={prom?.url || undefined}>
            query_range · {prom?.url || 'not configured'}
          </p>
          <div className="min-w-0 flex-1">
            <PromTimeChart
              query={
                'sum(rate(container_cpu_usage_seconds_total{container!="",container!="POD"}[5m]))'
              }
              hours={1}
              unit="cores"
              yFormatter={(v) => v.toFixed(2)}
            />
          </div>
        </Card>
        <Card className="flex min-w-0 flex-col p-5">
          <div className="mb-2 flex items-center justify-between gap-2">
            <h2 className="font-display text-lg font-bold tracking-[0.12em]">CLUSTER MEMORY (1h)</h2>
            <Badge tone="neutral">GiB</Badge>
          </div>
          <p className="mb-3 min-h-[1.25rem] truncate text-xs text-text-dim" title={prom?.url || undefined}>
            query_range · {prom?.url || 'not configured'}
          </p>
          <div className="min-w-0 flex-1">
            <PromTimeChart
              query={
                'sum(container_memory_working_set_bytes{container!="",container!="POD"}) / 1024 / 1024 / 1024'
              }
              hours={1}
              color="#ffae00"
              unit="GiB"
              yFormatter={(v) => v.toFixed(1)}
            />
          </div>
        </Card>
      </div>

      <div className="grid w-full gap-4 lg:grid-cols-[minmax(0,1.7fr)_minmax(0,1fr)]">
        <Card className="min-w-0 p-5">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">PROMETHEUS STATUS</h2>
          <p className="mt-1 text-sm text-text-dim">
            Endpoint: {prom?.url || 'not configured'}
          </p>
          {prom?.error ? (
            <div className="mt-3 rounded border border-warn/40 bg-warn/10 px-3 py-2 text-sm text-warn">
              {prom.error}
            </div>
          ) : (
            <div className="mt-3 text-sm text-text-dim">
              {prom?.mode === 'showcase' ? t('monitoring.showcaseHint') : t('monitoring.promHealthyHint')}
            </div>
          )}
        </Card>

        <Card className="min-w-0 p-5">
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

      <Card className="w-full min-w-0 overflow-hidden">
        <div className="border-b border-line px-5 py-4">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">NODE RESOURCE USAGE</h2>
          <p className="mt-1 text-sm text-text-dim">
            Point-in-time metrics-server snapshot (CPU / memory / requests)
          </p>
        </div>
        <HudTableScroll>
          <HudTable className="table-fixed">
            <colgroup>
              <col className="w-[40%]" />
              <col className="w-[15%]" />
              <col className="w-[15%]" />
              <col className="w-[15%]" />
              <col className="w-[15%]" />
            </colgroup>
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
                  <td className="truncate font-semibold text-cyan" title={n.nodeName}>
                    {n.nodeName}
                  </td>
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
          </HudTable>
        </HudTableScroll>
      </Card>
    </div>
  )
}
