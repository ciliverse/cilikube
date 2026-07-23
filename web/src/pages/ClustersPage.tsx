import { useMemo, useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { listClusters, updateCluster, type ClusterItem } from '@/api/cluster'
import { apiDelete, apiGet, apiPost } from '@/lib/api'
import { useAuth } from '@/store/auth'
import { useCluster } from '@/store/cluster'
import { Badge, Button, Card, EmptyState, Modal, PageHeader } from '@/components/ui'
import { HudTable, HudTableScroll } from '@/components/HudTableScroll'
import { ConfirmDialog } from '@/components/ConfirmDialog'

type LocalContext = {
  name: string
  cluster: string
  server: string
  user: string
  conflict_name?: boolean
  conflict_server?: boolean
  existing_name?: string
}

type LocalContextsResp = {
  path: string
  contexts: LocalContext[]
}

export function ClustersPage() {
  const { isAdmin, canMutate } = useAuth()
  const { clusterId, setClusterId } = useCluster()
  const queryClient = useQueryClient()
  const canWrite = isAdmin || canMutate('clusters')
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [kubeconfig, setKubeconfig] = useState('')
  const [err, setErr] = useState('')
  const [msg, setMsg] = useState('')
  const [busy, setBusy] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<ClusterItem | null>(null)
  const [editTarget, setEditTarget] = useState<ClusterItem | null>(null)
  const [editName, setEditName] = useState('')
  const [editDescription, setEditDescription] = useState('')
  const [editKubeconfig, setEditKubeconfig] = useState('')
  const [selectedContexts, setSelectedContexts] = useState<string[]>([])

  const q = useQuery({
    queryKey: ['clusters-mgmt'],
    queryFn: listClusters,
  })

  const localQ = useQuery({
    queryKey: ['clusters-local-contexts'],
    queryFn: () => apiGet<LocalContextsResp>('/api/v1/clusters/local-contexts'),
    enabled: canWrite,
    staleTime: 15_000,
  })

  const clusters = q.data || []
  const localContexts = useMemo(() => localQ.data?.contexts || [], [localQ.data])

  const create = async () => {
    if (!name.trim() || !kubeconfig.trim()) {
      setErr('Name and kubeconfig are required')
      return
    }
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      await apiPost('/api/v1/clusters', {
        name: name.trim(),
        description: description.trim() || undefined,
        kubeconfigData: btoa(unescape(encodeURIComponent(kubeconfig))),
      })
      setName('')
      setDescription('')
      setKubeconfig('')
      setMsg('Cluster created')
      await q.refetch()
      void queryClient.invalidateQueries({ queryKey: ['clusters'] })
      void localQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Create failed')
    } finally {
      setBusy(false)
    }
  }

  const toggleContext = (ctx: string) => {
    setSelectedContexts((prev) =>
      prev.includes(ctx) ? prev.filter((x) => x !== ctx) : [...prev, ctx],
    )
  }

  const importLocal = async () => {
    if (!selectedContexts.length) {
      setErr('Select at least one local context')
      return
    }
    setBusy(true)
    setErr('')
    setMsg('')
    try {
      const res = await apiPost<{
        imported: string[]
        skipped: Record<string, string>
        failed: Record<string, string>
      }>('/api/v1/clusters/import-local', {
        contexts: selectedContexts,
        skip_conflicts: true,
      })
      const imported = res?.imported || []
      const skipped = Object.keys(res?.skipped || {})
      const failed = Object.keys(res?.failed || {})
      setMsg(
        `Imported ${imported.length}` +
          (skipped.length ? `, skipped ${skipped.length}` : '') +
          (failed.length ? `, failed ${failed.length}` : ''),
      )
      if (failed.length) {
        setErr(Object.entries(res.failed).map(([k, v]) => `${k}: ${v}`).join('; '))
      }
      setSelectedContexts([])
      await q.refetch()
      void queryClient.invalidateQueries({ queryKey: ['clusters'] })
      void localQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Import failed')
    } finally {
      setBusy(false)
    }
  }

  const setActive = async (id: string) => {
    setBusy(true)
    setErr('')
    try {
      // Shared switcher: server active + topbar selection + query invalidate
      setClusterId(id)
      await q.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Set active failed')
    } finally {
      setBusy(false)
    }
  }

  const remove = async () => {
    if (!deleteTarget?.id) return
    setBusy(true)
    try {
      await apiDelete(`/api/v1/clusters/${deleteTarget.id}`)
      setDeleteTarget(null)
      await q.refetch()
      void queryClient.invalidateQueries({ queryKey: ['clusters'] })
    } catch (e: any) {
      setErr(e?.message || 'Delete failed')
    } finally {
      setBusy(false)
    }
  }

  const openEdit = (c: ClusterItem) => {
    setEditTarget(c)
    setEditName(c.name || '')
    setEditDescription(String(c.description || ''))
    setEditKubeconfig('')
    setErr('')
  }

  const saveEdit = async () => {
    if (!editTarget?.id) return
    if (!editName.trim()) {
      setErr('Name is required')
      return
    }
    setBusy(true)
    setErr('')
    try {
      const body: Parameters<typeof updateCluster>[1] = {
        name: editName.trim(),
        description: editDescription.trim(),
      }
      if (editKubeconfig.trim()) {
        body.kubeconfigData = btoa(unescape(encodeURIComponent(editKubeconfig)))
      }
      await updateCluster(String(editTarget.id), body)
      setEditTarget(null)
      await q.refetch()
      void queryClient.invalidateQueries({ queryKey: ['clusters'] })
    } catch (e: any) {
      setErr(e?.message || 'Update failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="flex w-full min-w-0 flex-col gap-4">
      <PageHeader
        title="CLUSTERS"
        subtitle="Paste kubeconfig or import contexts from the API host kubeconfig"
      />
      {err ? (
        <div className="rounded border border-danger/30 bg-danger/10 px-4 py-2 text-sm text-danger">{err}</div>
      ) : null}
      {msg ? (
        <div className="rounded border border-ok/30 bg-ok/10 px-4 py-2 text-sm text-ok">{msg}</div>
      ) : null}

      <Card className="overflow-hidden">
        <HudTableScroll>
          <HudTable>
            <thead>
              <tr>
                <th>Name</th>
                <th>Status</th>
                <th>Version</th>
                <th>Source</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {clusters.map((c) => {
                const id = String(c.id || c.name)
                const selected = id === clusterId || c.name === clusterId
                return (
                <tr key={id} className={selected ? 'bg-mist' : undefined}>
                  <td className="font-semibold text-cyan">
                    {c.name}
                    {selected ? (
                      <span className="ml-2 text-[10px] tracking-wider text-ok uppercase">
                        current
                      </span>
                    ) : null}
                  </td>
                  <td>
                    <Badge tone={selected || c.status === 'Available' || c.status === 'Active' ? 'ok' : 'neutral'}>
                      {selected ? 'Active' : c.status || '-'}
                    </Badge>
                  </td>
                  <td>{c.version || '-'}</td>
                  <td className="text-text-dim">{String(c.source || '-')}</td>
                  <td>
                    <div className="flex gap-1">
                      {canWrite ? (
                        <>
                          <Button
                            variant="outline"
                            className="px-2 py-1 text-xs"
                            type="button"
                            disabled={busy || selected}
                            onClick={() => void setActive(id)}
                          >
                            {selected ? 'Active' : 'Set active'}
                          </Button>
                          <Button
                            variant="ghost"
                            className="px-2 py-1 text-xs"
                            type="button"
                            onClick={() => openEdit(c)}
                          >
                            Edit
                          </Button>
                          <Button
                            variant="danger"
                            className="px-2 py-1 text-xs"
                            type="button"
                            onClick={() => setDeleteTarget(c)}
                          >
                            Delete
                          </Button>
                        </>
                      ) : (
                        <span className="text-xs text-text-dim">read-only</span>
                      )}
                    </div>
                  </td>
                </tr>
              )})}
              {!q.isLoading && !clusters.length ? (
                <tr>
                  <td colSpan={5}>
                    <EmptyState>No clusters configured.</EmptyState>
                  </td>
                </tr>
              ) : null}
            </tbody>
          </HudTable>
        </HudTableScroll>
      </Card>

      {canWrite ? (
        <Card className="space-y-3 p-5">
          <div className="flex flex-wrap items-end justify-between gap-2">
            <div>
              <h2 className="font-display text-lg font-bold tracking-[0.12em]">
                IMPORT LOCAL KUBECONFIG
              </h2>
              <p className="mt-1 text-[11px] text-text-dim">
                Source: {localQ.data?.path || '…'} — selected contexts are copied into the DB
                (not hot-mounted).
              </p>
            </div>
            <Button
              type="button"
              disabled={busy || !selectedContexts.length}
              onClick={() => void importLocal()}
            >
              Import selected ({selectedContexts.length})
            </Button>
          </div>
          {localQ.isError ? (
            <p className="text-sm text-warn">
              {(localQ.error as Error)?.message || 'Could not read local kubeconfig'}
            </p>
          ) : null}
          <div className="overflow-hidden rounded border border-line">
            <HudTableScroll maxHeightClass="max-h-64">
              <HudTable>
                <thead>
                  <tr>
                    <th className="w-10"></th>
                    <th>Context</th>
                    <th>Server</th>
                    <th>Conflict</th>
                  </tr>
                </thead>
                <tbody>
                  {localContexts.map((c) => {
                    const blocked = Boolean(c.conflict_name || c.conflict_server)
                    return (
                      <tr key={c.name} className={blocked ? 'opacity-60' : undefined}>
                        <td>
                          <input
                            type="checkbox"
                            disabled={blocked || busy}
                            checked={selectedContexts.includes(c.name)}
                            onChange={() => toggleContext(c.name)}
                          />
                        </td>
                        <td className="font-mono text-xs text-cyan">{c.name}</td>
                        <td className="font-mono text-[11px] text-text-dim">{c.server || '-'}</td>
                        <td className="text-xs">
                          {c.conflict_name ? (
                            <Badge tone="warn">name taken</Badge>
                          ) : c.conflict_server ? (
                            <Badge tone="warn">same as {c.existing_name}</Badge>
                          ) : (
                            <span className="text-text-dim">—</span>
                          )}
                        </td>
                      </tr>
                    )
                  })}
                  {!localQ.isLoading && !localContexts.length ? (
                    <tr>
                      <td colSpan={4}>
                        <EmptyState>No contexts found in local kubeconfig.</EmptyState>
                      </td>
                    </tr>
                  ) : null}
                </tbody>
              </HudTable>
            </HudTableScroll>
          </div>
        </Card>
      ) : null}

      {canWrite ? (
        <Card className="space-y-3 p-5">
          <h2 className="font-display text-lg font-bold tracking-[0.12em]">ADD CLUSTER</h2>
          <label className="block space-y-1">
            <span className="hud-label">Name</span>
            <input className="hud-field" value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Description</span>
            <input
              className="hud-field"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Kubeconfig</span>
            <textarea
              className="hud-field min-h-[160px] font-mono text-xs"
              value={kubeconfig}
              onChange={(e) => setKubeconfig(e.target.value)}
              placeholder="paste kubeconfig YAML"
            />
          </label>
          <Button type="button" disabled={busy} onClick={() => void create()}>
            Create cluster
          </Button>
        </Card>
      ) : null}

      <Modal
        open={Boolean(editTarget)}
        title="EDIT CLUSTER"
        subtitle={editTarget?.name}
        onClose={() => setEditTarget(null)}
      >
        <div className="space-y-3 px-5 py-4">
          <label className="block space-y-1">
            <span className="hud-label">Name</span>
            <input
              className="hud-field"
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Description</span>
            <input
              className="hud-field"
              value={editDescription}
              onChange={(e) => setEditDescription(e.target.value)}
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Kubeconfig (optional replace)</span>
            <textarea
              className="hud-field min-h-[140px] font-mono text-xs"
              value={editKubeconfig}
              onChange={(e) => setEditKubeconfig(e.target.value)}
              placeholder="leave empty to keep existing kubeconfig"
            />
          </label>
          <div className="flex justify-end gap-2">
            <Button variant="ghost" type="button" onClick={() => setEditTarget(null)}>
              Cancel
            </Button>
            <Button type="button" disabled={busy} onClick={() => void saveEdit()}>
              Save
            </Button>
          </div>
        </div>
      </Modal>

      <ConfirmDialog
        open={Boolean(deleteTarget)}
        title="DELETE CLUSTER"
        confirmText={deleteTarget?.name}
        confirmLabel="Delete cluster"
        busy={busy}
        description="Remove this cluster registration from CiliKube (does not delete the Kubernetes cluster)."
        onClose={() => setDeleteTarget(null)}
        onConfirm={remove}
      />
    </div>
  )
}
