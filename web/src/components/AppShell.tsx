import { NavLink, Outlet, useLocation } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Activity,
  Boxes,
  CalendarClock,
  Database,
  Gauge,
  HardDrive,
  KeyRound,
  LayoutDashboard,
  Lock,
  Network,
  Server,
  Settings,
  Shield,
  Layers,
  UserRound,
  Workflow,
  Cloud,
  FileCode2,
  ScrollText,
  Terminal,
} from 'lucide-react'
import { useAuth } from '@/store/auth'
import { useCluster } from '@/store/cluster'
import { ALL_NAMESPACES, useNamespace } from '@/store/namespace'
import { ConnDot, HudSelect } from './ui'
import { GlobalSearchPalette } from './GlobalSearchPalette'
import { OAuthAccountBanner } from './OAuthAccountBanner'
import { UserMenu } from './UserMenu'
import { cn } from '@/lib/utils'

type NavItem = {
  to: string
  label: string
  icon: typeof Server
  namespaced?: boolean
  resource?: string
}
type NavGroup = { title: string; items: NavItem[] }

const navGroups: NavGroup[] = [
  {
    title: 'Cluster',
    items: [
      { to: '/', label: 'Overview', icon: LayoutDashboard },
      { to: '/nodes', label: 'Nodes', icon: Server, resource: 'nodes' },
      { to: '/namespaces', label: 'Namespaces', icon: Layers, resource: 'namespaces' },
      { to: '/events', label: 'Events', icon: Activity },
      { to: '/clusters', label: 'Clusters', icon: Cloud, resource: 'clusters' },
      { to: '/crds', label: 'CRDs', icon: FileCode2 },
    ],
  },
  {
    title: 'Workloads',
    items: [
      { to: '/pods', label: 'Pods', icon: Boxes, namespaced: true, resource: 'pods' },
      { to: '/deployments', label: 'Deployments', icon: Boxes, namespaced: true, resource: 'deployments' },
      { to: '/statefulsets', label: 'StatefulSets', icon: Boxes, namespaced: true, resource: 'statefulsets' },
      { to: '/daemonsets', label: 'DaemonSets', icon: Boxes, namespaced: true, resource: 'daemonsets' },
      { to: '/jobs', label: 'Jobs', icon: Workflow, namespaced: true, resource: 'jobs' },
      { to: '/cronjobs', label: 'CronJobs', icon: CalendarClock, namespaced: true, resource: 'cronjobs' },
      {
        to: '/horizontalpodautoscalers',
        label: 'HPA',
        icon: Gauge,
        namespaced: true,
        resource: 'horizontalpodautoscalers',
      },
      {
        to: '/poddisruptionbudgets',
        label: 'PDB',
        icon: Shield,
        namespaced: true,
        resource: 'poddisruptionbudgets',
      },
    ],
  },
  {
    title: 'Network',
    items: [
      { to: '/services', label: 'Services', icon: Network, namespaced: true, resource: 'services' },
      { to: '/ingresses', label: 'Ingress', icon: Network, namespaced: true, resource: 'ingresses' },
      { to: '/networkpolicies', label: 'NetworkPolicies', icon: Shield, namespaced: true, resource: 'networkpolicies' },
    ],
  },
  {
    title: 'Config',
    items: [
      { to: '/configmaps', label: 'ConfigMaps', icon: Database, namespaced: true, resource: 'configmaps' },
      { to: '/secrets', label: 'Secrets', icon: KeyRound, namespaced: true, resource: 'secrets' },
      { to: '/serviceaccounts', label: 'ServiceAccounts', icon: UserRound, namespaced: true, resource: 'serviceaccounts' },
      {
        to: '/resourcequotas',
        label: 'ResourceQuotas',
        icon: Database,
        namespaced: true,
        resource: 'resourcequotas',
      },
      {
        to: '/limitranges',
        label: 'LimitRanges',
        icon: Gauge,
        namespaced: true,
        resource: 'limitranges',
      },
    ],
  },
  {
    title: 'Storage',
    items: [
      { to: '/persistentvolumes', label: 'PV', icon: HardDrive, resource: 'persistentvolumes' },
      {
        to: '/persistentvolumeclaims',
        label: 'PVC',
        icon: HardDrive,
        namespaced: true,
        resource: 'persistentvolumeclaims',
      },
      { to: '/storageclasses', label: 'StorageClass', icon: HardDrive, resource: 'storageclasses' },
    ],
  },
  {
    title: 'Access',
    items: [
      { to: '/roles', label: 'Roles', icon: Lock, namespaced: true, resource: 'roles' },
      { to: '/rolebindings', label: 'RoleBindings', icon: Lock, namespaced: true, resource: 'rolebindings' },
      { to: '/clusterroles', label: 'ClusterRoles', icon: Shield, resource: 'clusterroles' },
      { to: '/clusterrolebindings', label: 'ClusterRoleBindings', icon: Shield, resource: 'clusterrolebindings' },
    ],
  },
  {
    title: 'Observe',
    items: [
      { to: '/monitoring', label: 'Monitoring', icon: Activity },
      { to: '/audit', label: 'Audit', icon: ScrollText },
      { to: '/proxy', label: 'API Proxy', icon: Terminal },
      { to: '/helm', label: 'Helm', icon: Workflow, namespaced: true },
    ],
  },
  {
    title: 'Admin',
    items: [
      { to: '/admin/users', label: 'Users', icon: UserRound },
      { to: '/admin/roles', label: 'Roles', icon: Shield },
      { to: '/admin/settings', label: 'Settings', icon: Settings },
    ],
  },
]

