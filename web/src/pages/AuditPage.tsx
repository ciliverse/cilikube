import { useMemo, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiGet } from '@/lib/api'
import { useAuth } from '@/store/auth'
import { Badge, Button, Card, EmptyState, PageHeader, StatCard } from '@/components/ui'
import { HudTable, HudTableScroll } from '@/components/HudTableScroll'

/** Format Date as datetime-local value (local timezone, minute precision). */
function toLocalInput(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function defaultRange(hours: number): { start: string; end: string } {
  const end = new Date()
  const start = new Date(end.getTime() - hours * 3_600_000)
  return { start: toLocalInput(start), end: toLocalInput(end) }
}

function toRFC3339(localValue: string): string {
  if (!localValue) return ''
  const d = new Date(localValue)
  if (Number.isNaN(d.getTime())) return ''
  return d.toISOString()
}

/** Go-safe duration only (hours), never Nd. */
function periodFromRange(start: string, end: string): string {
  const a = new Date(start).getTime()
  const b = new Date(end).getTime()
  if (!a || !b || b <= a) return '24h'
  const hours = Math.max(1, Math.round((b - a) / 3_600_000))
  return `${hours}h`
}

function formatLocalWindow(startLocal: string, endLocal: string): string {
  const a = new Date(startLocal)
  const b = new Date(endLocal)
  if (Number.isNaN(a.getTime()) || Number.isNaN(b.getTime())) return ''
  const opts: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }
  return `${a.toLocaleString(undefined, opts)} → ${b.toLocaleString(undefined, opts)}`
}

const PRESETS: { label: string; hours: number }[] = [
  { label: '1h', hours: 1 },
  { label: '24h', hours: 24 },
  { label: '7d', hours: 24 * 7 },
  { label: '30d', hours: 24 * 30 },
]

