import { NavLink, Outlet, useLocation } from 'react-router-dom'
import { AnimatePresence, motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
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
  labelKey: string
  icon: typeof Server
  namespaced?: boolean
  resource?: string
}
type NavGroup = { titleKey: string; items: NavItem[] }

const navGroups: NavGroup[] = [
  {
    titleKey: 'nav.cluster',
    items: [
      { to: '/', labelKey: 'nav.overview', icon: LayoutDashboard },
      { to: '/nodes', labelKey: 'nav.nodes', icon: Server, resource: 'nodes' },
      { to: '/namespaces', labelKey: 'nav.namespaces', icon: Layers, resource: 'namespaces' },
      { to: '/events', labelKey: 'nav.events', icon: Activity },
      { to: '/clusters', labelKey: 'nav.clusters', icon: Cloud, resource: 'clusters' },
      { to: '/crds', labelKey: 'nav.crds', icon: FileCode2 },
    ],
  },
  {
    titleKey: 'nav.workloads',
    items: [
      { to: '/pods', labelKey: 'nav.pods', icon: Boxes, namespaced: true, resource: 'pods' },
      { to: '/deployments', labelKey: 'nav.deployments', icon: Boxes, namespaced: true, resource: 'deployments' },
      { to: '/statefulsets', labelKey: 'nav.statefulsets', icon: Boxes, namespaced: true, resource: 'statefulsets' },
      { to: '/daemonsets', labelKey: 'nav.daemonsets', icon: Boxes, namespaced: true, resource: 'daemonsets' },
      { to: '/jobs', labelKey: 'nav.jobs', icon: Workflow, namespaced: true, resource: 'jobs' },
      { to: '/cronjobs', labelKey: 'nav.cronjobs', icon: CalendarClock, namespaced: true, resource: 'cronjobs' },
      {
        to: '/horizontalpodautoscalers',
        labelKey: 'nav.hpa',
        icon: Gauge,
        namespaced: true,
        resource: 'horizontalpodautoscalers',
      },
      {
        to: '/poddisruptionbudgets',
        labelKey: 'nav.pdb',
        icon: Shield,
        namespaced: true,
        resource: 'poddisruptionbudgets',
      },
    ],
  },
  {
    titleKey: 'nav.network',
    items: [
      { to: '/services', labelKey: 'nav.services', icon: Network, namespaced: true, resource: 'services' },
      { to: '/ingresses', labelKey: 'nav.ingress', icon: Network, namespaced: true, resource: 'ingresses' },
      { to: '/gatewayclasses', labelKey: 'nav.gatewayclasses', icon: Layers, namespaced: false, resource: 'gatewayclasses' },
      { to: '/gateways', labelKey: 'nav.gateways', icon: Workflow, namespaced: true, resource: 'gateways' },
      { to: '/httproutes', labelKey: 'nav.httproutes', icon: Network, namespaced: true, resource: 'httproutes' },
      { to: '/networkpolicies', labelKey: 'nav.networkpolicies', icon: Shield, namespaced: true, resource: 'networkpolicies' },
    ],
  },
  {
    titleKey: 'nav.config',
    items: [
      { to: '/configmaps', labelKey: 'nav.configmaps', icon: Database, namespaced: true, resource: 'configmaps' },
      { to: '/secrets', labelKey: 'nav.secrets', icon: KeyRound, namespaced: true, resource: 'secrets' },
      { to: '/serviceaccounts', labelKey: 'nav.serviceaccounts', icon: UserRound, namespaced: true, resource: 'serviceaccounts' },
      {
        to: '/resourcequotas',
        labelKey: 'nav.resourcequotas',
        icon: Database,
        namespaced: true,
        resource: 'resourcequotas',
      },
      {
        to: '/limitranges',
        labelKey: 'nav.limitranges',
        icon: Gauge,
        namespaced: true,
        resource: 'limitranges',
      },
    ],
  },
  {
    titleKey: 'nav.storage',
    items: [
      { to: '/persistentvolumes', labelKey: 'nav.pv', icon: HardDrive, resource: 'persistentvolumes' },
      {
        to: '/persistentvolumeclaims',
        labelKey: 'nav.pvc',
        icon: HardDrive,
        namespaced: true,
        resource: 'persistentvolumeclaims',
      },
      { to: '/storageclasses', labelKey: 'nav.storageclass', icon: HardDrive, resource: 'storageclasses' },
    ],
  },
  {
    titleKey: 'nav.access',
    items: [
      { to: '/roles', labelKey: 'nav.roles', icon: Lock, namespaced: true, resource: 'roles' },
      { to: '/rolebindings', labelKey: 'nav.rolebindings', icon: Lock, namespaced: true, resource: 'rolebindings' },
      { to: '/clusterroles', labelKey: 'nav.clusterroles', icon: Shield, resource: 'clusterroles' },
      { to: '/clusterrolebindings', labelKey: 'nav.clusterrolebindings', icon: Shield, resource: 'clusterrolebindings' },
    ],
  },
  {
    titleKey: 'nav.observe',
    items: [
      { to: '/monitoring', labelKey: 'nav.monitoring', icon: Activity },
      { to: '/audit', labelKey: 'nav.audit', icon: ScrollText },
      { to: '/proxy', labelKey: 'nav.proxy', icon: Terminal },
      { to: '/helm', labelKey: 'nav.helm', icon: Workflow, namespaced: true },
    ],
  },
  {
    titleKey: 'nav.admin',
    items: [
      { to: '/admin/users', labelKey: 'nav.users', icon: UserRound },
      { to: '/admin/roles', labelKey: 'nav.adminRoles', icon: Shield },
      { to: '/admin/settings', labelKey: 'nav.settings', icon: Settings },
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
  const { t } = useTranslation()
  return (
    <nav className="min-h-0 flex-1 space-y-4 overflow-y-auto overscroll-contain px-2 py-4">
      {groups.map((group) => (
        <div key={group.titleKey}>
          <div className="hud-label mb-1.5 px-2">{t(group.titleKey)}</div>
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
                {t(item.labelKey)}
              </NavLink>
            ))}
          </div>
        </div>
      ))}
    </nav>
  )
}

