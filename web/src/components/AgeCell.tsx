import { useEffect, useState } from 'react'
import { formatAbsoluteTime, formatAge, formatCreatedAt, parseResourceTime } from '@/lib/time'
import { cn } from '@/lib/utils'

function useNow(intervalMs = 30_000) {
  const [now, setNow] = useState(() => Date.now())
  useEffect(() => {
    const id = window.setInterval(() => setNow(Date.now()), intervalMs)
    return () => window.clearInterval(id)
  }, [intervalMs])
  return now
}

/** Relative age — k9s/kubectl style (14h, 3d). */
export function AgeCell({
  value,
  className,
}: {
  value: unknown
  className?: string
}) {
  const now = useNow()
  const age = formatAge(value, now)
  const abs = formatAbsoluteTime(value)
  const missing = age === '-' || parseResourceTime(value) == null

  return (
    <span
      className={cn(
        'font-mono text-[12px] tabular-nums tracking-wide',
        missing ? 'text-text-dim' : 'text-text',
        className,
      )}
      title={abs || undefined}
    >
      {age}
    </span>
  )
}

/** Absolute creation time for table columns. */
export function CreatedCell({
  value,
  className,
}: {
  value: unknown
  className?: string
}) {
  const text = formatCreatedAt(value)
  const abs = formatAbsoluteTime(value)
  const missing = text === '-'

  return (
    <span
      className={cn(
        'font-mono text-[11px] tabular-nums',
        missing ? 'text-text-dim' : 'text-text',
        className,
      )}
      title={abs || undefined}
    >
      {missing ? '-' : text}
    </span>
  )
}

/**
 * k9s-inspired stacked time: Age on top, Created underneath.
 * Use when horizontal space is tight (optional alternative to two columns).
 */
export function AgeCreatedStack({ value }: { value: unknown }) {
  return (
    <div className="flex flex-col leading-tight">
      <AgeCell value={value} />
      <CreatedCell value={value} className="mt-0.5 opacity-90" />
    </div>
  )
}
