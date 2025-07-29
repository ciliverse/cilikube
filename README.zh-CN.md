<div align="center">
  <img alt="CiliKube Logo" width="500" height="100" src="docs/logo.png">
  <h1>CiliKube</h1>
  <span><a href="./README.md">English</a> | 中文</span>
</div>

<div align="center">
  <img src="https://img.shields.io/badge/Frontend-Vue3-blue?style=flat-square&logo=vue.js" alt="Vue3">
  <img src="https://img.shields.io/badge/Frontend-TypeScript-blue?style=flat-square&logo=typescript" alt="TypeScript">
  <img src="https://img.shields.io/badge/Frontend-Vite-blue?style=flat-square&logo=vite" alt="Vite">
  <img src="https://img.shields.io/badge/Frontend-Element%20Plus-blue?style=flat-square&logo=element-plus" alt="Element Plus">
  <img src="https://img.shields.io/badge/Backend-Go-blue?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Backend-Gin-blue?style=flat-square&logo=gin" alt="Gin">
  <img src="https://img.shields.io/badge/1.33.0-Kubernetes-blue?style=flat-square&logo=kubernetes" alt="Kubernetes">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue?style=flat-square" alt="License: Apache 2.0">
  <img src="https://img.shields.io/github/stars/ciliverse/cilikube?style=social" alt="GitHub Stars">
  <img src="https://img.shields.io/github/forks/ciliverse/cilikube?style=social" alt="GitHub Forks">
</div>

## ❤️ 支持项目

开源不易，如果您觉得 CiliKube 对您有帮助或启发，请不吝点亮 Star ⭐！这是对作者持续维护和更新的最大鼓励。

关注微信公众号**希里安**，获取项目最新动态和技术分享！

## ❤️ 感谢所有贡献者

<a href="https://github.com/ciliverse/cilikube/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ciliverse/cilikube" />
</a>

感谢所有为 CiliKube 贡献代码和建议的开发者们！你们的努力让这个项目变得更好。

## 📝 感谢赞助

本项目 CDN 加速及安全防护由 Tencent EdgeOne 赞助

