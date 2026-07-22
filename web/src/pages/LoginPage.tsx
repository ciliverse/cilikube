import { useState, type FormEvent } from 'react'
import { Navigate, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
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
            SIGN IN
          </h1>
          <p className="text-sm text-text-dim">
            Enter credentials to open the operator console.
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
          <label className="block space-y-1.5">
            <span className="hud-label">Password</span>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              required
            />
          </label>

          {error ? (
            <div className="rounded border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger">
              {error}
            </div>
          ) : null}

          <Button type="submit" className="w-full tracking-[0.14em] uppercase" disabled={loading}>
            {loading ? 'Connecting…' : 'Enter'}
          </Button>
        </form>

        <p className="mt-6 text-center text-[11px] tracking-[0.16em] text-text-dim uppercase">
          React · Vite · TanStack Query · Tailwind
        </p>
      </motion.div>
    </div>
  )
}
