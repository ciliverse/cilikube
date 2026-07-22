import type { ButtonHTMLAttributes, InputHTMLAttributes, ReactNode } from 'react'
import { cn } from '@/lib/utils'

export function Button({
  className,
  variant = 'primary',
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'ghost' | 'danger' | 'outline'
}) {
  const variants = {
    primary:
      'border border-cyan/50 bg-cyan/15 text-cyan hover:bg-cyan/25 shadow-[0_0_18px_rgba(53,230,255,0.18)]',
    ghost: 'bg-transparent text-text-dim hover:bg-cyan-faint hover:text-text',
    danger: 'border border-danger/40 bg-danger/15 text-danger hover:bg-danger/25',
    outline: 'border border-line bg-panel-solid text-text hover:border-cyan hover:text-cyan',
  }
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center gap-2 rounded px-4 py-2.5 text-sm font-semibold tracking-wide transition active:scale-[0.98] disabled:opacity-50',
        variants[variant],
        className,
      )}
      {...props}
    />
  )
}

export function Input({
  className,
  ...props
}: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      className={cn(
        'w-full rounded border border-line bg-panel-solid px-3.5 py-2.5 text-sm text-text outline-none transition placeholder:text-text-dim/60 focus:border-cyan focus:shadow-[0_0_0_3px_rgba(53,230,255,0.14)]',
        className,
      )}
      {...props}
    />
  )
}

export function Card({
  className,
  children,
}: {
  className?: string
  children: ReactNode
}) {
  return <div className={cn('hud-panel rounded', className)}>{children}</div>
}

export function Panel({
  className,
  children,
}: {
  className?: string
  children: ReactNode
}) {
  return <div className={cn('hud-panel rounded', className)}>{children}</div>
}

export function Badge({
  children,
  tone = 'neutral',
}: {
  children: ReactNode
  tone?: 'neutral' | 'ok' | 'warn' | 'danger' | 'accent'
}) {
  const tones = {
    neutral: 'border-line bg-mist text-text-dim',
    ok: 'border-ok/30 bg-ok/10 text-ok',
    warn: 'border-warn/30 bg-warn/10 text-warn',
    danger: 'border-danger/30 bg-danger/10 text-danger',
    accent: 'border-cyan/30 bg-cyan-faint text-cyan',
  }
  return (
    <span
      className={cn(
        'inline-flex items-center rounded border px-2 py-0.5 text-[11px] font-semibold tracking-wider uppercase',
        tones[tone],
      )}
    >
      {children}
    </span>
  )
}

export function PageHeader({
  title,
  subtitle,
  action,
}: {
  title: string
  subtitle?: string
  action?: ReactNode
}) {
  return (
    <div className="mb-6 flex flex-wrap items-end justify-between gap-4">
      <div>
        <div className="hud-label mb-2">Control plane</div>
        <h1 className="font-display text-2xl font-bold tracking-[0.12em] text-text md:text-3xl">
          {title}
        </h1>
        {subtitle ? <p className="mt-2 text-sm text-text-dim">{subtitle}</p> : null}
      </div>
      {action}
    </div>
  )
}

export function StatCard({
  label,
  value,
  icon,
}: {
  label: string
  value: ReactNode
  icon?: ReactNode
}) {
  return (
    <Card className="p-5">
      <div className="flex items-start justify-between gap-3">
        <div>
          <div className="hud-label">{label}</div>
          <div className="mt-2 font-display text-3xl font-bold tracking-wide text-cyan">
            {value}
          </div>
        </div>
        {icon ? (
          <div className="grid h-10 w-10 place-items-center rounded border border-line bg-cyan-faint text-cyan">
            {icon}
          </div>
        ) : null}
      </div>
    </Card>
  )
}

export function ConnDot({ online = true }: { online?: boolean }) {
  return (
    <span
      className={cn(
        'inline-block h-2 w-2 rounded-full',
        online
          ? 'bg-ok shadow-[0_0_8px_rgba(77,255,176,0.7)]'
          : 'bg-danger shadow-[0_0_8px_rgba(255,77,94,0.7)]',
      )}
    />
  )
}

export function EmptyState({ children }: { children: ReactNode }) {
  return (
    <div className="px-5 py-10 text-center text-sm text-text-dim">{children}</div>
  )
}
