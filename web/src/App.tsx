import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider, useAuth } from '@/store/auth'
import { ClusterProvider } from '@/store/cluster'
import { NamespaceProvider } from '@/store/namespace'
import { AppShell } from '@/components/AppShell'
import { LoginPage } from '@/pages/LoginPage'
import { OAuthCallbackPage } from '@/pages/OAuthCallbackPage'
import { OverviewPage } from '@/pages/OverviewPage'
import { NodesPage } from '@/pages/NodesPage'
import { EventsPage } from '@/pages/EventsPage'
import { MonitoringPage } from '@/pages/MonitoringPage'
import { ResourceDetailPage } from '@/pages/ResourceDetailPage'
import { ClustersPage } from '@/pages/ClustersPage'
import { AdminUsersPage } from '@/pages/AdminUsersPage'
import { AdminRolesPage } from '@/pages/AdminRolesPage'
import { CrdsPage } from '@/pages/CrdsPage'
import { AuditPage } from '@/pages/AuditPage'
import { GlobalSearchPage } from '@/pages/GlobalSearchPage'
import { HelmPage } from '@/pages/HelmPage'
import { ProxyConsolePage } from '@/pages/ProxyConsolePage'
import { ProfilePage } from '@/pages/ProfilePage'
import { AdminSettingsPage } from '@/pages/AdminSettingsPage'
import {
  ClusterRoleBindingsPage,
  ClusterRolesPage,
  ConfigMapsPage,
  CronJobsPage,
  DaemonSetsPage,
  DeploymentsPage,
  HPAPage,
  IngressPage,
  JobsPage,
  LimitRangesPage,
  NamespacesPage,
  NetworkPoliciesPage,
  PDBPage,
  PodsPage,
  PVPage,
  PVCPage,
  ResourceQuotasPage,
  RoleBindingsPage,
  RolesPage,
  SecretsPage,
  ServiceAccountsPage,
  ServicesPage,
  StatefulSetsPage,
  StorageClassPage,
} from '@/pages/resources'
import { useState, type ReactNode } from 'react'
import { BootScreen, shouldShowBootScreen } from '@/components/BootScreen'

function Detail({ resource, namespaced = true }: { resource: string; namespaced?: boolean }) {
  return <ResourceDetailPage resource={resource} namespaced={namespaced} />
}

function BootGate({ children }: { children: ReactNode }) {
  const [booting, setBooting] = useState(() => shouldShowBootScreen())
  if (booting) return <BootScreen onDone={() => setBooting(false)} />
  return children
}

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
        <BootGate>
          <BrowserRouter>
            <Routes>
              <Route path="/login" element={<LoginPage />} />
              <Route path="/login/oauth/callback" element={<OAuthCallbackPage />} />
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
              <Route path="pods/:namespace/:name" element={<Detail resource="pods" />} />
              <Route path="workloads" element={<Navigate to="/pods" replace />} />
              <Route path="deployments" element={<DeploymentsPage />} />
              <Route path="deployments/:namespace/:name" element={<Detail resource="deployments" />} />
              <Route path="statefulsets" element={<StatefulSetsPage />} />
              <Route
                path="statefulsets/:namespace/:name"
                element={<Detail resource="statefulsets" />}
              />
              <Route path="daemonsets" element={<DaemonSetsPage />} />
              <Route path="daemonsets/:namespace/:name" element={<Detail resource="daemonsets" />} />
              <Route path="jobs" element={<JobsPage />} />
              <Route path="jobs/:namespace/:name" element={<Detail resource="jobs" />} />
              <Route path="cronjobs" element={<CronJobsPage />} />
              <Route path="cronjobs/:namespace/:name" element={<Detail resource="cronjobs" />} />
              <Route path="horizontalpodautoscalers" element={<HPAPage />} />
              <Route
                path="horizontalpodautoscalers/:namespace/:name"
                element={<Detail resource="horizontalpodautoscalers" />}
              />
              <Route path="poddisruptionbudgets" element={<PDBPage />} />
              <Route
                path="poddisruptionbudgets/:namespace/:name"
                element={<Detail resource="poddisruptionbudgets" />}
              />
              <Route path="services" element={<ServicesPage />} />
              <Route path="services/:namespace/:name" element={<Detail resource="services" />} />
              <Route path="ingresses" element={<IngressPage />} />
              <Route path="ingresses/:namespace/:name" element={<Detail resource="ingresses" />} />
              <Route path="networkpolicies" element={<NetworkPoliciesPage />} />
              <Route
                path="networkpolicies/:namespace/:name"
                element={<Detail resource="networkpolicies" />}
              />
              <Route path="configmaps" element={<ConfigMapsPage />} />
              <Route path="configmaps/:namespace/:name" element={<Detail resource="configmaps" />} />
              <Route path="secrets" element={<SecretsPage />} />
              <Route path="secrets/:namespace/:name" element={<Detail resource="secrets" />} />
              <Route path="serviceaccounts" element={<ServiceAccountsPage />} />
              <Route
                path="serviceaccounts/:namespace/:name"
                element={<Detail resource="serviceaccounts" />}
              />
              <Route path="resourcequotas" element={<ResourceQuotasPage />} />
              <Route
                path="resourcequotas/:namespace/:name"
                element={<Detail resource="resourcequotas" />}
              />
              <Route path="limitranges" element={<LimitRangesPage />} />
              <Route
                path="limitranges/:namespace/:name"
                element={<Detail resource="limitranges" />}
              />
              <Route path="persistentvolumes" element={<PVPage />} />
              <Route
                path="persistentvolumes/:name"
                element={<Detail resource="persistentvolumes" namespaced={false} />}
              />
              <Route path="persistentvolumeclaims" element={<PVCPage />} />
              <Route
                path="persistentvolumeclaims/:namespace/:name"
                element={<Detail resource="persistentvolumeclaims" />}
              />
              <Route path="storageclasses" element={<StorageClassPage />} />
              <Route
                path="storageclasses/:name"
                element={<Detail resource="storageclasses" namespaced={false} />}
              />
              <Route path="roles" element={<RolesPage />} />
              <Route path="roles/:namespace/:name" element={<Detail resource="roles" />} />
              <Route path="rolebindings" element={<RoleBindingsPage />} />
              <Route
                path="rolebindings/:namespace/:name"
                element={<Detail resource="rolebindings" />}
              />
              <Route path="clusterroles" element={<ClusterRolesPage />} />
              <Route
                path="clusterroles/:name"
                element={<Detail resource="clusterroles" namespaced={false} />}
              />
              <Route path="clusterrolebindings" element={<ClusterRoleBindingsPage />} />
              <Route
                path="clusterrolebindings/:name"
                element={<Detail resource="clusterrolebindings" namespaced={false} />}
              />
              <Route path="nodes/:name" element={<Detail resource="nodes" namespaced={false} />} />
              <Route path="monitoring" element={<MonitoringPage />} />
              <Route path="clusters" element={<ClustersPage />} />
              <Route path="admin/users" element={<AdminUsersPage />} />
              <Route path="admin/roles" element={<AdminRolesPage />} />
              <Route path="admin/settings" element={<AdminSettingsPage />} />
              <Route path="profile" element={<ProfilePage />} />
              <Route path="crds" element={<CrdsPage />} />
              <Route path="audit" element={<AuditPage />} />
              <Route path="search" element={<GlobalSearchPage />} />
              <Route path="helm" element={<HelmPage />} />
              <Route path="proxy" element={<ProxyConsolePage />} />
              </Route>
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </BrowserRouter>
        </BootGate>
      </AuthProvider>
    </QueryClientProvider>
  )
}
