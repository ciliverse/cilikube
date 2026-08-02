import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { listNodes } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { AgeCell, CreatedCell } from '@/components/AgeCell'
import { HudTable, HudTablePanel, ListPageFrame } from '@/components/HudTableScroll'
import { Badge, EmptyState, PageHeader } from '@/components/ui'
import { metaCreated } from '@/api/resources'
import { shouldSkipEnterAnim } from '@/lib/motionPrefs'
import { useTranslation } from 'react-i18next'

function nodeReady(node: any) {
  const conditions = node?.status?.conditions || []
  const ready = conditions.find((c: any) => c.type === 'Ready')
  return ready?.status === 'True'
}

export function NodesPage() {
  const { t } = useTranslation()
  const { clusterId } = useCluster()
  const { data = [], isLoading } = useQuery({
    queryKey: ['nodes', clusterId],
    queryFn: listNodes,
    enabled: Boolean(clusterId),
  })

  return (
    <ListPageFrame>
      <PageHeader title={t('nodes.title')} subtitle={t('nodes.subtitle')} />
      <HudTablePanel>
          <HudTable>
            <thead>
              <tr>
                <th>Name</th>
                <th>Status</th>
                <th>Roles</th>
                <th>Version</th>
                <th>Age</th>
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
                    initial={shouldSkipEnterAnim() ? false : { opacity: 0, y: 8 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: shouldSkipEnterAnim() ? 0 : index * 0.03 }}
                  >
                    <td>
                      <Link className="font-semibold text-cyan hover:underline" to={`/nodes/${name}`}>
                        {name}
                      </Link>
                    </td>
                    <td>
                      <Badge tone={nodeReady(node) ? 'ok' : 'danger'}>
                        {nodeReady(node) ? 'Ready' : 'NotReady'}
                      </Badge>
                    </td>
                    <td>{roles.length ? roles.join(', ') : 'worker'}</td>
                    <td className="font-mono text-xs">
                      {node.status?.nodeInfo?.kubeletVersion || '-'}
                    </td>
                    <td>
                      <AgeCell value={metaCreated(node)} />
                    </td>
                    <td>
                      <CreatedCell value={metaCreated(node)} />
                    </td>
                  </motion.tr>
                )
              })}
              {!isLoading && !data.length ? (
                <tr>
                  <td colSpan={6}>
                    <EmptyState>No nodes found for this cluster.</EmptyState>
                  </td>
                </tr>
              ) : null}
            </tbody>
          </HudTable>
      </HudTablePanel>
    </ListPageFrame>
  )
}