[![EdgeOne](https://edgeone.ai/media/34fe3a45-492d-4ea4-ae5d-ea1087ca7b4b.png)](https://edgeone.ai/zh?from=github)

## 🤔 CiliKube 是什么？

CiliKube 是一个使用现代主流技术栈（Vue3, TypeScript, Go, Gin）构建的开源全栈 Kubernetes 多集群资源管理平台。它致力于提供一个简洁、优雅的界面，来简化 Kubernetes 资源的日常管理（增删改查）并支持功能拓展。是入门学习k8s开发的不二之选。

![架构图](docs/architech.png)

## ✨ CiliKube 的特色

与追求"大而全"的复杂系统不同，CiliKube 专注于"小而美"。它的核心目标是：

1. **核心功能**: 提供清晰、直观的界面来管理常用的 K8s 资源
2. **学习友好**: 代码结构清晰，技术栈现代，非常适合作为学习 **Vue3/Go Web 开发** 和 **Kubernetes 二次开发** 的入门项目
3. **易于拓展**: 预留了自定义功能的空间，方便用户根据自身需求进行扩展

## 🎯 目标用户

- 希望学习 **Vue3 + TypeScript + ElementPlus** 前端开发的开发者
- 希望学习 **Go + Gin** 后端开发的开发者
- 对 **Kubernetes API** 和 **client-go** 使用感兴趣的云原生爱好者
- 需要一个简洁 K8s 管理面板，并可能进行二次开发的团队或个人

## 💡 项目背景

CiliKube 起源于作者学习 Web 全栈开发的实践项目。在学习过程中，作者深入探索了 Kubernetes，并获得了相关认证。这个项目不仅是学习成果的体现，更希望成为一把"钥匙 (Key)"，帮助更多像作者一样的学习者打开开源世界的大门，参与贡献，共同成长。

## 🌐 在线预览

- 在线演示: http://cilikubedemo.cillian.website
- 演示账号:
  - 用户名: admin
  - 密码: 12345678

## 📚 文档

- 官方文档: [cilikube.cillian.website](https://cilikube.cillian.website)

## 🚀 技术栈

本项目采用了当前流行的前后端技术栈，确保开发者能够接触和使用最新的工具和库。

**环境要求 (推荐)**:
- Node.js >= 18.0.0 (项目当前使用 v22.14.0 开发)
- Go >= 1.20 (项目当前使用 v1.24.2 开发)
- PNPM >= 8.x

**前端**: `Vue3` `TypeScript` `Vite` `Element Plus` `Pinia` `Vue Router` `Axios` `UnoCSS` `Scss` `ESLint` `Prettier`
- 基于优秀的 [v3-admin-vite](https://github.com/un-pany/v3-admin-vite) 模板进行开发，感谢原作者 un-pany。

**后端**: `Go` `Gin` `Kubernetes client-go` `JWT` `Gorilla Websocket` `Viper` `Zap Logger`

## ✨ 主要功能

- **用户认证**: 基于 JWT 的登录和权限验证
- **仪表盘**: 集群资源概览
- **多集群管理**: 支持管理多个 Kubernetes 集群
- **资源管理**:
  - 节点 (Node) 管理
  - 命名空间 (Namespace) 管理
  - Pod 管理 (列表、详情、日志、终端)
  - 存储卷 (PV/PVC) 管理
  - 配置 (ConfigMap/Secret) 管理
  - 网络 (Service/Ingress) 管理
  - 工作负载 (Deployment/StatefulSet/DaemonSet) 管理
- **系统设置**: 主题切换、国际化支持

## 🛠️ 开发计划

**前端**
- [x] 登录界面
- [x] 基础布局 (侧边栏, 顶部导航, 标签栏)
- [x] 消息通知
- [x] 工作负载资源页面 (Deployment, StatefulSet, DaemonSet 等)
- [x] 配置管理页面 (ConfigMap, Secret)
- [x] 网络资源页面 (Service, Ingress)
- [x] 存储资源页面 (StorageClass, PV, PVC)
- [x] 访问控制页面 (RBAC - ServiceAccount, Role, ClusterRoleBinding 等)
- [x] 日志查看页面优化
- [x] Web Shell 终端集成
- [ ] 事件 (Events) 查看
- [ ] CRD 资源管理 (基础)
- [ ] 监控集成 (集成 Prometheus/Grafana 数据展示)

**后端**
- [x] Kubernetes 客户端初始化
- [x] 基础路由设置 (Gin)
- [x] CORS 跨域配置
- [x] JWT 认证中间件
- [x] WebSocket 接口 (用于日志和 Web Shell)
- [x] 多集群支持
- [x] Node (节点) 资源接口
- [x] Pod 资源接口 (列表, 详情, 删除, 日志, Exec)
- [x] PV/PVC 资源接口
- [x] Namespace 资源接口
- [x] Deployment / StatefulSet / DaemonSet 资源接口
- [x] Service / Ingress 资源接口
- [x] ConfigMap / Secret 资源接口
- [x] RBAC 相关资源接口
- [x] Event 资源接口

## 💻 本地开发

### 环境准备
1. 安装 [Node.js](https://nodejs.org/) (>=18) 和 [pnpm](https://pnpm.io/)
2. 安装 [Go](https://go.dev/) (>=1.20)
3. 拥有一个 Kubernetes 集群，并配置好 kubeconfig 文件 (默认读取 `~/.kube/config`)

### 运行前端
```bash
# 进入前端目录
cd cilikube-web
# 安装依赖
pnpm install
# 启动开发服务器
pnpm dev
```

访问 http://localhost:8888 即可看到前端界面。

### 运行后端
```bash
# 进入后端目录
cd cilikube
# (可选) 更新 Go 依赖
go mod tidy
# 运行后端服务 (默认监听 8080 端口)
# 配置文件在 configs/config.yaml 中修改
go run cmd/server/main.go
```

### 构建项目
```bash
# 构建前端生产环境包 (输出到 cilikube-web/dist)
cd cilikube-web
pnpm build

# 构建后端可执行文件
cd ../cilikube
go build -o cilikube cmd/server/main.go
```

## 🐳 Docker 部署

### 使用官方镜像
```bash
# 后端
docker run -d --name cilikube -p 8080:8080 -v ~/.kube:/root/.kube:ro cilliantech/cilikube:latest

# 前端
docker run -d --name cilikube-web -p 80:80 cilliantech/cilikube-web:latest
```

### 使用 Docker Compose
```bash
docker-compose up -d
```

访问 http://localhost 即可。

## ☸️ Kubernetes 部署 (Helm)

### 环境准备
- 安装 Helm (>=3.0)
- 拥有一个 Kubernetes 集群，并配置好 kubeconfig 文件
- 安装 kubectl (>=1.20)

### 部署步骤
```bash
# 添加 Helm 仓库
helm repo add cilikube https://charts.cillian.website

# 更新 Helm 仓库
helm repo update

# 安装 CiliKube
helm install cilikube cilikube/cilikube -n cilikube --create-namespace

# 查看服务状态
kubectl get svc cilikube -n cilikube
```

## 🎨 功能预览

### 全新 Antd 组件界面 (即将上线!)
![新界面](docs/newui.png)
![Antd 2](docs/antd-2.png)
![节点 Ant](docs/node-ant.png)
![集群 Ant](docs/cluster-ant.png)

### 当前界面
![登录](docs/login.png)
![仪表盘](docs/dashboard.png)
![集群](docs/cluster.png)
![Pod](docs/pod.png)
![终端](docs/shell.png)

## 🤝 贡献指南

我们欢迎各种形式的贡献！如果您想参与改进 CiliKube，请：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'feat: Add some AmazingFeature'`) - 请遵循 Git 提交规范
4. 将您的分支推送到 Github (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

### Git 提交规范

请遵循 Conventional Commits 规范：

- `feat`: 新增功能
- `fix`: 修复 Bug
- `perf`: 性能优化
- `style`: 代码样式调整（不影响逻辑）
- `refactor`: 代码重构
- `revert`: 撤销更改
- `test`: 添加或修改测试
- `docs`: 文档或注释修改
- `chore`: 构建流程、依赖管理等杂项更改
- `workflow`: 工作流改进
- `ci`: CI/CD 配置相关
- `types`: 类型定义修改
- `wip`: 开发中的提交（不建议合入主分支）

## 📞 联系方式

- Email: cilliantech@gmail.com
- Website: https://www.cillian.website
- WeChat: 希里安

<img src="docs/wechat400x400.png" width="100" height="100" />

## 📜 许可证

本项目基于 Apache 2.0 License 开源

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)

## 🌟 Star History

<a href="https://star-history.com/#ciliverse/cilikube&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=ciliverse/cilikube&type=Date" />
 </picture>
</a>