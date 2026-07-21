import { useState, type FormEvent } from 'react'
import { Navigate, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Hexagon, Sparkles } from 'lucide-react'
import { useAuth } from '@/store/auth'
import { Button, Input } from '@/components/ui'

export function LoginPage() {
  const { login, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  if (isAuthenticated) return <Navigate to="/" replace />

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      await login(username, password)
      navigate('/')
    } catch (err: any) {
      setError(err?.response?.data?.message || err?.message || 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="relative flex min-h-screen items-center justify-center px-4">
      <motion.div
        aria-hidden
        className="pointer-events-none absolute inset-0 overflow-hidden"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
      >
        <motion.div
          className="absolute -left-20 top-10 h-72 w-72 rounded-full bg-accent/20 blur-3xl"
          animate={{ x: [0, 40, 0], y: [0, 30, 0] }}
          transition={{ duration: 10, repeat: Infinity, ease: 'easeInOut' }}
        />
        <motion.div
          className="absolute bottom-10 right-0 h-80 w-80 rounded-full bg-signal/15 blur-3xl"
          animate={{ x: [0, -30, 0], y: [0, -20, 0] }}
          transition={{ duration: 12, repeat: Infinity, ease: 'easeInOut' }}
        />
      </motion.div>

      <motion.div
        initial={{ opacity: 0, y: 24, scale: 0.98 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.5, ease: 'easeOut' }}
        className="relative grid w-full max-w-5xl overflow-hidden rounded-[28px] border border-white/70 bg-white/70 shadow-[0_40px_120px_-50px_rgba(16,24,40,0.55)] backdrop-blur-xl md:grid-cols-[1.15fr_0.85fr]"
      >
        <section className="relative hidden overflow-hidden bg-ink px-10 py-12 text-white md:block">
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(20,184,166,0.35),transparent_45%),radial-gradient(circle_at_80%_70%,rgba(234,88,12,0.25),transparent_40%)]" />
          <div className="relative z-10 flex h-full flex-col justify-between">
            <div className="flex items-center gap-3">
              <div className="grid h-12 w-12 place-items-center rounded-2xl bg-accent">
                <Hexagon className="h-7 w-7" />
              </div>
              <div className="font-display text-3xl font-extrabold tracking-tight">CiliKube</div>
            </div>
            <div>
              <h2 className="font-display text-4xl font-extrabold leading-tight tracking-tight">
                Multi-cluster
                <br />
                control, refined.
              </h2>
              <p className="mt-4 max-w-md text-sm leading-relaxed text-white/70">
                Operate nodes, workloads, events, and live metrics from one atmospheric control
                surface — built for Kubernetes operators who care about clarity and pace.
              </p>
            </div>
            <div className="inline-flex items-center gap-2 text-sm text-accent-bright">
              <Sparkles className="h-4 w-4" />
              React · Vite · TanStack Query · Framer Motion
            </div>
          </div>
        </section>

        <section className="px-8 py-10 md:px-10">
          <div className="mb-8 md:hidden">
            <div className="font-display text-3xl font-extrabold">CiliKube</div>
          </div>
          <h1 className="font-display text-3xl font-extrabold tracking-tight">Sign in</h1>
          <p className="mt-2 text-sm text-ink-soft">Enter your credentials to open the control plane.</p>

          <form className="mt-8 space-y-4" onSubmit={onSubmit}>
            <label className="block space-y-1.5">
              <span className="text-xs font-semibold uppercase tracking-wide text-ink-soft">
                Username
              </span>
              <Input
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                autoComplete="username"
                required
              />
            </label>
            <label className="block space-y-1.5">
              <span className="text-xs font-semibold uppercase tracking-wide text-ink-soft">
                Password
              </span>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete="current-password"
                required
              />
            </label>

            {error ? (
              <div className="rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm text-danger">
                {error}
              </div>
            ) : null}

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? 'Signing in…' : 'Enter control plane'}
            </Button>
          </form>
        </section>
      </motion.div>
    </div>
  )
}
