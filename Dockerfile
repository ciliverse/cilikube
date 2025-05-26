# ---- Stage 1: BUILDER ----
# 使用官方的 Golang Alpine 镜像作为构建环境
# 你可以根据你的项目选择一个合适的 Go 版本，例如 1.22, 1.23, 1.24 等
FROM golang:1.24-alpine AS builder

# 可选: 设置 Go 代理以加速在中国大陆地区的依赖下载
# RUN go env -w GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /cilikube

# 复制 go.mod 和 go.sum 文件
# 优先复制它们可以利用 Docker 的层缓存机制，如果这些文件没有改变，则不需要重新下载依赖
COPY go.mod go.sum ./

# 下载 Go 模块依赖
RUN go mod download

# 可选: 验证依赖的完整性
# RUN go mod verify


# 复制项目其余的源代码到工作目录
COPY . .

# 编译 Go 应用
# CGO_ENABLED=0: 禁用 CGO，以便更容易地构建静态链接的二进制文件
# GOOS=linux: 指定目标操作系统为 Linux
# -ldflags="-w -s": 优化编译输出，-w 去除调试信息，-s 去除符号表，减小二进制文件体积
# -o cilikube: 指定输出的可执行文件名为 cilikube
# ./cmd/server/main.go: 你应用的入口点
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o cilikube ./cmd/server/main.go


# ---- Stage 2: RUNNER ----
# 使用官方的 Alpine Linux 镜像作为轻量级的运行环境
# alpine:3.19 是一个较新且稳定的版本 (请检查是否有更新的稳定版)
FROM alpine:3.19 AS runner
# 对于非常小的镜像，且你的应用是完全静态链接的，也可以考虑：
# FROM gcr.io/distroless/static-debian11 AS RUNNER
# FROM scratch AS RUNNER

# 设置工作目录
WORKDIR /cilikube

# 创建一个非 root 用户和组，以增强安全性
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 从 BUILDER 阶段复制编译好的可执行文件到当前阶段的 /cilikube/ 目录
# 从 BUILDER 阶段复制配置文件到当前阶段的 /cilikube/ 目录
# 这样，配置文件和可执行文件在同一个目录下
# (BUILDER 阶段通过 `COPY . .` 已经拥有了 configs/config.yaml，路径是 /app/configs/config.yaml)
COPY --from=BUILDER /cilikube/cilikube .
COPY --from=BUILDER /cilikube/configs/config.yaml .

# 如果你希望在最终镜像中保持 configs/config.yaml 的目录结构，可以使用:
# RUN mkdir ./configs
# COPY --from=BUILDER /app/configs/config.yaml ./configs/config.yaml
# 相应的 CMD 也需要修改为 ["--config", "./configs/config.yaml"]

# # 确保新创建的用户拥有对应用文件和配置文件的所有权
RUN chown appuser:appgroup cilikube config.yaml
# # 如果你保持了 ./configs/ 目录结构，则使用:
# # RUN chown -R appuser:appgroup cilikube configs

# 声明应用将监听的端口 (这主要是一个文档作用)
EXPOSE 8080


# 切换到非 root 用户运行应用
USER appuser

# 设置容器启动时执行的命令 (应用入口点)
ENTRYPOINT ["./cilikube"]

# 为入口点设置默认参数
# 假设你的应用通过 --config 参数指定配置文件，并且配置文件位于工作目录
CMD ["--config", "config.yaml"]
# 如果配置文件在 ./configs/config.yaml，则使用:
# CMD ["--config", "./configs/config.yaml"]
