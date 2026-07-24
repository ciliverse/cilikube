import { useEffect, useState, type FormEvent } from 'react'
import { Navigate, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { fetchOAuthProviders, register, type OAuthProviderInfo } from '@/api/auth'
import { fetchShowcaseInfo, type ShowcaseInfo } from '@/api/showcase'
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
  const [showcase, setShowcase] = useState<ShowcaseInfo | null>(null)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const [oauthData, showcaseData] = await Promise.all([
          fetchOAuthProviders(),
          fetchShowcaseInfo(),
        ])
        if (cancelled) return
        setShowcase(showcaseData.showcase ? showcaseData : null)
        const list = oauthData.providers || []
        setProviders(list)
        // Never offer self-registration on the public showcase.
        setAllowRegistration(Boolean(oauthData.allow_registration) && !showcaseData.showcase)
        const ready = list.filter((p) => p.login_ready && p.auth_url)
        const enabledOnly = list.filter((p) => p.enabled && !p.configured)
        if (showcaseData.showcase) {
          setOauthStatus('off')
          setOauthHint('')
        } else if (ready.length) {
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
    <div className="relative flex min-h-dvh w-full items-start justify-center overflow-y-auto px-4 pb-[max(1.5rem,env(safe-area-inset-bottom))] pt-[max(1.5rem,env(safe-area-inset-top))] sm:items-center sm:py-8">
      <motion.div
        aria-hidden
        className="pointer-events-none absolute inset-0 hidden sm:block"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
      >
        <div className="absolute left-[12%] top-[18%] h-64 w-64 rounded-full bg-cyan/10 blur-3xl" />
        <div className="absolute bottom-[12%] right-[10%] h-72 w-72 rounded-full bg-orange/10 blur-3xl" />
      </motion.div>

      <motion.div
        initial={{ opacity: 0, y: 14 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4, ease: 'easeOut' }}
        className="relative w-full max-w-md"
      >
        <div className="mb-6 text-center sm:mb-10">
          <div className="hud-brand text-2xl sm:text-3xl md:text-4xl">
            CILI<span className="accent">KUBE</span>
          </div>
          <p className="mt-2 text-xs tracking-[0.18em] text-text-dim uppercase sm:mt-3 sm:text-sm">
            Multi-cluster control plane
          </p>
        </div>

        <form
          className="hud-panel space-y-3.5 rounded p-5 sm:space-y-4 sm:p-6 md:p-8"
          onSubmit={onSubmit}
        >
          <div className="hud-label">Authenticate</div>
          <h1 className="font-display text-lg font-bold tracking-[0.14em] text-text sm:text-xl">
            {mode === 'login' ? 'SIGN IN' : 'REGISTER'}
          </h1>
          <p className="text-[13px] text-text-dim sm:text-sm">
            {mode === 'login'
              ? 'Enter credentials to open the operator console.'
              : 'Create an account (registration must be enabled by an admin).'}
          </p>

          {showcase?.accounts?.length ? (
            <div className="space-y-2 rounded border border-cyan/30 bg-cyan/5 px-3 py-3 text-sm">
              <div className="hud-label text-cyan">Accounts</div>
              {showcase.accounts.map((a) => (
                <button
                  key={a.username}
                  type="button"
                  className="flex w-full flex-col gap-0.5 rounded border border-line/60 bg-bg/40 px-3 py-2 text-left transition hover:border-cyan/50"
                  onClick={() => {
                    setUsername(a.username)
                    setPassword(a.password)
                    setMode('login')
                    setError('')
                  }}
                >
                  <span className="font-mono text-xs tracking-wide text-text">
                    {a.username}
                    <span className="mx-2 text-text-dim">/</span>
                    {a.password}
                  </span>
                  <span className="text-[11px] text-text-dim">role: {a.role}</span>
                </button>
              ))}
            </div>
          ) : null}

          <label className="block space-y-1.5">
            <span className="hud-label">Username</span>
            <Input
              className="h-11 text-base sm:h-auto sm:text-sm"
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
                className="h-11 text-base sm:h-auto sm:text-sm"
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
              className="h-11 text-base sm:h-auto sm:text-sm"
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

          <Button
            type="submit"
            className="h-11 w-full tracking-[0.14em] uppercase"
            disabled={loading}
          >
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
                  className="inline-flex h-11 w-full items-center justify-center gap-2 tracking-[0.14em] uppercase"
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

        <p className="mt-5 text-center text-[10px] tracking-[0.14em] text-text-dim uppercase sm:mt-6 sm:text-[11px] sm:tracking-[0.16em]">
          React · Vite · TanStack Query · Tailwind
        </p>
      </motion.div>
    </div>
  )
}
