import type { ReactNode } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import dayjs from 'dayjs'
import { useCluster } from '@/store/cluster'
import { useNamespace } from '@/store/namespace'
import { Badge, Card, EmptyState, PageHeader } from '@/components/ui'
import { metaCreated } from '@/api/resources'

export type ResourceColumn = {
  key: string
  header: string
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
}: {
  title: string
  subtitle?: string
  resourceKey: string
  namespaced?: boolean
  queryFn: () => Promise<any[]>
  columns: ResourceColumn[]
  extraAction?: ReactNode
}) {
  const { clusterId } = useCluster()
  const { namespace } = useNamespace()

  const { data = [], isLoading } = useQuery({
    queryKey: [resourceKey, clusterId, namespaced ? namespace : 'cluster'],
    queryFn,
    enabled: Boolean(clusterId) && (!namespaced || Boolean(namespace)),
  })

  return (
    <div>
      <PageHeader
        title={title}
        subtitle={
          subtitle ||
          (namespaced ? `Namespace: ${namespace}` : 'Cluster-scoped resources')
        }
        action={
          <div className="flex items-center gap-3">
            {extraAction}
            <Badge tone="accent">{data.length}</Badge>
          </div>
        }
      />

      <Card className="overflow-hidden">
        <div className="overflow-x-auto">
          <table className="hud-table">
            <thead>
              <tr>
                {columns.map((col) => (
                  <th key={col.key}>{col.header}</th>
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
                  {columns.map((col) => (
                    <td key={col.key}>{col.render(item)}</td>
                  ))}
                </motion.tr>
              ))}
              {!isLoading && !data.length ? (
                <tr>
                  <td colSpan={columns.length}>
                    <EmptyState>
                      {namespaced
                        ? `No ${resourceKey} in namespace ${namespace}.`
                        : `No ${resourceKey} found.`}
                    </EmptyState>
                  </td>
                </tr>
              ) : null}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  )
}

export function createdCell(item: any) {
  const ts = metaCreated(item)
  return (
    <span className="text-text-dim">
      {ts ? dayjs(ts).format('YYYY-MM-DD HH:mm') : '-'}
    </span>
  )
}
