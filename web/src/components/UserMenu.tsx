import { useEffect, useRef, useState, type CSSProperties } from 'react'
import { Link } from 'react-router-dom'
import { ChevronDown, LogOut, Palette, Settings, Type, UserRound } from 'lucide-react'
import { useAuth } from '@/store/auth'
import { useTheme } from '@/theme/useTheme'
import { useFont } from '@/theme/useFont'
import { switchTheme } from '@/theme/switchTheme'
import { cn } from '@/lib/utils'

export function UserMenu({
  primaryRole,
  isViewerOnly,
}: {
  primaryRole: string
  isViewerOnly: boolean
}) {
  const { user, logout, isAdmin } = useAuth()
  const { themeId, themes } = useTheme()
  const { fontId, fonts, setFont } = useFont()
  const [open, setOpen] = useState(false)
  const rootRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return
    const onDoc = (e: MouseEvent) => {
      if (!rootRef.current?.contains(e.target as Node)) setOpen(false)
    }
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setOpen(false)
    }
    document.addEventListener('mousedown', onDoc)
    document.addEventListener('keydown', onKey)
    return () => {
      document.removeEventListener('mousedown', onDoc)
      document.removeEventListener('keydown', onKey)
    }
  }, [open])

  const roleTone =
    primaryRole === 'admin'
      ? 'text-orange'
      : primaryRole === 'editor'
        ? 'text-cyan'
        : 'text-text-dim'

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        className="flex h-9 items-center gap-1.5 rounded border border-transparent px-1.5 text-xs text-text-dim transition hover:border-line hover:bg-mist hover:text-text sm:h-auto sm:gap-2 sm:px-2 sm:py-1"
        aria-haspopup="menu"
        aria-expanded={open}
        aria-label={`${user?.username || 'operator'} / ${primaryRole}`}
        onClick={() => setOpen((v) => !v)}
      >
        <UserRound className="h-4 w-4 shrink-0 text-cyan sm:hidden" />
        <span className="hidden max-w-[7rem] truncate text-text sm:inline">
          {user?.username || 'operator'}
        </span>
        <span className="hidden text-line sm:inline">/</span>
        <span className={cn('hidden sm:inline', roleTone)}>{primaryRole}</span>
        {isViewerOnly ? (
          <span className="hidden rounded border border-line px-1.5 py-0.5 text-[10px] tracking-wider uppercase md:inline">
            read-only
          </span>
        ) : null}
        <ChevronDown className={cn('h-3.5 w-3.5 transition', open && 'rotate-180 text-cyan')} />
      </button>

      {open ? (
        <div
          role="menu"
          className="absolute right-0 z-[210] mt-1 max-h-[min(70dvh,28rem)] min-w-[min(220px,calc(100vw-1rem))] overflow-y-auto overscroll-contain rounded border border-line bg-panel-solid py-1 shadow-[0_8px_28px_rgba(0,0,0,0.45)]"
        >
          <div className="border-b border-line px-3 py-2 sm:hidden">
            <div className="truncate text-[12px] font-semibold text-text">
              {user?.username || 'operator'}
            </div>
            <div className={cn('mt-0.5 text-[11px] uppercase tracking-wider', roleTone)}>
              {primaryRole}
              {isViewerOnly ? ' · read-only' : ''}
            </div>
          </div>
          <Link
            role="menuitem"
            to="/profile"
            className="flex min-h-11 items-center gap-2 px-3 py-2.5 text-[13px] text-text hover:bg-mist sm:min-h-0 sm:py-2 sm:text-[12px]"
            onClick={() => setOpen(false)}
          >
            <UserRound className="h-3.5 w-3.5 text-cyan" />
            Profile
          </Link>
          {isAdmin ? (
            <Link
              role="menuitem"
              to="/admin/settings"
              className="flex min-h-11 items-center gap-2 px-3 py-2.5 text-[13px] text-text hover:bg-mist sm:min-h-0 sm:py-2 sm:text-[12px]"
              onClick={() => setOpen(false)}
            >
              <Settings className="h-3.5 w-3.5 text-cyan" />
              Settings
            </Link>
          ) : null}

          <div className="my-1 border-t border-line" />
          <div className="flex items-center gap-2 px-3 pt-2 text-[10px] tracking-[0.14em] text-text-dim uppercase">
            <Palette className="h-3 w-3 text-cyan" />
            Theme
          </div>
          <div className="theme-swatches" role="group" aria-label="Theme">
            {themes.map((t) => (
              <button
                key={t.id}
                type="button"
                title={`${t.name} (${t.mode})`}
                aria-label={t.name}
                aria-pressed={themeId === t.id}
                className={cn('theme-swatch', themeId === t.id && 'active')}
                style={
                  {
                    '--swatch-bg': t.colors.bg,
                    '--swatch-accent': t.colors.primary,
                  } as CSSProperties
                }
                onClick={() => switchTheme(t.id)}
              />
            ))}
          </div>

          <div className="my-1 border-t border-line" />
          <div className="flex items-center gap-2 px-3 pt-2 text-[10px] tracking-[0.14em] text-text-dim uppercase">
            <Type className="h-3 w-3 text-cyan" />
            Font
          </div>
          <div className="flex flex-col gap-0.5 px-2 pb-2 pt-1">
            {fonts.map((f) => (
              <button
                key={f.id}
                type="button"
                role="menuitemradio"
                aria-checked={fontId === f.id}
                className={cn(
                  'min-h-10 rounded px-2 py-2 text-left text-[13px] transition sm:min-h-0 sm:py-1.5 sm:text-[12px]',
                  fontId === f.id
                    ? 'bg-mist text-cyan'
                    : 'text-text hover:bg-mist hover:text-text',
                )}
                style={{ fontFamily: f.sans }}
                onClick={() => setFont(f.id)}
              >
                {f.name}
              </button>
            ))}
          </div>

          <div className="my-1 border-t border-line" />
          <button
            type="button"
            role="menuitem"
            className="flex min-h-11 w-full items-center gap-2 px-3 py-2.5 text-left text-[13px] text-danger hover:bg-mist sm:min-h-0 sm:py-2 sm:text-[12px]"
            onClick={() => {
              setOpen(false)
              void logout()
            }}
          >
            <LogOut className="h-3.5 w-3.5" />
            Logout
          </button>
        </div>
      ) : null}
    </div>
  )
}
