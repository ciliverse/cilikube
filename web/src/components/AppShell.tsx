import { Link, NavLink, Outlet, useLocation } from 'react-router-dom'
import { AnimatePresence, motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { useTranslation } from 'react-i18next'
import {
  Activity,
  Boxes,
  CalendarClock,
  ChevronDown,
  Cloud,
  Database,
  FileCode2,
  Gauge,
  HardDrive,
  History,
  KeyRound,
  Layers,
  LayoutDashboard,
  LayoutGrid,
  Lock,
  Menu,
  Network,
  ScrollText,
  Server,
  Settings,
  Shield,
  Sparkles,
  Terminal,
  UserRound,
  Waypoints,
  Workflow,
  X,
} from 'lucide-react'
import { shouldSkipEnterAnim } from '@/lib/motionPrefs'
import { useAuth } from '@/store/auth'
import { ConnDot } from './ui'
import { BrandMark } from './BrandMark'
import { ClusterNamespaceControls } from './ClusterNamespaceControls'
import { GlobalSearchPalette } from './GlobalSearchPalette'
import { OAuthAccountBanner } from './OAuthAccountBanner'
import { StarSupportFloat } from './StarSupportCta'
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
      { to: '/fleet', labelKey: 'nav.fleet', icon: LayoutGrid },
      { to: '/overview', labelKey: 'nav.overview', icon: LayoutDashboard },
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
      { to: '/topology', labelKey: 'nav.topology', icon: Waypoints, namespaced: true },
      { to: '/timeline', labelKey: 'nav.timeline', icon: History },
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

function pathInGroup(pathname: string, group: NavGroup): boolean {
  return group.items.some(
    (item) => pathname === item.to || pathname.startsWith(`${item.to}/`),
  )
}

function NavBody({
  groups,
  onNavigate,
}: {
  groups: NavGroup[]
  onNavigate?: () => void
}) {
  const { t } = useTranslation()
  const location = useLocation()
  const activeGroupKey =
    groups.find((group) => pathInGroup(location.pathname, group))?.titleKey ?? null

  const [openByKey, setOpenByKey] = useState<Record<string, boolean>>(() => {
    const initial: Record<string, boolean> = {}
    for (const group of groups) {
      initial[group.titleKey] = pathInGroup(location.pathname, group)
    }
    return initial
  })

  useEffect(() => {
    if (!activeGroupKey) return
    setOpenByKey((prev) =>
      prev[activeGroupKey] ? prev : { ...prev, [activeGroupKey]: true },
    )
  }, [activeGroupKey])

  const toggleGroup = (titleKey: string) => {
    setOpenByKey((prev) => ({ ...prev, [titleKey]: !prev[titleKey] }))
  }

  return (
    <nav className="min-h-0 flex-1 space-y-1 overflow-y-auto overscroll-contain px-2 py-3">
      {groups.map((group) => {
        const open = Boolean(openByKey[group.titleKey])
        const hasActive = group.titleKey === activeGroupKey
        return (
          <div key={group.titleKey}>
            <button
              type="button"
              className={cn(
                'flex w-full items-center gap-1 rounded px-2 py-1.5 text-left transition',
                'hover:bg-mist/60',
                hasActive ? 'text-cyan' : 'text-text-dim',
              )}
              aria-expanded={open}
              onClick={() => toggleGroup(group.titleKey)}
            >
              <ChevronDown
                className={cn(
                  'h-3 w-3 shrink-0 opacity-70 transition-transform duration-150',
                  open ? 'rotate-0' : '-rotate-90',
                )}
              />
              <span className="hud-label min-w-0 flex-1 truncate">{t(group.titleKey)}</span>
            </button>
            {open ? (
              <div className="mb-2 flex flex-col gap-0.5">
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
            ) : null}
          </div>
        )
      })}
    </nav>
  )
}

export function AppShell() {
  const { t } = useTranslation()
  const { user, roles, checkPermission, isViewerOnly, isAdmin } = useAuth()
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

  const skipEnterAnim = shouldSkipEnterAnim()
  const aiHome = location.pathname === '/ai' || location.pathname === '/'

  return (
    <div className="relative flex h-dvh w-full flex-col overflow-hidden pb-[env(safe-area-inset-bottom)] pt-[env(safe-area-inset-top)]">
      <header className="app-topbar">
        <div className="app-topbar-left">
          {/* Mobile nav — left of brand, only when sidebar is hidden (< md) */}
          {!aiHome ? (
            <button
              type="button"
              className="app-topbar-nav-toggle md:hidden"
              aria-label={mobileNavOpen ? 'Close navigation' : 'Open navigation'}
              aria-expanded={mobileNavOpen}
              onClick={() => setMobileNavOpen((v) => !v)}
            >
              {mobileNavOpen ? <X className="h-3.5 w-3.5" /> : <Menu className="h-3.5 w-3.5" />}
            </button>
          ) : null}
          <BrandMark
            to="/ai"
            className="hidden min-w-0 lg:inline-flex"
            brandClassName="text-sm md:text-base"
          />
          <BrandMark
            to="/ai"
            compact
            className="min-w-0 lg:hidden"
            brandClassName="text-xs tracking-[0.12em]"
          />

          {/* Expanded search only when the bar has real width */}
          {!aiHome ? (
            <div className="app-topbar-search-expanded hidden min-w-0 flex-1 xl:block">
              <GlobalSearchPalette />
            </div>
          ) : null}
        </div>

        <div className="app-topbar-center">
          {aiHome ? (
            <Link to="/fleet" className="app-topbar-link is-emphasis">
              <Terminal className="h-3.5 w-3.5" />
              <span className="hidden sm:inline">{t('nav.console')}</span>
            </Link>
          ) : (
            <Link to="/ai" className="app-topbar-link" title={t('nav.ai')}>
              <Sparkles className="h-3.5 w-3.5" />
              <span className="hidden xl:inline">{t('nav.ai')}</span>
            </Link>
          )}

          {/* >= lg: context lives in topbar; < lg: main strip / drawer */}
          <ClusterNamespaceControls
            layout="inline"
            showLabels
            className="hidden lg:flex"
          />
        </div>

        <div className="app-topbar-right">
          {!aiHome ? (
            <div className="app-topbar-search-compact shrink-0 xl:hidden">
              <GlobalSearchPalette />
            </div>
          ) : null}
          <span className="hidden shrink-0 sm:inline-flex">
            <ConnDot online />
          </span>
          <div className="shrink-0">
            <UserMenu primaryRole={primaryRole} isViewerOnly={isViewerOnly} />
          </div>
        </div>
      </header>

      {/* Mobile drawer — portal to body so topbar (z-shell) cannot bury it */}
      {typeof document !== 'undefined'
        ? createPortal(
            <AnimatePresence>
              {!aiHome && mobileNavOpen ? (
                <>
                  <motion.button
                    type="button"
                    key="nav-backdrop"
                    className="fixed inset-0 z-[var(--z-drawer-backdrop)] bg-black/60 md:hidden"
                    aria-label="Close navigation overlay"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.18 }}
                    onClick={() => setMobileNavOpen(false)}
                  />
                  <motion.aside
                    key="nav-drawer"
                    className="fixed inset-y-0 left-0 z-[var(--z-drawer)] flex w-[min(86vw,300px)] flex-col border-r border-line bg-panel-solid pt-[env(safe-area-inset-top)] shadow-[8px_0_32px_rgba(0,0,0,0.45)] md:hidden"
                    initial={{ x: '-100%' }}
                    animate={{ x: 0 }}
                    exit={{ x: '-100%' }}
                    transition={{ type: 'tween', duration: 0.22, ease: 'easeOut' }}
                  >
                    <div className="flex h-14 shrink-0 items-center justify-between border-b border-line px-3">
                      <BrandMark brandClassName="text-sm" />
                      <button
                        type="button"
                        className="inline-flex h-10 w-10 items-center justify-center rounded border border-line text-cyan"
                        aria-label="Close navigation"
                        onClick={() => setMobileNavOpen(false)}
                      >
                        <X className="h-5 w-5" />
                      </button>
                    </div>

                    <div className="shrink-0 border-b border-line p-3">
                      <ClusterNamespaceControls layout="stack" showLabels />
                    </div>

                    <NavBody groups={visibleGroups} onNavigate={() => setMobileNavOpen(false)} />
                  </motion.aside>
                </>
              ) : null}
            </AnimatePresence>,
            document.body,
          )
        : null}

      <div className="relative flex min-h-0 w-full flex-1">
        {/* Desktop sidebar — hidden on AI home (console entry lives in the AI page) */}
        {!aiHome ? (
          <aside className="hidden h-full w-[232px] shrink-0 flex-col border-r border-line bg-panel-solid/50 md:flex">
            <NavBody groups={visibleGroups} />
          </aside>
        ) : null}

        <main
          className={cn(
            'relative flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden',
            aiHome
              ? 'px-0 py-0'
              : 'px-2.5 py-2.5 sm:px-4 sm:py-4 md:px-6 md:py-5 xl:px-8',
          )}
        >
          {!aiHome ? (
            <ClusterNamespaceControls
              layout="grid"
              showLabels
              className="mb-2.5 shrink-0 lg:hidden"
            />
          ) : null}

          {!aiHome ? (
            <div className="shrink-0">
              <OAuthAccountBanner />
            </div>
          ) : null}

          <motion.div
            key={location.pathname}
            className={cn(
              'flex h-full min-h-0 w-full min-w-0 flex-1 flex-col',
              aiHome ? 'overflow-hidden' : 'overflow-y-auto overscroll-contain',
            )}
            initial={skipEnterAnim ? false : { opacity: 0, y: 6 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: skipEnterAnim ? 0 : 0.2, ease: 'easeOut' }}
          >
            <Outlet />
          </motion.div>
        </main>
      </div>
      <StarSupportFloat />
    </div>
  )
}
