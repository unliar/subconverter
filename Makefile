.PHONY: build test clean run lint fmt vet help docker docker-build

# 变量定义
APP_NAME := subconverter-go
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go 相关变量
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet

# 构建标志
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -s -w"
BUILD_FLAGS := -trimpath $(LDFLAGS)

# 目标文件
BINARY_NAME := $(APP_NAME)
BINARY_DIR := bin
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)

help: ## 显示帮助信息
	@echo "可用的命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 构建应用程序
	@echo "构建 $(APP_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_PATH) ./cmd/subconverter

build-linux: ## 构建 Linux 版本
	@echo "构建 Linux 版本..."
	@mkdir -p $(BINARY_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/subconverter

build-windows: ## 构建 Windows 版本
	@echo "构建 Windows 版本..."
	@mkdir -p $(BINARY_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/subconverter

build-darwin: ## 构建 macOS 版本
	@echo "构建 macOS 版本..."
	@mkdir -p $(BINARY_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/subconverter

build-all: build-linux build-windows build-darwin ## 构建所有平台版本

run: ## 运行应用程序
	$(GOCMD) run ./cmd/subconverter

test: ## 运行测试
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## 运行测试并显示覆盖率
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "测试覆盖率报告已生成：coverage.html"

bench: ## 运行基准测试
	$(GOTEST) -bench=. -benchmem ./...

lint: ## 运行代码检查
	@which golangci-lint > /dev/null || (echo "请先安装 golangci-lint" && exit 1)
	golangci-lint run

fmt: ## 格式化代码
	$(GOFMT) -s -w .

vet: ## 运行 go vet
	$(GOVET) ./...

clean: ## 清理构建文件
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

deps: ## 下载依赖
	$(GOMOD) download
	$(GOMOD) tidy

deps-update: ## 更新依赖
	$(GOGET) -u ./...
	$(GOMOD) tidy

docker-build: ## 构建 Docker 镜像
	docker build -t $(APP_NAME):$(VERSION) -f docker/Dockerfile .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

docker-run: ## 运行 Docker 容器
	docker run --rm -p 25500:25500 $(APP_NAME):latest

install: ## 安装到系统
	$(GOBUILD) $(BUILD_FLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) ./cmd/subconverter

# 开发相关命令
dev: ## 启动开发模式（需要安装 air）
	@which air > /dev/null || (echo "请先安装 air: go install github.com/cosmtrek/air@latest" && exit 1)
	air

setup: ## 设置开发环境
	@echo "设置开发环境..."
	$(GOMOD) download
	@echo "安装开发工具..."
	$(GOGET) github.com/cosmtrek/air@latest
	@echo "设置完成！运行 'make dev' 启动开发模式"

# 版本信息
version: ## 显示版本信息
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

.DEFAULT_GOAL := help