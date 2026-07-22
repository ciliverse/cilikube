import { NavLink, Outlet, useLocation } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Activity,
  Boxes,
  Database,
  HardDrive,
  KeyRound,
  LayoutDashboard,
  Lock,
  LogOut,
  Network,
  Server,
  Shield,
  Layers,
} from 'lucide-react'
import { useAuth } from '@/store/auth'
import { useCluster } from '@/store/cluster'
import { useNamespace } from '@/store/namespace'
import { Button, ConnDot } from './ui'
import { cn } from '@/lib/utils'

type NavItem = { to: string; label: string; icon: typeof Server; namespaced?: boolean }
type NavGroup = { title: string; items: NavItem[] }

const navGroups: NavGroup[] = [
  {
    title: 'Cluster',
    items: [
      { to: '/', label: 'Overview', icon: LayoutDashboard },
      { to: '/nodes', label: 'Nodes', icon: Server },
      { to: '/namespaces', label: 'Namespaces', icon: Layers },
      { to: '/events', label: 'Events', icon: Activity },
    ],
  },
  {
    title: 'Workloads',
    items: [
      { to: '/pods', label: 'Pods', icon: Boxes, namespaced: true },
      { to: '/deployments', label: 'Deployments', icon: Boxes, namespaced: true },
    ],
  },
  {
    title: 'Network',
    items: [
      { to: '/services', label: 'Services', icon: Network, namespaced: true },
      { to: '/ingresses', label: 'Ingress', icon: Network, namespaced: true },
    ],
  },
  {
    title: 'Config',
    items: [
      { to: '/configmaps', label: 'ConfigMaps', icon: Database, namespaced: true },
      { to: '/secrets', label: 'Secrets', icon: KeyRound, namespaced: true },
    ],
  },
  {
    title: 'Storage',
    items: [
      { to: '/persistentvolumes', label: 'PV', icon: HardDrive },
      { to: '/persistentvolumeclaims', label: 'PVC', icon: HardDrive, namespaced: true },
      { to: '/storageclasses', label: 'StorageClass', icon: HardDrive },
    ],
  },
  {
    title: 'Access',
    items: [
      { to: '/roles', label: 'Roles', icon: Lock, namespaced: true },
      { to: '/rolebindings', label: 'RoleBindings', icon: Lock, namespaced: true },
      { to: '/clusterroles', label: 'ClusterRoles', icon: Shield },
      { to: '/clusterrolebindings', label: 'ClusterRoleBindings', icon: Shield },
    ],
  },
  {
    title: 'Observe',
    items: [{ to: '/monitoring', label: 'Monitoring', icon: Activity }],
  },
]

const namespacedPaths = new Set(
  navGroups.flatMap((g) => g.items.filter((i) => i.namespaced).map((i) => i.to)),
)

export function AppShell() {
  const { user, logout } = useAuth()
  const { clusters, clusterId, setClusterId, activeCluster } = useCluster()
  const { namespaces, namespace, setNamespace } = useNamespace()
  const location = useLocation()
  const showNamespace = namespacedPaths.has(location.pathname)

  return (
    <div className="relative flex min-h-screen flex-col">
      <header className="sticky top-0 z-20 flex h-14 items-center gap-4 border-b border-line bg-panel-solid/95 px-4 backdrop-blur md:px-5">
        <div className="hud-brand shrink-0 text-sm md:text-base">
          CILI<span className="accent">KUBE</span>
        </div>

        <div className="hidden items-center gap-2 border-l border-line pl-4 text-xs tracking-[0.16em] text-text-dim uppercase lg:flex">
          Control plane
        </div>

        <div className="ml-auto flex min-w-0 items-center gap-3">
          <label className="hidden items-center gap-2 sm:flex">
            <span className="hud-label">Cluster</span>
            <select
              value={clusterId}
              onChange={(e) => setClusterId(e.target.value)}
              className="hud-select w-auto min-w-[140px] max-w-[220px]"
            >
              {clusters.map((c) => (
                <option key={c.id || c.name} value={c.id || c.name}>
                  {c.name || c.id}
                </option>
              ))}
            </select>
          </label>

          {showNamespace ? (
            <label className="hidden items-center gap-2 md:flex">
              <span className="hud-label">NS</span>
              <select
                value={namespace}
                onChange={(e) => setNamespace(e.target.value)}
                className="hud-select w-auto min-w-[120px] max-w-[180px]"
              >
                {namespaces.map((ns) => (
                  <option key={ns} value={ns}>
                    {ns}
                  </option>
                ))}
              </select>
            </label>
          ) : null}

          <div className="hidden items-center gap-2 text-xs text-text-dim xl:flex">
            <ConnDot online />
            <span className="max-w-[160px] truncate tracking-wider uppercase">
              {activeCluster?.name || clusterId || 'no-cluster'}
            </span>
          </div>

          <div className="hidden text-xs text-text-dim md:block">
            <span className="text-text">{user?.username || 'operator'}</span>
            <span className="mx-1.5 text-line">/</span>
            <span>{user?.role || 'admin'}</span>
          </div>

          <Button variant="ghost" className="px-2" onClick={() => void logout()} title="Logout">
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
      </header>

      <div className="mx-auto flex w-full max-w-[1600px] flex-1">
        <aside className="sticky top-14 flex h-[calc(100vh-3.5rem)] w-[220px] shrink-0 flex-col border-r border-line bg-panel-solid/40">
          <nav className="flex-1 space-y-4 overflow-y-auto px-2 py-4">
            {navGroups.map((group) => (
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

          <div className="border-t border-line p-3 sm:hidden">
            <div className="hud-label mb-2">Cluster</div>
            <select
              value={clusterId}
              onChange={(e) => setClusterId(e.target.value)}
              className="hud-select"
            >
              {clusters.map((c) => (
                <option key={c.id || c.name} value={c.id || c.name}>
                  {c.name || c.id}
                </option>
              ))}
            </select>
          </div>
        </aside>

        <main className="relative min-w-0 flex-1 px-4 py-5 md:px-7 md:py-6">
          {showNamespace ? (
            <div className="mb-4 flex items-center gap-2 md:hidden">
              <span className="hud-label">Namespace</span>
              <select
                value={namespace}
                onChange={(e) => setNamespace(e.target.value)}
                className="hud-select"
              >
                {namespaces.map((ns) => (
                  <option key={ns} value={ns}>
                    {ns}
                  </option>
                ))}
              </select>
            </div>
          ) : null}

          <motion.div
            key={location.pathname}
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
