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
import { useTheme } from '@/theme/useTheme'

type Props = {
  title?: string
  query: string
  hours?: number
  step?: string
  color?: string
  /** Axis / tooltip unit, shown beside numeric ticks (e.g. cores, GiB). */
  unit?: string
  /** Format numeric tick/tooltip value only — do not append unit here. */
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
  color,
  unit,
  yFormatter = (v) => v.toFixed(2),
}: Props) {
  const { clusterId } = useCluster()
  const { theme } = useTheme()
  const stroke = color || theme.colors.primary
  const dim = theme.colors.textDim
  const axis = theme.colors.primaryDim
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
        Prometheus 未配置。管理员可到 Settings → Prometheus 填写 URL，或开启公网 Showcase
        模拟时序。
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

  const gradId = `prom-${theme.id}-${stroke.replace('#', '')}`
  const formatValue = (v: number) => {
    const n = yFormatter(v)
    return unit ? `${n} ${unit}` : n
  }

  return (
    <div className="h-64 w-full min-w-0">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 8, right: 12, left: 4, bottom: 4 }}>
          <defs>
            <linearGradient id={gradId} x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor={stroke} stopOpacity={0.4} />
              <stop offset="95%" stopColor={stroke} stopOpacity={0} />
            </linearGradient>
          </defs>
          <CartesianGrid stroke={theme.colors.primary} strokeOpacity={0.12} vertical={false} />
          <XAxis
            dataKey="t"
            tick={{ fontSize: 11, fill: dim }}
            stroke={axis}
            minTickGap={28}
            tickMargin={8}
          />
          <YAxis
            tick={{ fontSize: 11, fill: dim }}
            stroke={axis}
            tickFormatter={(v) => yFormatter(Number(v))}
            width={56}
            tickMargin={8}
            domain={['auto', 'auto']}
          />
          <Tooltip
            formatter={(value) => [formatValue(Number(value ?? 0)), 'value']}
            labelFormatter={(label) => `Time ${label}`}
            contentStyle={{
              background: theme.terminal.bg,
              border: `1px solid ${theme.colors.primary}40`,
              borderRadius: 4,
              color: theme.terminal.fg,
            }}
          />
          <Area
            type="monotone"
            dataKey="v"
            stroke={stroke}
            fill={`url(#${gradId})`}
            strokeWidth={2}
            isAnimationActive={false}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  )
}
