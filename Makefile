# Synapse — Personal Knowledge Hub
# https://github.com/tunsuy/synapse

BINARY_NAME := synapse
MODULE := github.com/tunsuy/synapse
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := bin
MAIN_PKG := ./cmd/synapse

GOFLAGS := -trimpath
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: all build clean test test-verbose test-cover lint vet fmt run-init help

## help: 显示帮助信息
help:
	@echo "Synapse Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build        编译 synapse 二进制"
	@echo "  make test         运行所有测试"
	@echo "  make test-verbose 运行所有测试（详细输出）"
	@echo "  make test-cover   运行测试并生成覆盖率报告"
	@echo "  make lint         运行静态分析"
	@echo "  make fmt          格式化代码"
	@echo "  make clean        清理构建产物"
	@echo "  make run-init     构建并运行 synapse init"
	@echo ""

## all: 编译
all: build

## build: 编译 synapse 二进制
build:
	@echo "==> Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PKG)
	@echo "==> Done: $(BUILD_DIR)/$(BINARY_NAME)"

## test: 运行所有测试
test:
	go test ./... -count=1

## test-verbose: 运行所有测试（详细输出）
test-verbose:
	go test ./... -v -count=1

## test-cover: 运行测试并生成覆盖率报告
test-cover:
	@mkdir -p $(BUILD_DIR)
	go test ./... -coverprofile=$(BUILD_DIR)/coverage.out
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "==> Coverage report: $(BUILD_DIR)/coverage.html"
	@go tool cover -func=$(BUILD_DIR)/coverage.out | tail -1

## lint: 运行静态分析（go vet）
lint: vet

## vet: 运行 go vet
vet:
	go vet ./...

## fmt: 格式化代码
fmt:
	gofmt -s -w .
	goimports -w . 2>/dev/null || true

## clean: 清理构建产物
clean:
	@rm -rf $(BUILD_DIR)
	@echo "==> Cleaned"

## run-init: 构建并运行 synapse init（测试用，输出到临时目录）
run-init: build
	@tmpdir=$$(mktemp -d) && \
	echo "==> Running synapse init at $$tmpdir" && \
	$(BUILD_DIR)/$(BINARY_NAME) init --path $$tmpdir --name "Demo User" && \
	echo "" && \
	echo "==> Files created:" && \
	find $$tmpdir -type f | sort && \
	rm -rf $$tmpdir
