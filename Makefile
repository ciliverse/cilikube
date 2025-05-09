.PHONY: help build-backend run-backend clean-backend build-ui dev-ui clean-ui all clean run docker-build docker-run docker-compose-dev docker-compose-prod docker-compose-down

# 颜色定义
YELLOW=\033[1;33m
GREEN=\033[1;32m
RED=\033[1;31m
BLUE=\033[1;34m
NC=\033[0m # No Color

# 代理配置
GO_PROXY=https://goproxy.cn,direct
NPM_REGISTRY=https://registry.npmmirror.com

# API配置，可在运行时通过环境变量覆盖
# 例如: API_HOST=172.18.100.194 API_PORT=8080 make docker-compose-dev

# 帮助信息
help:
	@echo "${YELLOW}CiliKube 构建辅助工具${NC}"
	@echo "${GREEN}后端命令:${NC}"
	@echo "  make build-backend    - 编译Go后端服务"
	@echo "  make run-backend      - 运行Go后端服务"
	@echo "  make clean-backend    - 清理后端构建产物"
	@echo ""
	@echo "${GREEN}前端命令:${NC}"
	@echo "  make build-ui         - 构建前端项目"
	@echo "  make dev-ui           - 开发模式运行前端"
	@echo "  make clean-ui         - 清理前端构建产物"
	@echo ""
	@echo "${GREEN}Docker命令:${NC}"
	@echo "  make docker-build     - 构建Docker镜像"
	@echo "  make docker-run       - 运行Docker镜像"
	@echo "  make docker-compose-dev  - 使用docker-compose启动开发环境"
	@echo "    可选参数: API_HOST=<ip地址> API_PORT=<端口>"
	@echo "    例如: API_HOST=172.18.100.194 API_PORT=8080 make docker-compose-dev"
	@echo "  make docker-compose-prod - 使用docker-compose启动生产环境"
	@echo "    可选参数: API_HOST=<ip地址> API_PORT=<端口>"
	@echo "  make docker-compose-down - 停止并移除所有容器"
	@echo ""
	@echo "${GREEN}通用命令:${NC}"
	@echo "  make all              - 构建前后端"
	@echo "  make clean            - 清理所有构建产物"

# ===== 后端部分 =====
build-backend:
	@echo "${BLUE}编译后端服务...${NC}"
	@echo "${BLUE}更新Go依赖...${NC}"
	@GOPROXY=$(GO_PROXY) go mod tidy
	@cd cmd/server && GOPROXY=$(GO_PROXY) go build -o ../../bin/cilikube main.go
	@echo "${GREEN}后端构建完成: bin/cilikube${NC}"

run-backend:
	@echo "${BLUE}运行后端服务...${NC}"
	@echo "${BLUE}更新Go依赖...${NC}"
	@GOPROXY=$(GO_PROXY) go mod tidy
	@GOPROXY=$(GO_PROXY) go run cmd/server/main.go

clean-backend:
	@echo "${BLUE}清理后端构建产物...${NC}"
	@rm -rf bin/
	@echo "${GREEN}后端清理完成${NC}"

# ===== 前端部分 =====
build-ui:
	@echo "${BLUE}构建前端应用...${NC}"
	@cd ui && npm config set registry $(NPM_REGISTRY) && yarn install && yarn build
	@echo "${GREEN}前端构建完成: ui/dist/${NC}"

dev-ui:
	@echo "${BLUE}开发模式运行前端...${NC}"
	@cd ui && npm config set registry $(NPM_REGISTRY) && yarn install && yarn dev

clean-ui:
	@echo "${BLUE}清理前端构建产物...${NC}"
	@rm -rf ui/dist
	@rm -rf ui/node_modules
	@echo "${GREEN}前端清理完成${NC}"

# ===== 通用部分 =====
all: build-backend build-ui
	@echo "${GREEN}全部构建完成${NC}"

clean: clean-backend clean-ui
	@echo "${GREEN}全部清理完成${NC}"