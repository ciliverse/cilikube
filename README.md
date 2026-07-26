<div align="center">
  <img alt="CiliKube Logo" width="200" height="200" src="docs/logo.png">
  <h1>CiliKube</h1>
  <span>English | <a href="./README.zh-CN.md">中文</a></span>
</div>

<div align="center">
  <img src="https://img.shields.io/badge/Release-v0.9.0-green?style=flat-square" alt="Release v0.9.0">
  <img src="https://img.shields.io/badge/Frontend-React%2019-61DAFB?style=flat-square&logo=react" alt="React 19">
  <img src="https://img.shields.io/badge/Frontend-TypeScript%207-blue?style=flat-square&logo=typescript" alt="TypeScript 7">
  <img src="https://img.shields.io/badge/Frontend-Vite%208-blue?style=flat-square&logo=vite" alt="Vite 8">
  <img src="https://img.shields.io/badge/Frontend-Tailwind%204-38BDF8?style=flat-square&logo=tailwindcss" alt="Tailwind 4">
  <img src="https://img.shields.io/badge/Backend-Go%201.26+-blue?style=flat-square&logo=go" alt="Go 1.26+">
  <img src="https://img.shields.io/badge/Backend-Gin-blue?style=flat-square&logo=gin" alt="Gin">
  <img src="https://img.shields.io/badge/Kubernetes-1.36.2-blue?style=flat-square&logo=kubernetes" alt="Kubernetes 1.36.2">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue?style=flat-square" alt="License: Apache 2.0">
  <img src="https://img.shields.io/github/stars/ciliverse/cilikube?style=social" alt="GitHub Stars">
  <img src="https://img.shields.io/github/forks/ciliverse/cilikube?style=social" alt="GitHub Forks">
</div>

## 🌟 Project Support

We appreciate your interest in CiliKube. If you find this project valuable for your Kubernetes management needs, please consider starring the repository ⭐. Community support drives continuous development and improvement.

Stay updated with the latest releases and technical insights by following our WeChat Official Account **cilliantech**.

## 🤝 Contributors

<a href="https://github.com/ciliverse/cilikube/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ciliverse/cilikube" />
</a>

We extend our gratitude to all contributors who have helped improve CiliKube through code contributions, bug reports, and feature suggestions.

## 🏢 Sponsorship

This project's CDN acceleration and security protection services are generously sponsored by Tencent EdgeOne.

<a href="https://edgeone.ai/zh?from=github">
  <img src="https://edgeone.ai/media/34fe3a45-492d-4ea4-ae5d-ea1087ca7b4b.png" alt="EdgeOne" width="350" height="50">
</a>

## 📖 Overview

CiliKube is an enterprise-grade, open-source Kubernetes multi-cluster management platform built with modern web technologies including React, TypeScript, Go, and Gin. The platform provides an intuitive, streamlined interface for comprehensive Kubernetes resource management while maintaining extensibility for custom requirements. CiliKube serves as an ideal foundation for organizations seeking efficient cluster operations and developers learning cloud-native technologies.

### v0.9.0 (current)

Highlights since v0.8.0: Gateway API (GatewayClass / Gateway / HTTPRoute), Chinese UI + Maple Mono CN, desktop apps (Windows / macOS / Linux), forced password change on first desktop login.

<div align="center">
  <img src="docs/v0.8/overview.png" alt="Cluster Overview" width="100%">
  <p><strong>Cluster Overview</strong></p>
</div>

<div align="center">
  <img src="docs/v0.8/admin-users.png" alt="Admin Users" width="100%">
  <p><strong>Admin · Users</strong></p>
</div>

<div align="center">
  <img src="docs/v0.8/admin-settings.png" alt="System Settings" width="100%">
  <p><strong>Admin · Settings</strong></p>
</div>

<details>
<summary>Cluster overview across themes (TRON / Paper / Matrix / Amber / Nord)</summary>

<br>

<div align="center">
  <img src="docs/v0.8/overview-tron.png" alt="TRON theme overview" width="100%">
  <p><strong>TRON</strong></p>
</div>

<div align="center">
  <img src="docs/v0.8/overview-paper.png" alt="Paper theme overview" width="100%">
  <p><strong>Paper</strong></p>
</div>

<div align="center">
  <img src="docs/v0.8/overview-matrix.png" alt="Matrix theme overview" width="100%">
  <p><strong>Matrix</strong></p>
</div>

<div align="center">
  <img src="docs/v0.8/overview-amber.png" alt="Amber theme overview" width="100%">
  <p><strong>Amber</strong></p>
</div>

<div align="center">
  <img src="docs/v0.8/overview-nord.png" alt="Nord theme overview" width="100%">
  <p><strong>Nord</strong></p>
</div>

</details>


## ✨ Key Differentiators

