<div align="center">
  <img alt="CiliKube Logo" width="500" height="100" src="ui/public/logo.png">
  <h1>CiliKube</h1>
  <span><a href="./README.zh-CN.md">中文纯净版</a> | English</span>
</div>

![Alt](https://repobeats.axiom.co/api/embed/97bc0de802d8faf0f90512177cb349c4e494d76d.svg "Repobeats analytics image")
<!-- [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fciliverse%2Fcilikube.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fciliverse%2Fcilikube?ref=badge_shield) -->
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)



## 🤔 What is CiliKube? | CiliKube 是什么？

CiliKube is an open-source, full-stack Kubernetes (k8s) resource management platform built with modern, mainstream technology stacks (Vue3, TypeScript, Go, Gin). It aims to provide a clean and elegant interface to simplify the daily management (CRUD operations) of Kubernetes resources and support feature expansion.

CiliKube 是一个使用现代主流技术栈（Vue3, TypeScript, Go, Gin）构建的开源全栈 Kubernetes (k8s) 资源管理平台。它致力于提供一个简洁、优雅的界面，来简化 Kubernetes 资源的日常管理（增删改查）并支持功能拓展。
![alt text](ui/src/docs/architech.png)




## ✨ What Makes CiliKube Special? | CiliKube 的特色

Unlike complex systems pursuing "large and comprehensive" features, CiliKube focuses on being "small and beautiful." Its core goals are:

1.  **Core Functionality**: Offers a clear, intuitive interface for managing common K8s resources.
2.  **Learning-Friendly**: Features clean code structure and a modern tech stack, making it an excellent starting point for learning **Vue3/Go Web Development** and **custom Kubernetes tooling/development**.
3.  **Easy to Extend**: Designed with extensibility in mind, allowing users to easily add custom features based on their needs.

与追求“大而全”的复杂系统不同，CiliKube 专注于“小而美”。它的核心目标是：

1.  **核心功能**: 提供清晰、直观的界面来管理常用的 K8s 资源。
2.  **学习友好**: 代码结构清晰，技术栈现代，非常适合作为学习 **Vue3/Go Web 开发** 和 **Kubernetes 二次开发** 的入门项目。
3.  **易于拓展**: 预留了自定义功能的空间，方便用户根据自身需求进行扩展。



## 🎯 Target Users | 目标用户

*   Developers looking to learn **Vue3 + TypeScript + ElementPlus** frontend development.
*   Developers looking to learn **Go + Gin** backend development.
*   Cloud-native enthusiasts interested in using the **Kubernetes API** and **client-go**.
*   Teams or individuals needing a concise K8s management dashboard with the potential for customization.
*   希望学习 **Vue3 + TypeScript + ElementPlus** 前端开发的开发者。
*   希望学习 **Go + Gin** 后端开发的开发者。
*   对 **Kubernetes API** 和 **client-go** 使用感兴趣的云原生爱好者。
*   需要一个简洁 K8s 管理面板，并可能进行二次开发的团队或个人。

## 💡 Project Background | 项目背景

CiliKube originated from the author's practical project while learning full-stack web development. During the learning process, the author delved deep into Kubernetes and obtained relevant certifications. This project is not just a demonstration of learning outcomes but also aims to be a "Key," helping more learners like the author open the door to the open-source world, contribute, and grow together.

CiliKube 起源于作者学习 Web 全栈开发的实践项目。在学习过程中，作者深入探索了 Kubernetes，并获得了相关认证。这个项目不仅是学习成果的体现，更希望成为一把“钥匙 (Key)”，帮助更多像作者一样的学习者打开开源世界的大门，参与贡献，共同成长。

