<div align="center">
  <img alt="CiliKube Logo" width="500" height="100" src="docs/logo.png">
  <h1>CiliKube</h1>
  <span><a href="./README.zh-CN.md">中文纯净版</a> | English</span>
</div>

![Alt](https://repobeats.axiom.co/api/embed/97bc0de802d8faf0f90512177cb349c4e494d76d.svg "Repobeats analytics image")
<!-- [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fciliverse%2Fcilikube.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fciliverse%2Fcilikube?ref=badge_shield) -->
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)

## 最新公告（2025年6月）

- 目前单集群维护在release-0.1.x最新分支
- 多集群版本在0.2.x最新分支，正在开发中
- 全新Antd组件web界面即将上线，敬请期待！

## 🤔 What is CiliKube? | CiliKube 是什么？

CiliKube is an open-source, full-stack Kubernetes (k8s) multi-cluster resource management platform built with modern mainstream technology stack (Vue3, TypeScript, Go, Gin). It aims to provide a concise and elegant interface to simplify daily Kubernetes resource management (CRUD operations) while supporting feature extensions. It's an excellent choice for beginners learning k8s development.

CiliKube 是一个使用现代主流技术栈（Vue3, TypeScript, Go, Gin）构建的开源全栈 Kubernetes (k8s) 多集群资源管理平台。它致力于提供一个简洁、优雅的界面，来简化 Kubernetes 资源的日常管理（增删改查）并支持功能拓展。是入门学习k8s开发的不二之选。
![alt text](docs/architech.png)




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




## 🌐 Online Demo | 在线预览

*   Online Demo: http://cilikubedemo.cillian.website
*   在线演示站点: http://cilikubedemo.cillian.website
*   Demo Credentials:
    *   Username: admin
    *   Password: 12345678
*   演示账号:
    *   用户名: admin
    *   密码: 12345678

## 🎨 Feature Preview | 功能预览

### **即将上线全新Antd组件web界面，敬请期待！**

![alt text](docs/newui.png)
![alt text](docs/antd-2.png)

### **原始UI**
![alt text](docs/login.png)
![alt text](docs/first.png)
![alt text](docs/minikube2.png)
![alt text](docs/techstack.png)
![alt text](docs/monitor1.png)
![alt text](docs/monitor2.png)
![alt text](docs/dashboard.png)
![alt text](docs/nav.png)
![alt text](docs/cluster.png)
![alt text](docs/ingress.png)
![alt text](docs/namespace.png)
![alt text](docs/svc.png)
![alt text](docs/deployment.png)
![alt text](docs/pod.png)
![alt text](docs/shell.png)
![alt text](docs/pv.png)
![alt text](docs/pvc.png)
![alt text](docs/secret.png)
![alt text](docs/configmap.png)

