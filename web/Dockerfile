# ---- Stage 1: BUILDER ----
FROM node:20-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git python3 make g++

# 设置工作目录
WORKDIR /build

# 复制包管理文件
COPY package.json pnpm-lock.yaml* ./

# 安装 pnpm 并安装依赖
RUN npm install -g pnpm@latest && \
    pnpm install --frozen-lockfile --prefer-offline

# 复制源代码
COPY . .

# 构建应用
RUN pnpm build

# ---- Stage 2: NGINX ----
FROM nginx:1.25-alpine AS runner

# 安装必要的工具
RUN apk add --no-cache curl

# 创建非 root 用户
RUN addgroup -g 1001 -S nginx && \
    adduser -u 1001 -S nginx -G nginx

# 复制构建产物
COPY --from=builder /build/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/nginx.conf

# 设置权限
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:80/ || exit 1

# 暴露端口
EXPOSE 80

# 切换到非 root 用户
USER nginx

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]
