import { cn } from '@/lib/utils'
import type { PodMetricsItem } from '@/api/cluster'

function toneForRatio(ratio?: number): string {
  if (ratio == null || !Number.isFinite(ratio) || ratio <= 0) return 'text-text-dim'
  if (ratio >= 1) return 'text-danger'
  if (ratio >= 0.8) return 'text-warn'
  return 'text-text'
}

function MetricValue({
  value,
  title,
  ratio,
}: {
  value?: string
  title?: string
  ratio?: number
}) {
  const text = value && value !== '' ? value : '-'
  return (
    <span
      className={cn('font-mono text-[12px] tabular-nums', toneForRatio(ratio))}
      title={title}
    >
      {text}
    </span>
  )
}

/** Compact k9s-style percent cell (%CPU/R, %MEM/L, …). */
export function PercentCell({
  percent,
  ratio,
  hint,
}: {
  percent?: string
  ratio?: number
  hint?: string
}) {
  return <MetricValue value={percent} ratio={ratio} title={hint} />
}

export function podMetricColumns(get: (item: any) => PodMetricsItem | undefined) {
  return [
    {
      key: 'cpu',
      header: 'CPU',
      title: 'Current CPU usage (metrics-server)',
      render: (item: any) => {
        const m = get(item)
        return (
          <MetricValue
            value={m?.cpu}
            title={m ? `req ${m.cpuRequest || '-'} / lim ${m.cpuLimit || '-'}` : undefined}
          />
        )
      },
    },
    {
      key: 'cpuR',
      header: '%CPU/R',
      title: 'CPU usage as % of Request',
      render: (item: any) => {
        const m = get(item)
        return (
          <PercentCell
            percent={m?.cpuRequestPercent}
            ratio={m?.cpuRequestRatio}
            hint={m?.cpuRequest ? `request ${m.cpuRequest}` : 'no CPU request'}
          />
        )
      },
    },
    {
      key: 'cpuL',
      header: '%CPU/L',
      title: 'CPU usage as % of Limit',
      render: (item: any) => {
        const m = get(item)
        return (
          <PercentCell
            percent={m?.cpuLimitPercent}
            ratio={m?.cpuLimitRatio}
            hint={m?.cpuLimit ? `limit ${m.cpuLimit}` : 'no CPU limit'}
          />
        )
      },
    },
    {
      key: 'mem',
      header: 'MEM',
      title: 'Current memory usage (metrics-server)',
      render: (item: any) => {
        const m = get(item)
        return (
          <MetricValue
            value={m?.memory}
            title={m ? `req ${m.memoryRequest || '-'} / lim ${m.memoryLimit || '-'}` : undefined}
          />
        )
      },
    },
    {
      key: 'memR',
      header: '%MEM/R',
      title: 'Memory usage as % of Request',
      render: (item: any) => {
        const m = get(item)
        return (
          <PercentCell
            percent={m?.memoryRequestPercent}
            ratio={m?.memoryRequestRatio}
            hint={m?.memoryRequest ? `request ${m.memoryRequest}` : 'no memory request'}
          />
        )
      },
    },
    {
      key: 'memL',
      header: '%MEM/L',
      title: 'Memory usage as % of Limit',
      render: (item: any) => {
        const m = get(item)
        return (
          <PercentCell
            percent={m?.memoryLimitPercent}
            ratio={m?.memoryLimitRatio}
            hint={m?.memoryLimit ? `limit ${m.memoryLimit}` : 'no memory limit'}
          />
        )
      },
    },
  ]
}
