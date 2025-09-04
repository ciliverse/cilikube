OUT_DIR := output
BINARY_NAME := cilikube
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -w -s"

.PHONY: build run build-linux build-mac build-windows build-all test lint clean dev docker help

# 默认目标
all: build

# 更新依赖
update-dependencies:
	@echo "Updating Go dependencies..."
	go mod tidy
	go mod download

# 开发环境构建
build: clean update-dependencies
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME) cmd/server/main.go

# 开发环境运行
dev: build
	@echo "Starting development server..."
	./$(OUT_DIR)/$(BINARY_NAME) --config configs/config.yaml

# 生产环境运行
run: build
	@echo "Starting production server..."
	./$(OUT_DIR)/$(BINARY_NAME)

# 清理构建文件
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(OUT_DIR)
	go clean -cache

# 运行测试
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 运行基准测试
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# 代码检查
lint:
	@echo "Running linters..."
	golangci-lint run ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# 安全检查
security:
	@echo "Running security checks..."
	gosec ./...

# Linux 构建
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-linux-amd64 cmd/server/main.go

# macOS 构建
build-mac:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-darwin-amd64 cmd/server/main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-darwin-arm64 cmd/server/main.go

# Windows 构建
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-windows-amd64.exe cmd/server/main.go

# 全平台构建
build-all: build-linux build-mac build-windows
	@echo "All builds completed!"

# Docker 构建
docker:
	@echo "Building Docker image..."
	docker build -t cilikube:$(VERSION) .
	docker build -t cilikube:latest .

# Docker 运行
docker-run:
	@echo "Running Docker container..."
	docker run -d --name cilikube -p 8080:8080 -v ~/.kube:/root/.kube:ro cilikube:latest

# 安装开发工具
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# 生成 API 文档
docs:
	@echo "Generating API documentation..."
	swag init -g cmd/server/main.go -o docs/swagger

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  dev            - Build and run in development mode"
	@echo "  run            - Build and run in production mode"
	@echo "  test           - Run tests with coverage"
	@echo "  bench          - Run benchmarks"
	@echo "  lint           - Run code linters"
	@echo "  fmt            - Format code"
	@echo "  security       - Run security checks"
	@echo "  clean          - Clean build artifacts"
	@echo "  build-all      - Build for all platforms"
	@echo "  docker         - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  install-tools  - Install development tools"
	@echo "  docs           - Generate API documentation"
	@echo "  help           - Show this help message"