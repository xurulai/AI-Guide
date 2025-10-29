.PHONY: build run test lint clean docker-build docker-up docker-down

# 变量定义
APP_NAME=shopping-guide-api
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(shell go version)

# 构建
build:
	@echo "Building $(APP_NAME)..."
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o bin/server cmd/server/main.go

# 运行
run:
	@echo "Running $(APP_NAME)..."
	go run cmd/server/main.go

# 测试
test:
	@echo "Running tests..."
	go test -v -cover ./...

# 测试覆盖率
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# 清理
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Docker构建
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) -f deploy/docker/Dockerfile .

# Docker启动
docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f deploy/docker/docker-compose.yml up -d

# Docker停止
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f deploy/docker/docker-compose.yml down

# 数据库迁移
migrate-up:
	@echo "Running database migrations..."
	mysql -u root -p < scripts/init_db.sql

# 生成mock
mock:
	@echo "Generating mocks..."
	mockgen -source=internal/service/chat_service.go -destination=internal/service/mock/chat_service_mock.go

# 帮助
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run code linter"
	@echo "  fmt            - Format code"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  migrate-up     - Run database migrations"
	@echo "  mock           - Generate mocks"

