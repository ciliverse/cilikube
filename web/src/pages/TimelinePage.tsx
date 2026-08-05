import { useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import dayjs from 'dayjs'
import { Loader2, RefreshCw } from 'lucide-react'
import { listNamespaces } from '@/api/cluster'
import { getTimeline, type TimelineRow } from '@/api/timeline'
import { useCluster } from '@/store/cluster'
import { ALL_NAMESPACES, useNamespace } from '@/store/namespace'
import { Button, EmptyState, HudSelect, PageHeader } from '@/components/ui'
import { cn } from '@/lib/utils'

type WindowKey = '15m' | '1h' | '6h'

type SpineItem = {
  id: string
  at: string
  marker: string
  title: string
  message: string
  kind: string
  name: string
  namespace: string
  appGroup: string
  href: string
  status?: string
  /** Unique Y-lane key — one resource per row so points don't stack on a single Y. */
  lane: string
  laneLabel: string
}

/** Per-resource lane (not app group) so the scatter Y-axis is staggered. */
function resourceLane(row: Pick<TimelineRow, 'namespace' | 'kind' | 'name'>) {
  const lane = `${row.namespace}/${row.kind}/${row.name}`
  const shortName =
    row.name.length > 22 ? `${row.name.slice(0, 10)}…${row.name.slice(-8)}` : row.name
  return { lane, laneLabel: `${row.kind}/${shortName}` }
}

function flattenSpine(rows: TimelineRow[]): SpineItem[] {
  const items: SpineItem[] = []
  for (const row of rows) {
    const { lane, laneLabel } = resourceLane(row)
    const base = {
      kind: row.kind,
      name: row.name,
      namespace: row.namespace,
      appGroup: row.appGroup,
      href: row.href,
      lane,
      laneLabel,
    }
    for (let i = 0; i < row.events.length; i++) {
      const ev = row.events[i]
      items.push({
        ...base,
        id: `e-${row.kind}-${row.namespace}-${row.name}-${ev.at}-${i}`,
        at: ev.at,
        marker: ev.marker || 'modified',
        title: ev.reason || ev.type,
        message: ev.message || '',
      })
    }
    for (let i = 1; i < row.segments.length; i++) {
      const prev = row.segments[i - 1]
      const seg = row.segments[i]
      if (prev.status === seg.status) continue
      items.push({
        ...base,
        id: `s-${row.kind}-${row.namespace}-${row.name}-${seg.from}`,
        at: seg.from,
        marker: seg.status === 'unhealthy' || seg.status === 'degraded' ? 'warning' : 'modified',
        title: `${prev.status} → ${seg.status}`,
        message: seg.reason || '',
        status: seg.status,
      })
    }
  }
  items.sort((a, b) => dayjs(b.at).valueOf() - dayjs(a.at).valueOf())
  return items
}

function TimeScatter({
  items,
  from,
  to,
  activeId,
  onSelect,
}: {
  items: SpineItem[]
  from: string
  to: string
  activeId: string | null
  onSelect: (id: string) => void
}) {
  const { t } = useTranslation()
  const t0 = dayjs(from).valueOf()
  const t1 = dayjs(to).valueOf()
  const span = Math.max(1, t1 - t0)

  const lanes = useMemo(() => {
    const counts = new Map<string, { count: number; label: string }>()
    for (const it of items) {
      const cur = counts.get(it.lane)
      if (cur) cur.count += 1
      else counts.set(it.lane, { count: 1, label: it.laneLabel })
    }
    return [...counts.entries()]
      .sort((a, b) => b[1].count - a[1].count || a[1].label.localeCompare(b[1].label))
      .slice(0, 14)
      .map(([key, v]) => ({ key, label: v.label }))
      .sort((a, b) => a.label.localeCompare(b.label))
  }, [items])

  const laneIndex = useMemo(() => {
    const m = new Map<string, number>()
    lanes.forEach((l, i) => m.set(l.key, i))
    return m
  }, [lanes])

  const ticks = useMemo(() => {
    const n = 5
    const out: { x: number; label: string }[] = []
    for (let i = 0; i <= n; i++) {
      const ts = t0 + (span * i) / n
      out.push({ x: (i / n) * 100, label: dayjs(ts).format('HH:mm') })
    }
    return out
  }, [t0, span])

  // Same-lane same-time points get a tiny Y jitter so they don't fully overlap.
  const plotItems = useMemo(() => {
    const bucket = new Map<string, number>()
    return items
      .filter((it) => laneIndex.has(it.lane))
      .map((it) => {
        const slot = `${it.lane}|${dayjs(it.at).format('HH:mm:ss')}`
        const n = bucket.get(slot) || 0
        bucket.set(slot, n + 1)
        return { ...it, jitter: n }
      })
  }, [items, laneIndex])

  const rowH = 28
  const padTop = 8
  const labelW = 156
  const chartH = padTop + Math.max(1, lanes.length) * rowH + 8

  if (lanes.length === 0) return null

  return (
    <div className="tlv-scatter">
      <div className="tlv-scatter-head">
        <span>{t('timeline.scatterTitle')}</span>
        <em>{t('timeline.scatterHint')}</em>
      </div>

      <div className="tlv-scatter-frame" style={{ height: chartH }}>
        <div className="tlv-scatter-ylabels" style={{ width: labelW, paddingTop: padTop }}>
          {lanes.map((lane) => (
            <button
              key={lane.key}
              type="button"
              className="tlv-scatter-ylabel"
              style={{ height: rowH }}
              title={lane.key}
              onClick={() => {
                const hit = items.find((it) => it.lane === lane.key)
                if (hit) onSelect(hit.id)
              }}
            >
              {lane.label}
            </button>
          ))}
        </div>
        <div className="tlv-scatter-plot" style={{ height: chartH }}>
          <svg
            className="tlv-scatter-svg"
            width="100%"
            height={chartH}
            viewBox={`0 0 860 ${chartH}`}
            preserveAspectRatio="none"
            aria-hidden
          >
            {lanes.map((lane, i) => {
              const y = padTop + i * rowH + rowH / 2
              return (
                <line
                  key={`g-${lane.key}`}
                  x1={0}
                  x2={860}
                  y1={y}
                  y2={y}
                  className="tlv-scatter-guide"
                />
              )
            })}
            <line
              x1={858}
              x2={858}
              y1={padTop - 2}
              y2={padTop + lanes.length * rowH}
              className="tlv-scatter-now"
            />
          </svg>
          <div className="tlv-scatter-marks" style={{ height: chartH }}>
            {plotItems.map((it) => {
              const li = laneIndex.get(it.lane)!
              const xPct = ((dayjs(it.at).valueOf() - t0) / span) * 100
              const yPx =
                padTop + li * rowH + rowH / 2 + (it.jitter % 3) * 3 - (it.jitter ? 3 : 0)
              return (
                <button
                  key={it.id}
                  type="button"
                  aria-label={`${dayjs(it.at).format('HH:mm:ss')} ${it.title}`}
                  className={cn(
                    'tlv-scatter-mark',
                    `is-${it.marker || 'modified'}`,
                    activeId === it.id && 'is-active',
                  )}
                  style={{ left: `${xPct}%`, top: yPx }}
                  title={`${dayjs(it.at).format('HH:mm:ss')} · ${it.laneLabel} · ${it.title}`}
                  onClick={(e) => {
                    e.stopPropagation()
                    onSelect(it.id)
                  }}
                />
              )
            })}
          </div>
        </div>
      </div>

      {/* Axis lives OUTSIDE fixed-height chart so labels are never clipped */}
      <div className="tlv-scatter-axis">
        <div className="tlv-scatter-axis-gutter" style={{ width: labelW }} />
        <div className="tlv-scatter-axis-ticks">
          {ticks.map((tk) => (
            <span key={tk.label + tk.x}>{tk.label}</span>
          ))}
          <em>{t('timeline.now')}</em>
        </div>
      </div>
    </div>
  )
}

export function TimelinePage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { clusterId } = useCluster()
  const { namespace: topNs } = useNamespace()
  const [windowKey, setWindowKey] = useState<WindowKey>('15m')
  const [filterNs, setFilterNs] = useState('')
  const [filterService, setFilterService] = useState('')
  const [activeId, setActiveId] = useState<string | null>(null)
  const [layout, setLayout] = useState<'both' | 'scatter' | 'spine'>('both')

  const effectiveNs =
    filterNs || (topNs && topNs !== ALL_NAMESPACES ? topNs : undefined)

  const nsQ = useQuery({
    queryKey: ['namespaces', clusterId],
    queryFn: listNamespaces,
    enabled: Boolean(clusterId),
  })

  const tlQ = useQuery({
    queryKey: ['timeline', clusterId, effectiveNs, windowKey],
    queryFn: () =>
      getTimeline({
        namespace: effectiveNs,
        window: windowKey,
        groupBy: 'app',
      }),
    enabled: Boolean(clusterId),
    refetchInterval: 20_000,
  })

  const rows = useMemo(
    () => tlQ.data?.groups.flatMap((g) => g.rows) ?? [],
    [tlQ.data],
  )

  const serviceOptions = useMemo(() => {
    const set = new Set<string>()
    for (const r of rows) {
      if (r.appGroup && r.appGroup !== '_ungrouped') set.add(r.appGroup)
      if (r.kind === 'service') set.add(r.name)
    }
    return [...set].sort((a, b) => a.localeCompare(b))
  }, [rows])

  const spine = useMemo(() => {
    let list = flattenSpine(rows)
    if (filterService) {
      const q = filterService.toLowerCase()
      list = list.filter(
        (it) =>
          it.appGroup.toLowerCase() === q ||
          it.name.toLowerCase() === q ||
          it.name.toLowerCase().includes(q),
      )
    }
    return list
  }, [rows, filterService])

  const focusItem = (id: string) => {
    setActiveId(id)
    if (layout === 'scatter') setLayout('both')
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        document.getElementById(`spine-${id}`)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
      })
    })
  }

  return (
    <div className="tlv-page">
      <PageHeader
        title={t('timeline.title')}
        subtitle={t('timeline.subtitleSpine')}
        action={
          <Button variant="ghost" type="button" onClick={() => void tlQ.refetch()}>
            <RefreshCw className={cn('h-3.5 w-3.5', tlQ.isFetching && 'animate-spin')} />
            {t('common.refresh')}
          </Button>
        }
      />

      <div className="tlv-toolbar">
        <div className="topo-seg tl-window">
          {(['15m', '1h', '6h'] as WindowKey[]).map((w) => (
            <button
              key={w}
              type="button"
              className={cn(windowKey === w && 'is-active')}
              onClick={() => setWindowKey(w)}
            >
              {t(`timeline.window.${w}`)}
            </button>
          ))}
        </div>
        <HudSelect
          aria-label={t('timeline.filterNs')}
          className="w-auto min-w-[10rem]"
          value={filterNs}
          onChange={setFilterNs}
          searchableWhen={0}
          options={[
            { value: '', label: t('timeline.allNamespaces') },
            ...(nsQ.data || []).map((ns) => ({ value: ns, label: ns })),
          ]}
        />
        <HudSelect
          aria-label={t('timeline.filterService')}
          className="w-auto min-w-[12rem]"
          value={filterService}
          onChange={setFilterService}
          searchableWhen={0}
          options={[
            { value: '', label: t('timeline.allServices') },
            ...serviceOptions.map((s) => ({ value: s, label: s })),
          ]}
        />
        <div className="topo-seg tl-window">
          {([
            ['both', 'timeline.layout.both'],
            ['scatter', 'timeline.layout.scatter'],
            ['spine', 'timeline.layout.spine'],
          ] as const).map(([k, label]) => (
            <button
              key={k}
              type="button"
              className={cn(layout === k && 'is-active')}
              onClick={() => setLayout(k)}
            >
              {t(label)}
            </button>
          ))}
        </div>
        <div className="tlv-legend">
          <span className="tlv-leg is-warning">{t('timeline.marker.warning')}</span>
          <span className="tlv-leg is-modified">{t('timeline.marker.modified')}</span>
          <span className="tlv-leg is-created">{t('timeline.marker.created')}</span>
          <span className="tlv-leg is-deleted">{t('timeline.marker.deleted')}</span>
        </div>
      </div>

      {tlQ.data?.sampling?.provisional ? (
        <p className="tlv-note">{t('timeline.provisionalSpine')}</p>
      ) : null}

      {tlQ.isLoading ? (
        <div className="tlv-loading">
          <Loader2 className="h-5 w-5 animate-spin text-cyan" />
          {t('common.loading')}
        </div>
      ) : tlQ.isError ? (
        <EmptyState>{(tlQ.error as Error)?.message || t('timeline.loadFailed')}</EmptyState>
      ) : spine.length === 0 ? (
        <EmptyState>{t('timeline.emptySpine')}</EmptyState>
      ) : (
        <div className="tlv-shell">
          {(layout === 'both' || layout === 'scatter') && tlQ.data ? (
            <TimeScatter
              items={spine}
              from={tlQ.data.from}
              to={tlQ.data.to}
              activeId={activeId}
              onSelect={focusItem}
            />
          ) : null}

          {(layout === 'both' || layout === 'spine') && (
            <div className="tlv-spine">
              <div className="tlv-now-row">
                <div className="tlv-time">{t('timeline.now')}</div>
                <div className="tlv-node is-now" />
                <div className="tlv-card is-now">{dayjs(tlQ.data!.to).format('HH:mm:ss')}</div>
              </div>
              {spine.map((item) => (
                <button
                  key={item.id}
                  id={`spine-${item.id}`}
                  type="button"
                  className={cn('tlv-item', activeId === item.id && 'is-active')}
                  onClick={() => {
                    setActiveId(item.id)
                    if (item.href) {
                      navigate(item.href, { state: { from: '/timeline' } })
                    }
                  }}
                >
                  <div className="tlv-time">
                    <strong>{dayjs(item.at).format('HH:mm:ss')}</strong>
                    <span>{dayjs(item.at).format('M/D')}</span>
                  </div>
                  <div className={cn('tlv-node', `is-${item.marker}`)} />
                  <div className={cn('tlv-card', `is-${item.marker}`)}>
                    <div className="tlv-card-top">
                      <span className={cn('tlv-mark-label', `is-${item.marker}`)}>
                        {t(`timeline.marker.${item.marker}`, { defaultValue: item.marker })}
                      </span>
                      <span className="tlv-card-title">{item.title}</span>
                      {item.status ? (
                        <span className={cn('tl-status-chip', `is-${item.status}`)}>
                          {t(`timeline.status.${item.status}`, { defaultValue: item.status })}
                        </span>
                      ) : null}
                    </div>
                    {item.message ? <p className="tlv-msg">{item.message}</p> : null}
                    <div className="tlv-meta">
                      <em>{item.kind}</em>
                      <span>{item.name}</span>
                      <span className="tlv-ns">{item.namespace}</span>
                      {item.appGroup && item.appGroup !== '_ungrouped' ? (
                        <span className="tlv-app">{item.appGroup}</span>
                      ) : null}
                    </div>
                  </div>
                </button>
              ))}
            </div>
          )}

          <div className="tlv-footer">
            {t('timeline.hintSpine')} ·{' '}
            <Link to="/events" className="text-cyan hover:underline">
              {t('nav.events')}
            </Link>
          </div>
        </div>
      )}
    </div>
  )
}
