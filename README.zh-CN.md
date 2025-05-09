<div align="center">
  <img alt="CiliKube Logo" width="500" height="100" src="ui/public/logo.png">
  <h1>CiliKube</h1>
  <span>中文 | <a href="./README.md">English</a></span>
</div>

[![Alt](https://repobeats.axiom.co/api/embed/4b23db6e62b6a072c36e2d37235d49c9bf08af5b.svg "Repobeats analytics image")](https://github.com/ciliverse/cilikube)
<!-- [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fciliverse%2Fcilikube.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fciliverse%2Fcilikube?ref=badge_shield) -->


## 🤔 CiliKube 是什么？

CiliKube 是一个使用现代主流技术栈（Vue3, TypeScript, Go, Gin）构建的开源全栈 Kubernetes (k8s) 资源管理平台。它致力于提供一个简洁、优雅的界面，来简化 Kubernetes 资源的日常管理（增删改查）并支持功能拓展。
![alt text](ui/src/docs/architech.png)

## ✨ CiliKube 的特色

与追求“大而全”的复杂系统不同，CiliKube 专注于“小而美”。它的核心目标是：

1.  **核心功能**: 提供清晰、直观的界面来管理常用的 K8s 资源。
2.  **学习友好**: 代码结构清晰，技术栈现代，非常适合作为学习 **Vue3/Go Web 开发** 和 **Kubernetes 二次开发** 的入门项目。
3.  **易于拓展**: 预留了自定义功能的空间，方便用户根据自身需求进行扩展。

## 🎯 目标用户

*   希望学习 **Vue3 + TypeScript + ElementPlus** 前端开发的开发者。
*   希望学习 **Go + Gin** 后端开发的开发者。
*   对 **Kubernetes API** 和 **client-go** 使用感兴趣的云原生爱好者。
*   需要一个简洁 K8s 管理面板，并可能进行二次开发的团队或个人。

## 💡 项目背景

CiliKube 起源于作者学习 Web 全栈开发的实践项目。在学习过程中，作者深入探索了 Kubernetes，并获得了相关认证。这个项目不仅是学习成果的体现，更希望成为一把“钥匙 (Key)”，帮助更多像作者一样的学习者打开开源世界的大门，参与贡献，共同成长。

## 📚 文档

*   官方文档: [cilikube.cillian.website](https://cilikube.cillian.website) 

## 🌐 在线预览

*   在线演示站点正在部署中，敬请期待！

## 🚀 技术栈

本项目采用了当前流行的前后端技术栈，确保开发者能够接触和使用最新的工具和库。

*   **环境要求 (推荐)**:
    *   Node.js >= 18.0.0 (项目当前使用 v22.14.0开发)
    *   Go >= 1.20 (项目当前使用 v1.24.2开发)
    *   PNPM >= 8.x

*   **前端**: `Vue3` `TypeScript` `Vite` `Element Plus` `Pinia` `Vue Router` `Axios` `UnoCSS` `Scss` `ESlint` `Prettier`
    *   基于优秀的 [v3-admin-vite](https://github.com/un-pany/v3-admin-vite) 模板进行开发，感谢原作者 un-pany。

*   **后端**: `Go` `Gin` `Kubernetes client-go` `JWT (dgrijalva/jwt-go)` `Gorilla Websocket` `Logger (wonderivan/logger)`

## ✨ 主要功能 ([查看详细开发计划](#️-开发计划-roadmap))

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

## 🛠️ 开发计划 (Roadmap)

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
*   [x] JWT 认证中间件
*   [x] WebSocket 接口 (用于日志和 Web Shell)
*   [x] Node (节点) 资源接口
*   [x] Pod 资源接口 (列表, 详情, 删除, 日志, Exec)
*   [x] PV/PVC 资源接口
*   [ ] Namespace 资源接口
*   [ ] Deployment / StatefulSet / DaemonSet 资源接口
*   [ ] Service / Ingress 资源接口
*   [ ] ConfigMap / Secret 资源接口
*   [ ] RBAC 相关资源接口
*   [ ] Event 资源接口

## 💻 本地开发

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

# (可选) 更新 Go 依赖
# go mod tidy

# 运行后端服务 (默认监听 8081 端口)
go run cmd/server/main.go
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
## 🎨 功能预览

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

## 🤝 贡献指南

我们欢迎各种形式的贡献！如果您想参与改进 CiliKube，请：

- Fork 本仓库

- 创建您的特性分支 (git checkout -b feature/AmazingFeature)

- 提交您的更改 (git commit -m 'feat: Add some AmazingFeature') - 请遵循 Git 提交规范

- 将您的分支推送到 Github (git push origin feature/AmazingFeature)

- 提交 Pull Request

## Git 提交规范

请遵循 Conventional Commits 规范：

- feat: 新增功能

- fix: 修复 Bug

- perf: 性能优化

- style: 代码样式调整（不影响逻辑）

- refactor: 代码重构

- revert: 撤销更改

- test: 添加或修改测试

- docs: 文档或注释修改

- chore: 构建流程、依赖管理等杂项更改

- workflow: 工作流改进

- ci: CI/CD 配置相关

- types: 类型定义修改

- wip: 开发中的提交（不建议合入主分支）

## ❤️ 支持项目

开源不易，如果您觉得 CiliKube 对您有帮助或启发，请不吝点亮 Star ⭐！这是对作者持续维护和更新的最大鼓励。

关注公众号 希里安，获取项目最新动态和技术分享！

(可以考虑放公众号二维码图片)

## 📞 联系方式

- Email: cilliantech@gmail.com

- Website: https://www.cillian.website

- WeChat

![alt text](ui/src/docs/wechat400x400.png)

## 📜 License

本项目基于 Apache 2.0 License 开源。[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)






## 🌟 Star History
<a href="https://star-history.com/#ciliverse/cilikube&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date" />
 </picture>
</a>

