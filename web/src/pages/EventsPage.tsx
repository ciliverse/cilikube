import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import dayjs from 'dayjs'
import { listEvents, listNamespaces } from '@/api/cluster'
import { useCluster } from '@/store/cluster'
import { Badge, Card, EmptyState, PageHeader } from '@/components/ui'

export function EventsPage() {
  const { clusterId } = useCluster()
  const [namespace, setNamespace] = useState('')
  const [type, setType] = useState('')

  const nsQ = useQuery({
    queryKey: ['namespaces', clusterId],
    queryFn: listNamespaces,
    enabled: Boolean(clusterId),
  })

  const eventsQ = useQuery({
    queryKey: ['events', clusterId, namespace],
    queryFn: () => listEvents({ namespace: namespace || undefined, limit: 100 }),
    enabled: Boolean(clusterId),
    refetchInterval: 30_000,
  })

  const events = (eventsQ.data?.events || []).filter((e: any) =>
    type ? e.type === type : true,
  )

  return (
    <div>
      <PageHeader
        title="EVENTS"
        subtitle="Cluster signal stream"
        action={<Badge tone="neutral">{events.length} events</Badge>}
      />

      <Card className="mb-4 flex flex-wrap gap-3 p-4">
        <select
          value={namespace}
          onChange={(e) => setNamespace(e.target.value)}
          className="hud-select w-auto min-w-[160px]"
        >
          <option value="">All namespaces</option>
          {(nsQ.data || []).map((ns) => (
            <option key={ns} value={ns}>
              {ns}
            </option>
          ))}
        </select>
        <select
          value={type}
          onChange={(e) => setType(e.target.value)}
          className="hud-select w-auto min-w-[140px]"
        >
          <option value="">All types</option>
          <option value="Normal">Normal</option>
          <option value="Warning">Warning</option>
        </select>
      </Card>

      <div className="space-y-3">
        {events.map((event: any, index: number) => (
          <motion.div
            key={event.id || `${event.name}-${index}`}
            initial={{ opacity: 0, x: -8 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: Math.min(index, 10) * 0.03 }}
          >
            <Card className="p-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="flex flex-wrap items-center gap-2">
                  <Badge tone={event.type === 'Warning' ? 'warn' : 'ok'}>
                    {event.type || 'Normal'}
                  </Badge>
                  <span className="font-semibold text-cyan">{event.reason}</span>
                  <span className="text-xs text-text-dim">
                    {event.objectKind}/{event.object || event.name}
                  </span>
                </div>
                <span className="text-xs text-text-dim">
                  {event.lastTime
                    ? dayjs(event.lastTime).format('YYYY-MM-DD HH:mm:ss')
                    : '-'}
                </span>
              </div>
              <p className="mt-2 text-sm text-text-dim">{event.message}</p>
              <div className="mt-2 text-xs text-text-dim">
                ns: {event.namespace || '-'} · count: {event.count ?? 1}
              </div>
            </Card>
          </motion.div>
        ))}
        {!eventsQ.isLoading && !events.length ? (
          <Card>
            <EmptyState>No events matched.</EmptyState>
          </Card>
        ) : null}
      </div>
    </div>
  )
}
