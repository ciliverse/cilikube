import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import {
  createAdminUser,
  deleteAdminUser,
  listAdminRoles,
  listAdminUsers,
  updateAdminUser,
  updateAdminUserStatus,
  type AdminUser,
} from '@/api/admin'
import { useAuth } from '@/store/auth'
import { Badge, Button, EmptyState, Modal, PageHeader } from '@/components/ui'
import { HudTable, HudTablePanel, ListPageFrame } from '@/components/HudTableScroll'
import { ConfirmDialog } from '@/components/ConfirmDialog'

const ROLE_CHOICES = ['viewer', 'editor', 'admin'] as const

export function AdminUsersPage() {
  const { isAdmin } = useAuth()
  const queryClient = useQueryClient()
  const [deleteTarget, setDeleteTarget] = useState<AdminUser | null>(null)
  const [editTarget, setEditTarget] = useState<AdminUser | null>(null)
  const [creating, setCreating] = useState(false)
  const [busy, setBusy] = useState(false)
  const [err, setErr] = useState('')

  const [form, setForm] = useState({
    username: '',
    email: '',
    display_name: '',
    password: '',
    confirmPassword: '',
    roles: ['viewer'] as string[],
  })

  const usersQ = useQuery({
    queryKey: ['admin-users'],
    enabled: isAdmin,
    queryFn: () => listAdminUsers(),
  })

  const rolesQ = useQuery({
    queryKey: ['admin-roles'],
    enabled: isAdmin,
    queryFn: listAdminRoles,
  })

  const roleNames =
    rolesQ.data?.map((r) => r.name).filter(Boolean) || [...ROLE_CHOICES]

  if (!isAdmin) {
    return (
      <div className="rounded border border-warn/40 bg-warn/10 px-5 py-8 text-sm text-warn">
        Admin privileges required.
      </div>
    )
  }

  const users = usersQ.data || []

  const resetForm = (user?: AdminUser | null) => {
    if (user) {
      setForm({
        username: user.username,
        email: user.email || '',
        display_name: user.display_name || '',
        password: '',
        confirmPassword: '',
        roles: user.roles?.length ? [...user.roles] : ['viewer'],
      })
    } else {
      setForm({
        username: '',
        email: '',
        display_name: '',
        password: '',
        confirmPassword: '',
        roles: ['viewer'],
      })
    }
  }

  const toggleRole = (name: string) => {
    setForm((f) => {
      const has = f.roles.includes(name)
      if (has) return { ...f, roles: f.roles.filter((r) => r !== name) }
      return { ...f, roles: [...f.roles, name] }
    })
  }

  const create = async () => {
    setBusy(true)
    setErr('')
    try {
      if (!form.username || !form.email || !form.password) {
        throw new Error('Username, email and password are required')
      }
      if (form.password !== form.confirmPassword) {
        throw new Error('Password confirmation does not match')
      }
      await createAdminUser({
        username: form.username.trim(),
        email: form.email.trim(),
        password: form.password,
        confirmPassword: form.confirmPassword,
        display_name: form.display_name.trim(),
        roles: form.roles.join(','),
      })
      // Assign roles via update (create API stores roles string but assignment is on update)
      const fresh = await listAdminUsers()
      const created = fresh.find((u) => u.username === form.username.trim())
      if (created && form.roles.length) {
        await updateAdminUser(created.id, {
          email: form.email.trim(),
          display_name: form.display_name.trim(),
          roles: form.roles,
        })
      }
      setCreating(false)
      resetForm()
      void queryClient.invalidateQueries({ queryKey: ['admin-users'] })
    } catch (e: any) {
      setErr(e?.message || 'Create failed')
    } finally {
      setBusy(false)
    }
  }

  const saveEdit = async () => {
    if (!editTarget) return
    setBusy(true)
    setErr('')
    try {
      await updateAdminUser(editTarget.id, {
        email: form.email.trim(),
        display_name: form.display_name.trim(),
        roles: form.roles.length ? form.roles : ['viewer'],
        is_active: editTarget.is_active,
      })
      setEditTarget(null)
      void queryClient.invalidateQueries({ queryKey: ['admin-users'] })
    } catch (e: any) {
      setErr(e?.message || 'Update failed')
    } finally {
      setBusy(false)
    }
  }

  const toggle = async (user: AdminUser) => {
    setBusy(true)
    setErr('')
    try {
      await updateAdminUserStatus(user.id, !(user.is_active ?? true))
      await usersQ.refetch()
    } catch (e: any) {
      setErr(e?.message || 'Update failed')
    } finally {
      setBusy(false)
    }
  }

  const remove = async () => {
    if (!deleteTarget) return
    setBusy(true)
    try {
      await deleteAdminUser(deleteTarget.id)
      setDeleteTarget(null)
      void queryClient.invalidateQueries({ queryKey: ['admin-users'] })
    } catch (e: any) {
      setErr(e?.message || 'Delete failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <ListPageFrame className="gap-3">
      <PageHeader
        title="ADMIN USERS"
        subtitle="Application user accounts"
        action={
          <Button
            type="button"
            className="px-3 py-1.5 text-xs"
            onClick={() => {
              resetForm()
              setErr('')
              setCreating(true)
            }}
          >
            Create user
          </Button>
        }
      />
      {err ? (
        <div className="shrink-0 rounded border border-danger/30 bg-danger/10 px-4 py-2 text-sm text-danger">{err}</div>
      ) : null}
      <HudTablePanel>
          <HudTable>
            <thead>
              <tr>
                <th>ID</th>
                <th>Username</th>
                <th>Email</th>
                <th>Roles</th>
                <th>Active</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id}>
                  <td>{u.id}</td>
                  <td className="font-semibold text-cyan">{u.username}</td>
                  <td>{u.email || '-'}</td>
                  <td>
                    <div className="flex flex-wrap gap-1">
                      {(u.roles?.length ? u.roles : ['-']).map((r) => (
                        <Badge key={r} tone="accent">
                          {r}
                        </Badge>
                      ))}
                    </div>
                  </td>
                  <td>
                    <Badge tone={u.is_active !== false ? 'ok' : 'danger'}>
                      {u.is_active !== false ? 'yes' : 'no'}
                    </Badge>
                  </td>
                  <td>
                    <div className="flex flex-wrap gap-1">
                      <Button
                        variant="outline"
                        className="px-2 py-1 text-xs"
                        type="button"
                        onClick={() => {
                          resetForm(u)
                          setErr('')
                          setEditTarget(u)
                        }}
                      >
                        Edit
                      </Button>
                      <Button
                        variant="outline"
                        className="px-2 py-1 text-xs"
                        type="button"
                        disabled={busy}
                        onClick={() => void toggle(u)}
                      >
                        Toggle
                      </Button>
                      <Button
                        variant="danger"
                        className="px-2 py-1 text-xs"
                        type="button"
                        onClick={() => setDeleteTarget(u)}
                      >
                        Delete
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
              {!usersQ.isLoading && !users.length ? (
                <tr>
                  <td colSpan={6}>
                    <EmptyState>No users found.</EmptyState>
                  </td>
                </tr>
              ) : null}
            </tbody>
          </HudTable>
      </HudTablePanel>

      <Modal
        open={creating || Boolean(editTarget)}
        title={creating ? 'CREATE USER' : 'EDIT USER'}
        subtitle={editTarget?.username}
        onClose={() => {
          setCreating(false)
          setEditTarget(null)
        }}
      >
        <div className="space-y-3 p-5">
          {creating ? (
            <label className="block space-y-1">
              <span className="hud-label">Username</span>
              <input
                className="hud-field"
                value={form.username}
                onChange={(e) => setForm((f) => ({ ...f, username: e.target.value }))}
              />
            </label>
          ) : null}
          <label className="block space-y-1">
            <span className="hud-label">Email</span>
            <input
              className="hud-field"
              value={form.email}
              onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
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
          {creating ? (
            <>
              <label className="block space-y-1">
                <span className="hud-label">Password</span>
                <input
                  type="password"
                  className="hud-field"
                  value={form.password}
                  onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
                />
              </label>
              <label className="block space-y-1">
                <span className="hud-label">Confirm password</span>
                <input
                  type="password"
                  className="hud-field"
                  value={form.confirmPassword}
                  onChange={(e) => setForm((f) => ({ ...f, confirmPassword: e.target.value }))}
                />
              </label>
            </>
          ) : null}
          <div>
            <div className="hud-label mb-2">Roles</div>
            <div className="flex flex-wrap gap-2">
              {roleNames.map((name) => {
                const on = form.roles.includes(name)
                return (
                  <button
                    key={name}
                    type="button"
                    className={
                      on
                        ? 'rounded border border-cyan/40 bg-cyan/15 px-2 py-1 text-xs text-cyan'
                        : 'rounded border border-line px-2 py-1 text-xs text-text-dim hover:border-cyan/40'
                    }
                    onClick={() => toggleRole(name)}
                  >
                    {name}
                  </button>
                )
              })}
            </div>
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button
              variant="ghost"
              type="button"
              onClick={() => {
                setCreating(false)
                setEditTarget(null)
              }}
            >
              Cancel
            </Button>
            <Button
              type="button"
              disabled={busy}
              onClick={() => void (creating ? create() : saveEdit())}
            >
              {busy ? 'Saving…' : creating ? 'Create' : 'Save'}
            </Button>
          </div>
        </div>
      </Modal>

      <ConfirmDialog
        open={Boolean(deleteTarget)}
        title="DELETE USER"
        confirmText={deleteTarget?.username}
        confirmLabel="Delete user"
        busy={busy}
        description="This permanently removes the application user account."
        onClose={() => setDeleteTarget(null)}
        onConfirm={remove}
      />
    </ListPageFrame>
  )
}
