import type { UserInfo } from '@/api/auth'

export type ResourceAction = 'read' | 'write' | 'delete' | 'exec'

export function userRoles(user: UserInfo | null | undefined): string[] {
  if (!user) return []
  const roles = [...(user.roles || [])]
  if (user.role && !roles.includes(user.role)) roles.push(user.role)
  return roles
}

export function isAdmin(user: UserInfo | null | undefined): boolean {
  return userRoles(user).includes('admin') || userRoles(user).includes('super_admin')
}

/** Admin or editor may mutate most workloads. */
export function canEdit(user: UserInfo | null | undefined): boolean {
  const roles = userRoles(user)
  return roles.some((r) => ['admin', 'super_admin', 'editor'].includes(r))
}

export function isViewerOnly(user: UserInfo | null | undefined): boolean {
  return userRoles(user).includes('viewer') && !canEdit(user)
}

const readOnlyForEditor = new Set([
  'secret',
  'secrets',
  'cluster',
  'clusters',
  'cluster-management',
  'role',
  'roles',
  'rolebinding',
  'rolebindings',
  'clusterrole',
  'clusterroles',
  'clusterrolebinding',
  'clusterrolebindings',
  'node',
  'nodes',
  'namespace',
  'namespaces',
])

const noAccessForEditor = new Set(['secret', 'secrets'])

/**
 * Frontend gate aligned with Vue helpers + Casbin default roles.
 * Backend remains the source of truth; this only hides unsafe controls.
 */
export function checkResourcePermission(
  user: UserInfo | null | undefined,
  resource: string,
  action: ResourceAction,
): boolean {
  if (isAdmin(user)) return true

  const roles = userRoles(user)
  const key = resource.toLowerCase()

  if (roles.includes('editor')) {
    if (noAccessForEditor.has(key)) return false
    if (readOnlyForEditor.has(key)) return action === 'read'
    if (action === 'exec') return key === 'pod' || key === 'pods'
    return action === 'read' || action === 'write' || action === 'delete'
  }

  if (roles.includes('viewer') || roles.includes('normal_user')) {
    if (action !== 'read') return false
    if (key === 'secret' || key === 'secrets') return false
    return true
  }

  // Unknown role: read-only in UI
  return action === 'read' && key !== 'secrets' && key !== 'secret'
}

export function canMutateResource(user: UserInfo | null | undefined, resource: string): boolean {
  return checkResourcePermission(user, resource, 'write')
}

export function canDeleteResource(user: UserInfo | null | undefined, resource: string): boolean {
  return checkResourcePermission(user, resource, 'delete')
}

export function canExecPods(user: UserInfo | null | undefined): boolean {
  return checkResourcePermission(user, 'pods', 'exec')
}
