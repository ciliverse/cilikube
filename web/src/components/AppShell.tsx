import { NavLink, Outlet, useLocation } from 'react-router-dom'
import { AnimatePresence, motion } from 'framer-motion'
import { useEffect, useState } from 'react'
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
  Menu,
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
  X,
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
      { to: '/gatewayclasses', label: 'GatewayClasses', icon: Layers, namespaced: false, resource: 'gatewayclasses' },
      { to: '/gateways', label: 'Gateways', icon: Workflow, namespaced: true, resource: 'gateways' },
      { to: '/httproutes', label: 'HTTPRoutes', icon: Network, namespaced: true, resource: 'httproutes' },
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

function NavBody({
  groups,
  onNavigate,
}: {
  groups: NavGroup[]
  onNavigate?: () => void
}) {
  return (
    <nav className="min-h-0 flex-1 space-y-4 overflow-y-auto overscroll-contain px-2 py-4">
      {groups.map((group) => (
        <div key={group.title}>
          <div className="hud-label mb-1.5 px-2">{group.title}</div>
          <div className="flex flex-col gap-0.5">
            {group.items.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.to === '/'}
                onClick={onNavigate}
                className={({ isActive }) =>
                  cn(
                    'flex min-h-11 items-center gap-2.5 rounded px-2.5 py-2.5 text-[13px] font-semibold tracking-wide transition md:min-h-0 md:py-2',
                    isActive
                      ? 'border border-cyan/40 bg-cyan/15 text-cyan shadow-[0_0_14px_rgba(53,230,255,0.12)]'
                      : 'border border-transparent text-text-dim hover:border-line hover:bg-mist hover:text-text',
                  )
                }
              >
                <item.icon className="h-3.5 w-3.5 shrink-0 opacity-90" />
                {item.label}
              </NavLink>
            ))}
          </div>
        </div>
      ))}
    </nav>
  )
}

