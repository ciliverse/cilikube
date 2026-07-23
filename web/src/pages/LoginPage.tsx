import { useEffect, useState, type FormEvent } from 'react'
import { Navigate, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { fetchOAuthProviders, register, type OAuthProviderInfo } from '@/api/auth'
import { useAuth } from '@/store/auth'
import { Button, Input } from '@/components/ui'

function GitHubMark({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor" aria-hidden>
      <path d="M12 0C5.37 0 0 5.37 0 12c0 5.3 3.44 9.8 8.21 11.39.6.11.82-.26.82-.58v-2.23c-3.34.73-4.04-1.41-4.04-1.41-.55-1.39-1.33-1.76-1.33-1.76-1.09-.75.08-.73.08-.73 1.2.08 1.84 1.24 1.84 1.24 1.07 1.83 2.81 1.3 3.5 1 .11-.78.42-1.3.76-1.6-2.67-.3-5.47-1.33-5.47-5.93 0-1.31.47-2.38 1.24-3.22-.12-.3-.54-1.52.12-3.18 0 0 1.01-.32 3.3 1.23a11.5 11.5 0 0 1 6 0c2.29-1.55 3.3-1.23 3.3-1.23.66 1.66.24 2.88.12 3.18.77.84 1.24 1.91 1.24 3.22 0 4.61-2.81 5.62-5.49 5.92.43.37.81 1.1.81 2.22v3.29c0 .32.22.7.82.58C20.56 21.8 24 17.3 24 12 24 5.37 18.63 0 12 0z" />
    </svg>
  )
}

export function LoginPage() {
  const { login, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [providers, setProviders] = useState<OAuthProviderInfo[]>([])
  const [allowRegistration, setAllowRegistration] = useState(false)
  const [oauthStatus, setOauthStatus] = useState<'loading' | 'ready' | 'partial' | 'off' | 'error'>(
    'loading',
  )
  const [oauthHint, setOauthHint] = useState('')

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const data = await fetchOAuthProviders()
        if (cancelled) return
        const list = data.providers || []
        setProviders(list)
        setAllowRegistration(Boolean(data.allow_registration))
        const ready = list.filter((p) => p.login_ready && p.auth_url)
        const enabledOnly = list.filter((p) => p.enabled && !p.configured)
        if (ready.length) {
          setOauthStatus('ready')
          setOauthHint('')
        } else if (enabledOnly.length) {
          setOauthStatus('partial')
          setOauthHint(
            'GitHub login is enabled but Client ID is missing. Set credentials in Admin → Settings.',
          )
        } else {
          setOauthStatus('off')
          setOauthHint('')
        }
      } catch {
        if (!cancelled) {
          setOauthStatus('error')
          setOauthHint('Could not load OAuth providers')
        }
      }
    })()
    return () => {
      cancelled = true
    }
  }, [])

  if (isAuthenticated) return <Navigate to="/" replace />

  const readyProviders = providers.filter((p) => p.login_ready && p.auth_url)

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      if (mode === 'register') {
        await register(username.trim(), email.trim(), password)
        await login(username.trim(), password)
      } else {
        await login(username, password)
      }
      navigate('/')
    } catch (err: any) {
      setError(err?.response?.data?.message || err?.message || 'Request failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="relative flex h-full min-h-full items-center justify-center overflow-y-auto px-4 py-8">
      <motion.div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
      >
        <div className="absolute left-[12%] top-[18%] h-64 w-64 rounded-full bg-cyan/10 blur-3xl" />
        <div className="absolute bottom-[12%] right-[10%] h-72 w-72 rounded-full bg-orange/10 blur-3xl" />
      </motion.div>

      <motion.div
        initial={{ opacity: 0, y: 18 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.45, ease: 'easeOut' }}
        className="relative w-full max-w-md"
      >
        <div className="mb-10 text-center">
          <div className="hud-brand text-3xl md:text-4xl">
            CILI<span className="accent">KUBE</span>
          </div>
          <p className="mt-3 text-sm tracking-[0.18em] text-text-dim uppercase">
            Multi-cluster control plane
          </p>
        </div>

        <form className="hud-panel space-y-4 rounded p-6 md:p-8" onSubmit={onSubmit}>
          <div className="hud-label">Authenticate</div>
          <h1 className="font-display text-xl font-bold tracking-[0.14em] text-text">
            {mode === 'login' ? 'SIGN IN' : 'REGISTER'}
          </h1>
          <p className="text-sm text-text-dim">
            {mode === 'login'
              ? 'Enter credentials to open the operator console.'
              : 'Create an account (registration must be enabled by an admin).'}
          </p>

          <label className="block space-y-1.5">
            <span className="hud-label">Username</span>
            <Input
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              autoComplete="username"
              required
            />
          </label>
          {mode === 'register' ? (
            <label className="block space-y-1.5">
              <span className="hud-label">Email</span>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="email"
                required
              />
            </label>
          ) : null}
          <label className="block space-y-1.5">
            <span className="hud-label">Password</span>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
              required
            />
          </label>

          {error ? (
            <div className="rounded border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
              {error}
            </div>
          ) : null}

          <Button type="submit" className="w-full tracking-[0.14em] uppercase" disabled={loading}>
            {loading ? 'Connecting…' : mode === 'login' ? 'Enter' : 'Create account'}
          </Button>

          {allowRegistration ? (
            <button
              type="button"
              className="w-full text-center text-xs tracking-[0.12em] text-cyan uppercase hover:underline"
              onClick={() => {
                setMode((m) => (m === 'login' ? 'register' : 'login'))
                setError('')
              }}
            >
              {mode === 'login' ? 'Need an account? Register' : 'Already registered? Sign in'}
            </button>
          ) : null}

          {oauthStatus === 'partial' || oauthStatus === 'error' ? (
            <div className="rounded border border-warn/40 bg-warn/10 px-3 py-2 text-xs text-warn">
              {oauthHint}
            </div>
          ) : null}

          {readyProviders.length ? (
            <>
              <div className="flex items-center gap-3 text-[11px] tracking-[0.14em] text-text-dim uppercase">
                <span className="h-px flex-1 bg-line" />
                or
                <span className="h-px flex-1 bg-line" />
              </div>
              {readyProviders.map((p) => (
                <Button
                  key={p.name}
                  type="button"
                  variant="outline"
                  className="inline-flex w-full items-center justify-center gap-2 tracking-[0.14em] uppercase"
                  onClick={() => {
                    window.location.href = p.auth_url!
                  }}
                >
                  {p.name === 'github' ? <GitHubMark className="h-4 w-4 shrink-0" /> : null}
                  Continue with {p.display_name || p.name}
                </Button>
              ))}
              {allowRegistration ? (
                <p className="text-center text-[10px] text-text-dim">
                  First-time GitHub sign-in creates a new <span className="text-warn">viewer</span>{' '}
                  account unless the email matches an existing user (auto-link). To keep admin
                  access, sign in as admin and link GitHub from Profile.
                </p>
              ) : (
                <p className="text-center text-[10px] text-text-dim">
                  GitHub accounts must already be linked or provisioned by an admin.
                </p>
              )}
            </>
          ) : null}
        </form>

        <p className="mt-6 text-center text-[11px] tracking-[0.16em] text-text-dim uppercase">
          React · Vite · TanStack Query · Tailwind
        </p>
      </motion.div>
    </div>
  )
}
