import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { getPrometheusStatus, prometheusQueryRange } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { EmptyState } from '@/components/ui'

type Props = {
  title?: string
  query: string
  hours?: number
  step?: string
  color?: string
  yFormatter?: (v: number) => string
}

function parseMatrix(result: any): Array<{ t: string; v: number; ts: number }> {
  const series = result?.data?.result || result?.result || []
  const first = Array.isArray(series) ? series[0] : null
  const values: Array<[number | string, string]> = first?.values || []
  return values.map(([ts, val]) => {
    const n = typeof ts === 'number' ? ts : Number(ts)
    const d = new Date(n * 1000)
    const label = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
    return { ts: n, t: label, v: parseFloat(val) || 0 }
  })
}

export function PromTimeChart({
  query,
  hours = 1,
  step = '60s',
  color = '#35e6ff',
  yFormatter = (v) => v.toFixed(2),
}: Props) {
  const { clusterId } = useCluster()
  const statusQ = useQuery({
    queryKey: ['prometheus-status'],
    queryFn: getPrometheusStatus,
    refetchInterval: 60_000,
  })

  const enabled = Boolean(statusQ.data?.enabled && statusQ.data?.healthy && query)

  const rangeQ = useQuery({
    queryKey: ['prom-range', clusterId, query, hours, step],
    enabled,
    refetchInterval: 30_000,
    queryFn: async () => {
      const end = new Date()
      const start = new Date(end.getTime() - hours * 3600_000)
      return prometheusQueryRange({ query, start, end, step })
    },
  })

  const data = useMemo(() => parseMatrix(rangeQ.data), [rangeQ.data])

  if (statusQ.isLoading) {
    return <EmptyState>Checking Prometheus…</EmptyState>
  }
  if (!statusQ.data?.enabled) {
    return (
      <EmptyState>
        Prometheus is not configured. Enable `prometheus.enabled` for time-series charts.
      </EmptyState>
    )
  }
  if (!statusQ.data?.healthy) {
    return <EmptyState>Prometheus is unreachable: {statusQ.data?.error || 'unhealthy'}</EmptyState>
  }
  if (rangeQ.isLoading) {
    return <EmptyState>Loading time series…</EmptyState>
  }
  if (!data.length) {
    return <EmptyState>No series returned for this query.</EmptyState>
  }

  return (
    <div className="h-56 w-full min-w-0">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
          <defs>
            <linearGradient id={`prom-${color.replace('#', '')}`} x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor={color} stopOpacity={0.4} />
              <stop offset="95%" stopColor={color} stopOpacity={0} />
            </linearGradient>
          </defs>
          <CartesianGrid stroke="rgba(53,230,255,0.08)" vertical={false} />
          <XAxis dataKey="t" tick={{ fontSize: 10, fill: '#6f98a3' }} stroke="#1a7f92" minTickGap={24} />
          <YAxis
            tick={{ fontSize: 11, fill: '#6f98a3' }}
            stroke="#1a7f92"
            tickFormatter={yFormatter}
            width={48}
          />
          <Tooltip
            formatter={(value) => [yFormatter(Number(value ?? 0)), 'value']}
            labelFormatter={(label) => `Time ${label}`}
            contentStyle={{
              background: '#071016',
              border: '1px solid rgba(53,230,255,0.25)',
              borderRadius: 4,
              color: '#cfeef5',
            }}
          />
          <Area
            type="monotone"
            dataKey="v"
            stroke={color}
            fill={`url(#prom-${color.replace('#', '')})`}
            strokeWidth={2}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  )
}
