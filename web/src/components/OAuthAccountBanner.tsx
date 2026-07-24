import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { clearOAuthNotice, peekOAuthNotice, type OAuthNotice } from '@/api/auth'
import { useAuth } from '@/store/auth'

function dismissKey(userId: number) {
  return `cilikube_viewer_banner_dismissed_${userId}`
}

export function OAuthAccountBanner() {
  const { user, isViewerOnly, logout } = useAuth()
  const [notice, setNotice] = useState<OAuthNotice | null>(null)
  const [dismissed, setDismissed] = useState(false)
  const [expanded, setExpanded] = useState(false)

  useEffect(() => {
    if (!user?.id) return
    setNotice(peekOAuthNotice())
    setDismissed(localStorage.getItem(dismissKey(user.id)) === '1')
  }, [user?.id])

  if (!user || !isViewerOnly || dismissed) return null

  const title =
    notice?.kind === 'new_viewer'
      ? `New account: ${user.username} (viewer)`
      : `Signed in as ${user.username} (viewer)`

  return (
    <div className="mb-2.5 rounded border border-warn/40 bg-warn/10 px-3 py-2.5 text-sm text-text sm:mb-4 sm:px-4 sm:py-3">
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0 space-y-1">
          <p className="text-[13px] font-semibold tracking-wide text-warn sm:text-sm">{title}</p>
          <p
            className={cnClamp(
              'text-[11px] text-text-dim sm:text-xs',
              expanded ? '' : 'line-clamp-2 sm:line-clamp-none',
            )}
          >
            Admin data stays on the admin account. Options: sign out and log in as admin; ask an
            admin to promote this user; or link GitHub from Profile while logged in as admin.
          </p>
          <button
            type="button"
            className="text-[11px] text-cyan hover:underline sm:hidden"
            onClick={() => setExpanded((v) => !v)}
          >
            {expanded ? 'Show less' : 'More info'}
          </button>
          <div className="flex flex-wrap gap-x-3 gap-y-1 pt-0.5 text-[11px] sm:text-xs">
            <Link to="/profile" className="text-cyan hover:underline">
              Open Profile
            </Link>
            <button
              type="button"
              className="text-cyan hover:underline"
              onClick={() => void logout()}
            >
              Sign out
            </button>
          </div>
        </div>
        <button
          type="button"
          className="shrink-0 text-[10px] tracking-[0.12em] text-text-dim uppercase hover:text-text sm:text-xs"
          onClick={() => {
            clearOAuthNotice()
            localStorage.setItem(dismissKey(user.id), '1')
            setDismissed(true)
          }}
        >
          Dismiss
        </button>
      </div>
    </div>
  )
}

function cnClamp(...parts: Array<string | false | undefined>) {
  return parts.filter(Boolean).join(' ')
}
