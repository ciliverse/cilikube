import { cn } from '@/lib/utils'

/** Compact Kubernetes wheel mark for fleet / cluster chrome. */
export function KubernetesMark({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 64 64"
      className={cn('shrink-0', className)}
      aria-hidden
      focusable="false"
    >
      <circle cx="32" cy="32" r="30" fill="currentColor" opacity="0.12" />
      <g fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinejoin="round">
        <polygon points="32,8 52,19.5 52,44.5 32,56 12,44.5 12,19.5" />
        <polygon points="32,18 44.5,25 44.5,39 32,46 19.5,39 19.5,25" opacity="0.85" />
      </g>
      <circle cx="32" cy="32" r="5.5" fill="currentColor" />
      {[0, 60, 120, 180, 240, 300].map((deg) => {
        const rad = (deg * Math.PI) / 180
        const x = 32 + Math.cos(rad) * 18
        const y = 32 + Math.sin(rad) * 18
        return <circle key={deg} cx={x} cy={y} r="2.2" fill="currentColor" opacity="0.9" />
      })}
    </svg>
  )
}
