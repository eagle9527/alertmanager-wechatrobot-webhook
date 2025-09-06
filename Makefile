# Makefile for alertmanager-wechatrobot-webhook

# 变量定义
APP_NAME = alertmanager-wechatbot-webhook
GO_VERSION = 1.19
DOCKER_IMAGE = swr.cn-east-3.myhuaweicloud.com/lunz-prometheus/alertmanager-wechatbot-webhook
DOCKER_TAG = latest

# 默认目标
.PHONY: all
all: build

# 编译
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o $(APP_NAME) .

# 清理
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -f $(APP_NAME)
	go clean

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

# 下载依赖
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# 构建Docker镜像
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build   --platform linux/amd64   -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# 运行Docker容器
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -d --name $(APP_NAME) -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# 停止Docker容器
.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# 安装
.PHONY: install
install: build
	@echo "Installing $(APP_NAME)..."
	cp $(APP_NAME) /usr/local/bin/

# 卸载
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	rm -f /usr/local/bin/$(APP_NAME)

# 运行
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

# 开发模式运行
.PHONY: dev
dev:
	@echo "Running in development mode..."
	go run .

# 检查代码质量
.PHONY: check
check: fmt vet test
	@echo "Code quality check completed"

# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - 编译应用程序"
	@echo "  clean        - 清理编译文件"
	@echo "  test         - 运行测试"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 代码检查"
	@echo "  deps         - 下载依赖"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-run   - 运行Docker容器"
	@echo "  docker-stop  - 停止Docker容器"
	@echo "  install      - 安装到系统"
	@echo "  uninstall    - 从系统卸载"
	@echo "  run          - 编译并运行"
	@echo "  dev          - 开发模式运行"
	@echo "  check        - 代码质量检查"
	@echo "  help         - 显示帮助信息"