import { useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { useQueryClient } from '@tanstack/react-query'
import { Network, RefreshCw, ScrollText, TerminalSquare } from 'lucide-react'
import { Badge, Button } from '@/components/ui'
import { ResourceListPage, ageCell, createdCell } from '@/components/ResourceListPage'
import { podMetricsKey, usePodMetricsMap } from '@/hooks/usePodMetricsMap'
import { podMetricColumns } from '@/components/PodMetricCells'
import { PodWorkbench } from '@/components/PodWorkbench'
import { PodPortForward } from '@/components/PodPortForward'
import { ResourceYamlModal } from '@/components/ResourceYamlModal'
import { ScaleDialog } from '@/components/ScaleDialog'
import { YamlActionButton } from '@/components/YamlActionButton'
import { ConfirmDialog } from '@/components/ConfirmDialog'
import {
  getNamespacedResource,
  listClusterResource,
  listNamespacedResource,
  metaName,
  metaNamespace,
  updateNamespacedResource,
} from '@/api/resources'
import { apiPatch } from '@/lib/api'
import { useNamespace } from '@/store/namespace'
import { useCluster } from '@/store/cluster'
import { useAuth } from '@/store/auth'
import { useYamlModal } from '@/hooks/useYamlModal'

function ReadyBadge({ ok, okText = 'Ready', badText = 'NotReady' }: { ok: boolean; okText?: string; badText?: string }) {
  return <Badge tone={ok ? 'ok' : 'danger'}>{ok ? okText : badText}</Badge>
}

function nameLink(resource: string, item: any, namespaced = true) {
  const name = metaName(item)
  const ns = metaNamespace(item)
  const to = namespaced ? `/${resource}/${ns}/${name}` : `/${resource}/${name}`
  return (
    <Link className="font-semibold text-cyan hover:underline" to={to} title={name}>
      {name}
    </Link>
  )
}

export function NamespacesPage() {
  const { yamlButton, yamlModal } = useYamlModal('namespaces', false)
  return (
    <>
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function PodsPage() {
  const { namespace } = useNamespace()
  const { clusterId } = useCluster()
  const { canExec } = useAuth()
  const queryClient = useQueryClient()
  const [target, setTarget] = useState<any | null>(null)
  const [tab, setTab] = useState<'logs' | 'exec' | 'attach' | 'yaml'>('logs')
  const [pfTarget, setPfTarget] = useState<any | null>(null)
  const metrics = usePodMetricsMap(namespace || undefined)

  return (
    <>
      <ResourceListPage
        title="PODS"
        namespaced
        resourceKey="pods"
        pinFirstColumn
        subtitle={
          metrics.available
            ? 'CPU/MEM + % of request/limit (metrics-server, ~15s)'
            : metrics.message || 'metrics-server unavailable — metrics show as -'
        }
        queryFn={() => listNamespacedResource(namespace, 'pods')}
        columns={[
          {
            key: 'name',
            header: 'Name',
            render: (item) => (
              <Link
                className="font-semibold text-cyan hover:underline"
                to={`/pods/${metaNamespace(item)}/${metaName(item)}`}
              >
                {metaName(item)}
              </Link>
            ),
          },
          { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
          {
            key: 'status',
            header: 'Status',
            render: (item) => {
              const phase = item.status?.phase || 'Unknown'
              const tone =
                phase === 'Running'
                  ? 'ok'
                  : phase === 'Pending'
                    ? 'warn'
                    : phase === 'Failed'
                      ? 'danger'
                      : 'neutral'
              return <Badge tone={tone}>{phase}</Badge>
            },
          },
          ...podMetricColumns((item) =>
            metrics.map.get(podMetricsKey(metaNamespace(item), metaName(item))),
          ),
          {
            key: 'restarts',
            header: 'Restarts',
            render: (item) =>
              (item.status?.containerStatuses || []).reduce(
                (sum: number, c: any) => sum + (c.restartCount || 0),
                0,
              ),
          },
          { key: 'age', header: 'Age', render: ageCell },
          { key: 'created', header: 'Created', render: createdCell },
        ]}
        actions={(item) => (
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              className="px-2 py-1"
              title="Logs"
              type="button"
              onClick={() => {
                setTab('logs')
                setTarget(item)
              }}
            >
              <ScrollText className="h-3.5 w-3.5" />
            </Button>
            {canExec ? (
              <>
                <Button
                  variant="ghost"
                  className="px-2 py-1"
                  title="Exec"
                  type="button"
                  onClick={() => {
                    setTab('exec')
                    setTarget(item)
                  }}
                >
                  <TerminalSquare className="h-3.5 w-3.5" />
                </Button>
                <Button
                  variant="ghost"
                  className="px-2 py-1"
                  title="Port Forward"
                  type="button"
                  onClick={() => setPfTarget(item)}
                >
                  <Network className="h-3.5 w-3.5" />
                </Button>
              </>
            ) : null}
            <YamlActionButton
              resource="pods"
              onClick={() => {
                setTab('yaml')
                setTarget(item)
              }}
            />
          </div>
        )}
      />
      <PodWorkbench
        open={Boolean(target)}
        pod={target}
        initialTab={tab}
        onClose={() => setTarget(null)}
        onDeleted={() => {
          void queryClient.invalidateQueries({ queryKey: ['pods', clusterId, namespace] })
        }}
      />
      <PodPortForward
        open={Boolean(pfTarget)}
        namespace={metaNamespace(pfTarget)}
        podName={metaName(pfTarget)}
        onClose={() => setPfTarget(null)}
      />
    </>
  )
}

export function DeploymentsPage() {
  const { namespace } = useNamespace()
  const { canMutate } = useAuth()
  const [yamlTarget, setYamlTarget] = useState<any | null>(null)
  const [scaleTarget, setScaleTarget] = useState<any | null>(null)
  const [restartTarget, setRestartTarget] = useState<any | null>(null)
  const [scaleBusy, setScaleBusy] = useState(false)
  const [scaleErr, setScaleErr] = useState('')
  const refetchRef = useRef<(() => void) | null>(null)

  return (
    <>
      <ResourceListPage
        title="DEPLOYMENTS"
        namespaced
        resourceKey="deployments"
        queryFn={() => listNamespacedResource(namespace, 'deployments')}
        columns={[
          {
            key: 'name',
            header: 'Name',
            render: (item) => (
              <Link
                className="font-semibold text-cyan hover:underline"
                to={`/deployments/${metaNamespace(item)}/${metaName(item)}`}
              >
                {metaName(item)}
              </Link>
            ),
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
          { key: 'age', header: 'Age', render: ageCell },
          { key: 'created', header: 'Created', render: createdCell },
        ]}
        actions={(item, { refetch }) => (
          <div className="flex items-center gap-1">
            {canMutate('deployments') ? (
              <>
                <Button
                  variant="ghost"
                  className="px-2 py-1 text-xs"
                  type="button"
                  title="Scale"
                  onClick={() => {
                    refetchRef.current = refetch
                    setScaleErr('')
                    setScaleTarget(item)
                  }}
                >
                  Scale
                </Button>
                <Button
                  variant="ghost"
                  className="px-2 py-1"
                  type="button"
                  title="Restart"
                  onClick={() => {
                    refetchRef.current = refetch
                    setRestartTarget(item)
                  }}
                >
                  <RefreshCw className="h-3.5 w-3.5" />
                </Button>
              </>
            ) : null}
            <YamlActionButton resource="deployments" onClick={() => setYamlTarget(item)} />
          </div>
        )}
      />
      <ResourceYamlModal
        open={Boolean(yamlTarget)}
        item={yamlTarget}
        resource="deployments"
        onClose={() => setYamlTarget(null)}
      />
      <ScaleDialog
        open={Boolean(scaleTarget)}
        resourceName={`${metaNamespace(scaleTarget)}/${metaName(scaleTarget)}`}
        current={scaleTarget?.spec?.replicas ?? 0}
        busy={scaleBusy}
        onClose={() => setScaleTarget(null)}
        onConfirm={async (replicas) => {
          setScaleBusy(true)
          setScaleErr('')
          try {
            await apiPatch(
              `/api/v1/namespaces/${metaNamespace(scaleTarget)}/deployments/${metaName(scaleTarget)}`,
              { spec: { replicas } },
            )
            refetchRef.current?.()
            setScaleTarget(null)
          } catch (e: any) {
            setScaleErr(e?.message || 'Scale failed')
          } finally {
            setScaleBusy(false)
          }
        }}
      />
      <ConfirmDialog
        open={Boolean(restartTarget)}
        title="RESTART DEPLOYMENT"
        danger={false}
        confirmLabel="Restart"
        busy={scaleBusy}
        description={`Roll restart ${metaNamespace(restartTarget)}/${metaName(restartTarget)}?`}
        onClose={() => setRestartTarget(null)}
        onConfirm={async () => {
          setScaleBusy(true)
          try {
            const ns = metaNamespace(restartTarget)
            const name = metaName(restartTarget)
            const live = await getNamespacedResource(ns, 'deployments', name)
            const clone = structuredClone(live)
            clone.spec.template.metadata = clone.spec.template.metadata || {}
            clone.spec.template.metadata.annotations = {
              ...(clone.spec.template.metadata.annotations || {}),
              'cilikube.io/restartedAt': new Date().toISOString(),
            }
            await updateNamespacedResource(ns, 'deployments', name, clone)
            refetchRef.current?.()
            setRestartTarget(null)
          } catch (e: any) {
            setScaleErr(e?.message || 'Restart failed')
            setRestartTarget(null)
          } finally {
            setScaleBusy(false)
          }
        }}
      />
      <ConfirmDialog
        open={Boolean(scaleErr)}
        title="OPERATION FAILED"
        danger={false}
        confirmLabel="OK"
        description={scaleErr}
        onClose={() => setScaleErr('')}
        onConfirm={() => setScaleErr('')}
      />
    </>
  )
}

export function StatefulSetsPage() {
  const { namespace } = useNamespace()
  const { canMutate } = useAuth()
  const { yamlButton, yamlModal } = useYamlModal('statefulsets', true)
  const [scaleTarget, setScaleTarget] = useState<any | null>(null)
  const [restartTarget, setRestartTarget] = useState<any | null>(null)
  const [scaleBusy, setScaleBusy] = useState(false)
  const [scaleErr, setScaleErr] = useState('')
  const refetchRef = useRef<(() => void) | null>(null)

  return (
    <>
      <ResourceListPage
        title="STATEFULSETS"
        namespaced
        resourceKey="statefulsets"
        queryFn={() => listNamespacedResource(namespace, 'statefulsets')}
        columns={[
          {
            key: 'name',
            header: 'Name',
            render: (item) => (
              <Link
                className="font-semibold text-cyan hover:underline"
                to={`/statefulsets/${metaNamespace(item)}/${metaName(item)}`}
              >
                {metaName(item)}
              </Link>
            ),
          },
          { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
          {
            key: 'ready',
            header: 'Ready',
            render: (item) => (
              <span className="font-mono text-xs">
                {item.status?.readyReplicas ?? 0}/{item.spec?.replicas ?? 0}
              </span>
            ),
          },
          {
            key: 'service',
            header: 'Service',
            render: (item) => item.spec?.serviceName || '-',
          },
          { key: 'age', header: 'Age', render: ageCell },
          { key: 'created', header: 'Created', render: createdCell },
        ]}
        actions={(item, { refetch }) => (
          <div className="flex items-center gap-1">
            {canMutate('statefulsets') ? (
              <>
                <Button
                  variant="ghost"
                  className="px-2 py-1 text-xs"
                  type="button"
                  onClick={() => {
                    refetchRef.current = refetch
                    setScaleTarget(item)
                  }}
                >
                  Scale
                </Button>
                <Button
                  variant="ghost"
                  className="px-2 py-1"
                  type="button"
                  title="Restart"
                  onClick={() => {
                    refetchRef.current = refetch
                    setRestartTarget(item)
                  }}
                >
                  <RefreshCw className="h-3.5 w-3.5" />
                </Button>
              </>
            ) : null}
            {yamlButton(item)}
          </div>
        )}
      />
      {yamlModal}
      <ScaleDialog
        open={Boolean(scaleTarget)}
        resourceName={`${metaNamespace(scaleTarget)}/${metaName(scaleTarget)}`}
        current={scaleTarget?.spec?.replicas ?? 0}
        busy={scaleBusy}
        onClose={() => setScaleTarget(null)}
        onConfirm={async (replicas) => {
          setScaleBusy(true)
          setScaleErr('')
          try {
            await apiPatch(
              `/api/v1/namespaces/${metaNamespace(scaleTarget)}/statefulsets/${metaName(scaleTarget)}`,
              { spec: { replicas } },
            )
            refetchRef.current?.()
            setScaleTarget(null)
          } catch (e: any) {
            setScaleErr(e?.message || 'Scale failed')
          } finally {
            setScaleBusy(false)
          }
        }}
      />
      <ConfirmDialog
        open={Boolean(restartTarget)}
        title="RESTART STATEFULSET"
        danger={false}
        confirmLabel="Restart"
        busy={scaleBusy}
        description={`Roll restart ${metaNamespace(restartTarget)}/${metaName(restartTarget)}?`}
        onClose={() => setRestartTarget(null)}
        onConfirm={async () => {
          setScaleBusy(true)
          try {
            const ns = metaNamespace(restartTarget)
            const name = metaName(restartTarget)
            const live = await getNamespacedResource(ns, 'statefulsets', name)
            const clone = structuredClone(live)
            clone.spec.template.metadata = clone.spec.template.metadata || {}
            clone.spec.template.metadata.annotations = {
              ...(clone.spec.template.metadata.annotations || {}),
              'cilikube.io/restartedAt': new Date().toISOString(),
            }
            await updateNamespacedResource(ns, 'statefulsets', name, clone)
            refetchRef.current?.()
            setRestartTarget(null)
          } catch (e: any) {
            setScaleErr(e?.message || 'Restart failed')
            setRestartTarget(null)
          } finally {
            setScaleBusy(false)
          }
        }}
      />
      <ConfirmDialog
        open={Boolean(scaleErr)}
        title="OPERATION FAILED"
        danger={false}
        confirmLabel="OK"
        description={scaleErr}
        onClose={() => setScaleErr('')}
        onConfirm={() => setScaleErr('')}
      />
    </>
  )
}

export function DaemonSetsPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('daemonsets', true)
  return (
    <>
    <ResourceListPage
      title="DAEMONSETS"
      namespaced
      resourceKey="daemonsets"
      queryFn={() => listNamespacedResource(namespace, 'daemonsets')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('daemonsets', item),
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'ready',
          header: 'Ready',
          render: (item) => (
            <span className="font-mono text-xs">
              {item.status?.numberReady ?? 0}/{item.status?.desiredNumberScheduled ?? 0}
            </span>
          ),
        },
        {
          key: 'updated',
          header: 'Updated',
          render: (item) => item.status?.updatedNumberScheduled ?? 0,
        },
        {
          key: 'available',
          header: 'Available',
          render: (item) => item.status?.numberAvailable ?? 0,
        },
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function JobsPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('jobs', true)
  return (
    <>
    <ResourceListPage
      title="JOBS"
      namespaced
      resourceKey="jobs"
      queryFn={() => listNamespacedResource(namespace, 'jobs')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('jobs', item),
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'completions',
          header: 'Completions',
          render: (item) => (
            <span className="font-mono text-xs">
              {item.status?.succeeded ?? 0}/{item.spec?.completions ?? 1}
            </span>
          ),
        },
        {
          key: 'duration',
          header: 'Active',
          render: (item) => item.status?.active ?? 0,
        },
        {
          key: 'failed',
          header: 'Failed',
          render: (item) => item.status?.failed ?? 0,
        },
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function CronJobsPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('cronjobs', true)
  return (
    <>
    <ResourceListPage
      title="CRONJOBS"
      namespaced
      resourceKey="cronjobs"
      queryFn={() => listNamespacedResource(namespace, 'cronjobs')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('cronjobs', item),
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'schedule',
          header: 'Schedule',
          render: (item) => <span className="font-mono text-xs">{item.spec?.schedule || '-'}</span>,
        },
        {
          key: 'suspend',
          header: 'Suspend',
          render: (item) => (
            <Badge tone={item.spec?.suspend ? 'warn' : 'ok'}>
              {item.spec?.suspend ? 'Yes' : 'No'}
            </Badge>
          ),
        },
        {
          key: 'last',
          header: 'Last schedule',
          render: (item) =>
            item.status?.lastScheduleTime ? (
              <span className="inline-flex flex-col">
                {ageCell({ metadata: { creationTimestamp: item.status.lastScheduleTime } })}
                {createdCell({ metadata: { creationTimestamp: item.status.lastScheduleTime } })}
              </span>
            ) : (
              '-'
            ),
        },
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function NetworkPoliciesPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('networkpolicies', true)
  return (
    <>
    <ResourceListPage
      title="NETWORK POLICIES"
      namespaced
      resourceKey="networkpolicies"
      queryFn={() => listNamespacedResource(namespace, 'networkpolicies')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('networkpolicies', item),
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'podSelector',
          header: 'Pod selector',
          render: (item) => {
            const labels = item.spec?.podSelector?.matchLabels || {}
            const entries = Object.entries(labels)
            return entries.length
              ? entries.map(([k, v]) => `${k}=${v}`).join(', ')
              : 'All pods'
          },
        },
        {
          key: 'policy',
          header: 'Policy types',
          render: (item) => (item.spec?.policyTypes || []).join(', ') || '-',
        },
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function ServiceAccountsPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('serviceaccounts', true)
  return (
    <>
    <ResourceListPage
      title="SERVICE ACCOUNTS"
      namespaced
      resourceKey="serviceaccounts"
      queryFn={() => listNamespacedResource(namespace, 'serviceaccounts')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('serviceaccounts', item),
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'secrets',
          header: 'Secrets',
          render: (item) => (item.secrets || []).length,
        },
        {
          key: 'automount',
          header: 'Automount',
          render: (item) =>
            item.automountServiceAccountToken === false ? (
              <Badge tone="warn">Off</Badge>
            ) : (
              <Badge tone="ok">On</Badge>
            ),
        },
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function ServicesPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('services', true)
  return (
    <>
    <ResourceListPage
      title="SERVICES"
      namespaced
      resourceKey="services"
      queryFn={() => listNamespacedResource(namespace, 'services')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('services', item),
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function IngressPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('ingresses', true)
  return (
    <>
    <ResourceListPage
      title="INGRESS"
      namespaced
      resourceKey="ingresses"
      queryFn={() => listNamespacedResource(namespace, 'ingresses')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('ingresses', item),
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function ConfigMapsPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('configmaps', true)
  return (
    <>
    <ResourceListPage
      title="CONFIGMAPS"
      namespaced
      resourceKey="configmaps"
      queryFn={() => listNamespacedResource(namespace, 'configmaps')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('configmaps', item),
        },
        { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
        {
          key: 'keys',
          header: 'Data keys',
          render: (item) => Object.keys(item.data || {}).length,
        },
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function SecretsPage() {
  const { namespace } = useNamespace()
  const { checkPermission } = useAuth()
  const { yamlButton, yamlModal } = useYamlModal('secrets', true)
  if (!checkPermission('secrets', 'read')) {
    return (
      <div className="rounded border border-warn/40 bg-warn/10 px-5 py-8 text-sm text-warn">
        Secrets are restricted to administrators. Your role cannot view or modify secret data.
      </div>
    )
  }
  return (
    <>
    <ResourceListPage
      title="SECRETS"
      namespaced
      resourceKey="secrets"
      queryFn={() => listNamespacedResource(namespace, 'secrets')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('secrets', item),
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function PVPage() {
  const { yamlButton, yamlModal } = useYamlModal('persistentvolumes', false)
  return (
    <>
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function PVCPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('persistentvolumeclaims', true)
  return (
    <>
    <ResourceListPage
      title="PERSISTENT VOLUME CLAIMS"
      namespaced
      resourceKey="persistentvolumeclaims"
      queryFn={() => listNamespacedResource(namespace, 'persistentvolumeclaims')}
      columns={[
        {
          key: 'name',
          header: 'Name',
          render: (item) => nameLink('persistentvolumeclaims', item),
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function StorageClassPage() {
  const { yamlButton, yamlModal } = useYamlModal('storageclasses', false)
  return (
    <>
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function RolesPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('roles', true)
  return (
    <>
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function RoleBindingsPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('rolebindings', true)
  return (
    <>
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function ClusterRolesPage() {
  const { yamlButton, yamlModal } = useYamlModal('clusterroles', false)
  return (
    <>
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function ClusterRoleBindingsPage() {
  const { yamlButton, yamlModal } = useYamlModal('clusterrolebindings', false)
  return (
    <>
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
        { key: 'age', header: 'Age', render: ageCell },
        { key: 'created', header: 'Created', render: createdCell },
      ]}
      actions={(item) => yamlButton(item)}
    />
    {yamlModal}
    </>
  )

}

export function HPAPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('horizontalpodautoscalers', true)
  return (
    <>
      <ResourceListPage
        title="HORIZONTAL POD AUTOSCALERS"
        namespaced
        resourceKey="horizontalpodautoscalers"
        queryFn={() => listNamespacedResource(namespace, 'horizontalpodautoscalers')}
        columns={[
          {
            key: 'name',
            header: 'Name',
            render: (item) => nameLink('horizontalpodautoscalers', item),
          },
          { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
          {
            key: 'target',
            header: 'Target',
            render: (item) =>
              item.spec?.scaleTargetRef
                ? `${item.spec.scaleTargetRef.kind}/${item.spec.scaleTargetRef.name}`
                : '-',
          },
          {
            key: 'min',
            header: 'Min',
            render: (item) => item.spec?.minReplicas ?? '-',
          },
          {
            key: 'max',
            header: 'Max',
            render: (item) => item.spec?.maxReplicas ?? '-',
          },
          {
            key: 'current',
            header: 'Current',
            render: (item) => item.status?.currentReplicas ?? '-',
          },
          { key: 'age', header: 'Age', render: ageCell },
          { key: 'created', header: 'Created', render: createdCell },
        ]}
        actions={(item) => yamlButton(item)}
      />
      {yamlModal}
    </>
  )
}

export function PDBPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('poddisruptionbudgets', true)
  return (
    <>
      <ResourceListPage
        title="POD DISRUPTION BUDGETS"
        namespaced
        resourceKey="poddisruptionbudgets"
        queryFn={() => listNamespacedResource(namespace, 'poddisruptionbudgets')}
        columns={[
          {
            key: 'name',
            header: 'Name',
            render: (item) => nameLink('poddisruptionbudgets', item),
          },
          { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
          {
            key: 'min',
            header: 'Min available',
            render: (item) =>
              item.spec?.minAvailable != null ? String(item.spec.minAvailable) : '-',
          },
          {
            key: 'max',
            header: 'Max unavailable',
            render: (item) =>
              item.spec?.maxUnavailable != null ? String(item.spec.maxUnavailable) : '-',
          },
          {
            key: 'allowed',
            header: 'Allowed disruptions',
            render: (item) => item.status?.disruptionsAllowed ?? '-',
          },
          { key: 'age', header: 'Age', render: ageCell },
          { key: 'created', header: 'Created', render: createdCell },
        ]}
        actions={(item) => yamlButton(item)}
      />
      {yamlModal}
    </>
  )
}

export function ResourceQuotasPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('resourcequotas', true)
  return (
    <>
      <ResourceListPage
        title="RESOURCE QUOTAS"
        namespaced
        resourceKey="resourcequotas"
        queryFn={() => listNamespacedResource(namespace, 'resourcequotas')}
        columns={[
          {
            key: 'name',
            header: 'Name',
            render: (item) => nameLink('resourcequotas', item),
          },
          { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
          {
            key: 'hard',
            header: 'Hard limits',
            render: (item) => Object.keys(item.status?.hard || item.spec?.hard || {}).length,
          },
          {
            key: 'used',
            header: 'Used keys',
            render: (item) => Object.keys(item.status?.used || {}).length,
          },
          { key: 'age', header: 'Age', render: ageCell },
          { key: 'created', header: 'Created', render: createdCell },
        ]}
        actions={(item) => yamlButton(item)}
      />
      {yamlModal}
    </>
  )
}

export function LimitRangesPage() {
  const { namespace } = useNamespace()
  const { yamlButton, yamlModal } = useYamlModal('limitranges', true)
  return (
    <>
      <ResourceListPage
        title="LIMIT RANGES"
        namespaced
        resourceKey="limitranges"
        queryFn={() => listNamespacedResource(namespace, 'limitranges')}
        columns={[
          {
            key: 'name',
            header: 'Name',
            render: (item) => nameLink('limitranges', item),
          },
          { key: 'ns', header: 'Namespace', render: (item) => metaNamespace(item) },
          {
            key: 'limits',
            header: 'Limits',
            render: (item) => (item.spec?.limits || []).length,
          },
          {
            key: 'types',
            header: 'Types',
            render: (item) =>
              (item.spec?.limits || [])
                .map((l: any) => l.type)
                .filter(Boolean)
                .join(', ') || '-',
          },
          { key: 'age', header: 'Age', render: ageCell },
          { key: 'created', header: 'Created', render: createdCell },
        ]}
        actions={(item) => yamlButton(item)}
      />
      {yamlModal}
    </>
  )
}
