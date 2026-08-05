# Topology View (Radar-style)

## Goal

Standalone **Topology** menu: Ingress → Service → Workload → Pod graph with filters, grouping, health colors, and a **Traffic** mode.

## Scope (C)

- Nav: Observe → Topology (`/topology`)
- Canvas: pan/zoom, group by `app` label or namespace
- Filters: resource kinds with counts
- Node click → existing detail routes
- Traffic: edge weight from Prometheus when configured; else showcase/synthetic rates from topology edges

## API

- `GET /api/v1/topology?namespace=&groupBy=app|namespace&kinds=`
- `GET /api/v1/topology/traffic?namespace=` → edge rates (`rps` / relative weight)

## Stack

- Backend: client-go typed lists + ownerReferences / selector / ingress backends
- Frontend: `@xyflow/react` + dagre layout

## Non-goals (later)

- Argo Rollouts / full mesh (Istio) topology
- Live packet capture
