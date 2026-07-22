import { NavLink, Outlet } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Activity,
  Boxes,
  LayoutDashboard,
  LogOut,
  Network,
  Server,
} from 'lucide-react'
import { useAuth } from '@/store/auth'
import { useCluster } from '@/store/cluster'
import { Button, ConnDot } from './ui'
import { cn } from '@/lib/utils'

const nav = [
  { to: '/', label: 'Overview', icon: LayoutDashboard },
  { to: '/nodes', label: 'Nodes', icon: Server },
  { to: '/workloads', label: 'Workloads', icon: Boxes },
  { to: '/events', label: 'Events', icon: Network },
  { to: '/monitoring', label: 'Monitoring', icon: Activity },
]

export function AppShell() {
  const { user, logout } = useAuth()
  const { clusters, clusterId, setClusterId, activeCluster } = useCluster()

  return (
    <div className="relative min-h-screen">
      <header className="sticky top-0 z-20 flex items-center gap-4 border-b border-line bg-panel-solid/90 px-4 py-3 backdrop-blur md:px-6">
        <div className="hud-brand text-sm md:text-base">
          CILI<span className="accent">KUBE</span>
        </div>
        <div className="hidden text-xs tracking-[0.2em] text-text-dim uppercase sm:block">
          Control plane
        </div>
        <div className="ml-auto flex items-center gap-4">
          <div className="hidden items-center gap-2 text-xs text-text-dim md:flex">
            <ConnDot online />
            <span className="tracking-wider uppercase">
              {activeCluster?.name || clusterId || 'no-cluster'}
            </span>
          </div>
          <div className="text-xs text-text-dim">
            <span className="text-text">{user?.username || 'operator'}</span>
            <span className="mx-2 text-line">/</span>
            <span>{user?.role || 'admin'}</span>
          </div>
          <Button variant="ghost" className="px-2" onClick={() => void logout()} title="Logout">
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
      </header>

      <div className="mx-auto flex min-h-[calc(100vh-57px)] max-w-[1500px]">
        <aside className="sticky top-[57px] flex h-[calc(100vh-57px)] w-[240px] shrink-0 flex-col border-r border-line bg-panel-solid/50 px-3 py-4">
          <nav className="flex flex-1 flex-col gap-1">
            {nav.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.to === '/'}
                className={({ isActive }) =>
                  cn(
                    'group flex items-center gap-3 rounded px-3 py-2.5 text-sm font-semibold tracking-wide transition',
                    isActive
                      ? 'border border-cyan/40 bg-cyan/15 text-cyan shadow-[0_0_16px_rgba(53,230,255,0.12)]'
                      : 'border border-transparent text-text-dim hover:border-line hover:bg-mist hover:text-text',
                  )
                }
              >
                <item.icon className="h-4 w-4 opacity-90" />
                {item.label}
              </NavLink>
            ))}
          </nav>

          <div className="hud-panel mt-4 space-y-3 rounded p-3">
            <div className="hud-label">Cluster</div>
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
            <div className="text-xs text-text-dim">
              {activeCluster?.status || activeCluster?.version || 'Ready'}
            </div>
          </div>
        </aside>

        <main className="relative flex-1 px-5 py-6 md:px-8">
          <motion.div
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.35, ease: 'easeOut' }}
          >
            <Outlet />
          </motion.div>
        </main>
      </div>
    </div>
  )
}
