# ================== 可配置参数 ==================
SCAN_DIR       ?= ./infrastructure,./pkg/sdk,./app
OUTPUT_FILE    ?= auto_import.go
MODULE_NAME    ?=
CONCURRENCY    ?= 8
IGNORE_PATTERNS ?= **/vendor/**,**/.git/**,**/*_test.go,**/testdata/**,**/node_modules/**,**/middleware/**

# 构建相关
version        ?= v0.0.1
GitCommit      ?= hhhhhh
COMPACT_TS     ?= $(shell date +"%Y%m%d%H%M%S")
BINARY_NAME    ?= xianFuBlogGo
BUILD_DIR      ?= bin
LDFLAGS        ?= -s -w -extldflags "-static" -X "main.UserName=" \
                  -X "github.com/hongmengzhu/xianfu-blog-go/cmd.BuildVersion=${version}" \
                  -X "github.com/hongmengzhu/xianfu-blog-go/cmd.BuildGitCommit=${GitCommit}" \
                  -X "github.com/hongmengzhu/xianfu-blog-go/cmd.BuildTime=${COMPACT_TS}"
GO_BUILD_CMD   = CGO_ENABLED=0 go build -v -a \
                 -ldflags '$(LDFLAGS)' \
                 -gcflags="all=-trimpath=$(CURDIR)" \
                 -asmflags="all=-trimpath=$(CURDIR)" \
                 -trimpath

# 自动检测当前系统环境
# ================== 自动检测系统（使用 uname，无需 go） ==================
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Linux)
    CURRENT_OS := linux
else ifeq ($(UNAME_S),Darwin)
    CURRENT_OS := darwin
else
    CURRENT_OS := unknown
endif

ifeq ($(UNAME_M),x86_64)
    CURRENT_ARCH := amd64
else ifeq ($(UNAME_M),arm64)
    CURRENT_ARCH := arm64
else ifeq ($(UNAME_M),aarch64)
    CURRENT_ARCH := arm64
else
    CURRENT_ARCH := unknown
endif

# ================== 设置默认目标 ==================
.DEFAULT_GOAL := run

# ================== 目标定义 ==================
.PHONY: all generate run build build-all clean index-coding upgrade

# 完整流程（生成 + 运行），不设为默认
all: generate run

# 生成导入文件
generate:
	@echo "===== 运行扫描工具 ====="
	go run ./pkg/tools/auto_import/auto-import.go \
		-module=$(MODULE_NAME) \
		-dir=$(SCAN_DIR) \
		-output=$(OUTPUT_FILE) \
		-concurrency=$(CONCURRENCY) \
		-ignore='$(IGNORE_PATTERNS)'

# 直接运行（默认目标）
run: generate
	@echo "===== 启动项目 ====="
	go env -w GOEXPERIMENT=jsonv2
	go run main.go

# 智能编译当前平台
build:
	@echo "===== 检测到当前系统: $(CURRENT_OS)/$(CURRENT_ARCH) ====="
	@echo "===== 编译当前平台可执行文件 ====="
	$(GO_BUILD_CMD) -o $(BINARY_NAME) .

# 交叉编译所有平台
build-all: build-linux build-mac-intel build-mac-arm

build-linux:
	@mkdir -p $(BUILD_DIR)
	@echo "===== 构建 Linux (amd64) ====="
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

build-mac-intel:
	@mkdir -p $(BUILD_DIR)
	@echo "===== 构建 macOS Intel (amd64) ====="
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

build-mac-arm:
	@mkdir -p $(BUILD_DIR)
	@echo "===== 构建 macOS Apple Silicon (arm64) ====="
	GOOS=darwin GOARCH=arm64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

# 清理
clean:
	@echo "===== 清理构建产物 ====="
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

# ================== 索引初始化 ==================

# 索引初始化
index-coding:
	@echo "===== 初始化：CodeGraph ====="
	# 初始化项目配置，生成 .codegraph/ 存储目录
	codegraph init
	# 全项目索引（解析AST、构建符号图）
	codegraph index full
    # 增量索引（开发常用，只扫描变更文件）
    #codegraph index incremental
    # 验证索引
    #codegraph stats
    #codegraph validate
	#@echo "===== 初始化：Graphify ====="
    # 扫描整个仓库构建代码图索引
	#graphify extract .
    # 启动图服务供智能体调用
    #graphify serve
    # 加载已构建代码图，开启上下文检索服务
    #codegraph-context start --index-path .codegraph/index
    # 绑定已有代码图，分析变更影响范围
    #code-review-graph bind --graph-path .codegraph/index
    # 针对当前Git变更做影响分析
    #code-review-graph analyze diff
	@echo "===== 索引初始化完成 ====="

upgrade:
	@echo "===== 升级依赖 ====="
	go get -u ./...