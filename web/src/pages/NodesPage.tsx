import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import dayjs from 'dayjs'
import { listNodes } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, EmptyState, PageHeader } from '@/components/ui'

function nodeReady(node: any) {
  const conditions = node?.status?.conditions || []
  const ready = conditions.find((c: any) => c.type === 'Ready')
  return ready?.status === 'True'
}

export function NodesPage() {
  const { clusterId } = useCluster()
  const { data = [], isLoading } = useQuery({
    queryKey: ['nodes', clusterId],
    queryFn: listNodes,
    enabled: Boolean(clusterId),
  })

  return (
    <div>
      <PageHeader title="NODES" subtitle="Cluster compute fabric" />
      <Card className="overflow-hidden">
        <div className="overflow-x-auto">
          <table className="hud-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Status</th>
                <th>Roles</th>
                <th>Version</th>
                <th>Created</th>
              </tr>
            </thead>
            <tbody>
              {data.map((node: any, index: number) => {
                const name = node.metadata?.name
                const roles = Object.keys(node.metadata?.labels || {})
                  .filter((k) => k.startsWith('node-role.kubernetes.io/'))
                  .map((k) => k.replace('node-role.kubernetes.io/', ''))
                return (
                  <motion.tr
                    key={name}
                    initial={{ opacity: 0, y: 8 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.03 }}
                  >
                    <td className="font-semibold text-cyan">{name}</td>
                    <td>
                      <Badge tone={nodeReady(node) ? 'ok' : 'danger'}>
                        {nodeReady(node) ? 'Ready' : 'NotReady'}
                      </Badge>
                    </td>
                    <td className="text-text-dim">
                      {roles.length ? roles.join(', ') : 'worker'}
                    </td>
                    <td className="font-mono text-xs">
                      {node.status?.nodeInfo?.kubeletVersion || '-'}
                    </td>
                    <td className="text-text-dim">
                      {node.metadata?.creationTimestamp
                        ? dayjs(node.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm')
                        : '-'}
                    </td>
                  </motion.tr>
                )
              })}
              {!isLoading && !data.length ? (
                <tr>
                  <td colSpan={5}>
                    <EmptyState>No nodes found for this cluster.</EmptyState>
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