export function AuditPage() {
  const { isAdmin } = useAuth()
  const [action, setAction] = useState('')
  const [userId, setUserId] = useState('')
  const [page, setPage] = useState(1)
  const initial = useMemo(() => defaultRange(24), [])
  const [startLocal, setStartLocal] = useState(initial.start)
  const [endLocal, setEndLocal] = useState(initial.end)

  const startTime = toRFC3339(startLocal)
  const endTime = toRFC3339(endLocal)
  const timesReady = Boolean(startTime && endTime && new Date(endTime) > new Date(startTime))

  const applyPreset = (hours: number) => {
    const range = defaultRange(hours)
    setStartLocal(range.start)
    setEndLocal(range.end)
    setPage(1)
  }

  const logsQ = useQuery({
    queryKey: ['audit-logs', page, action, userId, startTime, endTime],
    enabled: isAdmin,
    queryFn: () => {
      const params: Record<string, unknown> = {
        page,
        page_size: 50,
      }
      if (action.trim()) params.action = action.trim()
      if (userId.trim()) params.user_id = Number(userId.trim()) || userId.trim()
      if (startTime) params.start_time = startTime
      if (endTime) params.end_time = endTime
      return apiGet<{ logs?: any[]; items?: any[]; total?: number; page?: number; page_size?: number }>(
        '/api/v1/audit/logs',
        params,
      )
    },
  })

  const reportQ = useQuery({
    queryKey: ['audit-report', startTime, endTime, userId],
    enabled: isAdmin && timesReady,
    queryFn: () => {
      const params: Record<string, unknown> = {
        start_time: startTime,
        end_time: endTime,
      }
      if (userId.trim()) params.user_id = Number(userId.trim()) || userId.trim()
      return apiGet<any>('/api/v1/audit/report', params)
    },
  })

  const metricsQ = useQuery({
    queryKey: ['audit-metrics', startTime, endTime],
    enabled: isAdmin && timesReady,
    queryFn: () =>
      apiGet<any>('/api/v1/audit/metrics', {
        period: periodFromRange(startTime, endTime),
        start_time: startTime,
        end_time: endTime,
      }),
  })

  const logs = logsQ.data?.logs || logsQ.data?.items || (Array.isArray(logsQ.data) ? logsQ.data : [])
  const total = logsQ.data?.total ?? (logs as any[]).length
  const pageSize = logsQ.data?.page_size || 50
  const totalPages = Math.max(1, Math.ceil(Number(total) / pageSize))

  const topActions = useMemo(() => {
    const summary = reportQ.data?.action_summary || {}
    return Object.entries(summary)
      .sort((a, b) => Number(b[1]) - Number(a[1]))
      .slice(0, 5)
  }, [reportQ.data])

  if (!isAdmin) {
    return (
      <div className="rounded border border-warn/40 bg-warn/10 px-5 py-8 text-sm text-warn">
        Admin privileges required to view audit logs.
      </div>
    )
  }

  return (
    <div className="flex w-full min-w-0 flex-col gap-4">
      <PageHeader title="AUDIT LOGS" subtitle="API and security activity" />

      <Card className="space-y-3 p-5">
        <div className="flex flex-wrap items-center gap-2">
          <span className="hud-label">Range</span>
          {PRESETS.map((p) => (
            <Button
              key={p.label}
              type="button"
              variant="outline"
              className="px-2.5 py-1 text-xs tracking-[0.12em] uppercase"
              onClick={() => applyPreset(p.hours)}
            >
              {p.label}
            </Button>
          ))}
        </div>
        <div className="flex flex-wrap items-end gap-3">
          <label className="block space-y-1">
            <span className="hud-label">Action</span>
            <input
              className="hud-field w-40"
              value={action}
              onChange={(e) => {
                setPage(1)
                setAction(e.target.value)
              }}
              placeholder="login"
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">User ID</span>
            <input
              className="hud-field w-28"
              value={userId}
              onChange={(e) => {
                setPage(1)
                setUserId(e.target.value)
              }}
              placeholder="1"
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Start</span>
            <input
              type="datetime-local"
              lang="en"
              step={60}
              className="hud-field"
              value={startLocal}
              onChange={(e) => {
                setPage(1)
                setStartLocal(e.target.value)
              }}
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">End</span>
            <input
              type="datetime-local"
              lang="en"
              step={60}
              className="hud-field"
              value={endLocal}
              onChange={(e) => {
                setPage(1)
                setEndLocal(e.target.value)
              }}
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Page</span>
            <input
              type="number"
              min={1}
              className="hud-field w-20"
              value={page}
              onChange={(e) => setPage(Math.max(1, Number(e.target.value) || 1))}
            />
          </label>
          <Button
            variant="outline"
            className="px-3 py-1.5 text-xs"
            type="button"
            onClick={() => {
              void logsQ.refetch()
              if (timesReady) {
                void reportQ.refetch()
                void metricsQ.refetch()
              }
            }}
          >
            Refresh
          </Button>
        </div>
        {timesReady ? (
          <p className="text-xs text-text-dim">
            Window (local): {formatLocalWindow(startLocal, endLocal)}
          </p>
        ) : (
          <p className="text-xs text-warn">
            Set a valid start before end to load report and metrics summary cards.
          </p>
        )}
      </Card>

      {timesReady ? (
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <StatCard label="Total events" value={reportQ.data?.total_events ?? '—'} />
          <StatCard label="Login attempts" value={reportQ.data?.login_attempts ?? '—'} />
          <StatCard
            label="Failed logins"
            value={metricsQ.data?.failed_logins ?? reportQ.data?.failed_logins ?? '—'}
          />
          <StatCard
            label="Permission denials"
            value={metricsQ.data?.permission_denials ?? reportQ.data?.permission_denials ?? '—'}
          />
        </div>
      ) : null}

      {timesReady && topActions.length ? (
        <Card className="p-5">
          <div className="hud-label mb-2">Top actions</div>
          <div className="flex flex-wrap gap-2">
            {topActions.map(([name, count]) => (
              <Badge key={name} tone="accent">
                {name}: {String(count)}
              </Badge>
            ))}
          </div>
          {metricsQ.data ? (
            <p className="mt-3 text-xs text-text-dim">
              Metrics · events/h {Number(metricsQ.data.events_per_hour || 0).toFixed(1)} · violations{' '}
              {metricsQ.data.security_violations ?? 0}
            </p>
          ) : null}
        </Card>
      ) : null}

      <Card className="overflow-hidden">
        <HudTableScroll>
          <HudTable>
            <thead>
              <tr>
                <th>Time</th>
                <th>User</th>
                <th>Action</th>
                <th>Resource</th>
                <th>IP</th>
                <th>Client</th>
              </tr>
            </thead>
            <tbody>
              {(logs as any[]).map((log: any, i: number) => {
                let details: any = log.details
                if (typeof details === 'string' && details) {
                  try {
                    details = JSON.parse(details)
                  } catch {
                    details = null
                  }
                }
                const username =
                  log.username ||
                  details?.username ||
                  (log.user_id != null ? `#${log.user_id}` : '-')
                const ua = log.user_agent || details?.user_agent || ''
                return (
                  <tr key={log.id || i}>
                    <td className="text-text-dim">{log.created_at || log.timestamp || '-'}</td>
                    <td className="font-mono text-xs">{username}</td>
                    <td>
                      <Badge tone="accent">{log.action || '-'}</Badge>
                    </td>
                    <td className="max-w-xs truncate">
                      {log.resource || log.path || '-'}
                      {log.resource_id ? `/${log.resource_id}` : ''}
                    </td>
                    <td className="font-mono text-xs text-cyan">
                      {log.ip_address || log.ip || details?.ip || '-'}
                    </td>
                    <td className="max-w-[10rem] truncate text-[11px] text-text-dim" title={ua}>
                      {ua || '-'}
                    </td>
                  </tr>
                )
              })}
              {!logsQ.isLoading && !(logs as any[]).length ? (
                <tr>
                  <td colSpan={6}>
                    <EmptyState>
                      {logsQ.isError
                        ? 'Failed to load audit logs (is /api/v1/audit registered?)'
                        : 'No audit entries yet.'}
                    </EmptyState>
                  </td>
                </tr>
              ) : null}
            </tbody>
          </HudTable>
        </HudTableScroll>
        <div className="flex items-center justify-between border-t border-line px-4 py-3 text-xs text-text-dim">
          <span>
            Page {page} / {totalPages} · {total} total
          </span>
          <div className="flex gap-2">
            <Button
              variant="ghost"
              className="px-2 py-1 text-xs"
              type="button"
              disabled={page <= 1}
              onClick={() => setPage((p) => Math.max(1, p - 1))}
            >
              Prev
            </Button>
            <Button
              variant="ghost"
              className="px-2 py-1 text-xs"
              type="button"
              disabled={page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              Next
            </Button>
          </div>
        </div>
      </Card>
    </div>
  )
}
