<div align="center">
  <img alt="CiliKube Logo" width="150" height="150" src="public/logo.png">
  <h1>CiliKube Web</h1>
  <span>English | <a href="#中文">中文</a></span>
</div>

<div align="center">
  <img src="https://img.shields.io/badge/Frontend-Vue3-blue?style=flat-square&logo=vue.js" alt="Vue3">
  <img src="https://img.shields.io/badge/Frontend-TypeScript-blue?style=flat-square&logo=typescript" alt="TypeScript">
  <img src="https://img.shields.io/badge/Frontend-Vite-blue?style=flat-square&logo=vite" alt="Vite">
  <img src="https://img.shields.io/badge/Frontend-Element%20Plus-blue?style=flat-square&logo=element-plus" alt="Element Plus">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue?style=flat-square" alt="License: Apache 2.0">
</div>

> **Monorepo note:** This frontend now lives in the main [cilikube](https://github.com/cillianxtech/cilikube) repository under `web/`. Develop from the repo root with `cd web && pnpm dev` (API defaults to `http://localhost:8080`).

## 📖 Overview

CiliKube Web is the frontend interface for CiliKube, an enterprise-grade Kubernetes multi-cluster management platform. Built with Vue3, TypeScript, and Element Plus, it provides an intuitive, modern interface for comprehensive Kubernetes resource management.

## 🚀 Technology Stack

- **Core**: Vue3, TypeScript, Vite, Element Plus
- **State Management**: Pinia, Vue Router
- **HTTP Client**: Axios
- **Styling**: UnoCSS, Scss
- **Code Quality**: ESLint, Prettier

## 💻 Development

### Prerequisites
- Node.js >= 18.0.0
- PNPM >= 8.x

### Getting Started
```bash
# Install dependencies
pnpm install

# Start development server
pnpm dev

# Build for production
pnpm build
```

## 🐳 Docker Deployment

```bash
# Build image
docker build -t cilikube-web:latest .

# Run container
docker run -d --name cilikube-web -p 8888:8888 cilikube-web:latest
```

## 📚 Documentation

- Official Documentation: [cilikube.cillian.website](https://cilikube.cillian.website)
- Main Repository: [CiliKube](../cilikube)

---

## 中文

## 📖 产品概述

CiliKube Web 是 CiliKube 的前端界面，CiliKube 是一个企业级 Kubernetes 多集群管理平台。采用 Vue3、TypeScript 和 Element Plus 构建，为全面的 Kubernetes 资源管理提供直观、现代化的界面。

## 🚀 技术架构

- **核心技术**: Vue3, TypeScript, Vite, Element Plus
- **状态管理**: Pinia, Vue Router
- **HTTP 客户端**: Axios
- **样式系统**: UnoCSS, Scss
- **代码质量**: ESLint, Prettier

## 💻 本地开发

### 环境要求
- Node.js >= 18.0.0
- PNPM >= 8.x

### 快速开始
```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build
```

## 🐳 Docker 部署

```bash
# 构建镜像
docker build -t cilikube-web:latest .

# 运行容器
docker run -d --name cilikube-web -p 8888:8888 cilikube-web:latest
```

## 📚 文档

- 官方文档: [cilikube.cillian.website](https://cilikube.cillian.website)
- 主仓库: [CiliKube](../cilikube)