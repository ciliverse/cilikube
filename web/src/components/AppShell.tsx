import { NavLink, Outlet } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Activity,
  Boxes,
  Cpu,
  Hexagon,
  LayoutDashboard,
  LogOut,
  Network,
  Server,
} from 'lucide-react'
import { useAuth } from '@/store/auth'
import { useCluster } from '@/store/cluster'
import { Button } from './ui'
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
      <div className="mx-auto flex min-h-screen max-w-[1500px]">
        <aside className="sticky top-0 flex h-screen w-[260px] shrink-0 flex-col border-r border-line/70 bg-white/55 px-4 py-5 backdrop-blur-xl">
          <div className="mb-8 flex items-center gap-3 px-2">
            <div className="grid h-11 w-11 place-items-center rounded-2xl bg-accent text-white shadow-lg shadow-accent/30">
              <Hexagon className="h-6 w-6" />
            </div>
            <div>
              <div className="font-display text-xl font-extrabold tracking-tight">CiliKube</div>
              <div className="text-xs text-ink-soft">Control Plane</div>
            </div>
          </div>

          <nav className="flex flex-1 flex-col gap-1">
            {nav.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.to === '/'}
                className={({ isActive }) =>
                  cn(
                    'group flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-semibold transition',
                    isActive
                      ? 'bg-accent text-white shadow-md shadow-accent/25'
                      : 'text-ink-soft hover:bg-black/5 hover:text-ink',
                  )
                }
              >
                <item.icon className="h-4 w-4 opacity-90" />
                {item.label}
              </NavLink>
            ))}
          </nav>

          <div className="mt-4 space-y-3 rounded-2xl border border-line bg-mist/70 p-3">
            <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-ink-soft">
              <Cpu className="h-3.5 w-3.5" />
              Cluster
            </div>
            <select
              value={clusterId}
              onChange={(e) => setClusterId(e.target.value)}
              className="w-full rounded-xl border border-line bg-white px-3 py-2 text-sm outline-none focus:border-accent"
            >
              {clusters.map((c) => (
                <option key={c.id || c.name} value={c.id || c.name}>
                  {c.name || c.id}
                </option>
              ))}
            </select>
            <div className="text-xs text-ink-soft">
              {activeCluster?.status || activeCluster?.version || 'Ready'}
            </div>
          </div>

          <div className="mt-4 flex items-center justify-between gap-2 px-1">
            <div className="min-w-0">
              <div className="truncate text-sm font-semibold">{user?.username || 'operator'}</div>
              <div className="truncate text-xs text-ink-soft">{user?.role || 'admin'}</div>
            </div>
            <Button variant="ghost" className="px-2" onClick={() => void logout()} title="Logout">
              <LogOut className="h-4 w-4" />
            </Button>
          </div>
        </aside>

        <main className="relative flex-1 px-5 py-6 md:px-8">
          <motion.div
            initial={{ opacity: 0, y: 10 }}
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
