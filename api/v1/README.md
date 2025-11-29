# CiliKube API v1

This directory contains the API definitions and documentation for CiliKube v1.

## API Documentation

- **OpenAPI Specification**: `openapi.yaml` - Complete API specification in OpenAPI 3.0 format
- **API Version**: v1
- **Base URL**: `/api/v1`

## API Overview

CiliKube API provides comprehensive Kubernetes multi-cluster management capabilities including:

### Authentication
- User login/logout
- JWT token-based authentication
- User profile management
- Role-based access control

### Cluster Management
- Multi-cluster configuration
- Cluster switching
- Cluster health monitoring

### Kubernetes Resources
- Pod management and monitoring
- Service and Ingress management
- ConfigMap and Secret management
- Deployment, StatefulSet, DaemonSet operations
- Storage management (PV/PVC)
- RBAC management

### Real-time Features
- Pod logs streaming (WebSocket)
- Pod shell access (WebSocket)
- Resource monitoring and metrics

## Authentication

All API endpoints (except login) require JWT authentication:

```
Authorization: Bearer <jwt_token>
```

## API Endpoints Structure

```
/api/v1/
├── auth/                 # Authentication endpoints
├── clusters/             # Cluster management
├── proxy/               # Kubernetes API proxy
├── summary/             # Dashboard and summary data
└── installer/           # Installation utilities
```

## Usage Examples

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

### List Clusters
```bash
curl -X GET http://localhost:8080/api/v1/clusters \
  -H "Authorization: Bearer <token>"
```

### Proxy to Kubernetes API
```bash
curl -X GET "http://localhost:8080/api/v1/proxy/api/v1/pods?clusterId=<cluster-id>" \
  -H "Authorization: Bearer <token>"
```

## WebSocket Endpoints

### Pod Logs
```
ws://localhost:8080/api/v1/proxy/api/v1/namespaces/{namespace}/pods/{name}/logs/ws?clusterId=<cluster-id>&container=<container>
```

### Pod Shell
```
ws://localhost:8080/api/v1/proxy/api/v1/namespaces/{namespace}/pods/{name}/exec/ws?clusterId=<cluster-id>&container=<container>
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### Error Response
```json
{
  "code": 400,
  "message": "error description"
}
```

## Development

When adding new API endpoints:

1. Update the OpenAPI specification in `openapi.yaml`
2. Implement handlers in `internal/handlers/`
3. Add routes in `internal/routes/`
4. Update this documentation

## Tools

- **Swagger UI**: Use tools like Swagger Editor to view and test the API
- **Postman**: Import the OpenAPI spec for API testing
- **curl**: Command-line testing examples provided above