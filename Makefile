# Luna项目Makefile

# 变量定义
GO=go
GOFLAGS=-ldflags="-s -w"
BINARY_NAME=luna
BINARY_DIR=bin
PLUGIN_DIR=plugins
PLUGIN_BUILD_SCRIPT=examples/build_plugin.sh

# 确保目标目录存在
$(shell mkdir -p $(BINARY_DIR))

# 默认目标
.PHONY: all
all: build

# 构建主程序
.PHONY: build
build:
	@echo "正在编译 Luna..."
	$(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/lua
	@echo "编译完成: $(BINARY_DIR)/$(BINARY_NAME)"

# 安装依赖
.PHONY: deps
deps:
	@echo "正在安装依赖..."
	$(GO) mod download
	@echo "依赖安装完成"

# 运行测试
.PHONY: test
test:
	@echo "正在运行测试..."
	$(GO) test -v ./...

# 构建插件
.PHONY: plugins
plugins:
	@echo "正在编译插件..."
	if [ -f $(PLUGIN_BUILD_SCRIPT) ]; then \
		sh $(PLUGIN_BUILD_SCRIPT); \
	else \
		echo "警告: 插件构建脚本不存在 ($(PLUGIN_BUILD_SCRIPT))"; \
	fi

# 清理生成的文件
.PHONY: clean
clean:
	@echo "正在清理..."
	rm -rf $(BINARY_DIR)
	@echo "清理完成"

# 帮助信息
.PHONY: help
help:
	@echo "Luna 构建系统"
	@echo ""
	@echo "可用目标:"
	@echo "  all      - 默认目标，构建主程序"
	@echo "  build    - 仅构建主程序"
	@echo "  plugins  - 构建插件"
	@echo "  deps     - 安装依赖"
	@echo "  test     - 运行测试"
	@echo "  clean    - 清理生成的文件"
	@echo "  help     - 显示此帮助信息"

# 完整构建（主程序+插件）
.PHONY: full
full: deps build plugins
	@echo "完整构建完成"