# CiliKube Deployment

This directory contains deployment configurations and scripts for CiliKube across different platforms and environments.

## Directory Structure

```
deployments/
├── docker/           # Docker deployment files
├── kubernetes/       # Kubernetes manifests
├── helm/            # Helm charts
└── monitoring/      # Monitoring configurations
```

## Deployment Options

### 1. Docker Deployment

**Quick Start:**
```bash
# Using Docker Compose (Recommended)
docker-compose up -d

# Using Docker directly
docker run -d --name cilikube -p 8080:8080 cilliantech/cilikube:latest
```

**Features:**
- Multi-service orchestration with Docker Compose
- Production-ready configuration
- Built-in monitoring stack (Prometheus + Grafana)
- Persistent data storage

**Documentation:** [docker/README.md](docker/README.md)

### 2. Kubernetes Deployment

**Quick Start:**
```bash
# Apply Kubernetes manifests
kubectl apply -f deployments/kubernetes/

# Or using Helm
helm install cilikube deployments/helm/cilikube
```

**Features:**
- High availability setup
- Auto-scaling capabilities
- Service mesh integration
- Cloud-native monitoring

**Documentation:** [kubernetes/README.md](kubernetes/README.md)

### 3. Helm Deployment

**Quick Start:**
```bash
# Add Helm repository
helm repo add cilikube https://charts.cillian.website

# Install CiliKube
helm install cilikube cilikube/cilikube
```

**Features:**
- Parameterized deployments
- Easy upgrades and rollbacks
- Multi-environment support
- Custom value configurations

**Documentation:** [helm/README.md](helm/README.md)

## Environment Support

| Environment | Docker | Kubernetes | Helm |
|-------------|--------|------------|------|
| Development | ✅ | ✅ | ✅ |
| Staging | ✅ | ✅ | ✅ |
| Production | ✅ | ✅ | ✅ |

## Prerequisites

### Common Requirements
- **Kubernetes Cluster**: v1.20+ (for K8s/Helm deployments)
- **Docker**: v20.10+ (for Docker deployments)
- **kubectl**: v1.20+ (for K8s management)
- **Helm**: v3.0+ (for Helm deployments)

### Resource Requirements

| Component | CPU | Memory | Storage |
|-----------|-----|--------|---------|
| Backend | 500m | 512Mi | - |
| Frontend | 100m | 128Mi | - |
| MySQL | 500m | 1Gi | 10Gi |
| Redis | 100m | 256Mi | 1Gi |
| Prometheus | 500m | 1Gi | 20Gi |

## Configuration

### Environment Variables
All deployments support configuration through environment variables. See individual deployment documentation for specific variables.

### Secrets Management
- **Docker**: Use `.env` files or Docker secrets
- **Kubernetes**: Use Kubernetes secrets and ConfigMaps
- **Helm**: Use Helm values and external secret operators

### Persistent Storage
- **Docker**: Docker volumes
- **Kubernetes**: PersistentVolumes with StorageClasses
- **Helm**: Configurable storage options

## Monitoring

All deployment methods include monitoring capabilities:

- **Metrics**: Prometheus for metrics collection
- **Visualization**: Grafana dashboards
- **Alerting**: AlertManager for notifications
- **Logging**: Structured logging with log aggregation

Access monitoring:
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

## Security

### Network Security
- TLS/SSL encryption for all communications
- Network policies for Kubernetes deployments
- Firewall rules for Docker deployments

### Authentication & Authorization
- JWT-based authentication
- RBAC integration with Kubernetes
- Role-based access control

### Secrets
- Encrypted secret storage
- Rotation policies
- External secret management integration

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check port usage
   netstat -tulpn | grep :8080
   ```

2. **Resource Constraints**
   ```bash
   # Check resource usage
   docker stats
   kubectl top nodes
   ```

3. **Network Issues**
   ```bash
   # Test connectivity
   curl -f http://localhost:8080/health
   ```

### Logs

```bash
# Docker
docker logs cilikube-backend

# Kubernetes
kubectl logs -f deployment/cilikube-backend

# Helm
kubectl logs -f -l app.kubernetes.io/name=cilikube
```

## Contributing

When adding new deployment configurations:

1. Follow the existing directory structure
2. Include comprehensive documentation
3. Add example configurations
4. Test across different environments
5. Update this main README

## Support

- **Documentation**: [cilikube.cillian.website](https://cilikube.cillian.website)
- **Issues**: [GitHub Issues](https://github.com/ciliverse/cilikube/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ciliverse/cilikube/discussions)