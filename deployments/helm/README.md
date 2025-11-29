# Helm Deployment

This directory contains Helm charts for deploying CiliKube in Kubernetes clusters with parameterized configurations.

## Quick Start

### Using Helm Repository (Recommended)

```bash
# Add CiliKube Helm repository
helm repo add cilikube https://charts.cillian.website
helm repo update

# Install CiliKube
helm install cilikube cilikube/cilikube -n cilikube --create-namespace

# Upgrade CiliKube
helm upgrade cilikube cilikube/cilikube -n cilikube
```

### Using Local Charts

```bash
# Install from local chart
helm install cilikube ./deployments/helm/cilikube -n cilikube --create-namespace

# Install with custom values
helm install cilikube ./deployments/helm/cilikube -f custom-values.yaml -n cilikube --create-namespace
```

## Prerequisites

- Kubernetes cluster v1.20+
- Helm v3.0+
- kubectl configured to access your cluster

## Chart Structure

```
helm/
└── cilikube/
    ├── Chart.yaml          # Chart metadata
    ├── values.yaml         # Default configuration values
    ├── values-dev.yaml     # Development environment values
    ├── values-prod.yaml    # Production environment values
    └── templates/
        ├── deployment-backend.yaml
        ├── deployment-frontend.yaml
        ├── service-backend.yaml
        ├── service-frontend.yaml
        ├── configmap.yaml
        ├── secret.yaml
        ├── ingress.yaml
        ├── rbac.yaml
        ├── pvc.yaml
        └── hpa.yaml
```

## Configuration

### Default Values

The chart comes with sensible defaults in `values.yaml`:

```yaml
# Application configuration
app:
  name: cilikube
  version: latest

# Backend configuration
backend:
  image:
    repository: cilliantech/cilikube
    tag: latest
    pullPolicy: IfNotPresent
  replicas: 2
  resources:
    requests:
      cpu: 500m
      memory: 512Mi
    limits:
      cpu: 1000m
      memory: 1Gi

# Frontend configuration
frontend:
  image:
    repository: cilliantech/cilikube-web
    tag: latest
    pullPolicy: IfNotPresent
  replicas: 2
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi
```

### Environment-Specific Values

#### Development Environment

```bash
# Install for development
helm install cilikube ./deployments/helm/cilikube \
  -f ./deployments/helm/cilikube/values-dev.yaml \
  -n cilikube-dev --create-namespace
```

#### Production Environment

```bash
# Install for production
helm install cilikube ./deployments/helm/cilikube \
  -f ./deployments/helm/cilikube/values-prod.yaml \
  -n cilikube-prod --create-namespace
```

### Custom Configuration

Create your own values file:

```yaml
# custom-values.yaml
backend:
  replicas: 3
  resources:
    requests:
      cpu: 1000m
      memory: 1Gi
    limits:
      cpu: 2000m
      memory: 2Gi

ingress:
  enabled: true
  className: nginx
  hosts:
    - host: cilikube.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: cilikube-tls
      hosts:
        - cilikube.example.com

monitoring:
  prometheus:
    enabled: true
  grafana:
    enabled: true
```

## Features

### Auto-scaling

Enable Horizontal Pod Autoscaler:

```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80
```

### Persistence

Configure persistent storage:

```yaml
persistence:
  enabled: true
  storageClass: fast-ssd
  size: 10Gi
  accessMode: ReadWriteOnce
```

### Monitoring

Enable monitoring stack:

```yaml
monitoring:
  prometheus:
    enabled: true
    retention: 15d
    storageClass: fast-ssd
    storage: 20Gi
  
  grafana:
    enabled: true
    adminPassword: secure-password
    persistence:
      enabled: true
      size: 5Gi
```

### Security

Configure security settings:

```yaml
security:
  podSecurityContext:
    runAsNonRoot: true
    runAsUser: 1001
    fsGroup: 1001
  
  securityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    capabilities:
      drop:
        - ALL
  
  networkPolicy:
    enabled: true
```

