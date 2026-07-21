import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import dayjs from 'dayjs'
import { listNodes } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, PageHeader } from '@/components/ui'

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
      <PageHeader title="Nodes" subtitle="Cluster compute fabric" />
      <Card className="overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full text-left text-sm">
            <thead className="bg-mist/80 text-xs uppercase tracking-wide text-ink-soft">
              <tr>
                <th className="px-5 py-3">Name</th>
                <th className="px-5 py-3">Status</th>
                <th className="px-5 py-3">Roles</th>
                <th className="px-5 py-3">Version</th>
                <th className="px-5 py-3">Created</th>
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
                    className="border-t border-line/70"
                  >
                    <td className="px-5 py-3 font-semibold">{name}</td>
                    <td className="px-5 py-3">
                      <Badge tone={nodeReady(node) ? 'ok' : 'danger'}>
                        {nodeReady(node) ? 'Ready' : 'NotReady'}
                      </Badge>
                    </td>
                    <td className="px-5 py-3 text-ink-soft">
                      {roles.length ? roles.join(', ') : 'worker'}
                    </td>
                    <td className="px-5 py-3 font-mono text-xs">
                      {node.status?.nodeInfo?.kubeletVersion || '-'}
                    </td>
                    <td className="px-5 py-3 text-ink-soft">
                      {node.metadata?.creationTimestamp
                        ? dayjs(node.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm')
                        : '-'}
                    </td>
                  </motion.tr>
                )
              })}
              {!isLoading && !data.length ? (
                <tr>
                  <td colSpan={5} className="px-5 py-10 text-ink-soft">
                    No nodes found for this cluster.
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