CiliKube distinguishes itself from complex enterprise solutions by prioritizing simplicity and usability without sacrificing functionality:

1. **Streamlined Interface**: Provides an intuitive, clean interface for essential Kubernetes resource management operations.
2. **Developer-Centric Design**: Built with modern development practices and clean architecture, making it an excellent reference for **React/Go web development** and **Kubernetes API integration**.
3. **Extensible Architecture**: Designed with modularity in mind, enabling seamless integration of custom features and workflows.

## 🎯 Target Audience

- **Frontend Developers**: Seeking hands-on experience with **React + TypeScript + Tailwind** ecosystem
- **Backend Developers**: Learning **Go + Gin** web development and microservices architecture
- **Cloud-Native Engineers**: Exploring **Kubernetes API** integration and **client-go** library implementation
- **DevOps Teams**: Requiring a lightweight, customizable Kubernetes management interface
- **Educational Institutions**: Teaching modern web development and cloud-native technologies

## 💡 Project Genesis

CiliKube emerged from a comprehensive full-stack development learning initiative, combining practical web development skills with deep Kubernetes expertise. The project represents both a technical achievement and an educational resource, designed to serve as a gateway for developers entering the cloud-native ecosystem. Our mission extends beyond providing a management tool—we aim to foster a community of learners and contributors in the open-source landscape.

## 🌐 Online Demo

- Online Demo: https://cilikube.cillian.website

## 📚 Documentation