export function AppShell() {
  const { user, roles, checkPermission, isViewerOnly, isAdmin } = useAuth()
  const { clusters, clusterId, setClusterId, switching } = useCluster()
  const { namespaces, namespace, setNamespace } = useNamespace()
  const location = useLocation()

  const visibleGroups = navGroups
    .map((group) => ({
      ...group,
      items: group.items.filter((item) => {
        if (item.to.startsWith('/admin') || item.to === '/audit') return isAdmin
        if (item.to === '/proxy') return !isViewerOnly
        if (item.resource) return checkPermission(item.resource, 'read')
        return true
      }),
    }))
    .filter((group) => group.items.length > 0)

  const primaryRole = roles.includes('admin')
    ? 'admin'
    : roles.includes('editor')
      ? 'editor'
      : roles.includes('viewer')
        ? 'viewer'
        : user?.role || 'user'

  const nsOptions = [
    { value: ALL_NAMESPACES, label: 'All namespaces' },
    ...namespaces.map((ns) => ({ value: ns, label: ns })),
  ]

  return (
    <div className="relative flex h-dvh w-full flex-col overflow-hidden">
      <header className="z-20 flex h-14 w-full shrink-0 items-center gap-3 border-b border-line bg-panel-solid/95 px-3 backdrop-blur md:gap-4 md:px-5">
        <div className="hud-brand shrink-0 text-sm md:text-base">
          CILI<span className="accent">KUBE</span>
        </div>

        <div className="min-w-0 flex-1">
          <GlobalSearchPalette />
        </div>

        <div className="flex shrink-0 items-center gap-2 md:gap-3">
          <label className="hidden items-center gap-2 sm:flex">
            <span className="hud-label">Cluster</span>
            <HudSelect
              aria-label="Cluster"
              className="w-auto min-w-[120px] max-w-[220px]"
              value={clusterId}
              onChange={setClusterId}
              disabled={switching || !clusters.length}
              options={clusters.map((c) => ({
                value: String(c.id || c.name),
                label: String(c.name || c.id),
              }))}
            />
          </label>

          <label className="hidden items-center gap-2 md:flex">
            <span className="hud-label">NS</span>
            <HudSelect
              aria-label="Namespace"
              className="w-auto min-w-[120px] max-w-[200px]"
              value={namespace}
              onChange={setNamespace}
              searchableWhen={0}
              options={nsOptions}
            />
          </label>

          <ConnDot online />

          <UserMenu primaryRole={primaryRole} isViewerOnly={isViewerOnly} />
        </div>
      </header>

      <div className="flex min-h-0 w-full flex-1">
        <aside className="flex h-full w-[232px] shrink-0 flex-col border-r border-line bg-panel-solid/50">
          <nav className="min-h-0 flex-1 space-y-4 overflow-y-auto overscroll-contain px-2 py-4">
            {visibleGroups.map((group) => (
              <div key={group.title}>
                <div className="hud-label mb-1.5 px-2">{group.title}</div>
                <div className="flex flex-col gap-0.5">
                  {group.items.map((item) => (
                    <NavLink
                      key={item.to}
                      to={item.to}
                      end={item.to === '/'}
                      className={({ isActive }) =>
                        cn(
                          'flex items-center gap-2.5 rounded px-2.5 py-2 text-[13px] font-semibold tracking-wide transition',
                          isActive
                            ? 'border border-cyan/40 bg-cyan/15 text-cyan shadow-[0_0_14px_rgba(53,230,255,0.12)]'
                            : 'border border-transparent text-text-dim hover:border-line hover:bg-mist hover:text-text',
                        )
                      }
                    >
                      <item.icon className="h-3.5 w-3.5 opacity-90" />
                      {item.label}
                    </NavLink>
                  ))}
                </div>
              </div>
            ))}
          </nav>

          <div className="shrink-0 space-y-3 border-t border-line p-3 sm:hidden">
            <div>
              <div className="hud-label mb-2">Cluster</div>
              <HudSelect
                aria-label="Cluster"
                className="w-full"
                value={clusterId}
                onChange={setClusterId}
                disabled={switching || !clusters.length}
                options={clusters.map((c) => ({
                  value: String(c.id || c.name),
                  label: String(c.name || c.id),
                }))}
              />
            </div>
            <div>
              <div className="hud-label mb-2">Namespace</div>
              <HudSelect
                aria-label="Namespace"
                className="w-full"
                value={namespace}
                onChange={setNamespace}
                searchableWhen={0}
                options={nsOptions}
              />
            </div>
          </div>
        </aside>

        <main className="relative flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden px-4 py-4 md:px-6 md:py-5 xl:px-8">
          <div className="mb-3 flex shrink-0 items-center gap-2 md:hidden">
            <span className="hud-label">Namespace</span>
            <HudSelect
              aria-label="Namespace"
              className="min-w-0 flex-1"
              value={namespace}
              onChange={setNamespace}
              searchableWhen={0}
              options={nsOptions}
            />
          </div>

          <div className="shrink-0">
            <OAuthAccountBanner />
          </div>

          <motion.div
            key={location.pathname}
            className="flex h-full min-h-0 w-full min-w-0 flex-1 flex-col overflow-y-auto overscroll-contain"
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.28, ease: 'easeOut' }}
          >
            <Outlet />
          </motion.div>
        </main>
      </div>
    </div>
  )
}
