import { useState, type FormEvent } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '@/store/auth'
import { apiPost } from '@/lib/api'
import { login as loginApi } from '@/api/auth'
import { Button, Input, PageHeader } from '@/components/ui'

export function ForcePasswordChangePage({
  defaultOldPassword = '',
}: {
  defaultOldPassword?: string
}) {
  const { t } = useTranslation()
  const { user, applySession, logout } = useAuth()
  const [oldPassword, setOldPassword] = useState(defaultOldPassword)
  const [newPassword, setNewPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    if (newPassword.length < 8) {
      setError(t('forcePassword.tooShort'))
      return
    }
    if (newPassword === '12345678') {
      setError(t('forcePassword.notDefault'))
      return
    }
    if (newPassword !== confirm) {
      setError(t('forcePassword.mismatch'))
      return
    }
    setBusy(true)
    try {
      await apiPost('/api/v1/auth/change-password', {
        old_password: oldPassword,
        new_password: newPassword,
      })
      const username = user?.username || ''
      await logout()
      const result = await loginApi(username, newPassword)
      applySession(result)
    } catch (err: any) {
      setError(err?.response?.data?.message || err?.message || t('common.failed'))
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-bg px-4">
      <div className="hud-panel w-full max-w-md space-y-4 p-6">
        <PageHeader title={t('forcePassword.title')} subtitle={t('forcePassword.subtitle')} />
        {error ? (
          <div className="rounded border border-danger/30 bg-danger/10 px-3 py-2 text-sm text-danger">
            {error}
          </div>
        ) : null}
        <form className="space-y-3" onSubmit={(e) => void onSubmit(e)}>
          <label className="block space-y-1 text-sm">
            <span className="hud-label">{t('forcePassword.oldPassword')}</span>
            <Input
              type="password"
              autoComplete="current-password"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              required
            />
          </label>
          <label className="block space-y-1 text-sm">
            <span className="hud-label">{t('forcePassword.newPassword')}</span>
            <Input
              type="password"
              autoComplete="new-password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              required
              minLength={8}
            />
          </label>
          <label className="block space-y-1 text-sm">
            <span className="hud-label">{t('forcePassword.confirmPassword')}</span>
            <Input
              type="password"
              autoComplete="new-password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              required
              minLength={8}
            />
          </label>
          <Button type="submit" disabled={busy} className="w-full">
            {busy ? t('common.loading') : t('forcePassword.submit')}
          </Button>
        </form>
      </div>
    </div>
  )
}
