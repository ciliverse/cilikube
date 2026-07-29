import { useEffect, useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { completeOAuthLogin, linkOAuthAccount, setOAuthNotice } from '@/api/auth'
import { userRoles } from '@/lib/permissions'
import { useAuth } from '@/store/auth'
import { Button } from '@/components/ui'

function isLinkState(state: string | null) {
  return Boolean(state && (state === 'cilikube_link' || state.startsWith('cilikube_link')))
}

export function OAuthCallbackPage() {
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const { applySession, isAuthenticated } = useAuth()
  const [error, setError] = useState('')
  const [status, setStatus] = useState('Exchanging authorization code…')

  useEffect(() => {
    const code = params.get('code')
    const state = params.get('state')
    const oauthError = params.get('error')
    const provider = params.get('provider') || 'github'
    const linking = isLinkState(state)

    if (oauthError) {
      setError(params.get('error_description') || oauthError)
      return
    }
    if (!code) {
      setError('Missing authorization code from OAuth provider')
      return
    }

    // Link flow requires an existing session; login flow must not skip the exchange.
    if (isAuthenticated && !linking) {
      navigate('/ai', { replace: true })
      return
    }
    if (linking && !isAuthenticated) {
      setError('Sign in first, then link GitHub from Profile.')
      return
    }

    let cancelled = false
    ;(async () => {
      try {
        if (linking) {
          setStatus('Linking GitHub to your account…')
          await linkOAuthAccount(provider, code)
          if (cancelled) return
          navigate('/profile?oauth=linked', { replace: true })
          return
        }

        const result = await completeOAuthLogin(provider, code, state || undefined)
        if (cancelled) return
        applySession(result)
        const roles = userRoles(result.user)
        const isViewer =
          roles.includes('viewer') && !roles.includes('admin') && !roles.includes('editor')
        if (result.is_new_user || isViewer) {
          setOAuthNotice({
            kind: result.is_new_user ? 'new_viewer' : 'viewer',
            username: result.user.username,
          })
        }
        navigate('/ai', { replace: true })
      } catch (e: any) {
        if (!cancelled) {
          setError(
            e?.response?.data?.error ||
              e?.response?.data?.message ||
              e?.message ||
              'OAuth login failed',
          )
        }
      }
    })()

    return () => {
      cancelled = true
    }
  }, [applySession, isAuthenticated, navigate, params])

  return (
    <div className="flex min-h-full items-center justify-center px-4 py-12">
      <div className="hud-panel w-full max-w-md space-y-4 rounded p-6 text-center">
        <div className="hud-label">OAuth</div>
        <h1 className="font-display text-xl font-bold tracking-[0.14em]">
          {isLinkState(params.get('state')) ? 'LINKING ACCOUNT' : 'COMPLETING SIGN-IN'}
        </h1>
        {error ? (
          <>
            <p className="rounded border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
              {error}
            </p>
            <Link to={isAuthenticated ? '/profile' : '/login'}>
              <Button type="button" className="w-full tracking-[0.14em] uppercase">
                {isAuthenticated ? 'Back to profile' : 'Back to login'}
              </Button>
            </Link>
          </>
        ) : (
          <p className="text-sm text-text-dim">{status}</p>
        )}
      </div>
    </div>
  )
}
