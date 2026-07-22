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
      'bg-accent text-white hover:bg-accent-bright shadow-[0_10px_30px_-12px_rgba(15,118,110,0.7)]',
    ghost: 'bg-transparent text-ink-soft hover:bg-black/5',
    danger: 'bg-danger text-white hover:opacity-90',
    outline: 'border border-line bg-panel text-ink hover:border-accent hover:text-accent',
  }
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center gap-2 rounded-xl px-4 py-2.5 text-sm font-semibold transition active:scale-[0.98] disabled:opacity-50',
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
        'w-full rounded-xl border border-line bg-white/80 px-3.5 py-2.5 text-sm outline-none transition placeholder:text-ink-soft/50 focus:border-accent focus:ring-4 focus:ring-accent/15',
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
  return (
    <div
      className={cn(
        'rounded-2xl border border-white/70 bg-panel/85 shadow-[0_20px_60px_-40px_rgba(16,24,40,0.45)] backdrop-blur',
        className,
      )}
    >
      {children}
    </div>
  )
}

export function Badge({
  children,
  tone = 'neutral',
}: {
  children: ReactNode
  tone?: 'neutral' | 'ok' | 'warn' | 'danger' | 'accent'
}) {
  const tones = {
    neutral: 'bg-mist text-ink-soft',
    ok: 'bg-emerald-50 text-ok',
    warn: 'bg-amber-50 text-warn',
    danger: 'bg-red-50 text-danger',
    accent: 'bg-accent-soft text-accent',
  }
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold',
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
        <h1 className="font-display text-3xl font-extrabold tracking-tight text-ink md:text-4xl">
          {title}
        </h1>
        {subtitle ? <p className="mt-1.5 text-sm text-ink-soft">{subtitle}</p> : null}
      </div>
      {action}
    </div>
  )
}
