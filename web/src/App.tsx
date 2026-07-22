import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider, useAuth } from '@/store/auth'
import { ClusterProvider } from '@/store/cluster'
import { NamespaceProvider } from '@/store/namespace'
import { AppShell } from '@/components/AppShell'
import { LoginPage } from '@/pages/LoginPage'
import { OverviewPage } from '@/pages/OverviewPage'
import { NodesPage } from '@/pages/NodesPage'
import { EventsPage } from '@/pages/EventsPage'
import { MonitoringPage } from '@/pages/MonitoringPage'
import {
  ClusterRoleBindingsPage,
  ClusterRolesPage,
  ConfigMapsPage,
  DeploymentsPage,
  IngressPage,
  NamespacesPage,
  PodsPage,
  PVPage,
  PVCPage,
  RoleBindingsPage,
  RolesPage,
  SecretsPage,
  ServicesPage,
  StorageClassPage,
} from '@/pages/resources'
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
                    <NamespaceProvider>
                      <AppShell />
                    </NamespaceProvider>
                  </ClusterProvider>
                </Protected>
              }
            >
              <Route index element={<OverviewPage />} />
              <Route path="nodes" element={<NodesPage />} />
              <Route path="namespaces" element={<NamespacesPage />} />
              <Route path="events" element={<EventsPage />} />
              <Route path="pods" element={<PodsPage />} />
              <Route path="workloads" element={<Navigate to="/pods" replace />} />
              <Route path="deployments" element={<DeploymentsPage />} />
              <Route path="services" element={<ServicesPage />} />
              <Route path="ingresses" element={<IngressPage />} />
              <Route path="configmaps" element={<ConfigMapsPage />} />
              <Route path="secrets" element={<SecretsPage />} />
              <Route path="persistentvolumes" element={<PVPage />} />
              <Route path="persistentvolumeclaims" element={<PVCPage />} />
              <Route path="storageclasses" element={<StorageClassPage />} />
              <Route path="roles" element={<RolesPage />} />
              <Route path="rolebindings" element={<RoleBindingsPage />} />
              <Route path="clusterroles" element={<ClusterRolesPage />} />
              <Route path="clusterrolebindings" element={<ClusterRoleBindingsPage />} />
              <Route path="monitoring" element={<MonitoringPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  )
}