- Online demo / product site: [cilikube.cillian.website](https://cilikube.cillian.website)

## 🚀 Technology Stack

CiliKube leverages industry-standard technologies and frameworks to ensure reliability, maintainability, and developer productivity.

**System Requirements**:
- Node.js >= 20.0.0 (Developed and tested with v24.14.1)
- Go >= 1.26.0 (Developed and tested with v1.26.4)
- PNPM >= 10.x (Package management)
- Kubernetes / kubectl 1.36.2 (client-go aligned)

**Frontend Architecture**: 
- **Core**: `React 19` `TypeScript 7` `Vite 8` `Tailwind CSS 4`
- **Data / Routing**: `TanStack Query` `React Router`
- **Motion / Charts**: `Framer Motion` `Recharts`
- **HTTP Client**: `Axios`
- The previous Vue UI remains available in the separate [cilikube-web](https://github.com/cillianxtech/cilikube-web) repository.

**Backend Architecture**: 
- **Framework**: `Go` `Gin`
- **Kubernetes Integration**: `client-go`
- **Authentication**: `JWT`
- **Real-time Communication**: `Gorilla WebSocket`
- **Configuration**: `Viper`
- **Logging**: `Zap Logger`

## ✨ Current Features

- **Auth & RBAC**: JWT login, optional GitHub OAuth, Casbin roles; user / role / settings admin; security audit log
- **Multi-cluster**: Import and switch clusters, local kubeconfig context import, active cluster selection
- **Overview & search**: Cluster overview dashboard, global resource search, Events stream
- **Cluster resources**:
  - Nodes / Namespaces / CRDs
  - Workloads: Pods (logs / web terminal), Deployments, StatefulSets, DaemonSets, Jobs, CronJobs, HPA, PDB
  - Network: Services, Ingress, GatewayClasses, Gateways, HTTPRoutes, NetworkPolicies
  - Config: ConfigMaps, Secrets, ServiceAccounts, ResourceQuotas, LimitRanges
  - Storage: PV / PVC / StorageClass
  - Access: Roles, RoleBindings, ClusterRoles, ClusterRoleBindings
- **Observe & ops**: Monitoring (Prometheus metrics), Helm releases, API Proxy console, resource YAML view / edit

## 🖥️ Desktop (Windows / macOS / Linux)

Double-click app: Electron shell + embedded Go API, using your local kubeconfig (`~/.kube/config` or `%USERPROFILE%\.kube\config`).

Latest: **[v0.9.2-desktop.1](https://github.com/ciliverse/cilikube/releases/tag/v0.9.2-desktop.1)**

| Platform | Asset |
| --- | --- |
| Windows x64 | [setup](https://github.com/ciliverse/cilikube/releases/download/v0.9.2-desktop.1/CiliKube-0.9.2-win-x64-setup.exe) / [portable](https://github.com/ciliverse/cilikube/releases/download/v0.9.2-desktop.1/CiliKube-0.9.2-win-x64-portable.exe) |
| macOS Apple Silicon | [DMG](https://github.com/ciliverse/cilikube/releases/download/v0.9.2-desktop.1/CiliKube-0.9.2-mac-arm64.dmg) |
| Linux x64 | [AppImage](https://github.com/ciliverse/cilikube/releases/download/v0.9.2-desktop.1/CiliKube-0.9.2-linux-x86_64.AppImage) |

- **First login**: `admin` / `12345678` (forced password change on first login)
- **Notes**: builds are unsigned (Windows SmartScreen / macOS Gatekeeper may warn). Intel Mac DMG not yet published.
- **Build**: push a desktop tag (CI builds all three platforms; SQLite needs CGO on each runner):

```bash
git tag v0.9.2-desktop.1
git push origin v0.9.2-desktop.1
```

## 💻 Local Development

### Environment Preparation
1. Install [Node.js](https://nodejs.org/) (>=20) and [pnpm](https://pnpm.io/) (>=10)
2. Install [Go](https://go.dev/) (>=1.26)
3. Have a Kubernetes cluster and configure the kubeconfig file (defaults to reading `~/.kube/config`)

### Running the Frontend
```bash
# Navigate to the frontend directory (monorepo path: web/)
cd web
# Install dependencies
pnpm install
# Start the development server
pnpm dev
```

Visit http://localhost:8888 to see the frontend interface.

### Running the Backend
```bash
# From the repository root
# (Optional) Update Go dependencies
go mod tidy
# Run the backend service (listens on port 8080 by default)
# Configuration files are modified in configs/config.yaml
go run cmd/server/main.go
```

### Building the Project
```bash
# Build frontend production package (output to web/dist)
cd web
pnpm build

# Build backend executable
cd ..
go build -o bin/cilikube cmd/server/main.go
```

## 🐳 Docker Deployment

### Using Official Images
```bash
# Backend
docker run -d --name cilikube -p 8080:8080 -v ~/.kube:/root/.kube:ro cilliantech/cilikube:latest

# Frontend
docker run -d --name cilikube-web -p 80:80 cilliantech/cilikube-web:latest
```

### Using Docker Compose
```bash
docker-compose up -d
```

Visit http://localhost to access the interface.

## ☸️ Kubernetes Deployment (Helm)

### Environment Preparation
- Install Helm (>=3.0)
- Have a Kubernetes cluster and configure the kubeconfig file
- Install kubectl (>=1.36.2)

### Deployment Steps
```bash
# Add Helm repository
helm repo add ciliverse https://charts.cillian.website

# Update Helm repository
helm repo update

# Install CiliKube
helm install cilikube ciliverse/cilikube -n cilikube --create-namespace

# Check service status
kubectl get svc cilikube -n cilikube
```

## 🎨 Previous UI (Vue)

<details>
<summary>Click to view historical Vue screenshots</summary>

<table>
  <tr>
    <td width="50%">
      <img src="docs/login.png" alt="Login" width="100%">
      <p align="center"><strong>Login Interface</strong></p>
    </td>
    <td width="50%">
      <img src="docs/dashboard.png" alt="Dashboard" width="100%">
      <p align="center"><strong>Dashboard Overview</strong></p>
    </td>
  </tr>
  <tr>
    <td width="50%">
      <img src="docs/cilikube12.png" alt="Navigation" width="100%">
      <p align="center"><strong>Navigation Menu</strong></p>
    </td>
    <td width="50%">
      <img src="docs/cluster.png" alt="Cluster" width="100%">
      <p align="center"><strong>Cluster Management</strong></p>
    </td>
  </tr>
  <tr>
    <td width="50%">
      <img src="docs/pod.png" alt="Pods" width="100%">
      <p align="center"><strong>Pod Management</strong></p>
    </td>
    <td width="50%">
      <img src="docs/shell.png" alt="Shell" width="100%">
      <p align="center"><strong>Web Terminal</strong></p>
    </td>
  </tr>
</table>

</details>

## 🤝 Contribution Guide

We welcome contributions of all forms! If you'd like to help improve CiliKube, please:

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'feat: Add some AmazingFeature'`) - Please follow the Git Commit Guidelines
4. Push your branch to your fork (`git push origin feature/AmazingFeature`)
5. Submit a Pull Request

### Git Commit Guidelines

Please follow the Conventional Commits specification:

- `feat`: Add new features
- `fix`: Fix issues/bugs
- `perf`: Optimize performance
- `style`: Change the code style without affecting the running result
- `refactor`: Refactor code
- `revert`: Revert changes
- `test`: Test related, does not involve changes to business code
- `docs`: Documentation and Annotation
- `chore`: Updating dependencies/modifying scaffolding configuration, etc.
- `workflow`: Workflow Improvements
- `ci`: CICD related changes
- `types`: Type definition changes
- `wip`: Work in progress (should generally not be merged)

## 📞 Contact

- Email: cilliantech@gmail.com
- Website: https://www.cillian.website
- WeChat: cillianoffical

<img src="docs/wechat400x400.png" width="100" height="100" />

## 📜 License

This project is open-sourced under the Apache 2.0 License

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)