export function AppShell() {
  const { user, roles, checkPermission, isViewerOnly, isAdmin } = useAuth()
  const { clusters, clusterId, setClusterId, switching } = useCluster()
  const { namespaces, namespace, setNamespace } = useNamespace()
  const location = useLocation()
  const [mobileNavOpen, setMobileNavOpen] = useState(false)

  useEffect(() => {
    setMobileNavOpen(false)
  }, [location.pathname])

  useEffect(() => {
    if (!mobileNavOpen) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setMobileNavOpen(false)
    }
    window.addEventListener('keydown', onKey)
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      window.removeEventListener('keydown', onKey)
      document.body.style.overflow = prev
    }
  }, [mobileNavOpen])

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

  const clusterOptions = clusters.map((c) => ({
    value: String(c.id || c.name),
    label: String(c.name || c.id),
  }))

  const reduceMotion =
    typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches

  return (
    <div className="relative flex h-dvh w-full flex-col overflow-hidden pb-[env(safe-area-inset-bottom)] pt-[env(safe-area-inset-top)]">
      <header className="z-30 flex h-12 w-full shrink-0 items-center gap-1.5 border-b border-line bg-panel-solid/95 px-2 backdrop-blur sm:h-14 sm:gap-3 sm:px-3 md:gap-4 md:px-5">
        <button
          type="button"
          className="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded border border-line text-cyan sm:h-10 sm:w-10 md:hidden"
          aria-label={mobileNavOpen ? 'Close navigation' : 'Open navigation'}
          aria-expanded={mobileNavOpen}
          onClick={() => setMobileNavOpen((v) => !v)}
        >
          {mobileNavOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </button>

        <div className="hud-brand hidden shrink-0 text-sm sm:block md:text-base">
          CILI<span className="accent">KUBE</span>
        </div>
        <div className="hud-brand shrink-0 text-xs tracking-[0.12em] sm:hidden">
          C<span className="accent">K</span>
        </div>

        <div className="min-w-0 flex-1">
          <GlobalSearchPalette />
        </div>

        <div className="flex shrink-0 items-center gap-1.5 sm:gap-2 md:gap-3">
          <label className="hidden items-center gap-2 lg:flex">
            <span className="hud-label">Cluster</span>
            <HudSelect
              aria-label="Cluster"
              className="w-auto min-w-[120px] max-w-[220px]"
              value={clusterId}
              onChange={setClusterId}
              disabled={switching || !clusters.length}
              options={clusterOptions}
            />
          </label>

          <label className="hidden items-center gap-2 lg:flex">
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

          <span className="hidden sm:inline-flex">
            <ConnDot online />
          </span>

          <UserMenu primaryRole={primaryRole} isViewerOnly={isViewerOnly} />
        </div>
      </header>

      <div className="relative flex min-h-0 w-full flex-1">
        {/* Desktop sidebar */}
        <aside className="hidden h-full w-[232px] shrink-0 flex-col border-r border-line bg-panel-solid/50 md:flex">
          <NavBody groups={visibleGroups} />
        </aside>

        {/* Mobile drawer */}
        <AnimatePresence>
          {mobileNavOpen ? (
            <>
              <motion.button
                type="button"
                key="nav-backdrop"
                className="fixed inset-0 z-40 bg-black/60 md:hidden"
                aria-label="Close navigation overlay"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 0.18 }}
                onClick={() => setMobileNavOpen(false)}
              />
              <motion.aside
                key="nav-drawer"
                className="fixed inset-y-0 left-0 z-50 flex w-[min(86vw,300px)] flex-col border-r border-line bg-panel-solid pt-[env(safe-area-inset-top)] shadow-[8px_0_32px_rgba(0,0,0,0.45)] md:hidden"
                initial={{ x: '-100%' }}
                animate={{ x: 0 }}
                exit={{ x: '-100%' }}
                transition={{ type: 'tween', duration: 0.22, ease: 'easeOut' }}
              >
                <div className="flex h-14 shrink-0 items-center justify-between border-b border-line px-3">
                  <div className="hud-brand text-sm">
                    CILI<span className="accent">KUBE</span>
                  </div>
                  <button
                    type="button"
                    className="inline-flex h-10 w-10 items-center justify-center rounded border border-line text-cyan"
                    aria-label="Close navigation"
                    onClick={() => setMobileNavOpen(false)}
                  >
                    <X className="h-5 w-5" />
                  </button>
                </div>

                <div className="shrink-0 space-y-3 border-b border-line p-3">
                  <div>
                    <div className="hud-label mb-2">Cluster</div>
                    <HudSelect
                      aria-label="Cluster"
                      className="w-full"
                      value={clusterId}
                      onChange={setClusterId}
                      disabled={switching || !clusters.length}
                      options={clusterOptions}
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

                <NavBody groups={visibleGroups} onNavigate={() => setMobileNavOpen(false)} />
              </motion.aside>
            </>
          ) : null}
        </AnimatePresence>

        <main className="relative flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden px-2.5 py-2.5 sm:px-4 sm:py-4 md:px-6 md:py-5 xl:px-8">
          <div className="mb-2.5 grid shrink-0 grid-cols-2 gap-2 lg:hidden">
            <div className="min-w-0">
              <div className="hud-label mb-1">Cluster</div>
              <HudSelect
                aria-label="Cluster"
                className="w-full"
                value={clusterId}
                onChange={setClusterId}
                disabled={switching || !clusters.length}
                options={clusterOptions}
              />
            </div>
            <div className="min-w-0">
              <div className="hud-label mb-1">NS</div>
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

          <div className="shrink-0">
            <OAuthAccountBanner />
          </div>

          <motion.div
            key={location.pathname}
            className="flex h-full min-h-0 w-full min-w-0 flex-1 flex-col overflow-y-auto overscroll-contain"
            initial={reduceMotion ? false : { opacity: 0, y: 6 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: reduceMotion ? 0 : 0.2, ease: 'easeOut' }}
          >
            <Outlet />
          </motion.div>
        </main>
      </div>
    </div>
  )
}