## 📚 Documentation | 文档
![alt text](ui/src/docs/docs.png)
*   Official Documentation: [cilikube.cillian.website](https://cilikube.cillian.website)
*   官方文档: [cilikube.cillian.website](https://cilikube.cillian.website) 


## 🌐 Online Demo | 在线预览

*   An online demo site is being deployed, stay tuned!
*   在线演示站点正在部署中，敬请期待！

## 🚀 Tech Stack | 技术栈

This project utilizes popular frontend and backend technology stacks, ensuring developers can engage with the latest tools and libraries.

*   **Environment Requirements (Recommended)**:
    *   Node.js >= 18.0.0 (Project developed with v22.14.0)
    *   Go >= 1.20 (Project developed with v1.24.2)
    *   PNPM >= 8.x

*   **Frontend**: `Vue3` `TypeScript` `Vite` `Element Plus` `Pinia` `Vue Router` `Axios` `UnoCSS` `Scss` `ESlint` `Prettier`
    *   Developed based on the excellent [v3-admin-vite](https://github.com/un-pany/v3-admin-vite) template, thanks to the original author un-pany.

*   **Backend**: `Go` `Gin` `Kubernetes client-go` `JWT (dgrijalva/jwt-go)` `Gorilla Websocket` `Logger (wonderivan/logger)`

本项目采用了当前流行的前后端技术栈，确保开发者能够接触和使用最新的工具和库。

*   **环境要求 (推荐)**:
    *   Node.js >= 18.0.0 (项目当前使用 v22.14.0 开发)
    *   Go >= 1.20 (项目当前使用 v1.24.2 开发)
    *   PNPM >= 8.x

*   **前端**: `Vue3` `TypeScript` `Vite` `Element Plus` `Pinia` `Vue Router` `Axios` `UnoCSS` `Scss` `ESlint` `Prettier`
    *   基于优秀的 [v3-admin-vite](https://github.com/un-pany/v3-admin-vite) 模板进行开发，感谢原作者 un-pany。

*   **后端**: `Go` `Gin` `Kubernetes client-go` `JWT (dgrijalva/jwt-go)` `Gorilla Websocket` `Logger (wonderivan/logger)`

## ✨ Main Features ([See Detailed Roadmap](#️-roadmap)) | 主要功能 ([查看详细开发计划](#️-开发计划-roadmap))

*   **User Authentication**: JWT-based login and authorization.
*   **Dashboard**: Cluster resource overview.
*   **Cluster Management**:
    *   Node Management
    *   Namespace Management
    *   Pod Management (List, Details, Logs, Terminal)
    *   Volume Management (PV/PVC)
    *   Configuration Management (ConfigMap/Secret)
    *   Network Management (Service/Ingress)
    *   Workload Management (Deployment/StatefulSet/DaemonSet) (Partially implemented)
*   **System Settings**: Theme switching, Internationalization (i18n) support (Planned).


*   **用户认证**: 基于 JWT 的登录和权限验证。
*   **仪表盘**: 集群资源概览。
*   **集群管理**:
    *   节点 (Node) 管理
    *   命名空间 (Namespace) 管理
    *   Pod 管理 (列表、详情、日志、终端)
    *   存储卷 (PV/PVC) 管理
    *   配置 (ConfigMap/Secret) 管理
    *   网络 (Service/Ingress) 管理
    *   工作负载 (Deployment/StatefulSet/DaemonSet) 管理 (部分实现)
*   **系统设置**: 主题切换、国际化支持 (计划中)。


## 🛠️ Roadmap | 开发计划

**Frontend**

*   [x] Login Page
*   [x] Basic Layout (Sidebar, Topbar, Tabs)
*   [x] Notifications
*   [ ] Workload Resource Pages (Deployment, StatefulSet, DaemonSet, etc.)
*   [ ] Configuration Management Pages (ConfigMap, Secret)
*   [ ] Network Resource Pages (Service, Ingress)
*   [ ] Storage Resource Pages (StorageClass, PV, PVC) - *PV/PVC partially done*
*   [ ] Access Control Pages (RBAC - ServiceAccount, Role, ClusterRoleBinding, etc.)
*   [ ] Events Viewer
*   [ ] Basic CRD Resource Management
*   [ ] Monitoring Integration (Display data from Prometheus/Grafana)
*   [ ] Log Viewer Enhancements
*   [ ] Web Shell Terminal Integration

**Backend**

*   [x] Kubernetes Client Initialization
*   [x] Basic Routing Setup (Gin)
*   [x] CORS Configuration
*   [x] JWT Authentication Middleware
*   [x] WebSocket Endpoint (for Logs and Web Shell)
*   [x] Node Resource API
*   [x] Pod Resource API (List, Get, Delete, Logs, Exec)
*   [x] PV/PVC Resource API
*   [ ] Namespace Resource API
*   [ ] Deployment / StatefulSet / DaemonSet Resource API
*   [ ] Service / Ingress Resource API
*   [ ] ConfigMap / Secret Resource API
*   [ ] RBAC Related Resource API
*   [ ] Event Resource API

**前端**

*   [x] 登录界面
*   [x] 基础布局 (侧边栏, 顶部导航, 标签栏)
*   [x] 消息通知
*   [x] 工作负载资源页面 (Deployment, StatefulSet, DaemonSet 等)
*   [ ] 配置管理页面 (ConfigMap, Secret)
*   [ ] 网络资源页面 (Service, Ingress)
*   [ ] 存储资源页面 (StorageClass, PV, PVC) - *PV/PVC 部分完成*
*   [ ] 访问控制页面 (RBAC - ServiceAccount, Role, ClusterRoleBinding 等)
*   [ ] 事件 (Events) 查看
*   [ ] CRD 资源管理 (基础)
*   [ ] 监控集成 (集成 Prometheus/Grafana 数据展示)
*   [ ] 日志查看页面优化
*   [ ] Web Shell 终端集成

**后端**

*   [x] Kubernetes 客户端初始化
*   [x] 基础路由设置 (Gin)
*   [x] CORS 跨域配置
*   [x] Node (节点) 资源接口
*   [x] Pod 资源接口 (列表, 详情, 删除, 日志, Exec)
*   [x] PV/PVC 资源接口
*   [x] Namespace 资源接口
*   [x] Deployment / StatefulSet / DaemonSet 资源接口
*   [x] Service / Ingress 资源接口
*   [x] ConfigMap / Secret 资源接口
*   [ ] RBAC 相关资源接口
*   [ ] Event 资源接口
*   [ ] JWT 认证中间件
*   [ ] WebSocket 接口 (用于日志和 Web Shell)
## 💻 Local Development | 本地开发

**Prerequisites**

1.  Install [Node.js](https://nodejs.org/) (>=18) and [pnpm](https://pnpm.io/)
2.  Install [Go](https://go.dev/) (>=1.20)
3.  Have a Kubernetes cluster and configure your `kubeconfig` file (reads `~/.kube/config` by default)

**Run Frontend**

```bash
# Enter frontend directory
cd ui

# Install dependencies
pnpm install

# Start development server
pnpm dev
```
**Run Backend**
```bash
# Enter backend directory
cd cmd/server/

# (Optional) Update Go dependencies
# go mod tidy

# Run backend service (listens on port 8081 by default)
go run main.go
```
**Build Project**
```bash
# Build frontend for production (output to ui/dist)
cd ui
pnpm build

# Build backend executable
cd ../server
go build -o cilikube-server main.go
```

**Run Tests (Frontend)**
```bash
cd ui
pnpm test:unit
```

**Lint Code (Frontend)**
```bash
cd ui
pnpm lint
```

**环境准备**

1.  安装 [Node.js](https://nodejs.org/) (>=18) 和 [pnpm](https://pnpm.io/)
2.  安装 [Go](https://go.dev/) (>=1.20)
3.  拥有一个 Kubernetes 集群，并配置好 `kubeconfig` 文件 (默认读取 `~/.kube/config`)


## 本地启动开发环境

```bash
# 启动后端
make run-backend

# 启动前端
make dev-ui
```



**运行前端**

```bash
# 进入前端目录
cd ui
# 安装依赖
pnpm install
# 启动开发服务器
pnpm dev
```

**运行后端**
```bash
# 进入后端目录
# cd cmd/server
# # (可选) 更新 Go 依赖
# go mod tidy
# # 运行后端服务 (默认监听 8081 端口)
# go run main.go

go run cmd/server/main.go
```

**Docker 部署**
```bash
docker run -d --name cilikube-server -p 8080:8080 cilliantech/cilikube:v0.1.0
```

**构建项目**
```bash
# 构建前端生产环境包 (输出到 ui/dist)
cd ui
pnpm build
# 构建后端可执行文件
cd ../server
go build -o cilikube-server main.go
```
**运行测试 (前端)**

```bash
cd ui
pnpm test:unit
```

**代码规范检查 (前端)**
```bash
cd ui
pnpm lint
```
## 🎨 Feature Preview | 功能预览
![alt text](ui/src/docs/login.png)
![alt text](ui/src/docs/first.png)
![alt text](ui/src/docs/minikube2.png)
![alt text](ui/src/docs/techstack.png)
![alt text](ui/src/docs/monitor1.png)
![alt text](ui/src/docs/monitor2.png)
![alt text](ui/src/docs/dashboard.png)
![alt text](ui/src/docs/nav.png)
![alt text](ui/src/docs/cluster.png)
![alt text](ui/src/docs/ingress.png)
![alt text](ui/src/docs/namespace.png)
![alt text](ui/src/docs/svc.png)
![alt text](ui/src/docs/deployment.png)
![alt text](ui/src/docs/pod.png)
![alt text](ui/src/docs/shell.png)
![alt text](ui/src/docs/pv.png)
![alt text](ui/src/docs/pvc.png)
![alt text](ui/src/docs/secret.png)
![alt text](ui/src/docs/configmap.png)




## 🤝 Contribution Guide | 贡献指南

We welcome contributions of all forms! If you'd like to help improve CiliKube, please:

Fork this repository

Create your feature branch (git checkout -b feature/AmazingFeature)

Commit your changes (git commit -m 'feat: Add some AmazingFeature') - Please follow the Git Commit Guidelines

Push your branch to your fork (git push origin feature/AmazingFeature)

Submit a Pull Request


我们欢迎各种形式的贡献！如果您想参与改进 CiliKube，请：

- Fork 本仓库

- 创建您的特性分支 (git checkout -b feature/AmazingFeature)

- 提交您的更改 (git commit -m 'feat: Add some AmazingFeature') - 请遵循 Git 提交规范

- 将您的分支推送到 Github (git push origin feature/AmazingFeature)

- 提交 Pull Request

## 🤝 Git Commit Guidelines | Git 提交规范

Please follow the Conventional Commits specification:

- feat: Add new features

- fix: Fix issues/bugs

- perf: Optimize performance

- style: Change the code style without affecting the running result

- refactor: Refactor code

- revert: Revert changes

- test: Test related, does not involve changes to business code

- docs: Documentation and Annotation

- chore: Updating dependencies/modifying scaffolding configuration, etc.

- workflow: Workflow Improvements

- ci: CICD related changes

- types: Type definition changes

- wip: Work in progress (should generally not be merged)



请遵循 Conventional Commits 规范：

- feat: 新增功能

- fix: 修复 Bug

- perf: 性能优化

- style: 代码样式调整（不影响逻辑）

- refactor: 代码重构
- 
- revert: 撤销更改

- test: 添加或修改测试

- docs: 文档或注释修改

- chore: 构建流程、依赖管理等杂项更改

- workflow: 工作流改进

- ci: CI/CD 配置相关

- types: 类型定义修改

- wip: 开发中的提交（不建议合入主分支）

## ❤️ Support the Project | 支持项目

Open source is not easy. If you find CiliKube helpful or inspiring, please consider giving it a Star ⭐! Your encouragement is the primary motivation for the author to maintain and update the project regularly.

Follow the WeChat Official Account 希里安 (cilliantech) to get the latest project updates and tech sharing!


开源不易，如果您觉得 CiliKube 对您有帮助或启发，请不吝点亮 Star ⭐！这是对作者持续维护和更新的最大鼓励。

关注公众号 希里安，获取项目最新动态和技术分享！



## 📞 Contact | 联系方式

Email: cilliantech@gmail.com

Website: https://www.cillian.website

WeChat

![alt text](ui/src/docs/wechat400x400.png)

## 📜 License | 许可证

This project is open-sourced under the Apache 2.0 License.


本项目基于 Apache 2.0 License 开源。[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)

## 🌟 Star History
<a href="https://star-history.com/#ciliverse/cilikube&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date" />
 </picture>
</a>
