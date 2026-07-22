import { Badge } from '@/components/ui'
import { ResourceListPage, createdCell } from '@/components/ResourceListPage'
import { listClusterResource, listNamespacedResource, metaName, metaNamespace } from '@/api/resources'
import { useNamespace } from '@/store/namespace'
function ReadyBadge({ ok, okText = 'Ready', badText = 'NotReady' }: { ok: boolean; okText?: string; badText?: string }) {
  return <Badge tone={ok ? 'ok' : 'danger'}>{ok ? okText : badText}</Badge>
}

export function NamespacesPage() {
  return (
    <ResourceListPage
      title="NAMESPACES"
      subtitle="Cluster namespaces"
      resourceKey="namespaces"
      queryFn={() => listClusterResource('namespaces')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        {
          key: 'status',
          header: 'Status',
          render: (item) => <Badge tone="accent">{item.status?.phase || '-'}</Badge>,
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function PodsPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="PODS"
      namespaced
      resourceKey="pods"
      queryFn={() => listNamespacedResource(namespace, 'pods')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'status',
          header: 'Status',
          render: (item) => {
            const phase = item.status?.phase || 'Unknown'
            const tone =
              phase === 'Running' ? 'ok' : phase === 'Pending' ? 'warn' : phase === 'Failed' ? 'danger' : 'neutral'
            return <Badge tone={tone}>{phase}</Badge>
          },
        },
        {
          key: 'restarts',
          header: 'Restarts',
          render: (item) =>
            (item.status?.containerStatuses || []).reduce(
              (sum: number, c: any) => sum + (c.restartCount || 0),
              0,
            ),
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function DeploymentsPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="DEPLOYMENTS"
      namespaced
      resourceKey="deployments"
      queryFn={() => listNamespacedResource(namespace, 'deployments')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'ready',
          header: 'Ready',
          render: (item) => {
            const ready = item.status?.readyReplicas ?? 0
            const desired = item.spec?.replicas ?? 0
            return (
              <span className="font-mono text-xs">
                {ready}/{desired}
              </span>
            )
          },
        },
        {
          key: 'up_to_date',
          header: 'Up-to-date',
          render: (item) => item.status?.updatedReplicas ?? 0,
        },
        {
          key: 'available',
          header: 'Available',
          render: (item) => item.status?.availableReplicas ?? 0,
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function ServicesPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="SERVICES"
      namespaced
      resourceKey="services"
      queryFn={() => listNamespacedResource(namespace, 'services')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'type',
          header: 'Type',
          render: (item) => <Badge tone="accent">{item.spec?.type || 'ClusterIP'}</Badge>,
        },
        {
          key: 'clusterip',
          header: 'Cluster IP',
          render: (item) => <span className="font-mono text-xs">{item.spec?.clusterIP || '-'}</span>,
        },
        {
          key: 'ports',
          header: 'Ports',
          render: (item) =>
            (item.spec?.ports || [])
              .map((p: any) => `${p.port}${p.nodePort ? `:${p.nodePort}` : ''}/${p.protocol || 'TCP'}`)
              .join(', ') || '-',
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function IngressPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="INGRESS"
      namespaced
      resourceKey="ingresses"
      queryFn={() => listNamespacedResource(namespace, 'ingresses')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'class',
          header: 'Class',
          render: (item) => item.spec?.ingressClassName || item.metadata?.annotations?.['kubernetes.io/ingress.class'] || '-',
        },
        {
          key: 'hosts',
          header: 'Hosts',
          render: (item) =>
            (item.spec?.rules || []).map((r: any) => r.host).filter(Boolean).join(', ') || '-',
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function ConfigMapsPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="CONFIGMAPS"
      namespaced
      resourceKey="configmaps"
      queryFn={() => listNamespacedResource(namespace, 'configmaps')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'keys',
          header: 'Data keys',
          render: (item) => Object.keys(item.data || {}).length,
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function SecretsPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="SECRETS"
      namespaced
      resourceKey="secrets"
      queryFn={() => listNamespacedResource(namespace, 'secrets')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'type',
          header: 'Type',
          render: (item) => <span className="font-mono text-xs">{item.type || '-'}</span>,
        },
        {
          key: 'keys',
          header: 'Data keys',
          render: (item) => Object.keys(item.data || {}).length,
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function PVPage() {
  return (
    <ResourceListPage
      title="PERSISTENT VOLUMES"
      resourceKey="persistentvolumes"
      queryFn={() => listClusterResource('persistentvolumes')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        {
          key: 'capacity',
          header: 'Capacity',
          render: (item) => item.spec?.capacity?.storage || '-',
        },
        {
          key: 'access',
          header: 'Access',
          render: (item) => (item.spec?.accessModes || []).join(', ') || '-',
        },
        {
          key: 'status',
          header: 'Status',
          render: (item) => <Badge tone="accent">{item.status?.phase || '-'}</Badge>,
        },
        {
          key: 'claim',
          header: 'Claim',
          render: (item) =>
            item.spec?.claimRef
              ? `${item.spec.claimRef.namespace}/${item.spec.claimRef.name}`
              : '-',
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function PVCPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="PERSISTENT VOLUME CLAIMS"
      namespaced
      resourceKey="persistentvolumeclaims"
      queryFn={() => listNamespacedResource(namespace, 'persistentvolumeclaims')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'status',
          header: 'Status',
          render: (item) => {
            const phase = item.status?.phase || '-'
            return <Badge tone={phase === 'Bound' ? 'ok' : 'warn'}>{phase}</Badge>
          },
        },
        {
          key: 'volume',
          header: 'Volume',
          render: (item) => item.spec?.volumeName || '-',
        },
        {
          key: 'capacity',
          header: 'Capacity',
          render: (item) => item.status?.capacity?.storage || item.spec?.resources?.requests?.storage || '-',
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function StorageClassPage() {
  return (
    <ResourceListPage
      title="STORAGE CLASSES"
      resourceKey="storageclasses"
      queryFn={() => listClusterResource('storageclasses')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        {
          key: 'provisioner',
          header: 'Provisioner',
          render: (item) => <span className="font-mono text-xs">{item.provisioner || '-'}</span>,
        },
        {
          key: 'reclaim',
          header: 'Reclaim',
          render: (item) => item.reclaimPolicy || '-',
        },
        {
          key: 'binding',
          header: 'Binding',
          render: (item) => item.volumeBindingMode || '-',
        },
        {
          key: 'default',
          header: 'Default',
          render: (item) =>
            item.metadata?.annotations?.['storageclass.kubernetes.io/is-default-class'] === 'true' ? (
              <ReadyBadge ok okText="Yes" badText="No" />
            ) : (
              <Badge tone="neutral">No</Badge>
            ),
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function RolesPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="ROLES"
      namespaced
      resourceKey="roles"
      queryFn={() => listNamespacedResource(namespace, 'roles')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'rules',
          header: 'Rules',
          render: (item) => (item.rules || []).length,
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function RoleBindingsPage() {
  const { namespace } = useNamespace()
  return (
    <ResourceListPage
      title="ROLE BINDINGS"
      namespaced
      resourceKey="rolebindings"
      queryFn={() => listNamespacedResource(namespace, 'rolebindings')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'role',
          header: 'Role',
          render: (item) =>
            item.roleRef
              ? `${item.roleRef.kind}/${item.roleRef.name}`
              : '-',
        },
        {
          key: 'subjects',
          header: 'Subjects',
          render: (item) => (item.subjects || []).length,
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function ClusterRolesPage() {
  return (
    <ResourceListPage
      title="CLUSTER ROLES"
      resourceKey="clusterroles"
      queryFn={() => listClusterResource('clusterroles')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        {
          key: 'rules',
          header: 'Rules',
          render: (item) => (item.rules || []).length,
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}

export function ClusterRoleBindingsPage() {
  return (
    <ResourceListPage
      title="CLUSTER ROLE BINDINGS"
      resourceKey="clusterrolebindings"
      queryFn={() => listClusterResource('clusterrolebindings')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => <span className="font-semibold text-cyan">{metaName(item)}</span>,
        },
        {
          key: 'role',
          header: 'Role',
          render: (item) =>
            item.roleRef
              ? `${item.roleRef.kind}/${item.roleRef.name}`
              : '-',
        },
        {
          key: 'subjects',
          header: 'Subjects',
          render: (item) => (item.subjects || []).length,
        },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
    />
  )
}
