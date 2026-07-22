import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider, useAuth } from '@/store/auth'
import { ClusterProvider } from '@/store/cluster'
import { AppShell } from '@/components/AppShell'
import { LoginPage } from '@/pages/LoginPage'
import { OverviewPage } from '@/pages/OverviewPage'
import { NodesPage } from '@/pages/NodesPage'
import { WorkloadsPage } from '@/pages/WorkloadsPage'
import { EventsPage } from '@/pages/EventsPage'
import { MonitoringPage } from '@/pages/MonitoringPage'
import type { ReactNode } from 'react'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

function Protected({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth()
  if (!isAuthenticated) return <Navigate to="/login" replace />
  return children
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/"
              element={
                <Protected>
                  <ClusterProvider>
                    <AppShell />
                  </ClusterProvider>
                </Protected>
              }
            >
              <Route index element={<OverviewPage />} />
              <Route path="nodes" element={<NodesPage />} />
              <Route path="workloads" element={<WorkloadsPage />} />
              <Route path="events" element={<EventsPage />} />
              <Route path="monitoring" element={<MonitoringPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  )
}
