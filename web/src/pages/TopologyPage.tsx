import { useEffect, useMemo, useCallback, useState, type MouseEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  MarkerType,
  Handle,
  Position,
  type Node,
  type Edge,
  type NodeProps,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { Loader2, RefreshCw, Search } from 'lucide-react'
import {
  getTopologyGraph,
  getTopologyTraffic,
  type TopologyEdge,
  type TopologyNode,
} from '@/api/topology'
import { layoutTopology, layoutTopologyByGroups } from '@/lib/topologyLayout'
import { useCluster } from '@/store/cluster'
import { ALL_NAMESPACES, useNamespace } from '@/store/namespace'
import { Button, EmptyState, PageHeader } from '@/components/ui'
import { cn } from '@/lib/utils'

/** Undirected BFS: service + ingress / workloads / pods linked through edges. */
function connectedIds(seed: string, edges: TopologyEdge[]): Set<string> {
  const adj = new Map<string, string[]>()
  for (const e of edges) {
    const a = adj.get(e.source) || []
    a.push(e.target)
    adj.set(e.source, a)
    const b = adj.get(e.target) || []
    b.push(e.source)
    adj.set(e.target, b)
  }
  const seen = new Set<string>([seed])
  const q = [seed]
  while (q.length) {
    const cur = q.shift()!
    for (const next of adj.get(cur) || []) {
      if (seen.has(next)) continue
      seen.add(next)
      q.push(next)
    }
  }
  return seen
}

const KIND_ORDER = [
  'ingress',
  'service',
  'deployment',
  'statefulset',
  'daemonset',
  'pod',
  'job',
  'cronjob',
  'hpa',
  'configmap',
] as const

function TopologyResourceNode({ data }: NodeProps) {
  const d = data as TopologyNode
  return (
    <div className={cn('topo-node', `is-${d.status || 'unknown'}`, `kind-${d.kind}`)}>
      <Handle type="target" position={Position.Left} className="topo-handle" />
      <div className="topo-node-kind">{d.kind}</div>
      <div className="topo-node-name" title={d.name}>
        {d.name}
      </div>
      {d.subtitle ? <div className="topo-node-sub">{d.subtitle}</div> : null}
      <Handle type="source" position={Position.Right} className="topo-handle" />
    </div>
  )
}

const nodeTypes = { topo: TopologyResourceNode }

export function TopologyPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { clusterId } = useCluster()
  const { namespace } = useNamespace()
  const [groupBy, setGroupBy] = useState<'app' | 'namespace'>('app')
  const [mode, setMode] = useState<'resources' | 'traffic'>('resources')
  const [hiddenKinds, setHiddenKinds] = useState<Set<string>>(() => new Set(['configmap']))
  /** Focus one Service by node id; null = show all. */
  const [focusServiceId, setFocusServiceId] = useState<string | null>(null)
  const [serviceQuery, setServiceQuery] = useState('')

  const nsReady = Boolean(namespace && namespace !== ALL_NAMESPACES)

  useEffect(() => {
    setFocusServiceId(null)
    setServiceQuery('')
  }, [clusterId, namespace])

  const graphQ = useQuery({
    queryKey: ['topology', clusterId, namespace, groupBy],
    queryFn: () => getTopologyGraph({ namespace, groupBy }),
    enabled: nsReady,
    refetchInterval: 30_000,
  })

  const trafficQ = useQuery({
    queryKey: ['topology-traffic', clusterId, namespace],
    queryFn: () => getTopologyTraffic(namespace),
    enabled: nsReady && mode === 'traffic',
    refetchInterval: 15_000,
  })

  const services = useMemo(() => {
    const raw = graphQ.data?.nodes || []
    return raw
      .filter((n) => n.kind === 'service')
      .slice()
      .sort((a, b) => a.name.localeCompare(b.name))
  }, [graphQ.data])

  const filteredServices = useMemo(() => {
    const q = serviceQuery.trim().toLowerCase()
    if (!q) return services
    return services.filter((s) => s.name.toLowerCase().includes(q))
  }, [services, serviceQuery])

  const { nodes, edges } = useMemo(() => {
    const raw = graphQ.data
    if (!raw) return { nodes: [] as Node[], edges: [] as Edge[] }

    const trafficMap = new Map<string, number>()
    if (mode === 'traffic' && trafficQ.data?.edges) {
      for (const te of trafficQ.data.edges) {
        trafficMap.set(`${te.source}->${te.target}`, te.rps)
      }
    }

    let kindFiltered = raw.nodes.filter((n) => !hiddenKinds.has(n.kind))

    if (focusServiceId) {
      const keep = connectedIds(focusServiceId, raw.edges)
      // Always keep the focused service even if kind filter hid it.
      kindFiltered = raw.nodes.filter(
        (n) => keep.has(n.id) && (n.id === focusServiceId || !hiddenKinds.has(n.kind)),
      )
    }

    const idSet = new Set(kindFiltered.map((n) => n.id))
    const visibleEdges = raw.edges.filter((e) => idSet.has(e.source) && idSet.has(e.target))

    // Fan-out: Service→Workload RPS onto owner edges (Deployment/STS/DS → Pod)
    // so the path keeps flowing instead of stopping at the workload.
    if (mode === 'traffic' && trafficMap.size > 0) {
      const inbound = new Map<string, number>()
      for (const te of trafficQ.data?.edges || []) {
        if (!idSet.has(te.target)) continue
        inbound.set(te.target, (inbound.get(te.target) || 0) + te.rps)
      }
      const ownerOutCount = new Map<string, number>()
      for (const e of visibleEdges) {
        if (e.kind !== 'owner') continue
        ownerOutCount.set(e.source, (ownerOutCount.get(e.source) || 0) + 1)
      }
      for (const e of visibleEdges) {
        if (e.kind !== 'owner') continue
        const key = `${e.source}->${e.target}`
        if (trafficMap.has(key)) continue
        const total = inbound.get(e.source) || 0
        const n = ownerOutCount.get(e.source) || 1
        if (total > 0) trafficMap.set(key, total / n)
      }
    }

    const rfNodes: Node[] = kindFiltered.map((n) => ({
      id: n.id,
      type: 'topo',
      position: { x: 0, y: 0 },
      data: { ...n },
    }))

    const rfEdges: Edge[] = visibleEdges.map((e) => {
      const rps = trafficMap.get(`${e.source}->${e.target}`)
      const showTraffic = mode === 'traffic' && rps != null
      return {
        id: e.id,
        source: e.source,
        target: e.target,
        type: 'smoothstep',
        className: showTraffic ? 'topo-edge-traffic' : undefined,
        animated: showTraffic && (rps || 0) > 5,
        label: showTraffic ? `${Number((rps || 0).toFixed(1))} rps` : undefined,
        labelShowBg: true,
        labelBgPadding: [6, 4] as [number, number],
        labelBgBorderRadius: 4,
        labelBgStyle: {
          fill: 'var(--color-panel-solid)',
          stroke: 'color-mix(in srgb, var(--color-cyan) 35%, var(--color-line))',
          strokeWidth: 1,
        },
        labelStyle: {
          fill: 'var(--color-cyan)',
          fontSize: 10,
          fontWeight: 650,
          fontFamily: 'var(--font-mono)',
        },
        markerEnd: { type: MarkerType.ArrowClosed, width: 14, height: 14 },
        style: {
          strokeWidth: showTraffic ? Math.min(1.2 + (rps || 0) / 40, 4) : 1.25,
          stroke: showTraffic
            ? 'color-mix(in srgb, var(--color-cyan) 70%, var(--color-line))'
            : 'var(--color-line)',
        },
      } as Edge
    })

    // Wider ranks in traffic mode so edge labels sit mid-path, not on arrowheads.
    const laid = focusServiceId
      ? layoutTopology(rfNodes, rfEdges, 'LR', mode === 'traffic' ? 110 : 72)
      : layoutTopologyByGroups(rfNodes, rfEdges, 'LR', mode === 'traffic' ? 110 : 72)
    return { nodes: laid.nodes, edges: laid.edges }
  }, [graphQ.data, trafficQ.data, hiddenKinds, focusServiceId, mode])

  const toggleKind = (kind: string) => {
    setHiddenKinds((prev) => {
      const next = new Set(prev)
      if (next.has(kind)) next.delete(kind)
      else next.add(kind)
      return next
    })
  }

  const toggleService = (id: string) => {
    setFocusServiceId((prev) => (prev === id ? null : id))
  }

  const onNodeClick = useCallback(
    (_: MouseEvent, node: Node) => {
      const href = (node.data as TopologyNode)?.href
      if (href) {
        navigate(href, { state: { from: '/topology' } })
      }
    },
    [navigate],
  )

  const countOf = (kind: string) =>
    graphQ.data?.counts?.find((c) => c.kind === kind)?.count ?? 0

  return (
    <div className="topo-page">
      <PageHeader
        title={t('topology.title')}
        subtitle={t('topology.subtitle')}
        action={
          <div className="flex flex-wrap items-center gap-2">
            <Button
              variant="ghost"
              type="button"
              onClick={() => {
                void graphQ.refetch()
                if (mode === 'traffic') void trafficQ.refetch()
              }}
            >
              <RefreshCw className={cn('h-3.5 w-3.5', graphQ.isFetching && 'animate-spin')} />
              {t('common.refresh')}
            </Button>
          </div>
        }
      />

      {!nsReady ? (
        <EmptyState>
          {t('topology.pickNamespace')}{' '}
          <span className="text-text-dim">({t('topology.pickNamespaceHint')})</span>
        </EmptyState>
      ) : (
        <div className="topo-shell">
          <aside className="topo-filters">
            <div className="topo-filters-title">{t('topology.filters')}</div>
            <div className="topo-filter-section">
              <div className="topo-filter-label">{t('topology.viewMode')}</div>
              <div className="topo-seg">
                <button
                  type="button"
                  className={cn(mode === 'resources' && 'is-active')}
                  onClick={() => setMode('resources')}
                >
                  {t('topology.modeResources')}
                </button>
                <button
                  type="button"
                  className={cn(mode === 'traffic' && 'is-active')}
                  onClick={() => setMode('traffic')}
                >
                  {t('topology.modeTraffic')}
                </button>
              </div>
            </div>
            <div className="topo-filter-section">
              <div className="topo-filter-label topo-filter-label-row">
                <span>
                  {t('topology.services')} · {services.length}
                </span>
                {focusServiceId ? (
                  <button
                    type="button"
                    className="topo-link-btn"
                    onClick={() => setFocusServiceId(null)}
                  >
                    {t('topology.showAllGroups')}
                  </button>
                ) : null}
              </div>
              {services.length > 6 ? (
                <label className="topo-search">
                  <Search className="h-3 w-3 shrink-0 opacity-50" />
                  <input
                    value={serviceQuery}
                    onChange={(e) => setServiceQuery(e.target.value)}
                    placeholder={t('topology.serviceSearch')}
                  />
                </label>
              ) : null}
              <div className="topo-group-list topo-service-list">
                {filteredServices.length === 0 ? (
                  <div className="topo-empty-hint">{t('topology.noServices')}</div>
                ) : (
                  filteredServices.map((s) => {
                    const focused = focusServiceId === s.id
                    return (
                      <button
                        key={s.id}
                        type="button"
                        className={cn('topo-group-item', focused && 'is-focused')}
                        title={t('topology.serviceClickHint')}
                        onClick={() => toggleService(s.id)}
                      >
                        <span className="truncate">{s.name}</span>
                      </button>
                    )
                  })
                )}
              </div>
            </div>
            <div className="topo-filter-section">
              <div className="topo-filter-label">{t('topology.kinds')}</div>
              <div className="topo-kind-list">
                {KIND_ORDER.map((kind) => {
                  const n = countOf(kind)
                  if (!graphQ.data && n === 0) return null
                  if (graphQ.data && n === 0) return null
                  const on = !hiddenKinds.has(kind)
                  return (
                    <button
                      key={kind}
                      type="button"
                      className={cn('topo-kind-chip', on && 'is-on')}
                      onClick={() => toggleKind(kind)}
                    >
                      <span>{kind}</span>
                      <em>{n}</em>
                    </button>
                  )
                })}
              </div>
            </div>
            <div className="topo-filter-section">
              <div className="topo-filter-label">{t('topology.groupBy')}</div>
              <div className="topo-seg">
                <button
                  type="button"
                  className={cn(groupBy === 'app' && 'is-active')}
                  onClick={() => setGroupBy('app')}
                >
                  {t('topology.groupApp')}
                </button>
                <button
                  type="button"
                  className={cn(groupBy === 'namespace' && 'is-active')}
                  onClick={() => setGroupBy('namespace')}
                >
                  {t('topology.groupNamespace')}
                </button>
              </div>
            </div>
            {mode === 'traffic' && trafficQ.data ? (
              <div className="topo-traffic-note">
                <div className="topo-traffic-mode">
                  {t('topology.trafficMode')}: <code>{trafficQ.data.mode}</code>
                </div>
                <p>
                  {trafficQ.data.mode === 'prometheus'
                    ? t('topology.trafficHintProm')
                    : trafficQ.data.mode === 'showcase'
                      ? t('topology.trafficHintShowcase')
                      : t('topology.trafficHintSynthetic')}
                </p>
              </div>
            ) : null}
          </aside>

          <div className="topo-canvas-wrap">
            {graphQ.isLoading ? (
              <div className="topo-loading">
                <Loader2 className="h-5 w-5 animate-spin text-cyan" />
                {t('common.loading')}
              </div>
            ) : graphQ.isError ? (
              <EmptyState>{(graphQ.error as Error)?.message || t('topology.loadFailed')}</EmptyState>
            ) : nodes.length === 0 ? (
              <EmptyState>{t('topology.empty')}</EmptyState>
            ) : (
              <ReactFlow
                key={`${focusServiceId || 'all'}-${mode}-${nodes.length}`}
                nodes={nodes}
                edges={edges}
                nodeTypes={nodeTypes}
                onNodeClick={onNodeClick}
                fitView
                fitViewOptions={{ padding: 0.2 }}
                minZoom={0.2}
                maxZoom={1.6}
                proOptions={{ hideAttribution: true }}
              >
                <Background gap={18} color="var(--color-grid-line)" />
                <Controls showInteractive={false} position="bottom-left" />
                <MiniMap
                  pannable
                  zoomable
                  position="bottom-right"
                  maskColor="color-mix(in srgb, var(--color-bg) 55%, transparent)"
                  nodeColor={(n) => {
                    const st = (n.data as TopologyNode)?.status
                    if (st === 'danger') return '#e11d2e'
                    if (st === 'warn') return '#e6b400'
                    if (st === 'ok') return '#1f9d55'
                    return '#6f98a3'
                  }}
                />
              </ReactFlow>
            )}
            {graphQ.data?.truncated ? (
              <div className="topo-truncated">{t('topology.truncated')}</div>
            ) : null}
            <div className="topo-hint">
              {t('topology.clickHint')}
              {focusServiceId ? null : (
                <>
                  {' · '}
                  <Link to="/pods" className="text-cyan hover:underline">
                    {t('nav.pods')}
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      )}

    </div>
  )
}
