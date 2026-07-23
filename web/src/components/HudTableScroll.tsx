import type { ReactNode, TableHTMLAttributes } from 'react'
import { Card } from '@/components/ui'
import { cn } from '@/lib/utils'

/**
 * Scrollport for hud tables: sticky header (and optional pinned first column).
 * Height follows content up to maxHeight; does not stretch empty space.
 * Horizontal scroll only when `wide` — normal tables stay within the card.
 */
export function HudTableScroll({
  children,
  className,
  pinFirst = false,
  wide = false,
  maxHeightClass = 'max-h-[min(70vh,calc(100dvh-12.5rem))]',
}: {
  children: ReactNode
  className?: string
  /** Pin the first column while scrolling horizontally (use with wide). */
  pinFirst?: boolean
  /** Allow horizontal scroll for metric-heavy tables. */
  wide?: boolean
  maxHeightClass?: string
}) {
  const enableWide = wide || pinFirst
  return (
    <div
      className={cn(
        'hud-table-scroll overscroll-contain',
        maxHeightClass,
        enableWide && 'hud-table-scroll--wide',
        pinFirst && 'hud-table-scroll--pin-first',
        className,
      )}
    >
      {children}
    </div>
  )
}

/**
 * List-page table shell:
 * - few rows → border hugs content
 * - many rows → grows to bottom of remaining viewport, then scrolls (sticky header)
 * - normal tables → no horizontal scrollbar (ellipsis)
 * - wide/pinFirst → horizontal scroll when columns overflow
 */
export function HudTablePanel({
  children,
  className,
  pinFirst = false,
  wide = false,
}: {
  children: ReactNode
  className?: string
  pinFirst?: boolean
  wide?: boolean
}) {
  const enableWide = wide || pinFirst
  return (
    <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
      <Card
        className={cn(
          /* h-fit + max-h-full: hug rows when few; cap to viewport when many */
          'hud-table-scroll h-fit max-h-full w-full overscroll-contain p-0',
          enableWide && 'hud-table-scroll--wide',
          pinFirst && 'hud-table-scroll--pin-first',
          className,
        )}
      >
        {children}
      </Card>
    </div>
  )
}

/** Full-height list page frame so HudTablePanel can max out to the bottom. */
export function ListPageFrame({
  children,
  className,
}: {
  children: ReactNode
  className?: string
}) {
  return (
    <div
      className={cn(
        'flex h-full min-h-0 w-full min-w-0 flex-col overflow-hidden',
        className,
      )}
    >
      {children}
    </div>
  )
}

export function HudTable({
  className,
  pinFirst,
  wide,
  children,
  ...rest
}: TableHTMLAttributes<HTMLTableElement> & { pinFirst?: boolean; wide?: boolean }) {
  const enableWide = wide || pinFirst
  return (
    <table
      className={cn(
        'hud-table',
        enableWide && 'hud-table--wide',
        pinFirst && 'hud-table--pin-first',
        className,
      )}
      {...rest}
    >
      {children}
    </table>
  )
}
