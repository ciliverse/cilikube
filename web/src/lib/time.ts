/** Parse K8s-style timestamps (RFC3339 string, millis, or metav1-like objects). */
export function parseResourceTime(value: unknown): number | null {
  if (value == null || value === '') return null
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value < 1e12 ? value * 1000 : value
  }
  if (typeof value === 'string') {
    const ms = Date.parse(value)
    if (!Number.isNaN(ms) && ms > 0) return ms
    return null
  }
  if (typeof value === 'object') {
    const o = value as Record<string, unknown>
    if ('Time' in o) return parseResourceTime(o.Time)
    if ('seconds' in o) {
      const sec = Number(o.seconds)
      if (!Number.isFinite(sec)) return null
      const nanos = Number(o.nanos || 0)
      return sec * 1000 + Math.floor(nanos / 1e6)
    }
  }
  return null
}

/** k9s / kubectl short age: 45s / 12m / 5h / 3d / 1y */
export function formatAge(value: unknown, now = Date.now()): string {
  const ms = parseResourceTime(value)
  if (ms == null || ms <= 0) return '-'
  let sec = Math.floor((now - ms) / 1000)
  if (sec < 0) sec = 0
  if (sec < 60) return `${sec}s`
  const min = Math.floor(sec / 60)
  if (min < 60) return `${min}m`
  const hr = Math.floor(min / 60)
  if (hr < 24) return `${hr}h`
  const day = Math.floor(hr / 24)
  if (day < 365) return `${day}d`
  return `${Math.floor(day / 365)}y`
}

/** Compact absolute time for tables (local): 2026-07-22 20:33 */
export function formatCreatedAt(value: unknown): string {
  const ms = parseResourceTime(value)
  if (ms == null || ms <= 0) return '-'
  const d = new Date(ms)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

export function formatAbsoluteTime(value: unknown): string {
  const ms = parseResourceTime(value)
  if (ms == null) return ''
  try {
    return new Date(ms).toLocaleString()
  } catch {
    return ''
  }
}
