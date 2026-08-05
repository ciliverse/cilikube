import dagre from '@dagrejs/dagre'
import type { Edge, Node } from '@xyflow/react'
import type { TopologyNode } from '@/api/topology'

const NODE_W = 200
const NODE_H = 72
const GROUP_GAP_Y = 64
const GROUP_GAP_X = 48

function layoutOne(
  nodes: Node[],
  edges: Edge[],
  direction: 'LR' | 'TB',
  ranksep = 72,
): { nodes: Node[]; width: number; height: number } {
  const g = new dagre.graphlib.Graph()
  g.setDefaultEdgeLabel(() => ({}))
  g.setGraph({ rankdir: direction, nodesep: 36, ranksep, marginx: 24, marginy: 24 })

  for (const n of nodes) {
    g.setNode(n.id, { width: NODE_W, height: NODE_H })
  }
  for (const e of edges) {
    g.setEdge(e.source, e.target)
  }
  dagre.layout(g)

  let minX = Infinity
  let minY = Infinity
  let maxX = -Infinity
  let maxY = -Infinity

  const laid = nodes.map((n) => {
    const p = g.node(n.id)
    const x = (p?.x || 0) - NODE_W / 2
    const y = (p?.y || 0) - NODE_H / 2
    minX = Math.min(minX, x)
    minY = Math.min(minY, y)
    maxX = Math.max(maxX, x + NODE_W)
    maxY = Math.max(maxY, y + NODE_H)
    return { ...n, position: { x, y } }
  })

  if (!Number.isFinite(minX)) {
    return { nodes: laid, width: 0, height: 0 }
  }

  // Normalize each subgraph to origin so we can stack/offset cleanly.
  const normalized = laid.map((n) => ({
    ...n,
    position: { x: n.position.x - minX, y: n.position.y - minY },
  }))

  return {
    nodes: normalized,
    width: maxX - minX,
    height: maxY - minY,
  }
}

/** Layout whole graph left-to-right (legacy). */
export function layoutTopology(
  nodes: Node[],
  edges: Edge[],
  direction: 'LR' | 'TB' = 'LR',
  ranksep = 72,
): { nodes: Node[]; edges: Edge[] } {
  const { nodes: laid } = layoutOne(nodes, edges, direction, ranksep)
  return { nodes: laid, edges }
}

/**
 * Layout each app/service group as its own subgraph, stacked vertically.
 * Avoids one giant horizontal strip when a namespace has many apps.
 */
export function layoutTopologyByGroups(
  nodes: Node[],
  edges: Edge[],
  direction: 'LR' | 'TB' = 'LR',
  ranksep = 72,
): { nodes: Node[]; edges: Edge[] } {
  if (nodes.length === 0) return { nodes, edges }

  const groupIds = new Map<string, string[]>()
  for (const n of nodes) {
    const g = ((n.data as TopologyNode)?.group || '_ungrouped').trim() || '_ungrouped'
    const list = groupIds.get(g) || []
    list.push(n.id)
    groupIds.set(g, list)
  }

  const sortedGroups = [...groupIds.keys()].sort((a, b) => a.localeCompare(b))
  if (sortedGroups.length <= 1) {
    return layoutTopology(nodes, edges, direction, ranksep)
  }

  const nodeById = new Map(nodes.map((n) => [n.id, n]))
  const out: Node[] = []
  let offsetY = 0

  for (const group of sortedGroups) {
    const ids = new Set(groupIds.get(group) || [])
    const subNodes = [...ids].map((id) => nodeById.get(id)!).filter(Boolean)
    const subEdges = edges.filter((e) => ids.has(e.source) && ids.has(e.target))
    const { nodes: laid, height } = layoutOne(subNodes, subEdges, direction, ranksep)

    for (const n of laid) {
      out.push({
        ...n,
        position: { x: n.position.x + GROUP_GAP_X, y: n.position.y + offsetY },
      })
    }
    offsetY += Math.max(height, NODE_H) + GROUP_GAP_Y
  }

  // Cross-group edges keep endpoints where laid; React Flow draws them across stacks.
  return { nodes: out, edges }
}
