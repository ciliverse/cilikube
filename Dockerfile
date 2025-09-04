# ---- Stage 1: BUILDER ----
FROM golang:1.24-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置 Go 代理以加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 获取构建信息
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# 编译应用
RUN go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o cilikube \
    ./cmd/server/main.go


# ---- Stage 2: RUNNER ----
FROM alpine:3.19 AS runner

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和配置
COPY --from=builder /build/cilikube ./
COPY --from=builder /build/configs ./configs/

# 设置文件权限
RUN chown -R appuser:appgroup /app && \
    chmod +x /app/cilikube

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 暴露端口
EXPOSE 8080

# 切换到非 root 用户
USER appuser

# 设置入口点
ENTRYPOINT ["./cilikube"]
CMD ["--config", "configs/config.yaml"]
