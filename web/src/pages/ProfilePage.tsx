import { useEffect, useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import {
  fetchOAuthAuthURL,
  fetchOAuthProviders,
  unlinkOAuthAccount,
  type LinkedOAuthProvider,
} from '@/api/auth'
import { apiDelete, apiGet, apiPost, apiPut } from '@/lib/api'
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
  const passkeys = useQuery({
    queryKey: ['passkeys'],
    queryFn: () => apiGet<Array<{ id: number; created_at: string }>>('/api/v1/auth/passkeys'),
  })
  const queryClient = useQueryClient()
  const [totpSecret, setTotpSecret] = useState('')
  const [totpUrl, setTotpUrl] = useState('')
  const [totpCode, setTotpCode] = useState('')

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
      setMsg(t('profilePage.githubLinked'))
      const next = new URLSearchParams(params)
      next.delete('oauth')
      setParams(next, { replace: true })
      void q.refetch()
    }
  }, [params, q, setParams, t])

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
      setMsg(t('profilePage.profileUpdated'))
      await q.refetch()
    } catch (e: any) {
      setErr(e?.message || t('profilePage.updateFailed'))
    } finally {
      setBusy(false)
    }
  }

  const changePassword = async () => {
    if (!oldPassword || !newPassword) {
      setErr(t('profilePage.passwordsRequired'))
      return
    }
    if (newPassword !== confirmPassword) {
      setErr(t('forcePassword.mismatch'))
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
      setMsg(t('profilePage.passwordChanged'))
    } catch (e: any) {
      setErr(e?.message || t('profilePage.passwordChangeFailed'))
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
      if (!data?.auth_url) throw new Error(t('profilePage.githubAuthUrlFailed'))
      window.location.href = data.auth_url
    } catch (e: any) {
      setErr(e?.response?.data?.error || e?.message || t('profilePage.githubLinkFailed'))
      setBusy(false)
    }
  }

  const unlinkGitHub = async () => {
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      await unlinkOAuthAccount('github')
      setMsg(t('profilePage.githubUnlinked'))
      await q.refetch()
    } catch (e: any) {
      setErr(e?.response?.data?.error || e?.message || t('profilePage.unlinkFailed'))
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
          {t('profilePage.accountMeta', {
            lastLogin: profile?.last_login || '—',
            active: profile?.is_active === false ? t('common.no') : t('common.yes'),
          })}
        </p>
        <label className="block space-y-1">
          <span className="hud-label">{t('profilePage.email')}</span>
          <input className="hud-field" value={email} onChange={(e) => setEmail(e.target.value)} />
        </label>
        <label className="block space-y-1">
          <span className="hud-label">{t('profilePage.displayName')}</span>
          <input
            className="hud-field"
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
          />
        </label>
        <Button type="button" disabled={busy || q.isLoading} onClick={() => void saveProfile()}>
          {t('profilePage.saveProfile')}
        </Button>
      </Card>

      <Card className="space-y-3 p-5">
        <h2 className="font-display text-lg font-bold tracking-[0.12em]">
          {t('profilePage.linkedLogins')}
        </h2>
        <p className="text-xs text-text-dim">{t('profilePage.linkedLoginsHint')}</p>
        <div className="flex flex-wrap items-center gap-2">
          <span className="text-sm">GitHub</span>
          <Badge tone={githubLinked ? 'ok' : 'neutral'}>
            {githubLinked ? t('profilePage.linked') : t('profilePage.notLinked')}
          </Badge>
        </div>
        {githubLinked ? (
          <Button type="button" variant="outline" disabled={busy} onClick={() => void unlinkGitHub()}>
            {t('profilePage.unlinkGithub')}
          </Button>
        ) : (
          <Button
            type="button"
            disabled={busy || !githubReady}
            onClick={() => void linkGitHub()}
          >
            {t('profilePage.linkGithub')}
          </Button>
        )}
        {!githubReady ? (
          <p className="text-xs text-warn">{t('profilePage.githubNotReady')}</p>
        ) : null}
      </Card>

      <Card className="space-y-3 p-5">
        <h2 className="font-display text-lg font-bold tracking-[0.12em]">
          {t('profile.changePassword')}
        </h2>
        <label className="block space-y-1">
          <span className="hud-label">{t('forcePassword.oldPassword')}</span>
          <input
            type="password"
            className="hud-field"
            value={oldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
            autoComplete="current-password"
          />
        </label>
        <label className="block space-y-1">
          <span className="hud-label">{t('forcePassword.newPassword')}</span>
          <input
            type="password"
            className="hud-field"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            autoComplete="new-password"
          />
        </label>
        <label className="block space-y-1">
          <span className="hud-label">{t('forcePassword.confirmPassword')}</span>
          <input
            type="password"
            className="hud-field"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            autoComplete="new-password"
          />
        </label>
        <Button type="button" disabled={busy} onClick={() => void changePassword()}>
          {t('profilePage.updatePassword')}
        </Button>
      </Card>
      <Card className="space-y-3 p-4">
        <h2 className="text-sm font-semibold">{t('profilePage.security')}</h2>
        <p className="text-xs text-text-dim">{t('profilePage.securityHint')}</p>
        {totpUrl ? <p className="break-all font-mono text-xs">{totpUrl}</p> : null}
        {totpSecret ? <p className="font-mono text-xs">{totpSecret}</p> : null}
        <div className="flex flex-wrap gap-2">
          <Button
            type="button"
            variant="outline"
            disabled={busy}
            onClick={() => {
              setBusy(true)
              apiPost<{ secret: string; url: string }>('/api/v1/auth/mfa/setup', {})
                .then((data) => {
                  setTotpSecret(data.secret)
                  setTotpUrl(data.url)
                })
                .catch((e: any) => setErr(e?.message || t('login.requestFailed')))
                .finally(() => setBusy(false))
            }}
          >
            {t('profilePage.setupMfa')}
          </Button>
          <input className="hud-field w-32" value={totpCode} onChange={(e) => setTotpCode(e.target.value)} placeholder="123456" />
          <Button
            type="button"
            disabled={busy || !totpCode}
            onClick={() => {
              setBusy(true)
              apiPost('/api/v1/auth/mfa/enable', { code: totpCode })
                .then(() => setMsg(t('profilePage.mfaOn')))
                .catch((e: any) => setErr(e?.message || t('login.requestFailed')))
                .finally(() => setBusy(false))
            }}
          >
            {t('profilePage.enableMfa')}
          </Button>
          <Button
            type="button"
            variant="ghost"
            disabled={busy || !totpCode}
            onClick={() => {
              setBusy(true)
              apiPost('/api/v1/auth/mfa/disable', { code: totpCode })
                .then(() => setMsg(t('profilePage.mfaOff')))
                .catch((e: any) => setErr(e?.message || t('login.requestFailed')))
                .finally(() => setBusy(false))
            }}
          >
            {t('profilePage.disableMfa')}
          </Button>
          <Button
            type="button"
            variant="outline"
            disabled={busy}
            onClick={() => {
              setBusy(true)
              apiPost<{
                session_id: string
                options: { publicKey: PublicKeyCredentialCreationOptionsJSON }
              }>(
                '/api/v1/auth/passkey/register/begin',
                {},
              )
                .then(async (begin) => {
                  const cred = (await navigator.credentials.create({
                    publicKey: PublicKeyCredential.parseCreationOptionsFromJSON(begin.options.publicKey),
                  })) as PublicKeyCredential | null
                  if (!cred) throw new Error(t('login.passkeyCancelled'))
                  await apiPost(`/api/v1/auth/passkey/register/finish?session_id=${encodeURIComponent(begin.session_id)}`, cred.toJSON())
                  setMsg(t('profilePage.passkeyAdded'))
                  void queryClient.invalidateQueries({ queryKey: ['passkeys'] })
                })
                .catch((e: any) => setErr(e?.message || t('login.requestFailed')))
                .finally(() => setBusy(false))
            }}
          >
            {t('profilePage.addPasskey')}
          </Button>
        </div>
        <div className="space-y-1 text-xs">
          <div className="hud-label">{t('profilePage.passkeys')}</div>
          {(passkeys.data || []).length === 0 ? (
            <p className="text-text-dim">{t('profilePage.noPasskeys')}</p>
          ) : (
            (passkeys.data || []).map((item) => (
              <div key={item.id} className="flex items-center justify-between gap-2">
                <span className="font-mono">#{item.id}</span>
                <Button
                  type="button"
                  variant="ghost"
                  className="px-2 py-1 text-xs"
                  disabled={busy}
                  onClick={() => {
                    setBusy(true)
                    apiDelete(`/api/v1/auth/passkeys/${item.id}`)
                      .then(() => queryClient.invalidateQueries({ queryKey: ['passkeys'] }))
                      .catch((e: any) => setErr(e?.message || t('login.requestFailed')))
                      .finally(() => setBusy(false))
                  }}
                >
                  {t('profilePage.removePasskey')}
                </Button>
              </div>
            ))
          )}
        </div>
      </Card>
    </div>
  )
}
