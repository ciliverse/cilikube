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
        className="flex items-center gap-2 rounded border border-transparent px-2 py-1 text-xs text-text-dim transition hover:border-line hover:bg-mist hover:text-text"
        aria-haspopup="menu"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
      >
        <span className="text-text">{user?.username || 'operator'}</span>
        <span className="text-line">/</span>
        <span className={roleTone}>{primaryRole}</span>
        {isViewerOnly ? (
          <span className="rounded border border-line px-1.5 py-0.5 text-[10px] tracking-wider uppercase">
            read-only
          </span>
        ) : null}
        <ChevronDown className={cn('h-3.5 w-3.5 transition', open && 'rotate-180 text-cyan')} />
      </button>

      {open ? (
        <div
          role="menu"
          className="absolute right-0 z-[210] mt-1 min-w-[220px] overflow-hidden rounded border border-line bg-panel-solid py-1 shadow-[0_8px_28px_rgba(0,0,0,0.45)]"
        >
          <Link
            role="menuitem"
            to="/profile"
            className="flex items-center gap-2 px-3 py-2 text-[12px] text-text hover:bg-mist"
            onClick={() => setOpen(false)}
          >
            <UserRound className="h-3.5 w-3.5 text-cyan" />
            Profile
          </Link>
          {isAdmin ? (
            <Link
              role="menuitem"
              to="/admin/settings"
              className="flex items-center gap-2 px-3 py-2 text-[12px] text-text hover:bg-mist"
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
                  'rounded px-2 py-1.5 text-left text-[12px] transition',
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
            className="flex w-full items-center gap-2 px-3 py-2 text-left text-[12px] text-danger hover:bg-mist"
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