## 📚 Documentation | 文档
![alt text](docs/docs.png)
*   Official Documentation: [cilikube.cillian.website](https://cilikube.cillian.website)
*   官方文档: [cilikube.cillian.website](https://cilikube.cillian.website) 

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
    *   Node.js >= 18.0.0 (项目当前使用 v22.15.0 开发)
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
*   [x] Workload Resource Pages (Deployment, StatefulSet, DaemonSet, etc.)
*   [x] Configuration Management Pages (ConfigMap, Secret)
*   [x] Network Resource Pages (Service, Ingress)
*   [x] Storage Resource Pages (StorageClass, PV, PVC) - *PV/PVC partially done*
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
*   [x] Namespace Resource API
*   [x] Deployment / StatefulSet / DaemonSet Resource API
*   [x] Service / Ingress Resource API
*   [x] ConfigMap / Secret Resource API
*   [ ] RBAC Related Resource API
*   [x] Event Resource API

**前端**

*   [x] 登录界面
*   [x] 基础布局 (侧边栏, 顶部导航, 标签栏)
*   [x] 消息通知
*   [x] 工作负载资源页面 (Deployment, StatefulSet, DaemonSet 等)
*   [x] 配置管理页面 (ConfigMap, Secret)
*   [x] 网络资源页面 (Service, Ingress)
*   [x] 存储资源页面 (StorageClass, PV, PVC) - *PV/PVC 部分完成*
*   [x] 访问控制页面 (RBAC - ServiceAccount, Role, ClusterRoleBinding 等)
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

CiliKube offers multiple deployment methods, making it convenient for users who want to quickly experience its interface or perform secondary development.

## 1. Local Development
**Tip**: If you want to perform secondary development or quickly experience CiliKube's front-end and back-end features, local development is a great choice. CiliKube's front-end and back-end code can run locally, ideal for development and debugging. Below are the steps for local setup:

### Environment Preparation
- Install Node.js (>=18) and pnpm
- Install Go (>=1.20)
- Have a Kubernetes cluster and configure the kubeconfig file (defaults to reading `~/.kube/config`)

### Running the Front-End
```bash
# Navigate to the front-end directory
cd cilikube-web
# Install dependencies
pnpm install
# Start the development server
pnpm dev
````

Visit http://localhost:8888 to see the front-end interface. You can modify the port and service configuration in the front-end code to connect to different back-end services.

### Running the Back-End

```bash
# Navigate to the back-end directory
cd cilikube
# (Optional) Update Go dependencies
go mod tidy
# Run the back-end service (listens on port 8080 by default)
# Configuration files are modified in configs/config.yaml
go run cmd/server/main.go
```

## 2\. Docker Deployment

**Tip**: For a quick experience of CiliKube's front-end and back-end features, Docker deployment is recommended. Both front-end and back-end can run via Docker images, perfect for quick setup and testing. Below are the steps for Docker deployment:

### Environment Preparation

  - Install Docker (\>=20.10)
  - Install Docker Compose (\>=1.29)
  - Have a Kubernetes cluster and configure the kubeconfig file (defaults to reading \~/.kube/config)

**Note**: When using official images, the latest stable version on Docker Hub is v0.1.0. Features in v0.1.1 can be experienced by building from source locally, with updated images coming soon.

  - Back-end: `cilliantech/cilikube:v0.1.0`
  - Front-end: `cilliantech/cilikube-web:v0.1.0`

<!-- end list -->

```bash
# Assuming the host's kubeconfig is at ~/.kube/config, the container expects it at /root/.kube/config
docker run -d --name cilikube -p 8080:8080 -v ~/.kube:/root/.kube:ro cilliantech/cilikube:v0.1.0
docker run -d --name cilikube-web -p 80:80 cilliantech/cilikube-web:v0.1.0
```

Alternatively, use Docker Compose to run:
A sample `docker-compose.yml` file can be found in the project's GitHub repository root directory.

```bash
docker-compose up -d
```

Visit http://localhost to access the interface.

**Note**: The above commands run CiliKube's front-end on port 80 and back-end on port 8080 locally. You can modify port mappings as needed.

### Building Custom Docker Images

You can also pull the code, modify the Dockerfile, and build custom images.
**Note**: The following steps assume you've cloned the front-end and back-end projects and are operating in their respective root directories.

#### 1\. Obtain the Code

Clone the CiliKube front-end and back-end repositories and navigate to their root directories.

```bash
cd path/to/cilikube
cd path/to/cilikube-web
```

#### 2\. Build Docker Images

After modifying the Dockerfile, build the Docker images for the front-end and back-end.

```bash
docker build -t "cilikube-server:latest" .
docker build -t "cilikube-web:latest" .
```

#### 3\. Run Docker Containers

After building the images, run the containers.

```bash
docker run --name cilikube-server -p 8080:8080 -d cilikube-server:latest
docker run --name cilikube-web -p 80:80 -d cilikube-web:latest
```

Now, open your browser and visit `http://<your-host-ip>` or `http://localhost` (if running Docker locally) to see the CiliKube login interface\!

## 3\. Kubernetes Deployment (Helm)

**Tip**: To run CiliKube in a Kubernetes cluster, Helm deployment is recommended. Both front-end and back-end can be deployed via Helm Charts, suitable for production environments and large-scale clusters. Below are the steps for Helm deployment:

### Environment Preparation

  - Install Helm (\>=3.0)
  - Have a Kubernetes cluster and configure the kubeconfig file (defaults to reading \~/.kube/config)
  - Install kubectl (\>=1.20)

### Deployment Steps

#### 1\. Add Helm Repository

Add the CiliKube Helm repository.

```bash
helm repo add cilikube [https://charts.cillian.website](https://charts.cillian.website)
```

#### 2\. Update Helm Repository

Update local Helm repository information.

```bash
helm repo update
```

#### 3\. Install CiliKube

Use Helm to install CiliKube.

```bash
helm install cilikube cilikube/cilikube -n cilikube --create-namespace
```

#### 4\. Access CiliKube

After installation, check the CiliKube service details.

```bash
kubectl get svc cilikube -n cilikube
```

Example output:

```bash
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
cilikube     NodePort    <your-host-ip>:<NodePort>   8080:80/TCP   1m
```

This indicates CiliKube is deployed and accessible via NodePort. Visit `http://<your-host-ip>:<NodePort>` in your browser to access the front-end interface.

#### 5\. Uninstall CiliKube

To remove CiliKube, use:

```bash
helm uninstall cilikube -n cilikube
```

**Note**: This command removes all CiliKube-related resources.

## 4\. Additional Notes

We strongly recommend consulting the official CiliKube documentation for the most comprehensive and up-to-date deployment guides:

  - Official Documentation: [cilikube.cillian.website](https://www.google.com/search?q=cilikube.cillian.website)
  - GitHub Repository: [github.com/ciliverse/cilikube](github.com/ciliverse/cilikube)

The official documentation provides detailed information on back-end deployment, database configuration (if required), and potential future deployment options.


---
CILIKUBE有多重部署方式，这对于想要快速体验其界面或进行二次开发的同学来说非常方便。

####  1.本地开发

Tips:  如果你想进行二次开发，或者想要快速体验 CiliKube 的前后端功能，可以选择本地开发的方式。CiliKube 的前后端代码都可以在本地运行，适合开发和调试。以下是本地运行的步骤：

**环境准备**

1. 安装 Node.js (>=18) 和 pnpm
2. 安装 Go (>=1.20)
3. 拥有一个 Kubernetes 集群，并配置好 kubeconfig 文件 (默认读取 ~/.kube/config)

**运行前端**
```bash
# 进入前端目录
cd cilikube-web
# 安装依赖
pnpm install
# 启动开发服务器
pnpm dev
```
- 访问 http://localhost:8888 即可看到前端界面，你可以在前端代码中修改端口以及服务配置，来连接不同的后端服务。

**运行后端**
```bash
# 进入后端目录
cd cilikube
# (可选) 更新 Go 依赖
go mod tidy

# 进入后端目录
cd cmd/server

# (可选) 更新 Go 依赖
# go mod tidy

# 运行后端服务 (默认监听 8080 端口)
# 配置文件在 configs/config.yaml 中修改
go run main.go
```


#### 2.Docker部署运行

Tips: 如果你想快速体验 CiliKube 的前后端功能，可以选择 Docker 部署的方式。CiliKube 的前后端都可以通过 Docker 镜像来运行，适合快速上手和测试。以下是 Docker 部署的步骤：

**环境准备**
1. 安装 Docker (>=20.10)
2. 安装 Docker Compose (>=1.29)
3. 拥有一个 Kubernetes 集群，并配置好 kubeconfig 文件 (默认读取 ~/.kube/config)


**注意：使用官方镜像，目前 Docker Hub 最新稳定版为 v0.1.0，v0.1.1 的特性可通过本地源码构建体验，镜像即将更新**
- 后端： cilliantech/cilikube:v0.1.0
- 前端：cilliantech/cilikube-web:v0.1.0

```bash
# 假设宿主机的 kubeconfig 在 ~/.kube/config, 容器内应用期望在 /root/.kube/config 读取
docker run -d --name cilikube -p 8080:8080 -v ~/.kube:/root/.kube:ro cilliantech/cilikube:v0.1.0
docker run -d --name cilikube-web -p 80:80 cilliantech/cilikube-web:v0.1.0
```

- 也可以使用docker-compose来运行

你可以在项目的 GitHub 仓库根目录找到 docker-compose.yml 文件示例

```bash
docker-compose up -d
```
- 访问 http://localhost即可



- **注意：** 以上命令会在本地80端口运行CiliKube的前端应用，8080端口运行后端应用。你可以根据需要修改端口映射。


也可以自己拉取代码后更改dockerfile内容后打镜像

**注意： 以下操作在clone前后端项目后分别在其根目录下操作**


1.  **获取代码：**
    首先，你需要获取CiliKube前后端的代码。通常可以通过`git clone`项目仓库，然后进入根目录。
    
    ```bash
    cd path/to/cilikube
    cd path/to/cilikube-web
    ```

2.  **构建Docker镜像：**
    在根目录下，修改镜像内容后执行以下命令来构建前后端应用的Docker镜像。

    ```bash
    docker build -t "cilikube-server:latest" .
    docker build -t "cilikube-web:latest" .
    ```

3.  **运行Docker容器：**
    镜像构建成功后，使用以下命令来运行前后端容器。

    ```bash
    docker run --name cilikube-server -p 8080:8080 -d cilikube-server:latest
    docker run --name cilikube-web -p 80:80 -d cilikube-web:latest
    ```

现在，打开你的浏览器，访问 `http://<你的主机IP>` 或 `http://localhost` (如果你在本机运行Docker)，就应该能看到CiliKube的登录界面了！


#### 3. Kubernetes 部署 (Helm部署)
Tips: 如果你想在 Kubernetes 集群中运行 CiliKube，可以选择 Helm 部署的方式。CiliKube 的前后端都可以通过 Helm Chart 来部署，适合生产环境和大规模集群。以下是 Helm 部署的步骤：
**环境准备**
1. 安装 Helm (>=3.0)
2. 拥有一个 Kubernetes 集群，并配置好 kubeconfig 文件 (默认读取 ~/.kube/config)
3. 安装 kubectl (>=1.20)


**部署步骤：**
1.  **添加 Helm 仓库：**
    首先，你需要添加 CiliKube 的 Helm 仓库。

    ```bash
    helm repo add cilikube https://charts.cillian.website
    ```
2.  **更新 Helm 仓库：**
    更新本地的 Helm 仓库信息。

    ```bash
    helm repo update
    ```
3.  **安装 CiliKube：**
    使用 Helm 安装 CiliKube。

    ```bash
    helm install cilikube cilikube/cilikube -n cilikube --create-namespace
    ```
4.  **访问 CiliKube：**
    安装完成后，你可以使用以下命令查看 CiliKube 的服务信息。

    ```bash
    kubectl get svc cilikube -n cilikube
    ```
    你会看到类似以下的输出：
    ```bash
    NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
    cilikube     NodePort
    你的主机IP: <NodePort>   8080:80/TCP   1m
    ```
    这表示 CiliKube 的服务已经成功部署，并且可以通过 NodePort 访问。你可以通过浏览器访问 `http://<你的主机IP>:<NodePort>` 来访问 CiliKube 的前端界面。
5.  **卸载 CiliKube：**
    如果你想卸载 CiliKube，可以使用以下命令。

    ```bash
    helm uninstall cilikube -n cilikube
    ```

**注意：** 以上命令会卸载CiliKube的所有相关资源。


#### 4. 备注

**我们强烈建议查阅CiliKube的官方文档以获取最全面和最新的部署指南：**

* **官方文档:** [cilikube.cillian.website](https://cilikube.cillian.website)
* **GitHub仓库:** [github.com/ciliverse/cilikube](https://github.com/ciliverse/cilikube)

在官方文档中，你可能会找到关于后端部署、数据库配置（如果需要）、以及未来可能支持的更多部署选项的详细信息。


## 🤝 Contribution Guide | 贡献指南

We welcome contributions of all forms! If you'd like to help improve CiliKube, please:

Fork this repository

Create your feature branch (git checkout -b feature/AmazingFeature)

Commit your changes (git commit -m 'feat: Add some AmazingFeature') - Please follow the Git Commit Guidelines

Push your branch to your fork (git push origin feature/AmazingFeature)

Submit a Pull Request


我们欢迎各种形式的贡献！如果您想参与改进 CiliKube，请：

### 注意
- 目前单集群维护在release-0.1.0分支，后续会有多集群版本的发布
- 多集群版本在main分支，正在开发中
- The single-cluster version is maintained in the release-0.1.0 branch, and a multi-cluster version will be released later.
- The main branch is for the multi-cluster version, which is under development.

- Fork 本仓库

- 创建您的特性分支 (git checkout -b feature/AmazingFeature)

- 提交您的更改 (git commit -m 'feat: Add some AmazingFeature') - 请遵循 Git 提交规范

- 将您的分支推送到 Github (git push origin feature/AmazingFeature)

- 提交 Pull Request

请确保您的代码经过充分测试，并且符合项目的代码规范。

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

## ❤️ Thanks To All The Contributors | 感谢所有贡献者
<a href="https://github.com/ciliverse/cilikube/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ciliverse/cilikube" />
</a>

感谢所有为 CiliKube 贡献代码和建议的开发者们！你们的努力让这个项目变得更好。


## 📞 Contact | 联系方式

Email: cilliantech@gmail.com

Website: https://www.cillian.website

WeChat

![alt text](docs/wechat400x400.png)

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
