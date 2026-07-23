import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useMemo, useState } from 'react'
import {
  createAdminRole,
  deleteAdminRole,
  getAvailablePermissions,
  getRolePermissions,
  listAdminRoles,
  setRolePermissions,
  updateAdminRole,
  type AdminRole,
} from '@/api/admin'
import { useAuth } from '@/store/auth'
import { Badge, Button, Card, EmptyState, Modal, PageHeader } from '@/components/ui'
import { HudTable, HudTableScroll } from '@/components/HudTableScroll'
import { ConfirmDialog } from '@/components/ConfirmDialog'

export function AdminRolesPage() {
  const { isAdmin } = useAuth()
  const queryClient = useQueryClient()
  const [selected, setSelected] = useState<AdminRole | null>(null)
  const [creating, setCreating] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<AdminRole | null>(null)
  const [checked, setChecked] = useState<string[]>([])
  const [busy, setBusy] = useState(false)
  const [err, setErr] = useState('')
  const [form, setForm] = useState({ name: '', display_name: '', description: '' })

  const rolesQ = useQuery({
    queryKey: ['admin-roles'],
    enabled: isAdmin,
    queryFn: listAdminRoles,
  })

  const permsQ = useQuery({
    queryKey: ['admin-permissions'],
    enabled: isAdmin,
    queryFn: getAvailablePermissions,
  })

  const rolePermsQ = useQuery({
    queryKey: ['admin-role-perms', selected?.id],
    enabled: Boolean(selected?.id),
    queryFn: () => getRolePermissions(selected!.id),
  })

  useEffect(() => {
    if (rolePermsQ.data) setChecked(rolePermsQ.data)
  }, [rolePermsQ.data])

  const roles = rolesQ.data || []
  const categories = permsQ.data || []

  const allPermNames = useMemo(
    () => categories.flatMap((c) => c.permissions.map((p) => p.name)),
    [categories],
  )

  if (!isAdmin) {
    return (
      <div className="rounded border border-warn/40 bg-warn/10 px-5 py-8 text-sm text-warn">
        Admin privileges required.
      </div>
    )
  }

  const togglePerm = (name: string) => {
    if (selected?.is_system) return
    setChecked((prev) =>
      prev.includes(name) ? prev.filter((p) => p !== name) : [...prev, name],
    )
  }

  const savePerms = async () => {
    if (!selected || selected.is_system) return
    setBusy(true)
    setErr('')
    try {
      await setRolePermissions(selected.id, checked)
      await updateAdminRole(selected.id, {
        display_name: selected.display_name || selected.name,
        description: selected.description || '',
        permissions: checked,
      })
      void queryClient.invalidateQueries({ queryKey: ['admin-roles'] })
      void queryClient.invalidateQueries({ queryKey: ['admin-role-perms', selected.id] })
    } catch (e: any) {
      setErr(e?.message || 'Failed to save permissions')
    } finally {
      setBusy(false)
    }
  }

  const create = async () => {
    setBusy(true)
    setErr('')
    try {
      if (!form.name.trim() || !form.display_name.trim()) {
        throw new Error('Name and display name are required')
      }
      await createAdminRole({
        name: form.name.trim(),
        display_name: form.display_name.trim(),
        description: form.description.trim(),
        permissions: [],
      })
      setCreating(false)
      setForm({ name: '', display_name: '', description: '' })
      void queryClient.invalidateQueries({ queryKey: ['admin-roles'] })
    } catch (e: any) {
      setErr(e?.message || 'Create failed')
    } finally {
      setBusy(false)
    }
  }

  const remove = async () => {
    if (!deleteTarget) return
    setBusy(true)
    try {
      await deleteAdminRole(deleteTarget.id)
      if (selected?.id === deleteTarget.id) setSelected(null)
      setDeleteTarget(null)
      void queryClient.invalidateQueries({ queryKey: ['admin-roles'] })
    } catch (e: any) {
      setErr(e?.message || 'Delete failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="flex w-full min-w-0 flex-col gap-4">
      <PageHeader
        title="ADMIN ROLES"
        subtitle="Application roles and permissions"
        action={
          <Button
            type="button"
            className="px-3 py-1.5 text-xs"
            onClick={() => {
              setErr('')
              setCreating(true)
            }}
          >
            Create role
          </Button>
        }
      />
      {err ? (
        <div className="rounded border border-danger/30 bg-danger/10 px-4 py-2 text-sm text-danger">{err}</div>
      ) : null}

      <div className="grid gap-4 lg:grid-cols-2">
        <Card className="overflow-hidden">
          <div className="border-b border-line px-4 py-3 font-display text-sm tracking-[0.12em]">
            ROLES ({roles.length})
          </div>
          <HudTableScroll maxHeightClass="max-h-[60vh]">
            <HudTable>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>System</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {roles.map((r) => {
                  const active = selected?.id === r.id
                  return (
                    <tr
                      key={r.id}
                      className={active ? 'bg-cyan/10' : 'cursor-pointer'}
                      onClick={() => setSelected(r)}
                    >
                      <td>
                        <div className="font-semibold text-cyan">{r.name}</div>
                        <div className="text-xs text-text-dim">{r.display_name || ''}</div>
                      </td>
                      <td>
                        <Badge tone={r.is_system ? 'warn' : 'neutral'}>
                          {r.is_system ? 'system' : 'custom'}
                        </Badge>
                      </td>
                      <td>
                        {!r.is_system ? (
                          <Button
                            variant="danger"
                            className="px-2 py-1 text-xs"
                            type="button"
                            onClick={(e) => {
                              e.stopPropagation()
                              setDeleteTarget(r)
                            }}
                          >
                            Delete
                          </Button>
                        ) : null}
                      </td>
                    </tr>
                  )
                })}
                {!rolesQ.isLoading && !roles.length ? (
                  <tr>
                    <td colSpan={3}>
                      <EmptyState>No roles.</EmptyState>
                    </td>
                  </tr>
                ) : null}
              </tbody>
            </HudTable>
          </HudTableScroll>
        </Card>

        <Card className="overflow-hidden">
          <div className="flex items-center justify-between border-b border-line px-4 py-3">
            <div className="font-display text-sm tracking-[0.12em]">
              PERMISSIONS {selected ? `· ${selected.name}` : ''}
            </div>
            {selected && !selected.is_system ? (
              <Button
                type="button"
                className="px-3 py-1.5 text-xs"
                disabled={busy}
                onClick={() => void savePerms()}
              >
                {busy ? 'Saving…' : 'Save'}
              </Button>
            ) : null}
          </div>
          {!selected ? (
            <EmptyState>Select a role to view permissions.</EmptyState>
          ) : selected.is_system ? (
            <div className="space-y-3 p-5 text-sm text-text-dim">
              <p>System role permissions are read-only.</p>
              <p>
                Loaded: {checked.length || allPermNames.length ? `${checked.length} assigned` : '—'}
              </p>
              <div className="flex flex-wrap gap-1">
                {checked.map((p) => (
                  <Badge key={p} tone="neutral">
                    {p}
                  </Badge>
                ))}
              </div>
            </div>
          ) : (
            <div className="max-h-[60vh] space-y-4 overflow-auto p-4">
              {categories.map((cat) => (
                <div key={cat.name}>
                  <div className="hud-label mb-2">{cat.display_name || cat.name}</div>
                  <div className="space-y-1">
                    {cat.permissions.map((p) => {
                      const on = checked.includes(p.name)
                      return (
                        <label
                          key={p.name}
                          className="flex cursor-pointer items-start gap-2 rounded border border-transparent px-2 py-1.5 hover:border-line"
                        >
                          <input
                            type="checkbox"
                            checked={on}
                            onChange={() => togglePerm(p.name)}
                            className="mt-0.5"
                          />
                          <span>
                            <span className="font-mono text-xs text-cyan">{p.name}</span>
                            <span className="mt-0.5 block text-xs text-text-dim">
                              {p.display_name || p.description || ''}
                            </span>
                          </span>
                        </label>
                      )
                    })}
                  </div>
                </div>
              ))}
              {!categories.length ? <EmptyState>No permission catalog.</EmptyState> : null}
            </div>
          )}
        </Card>
      </div>

      <Modal open={creating} title="CREATE ROLE" onClose={() => setCreating(false)}>
        <div className="space-y-3 p-5">
          <label className="block space-y-1">
            <span className="hud-label">Name</span>
            <input
              className="hud-field"
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              placeholder="e.g. operator"
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Display name</span>
            <input
              className="hud-field"
              value={form.display_name}
              onChange={(e) => setForm((f) => ({ ...f, display_name: e.target.value }))}
            />
          </label>
          <label className="block space-y-1">
            <span className="hud-label">Description</span>
            <textarea
              className="hud-field min-h-[80px]"
              value={form.description}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            />
          </label>
          <div className="flex justify-end gap-2">
            <Button variant="ghost" type="button" onClick={() => setCreating(false)}>
              Cancel
            </Button>
            <Button type="button" disabled={busy} onClick={() => void create()}>
              Create
            </Button>
          </div>
        </div>
      </Modal>

      <ConfirmDialog
        open={Boolean(deleteTarget)}
        title="DELETE ROLE"
        confirmText={deleteTarget?.name}
        confirmLabel="Delete role"
        busy={busy}
        description="Remove this custom role. Users with only this role may lose access."
        onClose={() => setDeleteTarget(null)}
        onConfirm={remove}
      />
    </div>
  )
}
