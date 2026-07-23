import { useState, type ReactNode } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { useCluster } from '@/store/cluster'
import { useNamespace } from '@/store/namespace'
import { useAuth } from '@/store/auth'
import { Badge, Button, EmptyState, PageHeader } from '@/components/ui'
import { HudTable, HudTablePanel, ListPageFrame } from '@/components/HudTableScroll'
import { ResourceCreateModal } from '@/components/ResourceCreateModal'
import { AgeCell, CreatedCell } from '@/components/AgeCell'
import { metaCreated } from '@/api/resources'

export type ResourceColumn = {
  key: string
  header: string
  /** Hover hint — useful for abbreviated headers like %CPU/R */
  title?: string
  render: (item: any) => ReactNode
}

export function ResourceListPage({
  title,
  subtitle,
  resourceKey,
  namespaced,
  queryFn,
  columns,
  extraAction,
  actions,
  creatable = true,
  pinFirstColumn = false,
}: {
  title: string
  subtitle?: string
  resourceKey: string
  namespaced?: boolean
  queryFn: () => Promise<any[]>
  columns: ResourceColumn[]
  extraAction?: ReactNode
  actions?: (item: any, helpers: { refetch: () => void }) => ReactNode
  /** Show Create YAML button when user can mutate this resource */
  creatable?: boolean
  /** Wide metric tables only: pin Name + allow horizontal scroll */
  pinFirstColumn?: boolean
}) {
  const { clusterId } = useCluster()
  const { namespace, isAllNamespaces } = useNamespace()
  const { canMutate } = useAuth()
  const [createOpen, setCreateOpen] = useState(false)
  const showCreate = creatable && canMutate(resourceKey)

  const { data = [], isLoading, refetch } = useQuery({
    queryKey: [resourceKey, clusterId, namespaced ? (namespace || '__all__') : 'cluster'],
    queryFn,
    enabled: Boolean(clusterId) && (!namespaced || namespace !== undefined),
  })

  const cols = actions
    ? [
        ...columns,
        {
          key: '_actions',
          header: 'Actions',
          render: (item: any) => actions(item, { refetch: () => void refetch() }),
        },
      ]
    : columns

  return (
    <ListPageFrame>
      <PageHeader
        title={title}
        subtitle={
          subtitle ||
          (namespaced
            ? isAllNamespaces
              ? 'All namespaces'
              : `Namespace: ${namespace}`
            : 'Cluster-scoped resources')
        }
        action={
          <div className="flex items-center gap-3">
            {showCreate ? (
              <Button
                type="button"
                variant="outline"
                className="px-3 py-1.5 text-xs"
                onClick={() => setCreateOpen(true)}
              >
                Create YAML
              </Button>
            ) : null}
            {extraAction}
            <Badge tone="accent">{data.length}</Badge>
          </div>
        }
      />

      <ResourceCreateModal
        open={createOpen}
        resource={resourceKey}
        namespaced={namespaced}
        onClose={() => setCreateOpen(false)}
      />

      <HudTablePanel pinFirst={pinFirstColumn} wide={pinFirstColumn}>
        <HudTable pinFirst={pinFirstColumn} wide={pinFirstColumn}>
          <thead>
            <tr>
              {cols.map((col) => (
                <th key={col.key} title={col.title}>
                  {col.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.map((item: any, index: number) => (
              <motion.tr
                key={item?.metadata?.uid || item?.metadata?.name || index}
                initial={{ opacity: 0, y: 6 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: Math.min(index, 12) * 0.02 }}
              >
                {cols.map((col) => (
                  <td key={col.key}>{col.render(item)}</td>
                ))}
              </motion.tr>
            ))}
            {!isLoading && !data.length ? (
              <tr>
                <td colSpan={cols.length}>
                  <EmptyState>
                    {namespaced
                      ? isAllNamespaces
                        ? `No ${resourceKey} across all namespaces.`
                        : `No ${resourceKey} in namespace ${namespace}.`
                      : `No ${resourceKey} found.`}
                  </EmptyState>
                </td>
              </tr>
            ) : null}
          </tbody>
        </HudTable>
      </HudTablePanel>
    </ListPageFrame>
  )
}

/** Age column — relative (k9s style). */
export function ageCell(item: any) {
  return <AgeCell value={metaCreated(item)} />
}

/** Created column — absolute local time. */
export function createdCell(item: any) {
  return <CreatedCell value={metaCreated(item)} />
}

/** @deprecated use ageCell + createdCell; kept for call sites that only need age */
export function ageCellFrom(value: unknown) {
  return <AgeCell value={value} />
}

/** Pair of columns: Age then Created — drop into ResourceListPage columns arrays. */
export const ageAndCreatedColumns = [
  { key: 'age', header: 'Age', render: ageCell },
  { key: 'created', header: 'Created', render: createdCell },
] as const
