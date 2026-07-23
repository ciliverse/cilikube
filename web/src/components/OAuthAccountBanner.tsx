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

  useEffect(() => {
    if (!user?.id) return
    setNotice(peekOAuthNotice())
    setDismissed(localStorage.getItem(dismissKey(user.id)) === '1')
  }, [user?.id])

  if (!user || !isViewerOnly || dismissed) return null

  const title =
    notice?.kind === 'new_viewer'
      ? `New account created: ${user.username} (viewer)`
      : `Signed in as ${user.username} (viewer)`

  return (
    <div className="mb-4 rounded border border-warn/40 bg-warn/10 px-4 py-3 text-sm text-text">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="min-w-0 space-y-1">
          <p className="font-semibold tracking-wide text-warn">{title}</p>
          <p className="text-xs text-text-dim">
            Admin data and clusters stay on the original admin account. First-time GitHub sign-in
            creates a separate viewer user when emails do not match. Options: sign out and log in as
            admin; ask an admin to promote this user in Admin → Users; or link GitHub from Profile
            while logged in as admin (unlink here first if this GitHub is already attached).
          </p>
          <div className="flex flex-wrap gap-3 pt-1 text-xs">
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
          className="shrink-0 text-xs tracking-[0.12em] text-text-dim uppercase hover:text-text"
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
