import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import {
  fetchOAuthAuthURL,
  fetchOAuthProviders,
  unlinkOAuthAccount,
  type LinkedOAuthProvider,
} from '@/api/auth'
import { apiGet, apiPut } from '@/lib/api'
import { Badge, Button, Card, PageHeader } from '@/components/ui'

type Profile = {
  id?: number
  username?: string
  email?: string
  display_name?: string
  avatar_url?: string
  role?: string
  roles?: string[]
  is_active?: boolean
  email_verified?: boolean
  last_login?: string
  oauth_providers?: LinkedOAuthProvider[]
}

export function ProfilePage() {
  const { t } = useTranslation()
  const [params, setParams] = useSearchParams()
  const [email, setEmail] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [msg, setMsg] = useState('')
  const [err, setErr] = useState('')
  const [busy, setBusy] = useState(false)

  const q = useQuery({
    queryKey: ['profile'],
    queryFn: () => apiGet<Profile>('/api/v1/profile'),
  })

  const providersQ = useQuery({
    queryKey: ['oauth-providers-public'],
    queryFn: () => fetchOAuthProviders(),
  })

  useEffect(() => {
    if (q.data) {
      setEmail(q.data.email || '')
      setDisplayName(q.data.display_name || '')
    }
  }, [q.data])

  useEffect(() => {
    if (params.get('oauth') === 'linked') {
      setMsg('GitHub linked successfully')
      const next = new URLSearchParams(params)
      next.delete('oauth')
      setParams(next, { replace: true })
      void q.refetch()
    }
  }, [params, q, setParams])

  const saveProfile = async () => {
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      await apiPut('/api/v1/profile', {
        email: email.trim(),
        display_name: displayName.trim(),
        avatar_url: q.data?.avatar_url || '',
      })
      setMsg('Profile updated')
      await q.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Update failed')
    } finally {
      setBusy(false)
    }
  }

  const changePassword = async () => {
    if (!oldPassword || !newPassword) {
      setErr('Old and new password are required')
      return
    }
    if (newPassword !== confirmPassword) {
      setErr('Password confirmation does not match')
      return
    }
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      await apiPut('/api/v1/profile/password', {
        old_password: oldPassword,
        new_password: newPassword,
      })
      setOldPassword('')
      setNewPassword('')
      setConfirmPassword('')
      setMsg('Password changed')
    } catch (e: any) {
      setErr(e?.message || 'Password change failed')
    } finally {
      setBusy(false)
    }
  }

  const linkGitHub = async () => {
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      const data = await fetchOAuthAuthURL('github', 'cilikube_link')
      if (!data?.auth_url) throw new Error('Failed to get GitHub authorize URL')
      window.location.href = data.auth_url
    } catch (e: any) {
      setErr(e?.response?.data?.error || e?.message || 'Failed to start GitHub link')
      setBusy(false)
    }
  }

  const unlinkGitHub = async () => {
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      await unlinkOAuthAccount('github')
      setMsg('GitHub unlinked')
      await q.refetch()
    } catch (e: any) {
      setErr(e?.response?.data?.error || e?.message || 'Unlink failed')
    } finally {
      setBusy(false)
    }
  }

  const profile = q.data
  const roles = profile?.roles?.length ? profile.roles : profile?.role ? [profile.role] : []
  const linked = profile?.oauth_providers || []
  const githubLinked = linked.some((p) => p.provider === 'github')
  const githubReady = Boolean(
    providersQ.data?.providers?.some((p) => p.name === 'github' && p.login_ready),
  )

  return (
    <div className="flex w-full min-w-0 flex-col gap-4">
      <PageHeader title={t('profilePage.title')} subtitle={t('profilePage.subtitle')} />
      {err ? (
        <div className="rounded border border-danger/30 bg-danger/10 px-4 py-2 text-sm text-danger">{err}</div>
      ) : null}
      {msg ? (
        <div className="rounded border border-ok/30 bg-ok/10 px-4 py-2 text-sm text-ok">{msg}</div>
      ) : null}

      <Card className="space-y-3 p-5">
        <div className="flex flex-wrap items-center gap-2">
          <span className="font-display text-lg font-bold tracking-[0.12em]">
            {profile?.username || '—'}
          </span>
          {roles.map((r) => (
            <Badge key={r} tone="accent">
              {r}
            </Badge>
          ))}
        </div>
        <p className="text-xs text-text-dim">
          Last login: {profile?.last_login || '—'} · Active:{' '}
          {profile?.is_active === false ? 'no' : 'yes'}
        </p>
        <label className="block space-y-1">
          <span className="hud-label">Email</span>
          <input className="hud-field" value={email} onChange={(e) => setEmail(e.target.value)} />
        </label>
        <label className="block space-y-1">
          <span className="hud-label">Display name</span>
          <input
            className="hud-field"
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
          />
        </label>
        <Button type="button" disabled={busy || q.isLoading} onClick={() => void saveProfile()}>
          Save profile
        </Button>
      </Card>

      <Card className="space-y-3 p-5">
        <h2 className="font-display text-lg font-bold tracking-[0.12em]">LINKED LOGINS</h2>
        <p className="text-xs text-text-dim">
          Link GitHub to <span className="text-text">this</span> account so Continue with GitHub
          signs you in here (with your current roles). If GitHub already created a separate viewer
          user, unlink there first or ask an admin to remove that account.
        </p>
        <div className="flex flex-wrap items-center gap-2">
          <span className="text-sm">GitHub</span>
          <Badge tone={githubLinked ? 'ok' : 'neutral'}>
            {githubLinked ? 'linked' : 'not linked'}
          </Badge>
        </div>
        {githubLinked ? (
          <Button type="button" variant="outline" disabled={busy} onClick={() => void unlinkGitHub()}>
            Unlink GitHub
          </Button>
        ) : (
          <Button
            type="button"
            disabled={busy || !githubReady}
            onClick={() => void linkGitHub()}
          >
            Link GitHub
          </Button>
        )}
        {!githubReady ? (
          <p className="text-xs text-warn">
            GitHub OAuth is not login-ready. Configure Client ID/Secret in Admin → Settings.
          </p>
        ) : null}
      </Card>

      <Card className="space-y-3 p-5">
        <h2 className="font-display text-lg font-bold tracking-[0.12em]">CHANGE PASSWORD</h2>
        <label className="block space-y-1">
          <span className="hud-label">Current password</span>
          <input
            type="password"
            className="hud-field"
            value={oldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
            autoComplete="current-password"
          />
        </label>
        <label className="block space-y-1">
          <span className="hud-label">New password</span>
          <input
            type="password"
            className="hud-field"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            autoComplete="new-password"
          />
        </label>
        <label className="block space-y-1">
          <span className="hud-label">Confirm new password</span>
          <input
            type="password"
            className="hud-field"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            autoComplete="new-password"
          />
        </label>
        <Button type="button" disabled={busy} onClick={() => void changePassword()}>
          Update password
        </Button>
      </Card>
    </div>
  )
}
