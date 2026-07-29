import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiDelete, apiGet, apiPost } from '@/lib/api'
import { useAuth } from '@/store/auth'
import { useCluster } from '@/store/cluster'
import { ALL_NAMESPACES, useNamespace } from '@/store/namespace'
import { Badge, Button, Card, EmptyState, PageHeader } from '@/components/ui'
import { HudTable, HudTableScroll } from '@/components/HudTableScroll'
import { ConfirmDialog } from '@/components/ConfirmDialog'

type HelmRelease = {
  name: string
  namespace: string
  revision?: string
  updated?: string
  status?: string
  chart?: string
  app_version?: string
}

export function HelmPage() {
  const { canEdit, isViewerOnly } = useAuth()
  const { clusterId } = useCluster()
  const { namespace } = useNamespace()
  const [err, setErr] = useState('')
  const [busy, setBusy] = useState(false)
  const [uninstallTarget, setUninstallTarget] = useState<HelmRelease | null>(null)
  const [rollbackTarget, setRollbackTarget] = useState<HelmRelease | null>(null)
  const [form, setForm] = useState({
    name: '',
    namespace: namespace && namespace !== ALL_NAMESPACES ? namespace : 'default',
    chart: '',
    values: '',
  })

  const listQ = useQuery({
    queryKey: ['helm-releases', clusterId, namespace],
    queryFn: () =>
      apiGet<HelmRelease[]>('/api/v1/helm/releases', {
        namespace: namespace && namespace !== ALL_NAMESPACES ? namespace : undefined,
      }),
  })

  const releases = listQ.data || []
  const listErr = (listQ.error as Error | null)?.message || ''
  const combinedErr = err || listErr
  const missingHelmCli = /helm CLI not found/i.test(combinedErr)

  const install = async () => {
    if (!form.name.trim() || !form.chart.trim() || !form.namespace.trim()) {
      setErr('Name, namespace and chart are required')
      return
    }
    setBusy(true)
    setErr('')
    try {
      await apiPost('/api/v1/helm/releases', {
        name: form.name.trim(),
        namespace: form.namespace.trim(),
        chart: form.chart.trim(),
        values: form.values,
        createNamespace: true,
      })
      setForm((f) => ({ ...f, name: '', chart: '', values: '' }))
      await listQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Install failed')
    } finally {
      setBusy(false)
    }
  }

  const uninstall = async () => {
    if (!uninstallTarget) return
    setBusy(true)
    setErr('')
    try {
      await apiDelete(
        `/api/v1/helm/releases/${encodeURIComponent(uninstallTarget.namespace)}/${encodeURIComponent(uninstallTarget.name)}`,
      )
      setUninstallTarget(null)
      await listQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Uninstall failed')
    } finally {
      setBusy(false)
    }
  }

  const rollback = async () => {
    if (!rollbackTarget) return
    setBusy(true)
    setErr('')
    try {
      await apiPost(
        `/api/v1/helm/releases/${encodeURIComponent(rollbackTarget.namespace)}/${encodeURIComponent(rollbackTarget.name)}/rollback`,
        {},
      )
      setRollbackTarget(null)
      await listQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Rollback failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="flex w-full min-w-0 flex-col gap-3">
      <PageHeader title="HELM" subtitle="Release lifecycle via helm CLI on API host" />
      {missingHelmCli ? (
        <div className="rounded border border-orange/40 bg-orange/10 px-4 py-3 text-sm text-text">
          <p className="font-semibold text-orange">Helm CLI not found on API host</p>
          <p className="mt-1 text-text-dim">
            Install the{' '}
            <span className="font-mono text-cyan">helm</span> binary on the machine running CiliKube
            API and ensure it is on <span className="font-mono">PATH</span>, then reload this page.
          </p>
          {combinedErr ? (
            <p className="mt-2 font-mono text-[11px] text-text-dim">{combinedErr}</p>
          ) : null}
        </div>
      ) : null}
      {err && !missingHelmCli ? (
        <div className="rounded border border-danger/30 bg-danger/10 px-4 py-2 text-sm text-danger">{err}</div>
      ) : null}
      {listErr && !missingHelmCli && !err ? (
        <div className="rounded border border-danger/30 bg-danger/10 px-4 py-2 text-sm text-danger">
          {listErr}
        </div>
      ) : null}

      <Card className="overflow-hidden p-0">
        <HudTableScroll maxHeightClass="max-h-[calc(100dvh-14rem)]">
          <HudTable>
            <thead>
              <tr>
                <th>Name</th>
                <th>Namespace</th>
                <th>Chart</th>
                <th>Status</th>
                <th>Revision</th>
                <th>Updated</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {releases.map((r) => (
                <tr key={`${r.namespace}/${r.name}`}>
                  <td className="font-semibold text-cyan">{r.name}</td>
                  <td>{r.namespace}</td>
                  <td className="font-mono text-xs">{r.chart || '-'}</td>
                  <td>
                    <Badge tone={r.status === 'deployed' ? 'ok' : 'warn'}>{r.status || '-'}</Badge>
                  </td>
                  <td>{r.revision || '-'}</td>
                  <td className="text-text-dim">{r.updated || '-'}</td>
                  <td>
                    {canEdit ? (
                      <div className="flex gap-1">
                        <Button
                          variant="ghost"
                          className="px-2 py-1 text-xs"
                          type="button"
                          disabled={busy}
                          onClick={() => setRollbackTarget(r)}
                        >
                          Rollback
                        </Button>
                        <Button
                          variant="danger"
                          className="px-2 py-1 text-xs"
                          type="button"
                          onClick={() => setUninstallTarget(r)}
                        >
                          Uninstall
                        </Button>
                      </div>
                    ) : (
                      <span className="text-xs text-text-dim">read-only</span>
                    )}
                  </td>
                </tr>
              ))}
              {!listQ.isLoading && !releases.length ? (
                <tr>
                  <td colSpan={7}>
                    <EmptyState>
                      {listQ.isError
                        ? listQ.error instanceof Error
                          ? listQ.error.message
                          : 'Failed to list releases'
                        : 'No Helm releases found.'}
                    </EmptyState>
                  </td>
                </tr>
              ) : null}
            </tbody>
          </HudTable>
        </HudTableScroll>
      </Card>

      {canEdit && !isViewerOnly ? (
        <Card className="space-y-3 p-5">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">INSTALL RELEASE</h2>
          <div className="grid gap-3 md:grid-cols-3">
            <label className="block space-y-1">
              <span className="hud-label">Name</span>
              <input
                className="hud-field"
                value={form.name}
                onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              />
            </label>
            <label className="block space-y-1">
              <span className="hud-label">Namespace</span>
              <input
                className="hud-field"
                value={form.namespace}
                onChange={(e) => setForm((f) => ({ ...f, namespace: e.target.value }))}
              />
            </label>
            <label className="block space-y-1">
              <span className="hud-label">Chart</span>
              <input
                className="hud-field font-mono text-xs"
                value={form.chart}
                onChange={(e) => setForm((f) => ({ ...f, chart: e.target.value }))}
                placeholder="bitnami/nginx or ./chart"
              />
            </label>
          </div>
          <label className="block space-y-1">
            <span className="hud-label">Values (YAML)</span>
            <textarea
              className="term-surface min-h-[120px] w-full rounded border border-line px-3 py-2 font-mono text-xs outline-none"
              value={form.values}
              onChange={(e) => setForm((f) => ({ ...f, values: e.target.value }))}
              spellCheck={false}
              placeholder="# optional values.yaml"
            />
          </label>
          <Button type="button" disabled={busy} onClick={() => void install()}>
            Install
          </Button>
        </Card>
      ) : null}

      <ConfirmDialog
        open={Boolean(uninstallTarget)}
        title="UNINSTALL RELEASE"
        confirmText={uninstallTarget?.name}
        confirmLabel="Uninstall"
        busy={busy}
        description={
          <span>
            Uninstall Helm release{' '}
            <span className="font-semibold text-text">
              {uninstallTarget?.namespace}/{uninstallTarget?.name}
            </span>
            ?
          </span>
        }
        onClose={() => setUninstallTarget(null)}
        onConfirm={uninstall}
      />

      <ConfirmDialog
        open={Boolean(rollbackTarget)}
        title="ROLLBACK RELEASE"
        danger={false}
        confirmLabel="Rollback"
        busy={busy}
        description={
          <span>
            Rollback{' '}
            <span className="font-semibold text-text">
              {rollbackTarget?.namespace}/{rollbackTarget?.name}
            </span>{' '}
            to the previous revision?
          </span>
        }
        onClose={() => setRollbackTarget(null)}
        onConfirm={rollback}
      />
    </div>
  )
}
