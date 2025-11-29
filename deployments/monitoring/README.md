# Monitoring Configuration

This directory contains monitoring configurations for CiliKube, including Prometheus, Grafana, and AlertManager setups.

## Overview

The monitoring stack provides comprehensive observability for CiliKube deployments:

- **Prometheus**: Metrics collection and storage
- **Grafana**: Visualization and dashboards
- **AlertManager**: Alert routing and notifications
- **Node Exporter**: System metrics collection
- **Application Metrics**: Custom CiliKube metrics

## Quick Start

### Docker Compose

```bash
# Start monitoring stack with CiliKube
docker-compose --profile monitoring up -d

# Access services
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3000 (admin/admin)
```

### Kubernetes

```bash
# Apply monitoring manifests
kubectl apply -f deployments/monitoring/kubernetes/

# Port forward to access services
kubectl port-forward svc/prometheus 9090:9090
kubectl port-forward svc/grafana 3000:3000
```

## Configuration Files

| File | Description |
|------|-------------|
| `prometheus.yml` | Prometheus configuration |
| `grafana/` | Grafana dashboards and datasources |
| `alertmanager/` | AlertManager configuration |
| `rules/` | Prometheus alerting rules |

## Prometheus Configuration

### Scrape Configs

The Prometheus configuration includes scrape configs for:

```yaml
scrape_configs:
  # CiliKube backend metrics
  - job_name: 'cilikube-backend'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s

  # Kubernetes API server
  - job_name: 'kubernetes-apiservers'
    kubernetes_sd_configs:
      - role: endpoints
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token

  # Kubernetes nodes
  - job_name: 'kubernetes-nodes'
    kubernetes_sd_configs:
      - role: node
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token

  # Kubernetes pods with annotations
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
```

### Recording Rules

Define recording rules for common queries:

```yaml
# rules/cilikube.yml
groups:
- name: cilikube.rules
  rules:
  - record: cilikube:request_duration_seconds:rate5m
    expr: rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])
  
  - record: cilikube:request_rate:rate5m
    expr: rate(http_requests_total[5m])
  
  - record: cilikube:error_rate:rate5m
    expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])
```

## Grafana Dashboards

### Pre-built Dashboards

1. **CiliKube Overview**
   - Application metrics
   - Request rates and latency
   - Error rates
   - Resource usage

2. **Kubernetes Cluster**
   - Node metrics
   - Pod metrics
   - Namespace overview
   - Resource quotas

3. **Infrastructure**
   - System metrics
   - Network metrics
   - Storage metrics
   - Database metrics

### Dashboard Configuration

```json
{
  "dashboard": {
    "title": "CiliKube Overview",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total[5m])) by (method, status)",
            "legendFormat": "{{method}} - {{status}}"
          }
        ]
      }
    ]
  }
}
```

## AlertManager Configuration

### Alert Rules

```yaml
# rules/alerts.yml
groups:
- name: cilikube.alerts
  rules:
  - alert: CiliKubeDown
    expr: up{job="cilikube-backend"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "CiliKube backend is down"
      description: "CiliKube backend has been down for more than 1 minute"

  - alert: HighErrorRate
    expr: cilikube:error_rate:rate5m > 0.1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value | humanizePercentage }} for the last 5 minutes"

  - alert: HighLatency
    expr: cilikube:request_duration_seconds:rate5m > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High latency detected"
      description: "Average response time is {{ $value }}s for the last 5 minutes"
```

### Notification Configuration

```yaml
# alertmanager/alertmanager.yml
global:
  smtp_smarthost: 'localhost:587'
  smtp_from: 'alerts@cilikube.example.com'

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'

receivers:
- name: 'web.hook'
  email_configs:
  - to: 'admin@cilikube.example.com'
    subject: 'CiliKube Alert: {{ .GroupLabels.alertname }}'
    body: |
      {{ range .Alerts }}
      Alert: {{ .Annotations.summary }}
      Description: {{ .Annotations.description }}
      {{ end }}

  slack_configs:
  - api_url: 'YOUR_SLACK_WEBHOOK_URL'
    channel: '#alerts'
    title: 'CiliKube Alert'
    text: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'
```

## Custom Metrics

### Application Metrics

CiliKube exposes custom metrics:

```go
// Example metrics in Go application
var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "status", "endpoint"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)
```

### Kubernetes Metrics

Monitor Kubernetes-specific metrics:

```yaml
# ServiceMonitor for Kubernetes
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cilikube-metrics
spec:
  selector:
    matchLabels:
      app: cilikube
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

## Performance Tuning

### Prometheus Optimization

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'production'
    region: 'us-west-2'

# Storage optimization
storage:
  tsdb:
    retention.time: 15d
    retention.size: 50GB
    wal-compression: true
```

### Resource Limits

```yaml
# Kubernetes resource limits
resources:
  requests:
    cpu: 500m
    memory: 1Gi
  limits:
    cpu: 2000m
    memory: 4Gi
```

## Security

### Authentication

```yaml
# Basic auth for Prometheus
basic_auth_users:
  admin: $2b$12$hNf2lSsxfm0.i4a.1kVpSOVyBCfIB51VRjgBUyv6kdnyTlgWj81Ay

# OAuth for Grafana
auth:
  generic_oauth:
    enabled: true
    client_id: cilikube
    client_secret: your-secret
    auth_url: https://auth.example.com/oauth/authorize
    token_url: https://auth.example.com/oauth/token
```

### Network Security

```yaml
# Network policies
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: monitoring-network-policy
spec:
  podSelector:
    matchLabels:
      app: prometheus
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: grafana
    ports:
    - protocol: TCP
      port: 9090
```

## Backup and Recovery

### Prometheus Data Backup

```bash
# Create snapshot
curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot

# Backup snapshot
tar -czf prometheus-backup-$(date +%Y%m%d).tar.gz \
  /prometheus/snapshots/20231201T120000Z-1234567890abcdef
```

### Grafana Configuration Backup

```bash
# Export dashboards
curl -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:3000/api/dashboards/db/cilikube-overview \
  > cilikube-dashboard.json
```

## Troubleshooting

### Common Issues

1. **Metrics not appearing**
   ```bash
   # Check Prometheus targets
   curl http://localhost:9090/api/v1/targets
   
   # Check service discovery
   curl http://localhost:9090/api/v1/label/__name__/values
   ```

2. **High memory usage**
   ```bash
   # Check Prometheus memory usage
   curl http://localhost:9090/api/v1/status/tsdb
   
   # Reduce retention or increase resources
   ```

3. **Alert not firing**
   ```bash
   # Check alert rules
   curl http://localhost:9090/api/v1/rules
   
   # Check AlertManager status
   curl http://localhost:9093/api/v1/status
   ```

### Debug Commands

```bash
# Check Prometheus configuration
promtool check config prometheus.yml

# Check alert rules
promtool check rules rules/*.yml

# Test queries
curl -G http://localhost:9090/api/v1/query \
  --data-urlencode 'query=up{job="cilikube-backend"}'
```

## Integration

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Check monitoring health
  run: |
    curl -f http://prometheus:9090/-/healthy
    curl -f http://grafana:3000/api/health
```

### External Systems

```yaml
# Webhook integration
webhook_configs:
- url: 'http://external-system/webhook'
  http_config:
    bearer_token: 'your-token'
```