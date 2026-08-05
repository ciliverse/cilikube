# Timeline View (Radar-style, hybrid)

## Goal

Standalone **Timeline** page under Observe: resource rows × time axis with:

- **Status segments** (colored bars) from local SQLite samples
- **Event markers** from live Kubernetes Events
- Time windows, namespace / search filters, optional group-by-app

Inspired by [skyhook-io/radar](https://github.com/skyhook-io/radar) Timeline view.

## Decisions (approved)

| Topic | Choice |
|-------|--------|
| Scope | Full (segments + markers), not events-only MVP |
| History store | Local SQLite sampling |
| Write strategy | Hybrid: status **change + heartbeat**; Events stay live (not bulk-ingested) |
| Scan interval | 15s per active cluster |
| Heartbeat write | ≥ 30s since last sample for same resource (or status/reason change) |
| Retention | 7 days; background prune |
| Default window | 15m (also 1h, 6h) |
| Nav | Observe → Timeline (`/timeline`); keep existing Events list page |

## Non-goals (v1)

- Multi-cluster single merged axis (v1 = current/selected cluster only)
- Full Informer hot-path (list+poll sampler is enough)
- Export / download of timeline
- Istio / custom CRD status semantics beyond typed core workloads

## Data model

### Table `timeline_status_samples`

| Column | Type | Notes |
|--------|------|--------|
| `id` | integer PK | autoincrement |
| `cluster_id` | varchar(36) | indexed |
| `namespace` | varchar(253) | |
| `kind` | varchar(64) | lowercase: deployment, pod, … |
| `name` | varchar(253) | |
| `uid` | varchar(64) | optional |
| `app_group` | varchar(253) | `app.kubernetes.io/name` → `app` → `app.kubernetes.io/instance` → `k8s-app` → `_ungrouped` |
| `status` | varchar(32) | see enum below |
| `reason` | varchar(256) | short human hint |
| `sampled_at` | datetime | UTC; indexed |

Indexes:

- `(cluster_id, namespace, kind, name, sampled_at)`
- `(cluster_id, sampled_at)` for prune / window queries

Status enum: `healthy` | `rolling` | `degraded` | `unhealthy` | `unknown`.

Migration: GORM AutoMigrate via store/database layer (same path as other tables).

## Sampler

### Process

1. Background goroutine started with API process (when DB enabled).
2. Every **15s**, for each Active registered cluster with a working clientset:
   - List target kinds in all namespaces (or paginated; hard cap per kind to avoid blowups, e.g. 2000 pods/cluster).
   - Compute `status` + `reason` + `app_group`.
   - Load last sample per resource key from memory cache (or last row); write if status/reason changed **or** `now - last.sampled_at >= 30s`.
3. Every hour (or on ticker): `DELETE WHERE sampled_at < now - 7d`.
4. Showcase / unreachable clusters: skip or write synthetic rows so demo UI is non-empty.

### Target kinds (v1)

Deployment, StatefulSet, DaemonSet, Pod, Job, CronJob, Service.

### Status rules (summary)

| Kind | healthy | rolling | degraded | unhealthy |
|------|---------|---------|----------|-----------|
| Deployment / STS | Ready == Desired | Available progressing / updated replicas lag | Ready &lt; Desired, no Unavailable | UnavailableReplicas &gt; 0 or Failed condition |
| DaemonSet | NumberReady == Desired | Updating | Ready &lt; Desired | — |
| Pod | Running + all Ready | — | Pending / not Ready | Failed, CrashLoopBackOff, ImagePullBackOff |
| Job | Complete | Active &gt; 0 | — | Failed |
| CronJob | exists, not suspended | — | suspended | — |
| Service | exists | — | — | — → mostly `healthy` / `unknown` |

Exact mapping lives in one Go helper (`ClassifyTimelineStatus`) shared by sampler and tests.

## API

### `GET /api/v1/timeline`

Query:

- `namespace` — optional; empty = all (subject to RBAC/cluster size caps)
- `from`, `to` — RFC3339; default `to=now`, `from=now-15m`
- `groupBy` — `app` \| `none` (default `app`)
- `q` — name / app substring
- `kinds` — comma list; default all sampler kinds

Response (sketch):

```json
{
  "clusterId": "...",
  "from": "...",
  "to": "...",
  "groupBy": "app",
  "groups": [
    {
      "name": "my-app",
      "rows": [
        {
          "kind": "deployment",
          "namespace": "ns",
          "name": "foo",
          "uid": "...",
          "href": "/deployments/...",
          "segments": [
            { "from": "...", "to": "...", "status": "healthy", "reason": "1/1 ready" }
          ],
          "events": [
            {
              "at": "...",
              "type": "Warning",
              "reason": "BackOff",
              "message": "...",
              "marker": "warning"
            }
          ],
          "eventBadges": { "modified": 2, "warning": 1 }
        }
      ]
    }
  ],
  "sampling": { "enabled": true, "heartbeatSeconds": 30, "retentionDays": 7 }
}
```

Segment build: order samples in `[from, to]`; extend each sample forward until next sample or `to`; if no samples, emit one segment from `from→to` using **live** classify of current object (and flag `provisional: true` in sampling meta).

Event markers: reuse `EventService` list filtered to involved objects / window; map:

| K8s | marker |
|-----|--------|
| type=Normal, reason≈Created/Scheduled | created |
| type=Normal, other | modified |
| type=Warning | warning |
| reason≈Deleted / Killing | deleted |

### `GET /api/v1/timeline/meta`

Returns sampler config, last success time per cluster, row counts, retention.

### Auth / Casbin

Same class as Events: `read:timeline` + viewer/editor GET on `/api/v1/timeline` and `/api/v1/timeline/meta`.

## Frontend

- Route `/timeline`, nav next to Events / Topology.
- Controls: window chips (15m / 1h / 6h), search, group-by-app toggle, kind filters (optional v1.1).
- Legend: event dots + status bar colors (align Radar: healthy green, rolling blue, degraded orange, unhealthy pink/red).
- Layout: sticky left resource column + scrollable SVG/HTML timeline; “Now” line.
- Hover: tooltip with reason / event message.
- Click row → existing detail `href`.
- Empty / cold-start copy: sampling will fill history after the first heartbeats.

i18n: `en` + `zh` keys under `timeline.*` and `nav.timeline`.

## Implementation order

1. Model + migrate + store helpers + prune
2. Sampler goroutine + status classifier + unit tests for classifier
3. Timeline handler/service (segments + events merge) + routes + Casbin
4. React Timeline page + CSS + nav/i18n
5. Showcase synthetic samples; manual smoke (one ns, 15m window)

## Config (optional yaml knobs)

Defaults hardcoded first; later:

```yaml
timeline:
  enabled: true
  scanSeconds: 15
  heartbeatSeconds: 30
  retentionDays: 7
  maxPodsPerCluster: 2000
```

## Spec self-review

- [x] No TBD placeholders for v1 scope
- [x] Storage vs Events responsibility split is explicit
- [x] Cold-start behavior defined (provisional live segment)
- [x] Caps/retention called out to avoid SQLite blowup
- [x] Matches approved hybrid approach and UI goals
- [x] Non-goals listed so Informer/multi-cluster do not sneak into v1