## Installation Examples

### Minimal Installation

```bash
helm install cilikube cilikube/cilikube \
  --set backend.replicas=1 \
  --set frontend.replicas=1 \
  --set monitoring.prometheus.enabled=false \
  --set monitoring.grafana.enabled=false
```

### Production Installation

```bash
helm install cilikube cilikube/cilikube \
  --set backend.replicas=3 \
  --set frontend.replicas=2 \
  --set autoscaling.enabled=true \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=cilikube.example.com \
  --set monitoring.prometheus.enabled=true \
  --set monitoring.grafana.enabled=true \
  --set persistence.enabled=true
```

### High Availability Installation

```bash
helm install cilikube cilikube/cilikube \
  --set backend.replicas=5 \
  --set frontend.replicas=3 \
  --set autoscaling.enabled=true \
  --set autoscaling.maxReplicas=20 \
  --set persistence.enabled=true \
  --set persistence.storageClass=fast-ssd \
  --set monitoring.prometheus.enabled=true \
  --set monitoring.grafana.enabled=true \
  --set security.networkPolicy.enabled=true
```

## Upgrading

### Standard Upgrade

```bash
# Update repository
helm repo update

# Upgrade to latest version
helm upgrade cilikube cilikube/cilikube -n cilikube
```

### Upgrade with New Values

```bash
# Upgrade with custom values
helm upgrade cilikube cilikube/cilikube \
  -f new-values.yaml \
  -n cilikube
```

### Rollback

```bash
# List releases
helm history cilikube -n cilikube

# Rollback to previous version
helm rollback cilikube 1 -n cilikube
```

## Customization

### Adding Custom Resources

Create additional templates in the `templates/` directory:

```yaml
# templates/custom-configmap.yaml
{{- if .Values.customConfig.enabled }}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "cilikube.fullname" . }}-custom
  labels:
    {{- include "cilikube.labels" . | nindent 4 }}
data:
  custom.conf: |
    {{ .Values.customConfig.content | indent 4 }}
{{- end }}
```

### Custom Hooks

Add Helm hooks for custom operations:

```yaml
# templates/pre-install-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ include "cilikube.fullname" . }}-pre-install
  annotations:
    "helm.sh/hook": pre-install
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    spec:
      containers:
      - name: pre-install
        image: busybox
        command: ['sh', '-c', 'echo "Pre-install setup completed"']
      restartPolicy: Never
```

## Monitoring and Observability

### Prometheus Metrics

The chart automatically configures Prometheus scraping:

```yaml
# Automatic service monitor creation
monitoring:
  serviceMonitor:
    enabled: true
    interval: 30s
    path: /metrics
```

### Grafana Dashboards

Pre-configured dashboards are included:

```yaml
grafana:
  dashboards:
    enabled: true
    configMapName: cilikube-dashboards
```

## Troubleshooting

### Check Release Status

```bash
# List releases
helm list -n cilikube

# Get release status
helm status cilikube -n cilikube

# Get release values
helm get values cilikube -n cilikube
```

### Debug Templates

```bash
# Render templates without installing
helm template cilikube ./deployments/helm/cilikube \
  -f values.yaml \
  --debug

# Dry run installation
helm install cilikube ./deployments/helm/cilikube \
  --dry-run --debug
```

### Common Issues

1. **Values not applied**
   ```bash
   # Check current values
   helm get values cilikube -n cilikube
   ```

2. **Template errors**
   ```bash
   # Validate templates
   helm lint ./deployments/helm/cilikube
   ```

3. **Resource conflicts**
   ```bash
   # Check existing resources
   kubectl get all -n cilikube
   ```

## Uninstallation

```bash
# Uninstall release
helm uninstall cilikube -n cilikube

# Delete namespace (optional)
kubectl delete namespace cilikube
```

## Contributing

When modifying the Helm chart:

1. Update `Chart.yaml` version
2. Document changes in `CHANGELOG.md`
3. Test with different value combinations
4. Validate templates with `helm lint`
5. Update documentation