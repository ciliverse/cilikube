import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiGet, apiPut } from '@/lib/api'
import { useAuth } from '@/store/auth'
import { Badge, Button, Card, PageHeader } from '@/components/ui'
import { BUILTIN_THEMES } from '@/theme/themes'
import { FONT_PACKS, getStoredFontId, setFontId } from '@/theme/fonts'
import { switchTheme } from '@/theme/switchTheme'
import { useTheme } from '@/theme/useTheme'

export function AdminSettingsPage() {
  const { isAdmin } = useAuth()
  const { themeId } = useTheme()
  const [msg, setMsg] = useState('')
  const [err, setErr] = useState('')
  const [busy, setBusy] = useState(false)
  const [fontPref, setFontPref] = useState(getStoredFontId)

  const [oauth, setOauth] = useState({
    allow_registration: false,
    auto_link_accounts: false,
    github_enabled: false,
    github_client_id: '',
    github_client_secret: '',
    github_redirect_url: '',
    github_configured: false,
    github_secret_set: false,
  })
  const [security, setSecurity] = useState<any>(null)
  const [prefs, setPrefs] = useState<any>(null)
  const [ai, setAi] = useState({
    enabled: false,
    provider: 'mock',
    base_url: '',
    model: 'mock',
    api_key: '',
    api_key_set: false,
    ready: false,
  })

  const systemQ = useQuery({
    queryKey: ['settings-system'],
    enabled: isAdmin,
    queryFn: () => apiGet<any>('/api/v1/settings/system'),
  })
  const oauthQ = useQuery({
    queryKey: ['settings-oauth'],
    enabled: isAdmin,
    queryFn: () => apiGet<any>('/api/v1/settings/oauth'),
  })
  const securityQ = useQuery({
    queryKey: ['settings-security'],
    enabled: isAdmin,
    queryFn: () => apiGet<any>('/api/v1/settings/security'),
  })
  const prefsQ = useQuery({
    queryKey: ['settings-preferences'],
    enabled: isAdmin,
    queryFn: () => apiGet<any>('/api/v1/settings/preferences'),
  })
  const aiQ = useQuery({
    queryKey: ['settings-ai'],
    enabled: isAdmin,
    queryFn: () => apiGet<any>('/api/v1/settings/ai'),
  })

  useEffect(() => {
    if (oauthQ.data) {
      const github = (oauthQ.data.providers || []).find((p: any) => p.name === 'github')
      const defaultRedirect =
        typeof window !== 'undefined'
          ? `${window.location.origin}/login/oauth/callback`
          : 'http://localhost:8888/login/oauth/callback'
      setOauth((prev) => ({
        ...prev,
        allow_registration: Boolean(oauthQ.data.settings?.allow_registration),
        auto_link_accounts: Boolean(oauthQ.data.settings?.auto_link_accounts),
        github_enabled: Boolean(github?.enabled),
        github_client_id: github?.client_id || '',
        github_redirect_url: github?.redirect_url || defaultRedirect,
        github_configured: Boolean(github?.configured),
        github_secret_set: Boolean(github?.client_secret_set),
        // never echo secret back into the editable field
        github_client_secret: '',
      }))
    }
  }, [oauthQ.data])

  useEffect(() => {
    if (securityQ.data) setSecurity(structuredClone(securityQ.data))
  }, [securityQ.data])

  useEffect(() => {
    if (prefsQ.data) setPrefs(structuredClone(prefsQ.data))
  }, [prefsQ.data])

  useEffect(() => {
    if (aiQ.data) {
      setAi({
        enabled: Boolean(aiQ.data.enabled),
        provider: aiQ.data.provider || 'mock',
        base_url: aiQ.data.base_url || '',
        model: aiQ.data.model || '',
        api_key: '',
        api_key_set: Boolean(aiQ.data.api_key_set),
        ready: Boolean(aiQ.data.ready),
      })
    }
  }, [aiQ.data])

  if (!isAdmin) {
    return (
      <div className="rounded border border-warn/40 bg-warn/10 px-5 py-8 text-sm text-warn">
        Admin privileges required.
      </div>
    )
  }

  const saveOauth = async () => {
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      const defaultRedirect = `${window.location.origin}/login/oauth/callback`
      const body: Record<string, unknown> = {
        allow_registration: oauth.allow_registration,
        auto_link_accounts: oauth.auto_link_accounts,
        github_enabled: oauth.github_enabled,
        github_client_id: oauth.github_client_id,
        github_redirect_url: oauth.github_redirect_url.trim() || defaultRedirect,
      }
      if (oauth.github_client_secret.trim()) {
        body.github_client_secret = oauth.github_client_secret.trim()
      }
      await apiPut('/api/v1/settings/oauth', body)
      setMsg('OAuth settings saved')
      setOauth((o) => ({ ...o, github_client_secret: '' }))
      await oauthQ.refetch()
      await systemQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Save failed')
    } finally {
      setBusy(false)
    }
  }

  const oauthLoginReady = oauth.github_enabled && Boolean(oauth.github_client_id.trim())

  const saveSecurity = async () => {
    if (!security) return
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      await apiPut('/api/v1/settings/security', security)
      setMsg('Security settings saved')
      await securityQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Save failed')
    } finally {
      setBusy(false)
    }
  }

  const savePrefs = async () => {
    if (!prefs) return
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      const nextTheme = prefs.ui_settings?.default_theme || 'tron'
      await apiPut('/api/v1/settings/preferences', prefs)
      // Apply immediately — server prefs alone do not drive the live UI
      switchTheme(nextTheme)
      setFontId(fontPref)
      setMsg('Preferences saved — theme & font applied')
      await prefsQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Save failed')
    } finally {
      setBusy(false)
    }
  }

  const saveAi = async () => {
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      const body: Record<string, unknown> = {
        enabled: ai.enabled,
        provider: ai.provider,
        base_url: ai.base_url,
        model: ai.model,
      }
      if (ai.api_key.trim()) body.api_key = ai.api_key.trim()
      await apiPut('/api/v1/settings/ai', body)
      setMsg('AI settings saved')
      setAi((a) => ({ ...a, api_key: '' }))
      await aiQ.refetch()
      await systemQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Save failed')
    } finally {
      setBusy(false)
    }
  }

  const sys = systemQ.data

  return (
    <div className="flex w-full min-w-0 flex-col gap-4">
      <PageHeader title="SYSTEM SETTINGS" subtitle="OAuth, security and preferences" />
      {err ? (
        <div className="rounded border border-danger/30 bg-danger/10 px-4 py-2 text-sm text-danger">{err}</div>
      ) : null}
      {msg ? (
        <div className="rounded border border-ok/30 bg-ok/10 px-4 py-2 text-sm text-ok">{msg}</div>
      ) : null}

      <Card className="space-y-2 p-5">
        <h2 className="font-display text-lg font-bold tracking-[0.12em]">SYSTEM</h2>
        <div className="flex flex-wrap gap-2 text-sm">
          <Badge tone="accent">version {sys?.version || '—'}</Badge>
          <Badge tone="neutral">{sys?.environment || '—'}</Badge>
          <Badge tone="neutral">{sys?.go_version || '—'}</Badge>
        </div>
        {sys?.features ? (
          <p className="text-xs text-text-dim">
            OAuth {sys.features.oauth_enabled ? 'on' : 'off'} · RBAC{' '}
            {sys.features.rbac_enabled ? 'on' : 'off'} · Audit{' '}
            {sys.features.audit_log_enabled ? 'on' : 'off'} · AI{' '}
            {sys.features.ai_enabled ? 'on' : 'off'}
          </p>
        ) : null}
      </Card>

      <Card className="space-y-3 p-5">
        <div className="flex flex-wrap items-center gap-2">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">AI</h2>
          <Badge tone={ai.ready ? 'ok' : 'warn'}>{ai.ready ? 'ready' : 'not ready'}</Badge>
        </div>
        <p className="text-xs text-text-dim">
          默认 mock 可直接用。切到 openai 时再填 Base URL、Model、API Key（OpenAI-compatible）。
        </p>
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={ai.enabled}
            onChange={(e) => setAi((a) => ({ ...a, enabled: e.target.checked }))}
          />
          Enable AI assistant
        </label>
        <div className="grid gap-3 md:grid-cols-2">
          <label className="block space-y-1">
            <span className="hud-label">Provider</span>
            <select
              className="hud-field"
              value={ai.provider}
              onChange={(e) => setAi((a) => ({ ...a, provider: e.target.value }))}
            >
              <option value="mock">mock (local demo)</option>
              <option value="openai">openai-compatible</option>
            </select>
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Model</span>
            <input
              className="hud-field font-mono text-xs"
              value={ai.model}
              onChange={(e) => setAi((a) => ({ ...a, model: e.target.value }))}
              placeholder="gpt-4o-mini"
            />
          </label>
          <label className="block space-y-1 md:col-span-2">
            <span className="hud-label">Base URL (optional)</span>
            <input
              className="hud-field font-mono text-xs"
              value={ai.base_url}
              onChange={(e) => setAi((a) => ({ ...a, base_url: e.target.value }))}
              placeholder="https://api.openai.com/v1"
            />
          </label>
          <label className="block space-y-1 md:col-span-2">
            <span className="hud-label">
              API Key{ai.api_key_set ? ' (set — leave blank to keep)' : ''}
            </span>
            <input
              type="password"
              className="hud-field font-mono text-xs"
              value={ai.api_key}
              onChange={(e) => setAi((a) => ({ ...a, api_key: e.target.value }))}
              placeholder={ai.api_key_set ? '••••••••' : 'sk-…'}
              autoComplete="new-password"
            />
          </label>
        </div>
        <Button disabled={busy} onClick={() => void saveAi()}>
          Save AI
        </Button>
      </Card>

      <Card className="space-y-3 p-5">
        <div className="flex flex-wrap items-center gap-2">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">OAUTH</h2>
          <Badge tone={oauthLoginReady ? 'ok' : oauth.github_enabled ? 'warn' : 'neutral'}>
            {oauthLoginReady ? 'login ready' : oauth.github_enabled ? 'needs client id' : 'off'}
          </Badge>
        </div>
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={oauth.github_enabled}
            onChange={(e) => setOauth((o) => ({ ...o, github_enabled: e.target.checked }))}
          />
          GitHub login enabled
        </label>
        {oauth.github_enabled && !oauth.github_client_id.trim() ? (
          <div className="rounded border border-warn/40 bg-warn/10 px-3 py-2 text-xs text-warn">
            Enabled without Client ID — the login page will not show a GitHub button until credentials
            are saved.
          </div>
        ) : null}
        <div className="grid gap-3 md:grid-cols-2">
          <label className="block space-y-1">
            <span className="hud-label">GitHub Client ID</span>
            <input
              className="hud-field font-mono text-xs"
              value={oauth.github_client_id}
              onChange={(e) => setOauth((o) => ({ ...o, github_client_id: e.target.value }))}
              placeholder="Iv1.xxxxxxxx"
              autoComplete="off"
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">
              GitHub Client Secret{oauth.github_secret_set ? ' (set — leave blank to keep)' : ''}
            </span>
            <input
              type="password"
              className="hud-field font-mono text-xs"
              value={oauth.github_client_secret}
              onChange={(e) => setOauth((o) => ({ ...o, github_client_secret: e.target.value }))}
              placeholder={oauth.github_secret_set ? '••••••••' : 'ghsecret…'}
              autoComplete="new-password"
            />
          </label>
        </div>
        <label className="block space-y-1">
          <span className="hud-label">Redirect URL (must match GitHub App callback)</span>
          <input
            className="hud-field font-mono text-xs"
            value={oauth.github_redirect_url}
            onChange={(e) => setOauth((o) => ({ ...o, github_redirect_url: e.target.value }))}
            placeholder={`${window.location.origin}/login/oauth/callback`}
          />
        </label>
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={oauth.allow_registration}
            onChange={(e) => setOauth((o) => ({ ...o, allow_registration: e.target.checked }))}
          />
          Allow registration (OAuth first login + username/password register)
        </label>
        <p className="text-xs text-text-dim">
          New OAuth users always get the <span className="text-warn">viewer</span> role. Promote them
          in Admin → Users, or link GitHub to an existing admin from Profile (emails must match for
          auto-link).
        </p>
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={oauth.auto_link_accounts}
            onChange={(e) => setOauth((o) => ({ ...o, auto_link_accounts: e.target.checked }))}
          />
          Auto-link accounts by email
        </label>
        <Button type="button" disabled={busy} onClick={() => void saveOauth()}>
          Save OAuth
        </Button>
      </Card>

      {security ? (
        <Card className="space-y-3 p-5">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">SECURITY</h2>
          <div className="grid gap-3 md:grid-cols-3">
            <label className="block space-y-1">
              <span className="hud-label">Min password length</span>
              <input
                type="number"
                className="hud-field"
                value={security.password_policy?.min_length ?? 8}
                onChange={(e) =>
                  setSecurity((s: any) => ({
                    ...s,
                    password_policy: {
                      ...s.password_policy,
                      min_length: Number(e.target.value) || 0,
                    },
                  }))
                }
              />
            </label>
            <label className="block space-y-1">
              <span className="hud-label">Session timeout (s)</span>
              <input
                type="number"
                className="hud-field"
                value={security.session_settings?.session_timeout ?? 0}
                onChange={(e) =>
                  setSecurity((s: any) => ({
                    ...s,
                    session_settings: {
                      ...s.session_settings,
                      session_timeout: Number(e.target.value) || 0,
                    },
                  }))
                }
              />
            </label>
            <label className="block space-y-1">
              <span className="hud-label">Audit retention (days)</span>
              <input
                type="number"
                className="hud-field"
                value={security.audit_settings?.retention_days ?? 30}
                onChange={(e) =>
                  setSecurity((s: any) => ({
                    ...s,
                    audit_settings: {
                      ...s.audit_settings,
                      retention_days: Number(e.target.value) || 0,
                    },
                  }))
                }
              />
            </label>
          </div>
          <div className="flex flex-wrap gap-4 text-sm">
            {(
              [
                ['require_uppercase', 'Uppercase'],
                ['require_lowercase', 'Lowercase'],
                ['require_numbers', 'Numbers'],
                ['require_symbols', 'Symbols'],
              ] as const
            ).map(([key, label]) => (
              <label key={key} className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={Boolean(security.password_policy?.[key])}
                  onChange={(e) =>
                    setSecurity((s: any) => ({
                      ...s,
                      password_policy: { ...s.password_policy, [key]: e.target.checked },
                    }))
                  }
                />
                {label}
              </label>
            ))}
          </div>
          <div className="flex flex-wrap gap-4 text-sm">
            {(
              [
                ['log_login_attempts', 'Log logins'],
                ['log_api_calls', 'Log API calls'],
                ['log_admin_actions', 'Log admin actions'],
              ] as const
            ).map(([key, label]) => (
              <label key={key} className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={Boolean(security.audit_settings?.[key])}
                  onChange={(e) =>
                    setSecurity((s: any) => ({
                      ...s,
                      audit_settings: { ...s.audit_settings, [key]: e.target.checked },
                    }))
                  }
                />
                {label}
              </label>
            ))}
          </div>
          <Button type="button" disabled={busy} onClick={() => void saveSecurity()}>
            Save security
          </Button>
        </Card>
      ) : null}

      {prefs ? (
        <Card className="space-y-3 p-5">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">PREFERENCES</h2>
          <div className="grid gap-3 md:grid-cols-3">
            <label className="block space-y-1">
              <span className="hud-label">Theme</span>
              <select
                className="hud-field"
                value={prefs.ui_settings?.default_theme || themeId || 'tron'}
                onChange={(e) => {
                  const id = e.target.value
                  setPrefs((p: any) => ({
                    ...p,
                    ui_settings: { ...p.ui_settings, default_theme: id },
                  }))
                  switchTheme(id)
                }}
              >
                {BUILTIN_THEMES.map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.name} ({t.mode})
                  </option>
                ))}
              </select>
              <p className="text-[10px] text-text-dim">
                Applies now; Save also persists as the server default.
              </p>
            </label>
            <label className="block space-y-1">
              <span className="hud-label">Font</span>
              <select
                className="hud-field"
                value={fontPref}
                onChange={(e) => {
                  setFontPref(e.target.value)
                  setFontId(e.target.value)
                }}
              >
                {FONT_PACKS.map((f) => (
                  <option key={f.id} value={f.id}>
                    {f.name}
                  </option>
                ))}
              </select>
              <p className="text-[10px] text-text-dim">
                Default: Latin Maple (~75KB) + system CJK. Maple Mono CN loads subset fonts from this
                site on demand (no third-party CDN).
              </p>
            </label>
            <label className="block space-y-1">
              <span className="hud-label">Default language</span>
              <input
                className="hud-field"
                value={prefs.ui_settings?.default_language || ''}
                onChange={(e) =>
                  setPrefs((p: any) => ({
                    ...p,
                    ui_settings: { ...p.ui_settings, default_language: e.target.value },
                  }))
                }
              />
            </label>
            <label className="block space-y-1">
              <span className="hud-label">Items per page</span>
              <input
                type="number"
                className="hud-field"
                value={prefs.ui_settings?.items_per_page ?? 20}
                onChange={(e) =>
                  setPrefs((p: any) => ({
                    ...p,
                    ui_settings: {
                      ...p.ui_settings,
                      items_per_page: Number(e.target.value) || 0,
                    },
                  }))
                }
              />
            </label>
          </div>
          <div className="flex flex-wrap gap-4 text-sm">
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={Boolean(prefs.ui_settings?.auto_refresh)}
                onChange={(e) =>
                  setPrefs((p: any) => ({
                    ...p,
                    ui_settings: { ...p.ui_settings, auto_refresh: e.target.checked },
                  }))
                }
              />
              Auto refresh
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={Boolean(prefs.feature_flags?.advanced_metrics)}
                onChange={(e) =>
                  setPrefs((p: any) => ({
                    ...p,
                    feature_flags: { ...p.feature_flags, advanced_metrics: e.target.checked },
                  }))
                }
              />
              Advanced metrics
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={Boolean(prefs.feature_flags?.beta_features)}
                onChange={(e) =>
                  setPrefs((p: any) => ({
                    ...p,
                    feature_flags: { ...p.feature_flags, beta_features: e.target.checked },
                  }))
                }
              />
              Beta features
            </label>
          </div>
          <Button type="button" disabled={busy} onClick={() => void savePrefs()}>
            Save preferences
          </Button>
        </Card>
      ) : null}
    </div>
  )
}