export function AppShell() {
  const { t } = useTranslation()
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
    { value: ALL_NAMESPACES, label: t('common.allNamespaces') },
    ...namespaces.map((ns) => ({ value: ns, label: ns })),
  ]

  const clusterOptions = clusters.map((c) => ({
    value: String(c.id || c.name),
    label: String(c.name || c.id),
  }))

  const reduceMotion =
    typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches

  const clusterLabel = t('nav.cluster')
  const nsLabel = t('common.namespace')

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
            <span className="hud-label">{clusterLabel}</span>
            <HudSelect
              aria-label={clusterLabel}
              className="w-auto min-w-[120px] max-w-[220px]"
              value={clusterId}
              onChange={setClusterId}
              disabled={switching || !clusters.length}
              options={clusterOptions}
            />
          </label>

          <label className="hidden items-center gap-2 lg:flex">
            <span className="hud-label">{nsLabel}</span>
            <HudSelect
              aria-label={nsLabel}
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
                    <div className="hud-label mb-2">{clusterLabel}</div>
                    <HudSelect
                      aria-label={clusterLabel}
                      className="w-full"
                      value={clusterId}
                      onChange={setClusterId}
                      disabled={switching || !clusters.length}
                      options={clusterOptions}
                    />
                  </div>
                  <div>
                    <div className="hud-label mb-2">{nsLabel}</div>
                    <HudSelect
                      aria-label={nsLabel}
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
              <div className="hud-label mb-1">{clusterLabel}</div>
              <HudSelect
                aria-label={clusterLabel}
                className="w-full"
                value={clusterId}
                onChange={setClusterId}
                disabled={switching || !clusters.length}
                options={clusterOptions}
              />
            </div>
            <div className="min-w-0">
              <div className="hud-label mb-1">{nsLabel}</div>
              <HudSelect
                aria-label={nsLabel}
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
